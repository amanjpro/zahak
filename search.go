package main

import (
	"fmt"
	"math"
	"time"

	"github.com/notnil/chess"
	"github.com/patrickmn/go-cache"
)

var evalCache = cache.New(1*time.Hour, 10*time.Minute)

var nodesVisited int64 = 0
var nodesSearched int64 = 0
var cacheHits int64 = 0

type EvalMove struct {
	eval float64
	move *chess.Move
	line []chess.Move
}

func search(position *chess.Position, depth int8) *EvalMove {
	nodesVisited = 0
	nodesSearched = 0
	cacheHits = 0
	evalCache.Flush()
	var bestEval *EvalMove
	var isMaximizingPlayer = position.Turn() == chess.White
	validMoves := position.ValidMoves()
	evals := make(chan EvalMove)
	start := time.Now()
	for _, move := range validMoves {
		go parallelMinimax(position.Update(move), move, depth, true, evals)
	}
	for i := 0; i < len(validMoves); i++ {
		evalMove := <-evals

		fmt.Println("Tree")
		for _, mv := range evalMove.line {
			fmt.Printf("%s ", mv.String())
		}
		fmt.Println("Eval: ", evalMove.eval)
		if isMaximizingPlayer {
			if bestEval == nil || evalMove.eval > bestEval.eval {
				bestEval = &evalMove
			}
		} else {
			if bestEval == nil || evalMove.eval < bestEval.eval {
				bestEval = &evalMove
			}
		}
	}
	end := time.Now()
	close(evals)
	fmt.Printf("Visited: %d, Selected: %d, Cache-hit: %d\n", nodesVisited, nodesSearched, cacheHits)
	fmt.Printf("Took %d seconds\n", end.Sub(start).Seconds())
	return bestEval
}

func parallelMinimax(position *chess.Position, move *chess.Move, depth int8, isMaximizingPlayer bool, resultEval chan EvalMove) {
	eval, moves := minimax(position, depth, isMaximizingPlayer, math.Inf(-1), math.Inf(1), []chess.Move{})
	resultEval <- EvalMove{eval, move, moves}
}

func minimax(position *chess.Position, depth int8, isMaximizingPlayer bool, alpha float64, beta float64, line []chess.Move) (float64, []chess.Move) {

	nodesVisited += 1

	if position.Status() == chess.Checkmate {
		if isMaximizingPlayer {
			return math.Inf(-1), line
		}
		return math.Inf(1), line
	}

	if position.Status() == chess.Stalemate || position.Status() == chess.InsufficientMaterial {
		return 0.0, line
	}

	if depth == 0 {
		// TODO: Perform all captures before giving up, to avoid the horizon effect
		if isMaximizingPlayer {
			return eval(position), line
		}
		return -eval(position), line
	}

	nodesSearched += 1

	moves := position.ValidMoves()
	newLine := line
	if isMaximizingPlayer {
		for _, move := range moves {
			newPosition := position.Update(move)
			newPositionHash := fmt.Sprintf("%x", newPosition.Hash())
			cachedEval, found := evalCache.Get(newPositionHash)
			if found {
				cacheHits += 1
				v := cachedEval.(float64)
				if v > alpha {
					newLine = append(line, *move)
					alpha = v
				}
			} else {
				v, t := minimax(newPosition, depth-1, false, alpha, beta, append(line, *move))
				evalCache.Set(newPositionHash, v, cache.DefaultExpiration)
				if v > alpha {
					alpha = v
					newLine = t
				}
			}
			if alpha >= beta {
				return beta, newLine
			}
		}
		return alpha, newLine

	} else {
		for _, move := range moves {
			newPosition := position.Update(move)
			newPositionHash := fmt.Sprintf("%x", newPosition.Hash())
			cachedEval, found := evalCache.Get(newPositionHash)
			if found {
				cacheHits += 1
				v := cachedEval.(float64)
				if v < beta {
					beta = v
					newLine = append(line, *move)
				}
			} else {
				v, t := minimax(newPosition, depth-1, true, alpha, beta, append(line, *move))
				evalCache.Set(newPositionHash, v, cache.DefaultExpiration)
				if v < beta {
					beta = v
					newLine = t
				}
			}
			if beta <= alpha {
				return alpha, newLine
			}
		}
		return beta, newLine
	}
}
