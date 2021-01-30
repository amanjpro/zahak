package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/notnil/chess"
)

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
	var bestEval *EvalMove
	var isMaximizingPlayer = position.Turn() == chess.White
	validMoves := position.ValidMoves()
	evals := make(chan EvalMove)
	start := time.Now()
	for _, move := range validMoves {
		go parallelMinimax(position.Update(move), move, depth, isMaximizingPlayer, evals)
	}
	for i := 0; i < len(validMoves); i++ {
		evalMove := <-evals

		fmt.Println("Move", evalMove.move.String())
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
			return math.Inf(1), line
		}
		return math.Inf(-1), line
	}

	if depth == 0 {
		// TODO: Perform all captures before giving up, to avoid the horizon effect
		// if isMaximizingPlayer {
		return eval(position), line
		// }
		// return -eval(position), line
	}

	nodesSearched += 1

	moves := position.ValidMoves()
	newLine := line

	for _, move := range moves {
		score, computedLine := getEval(position, depth-1, !isMaximizingPlayer, alpha, beta, move, line)
		if isMaximizingPlayer {
			if score >= beta {
				return beta, line
			}
			if score > alpha {
				newLine = computedLine
				alpha = score
			}
		} else {
			if score <= alpha {
				return alpha, line
			}
			if score < beta {
				newLine = computedLine
				beta = score
			}
		}
	}
	if isMaximizingPlayer {
		return alpha, newLine
	} else {
		return beta, newLine
	}
}

func getEval(position *chess.Position, depth int8, isMaximizingPlayer bool, alpha float64, beta float64, move *chess.Move, line []chess.Move) (float64, []chess.Move) {
	var score float64
	var computedLine []chess.Move
	newPosition := position.Update(move)
	newHashArray := newPosition.Hash()
	newPositionHash := binary.BigEndian.Uint64(newHashArray[:])
	cachedEval, found := evalCache.Get(newPositionHash)
	if found && len(cachedEval.line) >= int(depth-1) {
		cacheHits += 1
		score = cachedEval.eval
		computedLine = append(append(line, *move), cachedEval.line...)
	} else {
		v, t := minimax(newPosition, depth, isMaximizingPlayer, alpha, beta, []chess.Move{})
		evalCache.Set(newPositionHash, CachedEval{v, t})
		computedLine = append(append(line, *move), t...)
		score = v
	}
	return score, computedLine
}
