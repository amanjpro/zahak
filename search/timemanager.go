package search

import (
	. "github.com/amanjpro/zahak/engine"
)

const COMMUNICATION_TIME_BUFFER = 50

func (e *Engine) InitiateTimer(game *Game, availableTimeInMillis int, isPerMove bool,
	increment int, movesToTimeControl int) int64 {
	maximumTimeToThink := 0
	if isPerMove {
		maximumTimeToThink = availableTimeInMillis
	} else {
		movestogo := 30
		if movesToTimeControl != 0 {
			movestogo = movesToTimeControl
		}
		availableTimeInMillis /= movestogo
		maximumTimeToThink = availableTimeInMillis
	}
	return int64(maximumTimeToThink - COMMUNICATION_TIME_BUFFER)
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
