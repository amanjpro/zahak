package search

import (
	"time"

	. "github.com/amanjpro/zahak/engine"
)

func InitiateTimer(game *Game, availableTimeInMillis int, isPerMove bool,
	increment int, movesToTimeControl int, stopTimer chan bool) {
	maximumTimeToThink := 0
	if isPerMove {
		maximumTimeToThink = availableTimeInMillis - 10 + increment
	} else {
		if movesToTimeControl == 0 {
			mlh := 30
			timeInMinute := time.Duration(availableTimeInMillis).Minutes()
			if timeInMinute <= 15 {
				mlh = 20
			}
			if game.Position().IsEndGame() {
				movesToTimeControl = abs(mlh - int(game.MoveClock()))
			} else {
				movesToTimeControl = abs(mlh + 10 - int(game.MoveClock()))
			}
		}
		if game.Position().IsEndGame() {
			maximumTimeToThink = max(availableTimeInMillis/movesToTimeControl, 1000) // Naiive, but works for now
		} else {
			maximumTimeToThink = max(availableTimeInMillis/movesToTimeControl, 1000) // Naiive, but works for now
		}
	}

	select {
	case <-stopTimer:
		break
	case <-time.After(time.Duration(maximumTimeToThink) * time.Millisecond):
		STOP_SEARCH_GLOBALLY = true
		break
	}
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
