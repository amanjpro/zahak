package search

import (
	. "github.com/amanjpro/zahak/cache"
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func quescence(position *Position, isMaximizingPlayer bool, alpha int, beta int, ply uint16) int {
	nodesVisited += 1
	nodesSearched += 1

	legalMoves := position.QuiesceneMoves(ply <= 6)
	orderedMoves := orderMoves(&ValidMoves{position, legalMoves})

	if STOP_SEARCH_GLOBALLY {
		return Evaluate(position)
	}

	for _, move := range orderedMoves {
		score := getQuescenceEval(position, !isMaximizingPlayer, alpha, beta, move, ply+1)
		if isMaximizingPlayer {
			if score >= beta {
				return beta
			}
			if score > alpha {
				alpha = score
			}
		} else {
			if score <= alpha {
				return alpha
			}
			if score < beta {
				beta = score
			}
		}
	}

	if len(orderedMoves) == 0 {
		return Evaluate(position)
	} else if isMaximizingPlayer {
		return alpha
	} else {
		return beta
	}
}

func getQuescenceEval(position *Position, isMaximizingPlayer bool,
	alpha int, beta int, move *Move, ply uint16) int {
	var score int
	capturedPiece, oldEnPassant, oldTag := position.MakeMove(move)
	newPositionHash := position.Hash()
	cachedEval, found := TranspositionTable.Get(newPositionHash)
	if found {
		cacheHits += 1
		score = cachedEval.Eval
	} else {
		v := quescence(position, isMaximizingPlayer, alpha, beta, ply)
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
		TranspositionTable.Set(newPositionHash, &CachedEval{position.Hash(), v, 0, tpe, ply})
		score = v
	}
	position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece)
	return score
}
