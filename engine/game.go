package engine

type Game struct {
	position      *Position
	startPosition Position
	moves         []Move
	numberOfMoves uint16
}

func (g *Game) Move(m Move) {
	pos := g.position

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
