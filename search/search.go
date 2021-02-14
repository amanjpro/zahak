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
var pv = make([]*Move, 100)

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
	return bestEval
}

func startMinimax(position *Position, depth int8, ply uint16) (*Move, int) {

	// Collect evaluation for moves per iteration to help us order moves for the next iteration
	legalMoves := position.LegalMoves()
	iterationEvals := make([]int, len(legalMoves))

	var bestMove *Move
	var previousBestMove *Move

	fruitlessIterations := 0
	bestScore := -MAX_INT

	timeForSearch := 2 * time.Minute // TODO: with time management this should go

	alpha := -MAX_INT
	beta := MAX_INT

	start := time.Now()
	for iterationDepth := int8(1); iterationDepth <= depth; iterationDepth++ {
		orderedMoves := orderIterationMoves(&IterationMoves{legalMoves, iterationEvals})
		bestScore = -MAX_INT
		for index, move := range orderedMoves {
			if time.Now().Sub(start) > timeForSearch {
				if index != 0 {
					return previousBestMove, bestScore
				} else {
					return bestMove, bestScore
				}
			}
			fmt.Printf("info currmove %s currmovenumber %d\n\n", move.ToString(), index+1)
			sendPv := false
			cp, ep, tg := position.MakeMove(move)
			score := withAspirationWindow(position, iterationDepth, alpha, beta, ply)
			iterationEvals[index] = score
			position.UnMakeMove(move, tg, ep, cp)
			if score > bestScore {
				sendPv = true
				bestScore = score
				bestMove = move
			}
			if score > beta {
				alpha = score
			}
			if score == CHECKMATE_EVAL {
				return move, score
			}
			timeSpent := time.Now().Sub(start)
			if sendPv {
				fmt.Printf("info depth %d nps %d tbhits %d nodes %d score cp %d time %d pv %s",
					iterationDepth, nodesVisited/1000*int64(timeSpent.Seconds()),
					cacheHits, nodesVisited, int(bestScore/100), timeSpent.Milliseconds(), bestMove.ToString())
				for i, move := range pv {
					if move == nil {
						break
					}
					fmt.Printf(" %s", move.ToString())
					pv[i] = nil
				}
			} else {
				fmt.Printf("info depth %d nps %d tbhits %d nodes %d score cp %d time %d",
					iterationDepth, nodesVisited/1000*int64(timeSpent.Seconds()),
					cacheHits, nodesVisited, int(bestScore/100), timeSpent.Milliseconds())
			}
			fmt.Printf("\n\n")
		}
		if iterationDepth > 3 && *previousBestMove == *bestMove {
			fruitlessIterations++
			if fruitlessIterations > 3 {
				return bestMove, bestScore
			}
		} else {
			fruitlessIterations = 0
		}
		previousBestMove = bestMove
	}
	return bestMove, bestScore
}

func withAspirationWindow(position *Position, depth int8, alpha int, beta int, ply uint16) int {

	wp := WhitePawn
	aspirationWindow := wp.Weight() / 25
	alpha -= aspirationWindow
	alpha += aspirationWindow
	for trials := 1; trials <= 3; trials++ {
		score := minimax(position, depth, 0, false, alpha, beta, ply)
		if score <= alpha {
			alpha -= aspirationWindow * trials
		} else if score >= beta {
			beta += aspirationWindow * trials
		} else {
			return score
		}
	}
	return minimax(position, depth, 0, false, -MAX_INT, MAX_INT, ply)
}

func minimax(position *Position, depthLeft int8, searchHeight int8,
	isMaximizingPlayer bool, alpha int, beta int, ply uint16) int {
	nodesVisited += 1

	if position.Status() == Checkmate {
		if isMaximizingPlayer {
			return -CHECKMATE_EVAL
		}
		return CHECKMATE_EVAL
	} else if position.Status() == Draw {
		return 0
	}

	if depthLeft == 0 || STOP_SEARCH_GLOBALLY {
		return Evaluate(position)
		// return quescence(position, isMaximizingPlayer, alpha, beta, 0)
	}

	nodesSearched += 1
	legalMoves := position.LegalMoves()
	orderedMoves := orderMoves(&ValidMoves{position, legalMoves})
	var bestMove *Move

	for _, move := range orderedMoves {
		score := getEval(position, depthLeft, searchHeight, isMaximizingPlayer, alpha, beta, move, ply)
		if isMaximizingPlayer {
			if score >= beta {
				return beta
			}
			if score > alpha {
				bestMove = move
				alpha = score
			}
		} else {
			if score <= alpha {
				return alpha
			}
			if score < beta {
				bestMove = move
				beta = score
			}
		}
	}
	if isMaximizingPlayer {
		pv[searchHeight] = bestMove
		return alpha
	} else {
		pv[searchHeight] = bestMove
		return beta
	}
}

func getEval(position *Position, depthLeft int8, searchHeight int8, isMaximizingPlayer bool,
	alpha int, beta int, move *Move, ply uint16) int {
	var score int
	capturedPiece, oldEnPassant, oldTag := position.MakeMove(move)
	newPositionHash := position.Hash()
	cachedEval, found := TranspositionTable.Get(newPositionHash)
	if found &&
		(cachedEval.Eval == CHECKMATE_EVAL ||
			cachedEval.Eval == -CHECKMATE_EVAL ||
			cachedEval.Depth >= depthLeft) {
		cacheHits += 1
		score = cachedEval.Eval
	} else {
		v := minimax(position, depthLeft-1, searchHeight+1, !isMaximizingPlayer, alpha, beta, ply)
		var tpe NodeType
		if isMaximizingPlayer {
			if v >= beta {
				tpe = LowerBound
			}
			if score > alpha {
				tpe = Exact
			}
		} else {
			if v <= alpha {
				tpe = UpperBound
			}
			if v < beta {
				tpe = Exact
			}
		}
		TranspositionTable.Set(newPositionHash, &CachedEval{position.Hash(), v, depthLeft, tpe, ply})
		score = v
	}
	position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece)
	return score
}
