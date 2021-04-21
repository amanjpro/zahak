package search

import (
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

const rank2 = uint64(0x000000000000FF00)
const rank7 = uint64(0x00FF000000000000)

func dynamicMargin(p *Position) int16 {

	color := p.Turn()
	delta := WhitePawn.Weight()
	if color == White {
		if p.Board.GetBitboardOf(WhitePawn)&rank7 != 0 {
			delta = WhiteQueen.Weight()
		}
	} else {
		if p.Board.GetBitboardOf(BlackPawn)&rank2 != 0 {
			delta = WhiteQueen.Weight()
		}
	}

	if p.Board.GetBitboardOf(GetPiece(Queen, color)) != 0 {
		return delta + WhiteQueen.Weight()
	}

	if p.Board.GetBitboardOf(GetPiece(Rook, color)) != 0 {
		return delta + WhiteRook.Weight()
	}

	if p.Board.GetBitboardOf(GetPiece(Bishop, color)) != 0 || p.Board.GetBitboardOf(GetPiece(Knight, color)) != 0 {
		return delta + WhiteBishop.Weight()
	}

	return delta + WhitePawn.Weight()
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
	if !isInCheck && standPat+dynamicMargin(position) < alpha {
		e.info.deltaPruningCounter += 1
		return alpha, true
	}

	if alpha < standPat {
		alpha = standPat
	}

	// withChecks := false && ply < 4
	isEndgame := position.IsEndGame()
	movePicker := e.MovePickers[searchHeight]
	movePicker.RecycleWith(position, e, -1, EmptyMove, isEndgame, !isInCheck)

	for i := 0; ; i++ {
		move := movePicker.Next()
		if move == EmptyMove {
			break
		}

		if ep, tg, hc, legal := position.MakeMove(&move); legal {
			isCheckMove := move.IsCheck()
			isCaptureMove := move.IsCapture()
			if !isCheckMove && isCaptureMove && movePicker.captureMoveList.Scores[i] < 0 {
				// SEE pruning
				e.info.seeQuiescenceCounter += 1
				position.UnMakeMove(move, tg, ep, hc)
				continue
			}
			// if !isCheckMove && isCaptureMove && move.PromoType() == NoType {
			// 	toPSQT := PSQT(move.MovingPiece(), move.Destination(), isEndgame)
			// 	margin := WhiteQueen.Weight() + toPSQT
			// 	if standPat+margin <= alpha {
			// 		e.info.fpCounter += 1
			// 		position.UnMakeMove(move, tg, ep, hc)
			// 		continue
			// 	}
			// }

			// if !isInCheck && !isCheckMove {
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
			// e.RemoveMoveHistory(move, move.MovingPiece(), move.Destination(), searchHeight)
		}
	}
	return alpha, true
}
