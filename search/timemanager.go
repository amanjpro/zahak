package search

import (
	"time"
)

const COMMUNICATION_TIME_BUFFER = 50
const MAX_TIME int64 = 922_337_203_685_477_580

// Implements this: http://talkchess.com/forum3/viewtopic.php?f=7&t=77396&p=894325&hilit=cold+turkey#p894294
type TimeManager struct {
	StartTime           time.Time
	HardLimit           int64
	SoftLimit           int64
	NodesSinceLastCheck int64
	AbruptStop          bool
	StopSearchNow       bool
	IsPerMove           bool
	ExtensionCounter    int
	Pondering           bool
}

func NewTimeManager(startTime time.Time, availableTimeInMillis int64, isPerMove bool,
	increment int64, movesToTimeControl int64, pondering bool) *TimeManager {
	softLimit := int64(0)
	hardLimit := int64(0)
	if isPerMove {
		softLimit = availableTimeInMillis - COMMUNICATION_TIME_BUFFER
		hardLimit = softLimit
	} else {
		movestogo := int64(30)
		if movesToTimeControl > 0 {
			movestogo = movesToTimeControl
		}
		softLimit = availableTimeInMillis / movestogo
		softLimit = min64(int64(softLimit+increment-COMMUNICATION_TIME_BUFFER), availableTimeInMillis-COMMUNICATION_TIME_BUFFER)
		hardLimit = min64(softLimit*10, availableTimeInMillis-COMMUNICATION_TIME_BUFFER)
	}

	return &TimeManager{
		HardLimit:           hardLimit,
		SoftLimit:           softLimit,
		StartTime:           startTime,
		NodesSinceLastCheck: 0,
		AbruptStop:          false,
		StopSearchNow:       false,
		IsPerMove:           isPerMove,
		ExtensionCounter:    0,
		Pondering:           pondering,
	}
}

func (tm *TimeManager) ShouldStop(isRoot bool, canCutNow bool) bool {
	if tm.Pondering {
		return false
	}
	if tm.NodesSinceLastCheck < 2000 {
		tm.NodesSinceLastCheck += 1
		tm.AbruptStop = tm.AbruptStop || tm.StopSearchNow
		return tm.AbruptStop
	}
	tm.NodesSinceLastCheck = 0
	if isRoot && canCutNow {
		tm.AbruptStop = tm.AbruptStop || tm.StopSearchNow || time.Since(tm.StartTime).Milliseconds() >= 2*tm.SoftLimit
		return tm.AbruptStop
	} else {
		tm.AbruptStop = tm.AbruptStop || tm.StopSearchNow || time.Since(tm.StartTime).Milliseconds() >= tm.HardLimit
		return tm.AbruptStop
	}
}

func (tm *TimeManager) CanStartNewIteration() bool {
	if tm.Pondering {
		return true
	}
	if tm.AbruptStop || tm.StopSearchNow {
		return false
	}

	if tm.IsPerMove {
		return time.Since(tm.StartTime).Milliseconds() <= tm.SoftLimit
	} else {
		limit := 70 * tm.SoftLimit / 100
		return time.Since(tm.StartTime).Milliseconds() <= limit
	}
}

func (tm *TimeManager) ExtraTime() {
	if tm.Pondering {
		return
	}
	if tm.ExtensionCounter < 5 {
		tm.SoftLimit = min64(tm.HardLimit, tm.SoftLimit+tm.SoftLimit/10)
		tm.ExtensionCounter += 1
	}
}

func min64(a int64, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
