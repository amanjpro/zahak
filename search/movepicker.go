package search

import (
	. "github.com/amanjpro/zahak/cache"
	. "github.com/amanjpro/zahak/engine"
)

type MovePicker struct {
	position  *Position
	engine    *Engine
	moves     []Move
	scores    []int32
	moveOrder int8
	next      int
}

func NewMovePicker(p *Position, e *Engine, moves []Move, moveOrder int8) *MovePicker {
	mp := &MovePicker{
		p,
		e,
		moves,
		make([]int32, len(moves)),
		moveOrder,
		0,
	}

	mp.score()
	return mp
}

func (mp *MovePicker) score() {
	pv := mp.engine.pv
	position := mp.position
	board := position.Board
	engine := mp.engine
	moveOrder := mp.moveOrder

	for i, move := range mp.moves {
		// Is in PV?
		if pv != nil && pv.moveCount > moveOrder {
			mv := pv.MoveAt(moveOrder)
			if mv == move {
				mp.scores[i] = 900_000_000
				continue
			}
		}

		// Is in Transition table ???
		// TODO: This is slow, that tells us either cache access is slow or has computation is
		// Or maybe (unlikely) make/unmake move is slow
		cp, ep, tg, hc := position.MakeMove(move)
		hash := position.Hash()
		position.UnMakeMove(move, tg, ep, cp, hc)
		eval, ok := TranspositionTable.Get(hash)

		if ok && eval.Type == Exact {
			mp.scores[i] = 500_000_000 + eval.Eval
			continue
		}

		if ok {
			mp.scores[i] = 400_000_000 + int32(eval.Depth)
			continue
		}

		piece := board.PieceAt(move.Source)
		//
		// capture ordering
		if move.HasTag(Capture) {
			capPiece := board.PieceAt(move.Destination)
			if !move.HasTag(EnPassant) {
				// SEE for ordering
				gain := board.StaticExchangeEval(move.Destination, capPiece, move.Source, piece)
				if gain < 0 {
					mp.scores[i] = -100_000_000 + gain
				} else if gain == 0 {
					mp.scores[i] = 100_000_000 + capPiece.Weight() - piece.Weight()
				} else {
					mp.scores[i] = 100_010_000 + gain
				}
			} else {
				p := GetPiece(move.PromoType, White)
				mp.scores[i] = 100_000_000 + capPiece.Weight() + p.Weight() - piece.Weight()
			}
			continue
		}

		killer := engine.KillerMoveScore(move, moveOrder)
		if killer != 0 {
			mp.scores[i] = killer
			continue
		}

		history := engine.MoveHistoryScore(piece, move.Destination, moveOrder)
		if history != 0 {
			mp.scores[i] = history
			continue
		}

		if move.PromoType != NoType {
			p := GetPiece(move.PromoType, White)
			mp.scores[i] = 50_000 + p.Weight()
			continue
		}

		// // prefer checks
		// if move.HasTag(Check) {
		// 	mp.scores[i] = 10_000
		// 	continue
		// }
		//
		// // King safety (castling)
		// castling := KingSideCastle | QueenSideCastle
		// moveIsCastling := move.HasTag(castling)
		// if moveIsCastling {
		// 	mp.scores[i] = 3_000
		// 	continue
		// }
		//
		// // Prefer smaller pieces
		// if piece.Type() == King {
		// 	mp.scores[i] = 0
		// 	continue
		// }
		//
		mp.scores[i] = 1000 //- piece.Weight()
	}
}

func (mp *MovePicker) Reset() {
	mp.next = 0
}

func (mp *MovePicker) Next() Move {
	if mp.next >= len(mp.moves) {
		return EmptyMove
	}
	var best Move = mp.moves[mp.next]
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
