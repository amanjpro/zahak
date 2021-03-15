package search

import (
	. "github.com/amanjpro/zahak/engine"
)

type MovePicker struct {
	position  *Position
	engine    *Engine
	hashmove  Move
	moves     []Move
	scores    []int32
	hasPvMove bool
	moveOrder int8
	next      int
}

func NewMovePicker(p *Position, e *Engine, moves []Move, moveOrder int8, hashmove Move) *MovePicker {
	mp := &MovePicker{
		p,
		e,
		hashmove,
		moves,
		make([]int32, len(moves)),
		hashmove != EmptyMove,
		moveOrder,
		0,
	}

	mp.score()
	return mp
}

func (mp *MovePicker) HasPVMove() bool {
	return mp.hasPvMove
}

func (mp *MovePicker) UpgradeToPvMove(pvMove Move) {
	mp.hasPvMove = true
	for i, move := range mp.moves {
		if move == pvMove {
			mp.scores[i] = 900_000_000
			break
		}
	}
}

func (mp *MovePicker) score() {
	// pv := mp.engine.pv
	position := mp.position
	hashmove := mp.hashmove
	board := position.Board
	engine := mp.engine
	moveOrder := mp.moveOrder

	for i, move := range mp.moves {

		if move == hashmove {
			mp.scores[i] = 900_000_000
			continue
		}

		// // Is in PV?
		// if pv != nil && pv.moveCount > moveOrder {
		// 	mv := pv.MoveAt(moveOrder)
		// 	if mv == move {
		// 		mp.scores[i] = 500_000_000
		// 		mp.hasPvMove = true
		// 		continue
		// 	}
		// }
		//
		source := move.Source()
		dest := move.Destination()
		piece := move.MovingPiece()
		//
		// capture ordering
		if move.IsCapture() {
			capPiece := move.CapturedPiece()
			if !move.IsEnPassant() {
				// SEE for ordering
				gain := int32(board.StaticExchangeEval(dest, capPiece, source, piece))
				if gain < 0 {
					mp.scores[i] = -90_000_000 + gain
				} else if gain == 0 {
					mp.scores[i] = 100_000_000 + int32(capPiece.Weight()-piece.Weight())
				} else {
					mp.scores[i] = 100_100_000 + gain
				}
			} else {
				mp.scores[i] = 100_100_000 + int32(capPiece.Weight()-piece.Weight())
			}
			continue
		}

		killer := engine.KillerMoveScore(move, moveOrder)
		if killer != 0 {
			mp.scores[i] = killer
			continue
		}

		history := engine.MoveHistoryScore(piece, dest, moveOrder)
		if history != 0 {
			mp.scores[i] = history
			continue
		}

		promoType := move.PromoType()
		if promoType != NoType {
			p := GetPiece(promoType, White)
			mp.scores[i] = 50_000 + int32(p.Weight())
			continue
		}

		// prefer checks
		if move.IsCheck() {
			mp.scores[i] = 10_000
			continue
		}

		// King safety (castling)
		isCastling := move.IsKingSideCastle() || move.IsQueenSideCastle()
		if isCastling {
			mp.scores[i] = 3_000
			continue
		}

		// Prefer smaller pieces
		if piece.Type() == King {
			mp.scores[i] = 0
			continue
		}

		mp.scores[i] = 1000 - int32(piece.Weight())
	}
}

func (mp *MovePicker) Reset() {
	mp.next = 0
}

func (mp *MovePicker) Next() Move {

	if mp.next >= len(mp.moves) {
		return EmptyMove
	}

	var best = mp.moves[mp.next]
	var bestIndex = mp.next
	for i := mp.next + 1; i < len(mp.moves); i++ {
		move := mp.moves[i]
		if mp.scores[i] > mp.scores[bestIndex] {
			best = move
			bestIndex = i
		}
	}
	mp.moves[mp.next], mp.moves[bestIndex] = mp.moves[bestIndex], mp.moves[mp.next]
	mp.scores[mp.next], mp.scores[bestIndex] = mp.scores[bestIndex], mp.scores[mp.next]
	mp.next += 1
	return best
}
