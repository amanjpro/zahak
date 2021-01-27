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
}

func search(position *chess.Position, depth int8) *chess.Move {
	nodesVisited = 0
	nodesSearched = 0
	cacheHits = 0
	var bestEval float64
	var bestMove *chess.Move
	var isMaximizingPlayer bool
	if position.Turn() == chess.Black {
		bestEval = math.Inf(1)
		isMaximizingPlayer = false
	} else {
		bestEval = math.Inf(-1)
		isMaximizingPlayer = true
	}
	validMoves := position.ValidMoves()
	evals := make(chan EvalMove)
	start := time.Now()
	for _, move := range validMoves {
		go parallelMinimax(position.Update(move), move, depth, isMaximizingPlayer, evals)
	}
	for i := 0; i < len(validMoves); i++ {
		evalMove := <-evals
		if !isMaximizingPlayer {
			if evalMove.eval <= bestEval {
				bestEval = evalMove.eval
				bestMove = evalMove.move
			}
		} else {
			if evalMove.eval >= bestEval {
				bestEval = evalMove.eval
				bestMove = evalMove.move
			}
		}
	}
	end := time.Now()
	close(evals)
	fmt.Printf("Visited: %d, Selected: %d, Cache-hit: %d\n", nodesVisited, nodesSearched, cacheHits)
	fmt.Printf("Took %d seconds", end.Sub(start).Seconds())
	return bestMove
}

func parallelMinimax(position *chess.Position, move *chess.Move, depth int8, isMaximizingPlayer bool, resultEval chan EvalMove) {
	eval := minimax(position, depth, isMaximizingPlayer, math.Inf(-1), math.Inf(1))
	resultEval <- EvalMove{eval, move}
}

func minimax(position *chess.Position, depth int8, isMaximizingPlayer bool, alpha float64, beta float64) float64 {

	nodesVisited += 1

	if position.Status() == chess.Checkmate {
		if isMaximizingPlayer {
			return math.Inf(-1)
		}
		return math.Inf(1)
	}

	if position.Status() == chess.Stalemate {
		return 0.0
	}

	positionHash := fmt.Sprintf("%x", position.Hash())
	cachedEval, found := evalCache.Get(positionHash)

	if found {
		cacheHits += 1
		return cachedEval.(float64)
	}

	if depth == 0 || position.Status() != chess.NoMethod {
		return eval(position)
	}

	nodesSearched += 1

	moves := position.ValidMoves()
	if isMaximizingPlayer {
		var bestEval float64 = math.Inf(-1)
		for _, move := range moves {
			newPosition := position.Update(move)
			newPositionHash := fmt.Sprintf("%x", newPosition.Hash())
			var value float64 = minimax(newPosition, depth-1, false, alpha, beta)
			bestEval = math.Max(bestEval, value)
			evalCache.Set(newPositionHash, value, cache.DefaultExpiration)
			alpha = math.Max(alpha, bestEval)
			if beta <= alpha {
				break
			}
		}
		return bestEval

	} else {
		var bestEval float64 = math.Inf(1)
		for _, move := range moves {
			newPosition := position.Update(move)
			newPositionHash := fmt.Sprintf("%x", newPosition.Hash())
			var value float64 = minimax(newPosition, depth-1, true, alpha, beta)
			bestEval = math.Min(bestEval, value)
			evalCache.Set(newPositionHash, value, cache.DefaultExpiration)
			beta = math.Min(beta, bestEval)
			if beta <= alpha {
				break
			}
		}
		return bestEval
	}
}
