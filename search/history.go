package search

import (
	. "github.com/amanjpro/zahak/engine"
)

type MoveHistory struct {
	killers        [MAX_DEPTH][2]Move
	history        [2][64][64]int32
	counters       [2][64][64]Move
	counterHistory [12][64][12][64]int32
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

func historyBonus(current int32, depthLeft int8, coeff int32) int32 {
	bonus := int32(depthLeft) * int32(depthLeft) * coeff
	return current + 32*bonus - current*abs32(bonus)/512
}

func (m *MoveHistory) AddHistory(move Move, previousMove Move, depthLeft int8, searchHeight int8, stm Color, moves []Move) {
	if depthLeft >= 0 && move.PromoType() == NoType && !move.IsCapture() {

		if m.killers[searchHeight][0] != move {
			m.killers[searchHeight][1] = m.killers[searchHeight][0]
			m.killers[searchHeight][0] = move
		}

		if depthLeft <= 1 {
			return
		}

		src := move.Source()
		dest := move.Destination()
		psrc := previousMove.Source()
		pdest := previousMove.Destination()
		mpiece := move.MovingPiece() - 1
		ppiece := previousMove.MovingPiece() - 1

		for _, move := range moves {
			var coeff int32 = -1
			if move == move {
				coeff = 1
			}
			entry := m.history[stm][src][dest]
			m.history[stm][src][dest] = historyBonus(entry, depthLeft, coeff)

			if previousMove != EmptyMove {
				entry = m.counterHistory[ppiece][pdest][mpiece][dest]
				m.counterHistory[ppiece][pdest][mpiece][dest] = historyBonus(entry, depthLeft, coeff)
			}
		}

		if previousMove != EmptyMove {
			m.counters[stm][psrc][pdest] = move
		}
	}
}
