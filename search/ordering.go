package search

import (
	"sort"

	. "github.com/amanjpro/zahak/cache"
	. "github.com/amanjpro/zahak/engine"
)

type MovePicker struct {
	position  *Position
	engine    *Engine
	moves     []*Move
	moveOrder int8
	ply       uint16
	next      int
}

func NewMovePicker(p *Position, e *Engine, moves []*Move, moveOrder int8, ply uint16) *MovePicker {
	return &MovePicker{
		p,
		e,
		moves,
		moveOrder,
		ply,
		0,
	}
}

func (mp *MovePicker) Reset() {
	mp.next = 0
}

func (mp *MovePicker) Next() *Move {
	var best *Move
	var bestIndex int
	for i := mp.next; i < len(mp.moves); i++ {
		move := mp.moves[i]
		if best == nil {
			best = move
			bestIndex = i
		} else if Better(move, best, mp.engine, mp.position, mp.moveOrder, mp.ply) {
			best = move
			bestIndex = i
		}
	}
	mp.moves[mp.next], mp.moves[bestIndex] = mp.moves[bestIndex], mp.moves[mp.next]
	mp.next += 1
	return best
}

func Better(move1 *Move, move2 *Move, engine *Engine, position *Position, moveOrder int8, ply uint16) bool {
	pv := engine.pv
	board := position.Board
	// Is in PV?
	if pv != nil && pv.moveCount > moveOrder {
		mv := pv.MoveAt(moveOrder)
		if *mv == *move1 {
			return true
		}
		if *mv == *move2 {
			return false
		}
	}

	// Is in Transition table ???
	// TODO: This is slow, that tells us either cache access is slow or has computation is
	// Or maybe (unlikely) make/unmake move is slow
	cp1, ep1, tg1, hc1 := position.MakeMove(move1)
	hash1 := position.Hash()
	position.UnMakeMove(move1, tg1, ep1, cp1, hc1)
	eval1, ok1 := TranspositionTable.Get(hash1)

	cp2, ep2, tg2, hc2 := position.MakeMove(move2)
	hash2 := position.Hash()
	position.UnMakeMove(move2, tg2, ep2, cp2, hc2)
	eval2, ok2 := TranspositionTable.Get(hash2)

	if ok1 && ok2 {
		// if eval1.Depth > eval2.Depth {
		// 	return true
		// } else if eval1.Depth < eval2.Depth {
		// 	return false
		// }
		if eval1.Type == Exact && eval2.Type != Exact {
			return true
		} else if eval2.Type == Exact && eval1.Type != Exact {
			return false
		}
		// if eval1.Eval > eval2.Eval {
		// 	return true
		// } else if eval1.Eval < eval2.Eval {
		// 	return false
		// }
	} else if ok1 {
		if eval1.Type == Exact {
			return true
		}
	} else if ok2 {
		if eval2.Type == Exact {
			return false
		}
	}

	killer1 := engine.KillerMoveScore(move1, ply)
	killer2 := engine.KillerMoveScore(move2, ply)
	if killer1 > killer2 {
		return true
	} else if killer2 > killer2 {
		return false
	}

	piece1 := board.PieceAt(move1.Source)
	piece2 := board.PieceAt(move2.Source)
	//
	// capture ordering
	if move1.HasTag(Capture) && move2.HasTag(Capture) {
		capPiece1 := board.PieceAt(move1.Destination)
		capPiece2 := board.PieceAt(move2.Destination)
		if !move1.HasTag(EnPassant) && !move2.HasTag(EnPassant) {
			// SEE for ordering
			gain1 := board.StaticExchangeEval(move1.Destination, capPiece1, move1.Source, piece1)
			gain2 := board.StaticExchangeEval(move2.Destination, capPiece2, move2.Source, piece2)

			if gain1 < 0 && gain2 >= 0 {
				return false
			}
			if gain1 >= 0 && gain2 < 0 {
				return true
			}
		} else if !move1.HasTag(EnPassant) {
			// SEE for ordering
			gain1 := board.StaticExchangeEval(move1.Destination, capPiece1, move1.Source, piece1)
			return gain1 > 0
		} else if !move2.HasTag(EnPassant) {
			// SEE for ordering
			gain2 := board.StaticExchangeEval(move2.Destination, capPiece2, move2.Source, piece2)
			return gain2 < 0
		}

		// What are we capturing?
		if capPiece1.Type() > capPiece2.Type() {
			return true
		}
		if capPiece2.Type() > capPiece1.Type() {
			return false
		}

		// Who is capturing?
		if piece1.Type() < piece2.Type() {
			return true
		}
		if piece2.Type() < piece1.Type() {
			return false
		}
	} else if move1.HasTag(Capture) {
		capPiece1 := board.PieceAt(move1.Destination)
		gain1 := board.StaticExchangeEval(move1.Destination, capPiece1, move1.Source, piece1)
		return gain1 > 0
	} else if move2.HasTag(Capture) {
		capPiece2 := board.PieceAt(move2.Destination)
		gain2 := board.StaticExchangeEval(move2.Destination, capPiece2, move2.Source, piece2)
		return gain2 < 0
	}

	// prefer checks
	if move1.HasTag(Check) {
		return true
	}
	if move2.HasTag(Check) {
		return false
	}
	// Prefer smaller pieces
	if piece1.Type() < piece2.Type() {
		return true
	}

	// King safety (castling)
	castling := KingSideCastle | QueenSideCastle
	move1IsCastling := move1.HasTag(castling)
	move2IsCastling := move2.HasTag(castling)
	if move1IsCastling && !move2IsCastling {
		return true
	} else if move2IsCastling && !move1IsCastling {
		return false
	}

	if move1.PromoType != NoType && move2.PromoType == NoType {
		return true
	} else if move2.PromoType != NoType && move1.PromoType == NoType {
		return false
	} else if move2.PromoType != NoType && move1.PromoType != NoType {
		p1 := GetPiece(move1.PromoType, White)
		p2 := GetPiece(move2.PromoType, White)
		return p1.Weight() > p2.Weight()
	}

	return false
}

type IterationMoves struct {
	moves  []*Move
	engine *Engine
	evals  []int32
}

func (iter *IterationMoves) Len() int {
	return len(iter.moves)
}

func (iter *IterationMoves) Swap(i, j int) {
	evals := iter.evals
	moves := iter.moves
	moves[i], moves[j] = moves[j], moves[i]
	evals[i], evals[j] = evals[j], evals[i]
}

func (iter *IterationMoves) Less(i, j int) bool {
	pv := iter.engine.pv
	eval1, eval2 := iter.evals[i], iter.evals[j]
	equal := eval1 == eval2
	if equal {
		move1, move2 := iter.moves[i], iter.moves[j]
		// Is in PV?
		if pv != nil && pv.moveCount > 0 {
			mv := pv.MoveAt(0)
			if *mv == *move1 {
				return true
			}
			if *mv == *move2 {
				return false
			}
		}
	}
	return eval1 > eval2
}

func orderIterationMoves(iter *IterationMoves) []*Move {
	sort.Sort(iter)
	return iter.moves
}
