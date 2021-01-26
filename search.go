package main

import (
	"math"

	"github.com/notnil/chess"
)

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
			if localEval > bestEval {
				bestEval = localEval
				bestMove = move
			}
		}
	}
	return bestMove
}

func eval(position *chess.Position) float64 {
	return 0.0
}

func minimax(position *chess.Position, depth int8, isMaximizingPlayer bool, alpha float64, beta float64) float64 {

	if depth == 0 || position.Status() != chess.NoMethod {
		return eval(position)
	}
	moves := position.ValidMoves()
	if isMaximizingPlayer {
		var bestEval float64 = math.Inf(1)
		for _, move := range moves {
			var value float64 = minimax(position.Update(move), depth+1, false, alpha, beta)
			bestEval = math.Max(bestEval, value)
			alpha = math.Max(alpha, bestEval)
			if beta <= alpha {
				break
			}
		}
		return bestEval

	} else {
		var bestEval float64 = math.Inf(-1)
		for _, move := range moves {
			var value float64 = minimax(position.Update(move), depth+1, true, alpha, beta)
			bestEval = math.Min(bestEval, value)
			beta = math.Min(beta, bestEval)
			if beta <= alpha {
				break
			}
		}
		return bestEval
	}
}
