package search

import (
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func quescence(position *Position, isMaximizingPlayer bool, alpha int, beta int, ply int) int {
	legalMoves := position.QuiesceneMoves(ply <= 6)
	orderedMoves := orderMoves(&ValidMoves{position, legalMoves})

	if position.Status() == Checkmate {
		if isMaximizingPlayer {
			return -CHECKMATE_EVAL
		}
		return CHECKMATE_EVAL
	}

	standPat := Evaluate(position)
	if standPat >= beta {
		return beta
	}
	if alpha < standPat {
		alpha = standPat
	}

	if STOP_SEARCH_GLOBALLY {
		if isMaximizingPlayer {
			return standPat
		} else {
			return -standPat
		}
	}

	bestScore := standPat

	for _, move := range orderedMoves {
		cp, ep, tg := position.MakeMove(move)
		score := -quescence(position, !isMaximizingPlayer, -beta, -alpha, ply+1)
		position.UnMakeMove(move, tg, ep, cp)
		if score > bestScore {
			bestScore = score
			if score > alpha {
				if score >= beta {
					break
				}
				alpha = score
			}
		}
	}
	return bestScore
}
