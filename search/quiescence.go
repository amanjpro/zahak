package search

import (
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func (e *Engine) quiescence(position *Position, alpha int16, beta int16, currentMove Move, ply int8,
	standPat int16, searchHeight int8) (int16, bool) {

	e.info.quiesceCounter += 1
	e.VisitNode()

	var isInCheck = currentMove.IsCheck()

	if standPat >= beta {
		return beta, true // fail hard
	}

	if e.ShouldStop() {
		return 0, false
	}

	// Delta Pruning
	deltaMargin := WhiteQueen.Weight()
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

	// withChecks := false && ply < 4
	movePicker := NewMovePicker(position, e, searchHeight, EmptyMove, !isInCheck)

	for i := 0; ; i++ {
		move := movePicker.Next()
		if move == EmptyMove {
			break
		}
		isCheckMove := move.IsCheck()
		isCaptureMove := move.IsCapture()
		if !isInCheck && isCaptureMove && !isCheckMove && !move.IsEnPassant() {
			if movePicker.captureScores[i] < 0 {
				// SEE pruning
				e.info.seeQuiescenceCounter += 1
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
			e.AddMoveHistory(move, move.MovingPiece(), move.Destination(), searchHeight)
			alpha = score
		}
	}
	return alpha, true
}
