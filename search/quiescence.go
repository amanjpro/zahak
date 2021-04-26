package search

import (
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

const blackMask = uint64(0x000000000000FF00)
const whiteMask = uint64(0x00FF000000000000)

func dynamicMargin(pos *Position) int16 {

	color := pos.Turn()
	delta := b
	if color == White {
		if pos.Board.GetBitboardOf(WhitePawn)&whiteMask != 0 {
			delta = q
		}
	} else {
		if pos.Board.GetBitboardOf(BlackPawn)&blackMask != 0 {
			delta = q
		}
	}

	other := color.Other()
	if pos.Board.GetBitboardOf(GetPiece(Queen, other)) != 0 {
		return delta + q
	}

	if pos.Board.GetBitboardOf(GetPiece(Rook, other)) != 0 {
		return delta + r
	}

	if pos.Board.GetBitboardOf(GetPiece(Bishop, other)) != 0 || pos.Board.GetBitboardOf(GetPiece(Knight, other)) != 0 {
		return delta + b
	}

	return delta + p
}

func (e *Engine) quiescence(alpha int16, beta int16, currentMove Move, standPat int16, searchHeight int8) (int16, bool) {

	e.info.quiesceCounter += 1
	e.VisitNode()

	var isInCheck = currentMove.IsCheck()

	if standPat >= beta {
		return beta, true // fail hard
	}

	if e.ShouldStop() {
		return 0, false
	}

	position := e.Position

	// Delta Pruning
	if standPat+dynamicMargin(position) < alpha {
		e.info.deltaPruningCounter += 1
		return alpha, true
	}

	if alpha < standPat {
		alpha = standPat
	}

	// withChecks := false && ply < 4
	movePicker := e.MovePickers[searchHeight]
	movePicker.RecycleWith(position, e, -1, EmptyMove, !isInCheck)

	// isEndgame := position.IsEndGame()

	for i := 0; ; i++ {
		move := movePicker.Next()
		if move == EmptyMove {
			break
		}
		// isCheckMove := move.IsCheck()
		// isCaptureMove := move.IsCapture()
		if /*!isCheckMove && !isInCheck && /* isCaptureMove && */ movePicker.captureMoveList.Scores[i] < 0 {
			// SEE pruning
			e.info.seeQuiescenceCounter += 1
			continue
		}

		// promoType := move.PromoType()
		if !IsPromoting(move) {
			margin := p + move.CapturedPiece().Weight()
			// promoType := move.PromoType()
			// if isCaptureMove {
			// 	margin += move.CapturedPiece().Weight()
			// }
			// if promoType != NoType {
			// 	margin += GetPiece(promoType, White).Weight()
			// }
			// toPSQT := PSQT(move.MovingPiece(), move.Destination(), isEndgame)
			if standPat+margin <= alpha {
				e.info.fpCounter += 1
				// position.UnMakeMove(move, tg, ep, hc)
				continue
			}
		}

		ep, tg, hc := position.MakeMove(move)
		sp := Evaluate(position)

		e.pred.Push(position.Hash())
		v, ok := e.quiescence(-beta, -alpha, move, sp, searchHeight+1)
		e.pred.Pop()
		position.UnMakeMove(move, tg, ep, hc)
		if !ok {
			return v, ok
		}
		score := -v
		if score >= beta {
			// e.AddKillerMove(move, searchHeight)
			// e.AddMoveHistory(move, move.MovingPiece(), move.Destination(), searchHeight)
			return beta, true
		}
		if score > alpha {
			alpha = score
		}
	}
	return alpha, true
}
