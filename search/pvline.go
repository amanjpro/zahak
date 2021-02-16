package search

import (
	"bytes"

	. "github.com/amanjpro/zahak/engine"
)

type PVLine struct {
	moveCount uint16  // Number of moves in the line.
	line      []*Move // The line.
}

func NewPVLine(initialSize int8) *PVLine {
	return &PVLine{
		0,
		make([]*Move, initialSize),
	}
}

func (thisLine *PVLine) ReplaceLine(otherLine *PVLine) {
	thisLine.moveCount += otherLine.moveCount
	copy(thisLine.line[1:], otherLine.line)
}

func (thisLine *PVLine) AddFirst(move *Move) {
	thisLine.moveCount++
	thisLine.line[0] = move
}

func (pv *PVLine) ToString() string {
	var buffer bytes.Buffer

	for i := uint16(0); i < pv.moveCount; i++ {
		if i != 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString(pv.line[i].ToString())
	}
	return buffer.String()
}
