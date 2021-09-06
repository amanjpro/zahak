package search

import (
	"fmt"
	"time"

	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

const startFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func GenerateEpds() {
	cacheSize := uint32(32)
	runner := NewRunner(NewCache(cacheSize), NewPawnCache(1), 1)
	runner.AddTimeManager(NewTimeManager(time.Now(), MAX_TIME, false, 0, 0, false))
	engine := runner.Engines[0]

	game := FromFen(startFen)
	engine.Position = game.Position()
	gen(engine, 8)
}

func gen(e *Engine, depthLeft int) {
	if depthLeft == 0 {
		eval := Evaluate(e.Position, e.Pawnhash)
		if abs16(eval) < 700 {
			fmt.Printf("%s 8 ce %d;\n", e.Position.Fen(), eval)
		}
	} else {
		moves := e.Position.PseudoLegalMoves()
		for _, move := range moves {
			if oldEnPassant, oldTag, hc, ok := e.Position.MakeMove(move); ok {
				gen(e, depthLeft-1)
				e.Position.UnMakeMove(move, oldTag, oldEnPassant, hc)
			}
		}
	}
}
