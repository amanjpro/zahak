package search

import (
	"fmt"
	"sort"
	"time"

	. "github.com/amanjpro/zahak/cache"
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

var STOP_SEARCH_GLOBALLY = false

var nodesVisited int64 = 0
var nodesSearched int64 = 0
var cacheHits int64 = 0
var pv []*Move

type EvalMove struct {
	eval int
	move *Move
	line []*Move
}

func (e *EvalMove) Move() *Move {
	return e.move
}

func (e *EvalMove) Eval() int {
	return e.eval
}

func Search(position *Position, depth int8) EvalMove {
	STOP_SEARCH_GLOBALLY = false
	nodesVisited = 0
	nodesSearched = 0
	cacheHits = 0
	var bestEval EvalMove
	var isMaximizingPlayer = position.Turn() == White
	var dir = -1
	if isMaximizingPlayer {
		dir = 1
	}
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

	eval, moves := minimax(position, depth, 1, isMaximizingPlayer, -MAX_INT,
		MAX_INT, []*Move{})
	move := moves[0]
	mvStr := move.ToString()
	fmt.Printf("info nodes %d score cp %d currmove %s pv",
		nodesVisited, int(eval*dir), mvStr)
	for i := 1; i < len(moves); i++ {
		mv := moves[i]
		fmt.Printf(" %s", mv.ToString())
	}
	fmt.Print("\n\n")

	bestEval = EvalMove{eval, move, moves}
	// }
	end := time.Now()
	// close(evals)
	fmt.Printf("Visited: %d, Selected: %d, Cache-hit: %d\n", nodesVisited, nodesSearched, cacheHits)
	fmt.Printf("Took %f seconds\n", end.Sub(start).Seconds())
	pv = bestEval.line
	TranspositionTable.Rotate()
	return bestEval
}

// func parallelMinimax(position *Position, move *Move, depth int8,
// 	isMaximizingPlayer bool, resultEval chan EvalMove) {
// 	eval, moves := minimax(position, depth, 2, isMaximizingPlayer, math.Inf(-1),
// 		math.Inf(1), move, []*Move{})
// 	resultEval <- EvalMove{eval, move, moves}
// }

func minimax(position *Position, depthLeft int8, pvDepth int8, isMaximizingPlayer bool,
	alpha int, beta int, line []*Move) (int, []*Move) {
	nodesVisited += 1

	if depthLeft == 0 || STOP_SEARCH_GLOBALLY {
		// TODO: Perform all captures before giving up, to avoid the horizon effect
		// var dir float64 = -1
		// if isMaximizingPlayer {
		// 	dir = 1
		// }
		evl := Evaluate(position)
		// fmt.Printf("info nodes %d score cp %d currmove %s pv",
		// 	nodesVisited, int(evl*100*dir), baseMove.ToString())
		// for _, mv := range line {
		// 	fmt.Printf(" %s", mv.ToString())
		// }
		// fmt.Print("\n\n")

		return evl, line
	}

	nodesSearched += 1
	legalMoves := position.LegalMoves()
	orderedMoves := orderMoves(&ValidMoves{position, legalMoves, int(pvDepth)})
	newLine := line

	for _, move := range orderedMoves {
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

	if len(orderedMoves) == 0 {
		return Evaluate(position), line
	} else if isMaximizingPlayer {
		return alpha, newLine
	} else {
		return beta, newLine
	}
}

func getEval(position *Position, depthLeft int8, pvDepth int8, isMaximizingPlayer bool,
	alpha int, beta int, move *Move, line []*Move) (int, []*Move) {
	var score int
	computedLine := []*Move{}
	capturedPiece, oldEnPassant, oldTag := position.MakeMove(move)
	newPositionHash := position.Hash()
	cachedEval, found := TranspositionTable.Get(newPositionHash)
	if found &&
		(cachedEval.Eval == CHECKMATE_EVAL ||
			cachedEval.Eval == -CHECKMATE_EVAL ||
			cachedEval.Depth >= depthLeft) {
		cacheHits += 1
		score = cachedEval.Eval
		computedLine = append(line, move)
	} else {
		v, t := minimax(position, depthLeft, pvDepth, isMaximizingPlayer, alpha, beta, append(line, move))
		TranspositionTable.Set(newPositionHash, &CachedEval{v, int8(len(t))})
		computedLine = t
		score = v
	}
	position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece)
	return score, computedLine
}

type ValidMoves struct {
	position *Position
	moves    []*Move
	depth    int
}

func (validMoves *ValidMoves) Len() int {
	return len(validMoves.moves)
}

func (validMoves *ValidMoves) Swap(i, j int) {
	moves := validMoves.moves
	moves[i], moves[j] = moves[j], moves[i]
}

func (validMoves *ValidMoves) Less(i, j int) bool {
	moves := validMoves.moves
	move1, move2 := moves[i], moves[j]
	board := validMoves.position.Board
	// Is in PV?
	if pv != nil && len(pv) > validMoves.depth {
		if pv[validMoves.depth] == move1 {
			return true
		}
	}

	// Is in Transition table ???
	// TODO: This is slow, that tells us either cache access is slow or has computation is
	// Or maybe (unlikely) make/unmake move is slow
	cp1, ep1, tg1 := validMoves.position.MakeMove(move1)
	hash1 := validMoves.position.Hash()
	validMoves.position.UnMakeMove(move1, tg1, ep1, cp1)
	eval1, ok1 := TranspositionTable.Get(hash1)

	cp2, ep2, tg2 := validMoves.position.MakeMove(move2)
	hash2 := validMoves.position.Hash()
	validMoves.position.UnMakeMove(move2, tg2, ep2, cp2)
	eval2, ok2 := TranspositionTable.Get(hash2)

	if ok1 && ok2 {
		if eval1.Eval > eval2.Eval ||
			(eval1.Eval == eval2.Eval && eval1.Depth >= eval2.Depth) {
			return true
		} else if eval1.Eval < eval2.Eval {
			return false
		}
	}

	// capture ordering
	if move1.HasTag(Capture) && move2.HasTag(Capture) {
		// What are we capturing?
		piece1 := board.PieceAt(move1.Destination)
		piece2 := board.PieceAt(move2.Destination)
		if piece1.Type() > piece2.Type() {
			return true
		}
		// Who is capturing?
		piece1 = board.PieceAt(move1.Source)
		piece2 = board.PieceAt(move2.Source)
		if piece1.Type() <= piece2.Type() {
			return true
		}
		return false
	} else if move1.HasTag(Capture) {
		return true
	}

	piece1 := board.PieceAt(move1.Source)
	piece2 := board.PieceAt(move2.Source)

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

func orderMoves(validMoves *ValidMoves) []*Move {
	sort.Sort(validMoves)
	return validMoves.moves
}
