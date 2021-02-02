package main

type Game struct {
	position      Position
	moves         []*Move
	positions     map[uint64]int8
	numberOfMoves uint16
	halfMoveClock uint16
}
