package search

import (
	"fmt"
	"time"

	. "github.com/amanjpro/zahak/cache"
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

var STOP_SEARCH_GLOBALLY = false

var nodesVisited int64 = 0
var nodesSearched int64 = 0
var cacheHits int64 = 0
var pv = NewPVLine(100)

type EvalMove struct {
	eval int
	move *Move
}

func (e *EvalMove) Move() *Move {
	return e.move
}

func (e *EvalMove) Eval() int {
	return e.eval
}

func Search(position *Position, depth int8, ply uint16) EvalMove {
	STOP_SEARCH_GLOBALLY = false
	nodesVisited = 0
	nodesSearched = 0
	cacheHits = 0
	var bestEval EvalMove
	start := time.Now()
	bestMove, score := startMinimax(position, depth, ply)
	bestEval = EvalMove{score, bestMove}
	end := time.Now()
	fmt.Printf("Visited: %d, Selected: %d, Cache-hit: %d\n\n", nodesVisited, nodesSearched, cacheHits)
	fmt.Printf("Took %f seconds\n\n", end.Sub(start).Seconds())
	pv.Pop() // pop our move
	pv.Pop() // pop our opponent's move
	return bestEval
}

func startMinimax(position *Position, depth int8, ply uint16) (*Move, int) {

	// Collect evaluation for moves per iteration to help us order moves for the next iteration
	legalMoves := position.LegalMoves()
	iterationEvals := make([]int, len(legalMoves))

	var bestMove *Move
	// var previousBestMove *Move

	bestScore := -MAX_INT

	timeForSearch := 3 * time.Minute // TODO: with time management this should go
	// fruitlessIterations := 0

	alpha := -MAX_INT
	beta := MAX_INT

	// wp := WhitePawn
	// aspirationWindow := wp.Weight() / 4
	start := time.Now()
	for iterationDepth := int8(1); iterationDepth <= depth; iterationDepth++ {
		currentBestScore := -MAX_INT
		orderedMoves := orderIterationMoves(&IterationMoves{legalMoves, iterationEvals})
		line := NewPVLine(iterationDepth + 1)
		searchPv := true
		for index, move := range orderedMoves {
			if time.Now().Sub(start) > timeForSearch {
				if index != 0 {
					return bestMove, currentBestScore
				} else {
					return bestMove, bestScore
				}
			}
			fmt.Printf("info currmove %s currmovenumber %d\n\n", move.ToString(), index+1)
			sendPv := false
			cp, ep, tg := position.MakeMove(move)
			// score, newAlpha, newBeta, set := withAspirationWindow(position, iterationDepth, alpha, beta, ply, aspirationWindow, line)
			// if set && iterationDepth >= 4 {
			// 	alpha = newAlpha
			// 	beta = newBeta
			// }
			// score := -MAX_INT
			// if iterationDepth == 1 {
			// 	score = -alphaBeta(position, iterationDepth, 1, -beta, -alpha, ply, line)
			// } else {

			score := -MAX_INT
			// if index == 0 { // First move of the iteration? establish the aspiration window
			// 	   score = -pvSearch(position, iterationDepth, 1, -beta, -alfa, ply, line)
			// 		 bestScore = score
			// 	 if( bestScore > alfa ) {
			// 			if( bestscore >= beta )
			// 				 return bestscore;
			// 			alfa = bestscore;
			// 	 }
			//
			// }
			if searchPv {
				score = -alphaBetaPVS(position, iterationDepth, 1, -beta, -alpha, ply, line)
			} else {
				score = -zeroWindowSearch(position, iterationDepth, 1, -alpha, ply, 0)
				if score > alpha { // in fail-soft ... && score < beta ) is common
					score = -alphaBetaPVS(position, iterationDepth, 1, -beta, -alpha, ply, line) // re-search
				}
			}
			// }
			// This only works, because checkmate eval is clearly distinguished from
			// maximum/minimum beta/alpha
			if score > alpha && score < beta { // no very hard alpha-beta cutoff
				iterationEvals[index] = score
				alpha = score
				if score > currentBestScore {
					currentBestScore = score
					sendPv = true
					pv.AddFirst(move)
					pv.ReplaceLine(line)
					bestMove = move
					bestScore = currentBestScore
					searchPv = false

				}
			} else {
				iterationEvals[index] = -MAX_INT // if it is, then too bad, that is a bad move
			}
			position.UnMakeMove(move, tg, ep, cp)

			if score == CHECKMATE_EVAL {
				return move, score
			}
			timeSpent := time.Now().Sub(start)
			if sendPv {
				fmt.Printf("info depth %d nps %d tbhits %d nodes %d score cp %d time %d pv %s\n\n",
					iterationDepth, nodesVisited/1000*int64(timeSpent.Seconds()),
					cacheHits, nodesVisited, currentBestScore, timeSpent.Milliseconds(), pv.ToString())
				// } else {
				// 	fmt.Printf("info depth %d nps %d tbhits %d nodes %d score cp %d time %d",
				// 		iterationDepth, nodesVisited/1000*int64(timeSpent.Seconds()),
				// 		cacheHits, nodesVisited, currentBestScore, timeSpent.Milliseconds())
			}
		}
		// if iterationDepth >= 5 && *previousBestMove == *bestMove {
		// 	if fruitlessIterations <= 3 {
		// 		fruitlessIterations++
		// 	} else {
		// 		break
		// 	}
		// } else {
		// 	fruitlessIterations = 0
		// }
		// previousBestMove = bestMove

		timeSpent := time.Now().Sub(start)
		fmt.Printf("info depth %d nps %d tbhits %d nodes %d score cp %d time %d pv %s\n\n",
			iterationDepth, nodesVisited/1000*int64(timeSpent.Seconds()),
			cacheHits, nodesVisited, currentBestScore, timeSpent.Milliseconds(), pv.ToString())
		alpha = -MAX_INT
		beta = MAX_INT
		currentBestScore = -MAX_INT
	}

	return bestMove, bestScore
}

func withAspirationWindow(position *Position, depth int8, alpha int, beta int, ply uint16, window int, pvline *PVLine) (int, int, int, bool) {

	for trials := 1; trials <= 3; trials++ {
		score := -alphaBeta(position, depth, 1, -beta, -alpha, ply, pvline)
		currentWindow := trials * window
		if score <= alpha {
			alpha -= currentWindow
		} else if score >= beta {
			beta += currentWindow
		} else {
			return score, alpha, beta, true
		}
	}
	return -alphaBeta(position, depth, 1, MAX_INT, -MAX_INT, ply, pvline), -MAX_INT, MAX_INT, false
}

func alphaBeta(position *Position, depthLeft int8, searchHeight int8, alpha int, beta int, ply uint16, pvline *PVLine) int {
	nodesVisited += 1
	outcome := position.Status()
	if outcome == Checkmate {
		return -CHECKMATE_EVAL
	} else if outcome == Draw {
		return 0
	}

	if depthLeft == 0 || STOP_SEARCH_GLOBALLY {
		// if searchHeight >= 4 {
		// 	return Evaluate(position)
		// }
		return quiescence(position, alpha, beta, 0, Evaluate(position))
	}

	nodesSearched += 1

	legalMoves := position.LegalMoves()
	orderedMoves := orderMoves(&ValidMoves{position, legalMoves, searchHeight + 1})

	hash := position.Hash()
	cachedEval, found := TranspositionTable.Get(hash)
	if found &&
		(cachedEval.Eval == CHECKMATE_EVAL ||
			cachedEval.Eval == -CHECKMATE_EVAL ||
			cachedEval.Depth >= depthLeft) {
		cacheHits += 1
		score := cachedEval.Eval
		if score == CHECKMATE_EVAL || score == -CHECKMATE_EVAL {
			return cachedEval.Eval
		}
		if score >= beta && (cachedEval.Type != LowerBound || cachedEval.Type == Exact) {
			return beta
		}
		if score <= alpha && (cachedEval.Type != UpperBound || cachedEval.Type == Exact) {
			return alpha
		}
		if cachedEval.Type == Exact && score < beta && score > alpha {
			return score
		}
	}

	foundExact := false
	for _, move := range orderedMoves {
		capturedPiece, oldEnPassant, oldTag := position.MakeMove(move)
		line := NewPVLine(depthLeft - 1)
		score := -alphaBeta(position, depthLeft-1, searchHeight+1, -beta, -alpha, ply, line)
		position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece)
		if score >= beta {
			TranspositionTable.Set(hash, &CachedEval{hash, score, depthLeft, UpperBound, ply})
			return score
		}
		if score > alpha {
			alpha = score
			TranspositionTable.Set(hash, &CachedEval{hash, score, depthLeft, LowerBound, ply})
			pvline.AddFirst(move)
			pvline.ReplaceLine(line)
			// Potential PV move, lets copy it to the current pv-line
			foundExact = true
		}
	}
	if foundExact {
		TranspositionTable.Set(hash, &CachedEval{hash, alpha, depthLeft, Exact, ply})
	}
	return alpha
}

