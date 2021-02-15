package engine

import (
	"fmt"

	. "github.com/amanjpro/zahak/cache"
)

type Game struct {
	position      *Position
	startPosition Position
	moves         []*Move
	positions     map[uint64]int8
	numberOfMoves uint16
	halfMoveClock uint8
}

// TODO: Implement me
func (g *Game) IsLegalMove(m *Move) bool {
	return true
}

func (g *Game) Move(m *Move) {
	pos := g.position

	if g.IsLegalMove(m) {
		board := pos.Board
		movingPiece := board.PieceAt(m.Source)
		if m.HasTag(Capture) || movingPiece.Type() == Pawn {
			g.halfMoveClock = 0
		} else {
			g.halfMoveClock += 1
		}
		g.numberOfMoves += 1
		g.moves = append(g.moves, m)
		pos.MakeMove(m)
		hash := pos.Hash()
		_, ok := g.positions[hash]
		if ok {
			g.positions[hash] += 1
		} else {
			g.positions[hash] = 0
		}
	} else {
		fmt.Printf("Illegal move, please try again: %s\n%s\n", m.ToString(), pos.Board.Draw())
	}
}

func (g *Game) Status() Status {
	if g.halfMoveClock >= 100 {
		return Draw
	}
	for _, c := range g.positions {
		if c >= 3 {
			fmt.Println(c)
			return Draw
		}
	}
	return g.position.Status()
}

func (g *Game) Position() *Position {
	return g.position
}

func (g *Game) HalfMoveClock() uint8 {
	return g.halfMoveClock
}

func (g *Game) MoveClock() uint16 {
	if g.position.Turn() == Black {
		return g.numberOfMoves + 1
	}
	return g.numberOfMoves
}

func NewGame(
	position *Position,
	startPosition Position,
	moves []*Move,
	positions map[uint64]int8,
	numberOfMoves uint16,
	halfMoveClock uint8,
	clearCache bool) Game {

	if clearCache {
		initZobrist()
		ResetCache()
		NewQCache(2)
	}

	return Game{
		position,
		startPosition,
		moves,
		positions,
		numberOfMoves,
		halfMoveClock,
	}
}
