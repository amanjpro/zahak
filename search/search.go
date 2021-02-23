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
var cacheHits int64 = 0
var pv = NewPVLine(100)

type EvalMove struct {
	eval int16
	move *Move
}

func (e *EvalMove) Move() *Move {
	return e.move
}

func (e *EvalMove) Eval() int16 {
	return e.eval
}

func Search(position *Position, depth int8, ply uint16) EvalMove {
	STOP_SEARCH_GLOBALLY = false
	nodesVisited = 0
	cacheHits = 0
	var bestEval EvalMove
	bestMove, score := startMinimax(position, depth, ply)
	bestEval = EvalMove{score, bestMove}
	pv.Pop() // pop our move
	pv.Pop() // pop our opponent's move
	return bestEval
}

func startMinimax(position *Position, depth int8, ply uint16) (*Move, int16) {

	// Collect evaluation for moves per iteration to help us order moves for the
	// next iteration
	legalMoves := position.LegalMoves()
	iterationEvals := make([]int16, len(legalMoves))

	var bestMove *Move
	var previousBestMove *Move
	bestScore := -MAX_INT
	alpha := -MAX_INT
	beta := MAX_INT

	fruitelessIterations := 0

	start := time.Now()
END_LOOP:
	for iterationDepth := int8(1); iterationDepth <= depth; iterationDepth++ {
		if STOP_SEARCH_GLOBALLY {
			break END_LOOP
		}
		currentBestScore := -MAX_INT
		orderedMoves := orderIterationMoves(&IterationMoves{legalMoves, iterationEvals})
		line := NewPVLine(iterationDepth + 1)
		searchPv := true
		for index, move := range orderedMoves {
			if STOP_SEARCH_GLOBALLY {
				break END_LOOP
			}
			fmt.Printf("info currmove %s currmovenumber %d\n\n", move.ToString(), index+1)
			sendPv := false
			cp, ep, tg, hc := position.MakeMove(move)
			score := -MAX_INT
			if searchPv {
				score = -alphaBeta(position, iterationDepth, 1, -beta, -alpha, ply, line)
			} else {
				score = -zeroWindowSearch(position, iterationDepth, 1, -alpha, ply, true)
				if score > alpha { // in fail-soft ... && score < beta ) is common
					score = -alphaBeta(position, iterationDepth, 1, -beta, -alpha, ply, line) // re-search
				}
			}
			// This only works, because checkmate eval is clearly distinguished from
			// maximum/minimum beta/alpha
			if score > beta {
				beta = score
			}
			if score == CHECKMATE_EVAL {
				alpha = score
				iterationEvals[index] = score
				if score > currentBestScore || pv == nil || pv.moveCount > line.moveCount+1 {
					currentBestScore = score
					sendPv = true
					pv.AddFirst(move)
					pv.ReplaceLine(line)
					bestMove = move
					bestScore = currentBestScore
					searchPv = false
				}
			} else if score > alpha && score < beta { // no very hard alpha-beta cutoff
				iterationEvals[index] = score
				alpha = score
				if score > currentBestScore {
					currentBestScore = score
					sendPv = true
					pv.AddFirst(move)
					pv.ReplaceLine(line)
					bestMove = move
					bestScore = currentBestScore
					searchPv = false
				}
			} else {
				iterationEvals[index] = -MAX_INT // if it is, then too bad, that is a bad move
			}
			position.UnMakeMove(move, tg, ep, cp, hc)

			timeSpent := time.Now().Sub(start)
			if sendPv {
				fmt.Printf("info depth %d nps %d tbhits %d hashfull %d nodes %d score cp %d time %d pv %s\n\n",
					pv.moveCount, nps(nodesVisited, timeSpent.Seconds()),
					cacheHits, TranspositionTable.Consumed(), nodesVisited, currentBestScore,
					timeSpent.Milliseconds(), pv.ToString())
			}
		}

		if iterationDepth >= 10 && *bestMove == *previousBestMove {
			fruitelessIterations++
			if fruitelessIterations > 4 {
				break
			}
		} else {
			fruitelessIterations = 0
		}
		timeSpent := time.Now().Sub(start)
		fmt.Printf("info depth %d nps %d tbhits %d hashfull %d nodes %d score cp %d time %d pv %s\n\n",
			pv.moveCount, nps(nodesVisited, timeSpent.Seconds()),
			cacheHits, TranspositionTable.Consumed(), nodesVisited, currentBestScore,
			timeSpent.Milliseconds(), pv.ToString())
		if bestScore == CHECKMATE_EVAL {
			break
		}
		previousBestMove = bestMove
		alpha = -MAX_INT
		beta = MAX_INT
	}

	timeSpent := time.Now().Sub(start)
	fmt.Printf("info depth %d nps %d tbhits %d hashfull %d nodes %d score cp %d time %d pv %s\n\n",
		pv.moveCount, nps(nodesVisited, timeSpent.Seconds()),
		cacheHits, TranspositionTable.Consumed(), nodesVisited, bestScore,
		timeSpent.Milliseconds(), pv.ToString())
	return bestMove, bestScore
}

