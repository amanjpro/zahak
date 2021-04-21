package search

import (
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

type MovePicker struct {
	position        *Position
	engine          *Engine
	hashmove        Move
	quietMoveList   *MoveList
	captureMoveList *MoveList
	moveOrder       int8
	isEndgame       bool
	canUseHashMove  bool
	isQuiescence    bool
}

func EmptyMovePicker() *MovePicker {
	qml := NewMoveList(150)
	cml := NewMoveList(150)
	mp := &MovePicker{
		position:        nil,
		engine:          nil,
		hashmove:        EmptyMove,
		quietMoveList:   qml,
		captureMoveList: cml,
		moveOrder:       0,
		isEndgame:       false,
		canUseHashMove:  false,
		isQuiescence:    false,
	}
	return mp

}

func (mp *MovePicker) RecycleWith(p *Position, e *Engine, moveOrder int8, hashmove Move, isEndgame bool, isQuiescence bool) {
	mp.engine = e
	mp.position = p
	mp.moveOrder = moveOrder
	hashmove.UnMarkCheckAndLegal()
	mp.hashmove = hashmove
	mp.isQuiescence = isQuiescence
	mp.canUseHashMove = hashmove != EmptyMove
	mp.isEndgame = isEndgame
	nextCapture := 0
	nextQuiet := 0
	if hashmove != EmptyMove {
		if hashmove.IsCapture() {
			nextCapture = 1
		} else {
			nextQuiet = 1
		}
	}
	mp.quietMoveList.Size = 0
	mp.quietMoveList.Next = nextQuiet
	mp.captureMoveList.Size = 0
	mp.captureMoveList.Next = nextCapture
}

func (mp *MovePicker) generateQuietMoves() {
	if mp.isQuiescence || !mp.quietMoveList.IsEmpty() {
		return
	}
	mp.position.GetQuietMoves(mp.quietMoveList)
	mp.scoreQuietMoves()
}

func (mp *MovePicker) generateCaptureMoves() {
	if !mp.captureMoveList.IsEmpty() || !mp.quietMoveList.IsEmpty() {
		return
	}
	mp.position.GetCaptureMoves(mp.captureMoveList)
	mp.scoreCaptureMoves()
}

func (mp *MovePicker) HasNoPVMove() bool {
	return mp.hashmove == EmptyMove
}

func (mp *MovePicker) UpgradeToPvMove(pvMove Move) {
	if pvMove == EmptyMove || !mp.captureMoveList.IsEmpty() || !mp.quietMoveList.IsEmpty() {
		return
	}
	mp.hashmove = pvMove
	mp.canUseHashMove = true
}

func (mp *MovePicker) scoreCaptureMoves() {
	position := mp.position
	board := position.Board

	for i := 0; i < mp.captureMoveList.Size; i++ {
		move := mp.captureMoveList.Moves[i]

		if move == mp.hashmove {
			mp.captureMoveList.Scores[i] = 900_000_000
			mp.captureMoveList.Swap(0, i)
			mp.captureMoveList.Next = 1
			continue
		}

		source := move.Source()
		dest := move.Destination()
		piece := move.MovingPiece()
		promoType := move.PromoType()
		//
		// capture ordering
		if move.IsCapture() {
			capPiece := move.CapturedPiece()
			cweight := capPiece.Weight()
			pweight := piece.Weight()
			if promoType != NoType {
				p := GetPiece(promoType, White)
				mp.captureMoveList.Scores[i] = 150_000_000 + int32(p.Weight()+cweight)
				// } else if !move.IsEnPassant() {
			} else if cweight <= pweight {
				// SEE for ordering
				gain := int32(board.StaticExchangeEval(dest, capPiece, source, piece))
				if gain < 0 {
					mp.captureMoveList.Scores[i] = -90_000_000 + gain
				} else if gain == 0 {
					mp.captureMoveList.Scores[i] = 100_000_000 + int32(cweight-pweight)
				} else {
					mp.captureMoveList.Scores[i] = 100_100_000 + gain
				}
			} else {
				mp.captureMoveList.Scores[i] = 100_100_000 + int32(cweight-pweight)
			}
			continue
		}

		if promoType != NoType {
			p := GetPiece(promoType, White)
			mp.captureMoveList.Scores[i] = 150_000_000 + int32(p.Weight())
			continue
		}
	}
}

func (mp *MovePicker) scoreQuietMoves() {
	isEndgame := mp.isEndgame
	engine := mp.engine
	moveOrder := mp.moveOrder

	for i := 0; i < mp.quietMoveList.Size; i++ {
		move := mp.quietMoveList.Moves[i]

		if move == mp.hashmove {
			mp.quietMoveList.Scores[i] = 900_000_000
			mp.quietMoveList.Swap(0, i)
			mp.quietMoveList.Next = 1
			continue
		}

		dest := move.Destination()
		piece := move.MovingPiece()

		killer := engine.KillerMoveScore(move, moveOrder)
		if killer != 0 {
			mp.quietMoveList.Scores[i] = killer
			continue
		}

		history := engine.MoveHistoryScore(piece, dest, moveOrder)
		if history != 0 {
			mp.quietMoveList.Scores[i] = history
			continue
		}

		// prefer checks
		// if move.IsCheck() {
		// 	mp.quietMoveList.Scores[i] = 10_000
		// 	continue
		// }

		// King safety (castling)
		isCastling := move.IsKingSideCastle() || move.IsQueenSideCastle()
		if isCastling {
			mp.quietMoveList.Scores[i] = 6_000
			continue
		}

		mp.quietMoveList.Scores[i] = int32(2000 + PSQT(piece, dest, isEndgame))
		//
		// // Prefer smaller pieces
		// if piece.Type() == King {
		// 	mp.quietMoveList.Scores[i] = 0
		// 	continue
		// }
		//
		// mp.quietMoveList.Scores[i] = 1100 - int32(piece.Weight())
	}
}

func (mp *MovePicker) Reset() {
	mp.canUseHashMove = mp.hashmove != EmptyMove
	mp.quietMoveList.Next = 0
	mp.captureMoveList.Next = 0
	if mp.canUseHashMove {
		if mp.hashmove.IsCapture() {
			mp.captureMoveList.Next = 1
		} else {
			mp.quietMoveList.Next = 1
		}
	}
}

func (mp *MovePicker) Next() Move {
	if mp.hashmove != EmptyMove && mp.canUseHashMove {
		mp.canUseHashMove = false
		return mp.hashmove
	}

	move := mp.getNextCapture()
	if move == EmptyMove {
		return mp.getNextQuiet()
	}
	return move
}

func (mp *MovePicker) getNextCapture() Move {
	if mp.captureMoveList.IsEmpty() {
		mp.generateCaptureMoves()
	}

	if mp.captureMoveList.Next >= mp.captureMoveList.Size {
		return EmptyMove
	}

	next := mp.captureMoveList.Next
	bestIndex := next
	for i := next + 1; i < mp.captureMoveList.Size; i++ {
		if mp.captureMoveList.Scores[i] > mp.captureMoveList.Scores[bestIndex] {
			bestIndex = i
		}
	}
	if mp.captureMoveList.Scores[bestIndex] < 0 {
		alt := mp.getNextQuiet()
		if alt != EmptyMove {
			return alt
		}
	}
	best := mp.captureMoveList.Moves[bestIndex]
	mp.captureMoveList.Swap(next, bestIndex)
	mp.captureMoveList.IncNext()
	return best
}

func (mp *MovePicker) getNextQuiet() Move {
	if mp.quietMoveList.IsEmpty() {
		mp.generateQuietMoves()
	}

	if mp.quietMoveList.Next >= mp.quietMoveList.Size {
		return EmptyMove
	}

	next := mp.quietMoveList.Next
	bestIndex := next
	for i := next + 1; i < mp.quietMoveList.Size; i++ {
		if mp.quietMoveList.Scores[i] > mp.quietMoveList.Scores[bestIndex] {
			bestIndex = i
		}
	}
	best := mp.quietMoveList.Moves[bestIndex]
	mp.quietMoveList.Swap(next, bestIndex)
	mp.quietMoveList.IncNext()
	return best
}
