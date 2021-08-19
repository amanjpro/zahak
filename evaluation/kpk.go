package evaluation

import (
	"math/bits"

	. "github.com/amanjpro/zahak/engine"
)

// This function expects to have only 1 P, and two Kings are on the board
func KpkProbe(board *Bitboard, winningSide Color, turn Color) int16 {

	var winningKingSq, losingKingSq, pawnSq Square
	if winningSide == White {
		winningKingSq = Square(bits.LeadingZeros64(board.GetBitboardOf(WhiteKing)))
		losingKingSq = Square(bits.LeadingZeros64(board.GetBitboardOf(BlackKing)))
		pawnSq = Square(bits.LeadingZeros64(board.GetBitboardOf(WhitePawn)))
	} else {
		winningKingSq = Square(bits.LeadingZeros64(board.GetBitboardOf(BlackKing)))
		losingKingSq = Square(bits.LeadingZeros64(board.GetBitboardOf(WhiteKing)))
		pawnSq = Square(bits.LeadingZeros64(board.GetBitboardOf(BlackPawn)))
	}

	if KpkIsDraw(winningSide, turn, winningKingSq, losingKingSq, pawnSq) {
		return 0
	}

	res := MAX_NON_CHECKMATE - 1000 // don't want the engine to not queen
	if turn == winningSide {
		return res
	}
	return -res
}

// Following is taken from Cheng
func GetKPKBit(bit uint32) uint8 {
	return kpkBitbase[bit>>3] & (1 << (bit & 7))
}

// stm: white = king with pawn
func KPKIndex(winningKingSq Square, losingKingSq Square, pawnSq Square, turn Color) uint32 {
	file := pawnSq.File()
	// mirror horizontally if necessary
	xm := Square(0)
	if file >= FileE {
		xm = 7
	}
	winningKingSq ^= xm
	losingKingSq ^= xm
	pawnSq ^= xm
	file ^= xm.File()
	// now we can build index
	pp := uint32(((pawnSq&0x38)-8)>>1) | uint32(file)
	return uint32(winningKingSq) | (uint32(losingKingSq) << 6) | (uint32(turn) << 12) | (pp << 13)
}

// is draw?
// color (with pawn), color (stm), king (with pawn) position, bare king position, pawn position
func KpkIsDraw(winningSide Color, turn Color, winningKingSq Square, losingKingSq Square, pawnSq Square) bool {
	xm := Square(0)
	if winningSide != White {
		xm = Square(0x38)
	}
	return GetKPKBit(KPKIndex(winningKingSq^xm, losingKingSq^xm, pawnSq^xm, winningSide^turn)) != 0
}
