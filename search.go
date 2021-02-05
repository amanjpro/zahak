package main

import (
	"fmt"
	"math"
	"sort"
	"time"
)

var nodesVisited int64 = 0
var nodesSearched int64 = 0
var cacheHits int64 = 0
var pv []Move

type EvalMove struct {
	eval float64
	move *Move
	line []Move
}

func search(position *Position, depth int8) *EvalMove {
	nodesVisited = 0
	nodesSearched = 0
	cacheHits = 0
	var bestEval *EvalMove
	var isMaximizingPlayer = position.Turn() == White
	validMoves := position.LegalMoves()
	evals := make(chan EvalMove)
	start := time.Now()
	for i := 0; i < len(validMoves); i++ {
		p := position.copy()
		move := validMoves[i]
		p.MakeMove(move)
		go parallelMinimax(p, &move, depth, !isMaximizingPlayer, evals)
	}
	for i := 0; i < len(validMoves); i++ {
		evalMove := <-evals

		fmt.Println("Move", evalMove.move.ToString())
		fmt.Println("Tree")
		for _, mv := range evalMove.line {
			fmt.Printf("%s ", mv.ToString())
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

func parallelMinimax(position *Position, move *Move, depth int8,
	isMaximizingPlayer bool, resultEval chan EvalMove) {
	eval, moves := minimax(position, depth, 2, isMaximizingPlayer, math.Inf(-1),
		math.Inf(1), []Move{})
	resultEval <- EvalMove{eval, move, moves}
}

func minimax(position *Position, depthLeft int8, pvDepth int8, isMaximizingPlayer bool,
	alpha float64, beta float64, line []Move) (float64, []Move) {
	nodesVisited += 1

	if depthLeft == 0 {
		// TODO: Perform all captures before giving up, to avoid the horizon effect
		eval := eval(position)
		fmt.Printf("%s, %f\n", position.Fen(), eval)
		return eval, line
	}

	nodesSearched += 1
	legalMoves := position.LegalMoves()
	orderedMoves := orderMoves(ValidMoves{position, &legalMoves, int(pvDepth)})
	newLine := line

	for _, move := range *orderedMoves {
		score, computedLine := getEval(position, depthLeft-1, pvDepth+1,
			!isMaximizingPlayer, alpha, beta, &move, line)
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

func getEval(position *Position, depthLeft int8, pvDepth int8, isMaximizingPlayer bool,
	alpha float64, beta float64, move *Move, line []Move) (float64, []Move) {
	var score float64
	var computedLine []Move
	oldTag := position.tag
	oldEnPassant := position.enPassant
	capturedPiece := position.board.PieceAt(move.destination)
	position.MakeMove(*move)
	newPositionHash := position.Hash()
	cachedEval, found := evalCache.Get(newPositionHash)
	if found &&
		(cachedEval.eval == math.Inf(-1) ||
			cachedEval.eval == math.Inf(1) ||
			len(cachedEval.line) >= int(depthLeft)) {
		cacheHits += 1
		score = cachedEval.eval
		computedLine = append(append(line, *move), cachedEval.line...)
	} else {
		v, t := minimax(position, depthLeft, pvDepth, isMaximizingPlayer, alpha, beta, []Move{})
		evalCache.Set(newPositionHash, CachedEval{v, t})
		computedLine = append(append(line, *move), t...)
		score = v
	}
	position.UnMakeMove(*move, oldTag, oldEnPassant, capturedPiece)
	return score, computedLine
}

type ValidMoves struct {
	position *Position
	moves    *[]Move
	depth    int
}

func (validMoves ValidMoves) Len() int {
	return len(*validMoves.moves)
}

func (validMoves ValidMoves) Swap(i, j int) {
	moves := *validMoves.moves
	moves[i], moves[j] = moves[j], moves[i]
}

func (validMoves ValidMoves) Less(i, j int) bool {
	moves := *validMoves.moves
	move1, move2 := moves[i], moves[j]
	board := validMoves.position.board
	// Is in PV?
	if pv != nil && len(pv) > validMoves.depth {
		if pv[validMoves.depth] == move1 {
			return true
		}
	}

	// FIXME
	// // Is in Transition table ???
	// pos1 := validMoves.position.Update(move1)
	// hashA1 := pos1.Hash()
	// hash1 := binary.BigEndian.Uint64(hashA1[:])
	//
	// pos2 := validMoves.position.Update(move2)
	// hashA2 := pos2.Hash()
	// hash2 := binary.BigEndian.Uint64(hashA2[:])

	// if _, ok := evalCache.Get(hash1); ok {
	// 	return true
	// }
	//
	// if _, ok := evalCache.Get(hash2); ok {
	// 	return false
	// }

	// capture ordering
	if move1.HasTag(Capture) && move2.HasTag(Capture) {
		// What are we capturing?
		piece1 := board.PieceAt(move1.destination)
		piece2 := board.PieceAt(move2.destination)
		if piece1.Type() > piece2.Type() {
			return true
		}
		// Who is capturing?
		piece1 = board.PieceAt(move1.source)
		piece2 = board.PieceAt(move2.source)
		if piece1.Type() <= piece2.Type() {
			return true
		}
		return false
	} else if move1.HasTag(Capture) {
		return true
	}

	piece1 := board.PieceAt(move1.source)
	piece2 := board.PieceAt(move2.source)

	// prefer checks
	if move1.HasTag(Check) {
		return true
	}
	if move2.HasTag(Check) {
		return false
	}
	// Prefer smaller pieces
	if piece1.Type() <= piece2.Type() {
		return true
	}

	return false
}

func orderMoves(validMoves ValidMoves) *[]Move {
	sort.Sort(validMoves)
	return validMoves.moves
}
