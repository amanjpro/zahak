package search

import (
	. "github.com/amanjpro/zahak/engine"
)

const HistoryMax int32 = 397
const HistoryMultiplier = 32
const HistoryDivisor = 512

type MoveHistory struct {
	killers         [][]Move
	history         [][]int32
	counters        [][]Move
	counterHistory  [][]int32
	followupHistory [][]int32
}

func NewMoveHistory() MoveHistory {
	mh := MoveHistory{}

	mh.killers = make([][]Move, MAX_DEPTH)
	for i := 0; i < len(mh.killers); i++ {
		mh.killers[i] = make([]Move, 2)
	}

	mh.history = make([][]int32, 2*64)
	mh.counters = make([][]Move, 2*64)
	for i := 0; i < len(mh.counters); i++ {
		mh.history[i] = make([]int32, 2*64)
		mh.counters[i] = make([]Move, 2*64)
	}

	mh.counterHistory = make([][]int32, 12*64)
	mh.followupHistory = make([][]int32, 12*64)
	for i := 0; i < len(mh.counterHistory); i++ {
		mh.counterHistory[i] = make([]int32, 12*64)
		mh.followupHistory[i] = make([]int32, 12*64)
	}

	return mh
}

func (m *MoveHistory) History(stm Color, gpMove Move, pMove Move, move Move) int32 {
	msrc := int(move.Source())
	mdest := int(move.Destination())
	mpiece := int(move.MovingPiece() - 1)
	value := m.history[int(stm)*msrc][mdest]
	if pMove != EmptyMove {
		pdest := int(move.Destination())
		ppiece := int(pMove.MovingPiece() - 1)
		value += m.counterHistory[ppiece*pdest][mpiece*mdest]
	}
	if gpMove != EmptyMove {
		gpiece := int(gpMove.MovingPiece() - 1)
		gdest := int(gpMove.Destination())
		value += m.followupHistory[gpiece*gdest][mpiece*mdest]
	}
	return value
}

func (m *MoveHistory) KillerMoveAt(searchHeight int8) (Move, Move) {
	if searchHeight < 0 {
		return EmptyMove, EmptyMove
	}
	return m.killers[searchHeight][0], m.killers[searchHeight][1]
}

func (m *MoveHistory) CounterMoveAt(stm Color, previousMove Move) Move {
	if previousMove == EmptyMove {
		return EmptyMove
	}
	return m.counters[int(stm)*int(previousMove.Source())][previousMove.Destination()]
}

func (m *MoveHistory) AddHistory(move Move, pMove Move, gpMove Move, depthLeft int8, searchHeight int8, stm Color, moves []Move) {
	if depthLeft >= 0 && move.PromoType() == NoType && !move.IsCapture() {

		if m.killers[searchHeight][0] != move && move != EmptyMove {
			m.killers[searchHeight][1], m.killers[searchHeight][0] = m.killers[searchHeight][0], move
		}

		if depthLeft <= 1 {
			return
		}

		unsignedBonus := min32(int32(depthLeft)*int32(depthLeft), HistoryMax)

		psrc := int(pMove.Source())
		pdest := int(pMove.Destination())
		ppiece := int(pMove.MovingPiece() - 1)
		gpiece := int(gpMove.MovingPiece() - 1)
		gdest := int(gpMove.Destination())
		for _, mv := range moves {
			src := int(mv.Source())
			dest := int(mv.Destination())
			mpiece := int(mv.MovingPiece() - 1)

			var signedBonus int32
			if move == mv {
				signedBonus = unsignedBonus
			} else {
				signedBonus = -unsignedBonus
			}
			entry := m.history[int(stm)*src][dest]
			entry += HistoryMultiplier*signedBonus - entry*unsignedBonus/HistoryDivisor
			m.history[int(stm)*src][dest] = entry

			if pMove != EmptyMove {
				entry = m.counterHistory[ppiece*pdest][mpiece*dest]
				entry += HistoryMultiplier*signedBonus - entry*unsignedBonus/HistoryDivisor
				m.counterHistory[ppiece*pdest][mpiece*dest] = entry
			}

			if gpMove != EmptyMove {
				entry = m.followupHistory[gpiece*gdest][mpiece*dest]
				entry += HistoryMultiplier*signedBonus - entry*unsignedBonus/HistoryDivisor
				m.followupHistory[gpiece*gdest][mpiece*dest] = entry
			}
		}

		if pMove != EmptyMove {
			m.counters[int(stm)*psrc][pdest] = move
		}
	}
}

func min32(x int32, y int32) int32 {
	if x < y {
		return x
	}
	return y
}
