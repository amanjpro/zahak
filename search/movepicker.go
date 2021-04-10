package search

import (
	. "github.com/amanjpro/zahak/engine"
)

type MoveList struct {
	moves  []Move
	scores []int32
	size   int
	next   int
}

func (ml *MoveList) IsEmpty() bool {
	return ml.moves == nil
}

func (ml *MoveList) SwapWith(best int) {
	ml.moves[ml.next], ml.moves[best] = ml.moves[best], ml.moves[ml.next]
	ml.scores[ml.next], ml.scores[best] = ml.scores[best], ml.scores[ml.next]
}

func (ml *MoveList) Swap(first int, second int) {
	ml.moves[first], ml.moves[second] = ml.moves[second], ml.moves[first]
	ml.scores[first], ml.scores[second] = ml.scores[second], ml.scores[first]
}

func (ml *MoveList) Next() {
	ml.next += 1
}

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
		quietMoveList:   MoveList{next: nextQuiet},
		captureMoveList: MoveList{next: nextCapture},
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
	mp.quietMoveList.moves = mp.position.GetQuietMoves()
	mp.quietMoveList.size = len(mp.quietMoveList.moves)
	mp.scoreQuietMoves()
}

func (mp *MovePicker) generateCaptureMoves() {
	if !mp.captureMoveList.IsEmpty() || !mp.quietMoveList.IsEmpty() {
		return
	}
	mp.captureMoveList.moves = mp.position.GetCaptureMoves()
	mp.captureMoveList.size = len(mp.captureMoveList.moves)
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
	mp.captureMoveList.scores = make([]int32, mp.captureMoveList.size)

	for i := 0; i < mp.captureMoveList.size; i++ {
		move := mp.captureMoveList.moves[i]

		if move == mp.hashmove {
			mp.captureMoveList.scores[i] = 900_000_000
			mp.captureMoveList.Swap(0, i)
			mp.captureMoveList.next = 1
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
				mp.captureMoveList.scores[i] = 150_000_000 + int32(p.Weight()+capPiece.Weight())
			} else if !move.IsEnPassant() {
				// SEE for ordering
				gain := int32(board.StaticExchangeEval(dest, capPiece, source, piece))
				if gain < 0 {
					mp.captureMoveList.scores[i] = -90_000_000 + gain
				} else if gain == 0 {
					mp.captureMoveList.scores[i] = 100_000_000 + int32(capPiece.Weight()-piece.Weight())
				} else {
					mp.captureMoveList.scores[i] = 100_100_000 + gain
				}
			} else {
				mp.captureMoveList.scores[i] = 100_100_000 + int32(capPiece.Weight()-piece.Weight())
			}
			continue
		}

		if promoType != NoType {
			p := GetPiece(promoType, White)
			mp.captureMoveList.scores[i] = 150_000_000 + int32(p.Weight())
			continue
		}
	}
}

func (mp *MovePicker) scoreQuietMoves() {
	engine := mp.engine
	moveOrder := mp.moveOrder
	mp.quietMoveList.scores = make([]int32, mp.quietMoveList.size)

	for i := 0; i < mp.quietMoveList.size; i++ {
		move := mp.quietMoveList.moves[i]

		if move == mp.hashmove {
			mp.quietMoveList.scores[i] = 900_000_000
			mp.quietMoveList.Swap(0, i)
			mp.quietMoveList.next = 1
			continue
		}

		dest := move.Destination()
		piece := move.MovingPiece()

		killer := engine.KillerMoveScore(move, moveOrder)
		if killer != 0 {
			mp.quietMoveList.scores[i] = killer
			continue
		}

		history := engine.MoveHistoryScore(piece, dest, moveOrder)
		if history != 0 {
			mp.quietMoveList.scores[i] = history
			continue
		}

		// prefer checks
		if move.IsCheck() {
			mp.quietMoveList.scores[i] = 10_000
			continue
		}

		// King safety (castling)
		isCastling := move.IsKingSideCastle() || move.IsQueenSideCastle()
		if isCastling {
			mp.quietMoveList.scores[i] = 3_000
			continue
		}

		// Prefer smaller pieces
		if piece.Type() == King {
			mp.quietMoveList.scores[i] = 0
			continue
		}

		mp.quietMoveList.scores[i] = 1100 - int32(piece.Weight())
	}
}

func (mp *MovePicker) Reset() {
	mp.canUseHashMove = mp.hashmove != EmptyMove
	mp.quietMoveList.next = 0
	mp.captureMoveList.next = 0
	if mp.canUseHashMove {
		if mp.hashmove.IsCapture() {
			mp.captureMoveList.next = 1
		} else {
			mp.quietMoveList.next = 1
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

	if mp.captureMoveList.next >= mp.captureMoveList.size {
		return EmptyMove
	}

	bestIndex := mp.captureMoveList.next
	for i := bestIndex + 1; i < mp.captureMoveList.size; i++ {
		if mp.captureMoveList.scores[i] > mp.captureMoveList.scores[bestIndex] {
			bestIndex = i
		}
	}
	if mp.captureMoveList.scores[bestIndex] < 0 {
		alt := mp.getNextQuiet()
		if alt != EmptyMove {
			return alt
		}
	}
	best := mp.captureMoveList.moves[bestIndex]
	mp.captureMoveList.SwapWith(bestIndex)
	mp.captureMoveList.Next()
	return best
}

func (mp *MovePicker) getNextQuiet() Move {
	if mp.quietMoveList.IsEmpty() {
		mp.generateQuietMoves()
	}

	if mp.quietMoveList.next >= mp.quietMoveList.size {
		return EmptyMove
	}

	bestIndex := mp.quietMoveList.next
	for i := bestIndex + 1; i < mp.quietMoveList.size; i++ {
		if mp.quietMoveList.scores[i] > mp.quietMoveList.scores[bestIndex] {
			bestIndex = i
		}
	}
	best := mp.quietMoveList.moves[bestIndex]
	mp.quietMoveList.SwapWith(bestIndex)
	mp.quietMoveList.Next()
	return best
}
