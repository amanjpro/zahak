package search

import (
	"math"

	. "github.com/amanjpro/zahak/engine"
)

type MovePicker struct {
	position        *Position
	engine          *Engine
	hashmove        Move
	quietMoveList   MoveList
	captureMoveList MoveList
	captureSees     []int32
	searchHeight    int8
	currentMove     Move
	canUseHashMove  bool
	isQuiescence    bool
	killer1         Move
	killer2         Move
	killerIndex     int
	counterMove     Move
	currentDepth    int8
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
		captureSees:     make([]int32, 250),
		searchHeight:    0,
		canUseHashMove:  false,
		isQuiescence:    false,
		killer1:         EmptyMove,
		killer2:         EmptyMove,
		killerIndex:     0,
		counterMove:     EmptyMove,
	}
	return mp
}

func (mp *MovePicker) RecycleWith(p *Position, e *Engine, searchHeight int8, hashmove Move, currentDepth int8, isQuiescence bool) {
	mp.engine = e
	mp.position = p
	mp.searchHeight = searchHeight
	mp.hashmove = hashmove
	mp.isQuiescence = isQuiescence
	mp.canUseHashMove = hashmove != EmptyMove
	if searchHeight >= 0 {
		mp.currentMove = e.positionMoves[searchHeight]
	} else {
		mp.currentMove = EmptyMove
	}

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
	mp.currentDepth = currentDepth

	if !isQuiescence {
		mp.killer1, mp.killer2 = mp.engine.searchHistory.KillerMoveAt(searchHeight)
		if mp.killer1 == hashmove {
			mp.killer1 = EmptyMove
		}
		if mp.killer2 == hashmove {
			mp.killer2 = EmptyMove
		}
		mp.killerIndex = 1
		if mp.currentMove != EmptyMove {
			counterMove := mp.engine.searchHistory.CounterMoveAt(mp.currentMove)
			if counterMove != mp.killer1 && counterMove != mp.killer2 && counterMove != hashmove {
				mp.counterMove = counterMove
			} else {
				mp.counterMove = EmptyMove
			}
		}
	} else {
		mp.killerIndex = 0
		mp.killer1, mp.killer2 = EmptyMove, EmptyMove
		mp.counterMove = EmptyMove
	}
}

func (mp *MovePicker) generateQuietMoves() {
	if mp.isQuiescence || !mp.quietMoveList.IsEmpty() {
		return
	}
	mp.position.GetQuietMoves(&mp.quietMoveList)
}

func (mp *MovePicker) generateCaptureMoves() {
	if !mp.captureMoveList.IsEmpty() || !mp.quietMoveList.IsEmpty() {
		return
	}
	mp.position.GetCaptureMoves(&mp.captureMoveList)
}

func (mp *MovePicker) HasNoPVMove() bool {
	return mp.hashmove == EmptyMove
}

func (mp *MovePicker) UpgradeToPvMove(pvMove Move) {
	if pvMove == EmptyMove || mp.captureMoveList.IsScored || mp.quietMoveList.IsScored {
		return
	}
	mp.hashmove = pvMove
	if mp.killer1 == pvMove {
		mp.killer1 = EmptyMove
	}
	if mp.killer2 == pvMove {
		mp.killer2 = EmptyMove
	}
	mp.canUseHashMove = true
	if pvMove.IsCapture() || pvMove.PromoType() != NoType {
		mp.captureMoveList.Next = 1
	} else {
		mp.quietMoveList.Next = 1
	}
}