func alphaBeta(position *Position, depthLeft int8, searchHeight int8, alpha int16, beta int16, ply uint16, pvline *PVLine) int16 {
	nodesVisited += 1
	outcome := position.Status()
	if outcome == Checkmate {
		return -CHECKMATE_EVAL
	} else if outcome == Draw {
		return 0
	}

	if STOP_SEARCH_GLOBALLY {
		return alpha
	}

	if depthLeft == 0 {
		return quiescence(position, alpha, beta, 0, Evaluate(position))
	}

	hash := position.Hash()
	cachedEval, found := TranspositionTable.Get(hash)
	if found && cachedEval.Depth >= depthLeft {
		score := cachedEval.Eval
		if score >= beta && (cachedEval.Type == UpperBound || cachedEval.Type == Exact) {
			cacheHits += 1
			return beta
		}
		if score <= alpha && (cachedEval.Type == LowerBound || cachedEval.Type == Exact) {
			cacheHits += 1
			return alpha
		}
	}

	// NullMove pruning
	isNullMoveAllowed := depthLeft >= 5 && !position.IsEndGame() && !position.IsInCheck()
	R := int8(3)
	if searchHeight > 6 {
		R = 2
	}

	if isNullMoveAllowed {
		tempo := int16(15)    // TODO: Make it variable with a formula like: 10*(numPGAM > 0) + 10* numPGAM > 15);
		bound := beta - tempo // variable bound
		position.NullMove()
		score := -zeroWindowSearch(position, depthLeft-R-1, searchHeight+1, 1-bound, ply, true)
		position.NullMove()
		if score >= bound {
			return beta // null move pruning
		}
	}

	searchPv := true

	legalMoves := position.LegalMoves()
	orderedMoves := orderMoves(&ValidMoves{position, legalMoves, searchHeight + 1})

	for _, move := range orderedMoves {
		capturedPiece, oldEnPassant, oldTag, hc := position.MakeMove(move)
		line := NewPVLine(depthLeft - 1)
		score := -MAX_INT
		if searchPv {
			score = -alphaBeta(position, depthLeft-1, searchHeight+1, -beta, -alpha, ply, line)
		} else {
			score = -zeroWindowSearch(position, depthLeft-1, searchHeight+1, -alpha, ply, true)
			if score > alpha { // in fail-soft ... && score < beta ) is common
				score = -alphaBeta(position, depthLeft-1, searchHeight+1, -beta, -alpha, ply, line) // re-search
			}
		}
		position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece, hc)
		if score >= beta {
			// Those scores are never useufl
			if score != -MAX_INT && score != MAX_INT {
				TranspositionTable.Set(hash, &CachedEval{hash, score, depthLeft, UpperBound, ply})
			}
			return beta
		}
		if score > alpha ||
			(score == CHECKMATE_EVAL && score >= alpha &&
				(pvline == nil || pvline.moveCount < line.moveCount+1)) { // shorter checkmate?
			alpha = score
			// Potential PV move, lets copy it to the current pv-line
			pvline.AddFirst(move)
			pvline.ReplaceLine(line)
			searchPv = false
		}
	}
	if !searchPv {
		TranspositionTable.Set(hash, &CachedEval{hash, alpha, depthLeft, LowerBound, ply})
	} else {
		TranspositionTable.Set(hash, &CachedEval{hash, alpha, depthLeft, Exact, ply})
	}
	return alpha
}

