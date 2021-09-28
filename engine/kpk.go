package engine

import (
	"math/bits"
)

var WinningKpk = WhiteQueen.Weight() - WhitePawn.Weight()

var Flip = [64]int16{
	56, 57, 58, 59, 60, 61, 62, 63,
	48, 49, 50, 51, 52, 53, 54, 55,
	40, 41, 42, 43, 44, 45, 46, 47,
	32, 33, 34, 35, 36, 37, 38, 39,
	24, 25, 26, 27, 28, 29, 30, 31,
	16, 17, 18, 19, 20, 21, 22, 23,
	8, 9, 10, 11, 12, 13, 14, 15,
	0, 1, 2, 3, 4, 5, 6, 7,
}

// This function expects to have only 1 P, and two Kings are on the board
func KpkProbe(board *Bitboard, winningSide Color, turn Color) int16 {

	var winningKingSq, losingKingSq, pawnSq, actualPsq Square
	// I am flipping the board, to adapt to the KPK bitbase, as it is designed for
	// an upside down board
	if winningSide == White {
		winningKingSq = Square(Flip[bits.TrailingZeros64(board.GetBitboardOf(WhiteKing))])
		losingKingSq = Square(Flip[bits.TrailingZeros64(board.GetBitboardOf(BlackKing))])
		pindex := bits.TrailingZeros64(board.GetBitboardOf(WhitePawn))
		actualPsq = Square(pindex)
		pawnSq = Square(Flip[pindex])
	} else {
		winningKingSq = Square(Flip[bits.TrailingZeros64(board.GetBitboardOf(BlackKing))])
		losingKingSq = Square(Flip[bits.TrailingZeros64(board.GetBitboardOf(WhiteKing))])
		pindex := bits.TrailingZeros64(board.GetBitboardOf(BlackPawn))
		actualPsq = Square(pindex)
		pawnSq = Square(Flip[pindex])
	}

	if KpkIsDraw(winningSide, turn, winningKingSq, losingKingSq, pawnSq) {
		return 0
	}

	res := WinningKpk
	if winningSide == White {
		distance := int16(Rank8 - actualPsq.Rank())
		res -= distance
	} else {
		distance := int16(actualPsq.Rank() - Rank1)
		res -= distance
	}
	if turn != winningSide {
		res = -res
	}
	return res
}

// Following is taken from Cheng
func GetKPKBit(bit uint32) uint32 {
	return uint32(kpkBitbase[bit>>3] & (1 << (bit & 7)))
}

func KPKIndex(winningKingSq Square, losingKingSq Square, pawnSq Square, winningSide Color) uint32 {
	file := pawnSq.File()
	// mirror horizontally if necessary
	var xm Square
	if file > FileD {
		xm = 7
	} else {
		xm = 0
	}

	winningKingSq ^= xm
	losingKingSq ^= xm
	pawnSq ^= xm

	c := Square(0x38)
	// now we can build index
	pp := (((pawnSq & c) - 8) >> 1) | Square(file) ^ xm

	return uint32(winningKingSq) | (uint32(losingKingSq) << 6) | (uint32(winningSide) << 12) | (uint32(pp) << 13)
}

// is draw?
// color (with pawn), color (stm), king (with pawn) position, bare king position, pawn position
func KpkIsDraw(winningSide Color, turn Color, winningKingSq Square, losingKingSq Square, pawnSq Square) bool {
	var xm Square
	if winningSide == White {
		xm = 0
	} else {
		xm = 0x38
	}

	index := KPKIndex(winningKingSq^xm, losingKingSq^xm, pawnSq^xm, winningSide^turn)
	res := GetKPKBit(index)
	return res != 0
}
