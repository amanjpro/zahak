package search

import (
	"bytes"

	. "github.com/amanjpro/zahak/engine"
)

type PVLine struct {
	moveCount int8   // Number of moves in the line.
	line      []Move // The line.
}

func NewPVLine(initialSize int8) PVLine {
	return PVLine{
		0,
		make([]Move, initialSize),
	}
}

func (thisLine *PVLine) ReplaceLine(firstMove Move, otherLine PVLine) {
	otherLineLen := int(otherLine.moveCount)
	thisLine.line[0] = firstMove
	for i := 0; i < otherLineLen; i++ {
		thisLine.line[i+1] = otherLine.line[i]
	}
	thisLine.moveCount = otherLine.moveCount + 1
}

func (thisLine *PVLine) Clone(otherLine PVLine) {
	otherLineLen := int(otherLine.moveCount)
	for i := 0; i < otherLineLen; i++ {
		thisLine.line[i] = otherLine.line[i]
	}
	thisLine.moveCount = otherLine.moveCount
}

func (thisLine *PVLine) AddFirst(move Move) {
	thisLine.line[0] = move
	thisLine.moveCount += 1
}

func (thisLine *PVLine) MoveAt(index int8) Move {
	return thisLine.line[index]
}

func (thisLine *PVLine) Recycle() {
	thisLine.moveCount = 0
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