func zeroWindowSearch(position *Position, depthLeft int8, searchHeight int8, beta int16, ply uint16,
	multiCutFlag bool) int16 {
	nodesVisited += 1

	if STOP_SEARCH_GLOBALLY {
		return beta - 1
	}

	if depthLeft <= 0 {
		return quiescence(position, beta-1, beta, 0, Evaluate(position))
	}

	legalMoves := position.LegalMoves()
	orderedMoves := orderMoves(&ValidMoves{position, legalMoves, searchHeight + 1})

	hash := position.Hash()
	cachedEval, found := TranspositionTable.Get(hash)
	if found &&
		cachedEval.Depth >= depthLeft {
		score := cachedEval.Eval
		if score >= beta && (cachedEval.Type != UpperBound || cachedEval.Type == Exact) {
			cacheHits += 1
			return beta
		}
		if score <= beta-1 && (cachedEval.Type != LowerBound || cachedEval.Type == Exact) {
			cacheHits += 1
			return beta - 1
		}
	}
	// Multi-Cut Pruning
	R := int8(3)
	M := 6
	C := 3

	if depthLeft >= R && multiCutFlag && len(legalMoves) > M {
		cutNodeCounter := 0
		for i := 0; i < M; i++ {
			move := legalMoves[i]
			capturedPiece, oldEnPassant, oldTag, hc := position.MakeMove(move)
			score := -zeroWindowSearch(position, depthLeft-1-R, searchHeight+1, 1-beta, ply, !multiCutFlag)
			position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece, hc)
			if score >= beta {
				cutNodeCounter++
				if cutNodeCounter == C {
					return beta // mc-prune
				}
			}
		}
	}

	eval := Evaluate(position)

	rook := WhiteRook
	pawn := WhitePawn
	margin := pawn.Weight() + rook.Weight() // Rook + Pawn
	futility := eval + rook.Weight()

	// Razoring
	if depthLeft < 2 && eval+margin < beta-1 {
		return quiescence(position, beta-1, beta, 0, eval)
	}

	// Reverse Futility Pruning
	if depthLeft < 5 && eval-margin >= beta {
		return eval - margin /* fail soft */
	}

	// Extended Futility Pruning
	isInCheck := position.IsInCheck()
	lastRank := Rank7
	if position.Turn() == Black {
		lastRank = Rank2
	}

	for i, move := range orderedMoves {
		LMR := int8(0)
		if !isInCheck && searchHeight >= 6 && depthLeft == 2 {
			board := position.Board
			movingPiece := board.PieceAt(move.Source)
			isPromoting := (movingPiece.Type() == Pawn && move.Destination.Rank() == lastRank)

			// Extended Futility Pruning
			gain := Evaluate(position) - eval
			if !move.HasTag(Check) && futility+gain <= beta-1 &&
				move.PromoType == NoType && !isPromoting {
				continue
			}

			// Late Move Reduction
			if i >= 5 && !move.HasTag(Check) && move.PromoType == NoType && !isPromoting {
				LMR = 1
			}
		}

		capturedPiece, oldEnPassant, oldTag, hc := position.MakeMove(move)
		score := -zeroWindowSearch(position, depthLeft-1-LMR, searchHeight+1, 1-beta, ply, !multiCutFlag)

		position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece, hc)
		if score >= beta {
			// Those scores are never useufl
			if score != -MAX_INT && score != MAX_INT {
				TranspositionTable.Set(hash, &CachedEval{hash, score, depthLeft, UpperBound, ply})
			}
			return beta // fail-hard beta-cutoff
		}
	}
	TranspositionTable.Set(hash, &CachedEval{hash, beta - 1, depthLeft, LowerBound, ply})
	return beta - 1 // fail-hard, return alpha
}

func nps(nodes int64, dur float64) int64 {
	if dur == 0 {
		return 0
	}
	return int64(float64(nodes) / 1000 * dur)
}
