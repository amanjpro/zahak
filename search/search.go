package search

import (
	"fmt"
	"math"
	"sync"

	. "github.com/amanjpro/zahak/book"
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/fathom"
)

const TB_WIN_BOUND int16 = 27000
const TB_LOSS_BOUND int16 = -27000

func (r *Runner) Search(depth int8) {
	e := r.Engines[0]
	if bestmove := ProbeDTZ(e.Position); bestmove != EmptyMove {
		for e.TimeManager().Pondering {
			// busy waiting
		}
		fmt.Printf("bestmove %s\n", bestmove.ToString())
		return
	}

	if len(r.Engines) == 1 {
		e := r.Engines[0]
		e.Search(depth)
	} else {
		r.ClearForSearch()
		var wg sync.WaitGroup
		for i := 0; i < len(r.Engines); i++ {
			wg.Add(1)
			go func(e *Engine, depth int8, i int) {
				e.ParallelSearch(depth, 1, 1)
				wg.Done()
			}(r.Engines[i], depth, i)
		}
		wg.Wait()
		for e.TimeManager().Pondering {
			// busy waiting
		}
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
	for e.TimeManager().Pondering {
		// busy waiting
	}
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

func (e *Engine) updatePv(pvLine PVLine, score int16, depth int8, isBookmove bool) {
	parent := e.parent
	parent.pv.Clone(pvLine)
	parent.move = parent.pv.MoveAt(0)
	parent.score = score
	parent.depth = depth
	parent.isBookmove = isBookmove
}

func (e *Engine) rootSearch(depth int8, startDepth int8, depthIncrement int8) {
	pv := NewPVLine(MAX_DEPTH)

	lastDepth := int8(1)

	bookmove := GetBookMove(e.Position)
	if e.isMainThread && bookmove != EmptyMove {
		pv.Recycle()
		pv.AddFirst(bookmove)
		e.updatePv(pv, 0, 1, true)
	} else {
		for iterationDepth := startDepth; iterationDepth <= depth; iterationDepth += depthIncrement {

			if e.isMainThread {
				if iterationDepth > 1 && !e.TimeManager().CanStartNewIteration() {
					break
				}
			} else if iterationDepth > 1 && e.parent.Stop {
				break
			}

			e.innerLines[0].Recycle()
			e.startDepth = iterationDepth
			newScore := e.aspirationWindow(e.score, iterationDepth)

			if (e.isMainThread && e.TimeManager().AbruptStop) || (!e.isMainThread && e.parent.Stop) {
				break
			}

			if e.isMainThread && iterationDepth >= 8 && e.score-newScore >= 30 { // Position degrading
				e.TimeManager().ExtraTime()
			}

			lastDepth = iterationDepth
			e.pred.Clear()
			e.score = newScore
			if e.isMainThread && !e.innerLines[0].IsEmpty() {
				pv.Clone(e.innerLines[0])
				e.SendPv(pv, e.score, iterationDepth)
			}
			if e.isMainThread && !e.TimeManager().Pondering && e.parent.DebugMode {
				e.info.Print()
			}
			// if isCheckmateEval(e.score) {
			// 	break
			// }
		}
	}
	if e.isMainThread {
		// e.TimeManager().Pondering = false
		e.parent.Stop = true
		e.updatePv(pv, e.score, lastDepth, false)
		e.SendPv(pv, e.score, lastDepth)
	}
}

func (e *Engine) aspirationWindow(score int16, iterationDepth int8) int16 {
	e.doPruning = iterationDepth > 3
	if iterationDepth <= 6 {
		e.seldepth = 0
		return e.alphaBeta(iterationDepth, 0, -MAX_INT, MAX_INT)
	} else {
		var initialWindow int16 = 12
		var delta int16 = 16

		alpha := max16(score-initialWindow, -MAX_INT)
		beta := min16(score+initialWindow, MAX_INT)
		originalDepth := iterationDepth
		maxSeldepth := int8(0)
		e.seldepth = 0

		for true {
			beta = min16(beta, MAX_INT)
			alpha = max16(alpha, -MAX_INT)

			score = e.alphaBeta(iterationDepth, 0, alpha, beta)
			if /* e.startDepth == 0 || */ e.TimeManager().AbruptStop || e.parent.Stop {
				e.seldepth = max8(e.seldepth, maxSeldepth)
				return -MAX_INT
			}
			if score <= alpha {
				alpha = max16(alpha-delta, -MAX_INT)
				beta = (alpha + 3*beta) / 4
				iterationDepth = originalDepth
			} else if score >= beta {
				beta = min16(beta+delta, MAX_INT)
				if abs16(score) < WIN_IN_MAX {
					iterationDepth = max8(1, iterationDepth-1)
				}
			} else {
				e.seldepth = max8(e.seldepth, maxSeldepth)
				return score
			}
			delta += delta * 2 / 3
			maxSeldepth = max8(e.seldepth, maxSeldepth)
		}
		e.seldepth = max8(e.seldepth, maxSeldepth)
	}
	// We should never get here
	return -MAX_INT
}

func (e *Engine) alphaBeta(depthLeft int8, searchHeight int8, alpha int16, beta int16) int16 {
	e.VisitNode(searchHeight)

	isRootNode := searchHeight == 0
	isPvNode := alpha != beta-1

	position := e.Position
	TranspositionTable.Prefetch(position.Hash())

	currentMove := e.positionMoves[searchHeight]
	var gpMove Move
	if searchHeight > 1 {
		gpMove = e.positionMoves[searchHeight-1]
	}
	if !isRootNode {
		// Position is drawn
		if IsRepetition(position, e.pred, currentMove) || position.IsDraw() {
			return 0
		}

		// Mate distance pruning
		alpha = max16(alpha, -CHECKMATE_EVAL+int16(searchHeight))
		beta = min16(beta, CHECKMATE_EVAL-int16(searchHeight)-1)
		if alpha >= beta {
			return alpha
		}
	}

	if searchHeight >= MAX_DEPTH-1 {
		eval := position.Evaluate()
		e.staticEvals[searchHeight] = eval
		return eval
	}

	var isInCheck = position.IsInCheck()
	if isInCheck {
		e.info.checkExtentionCounter += 1
		depthLeft += 1 // Check Extension
	}

	if depthLeft <= 0 {
		e.staticEvals[searchHeight] = position.Evaluate()
		return e.quiescence(alpha, beta, searchHeight)
	}

	if isPvNode {
		e.info.mainSearchCounter += 1
	} else {
		e.info.zwCounter += 1
	}

	firstLayerOfSingularity := e.skipHeight == searchHeight && e.skipMove != EmptyMove
	hash := position.Hash()
	nHashMove, nEval, nDepth, nType, ttHit := TranspositionTable.Get(hash)
	if ttHit {
		ttHit = position.IsPseudoLegal(nHashMove)
		nEval = evalFromTT(nEval, searchHeight)
	}
	if !ttHit {
		nHashMove = EmptyMove
	}

	if !isPvNode && ttHit && nDepth >= depthLeft && !firstLayerOfSingularity {
		if nEval >= beta && nType == LowerBound {
			e.CacheHit()
			e.searchHistory.AddHistory(nHashMove, currentMove, gpMove, depthLeft, searchHeight, nil)
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

	// tablebase - we do not do this at root
	if !isRootNode {
		tbResult := ProbeWDL(position, depthLeft)

		if tbResult != TB_RESULT_FAILED {
			e.info.tbHit += 1

			var flag NodeType
			var score int16
			switch tbResult {
			case TB_WIN:
				score = TB_WIN_BOUND - int16(searchHeight)
				flag = LowerBound
			case TB_LOSS:
				score = -TB_WIN_BOUND + int16(searchHeight)
				flag = UpperBound
			default:
				score = 0
				flag = Exact
			}

			// if the tablebase gives us what we want, then we accept it's score and return
			if flag == Exact || (flag == LowerBound && score >= beta) || (flag == UpperBound && score <= alpha) {
				TranspositionTable.Set(hash, EmptyMove, score, depthLeft, flag, e.Ply)
				return score
			}

			// for pv node searches we adjust our a/b search accordingly
			if isPvNode {
				if flag == LowerBound {
					alpha = max16(alpha, score)
				}
			}
		}
	}

	// Internal iterative reduction based on Rebel's idea
	// if /* !isPvNode && */ !ttHit && depthLeft >= 3 {
	// 	e.info.internalIterativeReduction += 1
	// 	depthLeft -= 1
	// }

	if !isRootNode {
		if e.isMainThread && e.TimeManager().ShouldStop(false, false) {
			return -MAX_INT
		}
	}

	if !e.isMainThread && e.parent.Stop {
		return -MAX_INT
	}

	var eval int16 = position.Evaluate() //-MAX_INT
	// if !isRootNode && currentMove == EmptyMove {
	// 	eval = -1 * e.staticEvals[searchHeight-1]
	// } else {
	// eval =
	// }

	e.staticEvals[searchHeight] = eval
	improving := currentMove == EmptyMove ||
		(searchHeight > 2 && e.staticEvals[searchHeight] > e.staticEvals[searchHeight-2])

	// Pruning
	pruningAllowed := !isPvNode && !isInCheck && e.doPruning && !firstLayerOfSingularity

	// Idea taken from Berserk
	histDepth := depthLeft
	if eval+p > beta {
		histDepth += 1
	}

	if pruningAllowed {
		// Razoring
		if depthLeft < 2 && eval+350 <= alpha {
			newEval := e.quiescence(alpha, beta, searchHeight)
			e.info.razoringCounter += 1
			return newEval
		}

		// Reverse Futility Pruning
		reverseFutilityMargin := int16(depthLeft) * 85 //(b - p)
		if improving {
			reverseFutilityMargin -= 85 // int16(depthLeft) * p
		}
		if depthLeft < 8 && eval-reverseFutilityMargin >= beta {
			e.info.rfpCounter += 1
			return eval - reverseFutilityMargin /* fail soft */
		}

		// NullMove pruning
		isNullMoveAllowed := currentMove != EmptyMove && !position.IsEndGame()
		if isNullMoveAllowed && depthLeft >= 2 && eval > beta {
			var R = 4 + min8(depthLeft/4, 3)
			if eval >= beta+100 {
				R += 1
			}
			R = min8(R, depthLeft)
			// } else {
			// 	R = min8(R, depthLeft-1)
			// }
			// if R >= 2 {
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
			// }
		}

		// Threat pruning, idea from Koivisto
		tpMargin := int16(0)
		if improving {
			tpMargin = 30
		}
		if depthLeft == 1 && eval > beta+tpMargin && (!position.Board.HasThreats(position.Turn().Other()) || position.Board.HasThreats(position.Turn())) {
			return beta
		}

		// Prob cut
		// The idea is basically cherry picked from multiple engines, Weiss, Ethereal and Berserk for example
		// probBeta := min16(beta+120, WIN_IN_MAX)
		// if depthLeft > 4 && abs16(beta) < WIN_IN_MAX && !(ttHit && nDepth >= depthLeft-3 && nEval < probBeta) {
		//
		// 	hashMove := EmptyMove
		// 	if hashMove.IsCapture() || hashMove.PromoType() != NoType {
		// 		hashMove = nHashMove
		// 	}
		// 	movePicker := e.MovePickers[searchHeight]
		// 	movePicker.RecycleWith(position, e, depthLeft, searchHeight, hashMove, true)
		// 	seeScores := movePicker.captureMoveList.Scores
		// 	i := 0
		// 	for true {
		// 		move := movePicker.Next()
		// 		if move == EmptyMove || seeScores[i] < 0 {
		// 			break
		// 		}
		// 		if oldEnPassant, oldTag, hc, ok := position.MakeMove(move); ok {
		// 			var score int16
		// 			if depthLeft >= 8 {
		// 				e.innerLines[searchHeight+1].Recycle()
		// 				e.pred.Push(position.Hash())
		// 				e.positionMoves[searchHeight+1] = move
		// 				childEval := position.Evaluate()
		// 				e.staticEvals[searchHeight+1] = childEval
		// 				score = -e.quiescence(-probBeta, -probBeta+1, searchHeight+1)
		// 				e.pred.Pop()
		// 			}
		//
		// 			if depthLeft < 8 || score >= probBeta {
		// 				e.innerLines[searchHeight+1].Recycle()
		// 				e.pred.Push(position.Hash())
		// 				e.positionMoves[searchHeight+1] = move
		// 				score = -e.alphaBeta(depthLeft-4, searchHeight+1, -probBeta, -probBeta+1)
		// 				e.pred.Pop()
		// 			}
		// 			position.UnMakeMove(move, oldTag, oldEnPassant, hc)
		//
		// 			if score >= probBeta {
		// 				e.info.probCutCounter += 1
		// 				return score
		// 			}
		// 		}
		// 		i += 1
		// 	}
		// }
	}

	// Internal Iterative Deepening
	if /* isPvNode && */ depthLeft >= 8 && !ttHit {
		e.innerLines[searchHeight].Recycle()
		score := e.alphaBeta(depthLeft-7, searchHeight, alpha, beta)
		if e.isMainThread && e.TimeManager().AbruptStop {
			return score
		} else if !e.isMainThread && e.parent.Stop {
			return score
		}
		line := e.innerLines[searchHeight]
		if line.moveCount != 0 {
			nHashMove = e.innerLines[searchHeight].MoveAt(0)
		}
		e.innerLines[searchHeight].Recycle()
	}

	movePicker := e.MovePickers[searchHeight]
	movePicker.RecycleWith(position, e, searchHeight, nHashMove, false)
	oldAlpha := alpha

	// using fail soft with negamax:
	var bestscore int16
	var hashmove Move
	legalMoves := 0
	quietMoves := -1
	legalQuietMove := 0
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
		if isRootNode && e.mustSkip(hashmove) {
			continue
		}
		if oldEnPassant, oldTag, hc, ok := position.MakeMove(hashmove); ok {
			legalMoves += 1
			if isQuiet {
				legalQuietMove += 1
			}
			// Singular Extension
			var extension int8
			if depthLeft >= 8 &&
				hashmove == nHashMove &&
				ttHit &&
				e.skipMove == EmptyMove &&
				nDepth >= depthLeft-3 &&
				nType != UpperBound &&
				!position.IsInCheck() && // Check moves are automatically extended
				abs16(nEval) < WIN_IN_MAX &&
				!isRootNode {

				// ttMove has been made to check legality
				position.UnMakeMove(hashmove, oldTag, oldEnPassant, hc)

				// Search to reduced depth with a zero window a bit lower than ttScore
				threshold := max16(nEval-3*int16(depthLeft)/2, -CHECKMATE_EVAL)

				e.skipMove = hashmove
				e.skipHeight = searchHeight
				e.innerLines[searchHeight].Recycle()
				e.MovePickers[searchHeight] = e.TempMovePicker
				score := e.alphaBeta((depthLeft-1)/2, searchHeight, threshold-1, threshold)
				e.MovePickers[searchHeight] = movePicker
				e.innerLines[searchHeight].Recycle()
				e.skipMove = EmptyMove
				e.skipHeight = MAX_DEPTH

				// Extend as this move seems forced
				if score < threshold {
					e.info.singularExtensionCounter += 1
					extension += 1
				} else {

					// Multi-Cut, at least 2 moves beat beta, idea is taken from Stockfish
					if pruningAllowed {
						if threshold >= beta {
							e.info.multiCutCounter += 1
							return beta
						} else if score >= beta {
							e.skipHeight = searchHeight
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
				}

				// Replay ttMove
				position.MakeMove(hashmove)
			}

			e.pred.Push(position.Hash())
			e.innerLines[searchHeight+1].Recycle()
			e.positionMoves[searchHeight+1] = hashmove
			e.NoteMove(hashmove, legalQuietMove-1, searchHeight)
			bestscore = -e.alphaBeta(depthLeft-1+extension, searchHeight+1, -beta, -alpha)
			e.pred.Pop()
			position.UnMakeMove(hashmove, oldTag, oldEnPassant, hc)
			if bestscore > alpha {
				if bestscore >= beta {
					if (e.isMainThread && !e.TimeManager().AbruptStop) || (!e.isMainThread && !e.parent.Stop) {
						if !firstLayerOfSingularity {
							TranspositionTable.Set(hash, hashmove, evalToTT(bestscore, searchHeight), depthLeft, LowerBound, e.Ply)
						}
						quietMoves := e.triedQuietMoves[searchHeight][:legalQuietMove]
						e.searchHistory.AddHistory(hashmove, currentMove, gpMove, histDepth, searchHeight, quietMoves)
					}
					return bestscore
				}
				// Potential PV move, lets copy it to the current pv-line
				e.innerLines[searchHeight].ReplaceLine(hashmove, e.innerLines[searchHeight+1])
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
	if !improving && !isPvNode {
		pruningThreashold = pruningThreashold/2 - 1
	}

	lmrThreashold := 2
	if isPvNode {
		lmrThreashold += 1
	}
	fpMargin := eval + p*int16(depthLeft)
	rangeReduction := 0
	if eval-bestscore < 30 && depthLeft > 7 {
		rangeReduction += 1
	}

	seeScores := movePicker.captureMoveList.Scores
	// quietScores := movePicker.quietMoveList.Scores
	// var historyThreashold int32 = int32(depthLeft) * -1024
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
		if isRootNode && e.mustSkip(move) {
			continue
		}
		isKiller := movePicker.killer1 == move || movePicker.killer2 == move || movePicker.counterMove == move

		if e.doPruning && !isRootNode && bestscore > -WIN_IN_MAX {

			if depthLeft < 8 && isQuiet && !isKiller && fpMargin <= alpha && abs16(alpha) < WIN_IN_MAX {
				e.info.efpCounter += 1
				continue
			}

			// Late Move Pruning
			if isQuiet && depthLeft < 8 &&
				legalMoves+1 > pruningThreashold && !isKiller && abs16(alpha) < WIN_IN_MAX {
				e.info.lmpCounter += 1
				// This is a hack really, mp.Next() won't return any quiets, and I am hacking this
				// to avoid returning quiets, after the first LMP cut
				movePicker.isQuiescence = true
				continue // LMP
			}

			// SEE pruning
			if isCaptureMove && seeScores[noisyMoves] < 0 &&
				depthLeft <= 4 && eval <= alpha && abs16(alpha) < WIN_IN_MAX {
				e.info.seeCounter += 1
				break
			}
			//
			// // History pruning
			// lmrDepth := depthLeft - int8(lmrReductions[min8(31, depthLeft)][min(31, legalMoves+1)])
			// if isQuiet && quietScores[quietMoves] < historyThreashold && lmrDepth < 3 && legalMoves+1 > lmrThreashold {
			// 	e.info.historyPruningCounter += 1
			// 	continue
			// }
		}

		if oldEnPassant, oldTag, hc, ok := position.MakeMove(move); ok {
			legalMoves += 1
			if isQuiet {
				legalQuietMove += 1
			}

			if e.isMainThread && e.parent.DebugMode && isRootNode {
				fmt.Printf("info depth %d currmove %s currmovenumber %d\n", depthLeft, move.ToString(), legalMoves)
			}

			e.NoteMove(move, legalQuietMove-1, searchHeight)
			LMR := int8(0)

			// Late Move Reduction
			if e.doPruning && (isQuiet || isCaptureMove && seeScores[noisyMoves] < 0) && depthLeft > 2 && legalMoves > lmrThreashold {
				e.info.lmrCounter += 1
				LMR = int8(lmrReductions[min8(31, depthLeft)][min(31, legalMoves)])

				if isInCheck {
					LMR -= 1
				}

				if isKiller {
					LMR -= 1
				}

				if isPvNode {
					LMR -= 1
				}

				if improving {
					LMR -= 1
				}

				// Credit to Ofek Shochat
				if rangeReduction > 3 {
					LMR += 1
				}

				if isQuiet {
					LMR -= int8(e.searchHistory.History(gpMove, currentMove, move) / 10649) //12288)
				} else {
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
					if (e.isMainThread && !e.TimeManager().AbruptStop) || (!e.isMainThread && !e.parent.Stop) {
						if !firstLayerOfSingularity {
							TranspositionTable.Set(hash, move, evalToTT(score, searchHeight), depthLeft, LowerBound, e.Ply)
						}
						quietMoves := e.triedQuietMoves[searchHeight][:legalQuietMove]
						e.searchHistory.AddHistory(move, currentMove, gpMove, histDepth, searchHeight, quietMoves)
					}
					return score
				}
				// Potential PV move, lets copy it to the current pv-line
				e.innerLines[searchHeight].ReplaceLine(move, e.innerLines[searchHeight+1])
				bestscore = score
				hashmove = move
			}

			if eval-score < 30 && depthLeft > 7 {
				rangeReduction += 1
			}
		}

		// if !isRootNode && e.startDepth > 6 {
		// 	e.parent.mu.RLock()
		// 	if e.startDepth <= e.parent.depth {
		// 		e.startDepth = 0
		// 		e.parent.mu.RUnlock()
		// 		return -MAX_INT
		// 	}
		// 	e.parent.mu.RUnlock()
		// }
	}
	if ((e.isMainThread && !e.TimeManager().AbruptStop) || (!e.isMainThread && !e.parent.Stop)) && !firstLayerOfSingularity {
		if alpha > oldAlpha {
			TranspositionTable.Set(hash, hashmove, evalToTT(bestscore, searchHeight), depthLeft, Exact, e.Ply)
		} else {
			TranspositionTable.Set(hash, hashmove, evalToTT(bestscore, searchHeight), depthLeft, UpperBound, e.Ply)
		}
	}
	if e.isMainThread && isRootNode && legalMoves == 1 && len(e.MovesToSearch) != 1 {
		e.TimeManager().StopSearchNow = true
	}
	return bestscore
}
