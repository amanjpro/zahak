package search

import (
	"time"

	. "github.com/amanjpro/zahak/engine"
)

func (e *Engine) InitiateTimer(game *Game, availableTimeInMillis int, isPerMove bool,
	increment int, movesToTimeControl int) {
	maximumTimeToThink := 0
	numberOfMovesOutOfBook := int(game.MoveClock()) // FIXME: Yup, fix it
	nMoves := min(numberOfMovesOutOfBook, 10)
	factor := 2 - nMoves/10
	if isPerMove {
		maximumTimeToThink = availableTimeInMillis - 100 + increment
	} else {
		if movesToTimeControl == 0 {
			mlh := max(60-int(game.MoveClock()), 20) // We assume that there are 60 more moves to go
			timeInMinute := time.Duration(availableTimeInMillis).Minutes()
			if timeInMinute <= 15 {
				mlh = max(50-int(game.MoveClock()), 20) // We assume that there are 60 more moves to go
			}
			if game.Position().IsEndGame() {
				movesToTimeControl = abs(mlh)
			} else {
				movesToTimeControl = abs(mlh + 10) // add 10 more moves in the early stage
			}
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
