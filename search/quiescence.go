package search

import (
	. "github.com/amanjpro/zahak/cache"
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func quiescence(position *Position, alpha int, beta int, ply int) int {

	hash := position.Hash()

	item, found := QCache.Get(hash)
	if found {
		return item.Eval
	}

	outcome := position.Status()
	if outcome == Checkmate {
		return -CHECKMATE_EVAL
	} else if outcome == Draw {
		return 0
	}

	legalMoves := position.QuiesceneMoves(ply <= 4)
	orderedMoves := orderMoves(&ValidMoves{position, legalMoves})

	standPat := Evaluate(position)
	if standPat >= beta {
		return standPat
	}

	// Delta pruning
	// w := WhitePawn
	// deltaMargin := w.Weight() * 2 // 200 centipawns
	//
	// if standPat < alpha-deltaMargin {
	// 	return alpha
	// }

	if alpha < standPat {
		alpha = standPat
	}

	if STOP_SEARCH_GLOBALLY {
		return standPat
	}

	for _, move := range orderedMoves {
		cp, ep, tg := position.MakeMove(move)
		score := -quiescence(position, -beta, -alpha, ply+1)
		position.UnMakeMove(move, tg, ep, cp)
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}
	QCache.Set(hash, &QuiescenceEval{hash, alpha})
	return alpha
}
