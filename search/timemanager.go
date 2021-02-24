package search

import (
	"sync"
	"time"

	. "github.com/amanjpro/zahak/engine"
)

func (e *Engine) InitiateTimer(game *Game, availableTimeInMillis int, isPerMove bool,
	increment int, movesToTimeControl int, done *sync.WaitGroup, stopTimer chan bool) {
	maximumTimeToThink := 0
	if isPerMove {
		maximumTimeToThink = availableTimeInMillis - 10 + increment
	} else {
		if movesToTimeControl == 0 {
			mlh := 60 // We assume that there are 60 more moves to go
			timeInMinute := time.Duration(availableTimeInMillis).Minutes()
			if timeInMinute <= 15 {
				mlh = 50 // shorter games have shorter moves, hopefully
			}
			if game.Position().IsEndGame() {
				movesToTimeControl = abs(mlh)
			} else {
				movesToTimeControl = abs(mlh + 10) // add 10 more moves in the early stage
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
		e.StopSearchFlag = true
		break
	}
	done.Done()
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
