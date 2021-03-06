package search

import (
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func (e *Engine) quiescence(position *Position, alpha int32, beta int32, currentMove Move, ply int8,
	standPat int32, searchHeight int8) (int32, bool) {

	e.info.quiesceCounter += 1
	e.VisitNode()

	if IsRepetition(position, e.pred, currentMove) {
		return 0, true
	}
	var isInCheck bool
	if currentMove == EmptyMove {
		isInCheck = position.IsInCheck()
	} else {
		isInCheck = currentMove.IsCheck()
	}

	outcome := position.Status(isInCheck)
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

	// Delta Pruning
	q := WhiteQueen
	deltaMargin := q.Weight()
	promoType := currentMove.PromoType()
	if promoType != NoType {
		promo := GetPiece(promoType, White)
		deltaMargin += promo.Weight()
	}
	if !isInCheck && standPat+deltaMargin < alpha {
		e.info.deltaPruningCounter += 1
		return alpha, true
	}

	if alpha < standPat {
		alpha = standPat
	}

	withChecks := false && ply < 4
	legalMoves := position.QuiesceneMoves(withChecks)

	movePicker := NewMovePicker(position, e, legalMoves, searchHeight)

	for i := 0; i < len(legalMoves); i++ {
		move := movePicker.Next()
		isCheckMove := move.IsCheck()
		isCaptureMove := move.IsCapture()
		if !isInCheck && isCaptureMove && !isCheckMove && !move.IsEnPassant() {
			// SEE pruning
			e.info.seeQuiescenceCounter += 1
			if movePicker.scores[i] < 0 {
				continue
			}
		}

		ep, tg, hc := position.MakeMove(move)
		sp := Evaluate(position)

		e.pred.Push(position.Hash())
		v, ok := e.quiescence(position, -beta, -alpha, move, ply+1, sp, searchHeight+1)
		e.pred.Pop()
		position.UnMakeMove(move, tg, ep, hc)
		if !ok {
			return v, ok
		}
		score := -v
		if score >= beta {
			e.AddKillerMove(move, searchHeight)
			return beta, true
		}
		if score > alpha {
			e.AddMoveHistory(move, position.Board.PieceAt(move.Source()), move.Destination(), searchHeight)
			alpha = score
		}
	}
	return alpha, true
}
