package search

import (
	. "github.com/amanjpro/zahak/engine"
)

type MoveHistory struct {
	killers         [MAX_DEPTH][2]Move
	history         [12][64]int32
	counters        [12][64]Move
	counterHistory  [12 * 64][12 * 64]int32
	followupHistory [12 * 64][12 * 64]int32
}

func (m *MoveHistory) History(gpMove Move, pMove Move, move Move) int32 {
	mdest := int(move.Destination())
	mpiece := int(move.MovingPiece() - 1)
	value := m.history[mpiece][mdest]
	if pMove != EmptyMove {
		pdest := int(pMove.Destination())
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

func (m *MoveHistory) CounterMoveAt(previousMove Move) Move {
	if previousMove == EmptyMove {
		return EmptyMove
	}
	return m.counters[previousMove.MovingPiece()-1][previousMove.Destination()]
}

func historyBonus(current int32, bonus int32) int32 {
	return current + 32*bonus - current*abs32(bonus)/512
}

func (m *MoveHistory) AddHistory(move Move, pMove Move, gpMove Move, depthLeft int8, searchHeight int8, moves []Move) {
	if move != EmptyMove && depthLeft >= 0 && move.PromoType() == NoType && !move.IsCapture() {

		if m.killers[searchHeight][0] != move {
			m.killers[searchHeight][1], m.killers[searchHeight][0] = m.killers[searchHeight][0], move
		}

		if depthLeft > 1 {
			bonus := int32(depthLeft) * int32(depthLeft)

			pdest := int(pMove.Destination())
			ppiece := int(pMove.MovingPiece() - 1)
			gdest := int(gpMove.Destination())
			gpiece := int(gpMove.MovingPiece() - 1)
			for i := 0; i < len(moves); i++ {
				mv := moves[i]
				// src := int(mv.Source())

				if move != mv {
					mdest := int(mv.Destination())
					mpiece := int(mv.MovingPiece() - 1)
					m.history[mpiece][mdest] = historyBonus(m.history[mpiece][mdest], -bonus)

					if pMove != EmptyMove {
						m.counterHistory[ppiece*64+pdest][mpiece*64+mdest] = historyBonus(m.counterHistory[ppiece*64+pdest][mpiece*64+mdest], -bonus)
					}

					if gpMove != EmptyMove {
						m.followupHistory[gpiece*64+gdest][mpiece*64+mdest] = historyBonus(m.followupHistory[gpiece*64+gdest][mpiece*64+mdest], -bonus)
					}
				}
			}

			piece := int(move.MovingPiece() - 1)
			dest := int(move.Destination())
			m.history[piece][dest] = historyBonus(m.history[piece][dest], bonus)

			if pMove != EmptyMove {
				m.counterHistory[ppiece*64+pdest][piece*64+dest] = historyBonus(m.counterHistory[ppiece*64+pdest][piece*64+dest], bonus)
				m.counters[ppiece][pdest] = move
			}

			if gpMove != EmptyMove {
				m.followupHistory[gpiece*64+gdest][piece*64+dest] = historyBonus(m.followupHistory[gpiece*64+gdest][piece*64+dest], bonus)
			}
		}
	}
}

func min32(x int32, y int32) int32 {
	if x < y {
		return x
	}
	return y
}
