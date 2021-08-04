package search

import (
	"fmt"
	"math"

	. "github.com/amanjpro/zahak/book"
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func (e *Engine) Search(depth int8) {
	e.ClearForSearch()
	e.rootSearch(depth)
}

var p = WhitePawn.Weight()
var r = WhiteRook.Weight()
var b = WhiteBishop.Weight()
var q = WhiteQueen.Weight()
var WIN_IN_MAX = CHECKMATE_EVAL - int16(MAX_DEPTH)

var lmrReductions [32][32]int = initLMR()

// This idea is taken from Weiss, which I believe in turn is taken from many open source
// engines.
func initLMR() [32][32]int {
	var reductions [32][32]int
	for depth := 1; depth < 32; depth++ {
		for moves := 1; moves < 32; moves++ {
			reductions[depth][moves] = int(0.8 + math.Log(float64(depth))*math.Log(1.2*float64(moves))/2.5)
		}
	}
	return reductions
}

func (e *Engine) rootSearch(depth int8) {

	var previousBestMove Move

	e.move = EmptyMove
	e.score = -MAX_INT
	fruitelessIterations := 0

	bookmove := GetBookMove(e.Position)
	lastDepth := int8(1)

	if bookmove != EmptyMove {
		e.move = bookmove
		e.pv.Recycle()
		e.pv.AddFirst(bookmove)
	}

	if e.move == EmptyMove {
		e.pv.Recycle()
		for iterationDepth := int8(1); iterationDepth <= depth; iterationDepth++ {

			if iterationDepth > 1 && !e.TimeManager.CanStartNewIteration() {
				break
			}

			e.innerLines[0].Recycle()
			score := e.aspirationWindow(e.score, iterationDepth)

			if e.TimeManager.AbruptStop {
				break
			}

			if iterationDepth >= 8 && e.score-score >= 30 { // Position degrading
				e.TimeManager.ExtraTime()
			}

			e.pv.Clone(e.innerLines[0])
			e.score = score
			e.move = e.pv.MoveAt(0)
			e.SendPv(iterationDepth)
			lastDepth = iterationDepth
			if !e.Pondering && iterationDepth >= 35 && e.move == previousBestMove {
				fruitelessIterations++
				if fruitelessIterations > 4 {
					break
				}
			} else {
				fruitelessIterations = 0
			}
			if isCheckmateEval(e.score) {
				break
			}
			previousBestMove = e.move
			e.pred.Clear()
			if !e.Pondering && e.DebugMode {
				e.info.Print()
			}

		}

	}

	e.SendPv(lastDepth)
}

func (e *Engine) aspirationWindow(score int16, iterationDepth int8) int16 {
	e.doPruning = iterationDepth > 3
	if iterationDepth <= 6 {
		alpha := -MAX_INT
		beta := MAX_INT
		score = e.alphaBeta(iterationDepth, 0, alpha, beta)
	} else {
		var alpha, beta int16
		alphaMargin := int16(25)
		betaMargin := int16(25)
		for i := 0; i < 3; i++ {
			if i < 2 {
				alpha = max16(score-alphaMargin, -MAX_INT)
				beta = min16(score+betaMargin, MAX_INT)
			} else {
				alpha = -MAX_INT
				beta = MAX_INT
			}
			score = e.alphaBeta(iterationDepth, 0, alpha, beta)
			if score <= alpha {
				alphaMargin *= 2
			} else if score >= beta {
				betaMargin *= 2
			} else {
				return score
			}
		}
	}
	return score
}

func (e *Engine) alphaBeta(depthLeft int8, searchHeight int8, alpha int16, beta int16) int16 {
	e.VisitNode()

	isRootNode := searchHeight == 0
	isPvNode := alpha != beta-1

	position := e.Position

	currentMove := e.positionMoves[searchHeight]
	// Position is drawn
	if IsRepetition(position, e.pred, currentMove) || position.IsDraw() {
		return 0
	}

	if searchHeight >= MAX_DEPTH-1 {
		eval := Evaluate(position)
		e.staticEvals[searchHeight] = eval
		return eval
	}

	var isInCheck = position.IsInCheck()
	if isInCheck {
		e.info.checkExtentionCounter += 1
		depthLeft += 1 // Check Extension
	}

	if depthLeft <= 0 {
		e.staticEvals[searchHeight] = Evaluate(position)
		return e.quiescence(alpha, beta, searchHeight)
	}

	if isPvNode {
		e.info.mainSearchCounter += 1
	} else {
		e.info.zwCounter += 1
	}

	hash := position.Hash()
	nHashMove, nEval, nDepth, nType, ttHit := e.TranspositionTable.Get(hash)
	if !isPvNode && ttHit && nDepth >= depthLeft {
		if nEval >= beta && nType == LowerBound {
			e.CacheHit()
			return nEval
		}
		if nEval <= alpha && nType == UpperBound {
			e.CacheHit()
			return nEval
		}
		if nType == Exact {
			e.CacheHit()
			return nEval
		}
	}

	// Internal iterative reduction based on Rebel's idea
	if !isPvNode && !ttHit && depthLeft >= 3 {
		e.info.internalIterativeReduction += 1
		depthLeft -= 1
	}

	if !isRootNode {
		if e.TimeManager.ShouldStop(false, false) {
			return -MAX_INT
		}
	}

	var eval int16 = -MAX_INT
	if !isRootNode && currentMove == EmptyMove {
		eval = -1 * (e.staticEvals[searchHeight-1] + Tempo + Tempo)
	} else {
		eval = Evaluate(position)
	}

	e.staticEvals[searchHeight] = eval
	improving := currentMove == EmptyMove ||
		(searchHeight > 2 && e.staticEvals[searchHeight] > e.staticEvals[searchHeight-2])

	// Pruning
	pruningAllowed := !isPvNode && !isInCheck && e.doPruning

	if pruningAllowed {
		// Razoring
		if depthLeft < 3 && eval+b < beta {
			newEval := e.quiescence(alpha, beta, searchHeight)
			e.info.razoringCounter += 1
			return newEval
		}

		// Reverse Futility Pruning
		reverseFutilityMargin := int16(depthLeft) * p //(b - p)
		if improving {
			reverseFutilityMargin += p // int16(depthLeft) * p
		}
		if depthLeft < 7 && eval-reverseFutilityMargin >= beta {
			e.info.rfpCounter += 1
			return eval - reverseFutilityMargin /* fail soft */
		}

		// NullMove pruning
		isNullMoveAllowed := currentMove != EmptyMove && !position.IsEndGame()
		if isNullMoveAllowed && depthLeft >= 2 && eval > beta {
			var R = 4 + depthLeft/4
			if eval >= beta+50 {
				R = min8(R, depthLeft)
			} else {
				R = min8(R, depthLeft-1)
			}
			if R >= 2 {
				ep := position.MakeNullMove()
				e.pred.Push(position.Hash())
				e.innerLines[searchHeight+1].Recycle()
				e.positionMoves[searchHeight+1] = EmptyMove
				score := -e.alphaBeta(depthLeft-R, searchHeight+1, -beta, -beta+1)
				e.pred.Pop()
				position.UnMakeNullMove(ep)
				if score >= beta {
					e.info.nullMoveCounter += 1
					return score
				}
			}
		}

		// Prob cut
		// The idea is basically cherry picked from multiple engines, Weiss, Ethereal and Berserk for example
		probBeta := min16(beta+110, WIN_IN_MAX)
		if depthLeft > 4 && abs16(beta) < WIN_IN_MAX && !(ttHit && nDepth >= depthLeft-3 && nEval < probBeta) {

			hashMove := EmptyMove
			if hashMove.IsCapture() || hashMove.PromoType() != NoType {
				hashMove = nHashMove
			}
			movePicker := e.MovePickers[searchHeight]
			movePicker.RecycleWith(position, e, depthLeft, hashMove, true)
			seeScores := movePicker.captureMoveList.Scores
			i := 0
			for true {
				move := movePicker.Next()
				if move == EmptyMove || seeScores[i] < 0 {
					break
				}
				if oldEnPassant, oldTag, hc, ok := position.MakeMove(move); ok {
					var score int16
					if depthLeft >= 8 {
						e.innerLines[searchHeight+1].Recycle()
						e.pred.Push(position.Hash())
						e.positionMoves[searchHeight+1] = move
						childEval := Evaluate(position)
						e.staticEvals[searchHeight] = childEval
						score = -e.quiescence(-probBeta, -probBeta+1, searchHeight+1)
						e.pred.Pop()
					}

					if depthLeft < 8 || score >= probBeta {
						e.innerLines[searchHeight+1].Recycle()
						e.pred.Push(position.Hash())
						e.positionMoves[searchHeight+1] = move
						score = -e.alphaBeta(depthLeft-4, searchHeight+1, -probBeta, -probBeta+1)
						e.pred.Pop()
					}
					position.UnMakeMove(move, oldTag, oldEnPassant, hc)

					if score >= probBeta {
						e.info.probCutCounter += 1
						return score
					}
				}
				i += 1
			}
		}
	}

	// Internal Iterative Deepening
	if isPvNode && depthLeft >= 8 && !ttHit {
		e.innerLines[searchHeight].Recycle()
		score := e.alphaBeta(depthLeft-7, searchHeight, alpha, beta)
		if e.TimeManager.AbruptStop {
			return score
		}
		line := e.innerLines[searchHeight]
		if line.moveCount != 0 {
			hashmove := e.innerLines[searchHeight].MoveAt(0)
			nHashMove = hashmove
		}
		e.innerLines[searchHeight].Recycle()
	}

	movePicker := e.MovePickers[searchHeight]
	movePicker.RecycleWith(position, e, depthLeft, nHashMove, false)

	futilityMargin := eval + int16(depthLeft)*p
	if improving {
		futilityMargin += p
	}
	allowFutilityPruning := false
	if depthLeft < 7 && pruningAllowed &&
		abs16(alpha) < WIN_IN_MAX &&
		abs16(beta) < WIN_IN_MAX && futilityMargin <= alpha {
		allowFutilityPruning = true
	}

	oldAlpha := alpha

	// using fail soft with negamax:
	var bestscore int16
	var hashmove Move
	legalMoves := 1
	quietMoves := -1
	legalQuiteMove := -1
	noisyMoves := -1
	for true {
		hashmove = movePicker.Next()
		if hashmove == EmptyMove {
			break
		}
		isQuiet := false
		if hashmove.IsCapture() || hashmove.PromoType() != NoType {
			noisyMoves += 1
		} else {
			isQuiet = true
			quietMoves += 1
		}
		if oldEnPassant, oldTag, hc, ok := position.MakeMove(hashmove); ok {
			if isQuiet {
				legalQuiteMove += 1
			}
			e.pred.Push(position.Hash())
			e.innerLines[searchHeight+1].Recycle()
			e.positionMoves[searchHeight+1] = hashmove
			e.NoteMove(hashmove, legalQuiteMove, searchHeight)
			bestscore = -e.alphaBeta(depthLeft-1, searchHeight+1, -beta, -alpha)
			e.pred.Pop()
			position.UnMakeMove(hashmove, oldTag, oldEnPassant, hc)
			if bestscore > alpha {
				if bestscore >= beta {
					if !e.TimeManager.AbruptStop {
						e.TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, LowerBound, e.Ply)
						e.AddHistory(hashmove, hashmove.MovingPiece(), hashmove.Destination(), depthLeft, searchHeight, legalQuiteMove)
					}
					return bestscore
				}
				// Potential PV move, lets copy it to the current pv-line
				e.innerLines[searchHeight].AddFirst(hashmove)
				e.innerLines[searchHeight].ReplaceLine(e.innerLines[searchHeight+1])
				alpha = bestscore
			}
			break
		}
	}

	if hashmove == EmptyMove {
		if isInCheck {
			return -CHECKMATE_EVAL + int16(searchHeight)
		} else {
			return 0
		}
	}
	pruningThreashold := int(5 + depthLeft*depthLeft)
	if !improving {
		pruningThreashold /= 2
	}

	lmrThreashold := 2
	if isPvNode {
		lmrThreashold += 1
	}
	seeScores := movePicker.captureMoveList.Scores
	quietScores := movePicker.quietMoveList.Scores
	var move Move
	for true {

		if isRootNode {
			if e.TimeManager.ShouldStop(true, bestscore-e.score >= -20) {
				break
			}
		}

		move = movePicker.Next()
		if move == EmptyMove {
			break
		}

		isCaptureMove := move.IsCapture()
		promoType := move.PromoType()
		isQuiet := false
		if isCaptureMove || promoType != NoType {
			noisyMoves += 1
		} else {
			isQuiet = true
			quietMoves += 1
		}

		if oldEnPassant, oldTag, hc, ok := position.MakeMove(move); ok {
			legalMoves += 1
			if isQuiet {
				legalQuiteMove += 1
			}

			if e.DebugMode && isRootNode {
				fmt.Printf("info depth %d currmove %s currmovenumber %d\n", depthLeft, move.ToString(), legalMoves)
			}

			e.NoteMove(move, legalQuiteMove, searchHeight)
			isCheckMove := position.IsInCheck()
			notPromoting := !IsPromoting(move)
			LMR := int8(0)

			killerScore := e.KillerMoveScore(move, searchHeight)
			if !isInCheck && e.doPruning && !isRootNode && bestscore > -WIN_IN_MAX {

				if allowFutilityPruning &&
					!isCheckMove && notPromoting &&
					(!isCaptureMove ||
						depthLeft <= 1 && isCaptureMove && seeScores[noisyMoves] < 0) {
					e.info.efpCounter += 1
					position.UnMakeMove(move, oldTag, oldEnPassant, hc)
					continue
				}

				// Late Move Pruning
				if notPromoting && !isCaptureMove && !isCheckMove && depthLeft <= 8 &&
					legalMoves > pruningThreashold && killerScore <= 0 && abs16(alpha) < WIN_IN_MAX {
					e.info.lmpCounter += 1
					position.UnMakeMove(move, oldTag, oldEnPassant, hc)
					continue // LMP
				}

				// SEE pruning
				if isCaptureMove && seeScores[noisyMoves] < 0 &&
					!isCheckMove && depthLeft <= 2 && eval <= alpha && abs16(alpha) < WIN_IN_MAX {
					e.info.seeCounter += 1
					position.UnMakeMove(move, oldTag, oldEnPassant, hc)
					continue
				}
			}

			// Late Move Reduction
			if !isInCheck && e.doPruning && isQuiet && !isCheckMove && depthLeft > 2 && legalMoves > lmrThreashold {
				e.info.lmrCounter += 1
				LMR = int8(lmrReductions[min8(31, depthLeft)][min(31, legalMoves)])

				if killerScore > 0 {
					LMR -= 1
				}

				if isPvNode {
					LMR -= 1
				}

				if improving {
					LMR -= 1
				}

				if killerScore <= 0 && quietScores[quietMoves] > 0 {
					LMR -= 1
				}

				LMR = min8(depthLeft-2, max8(LMR, 1))
			}

			e.pred.Push(position.Hash())
			e.innerLines[searchHeight+1].Recycle()
			e.positionMoves[searchHeight+1] = move
			score := -e.alphaBeta(depthLeft-1-LMR, searchHeight+1, -alpha-1, -alpha)
			e.pred.Pop()
			if score > alpha && score < beta {
				e.info.researchCounter += 1
				// research with window [alpha;beta]
				e.pred.Push(position.Hash())
				e.innerLines[searchHeight+1].Recycle()
				score = -e.alphaBeta(depthLeft-1, searchHeight+1, -beta, -alpha)
				e.pred.Pop()
				if score > alpha {
					alpha = score
				}
			}
			position.UnMakeMove(move, oldTag, oldEnPassant, hc)

			if score > bestscore {
				if score >= beta {
					if !e.TimeManager.AbruptStop {
						e.TranspositionTable.Set(hash, move, score, depthLeft, LowerBound, e.Ply)
						e.AddHistory(move, move.MovingPiece(), move.Destination(), depthLeft, searchHeight, legalQuiteMove)
					}
					return score
				}
				// Potential PV move, lets copy it to the current pv-line
				e.innerLines[searchHeight].AddFirst(move)
				e.innerLines[searchHeight].ReplaceLine(e.innerLines[searchHeight+1])
				bestscore = score
				hashmove = move
			}
		}
	}
	if !e.TimeManager.AbruptStop {
		if alpha > oldAlpha {
			e.TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, Exact, e.Ply)
		} else {
			e.TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, UpperBound, e.Ply)
		}
	}
	if isRootNode && legalMoves == 1 {
		e.TimeManager.StopSearchNow = true
	}
	return bestscore
}
