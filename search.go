package main

import (
	"fmt"
	"math"
	"time"

	"github.com/notnil/chess"
	"github.com/patrickmn/go-cache"
)

var evalCache = cache.New(1*time.Hour, 10*time.Minute)

func search(position *chess.Position, depth int8) *chess.Move {
	var bestEval float64
	var bestMove *chess.Move
	var isMaximizingPlayer bool
	if position.Turn() == chess.Black {
		bestEval = math.Inf(1)
		isMaximizingPlayer = false
	} else {
		bestEval = math.Inf(-1)
		isMaximizingPlayer = false
	}
	for _, move := range position.ValidMoves() {
		localEval := minimax(position.Update(move), depth, isMaximizingPlayer, math.Inf(-1), math.Inf(1))
		if position.Turn() == chess.Black {
			if localEval <= bestEval {
				bestEval = localEval
				bestMove = move
			}
		} else {
			if localEval >= bestEval {
				bestEval = localEval
				bestMove = move
			}
		}
	}
	return bestMove
}

func minimax(position *chess.Position, depth int8, isMaximizingPlayer bool, alpha float64, beta float64) float64 {

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
		return cachedEval.(float64)
	}

	if depth == 0 || position.Status() != chess.NoMethod {
		return eval(position)
	}

	moves := position.ValidMoves()
	if isMaximizingPlayer {
		var bestEval float64 = math.Inf(1)
		for _, move := range moves {
			newPosition := position.Update(move)
			newPositionHash := fmt.Sprintf("%x", newPosition.Hash())
			var value float64 = minimax(newPosition, depth+1, false, alpha, beta)
			bestEval = math.Max(bestEval, value)
			evalCache.Set(newPositionHash, bestEval, cache.DefaultExpiration)
			alpha = math.Max(alpha, bestEval)
			if beta <= alpha {
				break
			}
		}
		return bestEval

	} else {
		var bestEval float64 = math.Inf(-1)
		for _, move := range moves {
			newPosition := position.Update(move)
			newPositionHash := fmt.Sprintf("%x", newPosition.Hash())
			var value float64 = minimax(newPosition, depth+1, true, alpha, beta)
			bestEval = math.Min(bestEval, value)
			evalCache.Set(newPositionHash, bestEval, cache.DefaultExpiration)
			beta = math.Min(beta, bestEval)
			if beta <= alpha {
				break
			}
		}
		return bestEval
	}
}
