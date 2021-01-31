package main

import (
	"fmt"
)

type Rank int8

const (
	Rank1 Rank = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
)

type File int8

const (
	FileA File = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
)

type Square struct {
	file File
	rank Rank
}

func (s *Square) BitboardIndex() int8 {
	return (int8(s.rank) * 8) + int8(s.file)
}

func (s *Square) Name() string {
	var fileName byte = ('a' + byte(s.file))
	var rankName int = (int(s.rank) + 1)
	return fmt.Sprint(string(fileName), rankName)
}

func SquareOf(file File, rank Rank) Square {
	if file == FileA {
		switch rank {
		case Rank1:
			return A1
		case Rank2:
			return A2
		case Rank3:
			return A3
		case Rank4:
			return A4
		case Rank5:
			return A5
		case Rank6:
			return A6
		case Rank7:
			return A7
		case Rank8:
			return A8
		}
	}

	if file == FileB {
		switch rank {
		case Rank1:
			return B1
		case Rank2:
			return B2
		case Rank3:
			return B3
		case Rank4:
			return B4
		case Rank5:
			return B5
		case Rank6:
			return B6
		case Rank7:
			return B7
		case Rank8:
			return B8
		}
	}

	if file == FileC {
		switch rank {
		case Rank1:
			return C1
		case Rank2:
			return C2
		case Rank3:
			return C3
		case Rank4:
			return C4
		case Rank5:
			return C5
		case Rank6:
			return C6
		case Rank7:
			return C7
		case Rank8:
			return C8
		}
	}

	if file == FileD {
		switch rank {
		case Rank1:
			return D1
		case Rank2:
			return D2
		case Rank3:
			return D3
		case Rank4:
			return D4
		case Rank5:
			return D5
		case Rank6:
			return D6
		case Rank7:
			return D7
		case Rank8:
			return D8
		}
	}

	if file == FileE {
		switch rank {
		case Rank1:
			return E1
		case Rank2:
			return E2
		case Rank3:
			return E3
		case Rank4:
			return E4
		case Rank5:
			return E5
		case Rank6:
			return E6
		case Rank7:
			return E7
		case Rank8:
			return E8
		}
	}

	if file == FileF {
		switch rank {
		case Rank1:
			return F1
		case Rank2:
			return F2
		case Rank3:
			return F3
		case Rank4:
			return F4
		case Rank5:
			return F5
		case Rank6:
			return F6
		case Rank7:
			return F7
		case Rank8:
			return F8
		}
	}

	if file == FileG {
		switch rank {
		case Rank1:
			return G1
		case Rank2:
			return G2
		case Rank3:
			return G3
		case Rank4:
			return G4
		case Rank5:
			return G5
		case Rank6:
			return G6
		case Rank7:
			return G7
		case Rank8:
			return G8
		}
	}

	if file == FileH {
		switch rank {
		case Rank1:
			return H1
		case Rank2:
			return H2
		case Rank3:
			return H3
		case Rank4:
			return H4
		case Rank5:
			return H5
		case Rank6:
			return H6
		case Rank7:
			return H7
		case Rank8:
			return H8
		}
	}

	return NoSquare // should never happen
}

var NoSquare Square = Square{-1, -1}
var A1 Square = Square{FileA, Rank1}
var A2 Square = Square{FileA, Rank2}
var A3 Square = Square{FileA, Rank3}
var A4 Square = Square{FileA, Rank4}
var A5 Square = Square{FileA, Rank5}
var A6 Square = Square{FileA, Rank6}
var A7 Square = Square{FileA, Rank7}
var A8 Square = Square{FileA, Rank8}
var B1 Square = Square{FileB, Rank1}
var B2 Square = Square{FileB, Rank2}
var B3 Square = Square{FileB, Rank3}
var B4 Square = Square{FileB, Rank4}
var B5 Square = Square{FileB, Rank5}
var B6 Square = Square{FileB, Rank6}
var B7 Square = Square{FileB, Rank7}
var B8 Square = Square{FileB, Rank8}
var C1 Square = Square{FileC, Rank1}
var C2 Square = Square{FileC, Rank2}
var C3 Square = Square{FileC, Rank3}
var C4 Square = Square{FileC, Rank4}
var C5 Square = Square{FileC, Rank5}
var C6 Square = Square{FileC, Rank6}
var C7 Square = Square{FileC, Rank7}
var C8 Square = Square{FileC, Rank8}
var D1 Square = Square{FileD, Rank1}
var D2 Square = Square{FileD, Rank2}
var D3 Square = Square{FileD, Rank3}
var D4 Square = Square{FileD, Rank4}
var D5 Square = Square{FileD, Rank5}
var D6 Square = Square{FileD, Rank6}
var D7 Square = Square{FileD, Rank7}
var D8 Square = Square{FileD, Rank8}
var E1 Square = Square{FileE, Rank1}
var E2 Square = Square{FileE, Rank2}
var E3 Square = Square{FileE, Rank3}
var E4 Square = Square{FileE, Rank4}
var E5 Square = Square{FileE, Rank5}
var E6 Square = Square{FileE, Rank6}
var E7 Square = Square{FileE, Rank7}
var E8 Square = Square{FileE, Rank8}
var F1 Square = Square{FileF, Rank1}
var F2 Square = Square{FileF, Rank2}
var F3 Square = Square{FileF, Rank3}
var F4 Square = Square{FileF, Rank4}
var F5 Square = Square{FileF, Rank5}
var F6 Square = Square{FileF, Rank6}
var F7 Square = Square{FileF, Rank7}
var F8 Square = Square{FileF, Rank8}
var G1 Square = Square{FileG, Rank1}
var G2 Square = Square{FileG, Rank2}
var G3 Square = Square{FileG, Rank3}
var G4 Square = Square{FileG, Rank4}
var G5 Square = Square{FileG, Rank5}
var G6 Square = Square{FileG, Rank6}
var G7 Square = Square{FileG, Rank7}
var G8 Square = Square{FileG, Rank8}
var H1 Square = Square{FileH, Rank1}
var H2 Square = Square{FileH, Rank2}
var H3 Square = Square{FileH, Rank3}
var H4 Square = Square{FileH, Rank4}
var H5 Square = Square{FileH, Rank5}
var H6 Square = Square{FileH, Rank6}
var H7 Square = Square{FileH, Rank7}
var H8 Square = Square{FileH, Rank8}
