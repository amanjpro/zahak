package search

import (
	. "github.com/amanjpro/zahak/engine"
)

const HistoryMax int32 = 612
const HistoryMultiplier = 32
const HistoryDivisor = 512

type MoveHistory struct {
	killers        [][]Move
	history        [][][]int32
	counters       [][][]Move
	counterHistory [][][][]int32
}

func NewMoveHistory() MoveHistory {
	mh := MoveHistory{}

	mh.killers = make([][]Move, MAX_DEPTH)
	for i := 0; i < len(mh.killers); i++ {
		mh.killers[i] = make([]Move, 2)
	}

	mh.history = make([][][]int32, 2)
	mh.counters = make([][][]Move, 2)
	for i := 0; i < len(mh.history); i++ {
		mh.history[i] = make([][]int32, 64)
		mh.counters[i] = make([][]Move, 64)
		for j := 0; j < len(mh.history[i]); j++ {
			mh.history[i][j] = make([]int32, 64)
			mh.counters[i][j] = make([]Move, 64)
		}
	}

	mh.counterHistory = make([][][][]int32, 12)
	for i := 0; i < len(mh.counterHistory); i++ {
		mh.counterHistory[i] = make([][][]int32, 64)
		for j := 0; j < len(mh.counterHistory[i]); j++ {
			mh.counterHistory[i][j] = make([][]int32, 12)
			for k := 0; k < len(mh.counterHistory[i][j]); k++ {
				mh.counterHistory[i][j][k] = make([]int32, 64)
			}
		}
	}
	return mh
}

func (m *MoveHistory) History(stm Color, previousMove Move, move Move) int32 {
	msrc := move.Source()
	mdest := move.Destination()
	pdest := move.Destination()
	mpiece := move.MovingPiece() - 1
	ppiece := previousMove.MovingPiece() - 1
	value := m.history[stm][msrc][mdest]
	if previousMove != EmptyMove {
		value += m.counterHistory[ppiece][pdest][mpiece][mdest]
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
	return m.counters[stm][previousMove.Source()][previousMove.Destination()]
}

func (m *MoveHistory) AddHistory(move Move, previousMove Move, depthLeft int8, searchHeight int8, stm Color, moves []Move) {
	if depthLeft >= 0 && move.PromoType() == NoType && !move.IsCapture() {

		if m.killers[searchHeight][0] != move && move != EmptyMove {
			m.killers[searchHeight][1], m.killers[searchHeight][0] = m.killers[searchHeight][0], move
		}

		if depthLeft <= 1 {
			return
		}

		unsignedBonus := min32(int32(depthLeft)*int32(depthLeft), HistoryMax)

		psrc := previousMove.Source()
		pdest := previousMove.Destination()
		ppiece := previousMove.MovingPiece() - 1
		for _, mv := range moves {
			src := mv.Source()
			dest := mv.Destination()
			mpiece := mv.MovingPiece() - 1

			var signedBonus int32
			if move == mv {
				signedBonus = unsignedBonus
			} else {
				signedBonus = -unsignedBonus
			}
			entry := m.history[stm][src][dest]
			entry += HistoryMultiplier*signedBonus - entry*unsignedBonus/HistoryDivisor
			m.history[stm][src][dest] = entry

			if previousMove != EmptyMove {
				entry = m.counterHistory[ppiece][pdest][mpiece][dest]
				entry += HistoryMultiplier*signedBonus - entry*unsignedBonus/HistoryDivisor
				m.counterHistory[ppiece][pdest][mpiece][dest] = entry
			}
		}

		if previousMove != EmptyMove {
			m.counters[stm][psrc][pdest] = move
		}
	}
}

func min32(x int32, y int32) int32 {
	if x < y {
		return x
	}
	return y
}
