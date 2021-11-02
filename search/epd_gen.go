package search

import (
	"fmt"
	"math/rand"
	"time"

	. "github.com/amanjpro/zahak/engine"
)

const startFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var randgen *rand.Rand

func init() {
	src := rand.NewSource(time.Now().UnixNano())
	randgen = rand.New(src)
}

func GenerateEpds() {
	cacheSize := 32
	TranspositionTable = NewCache(cacheSize)
	runner := NewRunner(1)
	runner.AddTimeManager(NewTimeManager(time.Now(), MAX_TIME, false, 0, 0, false))
	engine := runner.Engines[0]

	game := FromFen(startFen)
	engine.Position = game.Position()
	for true {
		gen(engine, 8)
	}
}

func gen(e *Engine, depthLeft int) {
	if depthLeft == 0 {
		eval := e.Position.Evaluate()
		if abs16(eval) < 700 {
			fmt.Printf("%s 8 ce %d;\n", e.Position.Fen(), eval)
		}
	} else {
		moves := e.Position.PseudoLegalMoves()
		for i := 0; i < 4*len(moves); i++ {
			move := moves[randgen.Intn(len(moves))]
			if oldEnPassant, oldTag, hc, ok := e.Position.MakeMove(move); ok {
				gen(e, depthLeft-1)
				e.Position.UnMakeMove(move, oldTag, oldEnPassant, hc)
				return
			}
		}
	}
}
