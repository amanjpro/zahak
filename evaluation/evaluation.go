package evaluation

import (
	"math/bits"

	. "github.com/amanjpro/zahak/engine"
)

const CHECKMATE_EVAL int32 = 400_000

// Piece Square Tables

// Middle game
var earlyPawnPst = []int32{
	0, 0, 0, 0, 0, 0, 0, 0,
	98, 134, 61, 95, 68, 126, 34, -11,
	-6, 7, 26, 31, 65, 56, 25, -20,
	-14, 13, 6, 21, 23, 12, 17, -23,
	-27, -2, -5, 12, 17, 6, 10, -25,
	-26, -4, -4, -10, 3, 3, 33, -12,
	-35, -1, -20, -23, -15, 24, 38, -22,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var earlyKnightPst = []int32{
	-167, -89, -34, -49, 61, -97, -15, -107,
	-73, -41, 72, 36, 23, 62, 7, -17,
	-47, 60, 37, 65, 84, 129, 73, 44,
	-9, 17, 19, 53, 37, 69, 18, 22,
	-13, 4, 16, 13, 28, 19, 21, -8,
	-23, -9, 12, 10, 19, 17, 25, -16,
	-29, -53, -12, -3, -1, 18, -14, -19,
	-105, -21, -58, -33, -17, -28, -19, -23,
}

var earlyBishopPst = []int32{
	-29, 4, -82, -37, -25, -42, 7, -8,
	-26, 16, -18, -13, 30, 59, 18, -47,
	-16, 37, 43, 40, 35, 50, 37, -2,
	-4, 5, 19, 50, 37, 37, 7, -2,
	-6, 13, 13, 26, 34, 12, 10, 4,
	0, 15, 15, 15, 14, 27, 18, 10,
	4, 15, 16, 0, 7, 21, 33, 1,
	-33, -3, -14, -21, -13, -12, -39, -21,
}

var earlyRookPst = []int32{
	32, 42, 32, 51, 63, 9, 31, 43,
	27, 32, 58, 62, 80, 67, 26, 44,
	-5, 19, 26, 36, 17, 45, 61, 16,
	-24, -11, 7, 26, 24, 35, -8, -20,
	-36, -26, -12, -1, 9, -7, 6, -23,
	-45, -25, -16, -17, 3, 0, -5, -33,
	-44, -16, -20, -9, -1, 11, -6, -71,
	-19, -13, 1, 17, 16, 7, -37, -26,
}

var earlyQueenPst = []int32{
	-28, 0, 29, 12, 59, 44, 43, 45,
	-24, -39, -5, 1, -16, 57, 28, 54,
	-13, -17, 7, 8, 29, 56, 47, 57,
	-27, -27, -16, -16, -1, 17, -2, 1,
	-9, -26, -9, -10, -2, -4, 3, -3,
	-14, 2, -11, -2, -5, 2, 14, 5,
	-35, -8, 11, 2, 8, 15, -3, 1,
	-1, -18, -9, 10, -15, -25, -31, -50,
}

var earlyKingPst = []int32{
	-65, 23, 16, -15, -56, -34, 2, 13,
	29, -1, -20, -7, -8, -4, -38, -29,
	-9, 24, 2, -16, -20, 6, 22, -22,
	-17, -20, -12, -27, -30, -25, -14, -36,
	-49, -1, -27, -39, -46, -44, -33, -51,
	-14, -14, -22, -46, -44, -30, -15, -27,
	1, 7, -8, -64, -43, -16, 9, 8,
	-15, 36, 12, -54, 8, -28, 24, 14,
}

// Endgame

var latePawnPst = []int32{
	0, 0, 0, 0, 0, 0, 0, 0,
	178, 173, 158, 134, 147, 132, 165, 187,
	94, 100, 85, 67, 56, 53, 82, 84,
	32, 24, 13, 5, -2, 4, 17, 17,
	13, 9, -3, -7, -7, -8, 3, -1,
	4, 7, -6, 1, 0, -5, -1, -8,
	13, 8, 8, 10, 13, 0, 2, -7,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var lateKnightPst = []int32{
	-58, -38, -13, -28, -31, -27, -63, -99,
	-25, -8, -25, -2, -9, -25, -24, -52,
	-24, -20, 10, 9, -1, -9, -19, -41,
	-17, 3, 22, 22, 22, 11, 8, -18,
	-18, -6, 16, 25, 16, 17, 4, -18,
	-23, -3, -1, 15, 10, -3, -20, -22,
	-42, -20, -10, -5, -2, -20, -23, -44,
	-29, -51, -23, -15, -22, -18, -50, -64,
}

var lateBishopPst = []int32{
	-14, -21, -11, -8, -7, -9, -17, -24,
	-8, -4, 7, -12, -3, -13, -4, -14,
	2, -8, 0, -1, -2, 6, 0, 4,
	-3, 9, 12, 9, 14, 10, 3, 2,
	-6, 3, 13, 19, 7, 10, -3, -9,
	-12, -3, 8, 10, 13, 3, -7, -15,
	-14, -18, -7, -1, 4, -9, -15, -27,
	-23, -9, -23, -5, -9, -16, -5, -17,
}

var lateRookPst = []int32{
	13, 10, 18, 15, 12, 12, 8, 5,
	11, 13, 13, 11, -3, 3, 8, 3,
	7, 7, 7, 5, 4, -3, -5, -3,
	4, 3, 13, 1, 2, 1, -1, 2,
	3, 5, 8, 4, -5, -6, -8, -11,
	-4, 0, -5, -1, -7, -12, -8, -16,
	-6, -6, 0, 2, -9, -9, -11, -3,
	-9, 2, 3, -1, -5, -13, 4, -20,
}

var lateQueenPst = []int32{
	-9, 22, 22, 27, 27, 19, 10, 20,
	-17, 20, 32, 41, 58, 25, 30, 0,
	-20, 6, 9, 49, 47, 35, 19, 9,
	3, 22, 24, 45, 57, 40, 57, 36,
	-18, 28, 19, 47, 31, 34, 39, 23,
	-16, -27, 15, 6, 9, 17, 10, 5,
	-22, -23, -30, -16, -16, -23, -36, -32,
	-33, -28, -22, -43, -5, -32, -20, -41,
}

var lateKingPst = []int32{
	-74, -35, -18, -18, -11, 15, 4, -17,
	-12, 17, 14, 17, 17, 38, 23, 11,
	10, 17, 23, 15, 20, 45, 44, 13,
	-8, 22, 24, 27, 26, 33, 26, 3,
	-18, -4, 21, 24, 27, 23, 9, -11,
	-19, -3, 11, 21, 23, 16, 7, -9,
	-27, -11, 4, 13, 14, 4, -5, -17,
	-53, -34, -21, -11, -28, -14, -24, -43,
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

func Evaluate(position *Position) int32 {
	board := position.Board
	p := BlackPawn
	n := BlackKnight
	b := BlackBishop
	r := BlackRook
	q := BlackQueen
	turn := position.Turn()

	isEndgame := board.IsEndGame()

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
		if isEndgame {
			blackCentipawns += latePawnPst[flip[index]]
		} else {
			blackCentipawns += earlyPawnPst[flip[index]]
		}
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
		if isEndgame {
			whiteCentipawns += latePawnPst[index]
		} else {
			whiteCentipawns += earlyPawnPst[index]
		}
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
		if isEndgame {
			blackCentipawns += lateKnightPst[flip[index]]
		} else {
			blackCentipawns += earlyKnightPst[flip[index]]
		}
		pieceIter ^= mask
	}

	pieceIter = bbBlackBishop
	for pieceIter != 0 {
		blackBishopsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		if isEndgame {
			blackCentipawns += lateBishopPst[flip[index]]
		} else {
			blackCentipawns += earlyBishopPst[flip[index]]
		}
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
		if isEndgame {
			blackCentipawns += lateRookPst[flip[index]]
		} else {
			blackCentipawns += earlyRookPst[flip[index]]
		}
		pieceIter ^= mask
	}

	pieceIter = bbBlackQueen
	for pieceIter != 0 {
		blackQueensCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		if isEndgame {
			blackCentipawns += lateQueenPst[flip[index]]
		} else {
			blackCentipawns += earlyQueenPst[flip[index]]
		}
		pieceIter ^= mask
	}

	pieceIter = bbBlackKing
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		if isEndgame {
			blackCentipawns += lateKingPst[flip[index]]
		} else {
			award := earlyKingPst[flip[index]]
			if award <= 0 {
				if !position.HasTag(BlackCanCastleKingSide) {
					award -= 10
				} else if !position.HasTag(BlackCanCastleQueenSide) {
					award -= 10
				}
			}
			blackCentipawns += award
		}
		pieceIter ^= mask
	}

	// PST for other white pieces
	pieceIter = bbWhiteKnight
	for pieceIter != 0 {
		whiteKnightsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		if isEndgame {
			whiteCentipawns += lateKnightPst[index]
		} else {
			whiteCentipawns += earlyKnightPst[index]
		}
		pieceIter ^= mask
	}

	pieceIter = bbWhiteBishop
	for pieceIter != 0 {
		whiteBishopsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		if isEndgame {
			whiteCentipawns += lateBishopPst[index]
		} else {
			whiteCentipawns += earlyBishopPst[index]
		}
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
		if isEndgame {
			whiteCentipawns += lateRookPst[index]
		} else {
			whiteCentipawns += earlyRookPst[index]
		}
		pieceIter ^= mask
	}

	pieceIter = bbWhiteQueen
	for pieceIter != 0 {
		whiteQueensCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		if isEndgame {
			whiteCentipawns += lateQueenPst[index]
		} else {
			whiteCentipawns += earlyQueenPst[index]
		}
		pieceIter ^= mask
	}

	pieceIter = bbWhiteKing
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		if isEndgame {
			whiteCentipawns += lateKingPst[index]
		} else {
			award := earlyKingPst[index]
			if award < 0 {
				if !position.HasTag(WhiteCanCastleKingSide) {
					award -= 10
				} else if !position.HasTag(WhiteCanCastleQueenSide) {
					award -= 10
				}
			}
			whiteCentipawns += award
		}
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
