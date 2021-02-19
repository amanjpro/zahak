package engine

import (
	"fmt"

	. "github.com/amanjpro/zahak/cache"
)

type Game struct {
	position      *Position
	startPosition Position
	moves         []*Move
	numberOfMoves uint16
}

func (g *Game) IsLegalMove(m *Move) bool {
	// Very inefficient, but doesn't really matter
	for _, move := range g.position.LegalMoves() {
		if *move == *m {
			return true
		}
	}
	return false
}

func (g *Game) Move(m *Move) {
	pos := g.position

	if g.IsLegalMove(m) {
		g.numberOfMoves += 1
		g.moves = append(g.moves, m)
		pos.MakeMove(m)
	} else {
		fmt.Printf("Illegal move, please try again: %s\n%s\n", m.ToString(), pos.Board.Draw())
	}
}

func (g *Game) Status() Status {
	return g.position.Status()
}

func (g *Game) Position() *Position {
	return g.position
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
	numberOfMoves uint16,
	clearCache bool) Game {

	if clearCache {
		initZobrist()
		ResetCache()
	}

	return Game{
		position,
		startPosition,
		moves,
		numberOfMoves,
	}
}
