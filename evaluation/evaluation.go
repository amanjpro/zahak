package evaluation

import (
	"math/bits"

	. "github.com/amanjpro/zahak/engine"
)

func Evaluate(position *Position) int32 {
	// board := position.Board
	// allPieces := board.AllPieces()
	return middlegameEval(position)
}

const CHECKMATE_EVAL int32 = 400_000

// Piece Square Tables
var pawnPst = []int32{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 15, 15, 0, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 0, 5, 5, 0, 0, 0,
	0, 0, 0, -25, -25, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var knightPst = []int32{
	-40, -25, -25, -25, -25, -25, -25, -40,
	-30, 0, 0, 0, 0, 0, 0, -30,
	-30, 0, 0, 0, 0, 0, 0, -30,
	-30, 0, 0, 15, 15, 0, 0, -30,
	-30, 0, 0, 15, 15, 0, 0, -30,
	-30, 0, 10, 0, 0, 10, 0, -30,
	-30, 0, 0, 5, 5, 0, 0, -30,
	-40, -30, -25, -25, -25, -25, -30, -40,
}

var bishopPst = []int32{
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-10, 0, 5, 0, 0, 5, 0, -10,
	-10, 0, 0, 10, 10, 0, 0, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 0, 5, 0, 0, 5, 0, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-10, -20, -20, -20, -20, -20, -20, -10,
}

var rookPst = []int32{
	0, 0, 0, 0, 0, 0, 0, 0,
	10, 10, 10, 10, 10, 10, 10, 10,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 5, 5, 0, 0, 0,
}

var queenPst = []int32{
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	5, 5, 10, 10, 10, 10, 5, 5,
	5, 5, 10, 15, 15, 10, 5, 5,
}

var kingPst = []int32{
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	20, 25, 25, -15, -15, 20, 25, 20,
}

var flip = []int32{
	56, 57, 58, 59, 60, 61, 62, 63,
	48, 49, 50, 51, 52, 53, 54, 55,
	40, 41, 42, 43, 44, 45, 46, 47,
	32, 33, 34, 35, 36, 37, 38, 39,
	24, 25, 26, 27, 28, 29, 30, 31,
	16, 17, 18, 19, 20, 21, 22, 23,
	8, 9, 10, 11, 12, 13, 14, 15,
	0, 1, 2, 3, 4, 5, 6, 7,
}

func middlegameEval(position *Position) int32 {
	board := position.Board
	p := BlackPawn
	n := BlackKnight
	b := BlackBishop
	r := BlackRook
	q := BlackQueen
	turn := position.Turn()

	// Compute material balance
	bbBlackPawn := board.GetBitboardOf(BlackPawn)
	bbBlackKnight := board.GetBitboardOf(BlackKnight)
	bbBlackBishop := board.GetBitboardOf(BlackBishop)
	bbBlackRook := board.GetBitboardOf(BlackRook)
	bbBlackQueen := board.GetBitboardOf(BlackQueen)
	bbBlackKing := board.GetBitboardOf(BlackKing)

	bbWhitePawn := board.GetBitboardOf(WhitePawn)
	bbWhiteKnight := board.GetBitboardOf(WhiteKnight)
	bbWhiteBishop := board.GetBitboardOf(WhiteBishop)
	bbWhiteRook := board.GetBitboardOf(WhiteRook)
	bbWhiteQueen := board.GetBitboardOf(WhiteQueen)
	bbWhiteKing := board.GetBitboardOf(WhiteKing)

	blackPawnsCount := int32(0)
	blackKnightsCount := int32(0)
	blackBishopsCount := int32(0)
	blackRooksCount := int32(0)
	blackQueensCount := int32(0)

	whitePawnsCount := int32(0)
	whiteKnightsCount := int32(0)
	whiteBishopsCount := int32(0)
	whiteRooksCount := int32(0)
	whiteQueensCount := int32(0)

	blackCentipawns := int32(0)
	whiteCentipawns := int32(0)
	whites := board.GetWhitePieces()
	blacks := board.GetBlackPieces()
	all := whites | blacks

	// PST for black pawns
	pieceIter := bbBlackPawn
	blackPawnsPerFile := [8]int8{0, 0, 0, 0, 0, 0, 0, 0}
	blackLeastAdvancedPawnsPerFile := [8]Rank{Rank1, Rank1, Rank1, Rank1, Rank1, Rank1, Rank1, Rank1}
	blackMostAdvancedPawnsPerFile := [8]Rank{Rank8, Rank8, Rank8, Rank8, Rank8, Rank8, Rank8, Rank8}
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		blackPawnsCount++
		// backwards pawn
		if board.IsBackwardPawn(mask, bbBlackPawn, Black) {
			blackCentipawns -= 25
		}
		// pawn map
		sq := Square(mask)
		file := sq.File()
		rank := sq.Rank()
		blackPawnsPerFile[int(file)] += 1
		if rank > blackLeastAdvancedPawnsPerFile[file] {
			blackLeastAdvancedPawnsPerFile[int(file)] = rank
		}
		if rank < blackMostAdvancedPawnsPerFile[file] {
			blackMostAdvancedPawnsPerFile[int(file)] = rank
		}
		blackCentipawns += pawnPst[flip[index]]
		pieceIter ^= mask
	}

	// PST for white pawns
	pieceIter = bbWhitePawn
	whitePawnsPerFile := [8]int8{0, 0, 0, 0, 0, 0, 0, 0}
	whiteLeastAdvancedPawnsPerFile := [8]Rank{Rank8, Rank8, Rank8, Rank8, Rank8, Rank8, Rank8, Rank8}
	whiteMostAdvancedPawnsPerFile := [8]Rank{Rank1, Rank1, Rank1, Rank1, Rank1, Rank1, Rank1, Rank1}
	for pieceIter != 0 {
		whitePawnsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		// backwards pawn
		if board.IsBackwardPawn(mask, bbBlackPawn, White) {
			whiteCentipawns -= 25
		}
		// pawn map
		sq := Square(mask)
		file := sq.File()
		rank := sq.Rank()
		whitePawnsPerFile[int(file)] += 1
		if rank < whiteLeastAdvancedPawnsPerFile[file] {
			whiteLeastAdvancedPawnsPerFile[int(file)] = rank
		}
		if rank > whiteMostAdvancedPawnsPerFile[file] {
			whiteMostAdvancedPawnsPerFile[int(file)] = rank
		}
		whiteCentipawns += pawnPst[index]
		pieceIter ^= mask
	}

	for i := 0; i < 8; i++ {
		// isolated pawn penalty - white
		if whitePawnsPerFile[i] > 0 {
			isIsolated := false
			if i == 0 && whitePawnsPerFile[i+1] <= 0 {
				isIsolated = true
			} else if i == 7 && whitePawnsPerFile[i-1] <= 0 {
				isIsolated = true
			} else if whitePawnsPerFile[i-1] <= 0 && whitePawnsPerFile[i+1] <= 0 {
				isIsolated = true
			}
			if isIsolated {
				whiteCentipawns -= 35
			}
		}

		// isolated pawn penalty - black
		if blackPawnsPerFile[i] > 0 {
			isIsolated := false
			if i == 0 && blackPawnsPerFile[i+1] <= 0 {
				isIsolated = true
			} else if i == 7 && blackPawnsPerFile[i-1] <= 0 {
				isIsolated = true
			} else if blackPawnsPerFile[i-1] <= 0 && blackPawnsPerFile[i+1] <= 0 {
				isIsolated = true
			}
			if isIsolated {
				blackCentipawns -= 35
			}
		}
		// double pawn penalty - black
		if blackPawnsPerFile[i] > 1 {
			blackCentipawns -= 35
		}
		// double pawn penalty - white
		if whitePawnsPerFile[i] > 1 {
			whiteCentipawns -= 35
		}
		// passed and candidate passed pawn award
		rank := whiteMostAdvancedPawnsPerFile[i]
		if rank != Rank1 {
			if blackLeastAdvancedPawnsPerFile[i] == Rank8 || blackLeastAdvancedPawnsPerFile[i] < rank { // candidate
				if i == 0 {
					if blackLeastAdvancedPawnsPerFile[i+1] == Rank8 || blackLeastAdvancedPawnsPerFile[i+1] < rank { // passed pawn
						whiteCentipawns += 50 //passed pawn
						if rank >= Rank5 {
							whiteCentipawns += int32(rank) * 50 // advanced pawns are better
						}
					} else {
						whiteCentipawns += 25 // candidate passed pawn
					}
				} else if i == 7 {
					if blackLeastAdvancedPawnsPerFile[i-1] == Rank8 || blackLeastAdvancedPawnsPerFile[i-1] < rank { // passed pawn
						whiteCentipawns += 50 //passed pawn
						if rank >= Rank5 {
							whiteCentipawns += int32(rank) * 50 // advanced pawns are better
						}
					} else {
						whiteCentipawns += 25 // candidate passed pawn
					}
				} else {
					if (blackLeastAdvancedPawnsPerFile[i-1] == Rank8 || blackLeastAdvancedPawnsPerFile[i-1] < rank) &&
						(blackLeastAdvancedPawnsPerFile[i+1] == Rank8 || blackLeastAdvancedPawnsPerFile[i+1] < rank) { // passed pawn
						whiteCentipawns += 50 //passed pawn
						if rank >= Rank5 {
							whiteCentipawns += int32(rank) * 50 // advanced pawns are better
						}
					} else {
						whiteCentipawns += 25 // candidate passed pawn
					}
				}
			}
		}

		rank = blackMostAdvancedPawnsPerFile[i]
		if rank != Rank8 {
			if whiteLeastAdvancedPawnsPerFile[i] == Rank1 || whiteLeastAdvancedPawnsPerFile[i] > rank { // candidate
				if i == 0 {
					if whiteLeastAdvancedPawnsPerFile[i+1] == Rank1 || whiteLeastAdvancedPawnsPerFile[i+1] > rank { // passed pawn
						blackCentipawns += 50 //passed pawn
						if rank <= Rank4 {
							blackCentipawns += int32(rank) * 50 // advanced pawns are better
						}
					} else {
						blackCentipawns += 25 // candidate passed pawn
					}
				} else if i == 7 {
					if whiteLeastAdvancedPawnsPerFile[i-1] == Rank1 || whiteLeastAdvancedPawnsPerFile[i-1] > rank { // passed pawn
						blackCentipawns += 50 //passed pawn
						if rank <= Rank4 {
							blackCentipawns += int32(rank) * 50 // advanced pawns are better
						}
					} else {
						blackCentipawns += 25 // candidate passed pawn
					}
				} else {
					if (whiteLeastAdvancedPawnsPerFile[i-1] == Rank1 || whiteLeastAdvancedPawnsPerFile[i-1] > rank) &&
						(whiteLeastAdvancedPawnsPerFile[i+1] == Rank1 || whiteLeastAdvancedPawnsPerFile[i+1] > rank) { // passed pawn
						blackCentipawns += 50 //passed pawn
						if rank <= Rank4 {
							blackCentipawns += int32(rank) * 50 // advanced pawns are better
						}
					} else {
						blackCentipawns += 25 // candidate passed pawn
					}
				}
			}
		}
	}

	// PST for other black pieces
	pieceIter = bbBlackKnight
	for pieceIter != 0 {
		blackKnightsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		blackCentipawns += knightPst[flip[index]]
		pieceIter ^= mask
	}

	pieceIter = bbBlackBishop
	for pieceIter != 0 {
		blackBishopsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		blackCentipawns += bishopPst[flip[index]]
		pieceIter ^= mask
	}

	pieceIter = bbBlackRook
	for pieceIter != 0 {
		blackRooksCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		file := Square(index).File()
		if blackPawnsPerFile[file] == 0 {
			if whitePawnsPerFile[file] == 0 { // open file
				blackCentipawns += 50
			} else { // semi-open file
				blackCentipawns += 25
			}
		}
		sq := Square(index)
		if board.IsVerticalDoubleRook(sq, bbBlackRook, all) {
			// double-rook vertical
			blackCentipawns += 25
		} else if board.IsHorizontalDoubleRook(sq, bbBlackRook, all) {
			// double-rook horizontal
			blackCentipawns += 15
		}
		blackCentipawns += rookPst[flip[index]]
		pieceIter ^= mask
	}

	pieceIter = bbBlackQueen
	for pieceIter != 0 {
		blackQueensCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		blackCentipawns += queenPst[flip[index]]
		pieceIter ^= mask
	}

	pieceIter = bbBlackKing
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		award := kingPst[flip[index]]
		if award <= 0 {
			if !position.HasTag(BlackCanCastleKingSide) {
				award -= 10
			} else if !position.HasTag(BlackCanCastleQueenSide) {
				award -= 10
			}
		}
		blackCentipawns += award
		pieceIter ^= mask
	}

	// PST for other white pieces
	pieceIter = bbWhiteKnight
	for pieceIter != 0 {
		whiteKnightsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		whiteCentipawns += knightPst[index]
		pieceIter ^= mask
	}

	pieceIter = bbWhiteBishop
	for pieceIter != 0 {
		whiteBishopsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		whiteCentipawns += bishopPst[index]
		pieceIter ^= mask
	}

	pieceIter = bbWhiteRook
	for pieceIter != 0 {
		whiteRooksCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		file := Square(index).File()
		if whitePawnsPerFile[file] == 0 {
			if blackPawnsPerFile[file] == 0 { // open file
				whiteCentipawns += 50
			} else { // semi-open file
				whiteCentipawns += 25
			}
		}
		sq := Square(index)
		if board.IsVerticalDoubleRook(sq, bbWhiteRook, all) {
			// double-rook vertical
			whiteCentipawns += 25
		} else if board.IsHorizontalDoubleRook(sq, bbWhiteRook, all) {
			// double-rook horizontal
			whiteCentipawns += 15
		}
		whiteCentipawns += rookPst[index]
		pieceIter ^= mask
	}

	pieceIter = bbWhiteQueen
	for pieceIter != 0 {
		whiteQueensCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		whiteCentipawns += queenPst[index]
		pieceIter ^= mask
	}

	pieceIter = bbWhiteKing
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		award := kingPst[index]
		if award < 0 {
			if !position.HasTag(WhiteCanCastleKingSide) {
				award -= 10
			} else if !position.HasTag(WhiteCanCastleQueenSide) {
				award -= 10
			}
		}
		whiteCentipawns += award
		pieceIter ^= mask
	}

	blackCentipawns += blackPawnsCount * p.Weight()
	blackCentipawns += blackKnightsCount * n.Weight()
	blackCentipawns += blackBishopsCount * b.Weight()
	blackCentipawns += blackRooksCount * r.Weight()
	blackCentipawns += blackQueensCount * q.Weight()

	whiteCentipawns += whitePawnsCount * p.Weight()
	whiteCentipawns += whiteKnightsCount * n.Weight()
	whiteCentipawns += whiteBishopsCount * b.Weight()
	whiteCentipawns += whiteRooksCount * r.Weight()
	whiteCentipawns += whiteQueensCount * q.Weight()

	// 2 Bishops vs 2 Knights
	if whiteBishopsCount >= 2 && blackBishopsCount < 2 {
		whiteCentipawns += 25
	}
	if whiteBishopsCount < 2 && blackBishopsCount >= 2 {
		blackCentipawns += 25
	}

	// mobility and attacks
	whiteAttacks := board.AllAttacks(Black) // get the squares that are taboo for black (white's reach)
	blackAttacks := board.AllAttacks(White) // get the squares that are taboo for whtie (black's reach)
	wAttackCounts := bits.OnesCount64(whiteAttacks)
	bAttackCounts := bits.OnesCount64(blackAttacks)

	whiteAggressivity := bits.OnesCount64(whiteAttacks >> 32) // keep hi-bits only (black's half)
	blackAggressivity := bits.OnesCount64(blackAttacks << 32) // keep lo-bits only (white's half)

	whiteCentipawns += int32(wAttackCounts - bAttackCounts)
	blackCentipawns += int32(bAttackCounts - wAttackCounts)

	whiteCentipawns += int32(2 * (whiteAggressivity - blackAggressivity))
	blackCentipawns += int32(2 * (blackAggressivity - whiteAggressivity))

	if turn == White {
		return whiteCentipawns - blackCentipawns
	} else {
		return blackCentipawns - whiteCentipawns
	}
}
