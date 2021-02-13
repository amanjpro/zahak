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
	var isMaximizingPlayer = position.Turn() == White
	// validMoves := position.LegalMoves()
	// evals := make(chan EvalMove)
	// evalIsSet := false
	start := time.Now()
	// for i := 0; i < len(validMoves); i++ {
	// 	p := position.copy()
	// 	move := validMoves[i]
	// 	p.MakeMove(move)
	// 	go parallelMinimax(p, move, depth, !isMaximizingPlayer, evals)
	// }
	// for i := 0; i < len(validMoves); i++ {
	// 	evalMove := <-evals
	//
	// 	mvStr := evalMove.move.ToString()
	// 	fmt.Printf("info nodes %d score cp %d currmove %s pv %s",
	// 		nodesVisited, int(evalMove.eval*100*dir), mvStr, mvStr)
	// 	for _, mv := range evalMove.line {
	// 		fmt.Printf(" %s", mv.ToString())
	// 	}
	// 	fmt.Print("\n\n")
	// 	if isMaximizingPlayer {
	// 		if !evalIsSet || evalMove.eval > bestEval.eval {
	// 			bestEval = evalMove
	// 			evalIsSet = true
	// 		}
	// 	} else {
	// 		if !evalIsSet || evalMove.eval < bestEval.eval {
	// 			bestEval = evalMove
	// 			evalIsSet = true
	// 		}
	// 	}

	bestMove, score := startMinimax(position, depth, isMaximizingPlayer, ply)
	// fmt.Printf("info nodes %d score cp %d currmove %s",
	// nodesVisited, int(score*dir), bestMove.ToString())
	// for i := 1; i < len(moves); i++ {
	// 	mv := moves[i]
	// 	fmt.Printf(" %s", mv.ToString())
	// }
	fmt.Print("\n\n")

	bestEval = EvalMove{score, bestMove}
	// }
	end := time.Now()
	// close(evals)
	fmt.Printf("Visited: %d, Selected: %d, Cache-hit: %d\n\n", nodesVisited, nodesSearched, cacheHits)
	fmt.Printf("Took %f seconds\n\n", end.Sub(start).Seconds())
	return bestEval
}

func startMinimax(position *Position, depth int8,
	isMaximizingPlayer bool, ply uint16) (*Move, int) {
	legalMoves := position.LegalMoves()
	iterationEvals := make([]int, len(legalMoves))
	var bestMove *Move
	var bestScore int
	var dir = -1
	if isMaximizingPlayer {
		dir = 1
	}

	timeForSearch := 5 * time.Minute
	wp := WhitePawn
	aspirationWindow := wp.Weight()
	start := time.Now()
	alpha := -MAX_INT
	beta := MAX_INT
	iterAlpha := alpha
	iterBeta := beta

	for iterationDepth := int8(1); iterationDepth <= depth; iterationDepth++ {
		bestScore = -1 * dir * MAX_INT
		orderedMoves := orderIterationMoves(&IterationMoves{legalMoves, iterationEvals})
		for index, move := range orderedMoves {
			fmt.Printf("info currmove %s\n\n", move.ToString())
			sendPv := false
			if time.Now().Sub(start) > timeForSearch {
				return bestMove, bestScore
			}
			cp, ep, tg := position.MakeMove(move)
			score := minimax(position, iterationDepth, 0, !isMaximizingPlayer, iterAlpha, iterBeta, ply)
			trials := 0
			for trials <= 3 {
				if score <= iterAlpha {
					iterAlpha = score
					score = minimax(position, iterationDepth, 0, !isMaximizingPlayer, iterAlpha, iterBeta, ply)
				} else if score >= iterBeta {
					iterBeta = score
					score = minimax(position, iterationDepth, 0, !isMaximizingPlayer, iterAlpha, iterBeta, ply)
				} else {
					break
				}
				trials++
				if trials == 2 {
					iterAlpha = alpha
					iterBeta = beta
				}
			}
			iterationEvals[index] = score
			position.UnMakeMove(move, tg, ep, cp)
			if isMaximizingPlayer {
				if score > bestScore {
					sendPv = true
					bestScore = score
					bestMove = move
				}
				if score > iterAlpha {
					iterAlpha = score
				}
				if score == CHECKMATE_EVAL {
					return move, CHECKMATE_EVAL
				}
			} else {
				if score < bestScore {
					sendPv = true
					bestScore = score
					bestMove = move
				}
				if score < iterBeta {
					iterBeta = score
				}
				if score == -CHECKMATE_EVAL {
					return move, -CHECKMATE_EVAL
				}
			}
			timeSpent := time.Now().Sub(start)
			if sendPv {
				fmt.Printf("info depth %d nps %d tbhits %d nodes %d score cp %d time %d pv %s",
					iterationDepth, nodesVisited/1000*int64(timeSpent.Seconds()),
					cacheHits, nodesVisited, int(bestScore*dir/100), timeSpent.Milliseconds(), bestMove.ToString())
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
					cacheHits, nodesVisited, int(bestScore*dir/100), timeSpent.Milliseconds())
			}
			fmt.Printf("\n\n")
		}
		iterAlpha = bestScore - aspirationWindow
		iterBeta = bestScore + aspirationWindow
	}
	return bestMove, bestScore
}

func minimax(position *Position, depthLeft int8, searchHeight int8,
	isMaximizingPlayer bool, alpha int, beta int, ply uint16) int {
	nodesVisited += 1

	if depthLeft == 0 || STOP_SEARCH_GLOBALLY {
		evl := Evaluate(position)
		if !STOP_SEARCH_GLOBALLY {
			if isMaximizingPlayer {
				evl = quescence(position, isMaximizingPlayer, evl, beta, ply)
			} else {
				evl = quescence(position, isMaximizingPlayer, alpha, evl, ply)
			}
		}
		return evl
	}

	nodesSearched += 1
	legalMoves := position.LegalMoves()
	orderedMoves := orderMoves(&ValidMoves{position, legalMoves})
	var bestMove *Move

	for _, move := range orderedMoves {
		score := getEval(position, depthLeft-1, searchHeight+1, !isMaximizingPlayer, alpha, beta, move, ply)
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

	if len(orderedMoves) == 0 {
		return Evaluate(position)
	} else if isMaximizingPlayer {
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
		v := minimax(position, depthLeft, searchHeight, isMaximizingPlayer, alpha, beta, ply)
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