func (mp *MovePicker) scoreCaptureMoves() int {
	// position := mp.position
	// board := position.Board
	var highestNonHashIndex int = -1
	var highestNonHashScore int32 = math.MinInt32
	engine := mp.engine

	scores := mp.captureMoveList.Scores
	moves := mp.captureMoveList.Moves
	size := mp.captureMoveList.Size

	_ = scores[size-1]
	_ = moves[size-1]

	for i := 0; i < size; i++ {
		move := moves[i]
		mp.captureSees[i] = 0

		if move == mp.hashmove {
			scores[i] = 900_000_000
			mp.captureMoveList.Swap(0, i)
			mp.captureSees[0], mp.captureSees[i] = mp.captureSees[i], mp.captureSees[0]
			if highestNonHashIndex == 0 {
				highestNonHashIndex = i
			}
			continue
		}

		// source := move.Source()
		// dest := move.Destination()
		board := mp.position.Board
		history := engine.searchHistory.TacticalHistory(move) / 128
		promoType := move.PromoType()

		// capture ordering
		if move.IsCapture() {
			capPiece := move.CapturedPiece()
			promoType := move.PromoType()
			source := move.Source()
			dest := move.Destination()
			piece := move.MovingPiece()
			if promoType != NoType {
				p := GetPiece(promoType, White)
				scores[i] = history + 2000 + int32(p.Weight()+capPiece.Weight())
			} else if !move.IsEnPassant() {
				// SEE for ordering
				gain := int32(board.SeeGe(dest, capPiece, source, piece, -50*int16(mp.currentDepth)))
				mp.captureSees[i] = gain
				if gain < 0 {
					scores[i] = history + /* -90_000_000 + */ gain
				} else if gain == 0 {
					scores[i] = history + /* 100_000_000 + */ int32(capPiece.Weight()-piece.Weight())
				} else {
					scores[i] = history + /* 100_100_000 + */ gain
				}
			} else {
				scores[i] = history + /* 100_100_000 + */ int32(capPiece.Weight()-piece.Weight())
			}
			goto end
		}

		if promoType != NoType {
			p := GetPiece(promoType, White)
			scores[i] = history + 2000 + int32(p.Weight())
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

	var highestNonSpecialIndex int = -1
	var highestNonSpecialScore int32 = math.MinInt32
	engine := mp.engine
	scores := mp.quietMoveList.Scores
	moves := mp.quietMoveList.Moves
	size := mp.quietMoveList.Size

	var gpMove Move
	if mp.searchHeight > 1 {
		gpMove = engine.positionMoves[mp.searchHeight-1]
	}

	nextSpecialIndex := 0
	_ = scores[size-1]
	_ = moves[size-1]

	for i := 0; i < size; i++ {
		move := moves[i]

		if move == mp.hashmove {
			scores[i] = 900_000_000
			mp.quietMoveList.Swap(nextSpecialIndex, i)
			if highestNonSpecialIndex == nextSpecialIndex {
				highestNonSpecialIndex = i
			}
			nextSpecialIndex += 1
		} else if mp.killer1 == move || mp.killer2 == move {
			score := int32(80_000_000)
			if mp.killer1 == move {
				score = 90_000_000
			}
			scores[i] = score
			mp.quietMoveList.Swap(nextSpecialIndex, i)
			if highestNonSpecialIndex == nextSpecialIndex {
				highestNonSpecialIndex = i
			}
			nextSpecialIndex += 1
		} else if move == mp.counterMove {
			scores[i] = 70_000_000
			mp.quietMoveList.Swap(nextSpecialIndex, i)
			if highestNonSpecialIndex == nextSpecialIndex {
				highestNonSpecialIndex = i
			}
			nextSpecialIndex += 1
		} else {
			history := engine.searchHistory.QuietHistory(gpMove, mp.currentMove, move)
			scores[i] = history

			if highestNonSpecialScore < history {
				highestNonSpecialIndex = i
				highestNonSpecialScore = history
			}
		}
	}
	mp.quietMoveList.IsScored = true
	return highestNonSpecialIndex
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
	if mp.captureMoveList.Scores[bestIndex] < 0 && !mp.isQuiescence {
		alt := mp.getNextQuiet()
		if alt != EmptyMove {
			return alt
		}
	}
	best := mp.captureMoveList.Moves[bestIndex]
	mp.captureMoveList.Swap(next, bestIndex)
	mp.captureSees[next], mp.captureSees[bestIndex] = mp.captureSees[bestIndex], mp.captureSees[next]
	mp.captureMoveList.IncNext()
	return best
}

func (mp *MovePicker) getNextQuiet() Move {
	if mp.killerIndex == 1 {
		mp.killerIndex += 1
		if mp.position.IsPseudoLegal(mp.killer1) {
			mp.quietMoveList.IncNext()
			return mp.killer1
		}
	}

	if mp.killerIndex == 2 {
		mp.killerIndex += 1
		if mp.position.IsPseudoLegal(mp.killer2) {
			mp.quietMoveList.IncNext()
			return mp.killer2
		}
	}

	if mp.killerIndex == 3 {
		mp.killerIndex += 1
		if mp.position.IsPseudoLegal(mp.counterMove) {
			mp.quietMoveList.IncNext()
			return mp.counterMove
		}
	}

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
