package search

import (
	. "github.com/amanjpro/zahak/engine"
)

type MovePicker struct {
	position        *Position
	engine          *Engine
	hashmove        Move
	quietMoveList   MoveList
	captureMoveList MoveList
	moveOrder       int8
	canUseHashMove  bool
	isQuiescence    bool
}

func NewMovePicker(p *Position, e *Engine, moveOrder int8, hashmove Move, isQuiescence bool) *MovePicker {
	nextCapture := 0
	nextQuiet := 0
	if hashmove != EmptyMove {
		if hashmove.IsCapture() {
			nextCapture = 1
		} else {
			nextQuiet = 1
		}
	}
	mp := &MovePicker{
		position:        p,
		engine:          e,
		hashmove:        hashmove,
		quietMoveList:   MoveList{Next: nextQuiet},
		captureMoveList: MoveList{Next: nextCapture},
		moveOrder:       moveOrder,
		canUseHashMove:  hashmove != EmptyMove,
		isQuiescence:    isQuiescence,
	}
	return mp
}

func (mp *MovePicker) generateQuietMoves() {
	if mp.isQuiescence || !mp.quietMoveList.IsEmpty() {
		return
	}
	mp.quietMoveList.Moves = mp.position.GetQuietMoves()
	mp.quietMoveList.Size = len(mp.quietMoveList.Moves)
	mp.scoreQuietMoves()
}

func (mp *MovePicker) generateCaptureMoves() {
	if !mp.captureMoveList.IsEmpty() || !mp.quietMoveList.IsEmpty() {
		return
	}
	mp.captureMoveList.Moves = mp.position.GetCaptureMoves()
	mp.captureMoveList.Size = len(mp.captureMoveList.Moves)
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
	mp.captureMoveList.Scores = make([]int32, mp.captureMoveList.Size)

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
			if promoType != NoType {
				p := GetPiece(promoType, White)
				mp.captureMoveList.Scores[i] = 150_000_000 + int32(p.Weight()+capPiece.Weight())
			} else if !move.IsEnPassant() {
				// SEE for ordering
				gain := int32(board.StaticExchangeEval(dest, capPiece, source, piece))
				if gain < 0 {
					mp.captureMoveList.Scores[i] = -90_000_000 + gain
				} else if gain == 0 {
					mp.captureMoveList.Scores[i] = 100_000_000 + int32(capPiece.Weight()-piece.Weight())
				} else {
					mp.captureMoveList.Scores[i] = 100_100_000 + gain
				}
			} else {
				mp.captureMoveList.Scores[i] = 100_100_000 + int32(capPiece.Weight()-piece.Weight())
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
	engine := mp.engine
	moveOrder := mp.moveOrder
	mp.quietMoveList.Scores = make([]int32, mp.quietMoveList.Size)

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
		if move.IsCheck() {
			mp.quietMoveList.Scores[i] = 10_000
			continue
		}

		// King safety (castling)
		isCastling := move.IsKingSideCastle() || move.IsQueenSideCastle()
		if isCastling {
			mp.quietMoveList.Scores[i] = 3_000
			continue
		}

		// Prefer smaller pieces
		if piece.Type() == King {
			mp.quietMoveList.Scores[i] = 0
			continue
		}

		mp.quietMoveList.Scores[i] = 1100 - int32(piece.Weight())
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

	bestIndex := mp.captureMoveList.Next
	for i := bestIndex + 1; i < mp.captureMoveList.Size; i++ {
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
	mp.captureMoveList.SwapWith(bestIndex)
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

	bestIndex := mp.quietMoveList.Next
	for i := bestIndex + 1; i < mp.quietMoveList.Size; i++ {
		if mp.quietMoveList.Scores[i] > mp.quietMoveList.Scores[bestIndex] {
			bestIndex = i
		}
	}
	best := mp.quietMoveList.Moves[bestIndex]
	mp.quietMoveList.SwapWith(bestIndex)
	mp.quietMoveList.IncNext()
	return best
}
