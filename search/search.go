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

var RazoringMargin int16 = 253
var TPMargin int16 = 144
var RFPMargin int16 = 53
var FPMargin int16 = 148
var RangeReductionMargin int16 = 44
var DeltaMargin int16 = 320
var LMRCaptureMargin int16 = 154

func (r *Runner) Search(depth int8, mateIn int16, nodes int64) {
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
		e.Search(depth, mateIn, nodes)
	} else {
		r.ClearForSearch()
		var wg sync.WaitGroup
		for i := 0; i < len(r.Engines); i++ {
			wg.Add(1)
			go func(e *Engine, depth int8, i int) {
				e.ParallelSearch(depth, mateIn, nodes)
				wg.Done()
			}(r.Engines[i], depth, i)
		}
		wg.Wait()
		for e.TimeManager().Pondering {
			// busy waiting
		}
		r.SendBestMove()
	}
	TranspositionTable.AdvanceAge()
}

func (e *Engine) ParallelSearch(depth int8, mateIn int16, nodes int64) {
	e.ClearForSearch()
	e.rootSearch(depth, mateIn, nodes)
}

func (e *Engine) Search(depth int8, mateIn int16, nodes int64) {
	e.parent.ClearForSearch()
	e.ClearForSearch()
	e.rootSearch(depth, mateIn, nodes)
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

var quietLmrReductions [32][32]int = initLMR(true)
var noisyLmrReductions [32][32]int = initLMR(false)

// This idea is taken from Weiss, which I believe in turn is taken from many open source
// engines.
func initLMR(isQuiet bool) [32][32]int {
	var reductions [32][32]int
	for depth := 1; depth < 32; depth++ {
		for moves := 1; moves < 32; moves++ {
			if isQuiet {
				reductions[depth][moves] = int(0.8 + math.Log(float64(depth))*math.Log(1.2*float64(moves))/2.5)
			} else {
				reductions[depth][moves] = int(math.Log(float64(depth)) * math.Log(1.2*float64(moves)) / 3.5)
			}
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

func (e *Engine) rootSearch(depth int8, mateIn int16, nodes int64) {
	pv := NewPVLine(MAX_DEPTH)

	lastDepth := int8(1)

	bookmove := GetBookMove(e.Position)
	if e.isMainThread && bookmove != EmptyMove {
		pv.Recycle()
		pv.AddFirst(bookmove)
		e.updatePv(pv, 0, 1, true)
	} else {
		for iterationDepth := int8(1); iterationDepth <= depth; iterationDepth += 1 {
			if e.isMainThread && nodes > 0 && nodes <= e.nodesVisited {
				break
			}

			if e.isMainThread {
				if iterationDepth > 1 && !e.TimeManager().CanStartNewIteration() {
					break
				}
			} else if iterationDepth > 1 && e.parent.Stop {
				break
			}

			e.startDepth = iterationDepth
			e.aspirationWindow(iterationDepth, mateIn != -2)
			newScore := e.Scores[0]

			if (e.isMainThread && e.TimeManager().AbruptStop) || (!e.isMainThread && e.parent.Stop) {
				break
			}

			if e.isMainThread && iterationDepth >= 8 && e.score-newScore >= 30 { // Position degrading
				e.TimeManager().ExtraTime()
			}

			lastDepth = iterationDepth
			e.pred.Clear()
			e.score = newScore
			e.rootMove = e.MultiPVs[0].MoveAt(0)
			if e.isMainThread && !e.MultiPVs[0].IsEmpty() {
				pv.Clone(e.MultiPVs[0])
				if e.MultiPV > 1 {
					e.SendMultiPv(pv, e.score, lastDepth)
				} else {
					e.SendPv(pv, e.score, iterationDepth)
				}
			}

			if e.isMainThread && mateIn != -1 && (newScore == -CHECKMATE_EVAL+mateIn || newScore == CHECKMATE_EVAL-mateIn ||
				newScore == -CHECKMATE_EVAL+mateIn-1 || newScore == CHECKMATE_EVAL-mateIn+1) {
				break
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
		if e.MultiPV > 1 {
			e.SendMultiPv(pv, e.score, lastDepth)
		} else {
			e.SendPv(pv, e.score, lastDepth)
		}
	}
}

func (e *Engine) aspirationWindow(iterationDepth int8, mateFinderMode bool) {
	e.doPruning = iterationDepth > 3 && !mateFinderMode

	var initialWindow int16 = 12
	var delta int16 = 16
	var alpha, beta, score int16
	originalDepth := iterationDepth
	maxSeldepth := int8(0)
	e.seldepth = 0
	lsm := len(e.MovesToSearch)

	for i := 0; i < e.MultiPV && (lsm == 0 || i < lsm); i++ {
		e.CurrentPV = i
		firstIteration := true
		e.NoMoves = false
		for !e.NoMoves {
			e.innerLines[0].Recycle()
			if firstIteration {
				score = e.Scores[i]
				alpha = max16(score-initialWindow, -MAX_INT)
				beta = min16(score+initialWindow, MAX_INT)
			}
			beta = min16(beta, MAX_INT)
			alpha = max16(alpha, -MAX_INT)

			if originalDepth <= 6 {
				beta = MAX_INT
				alpha = -MAX_INT
				delta = beta
			}

			score = e.alphaBeta(iterationDepth, 0, alpha, beta)
			if /* e.startDepth == 0 || */ e.TimeManager().AbruptStop || e.parent.Stop {
				e.seldepth = max8(e.seldepth, maxSeldepth)
				e.Scores[i] = -MAX_INT
				goto sortPVs
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
				e.Scores[i] = score
				e.MultiPVs[i].Clone(e.innerLines[0])
				break
			}
			delta += delta * 2 / 3
			maxSeldepth = max8(e.seldepth, maxSeldepth)
			firstIteration = false
		}
	}

sortPVs:
	for i := 0; i < e.MultiPV; i++ {
		for j := i + 1; j < e.MultiPV; j++ {
			if e.Scores[i] < e.Scores[j] {
				e.Scores[i], e.Scores[j] = e.Scores[j], e.Scores[i]
				e.MultiPVs[i], e.MultiPVs[j] = e.MultiPVs[j], e.MultiPVs[i]
			}
		}
	}
	e.seldepth = max8(e.seldepth, maxSeldepth)
	// We should never get here
}

func (e *Engine) alphaBeta(depthLeft int8, searchHeight int8, alpha int16, beta int16) int16 {
	e.VisitNode(searchHeight)

	isRootNode := searchHeight == 0
	isPvNode := alpha != beta-1

	position := e.Position
	e.tt.Prefetch(position.Hash())

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
		depthLeft += 1 // Check Extension
	}

	if depthLeft <= 0 {
		e.staticEvals[searchHeight] = position.Evaluate()
		return e.quiescence(alpha, beta, searchHeight)
	}

	firstLayerOfSingularity := e.skipHeight == searchHeight && e.skipMove != EmptyMove
	hash := position.Hash()
	nHashMove, nEval, nDepth, nType, ttHit := e.tt.Get(hash)
	if ttHit {
		ttHit = position.IsPseudoLegal(nHashMove)
		nEval = evalFromTT(nEval, searchHeight)
	}
	if !ttHit {
		if isRootNode {
			nHashMove = e.rootMove
		} else {
			nHashMove = EmptyMove
		}
	}

	if !isPvNode && ttHit && nDepth >= depthLeft && !firstLayerOfSingularity {
		if nEval >= beta && nType == LowerBound {
			e.CacheHit()
			e.searchHistory.AddHistory(nHashMove, currentMove, gpMove, depthLeft, searchHeight, nil, nil)
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
			e.tbHit += 1

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
				e.tt.Set(hash, EmptyMove, score, depthLeft, flag)
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
	if /* !isPvNode && */ !ttHit && depthLeft >= 5 {
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

	var eval int16 = position.Evaluate() //-MAX_INT
	// if !isRootNode && currentMove == EmptyMove {
	// 	eval = -1 * e.staticEvals[searchHeight-1]
	// } else {
	// eval =
	// }

	e.staticEvals[searchHeight] = eval
	improving := currentMove == EmptyMove ||
		(searchHeight > 2 && e.staticEvals[searchHeight] > e.staticEvals[searchHeight-2])

	e.searchHistory.ResetKillers(searchHeight + 1)
	// Pruning
	pruningAllowed := !isPvNode && !isInCheck && e.doPruning && !firstLayerOfSingularity

	// Idea taken from Berserk
	histDepth := depthLeft
	if eval+p > beta {
		histDepth += 1
	}

	if pruningAllowed {
		if ttHit && ((nEval > eval && nType == LowerBound) ||
			(nEval < eval && nType == UpperBound)) {
			eval = nEval
		}

		// Razoring
		if depthLeft < 2 && eval+RazoringMargin <= alpha {
			newEval := e.quiescence(alpha, beta, searchHeight)
			return newEval
		}

		// Reverse Futility Pruning
		reverseFutilityMargin := int16(depthLeft) * RFPMargin //(b - p)
		if improving {
			reverseFutilityMargin -= RFPMargin // int16(depthLeft) * p
		}
		if depthLeft < 9 && eval-reverseFutilityMargin >= beta {
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
				return score
			}
			// }
		}

		// Threat pruning, idea from Koivisto
		// if improving {
		// 	tpMargin += 30
		// }
		if depthLeft == 1 && eval > beta+TPMargin && (!position.Board.HasThreats(position.Turn().Other()) || position.Board.HasThreats(position.Turn())) {
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
		// 				return score
		// 			}
		// 		}
		// 		i += 1
		// 	}
		// }
	}

	// Internal Iterative Deepening
	// if /* isPvNode && */ depthLeft >= 8 && !ttHit {
	// 	e.innerLines[searchHeight].Recycle()
	// 	score := e.alphaBeta(depthLeft-7, searchHeight, alpha, beta)
	// 	if e.isMainThread && e.TimeManager().AbruptStop {
	// 		return score
	// 	} else if !e.isMainThread && e.parent.Stop {
	// 		return score
	// 	}
	// 	line := e.innerLines[searchHeight]
	// 	if line.moveCount != 0 {
	// 		nHashMove = e.innerLines[searchHeight].MoveAt(0)
	// 	}
	// 	e.innerLines[searchHeight].Recycle()
	// }
	//
	movePicker := e.MovePickers[searchHeight]
	movePicker.RecycleWith(position, e, searchHeight, nHashMove, depthLeft, false)
	oldAlpha := alpha

	// using fail soft with negamax:
	var bestscore int16
	var hashmove Move
	legalMoves := 0
	quietMoves := -1
	legalQuietMoves := 0
	legalNoisyMoves := 0
	noisyMoves := -1
	searchedAMove := false
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
		searchedAMove = true
		if oldEnPassant, oldTag, hc, ok := position.MakeMove(hashmove); ok {
			legalMoves += 1
			if isQuiet {
				legalQuietMoves += 1
			} else {
				legalNoisyMoves += 1
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
					extension += 1
				} else {

					// Multi-Cut, at least 2 moves beat beta, idea is taken from Stockfish
					if pruningAllowed {
						if threshold >= beta {
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
			e.NoteMove(hashmove, legalQuietMoves-1, legalNoisyMoves-1, searchHeight)
			bestscore = -e.alphaBeta(depthLeft-1+extension, searchHeight+1, -beta, -alpha)
			e.pred.Pop()
			position.UnMakeMove(hashmove, oldTag, oldEnPassant, hc)
			if bestscore > alpha {
				if bestscore >= beta {
					if (e.isMainThread && !e.TimeManager().AbruptStop) || (!e.isMainThread && !e.parent.Stop) {
						if !firstLayerOfSingularity {
							e.tt.Set(hash, hashmove, evalToTT(bestscore, searchHeight), depthLeft, LowerBound)
						}
						quietMoves := e.triedQuietMoves[searchHeight][:legalQuietMoves]
						noisyMoves := e.triedNoisyMoves[searchHeight][:legalNoisyMoves]
						e.searchHistory.AddHistory(hashmove, currentMove, gpMove, histDepth, searchHeight, quietMoves, noisyMoves)
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
	d := int(depthLeft)
	pruningThreashold := 5 + d*d
	if !improving && !isPvNode {
		pruningThreashold = pruningThreashold/2 - 1
	}

	lmrThreashold := 2
	if isPvNode {
		lmrThreashold += 1
	}
	fpMargin := eval + FPMargin*int16(depthLeft)
	rangeReduction := 0
	if eval-bestscore < RangeReductionMargin && depthLeft > 7 {
		rangeReduction += 1
	}

	// seeScores := movePicker.captureMoveList.Scores
	// quietScores := movePicker.quietMoveList.Scores
	// var historyThreashold int32 = int32(depthLeft) * -1024
	var move Move
	var seeScore int16
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
				continue
			}

			// Late Move Pruning
			if isQuiet && depthLeft < 8 &&
				legalMoves+1 > pruningThreashold && !isKiller && abs16(alpha) < WIN_IN_MAX {
				// This is a hack really, mp.Next() won't return any quiets, and I am hacking this
				// to avoid returning quiets, after the first LMP cut
				movePicker.isQuiescence = true
				continue // LMP
			}

			// SEE pruning
			seeBound := -50 * int16(depthLeft)
			if isCaptureMove && depthLeft < 7 && eval <= alpha && abs16(alpha) < WIN_IN_MAX {
				seeScore = int16(movePicker.captureSees[noisyMoves])
				if seeScore < seeBound {
					break
				}
			}

			if isQuiet && depthLeft < 7 && !isKiller {
				seeScore = position.Board.SeeGe(move.Destination(), move.CapturedPiece(), move.Source(), move.MovingPiece(), seeBound)
				if seeScore < seeBound {
					continue // Quiet SEE
				}

			}
			//
			// // History pruning
			// lmrDepth := depthLeft - int8(lmrReductions[min8(31, depthLeft)][min(31, legalMoves+1)])
			// if isQuiet && quietScores[quietMoves] < historyThreashold && lmrDepth < 3 && legalMoves+1 > lmrThreashold {
			// 	continue
			// }
		}

		if oldEnPassant, oldTag, hc, ok := position.MakeMove(move); ok {
			legalMoves += 1
			if isQuiet {
				legalQuietMoves += 1
			} else {
				legalNoisyMoves += 1
			}

			if e.isMainThread && e.parent.DebugMode && isRootNode {
				fmt.Printf("info depth %d currmove %s currmovenumber %d\n", depthLeft, move.ToString(), legalMoves)
			}

			e.NoteMove(move, legalQuietMoves-1, legalNoisyMoves-1, searchHeight)
			LMR := int8(0)

			// Late Move Reduction
			if e.doPruning && (isQuiet || isCaptureMove && seeScore < 0) && depthLeft > 2 && legalMoves > lmrThreashold {
				if isQuiet {
					LMR = int8(quietLmrReductions[min8(31, depthLeft)][min(31, legalMoves)])
					LMR -= int8(e.searchHistory.QuietHistory(gpMove, currentMove, move) / 10649) //12288)
				} else {
					LMR = int8(noisyLmrReductions[min8(31, depthLeft)][min(31, legalMoves)])
					if eval+move.CapturedPiece().Weight()+LMRCaptureMargin < beta {
						LMR += 1
					}
				}

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

				LMR = min8(depthLeft-2, max8(LMR, 1))
			}

			e.pred.Push(position.Hash())
			e.innerLines[searchHeight+1].Recycle()
			e.positionMoves[searchHeight+1] = move
			score := -e.alphaBeta(depthLeft-1-LMR, searchHeight+1, -alpha-1, -alpha)
			e.pred.Pop()
			if score > alpha && score < beta {
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
							e.tt.Set(hash, move, evalToTT(score, searchHeight), depthLeft, LowerBound)
						}
						quietMoves := e.triedQuietMoves[searchHeight][:legalQuietMoves]
						noisyMoves := e.triedNoisyMoves[searchHeight][:legalNoisyMoves]
						e.searchHistory.AddHistory(move, currentMove, gpMove, histDepth, searchHeight, quietMoves, noisyMoves)
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
			e.tt.Set(hash, hashmove, evalToTT(bestscore, searchHeight), depthLeft, Exact)
		} else {
			e.tt.Set(hash, hashmove, evalToTT(bestscore, searchHeight), depthLeft, UpperBound)
		}
	}
	if e.isMainThread && isRootNode && legalMoves == 1 && len(e.MovesToSearch) == 0 && e.MultiPV == 1 {
		e.TimeManager().StopSearchNow = true
	} else if e.isMainThread && isRootNode && !searchedAMove {
		e.NoMoves = true
	}
	return bestscore
}
