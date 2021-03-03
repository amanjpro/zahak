package search

import (
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func (e *Engine) quiescence(position *Position, alpha int32, beta int32, ply int8,
	standPat int32, searchHeight int8) (int32, bool) {

	e.VisitNode()
	outcome := position.Status()
	if outcome == Checkmate {
		return -CHECKMATE_EVAL, true
	} else if outcome == Draw {
		return 0, true
	}

	if standPat >= beta {
		return beta, true // fail hard
	}

	if e.ShouldStop() {
		return 0, false
	}

	isInCheck := position.IsInCheck()

	if alpha < standPat {
		alpha = standPat
	}

	withChecks := ply < 4
	legalMoves := position.QuiesceneMoves(withChecks)

	movePicker := NewMovePicker(position, e, legalMoves, searchHeight)

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

		var score int32
		callQuiescence := true
		if !isInCheck && !move.HasTag(Check) {
			// The logic looks difficult, but it is not
			// I basically pretend that I have called quiescence
			// with the reversed alpha/beta (like normal in negamax)
			// and then, if the bounds exceeded I pretend that I return
			// either alpha or beta
			newAlpha := -beta
			newBeta := -alpha
			q := WhiteQueen
			deltaMargin := q.Weight()
			if move.PromoType != NoType {
				promo := GetPiece(move.PromoType, White)
				deltaMargin += promo.Weight()
			}
			if sp >= newBeta {
				score = -newBeta
				position.UnMakeMove(move, tg, ep, cp, hc)
				callQuiescence = false
			}
			if sp+deltaMargin < newAlpha { // is capture
				position.UnMakeMove(move, tg, ep, cp, hc)
				callQuiescence = false
				score = -newAlpha
			}
		}

		if callQuiescence {
			v, ok := e.quiescence(position, -beta, -alpha, ply+1, sp, searchHeight+1)
			position.UnMakeMove(move, tg, ep, cp, hc)
			if !ok {
				return v, ok
			}
			score = -v
		}
		if score >= beta {
			e.AddKillerMove(move, searchHeight)
			return beta, true
		}
		if score > alpha {
			e.AddMoveHistory(move, position.Board.PieceAt(move.Source), move.Destination, searchHeight)
			alpha = score
		}
	}
	return alpha, true
}
