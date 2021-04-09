package engine

import (
	"fmt"
)

type Game struct {
	position      *Position
	startPosition Position
	moves         []Move
	numberOfMoves uint16
}

func (g *Game) IsLegalMove(m Move) bool {
	// Very inefficient, but doesn't really matter
	for _, move := range g.position.LegalMoves() {
		if move.EqualTo(m) {
			return true
		}
	}
	return false
}

func (g *Game) Move(m Move) {
	pos := g.position

	if g.IsLegalMove(m) {
		g.moves = append(g.moves, m)
		pos.MakeMove(&m)
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
	startPosition Position,
	moves []Move,
	numberOfMoves uint16,
	clearCache bool) Game {

	if clearCache {
		initZobrist()
	}

	return Game{
		position,
		startPosition,
		moves,
		numberOfMoves,
	}
}