func alphaBetaPVS(position *Position, depthLeft int8, searchHeight int8, alpha int, beta int, ply uint16, pvline *PVLine) int {
	nodesVisited += 1
	outcome := position.Status()
	if outcome == Checkmate {
		return -CHECKMATE_EVAL
	} else if outcome == Draw {
		return 0
	}

	if depthLeft == 0 || STOP_SEARCH_GLOBALLY {
		// if searchHeight >= 4 {
		// return Evaluate(position)
		// }
		return quiescence(position, alpha, beta, 0, Evaluate(position))
	}

	nodesSearched += 1

	legalMoves := position.LegalMoves()
	orderedMoves := orderMoves(&ValidMoves{position, legalMoves, searchHeight + 1})

	hash := position.Hash()
	cachedEval, found := TranspositionTable.Get(hash)
	if found &&
		(cachedEval.Eval == CHECKMATE_EVAL ||
			cachedEval.Eval == -CHECKMATE_EVAL ||
			cachedEval.Depth >= depthLeft) {
		cacheHits += 1
		score := cachedEval.Eval
		if score == CHECKMATE_EVAL || score == -CHECKMATE_EVAL {
			return cachedEval.Eval
		}
		if score >= beta && (cachedEval.Type != UpperBound || cachedEval.Type == Exact) {
			return beta
		}
		if score <= alpha && (cachedEval.Type != LowerBound || cachedEval.Type == Exact) {
			return alpha
		}
		// if cachedEval.Type == Exact {
		// 	return score
		// }
	}

	searchPv := true

	foundExact := false
	for _, move := range orderedMoves {
		capturedPiece, oldEnPassant, oldTag := position.MakeMove(move)
		line := NewPVLine(depthLeft - 1)
		score := -MAX_INT
		if searchPv {
			score = -alphaBetaPVS(position, depthLeft-1, searchHeight+1, -beta, -alpha, ply, line)
		} else {
			score = -zeroWindowSearch(position, depthLeft-1, searchHeight+1, -alpha, ply, 0)
			if score > alpha { // in fail-soft ... && score < beta ) is common
				score = -alphaBetaPVS(position, depthLeft-1, searchHeight+1, -beta, -alpha, ply, line) // re-search
			}
		}
		position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece)
		if score >= beta {
			TranspositionTable.Set(hash, &CachedEval{hash, score, depthLeft, LowerBound, ply})
			return beta
		}
		if score > alpha {
			alpha = score
			TranspositionTable.Set(hash, &CachedEval{hash, score, depthLeft, Exact, ply})
			// Potential PV move, lets copy it to the current pv-line
			pvline.AddFirst(move)
			pvline.ReplaceLine(line)
			foundExact = true
			searchPv = false
		}
	}
	if !foundExact {
		TranspositionTable.Set(hash, &CachedEval{hash, alpha, depthLeft, UpperBound, ply})
	}
	return alpha
}

