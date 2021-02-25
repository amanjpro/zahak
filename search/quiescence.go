package search

import (
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func (e *Engine) quiescence(position *Position, alpha int32, beta int32, ply int8, standPat int32, searchHeight int8) int32 {

	e.VisitNode()

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

	if standPat >= beta {
		return beta // fail hard
	}

	if alpha < standPat {
		alpha = standPat
	}

	if e.ShouldStop() {
		return standPat
	}

	isInCheck := position.IsInCheck()

	w := WhitePawn
	deltaMargin := w.Weight() * 2 // 200 centipawns
	for i := 0; i < len(legalMoves); i++ {
		move := movePicker.Next()
		if !isInCheck && move.HasTag(Capture) && !move.HasTag(EnPassant) {
			// SEE pruning
			if movePicker.scores[i] < 0 {
				continue
			}
		}
		cp, ep, tg, hc := position.MakeMove(move)
		if !isInCheck && cp != NoPiece && standPat < alpha-deltaMargin { // is capture
			// Delta pruning meaningless captures
			position.UnMakeMove(move, tg, ep, cp, hc)
			continue
		}
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
