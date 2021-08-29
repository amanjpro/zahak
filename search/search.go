package search

import (
	"fmt"
	"math"
	"sync"

	. "github.com/amanjpro/zahak/book"
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func (r *Runner) Search(depth int8) {

	if len(r.Engines) == 1 {
		e := r.Engines[0]
		e.Search(depth)
	} else {
		r.ClearForSearch()
		var wg sync.WaitGroup
		for i := 0; i < len(r.Engines); i++ {
			wg.Add(1)
			go func(e *Engine, depth int8, i int) {
				e.ParallelSearch(depth, int8(1+i%2), 2)
				wg.Done()
			}(r.Engines[i], depth, i)
		}
		wg.Wait()
		r.SendBestMove()
	}
}

func (e *Engine) ParallelSearch(depth int8, start int8, inc int8) {
	e.ClearForSearch()
	e.rootSearch(depth, start, inc)
}

func (e *Engine) Search(depth int8) {
	e.parent.ClearForSearch()
	e.ClearForSearch()
	e.rootSearch(depth, 1, 1)
	e.parent.SendBestMove()
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

func (e *Engine) updatePv(pvLine PVLine, score int16, depth int8, isBookmove bool) (PVLine, int16, int8, bool) {
	parent := e.parent
	parent.mu.Lock()
	updated := false
	if isBookmove || (parent.depth < depth && !parent.isBookmove) {
		parent.pv.Clone(pvLine)
		parent.move = parent.pv.MoveAt(0)
		parent.score = score
		parent.depth = depth
		parent.isBookmove = isBookmove
		updated = true
	} else {
		score = parent.score
		depth = parent.depth
		pvLine.Clone(parent.pv)
	}
	parent.mu.Unlock()
	return pvLine, score, depth, updated
}

func (e *Engine) rootSearch(depth int8, startDepth int8, depthIncrement int8) {
	pv := NewPVLine(MAX_DEPTH)

	lastDepth := int8(1)

	bookmove := GetBookMove(e.Position)
	if e.isMainThread && bookmove != EmptyMove {
		pv.Recycle()
		pv.AddFirst(bookmove)
		pv, e.score, _, _ = e.updatePv(pv, 0, 1, true)
	} else {
		for iterationDepth := startDepth; iterationDepth <= depth; iterationDepth += depthIncrement {

			if e.isMainThread {
				if iterationDepth > 1 && !e.TimeManager().CanStartNewIteration() {
					break
				}
			} else if iterationDepth > 1 && e.parent.Stop {
				break
			}

			var bookmove bool
			var globalDepth int8
			e.parent.mu.RLock()
			globalDepth = e.parent.depth
			bookmove = e.parent.isBookmove
			e.score = e.parent.score
			e.parent.mu.RUnlock()

			if bookmove {
				break
			}

			if iterationDepth <= globalDepth {
				continue
			}

			e.innerLines[0].Recycle()
			e.startDepth = iterationDepth
			newScore := e.aspirationWindow(e.score, iterationDepth)

			if (e.isMainThread && e.TimeManager().AbruptStop) || (!e.isMainThread && e.parent.Stop) {
				break
			}
			if e.startDepth == 0 {
				continue
			}
			pv.Clone(e.innerLines[0])

			if e.isMainThread && iterationDepth >= 8 && e.score-newScore >= 30 { // Position degrading
				e.TimeManager().ExtraTime()
			}

			var newDepth int8
			var updated bool
			pv, e.score, newDepth, updated = e.updatePv(pv, newScore, iterationDepth, false)

			lastDepth = newDepth
			e.pred.Clear()
			e.ShareInfo()
			if updated {
				e.SendPv(pv, e.score, newDepth)
			}
			if e.isMainThread && !e.TimeManager().Pondering && e.parent.DebugMode {
				e.parent.globalInfo.Print()
			}
			if isCheckmateEval(e.score) {
				break
			}
		}
	}

	if e.isMainThread {
		e.TimeManager().Pondering = false
		e.parent.Stop = true
		e.SendPv(pv, e.score, lastDepth)
	}
}

func (e *Engine) aspirationWindow(prevScore int16, iterationDepth int8) int16 {
	e.doPruning = iterationDepth > 3
	if iterationDepth <= 6 {
		return e.alphaBeta(iterationDepth, 0, -MAX_INT, MAX_INT)
	} else {
		alphaMargin := int16(25)
		betaMargin := int16(25)
		for i := 0; i < 2; i++ {
			alpha := max16(prevScore-alphaMargin, -MAX_INT)
			beta := min16(prevScore+betaMargin, MAX_INT)
			score := e.alphaBeta(iterationDepth, 0, alpha, beta)
			if e.startDepth == 0 {
				return -MAX_INT
			}
			if score <= alpha {
				alphaMargin *= 2
			} else if score >= beta {
				betaMargin *= 2
			} else {
				return score
			}
		}
	}
	return e.alphaBeta(iterationDepth, 0, -MAX_INT, MAX_INT)
}

func (e *Engine) alphaBeta(depthLeft int8, searchHeight int8, alpha int16, beta int16) int16 {
	e.VisitNode()

	isRootNode := searchHeight == 0
	isPvNode := alpha != beta-1

	position := e.Position
	pawnhash := e.Pawnhash

	currentMove := e.positionMoves[searchHeight]
	// Position is drawn
	if IsRepetition(position, e.pred, currentMove) || position.IsDraw() {
		return 0
	}

	if searchHeight >= MAX_DEPTH-1 {
		eval := Evaluate(position, pawnhash)
		e.staticEvals[searchHeight] = eval
		return eval
	}

	var isInCheck = position.IsInCheck()
	if isInCheck {
		e.info.checkExtentionCounter += 1
		depthLeft += 1 // Check Extension
	}

	if depthLeft <= 0 {
		e.staticEvals[searchHeight] = Evaluate(position, pawnhash)
		return e.quiescence(alpha, beta, searchHeight)
	}

	if isPvNode {
		e.info.mainSearchCounter += 1
	} else {
		e.info.zwCounter += 1
	}

	firstLayerOfSingularity := e.skipHeight == searchHeight && e.skipMove != EmptyMove
	hash := position.Hash()
	nHashMove, nEval, nDepth, nType, ttHit := e.TranspositionTable.Get(hash)
	ttHit = ttHit && position.IsPseudoLegal(nHashMove)
	if !ttHit {
		nHashMove = EmptyMove
	}

	if !isPvNode && ttHit && nDepth >= depthLeft && !firstLayerOfSingularity {
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
		if e.isMainThread && e.TimeManager().ShouldStop(false, false) {
			return -MAX_INT
		}
	}

	if !e.isMainThread && e.parent.Stop {
		return -MAX_INT
	}

	var eval int16 = -MAX_INT
	if !isRootNode && currentMove == EmptyMove {
		eval = -1 * (e.staticEvals[searchHeight-1] + Tempo + Tempo)
	} else {
		eval = Evaluate(position, pawnhash)
	}

	e.staticEvals[searchHeight] = eval
	improving := currentMove == EmptyMove ||
		(searchHeight > 2 && e.staticEvals[searchHeight] > e.staticEvals[searchHeight-2])

	// Pruning
	pruningAllowed := !isPvNode && !isInCheck && e.doPruning && !firstLayerOfSingularity

	if pruningAllowed {
		// Razoring
		razoringMargin := eval + int16(depthLeft)*p + p
		if depthLeft < 3 && eval+razoringMargin < beta {
			newEval := e.quiescence(alpha, beta, searchHeight)
			e.info.razoringCounter += 1
			return newEval
		}

		// Reverse Futility Pruning
		reverseFutilityMargin := int16(depthLeft) * p //(b - p)
		if improving {
			reverseFutilityMargin += p // int16(depthLeft) * p
		}
		if depthLeft < 8 && eval-reverseFutilityMargin >= beta {
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
			movePicker.RecycleWith(position, e, depthLeft, searchHeight, hashMove, true)
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
						childEval := Evaluate(position, pawnhash)
						e.staticEvals[searchHeight+1] = childEval
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
		if e.isMainThread && e.TimeManager().AbruptStop {
			return score
		} else if !e.isMainThread && e.parent.Stop {
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
	movePicker.RecycleWith(position, e, depthLeft, searchHeight, nHashMove, false)
	oldAlpha := alpha

	// using fail soft with negamax:
	var bestscore int16
	var hashmove Move
	legalMoves := 0
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
		if e.skipMove == hashmove && firstLayerOfSingularity {
			legalMoves += 1
			continue
		}
		if oldEnPassant, oldTag, hc, ok := position.MakeMove(hashmove); ok {
			legalMoves += 1
			if isQuiet {
				legalQuiteMove += 1
			}
			// Singular Extension
			var extension int8
			if depthLeft > 8 &&
				hashmove == nHashMove &&
				ttHit &&
				e.skipMove == EmptyMove &&
				nDepth > depthLeft-3 &&
				nType != UpperBound &&
				!position.IsInCheck() && // Check moves are automatically extended
				abs16(nEval) < WIN_IN_MAX &&
				!isRootNode {

				// ttMove has been made to check legality
				position.UnMakeMove(hashmove, oldTag, oldEnPassant, hc)

				// Search to reduced depth with a zero window a bit lower than ttScore
				threshold := max16(nEval-2*int16(depthLeft), -CHECKMATE_EVAL)

				e.skipMove = hashmove
				e.skipHeight = searchHeight
				e.innerLines[searchHeight].Recycle()
				e.MovePickers[searchHeight] = e.TempMovePicker
				score := e.alphaBeta(depthLeft/2-1, searchHeight, threshold-1, threshold)
				e.MovePickers[searchHeight] = movePicker
				e.innerLines[searchHeight].Recycle()
				e.skipMove = EmptyMove
				e.skipHeight = MAX_DEPTH

				// Extend as this move seems forced
				if score < threshold {
					e.info.singularExtensionCounter += 1
					extension += 1
				}

				// Multi-Cut, at least 2 moves beat beta, idea is taken from Stockfish
				if pruningAllowed {
					if threshold >= beta {
						e.info.multiCutCounter += 1
						return beta
					} else if score >= beta {
						e.skipHeight = 0
						e.skipMove = hashmove
						e.innerLines[searchHeight].Recycle()
						e.MovePickers[searchHeight] = e.TempMovePicker
						score = e.alphaBeta((depthLeft+3)/2, searchHeight, beta-1, beta)
						e.MovePickers[searchHeight] = movePicker
						e.innerLines[searchHeight].Recycle()
						e.skipMove = EmptyMove
						e.skipHeight = MAX_DEPTH
						e.info.multiCutCounter += 1

						if score >= beta {
							return beta
						}
					}
				}

				// Replay ttMove
				position.MakeMove(hashmove)
			}

			e.pred.Push(position.Hash())
			e.innerLines[searchHeight+1].Recycle()
			e.positionMoves[searchHeight+1] = hashmove
			e.NoteMove(hashmove, legalQuiteMove, searchHeight)
			bestscore = -e.alphaBeta(depthLeft-1+extension, searchHeight+1, -beta, -alpha)
			e.pred.Pop()
			position.UnMakeMove(hashmove, oldTag, oldEnPassant, hc)
			if bestscore > alpha {
				if bestscore >= beta {
					if (e.isMainThread && !e.TimeManager().AbruptStop) || (!e.isMainThread && !e.parent.Stop) {
						if !firstLayerOfSingularity {
							e.TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, LowerBound, e.Ply)
						}
						e.AddHistory(hashmove, currentMove, hashmove.MovingPiece(), hashmove.Destination(), depthLeft, searchHeight, legalQuiteMove)
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
			if firstLayerOfSingularity {
				return alpha
			}
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
	var historyThreashold int32 = int32(depthLeft) * -1024
	var move Move
	for true {

		if isRootNode {
			if e.isMainThread && e.TimeManager().ShouldStop(true, bestscore-e.score >= -20) {
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
		notPromoting := move.PromoType() == NoType
		isKiller := movePicker.killer1 == move || movePicker.killer2 == move

		if !isInCheck && e.doPruning && !isRootNode && bestscore > -WIN_IN_MAX {

			// Late Move Pruning
			if notPromoting && !isCaptureMove && depthLeft <= 8 &&
				legalMoves+1 > pruningThreashold && !isKiller && abs16(alpha) < WIN_IN_MAX {
				e.info.lmpCounter += 1
				// This is a hack really, mp.Next() won't return any quiets, and I am hacking this
				// to avoid returning quiets, after the first LMP cut
				movePicker.isQuiescence = true
				continue // LMP
			}

			// SEE pruning
			if isCaptureMove && seeScores[noisyMoves] < 0 &&
				/* !isCheckMove && */ depthLeft <= 2 && eval <= alpha && abs16(alpha) < WIN_IN_MAX {
				e.info.seeCounter += 1
				// position.UnMakeMove(move, oldTag, oldEnPassant, hc)
				continue
			}

			// History pruning
			lmrDepth := depthLeft - int8(lmrReductions[min8(31, depthLeft)][min(31, legalMoves+1)])
			if /* !isKiller && */ /* !isCheckMove && */ isQuiet && quietScores[quietMoves] < historyThreashold && lmrDepth < 3 && legalMoves+1 > lmrThreashold {
				e.info.historyPruningCounter += 1
				// position.UnMakeMove(move, oldTag, oldEnPassant, hc)
				continue
			}
		}

		if oldEnPassant, oldTag, hc, ok := position.MakeMove(move); ok {
			legalMoves += 1
			if isQuiet {
				legalQuiteMove += 1
			}

			if e.isMainThread && e.parent.DebugMode && isRootNode {
				fmt.Printf("info depth %d currmove %s currmovenumber %d\n", depthLeft, move.ToString(), legalMoves)
			}

			e.NoteMove(move, legalQuiteMove, searchHeight)
			LMR := int8(0)

			// Late Move Reduction
			if !isInCheck && e.doPruning && isQuiet && depthLeft > 2 && legalMoves > lmrThreashold {
				e.info.lmrCounter += 1
				LMR = int8(lmrReductions[min8(31, depthLeft)][min(31, legalMoves)])

				// if killerScore > 0 {
				// 	LMR -= 1
				// }

				if isPvNode {
					LMR -= 1
				}

				if improving {
					LMR -= 1
				}

				LMR -= int8(e.MoveHistoryScore(move.MovingPiece(), move.Destination(), depthLeft) / 24576)

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
					if (e.isMainThread && !e.TimeManager().AbruptStop) || (!e.isMainThread && !e.parent.Stop) {
						if !firstLayerOfSingularity {
							e.TranspositionTable.Set(hash, move, score, depthLeft, LowerBound, e.Ply)
						}
						e.AddHistory(move, currentMove, move.MovingPiece(), move.Destination(), depthLeft, searchHeight, legalQuiteMove)
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

		if !isRootNode && e.startDepth > 6 {
			e.parent.mu.RLock()
			if e.startDepth <= e.parent.depth {
				e.startDepth = 0
				e.parent.mu.RUnlock()
				return -MAX_INT
			}
			e.parent.mu.RUnlock()
		}
	}
	if (e.isMainThread && !e.TimeManager().AbruptStop) || (!e.isMainThread && !e.parent.Stop) && !firstLayerOfSingularity {
		if alpha > oldAlpha {
			e.TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, Exact, e.Ply)
		} else {
			e.TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, UpperBound, e.Ply)
		}
	}
	if e.isMainThread && isRootNode && legalMoves == 1 {
		e.TimeManager().StopSearchNow = true
	}
	return bestscore
}
