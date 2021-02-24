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
	scores    []int32
	moveOrder int8
	ply       uint16
	next      int
}

func NewMovePicker(p *Position, e *Engine, moves []*Move, moveOrder int8, ply uint16) *MovePicker {
	mp := &MovePicker{
		p,
		e,
		moves,
		make([]int32, len(moves)),
		moveOrder,
		ply,
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
	ply := mp.ply

	for i, move := range mp.moves {
		// Is in PV?
		if pv != nil && pv.moveCount > moveOrder {
			mv := pv.MoveAt(moveOrder)
			if *mv == *move {
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
				} else {
					mp.scores[i] = 100_000_000 + gain
				}
			} else {
				mp.scores[i] = 100_000_100
			}
			continue
		}

		killer := engine.KillerMoveScore(move, ply)
		if killer != 0 {
			mp.scores[i] = killer
			continue
		}

		history := engine.MoveHistoryScore(piece, move.Destination, ply)
		if history != 0 {
			mp.scores[i] = history
			continue
		}

		if move.PromoType != NoType {
			p := GetPiece(move.PromoType, White)
			mp.scores[i] = 50_000 + p.Weight()
			continue
		}

		// prefer checks
		if move.HasTag(Check) {
			mp.scores[i] = 10_000
			continue
		}

		// King safety (castling)
		castling := KingSideCastle | QueenSideCastle
		moveIsCastling := move.HasTag(castling)
		if moveIsCastling {
			mp.scores[i] = 3_000
			continue
		}

		// Prefer smaller pieces
		if piece.Type() == King {
			mp.scores[i] = 0
			continue
		}
		mp.scores[i] = 1000 - piece.Weight()
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
		} else if mp.scores[i] > mp.scores[bestIndex] {
			best = move
			bestIndex = i
		}
	}
	mp.moves[mp.next], mp.moves[bestIndex] = mp.moves[bestIndex], mp.moves[mp.next]
	mp.scores[mp.next], mp.scores[bestIndex] = mp.scores[bestIndex], mp.scores[mp.next]
	mp.next += 1
	return best
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
