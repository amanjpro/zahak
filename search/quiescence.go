package search

import (
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func quiescence(position *Position, alpha int, beta int, ply int, standPat int) int {

	outcome := position.Status()
	if outcome == Checkmate {
		return -CHECKMATE_EVAL
	} else if outcome == Draw {
		return 0
	}

	withChecks := ply <= 4
	legalMoves := position.QuiesceneMoves(withChecks)
	orderedMoves := orderMoves(&ValidMoves{position, legalMoves, 125})

	if standPat >= beta {
		return beta // fail hard
	}

	if alpha < standPat {
		alpha = standPat
	}

	if STOP_SEARCH_GLOBALLY {
		return standPat
	}

	w := WhitePawn
	deltaMargin := w.Weight() * 2 // 200 centipawns
	for _, move := range orderedMoves {
		cp, ep, tg := position.MakeMove(move)
		sp := Evaluate(position)
		if cp != NoPiece && standPat < alpha-deltaMargin { // is capture
			// Delta pruning meaningless captures
			position.UnMakeMove(move, tg, ep, cp)
			return alpha
		}
		score := -quiescence(position, -beta, -alpha, ply+1, sp)
		position.UnMakeMove(move, tg, ep, cp)
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}
	return alpha
}
