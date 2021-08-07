package search

import (
	"sync/atomic"
	"time"
)

const COMMUNICATION_TIME_BUFFER = 50

// Implements this: http://talkchess.com/forum3/viewtopic.php?f=7&t=77396&p=894325&hilit=cold+turkey#p894294
type TimeManager struct {
	StartTime           time.Time
	HardLimit           int64
	SoftLimit           int64
	NodesSinceLastCheck int64
	AbruptStop          bool
	StopSearchNow       atomic.Value
	IsPerMove           bool
	ExtensionCounter    int
	pondering           atomic.Value
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

	var s atomic.Value
	var p atomic.Value
	s.Store(false)
	p.Store(pondering)

	return &TimeManager{
		HardLimit:           hardLimit,
		SoftLimit:           softLimit,
		StartTime:           startTime,
		NodesSinceLastCheck: 0,
		AbruptStop:          false,
		StopSearchNow:       s,
		IsPerMove:           isPerMove,
		ExtensionCounter:    0,
		pondering:           p,
	}
}

func (tm *TimeManager) Pondering() bool {
	return tm.pondering.Load().(bool)
}

func (tm *TimeManager) ShouldStop(isRoot bool, canCutNow bool) bool {
	if tm.Pondering() {
		return false
	}
	stopSearchNow := tm.StopSearchNow.Load().(bool)
	if tm.NodesSinceLastCheck < 2000 {
		tm.NodesSinceLastCheck += 1
		tm.AbruptStop = tm.AbruptStop || stopSearchNow
		return tm.AbruptStop
	}
	tm.NodesSinceLastCheck = 0
	if isRoot && canCutNow {
		return stopSearchNow || time.Since(tm.StartTime).Milliseconds() >= 2*tm.SoftLimit
	} else {
		tm.AbruptStop = tm.AbruptStop || stopSearchNow || time.Since(tm.StartTime).Milliseconds() >= tm.HardLimit
		return tm.AbruptStop
	}
}

func (tm *TimeManager) UpdatePondering(flag bool) {
	tm.pondering.Store(flag)
}

func (tm *TimeManager) UpdateStopSearchNow(flag bool) {
	tm.StopSearchNow.Store(flag)
}

func (tm *TimeManager) CanStartNewIteration() bool {
	if tm.Pondering() {
		return true
	}
	stopSearchNow := tm.StopSearchNow.Load().(bool)
	if tm.AbruptStop || stopSearchNow {
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
	if tm.Pondering() {
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
