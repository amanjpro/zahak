package engine

import (
	"fmt"
)

type Game struct {
	position      *Position
	moves         []Move
	numberOfMoves uint16
}

func (g *Game) Move(m Move) {
	pos := g.position

	if pos.IsPseudoLegal(m) {
		g.moves = append(g.moves, m)
		pos.GameMakeMove(m)
		if pos.Turn() == White {
			g.numberOfMoves += 1
		}
		v, ok := pos.Positions[pos.Hash()]
		if ok {
			pos.Positions[pos.Hash()] = v + 1
		} else {
			pos.Positions[pos.Hash()] = 1
		}
	} else {
		fmt.Printf("Illegal move, please try again: %s\n%s\n", m.ToString(), pos.Board.Draw())
	}
}

func (g *Game) Position() *Position {
	return g.position
}

func (g *Game) MoveClock() uint16 {
	return g.numberOfMoves
}

func NewGame(
	position *Position,
	moves []Move,
	numberOfMoves uint16) Game {

	return Game{
		position,
		moves,
		numberOfMoves,
	}
}
