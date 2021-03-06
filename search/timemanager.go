package search

import (
	. "github.com/amanjpro/zahak/engine"
)

func (e *Engine) InitiateTimer(game *Game, availableTimeInMillis int, isPerMove bool,
	increment int, movesToTimeControl int) {
	maximumTimeToThink := 0
	numberOfMovesOutOfBook := int(game.MoveClock()) // / 2 // FIXME: Yup, fix it
	nMoves := min(numberOfMovesOutOfBook, 10)
	availableTimeInMillis += increment
	factor := 2 - nMoves/10
	if isPerMove {
		maximumTimeToThink = availableTimeInMillis - 100
	} else {
		if movesToTimeControl == 0 {
			mlh := max(50-int(game.MoveClock()), 20) // We assume that there are 40 moves to go
			movesToTimeControl = mlh
		}

		target := availableTimeInMillis / movesToTimeControl
		maximumTimeToThink = factor * target
	}

	e.ThinkTime = int64(maximumTimeToThink)
}

func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
