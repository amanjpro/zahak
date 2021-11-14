package search

import (
	. "github.com/amanjpro/zahak/engine"
)

const HistoryMax int32 = 397
const HistoryMultiplier = 32
const HistoryDivisor = 512

type MoveHistory struct {
	killers         [MAX_DEPTH][2]Move
	history         [2 * 64][64]int32
	counters        [2 * 64][64]Move
	counterHistory  [12 * 64][12 * 64]int32
	followupHistory [12 * 64][12 * 64]int32
}

func (mh *MoveHistory) Reset() {

	for i := 0; i < len(mh.killers); i++ {
		for j := 0; j < len(mh.killers[i]); j++ {
			mh.killers[i][j] = EmptyMove
		}
	}

	for i := 0; i < len(mh.counters); i++ {
		for j := 0; j < len(mh.counters[i]); j++ {
			mh.counters[i][j] = EmptyMove
		}
	}

	for i := 0; i < len(mh.counterHistory); i++ {
		for j := 0; j < len(mh.counterHistory[i]); j++ {
			mh.counterHistory[i][j] = 0
			mh.followupHistory[i][j] = 0
		}
	}
}

func (m *MoveHistory) History(stm Color, gpMove Move, pMove Move, move Move) int32 {
	msrc := int(move.Source())
	mdest := int(move.Destination())
	mpiece := int(move.MovingPiece() - 1)
	value := m.history[int(stm)*64+msrc][mdest]
	if pMove != EmptyMove {
		pdest := int(move.Destination())
		ppiece := int(pMove.MovingPiece() - 1)
		value += m.counterHistory[ppiece*64+pdest][mpiece*64+mdest]
	}
	if gpMove != EmptyMove {
		gpiece := int(gpMove.MovingPiece() - 1)
		gdest := int(gpMove.Destination())
		value += m.followupHistory[gpiece*64+gdest][mpiece*64+mdest]
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
	return m.counters[int(stm)*64+int(previousMove.Source())][previousMove.Destination()]
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
			entry := m.history[int(stm)*64+src][dest]
			entry += HistoryMultiplier*signedBonus - entry*unsignedBonus/HistoryDivisor
			m.history[int(stm)*64+src][dest] = entry

			if pMove != EmptyMove {
				entry = m.counterHistory[ppiece*64+pdest][mpiece*64+dest]
				entry += HistoryMultiplier*signedBonus - entry*unsignedBonus/HistoryDivisor
				m.counterHistory[ppiece*64+pdest][mpiece*64+dest] = entry
			}

			if gpMove != EmptyMove {
				entry = m.followupHistory[gpiece*64+gdest][mpiece*64+dest]
				entry += HistoryMultiplier*signedBonus - entry*unsignedBonus/HistoryDivisor
				m.followupHistory[gpiece*64+gdest][mpiece*64+dest] = entry
			}
		}

		if pMove != EmptyMove {
			m.counters[int(stm)*64+psrc][pdest] = move
		}
	}
}

func min32(x int32, y int32) int32 {
	if x < y {
		return x
	}
	return y
}
