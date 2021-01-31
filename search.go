package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/notnil/chess"
)

var nodesVisited int64 = 0
var nodesSearched int64 = 0
var cacheHits int64 = 0
var pv []chess.Move

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
		go parallelMinimax(position.Update(move), move, depth, !isMaximizingPlayer, evals)
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
	fmt.Printf("Took %f seconds\n", end.Sub(start).Seconds())
	pv = bestEval.line
	return bestEval
}

func parallelMinimax(position *chess.Position, move *chess.Move, depth int8,
	isMaximizingPlayer bool, resultEval chan EvalMove) {
	eval, moves := minimax(position, depth, 2, isMaximizingPlayer, math.Inf(-1),
		math.Inf(1), []chess.Move{})
	resultEval <- EvalMove{eval, move, moves}
}

func minimax(position *chess.Position, depthLeft int8, pvDepth int8, isMaximizingPlayer bool,
	alpha float64, beta float64, line []chess.Move) (float64, []chess.Move) {
	nodesVisited += 1

	if depthLeft == 0 {
		// TODO: Perform all captures before giving up, to avoid the horizon effect
		return eval(position), line
	}

	nodesSearched += 1

	moves := orderMoves(ValidMoves{position, position.ValidMoves(), int(pvDepth)})
	newLine := line

	for _, move := range moves {
		score, computedLine := getEval(position, depthLeft-1, pvDepth+1,
			!isMaximizingPlayer, alpha, beta, move, line)
		if isMaximizingPlayer {
			if score >= beta {
				return beta, newLine
			}
			if score > alpha {
				newLine = computedLine
				alpha = score
			}
		} else {
			if score <= alpha {
				return alpha, newLine
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

func getEval(position *chess.Position, depthLeft int8, pvDepth int8, isMaximizingPlayer bool,
	alpha float64, beta float64, move *chess.Move, line []chess.Move) (float64,
	[]chess.Move) {
	var score float64
	var computedLine []chess.Move
	newPosition := position.Update(move)
	newHashArray := newPosition.Hash()
	newPositionHash := binary.BigEndian.Uint64(newHashArray[:])
	cachedEval, found := evalCache.Get(newPositionHash)
	if found &&
		(cachedEval.eval == math.Inf(-1) ||
			cachedEval.eval == math.Inf(1) ||
			len(cachedEval.line) >= int(depthLeft)) {
		cacheHits += 1
		score = cachedEval.eval
		computedLine = append(append(line, *move), cachedEval.line...)
	} else {
		v, t := minimax(newPosition, depthLeft, pvDepth, isMaximizingPlayer, alpha, beta, []chess.Move{})
		evalCache.Set(newPositionHash, CachedEval{v, t})
		computedLine = append(append(line, *move), t...)
		score = v
	}
	return score, computedLine
}

type ValidMoves struct {
	position *chess.Position
	moves    []*chess.Move
	depth    int
}

func (validMoves ValidMoves) Len() int {
	return len(validMoves.moves)
}

func (validMoves ValidMoves) Swap(i, j int) {
	validMoves.moves[i], validMoves.moves[j] = validMoves.moves[j], validMoves.moves[i]
}

func (validMoves ValidMoves) Less(i, j int) bool {
	move1, move2 := validMoves.moves[i], validMoves.moves[j]
	board := validMoves.position.Board()
	// Is in PV?
	if pv != nil && len(pv) > validMoves.depth {
		if pv[validMoves.depth] == *move1 {
			return true
		}
	}

	// Is in Transition table ???
	pos1 := validMoves.position.Update(move1)
	hashA1 := pos1.Hash()
	hash1 := binary.BigEndian.Uint64(hashA1[:])

	pos2 := validMoves.position.Update(move2)
	hashA2 := pos2.Hash()
	hash2 := binary.BigEndian.Uint64(hashA2[:])

	if _, ok := evalCache.Get(hash1); ok {
		return true
	}

	if _, ok := evalCache.Get(hash2); ok {
		return false
	}

	// capture ordering
	if move1.HasTag(chess.Capture) && move2.HasTag(chess.Capture) {
		// What are we capturing?
		piece1 := board.Piece(move1.S2())
		piece2 := board.Piece(move2.S2())
		if piece1.Type() > piece2.Type() {
			return true
		}
		// Who is capturing?
		piece1 = board.Piece(move1.S1())
		piece2 = board.Piece(move2.S1())
		if piece1.Type() <= piece2.Type() {
			return true
		}
		return false
	} else if move1.HasTag(chess.Capture) {
		return true
	}

	piece1 := board.Piece(move1.S1())
	piece2 := board.Piece(move2.S1())

	// prefer checks
	if move1.HasTag(chess.Check) {
		return true
	}
	if move2.HasTag(chess.Check) {
		return false
	}
	// Prefer smaller pieces
	if piece1.Type() <= piece2.Type() {
		return true
	}

	return false
}

func orderMoves(validMoves ValidMoves) []*chess.Move {
	sort.Sort(validMoves)
	return validMoves.moves
}
