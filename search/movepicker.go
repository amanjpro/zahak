package search

import (
	"math"

	. "github.com/amanjpro/zahak/engine"
)

type MovePicker struct {
	position        *Position
	engine          *Engine
	hashmove        Move
	quietMoveList   *MoveList
	captureMoveList *MoveList
	moveOrder       int8
	canUseHashMove  bool
	isQuiescence    bool
}

func EmptyMovePicker() *MovePicker {
	qml := NewMoveList(250)
	cml := NewMoveList(250)
	mp := &MovePicker{
		position:        nil,
		engine:          nil,
		hashmove:        EmptyMove,
		quietMoveList:   qml,
		captureMoveList: cml,
		moveOrder:       0,
		canUseHashMove:  false,
		isQuiescence:    false,
	}
	return mp

}

func (mp *MovePicker) RecycleWith(p *Position, e *Engine, moveOrder int8, hashmove Move, isQuiescence bool) {
	mp.engine = e
	mp.position = p
	mp.moveOrder = moveOrder
	mp.hashmove = hashmove
	mp.isQuiescence = isQuiescence
	mp.canUseHashMove = hashmove != EmptyMove
	nextCapture := 0
	nextQuiet := 0
	if hashmove != EmptyMove {
		if hashmove.IsCapture() || hashmove.PromoType() != NoType {
			nextCapture = 1
		} else {
			nextQuiet = 1
		}
	}
	mp.quietMoveList.Size = 0
	mp.quietMoveList.Next = nextQuiet
	mp.quietMoveList.IsScored = false
	mp.captureMoveList.Size = 0
	mp.captureMoveList.Next = nextCapture
	mp.captureMoveList.IsScored = false
}

func (mp *MovePicker) generateQuietMoves() {
	if mp.isQuiescence || !mp.quietMoveList.IsEmpty() {
		return
	}
	mp.position.GetQuietMoves(mp.quietMoveList)
}

func (mp *MovePicker) generateCaptureMoves() {
	if !mp.captureMoveList.IsEmpty() || !mp.quietMoveList.IsEmpty() {
		return
	}
	mp.position.GetCaptureMoves(mp.captureMoveList)
}

func (mp *MovePicker) HasNoPVMove() bool {
	return mp.hashmove == EmptyMove
}

func (mp *MovePicker) UpgradeToPvMove(pvMove Move) {
	if pvMove == EmptyMove || mp.captureMoveList.IsScored || mp.quietMoveList.IsScored {
		return
	}
	mp.hashmove = pvMove
	mp.canUseHashMove = true
	if pvMove.IsCapture() || pvMove.PromoType() != NoType {
		mp.captureMoveList.Next = 1
	} else {
		mp.quietMoveList.Next = 1
	}
}

func (mp *MovePicker) scoreCaptureMoves() int {
	position := mp.position
	board := position.Board
	var highestNonHashIndex int = -1
	var highestNonHashScore int32 = math.MinInt32

	scores := mp.captureMoveList.Scores
	moves := mp.captureMoveList.Moves
	size := mp.captureMoveList.Size

	_ = scores[size-1]
	_ = moves[size-1]

	for i := 0; i < size; i++ {
		move := moves[i]

		if move == mp.hashmove {
			scores[i] = 900_000_000
			mp.captureMoveList.Swap(0, i)
			mp.captureMoveList.Next = 1
			if highestNonHashIndex == 0 {
				highestNonHashIndex = i
			}
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
				scores[i] = 150_000_000 + int32(p.Weight()+capPiece.Weight())
			} else if !move.IsEnPassant() {
				// SEE for ordering
				gain := int32(board.StaticExchangeEval(dest, capPiece, source, piece))
				if gain < 0 {
					scores[i] = -90_000_000 + gain
				} else if gain == 0 {
					scores[i] = 100_000_000 + int32(capPiece.Weight()-piece.Weight())
				} else {
					scores[i] = 100_100_000 + gain
				}
			} else {
				scores[i] = 100_100_000 + int32(capPiece.Weight()-piece.Weight())
			}
			goto end
		}

		if promoType != NoType {
			p := GetPiece(promoType, White)
			scores[i] = 150_000_000 + int32(p.Weight())
			goto end
		}

	end:
		if highestNonHashScore < scores[i] {
			highestNonHashIndex = i
			highestNonHashScore = scores[i]
		}
	}

	mp.captureMoveList.IsScored = true

	return highestNonHashIndex
}

