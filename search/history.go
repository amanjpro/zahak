package search

import (
	. "github.com/amanjpro/zahak/engine"
)

const HistoryMax int32 = 397
const HistoryMultiplier = 32
const HistoryDivisor = 512

type MoveHistory struct {
	killers         [MAX_DEPTH][2]Move
	history         [2][64][64]int32
	counters        [2][64][64]Move
	counterHistory  [12][64][12][64]int32
	followupHistory [12][64][12][64]int32
}

func (m *MoveHistory) History(stm Color, gpMove Move, pMove Move, move Move) int32 {
	msrc := move.Source()
	mdest := move.Destination()
	pdest := move.Destination()
	mpiece := move.MovingPiece() - 1
	ppiece := pMove.MovingPiece() - 1
	value := m.history[stm][msrc][mdest]
	if pMove != EmptyMove {
		value += m.counterHistory[ppiece][pdest][mpiece][mdest]
	}
	if gpMove != EmptyMove {
		gpiece := gpMove.MovingPiece() - 1
		gdest := gpMove.Destination()
		value += m.followupHistory[gpiece][gdest][mpiece][mdest]
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

func (m *MoveHistory) AddHistory(move Move, pMove Move, gpMove Move, depthLeft int8, searchHeight int8, stm Color, moves []Move) {
	if depthLeft >= 0 && move.PromoType() == NoType && !move.IsCapture() {

		if m.killers[searchHeight][0] != move && move != EmptyMove {
			m.killers[searchHeight][1], m.killers[searchHeight][0] = m.killers[searchHeight][0], move
		}

		if depthLeft <= 1 {
			return
		}

		unsignedBonus := min32(int32(depthLeft)*int32(depthLeft), HistoryMax)

		psrc := pMove.Source()
		pdest := pMove.Destination()
		ppiece := pMove.MovingPiece() - 1
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

			if pMove != EmptyMove {
				entry = m.counterHistory[ppiece][pdest][mpiece][dest]
				entry += HistoryMultiplier*signedBonus - entry*unsignedBonus/HistoryDivisor
				m.counterHistory[ppiece][pdest][mpiece][dest] = entry
			}

			if gpMove != EmptyMove {
				gpiece := gpMove.MovingPiece() - 1
				gdest := gpMove.Destination()
				entry = m.followupHistory[gpiece][gdest][mpiece][dest]
				entry += HistoryMultiplier*signedBonus - entry*unsignedBonus/HistoryDivisor
				m.followupHistory[gpiece][gdest][mpiece][dest] = entry
			}
		}

		if pMove != EmptyMove {
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
