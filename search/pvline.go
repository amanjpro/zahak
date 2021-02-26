package search

import (
	"bytes"

	. "github.com/amanjpro/zahak/engine"
)

type PVLine struct {
	moveCount int8   // Number of moves in the line.
	line      []Move // The line.
	hasFirst  bool
}

func NewPVLine(initialSize int8) *PVLine {
	return &PVLine{
		0,
		make([]Move, initialSize),
		false,
	}
}

func (thisLine *PVLine) ReplaceLine(otherLine *PVLine) {
	otherLineLen := int(otherLine.moveCount)
	if thisLine.hasFirst {
		thisLine.moveCount = 1
	} else {
		thisLine.moveCount = 0
	}
	for i := 0; i < otherLineLen; i++ {
		thisLine.moveCount += 1
		thisLine.line[i+1] = otherLine.line[i]
	}
}

func (thisLine *PVLine) AddFirst(move Move) {
	if !thisLine.hasFirst {
		thisLine.moveCount += 1
		thisLine.hasFirst = true
	}
	thisLine.line[0] = move
}

func (thisLine *PVLine) MoveAt(index int8) Move {
	return thisLine.line[index]
}

func (pv *PVLine) Pop() Move {
	var toReturn Move
	if pv.moveCount >= 0 {
		emptySlice := make([]Move, len(pv.line))
		mv, newSlice := pv.line[0], pv.line[1:]
		toReturn = mv
		copy(emptySlice, newSlice)
		pv.line = emptySlice
		pv.moveCount -= 1
	}
	return toReturn
}

func (pv *PVLine) ToString() string {
	var buffer bytes.Buffer

	for i := int8(0); i < pv.moveCount; i++ {
		if i != 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString(pv.line[i].ToString())
	}
	return buffer.String()
}