func (mp *MovePicker) scoreQuietMoves() int {

	var highestNonHashIndex int = -1
	var highestNonHashScore int32 = math.MinInt32
	engine := mp.engine
	moveOrder := mp.moveOrder
	scores := mp.quietMoveList.Scores
	moves := mp.quietMoveList.Moves
	size := mp.quietMoveList.Size

	_ = scores[size-1]
	_ = moves[size-1]

	for i := 0; i < size; i++ {
		move := moves[i]

		if move == mp.hashmove {
			scores[i] = 900_000_000
			mp.quietMoveList.Swap(0, i)
			mp.quietMoveList.Next = 1
			if highestNonHashIndex == 0 {
				highestNonHashIndex = i
			}
			continue
		}

		dest := move.Destination()
		piece := move.MovingPiece()
		killer := engine.KillerMoveScore(move, moveOrder)
		var history int32
		var isCastling bool

		if killer != 0 {
			scores[i] = killer
			goto end
		}

		history = engine.MoveHistoryScore(piece, dest, moveOrder)
		if history != 0 {
			scores[i] = history
			goto end
		}

		// prefer checks
		// if move.IsCheck() {
		// 	mp.quietMoveList.Scores[i] = 10_000
		// 	goto end
		// }

		// King safety (castling)
		isCastling = move.IsCastle()
		if isCastling {
			scores[i] = 3_000
			goto end
		}

		// Prefer smaller pieces
		if piece.Type() == King {
			scores[i] = 0
			goto end
		}

		scores[i] = 1100 - int32(piece.Weight())
	end:
		if highestNonHashScore < scores[i] {
			highestNonHashIndex = i
			highestNonHashScore = scores[i]
		}
	}
	mp.quietMoveList.IsScored = true
	return highestNonHashIndex
}

func (mp *MovePicker) Reset() {
	mp.canUseHashMove = mp.hashmove != EmptyMove
	mp.quietMoveList.Next = 0
	mp.captureMoveList.Next = 0
	if mp.canUseHashMove {
		if mp.hashmove.IsCapture() || mp.hashmove.PromoType() != NoType {
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

	size := mp.captureMoveList.Size
	if mp.captureMoveList.Next >= size {
		return EmptyMove
	}

	next := mp.captureMoveList.Next
	var bestIndex int
	scores := mp.captureMoveList.Scores
	_ = scores[size-1]
	if mp.captureMoveList.IsScored {
		bestIndex = next
		_ = scores[bestIndex]
		for i := next + 1; i < size; i++ {
			if scores[i] > scores[bestIndex] {
				bestIndex = i
			}
		}
	} else {
		bestIndex = mp.scoreCaptureMoves()
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

	size := mp.quietMoveList.Size
	if mp.quietMoveList.Next >= size {
		return EmptyMove
	}

	next := mp.quietMoveList.Next
	var bestIndex int
	scores := mp.quietMoveList.Scores
	_ = scores[size-1]
	if mp.quietMoveList.IsScored {
		bestIndex = next
		_ = scores[bestIndex]
		for i := next + 1; i < size; i++ {
			if scores[i] > scores[bestIndex] {
				bestIndex = i
			}
		}
	} else {
		bestIndex = mp.scoreQuietMoves()
	}
	best := mp.quietMoveList.Moves[bestIndex]
	mp.quietMoveList.Swap(next, bestIndex)
	mp.quietMoveList.IncNext()
	return best
}
