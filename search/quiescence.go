package search

import (
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func (e *Engine) quiescence(position *Position, alpha int32, beta int32, ply int8, standPat int32, searchHeight int8) int32 {

	e.VisitNode()

	if standPat >= beta {
		return beta // fail hard
	}

	isInCheck := position.IsInCheck()
	// Delta pruning is slowing things down
	// p := WhitePawn
	// deltaMargin := int32(p.Weight() * 1)
	// if !isInCheck && standPat < alpha-deltaMargin { // is capture
	// 	return alpha
	// }

	if alpha < standPat {
		alpha = standPat
	}

	withChecks := false && ply <= 4
	legalMoves := position.QuiesceneMoves(withChecks)

	if len(legalMoves) == 0 {
		outcome := position.Status()
		if outcome == Checkmate {
			return -CHECKMATE_EVAL
		} else if outcome == Draw {
			return 0
		}
	}

	movePicker := NewMovePicker(position, e, legalMoves, searchHeight)

	if e.ShouldStop() {
		return standPat
	}

	for i := 0; i < len(legalMoves); i++ {
		move := movePicker.Next()
		if !isInCheck && move.HasTag(Capture) && !move.HasTag(EnPassant) {
			// SEE pruning
			if movePicker.scores[i] < 0 {
				continue
			}
		}
		cp, ep, tg, hc := position.MakeMove(move)
		sp := Evaluate(position)
		score := -e.quiescence(position, -beta, -alpha, ply+1, sp, searchHeight+1)
		position.UnMakeMove(move, tg, ep, cp, hc)
		if score >= beta {
			e.AddKillerMove(move, searchHeight)
			return beta
		}
		if score > alpha {
			e.AddMoveHistory(move, position.Board.PieceAt(move.Source), move.Destination, searchHeight)
			alpha = score
		}
	}
	return alpha
}