func zeroWindowSearch(position *Position, depthLeft int8, searchHeight int8, beta int, ply uint16, inNullMoveSearch int) int {
	nodesVisited += 1
	// outcome := position.Status()
	// if outcome == Checkmate {
	// 	return -CHECKMATE_EVAL
	// } else if outcome == Draw {
	// 	return 0
	// }

	if depthLeft == 0 || STOP_SEARCH_GLOBALLY {
		// return quiescence(p, beta-1, beta, 0)
		// if searchHeight >= 4 {
		// return Evaluate(position)
		// }
		return quiescence(position, beta-1, beta, 0, Evaluate(position))
	}

	nodesSearched += 1

	legalMoves := position.LegalMoves()
	orderedMoves := orderMoves(&ValidMoves{position, legalMoves, searchHeight + 1})

	hash := position.Hash()
	cachedEval, found := TranspositionTable.Get(hash)
	if found &&
		// (cachedEval.Eval == CHECKMATE_EVAL ||
		// 	cachedEval.Eval == -CHECKMATE_EVAL ||
		cachedEval.Depth >= depthLeft {
		cacheHits += 1
		score := cachedEval.Eval
		// return cachedEval.Eval
		// }
		// if cachedEval.Type != Exact {
		// 	return score
		// }
		if score >= beta && (cachedEval.Type != LowerBound || cachedEval.Type == Exact) {
			return beta
		}
		if score <= beta-1 && (cachedEval.Type != UpperBound || cachedEval.Type == Exact) {
			return beta - 1
		}
		// if cachedEval.Type == Exact {
		// 	return score
		// }
	}

	isNullMoveAllowed := true // Null-Move pruning is always activated for now
	R := int8(3)
	if searchHeight > 6 {
		R = 2
	}

	for _, move := range orderedMoves {
		if isNullMoveAllowed && depthLeft >= 5 {
			bound := beta
			if inNullMoveSearch == 0 {
				tempo := 20    // TODO: Make it variable with a formula like: 10*(numPGAM > 0) + 10* numPGAM > 15);
				bound -= tempo // variable bound
			}
			position.NullMove()
			inNullMoveSearch++
			score := -zeroWindowSearch(position, depthLeft-R-1, searchHeight+1, 1-bound, ply, inNullMoveSearch)
			position.NullMove()
			inNullMoveSearch--
			if score >= bound {
				return beta // null move pruning
			}
		}
		capturedPiece, oldEnPassant, oldTag := position.MakeMove(move)
		score := -zeroWindowSearch(position, depthLeft-1, searchHeight+1, 1-beta, ply, inNullMoveSearch)
		position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece)
		if score >= beta {
			return beta // fail-hard beta-cutoff
		}
	}
	return beta - 1 // fail-hard, return alpha
}
