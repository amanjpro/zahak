package evaluation

import (
	"math/bits"

	. "github.com/amanjpro/zahak/cache"
	. "github.com/amanjpro/zahak/engine"
)

const CHECKMATE_EVAL int32 = 400_000

// Piece Square Tables

// Middle game
var earlyPawnPst = [64]int32{
	0, 0, 0, 0, 0, 0, 0, 0,
	80, 80, 80, 80, 80, 80, 80, 80,
	50, 50, 50, 50, 50, 50, 50, 50,
	30, 30, 30, 30, 30, 30, 30, 30,
	-10, -10, 0, 20, 20, 0, -10, -10,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 0, -5, -5, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var earlyKnightPst = [64]int32{
	-40, -25, -25, -25, -25, -25, -25, -40,
	-30, 0, 0, 0, 0, 0, 0, -30,
	-30, 0, 0, 0, 0, 0, 0, -30,
	-30, 0, 0, 15, 15, 0, 0, -30,
	-30, 0, 0, 15, 15, 0, 0, -30,
	-30, 0, 10, 0, 0, 10, 0, -30,
	-30, 0, 0, 5, 5, 0, 0, -30,
	-40, -30, -25, -25, -25, -25, -30, -40,
}

var earlyBishopPst = [64]int32{
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-10, 0, 5, 0, 0, 5, 0, -10,
	-10, 0, 0, 10, 10, 0, 0, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-10, -20, -20, -20, -20, -20, -20, -10,
}

var earlyRookPst = [64]int32{
	0, 0, 0, 0, 0, 0, 0, 0,
	10, 10, 10, 10, 10, 10, 10, 10,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 5, 5, 0, 0, 0,
}

var earlyQueenPst = [64]int32{
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	5, 5, 10, 10, 10, 10, 5, 5,
	5, 5, 10, 15, 15, 10, 5, 5,
}

var earlyKingPst = [64]int32{
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	20, 25, 25, -25, -25, 20, 25, 20,
}

// Endgame

var latePawnPst = [64]int32{
	0, 0, 0, 0, 0, 0, 0, 0,
	200, 200, 200, 200, 200, 200, 200, 200,
	150, 150, 150, 150, 150, 150, 150, 150,
	50, 50, 50, 50, 50, 50, 50, 50,
	10, 10, 10, 10, 10, 10, 10, 10,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var lateKnightPst = [64]int32{
	-30, -20, -10, -20, -20, -20, -30, -30,
	-10, -10, -10, -5, -5, -10, -10, -10,
	-10, -10, 10, 10, -10, -10, -10, -10,
	-10, 5, 10, 10, 10, 10, 10, -10,
	-10, -5, 10, 15, 10, 15, 5, -10,
	-10, -5, 0, 10, 10, 0, -10, -10,
	-25, -20, -10, -5, -5, -20, -20, -25,
	-30, -30, -30, -10, -10, -30, -30, -30,
}

var lateBishopPst = [64]int32{
	-10, -10, -10, -10, -10, -10, -10, -10,
	-10, -5, 5, -10, -5, -10, -5, -10,
	5, -10, 0, 0, 0, 5, 0, 5,
	-5, 10, 10, 10, 10, 10, 5, 0,
	-5, 5, 10, 15, 5, 10, -5, -10,
	-10, -5, 10, 10, 15, 5, -5, -10,
	-10, -15, -5, 0, 5, -5, -10, -15,
	-15, -10, -15, -5, -10, -10, -5, -15,
}

var lateRookPst = [64]int32{
	15, 10, 15, 15, 15, 15, 10, 5,
	10, 10, 10, 10, 5, 5, 10, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
	-5, -5, -5, -5, -5, -5, -5, -5,
	-10, -10, -10, -10, -10, -10, -10, -10,
}

var lateQueenPst = [64]int32{
	-10, 20, 20, 25, 25, 20, 10, 20,
	-15, 20, 30, 40, 40, 20, 20, 0,
	-20, 5, 10, 30, 30, 30, 5, -20,
	5, 20, 20, 30, 30, 20, 20, 5,
	-15, 25, 20, 30, 30, 20, 25, -15,
	-15, -25, 10, 5, 10, 15, 10, 5,
	-20, -20, -30, -15, -15, -20, -20, -20,
	-30, -30, -20, -30, -5, -20, -20, -20,
}

var lateKingPst = [64]int32{
	-50, -50, -50, -50, -50, -50, -50, -50,
	-15, 15, 15, 15, 15, 15, 15, -15,
	10, 15, 20, 15, 20, 20, 15, 10,
	-10, 20, 20, 20, 20, 20, 20, -10,
	-15, -5, 20, 20, 20, 20, -5, -15,
	-15, -5, 10, 20, 20, 10, -5, -15,
	-20, -10, 5, 15, 15, 5, -10, -20,
	-40, -40, -20, -10, -10, -20, -40, -40,
}

var flip = [64]int32{
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
	hash := position.Hash()
	value, ok := EvalTable.Get(hash)
	if ok {
		return value
	}

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

	blackKingSafetyCentiPawns := int32(0)
	whiteKingSafetyCentiPawns := int32(0)
	var whiteKingSquare Square
	var blackKingSquare Square

	// PST for black pawns
	blackPawnsPerFile := [8]int8{0, 0, 0, 0, 0, 0, 0, 0}
	blackLeastAdvancedPawnsPerFile := [8]Rank{Rank1, Rank1, Rank1, Rank1, Rank1, Rank1, Rank1, Rank1}
	blackMostAdvancedPawnsPerFile := [8]Rank{Rank8, Rank8, Rank8, Rank8, Rank8, Rank8, Rank8, Rank8}
	pieceIter := bbBlackPawn
	for pieceIter > 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		blackPawnsCount++
		// backwards pawn
		if board.IsBackwardPawn(mask, bbBlackPawn, Black) {
			blackCentipawns -= 15
		}
		// pawn map
		sq := Square(index)
		file := sq.File()
		rank := sq.Rank()
		blackPawnsPerFile[file] += 1
		if rank > blackLeastAdvancedPawnsPerFile[file] {
			blackLeastAdvancedPawnsPerFile[file] = rank
		}
		if rank < blackMostAdvancedPawnsPerFile[file] {
			blackMostAdvancedPawnsPerFile[file] = rank
		}
		if isEndgame {
			blackCentipawns += latePawnPst[index]
		} else {
			blackCentipawns += earlyPawnPst[index]
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
		if board.IsBackwardPawn(mask, bbWhitePawn, White) {
			whiteCentipawns -= 15
		}
		// pawn map
		sq := Square(index)
		file := sq.File()
		rank := sq.Rank()
		whitePawnsPerFile[file] += 1
		if rank < whiteLeastAdvancedPawnsPerFile[file] {
			whiteLeastAdvancedPawnsPerFile[file] = rank
		}
		if rank > whiteMostAdvancedPawnsPerFile[file] {
			whiteMostAdvancedPawnsPerFile[file] = rank
		}
		if isEndgame {
			whiteCentipawns += latePawnPst[flip[index]]
		} else {
			whiteCentipawns += earlyPawnPst[flip[index]]
		}
		pieceIter ^= mask
	}

	// black king
	pieceIter = bbBlackKing
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		if isEndgame {
			blackCentipawns += lateKingPst[index]
		} else {
			award := earlyKingPst[index]
			if award <= 0 {
				if !position.HasTag(BlackCanCastleKingSide) {
					award -= 10
				} else if !position.HasTag(BlackCanCastleQueenSide) {
					award -= 10
				}
			}
			blackCentipawns += award
		}
		blackKingSquare = Square(index)
		pieceIter ^= mask
	}

	// white king
	pieceIter = bbWhiteKing
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		if isEndgame {
			whiteCentipawns += lateKingPst[flip[index]]
		} else {
			award := earlyKingPst[flip[index]]
			if award <= 0 {
				if !position.HasTag(WhiteCanCastleKingSide) {
					award -= 10
				} else if !position.HasTag(WhiteCanCastleQueenSide) {
					award -= 10
				}
			}
			whiteCentipawns += award
		}

		whiteKingSquare = Square(index)
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
			} else if i != 7 && i != 0 && whitePawnsPerFile[i-1] <= 0 && whitePawnsPerFile[i+1] <= 0 {
				isIsolated = true
			}
			if isIsolated {
				whiteCentipawns -= 15
			}
		}

		// isolated pawn penalty - black
		if blackPawnsPerFile[i] > 0 {
			isIsolated := false
			if i == 0 && blackPawnsPerFile[i+1] <= 0 {
				isIsolated = true
			} else if i == 7 && blackPawnsPerFile[i-1] <= 0 {
				isIsolated = true
			} else if i != 0 && i != 7 && blackPawnsPerFile[i-1] <= 0 && blackPawnsPerFile[i+1] <= 0 {
				isIsolated = true
			}
			if isIsolated {
				blackCentipawns -= 15
			}
		}

		// double pawn penalty - black
		if blackPawnsPerFile[i] > 1 {
			blackCentipawns -= 15
		}
		// double pawn penalty - white
		if whitePawnsPerFile[i] > 1 {
			whiteCentipawns -= 15
		}
		// passed and candidate passed pawn award
		rank := whiteMostAdvancedPawnsPerFile[i]
		if rank != Rank1 {
			if blackLeastAdvancedPawnsPerFile[i] == Rank8 || blackLeastAdvancedPawnsPerFile[i] < rank { // candidate
				if i == 0 {
					if blackLeastAdvancedPawnsPerFile[i+1] == Rank8 || blackLeastAdvancedPawnsPerFile[i+1] < rank { // passed pawn
						if isEndgame {
							whiteCentipawns += 50 //passed pawn
						} else {
							whiteCentipawns += 20 //passed pawn
						}
					} else {
						if isEndgame {
							whiteCentipawns += 25 // candidate passed pawn
						} else {
							whiteCentipawns += 10
						}
					}
				} else if i == 7 {
					if blackLeastAdvancedPawnsPerFile[i-1] == Rank8 || blackLeastAdvancedPawnsPerFile[i-1] < rank { // passed pawn
						if isEndgame {
							whiteCentipawns += 50 //passed pawn
						} else {
							whiteCentipawns += 20
						}
					} else {
						whiteCentipawns += 25 // candidate passed pawn
					}
				} else {
					if (blackLeastAdvancedPawnsPerFile[i-1] == Rank8 || blackLeastAdvancedPawnsPerFile[i-1] < rank) &&
						(blackLeastAdvancedPawnsPerFile[i+1] == Rank8 || blackLeastAdvancedPawnsPerFile[i+1] < rank) { // passed pawn
						if isEndgame {
							whiteCentipawns += 50 //passed pawn
						} else {
							whiteCentipawns += 20 //passed pawn
						}
					} else {
						if isEndgame {
							whiteCentipawns += 25 // candidate passed pawn
						} else {
							whiteCentipawns += 10 // candidate passed pawn
						}
					}
				}
			}
		}

		rank = blackMostAdvancedPawnsPerFile[i]
		if rank != Rank8 {
			if whiteLeastAdvancedPawnsPerFile[i] == Rank1 || whiteLeastAdvancedPawnsPerFile[i] > rank { // candidate
				if i == 0 {
					if whiteLeastAdvancedPawnsPerFile[i+1] == Rank1 || whiteLeastAdvancedPawnsPerFile[i+1] > rank { // passed pawn
						if isEndgame {
							blackCentipawns += 50 //passed pawn
						} else {
							blackCentipawns += 20 //passed pawn
						}
					} else {
						if isEndgame {
							blackCentipawns += 25 // candidate passed pawn
						} else {
							blackCentipawns += 10 // candidate passed pawn
						}
					}
				} else if i == 7 {
					if whiteLeastAdvancedPawnsPerFile[i-1] == Rank1 || whiteLeastAdvancedPawnsPerFile[i-1] > rank { // passed pawn
						if isEndgame {
							blackCentipawns += 50 //passed pawn
						} else {
							blackCentipawns += 20 //passed pawn
						}
					} else {
						if isEndgame {
							blackCentipawns += 25 // candidate passed pawn
						} else {
							blackCentipawns += 10 // candidate passed pawn
						}
					}
				} else {
					if (whiteLeastAdvancedPawnsPerFile[i-1] == Rank1 || whiteLeastAdvancedPawnsPerFile[i-1] > rank) &&
						(whiteLeastAdvancedPawnsPerFile[i+1] == Rank1 || whiteLeastAdvancedPawnsPerFile[i+1] > rank) { // passed pawn
						if isEndgame {
							blackCentipawns += 50 //passed pawn
						} else {
							blackCentipawns += 20 //passed pawn
						}
					} else {
						if isEndgame {
							blackCentipawns += 25 // candidate passed pawn
						} else {
							blackCentipawns += 10 // candidate passed pawn
						}
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
			blackCentipawns += lateKnightPst[index]
		} else {
			blackCentipawns += earlyKnightPst[index]
		}
		pieceIter ^= mask
	}

	pieceIter = bbBlackBishop
	for pieceIter != 0 {
		blackBishopsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		if isEndgame {
			blackCentipawns += lateBishopPst[index]
		} else {
			blackCentipawns += earlyBishopPst[index]
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
				blackCentipawns += 25
			} else { // semi-open file
				blackCentipawns += 15
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
			blackCentipawns += lateRookPst[index]
		} else {
			blackCentipawns += earlyRookPst[index]
		}
		pieceIter ^= mask
	}

	pieceIter = bbBlackQueen
	for pieceIter != 0 {
		blackQueensCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		if isEndgame {
			blackCentipawns += lateQueenPst[index]
		} else {
			blackCentipawns += earlyQueenPst[index]
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
			whiteCentipawns += lateKnightPst[flip[index]]
		} else {
			whiteCentipawns += earlyKnightPst[flip[index]]
		}
		pieceIter ^= mask
	}

	pieceIter = bbWhiteBishop
	for pieceIter != 0 {
		whiteBishopsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		if isEndgame {
			whiteCentipawns += lateBishopPst[flip[index]]
		} else {
			whiteCentipawns += earlyBishopPst[flip[index]]
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
				whiteCentipawns += 25
			} else { // semi-open file
				whiteCentipawns += 15
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
			whiteCentipawns += lateRookPst[flip[index]]
		} else {
			whiteCentipawns += earlyRookPst[flip[index]]
		}
		pieceIter ^= mask
	}

	pieceIter = bbWhiteQueen
	for pieceIter != 0 {
		whiteQueensCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		if isEndgame {
			whiteCentipawns += lateQueenPst[flip[index]]
		} else {
			whiteCentipawns += earlyQueenPst[flip[index]]
		}
		pieceIter ^= mask
	}

	{
		// Black's Middle-game king safety
		square := blackKingSquare
		file := square.File()
		rank := square.Rank()

		blackKingSafetyCentiPawns -= (int32(Rank8 - rank)) * 5 //ranks start from 0

		var files = [3]int32{-1, -1, -1}
		if file == FileH {
			files[0] = int32(FileH)
			files[1] = int32(FileG)
			files[2] = int32(FileF)
		} else if file == FileA {
			files[0] = int32(FileA)
			files[1] = int32(FileB)
			files[2] = int32(FileC)
		} else {
			files[0] = int32(file) - 1
			files[1] = int32(file)
			files[2] = int32(file) + 1
		}

		for f := range files {
			if f == int(FileE) || f == int(FileD) { // Let's encourage e5 and d5
				continue
			}
			if blackPawnsPerFile[f] == 0 { // no pawn here
				if whitePawnsPerFile[f] == 0 { // open file!!
					blackKingSafetyCentiPawns -= 60
				} else {
					blackKingSafetyCentiPawns -= 50
				}
			} else {
				if blackLeastAdvancedPawnsPerFile[f] == Rank5 {
					blackKingSafetyCentiPawns -= 25
				} else if blackLeastAdvancedPawnsPerFile[f] <= Rank4 {
					blackKingSafetyCentiPawns -= 35 + 8 - int32(blackLeastAdvancedPawnsPerFile[f])
				}
			}

			if whitePawnsPerFile[f] != 0 {
				if whiteMostAdvancedPawnsPerFile[f] >= Rank5 {
					blackKingSafetyCentiPawns -= 25
				}
			} else {
				wfile := int8(whiteKingSquare.File())
				if File(wfile-1) != file &&
					File(wfile) != file &&
					File(wfile)+1 != file {
					blackKingSafetyCentiPawns -= 40 // white can pile up
				}
			}
		}
	}

	{
		// White's Middle-game king safety
		square := whiteKingSquare
		file := square.File()
		rank := square.Rank()

		whiteKingSafetyCentiPawns -= int32(rank) * 5 //ranks start from 0

		var files = [3]int32{-1, -1, -1}
		if file == FileH {
			files[0] = int32(FileH)
			files[1] = int32(FileG)
			files[2] = int32(FileF)
		} else if file == FileA {
			files[0] = int32(FileA)
			files[1] = int32(FileB)
			files[2] = int32(FileC)
		} else {
			files[0] = int32(file) - 1
			files[1] = int32(file)
			files[2] = int32(file) + 1
		}

		for f := range files {
			if f == int(FileE) || f == int(FileD) { // Let's encourage e4 and d4
				continue
			}
			if whitePawnsPerFile[f] == 0 { // no pawn here
				if blackPawnsPerFile[f] == 0 { // open file!!
					whiteKingSafetyCentiPawns -= 60
				} else {
					whiteKingSafetyCentiPawns -= 50
				}
			} else {
				if whiteLeastAdvancedPawnsPerFile[f] == Rank4 {
					whiteKingSafetyCentiPawns -= 25
				} else if whiteLeastAdvancedPawnsPerFile[f] >= Rank4 {
					whiteKingSafetyCentiPawns -= 35 + int32(whiteLeastAdvancedPawnsPerFile[f])
				}
			}

			if blackPawnsPerFile[f] != 0 {
				if blackMostAdvancedPawnsPerFile[f] <= Rank4 {
					whiteKingSafetyCentiPawns -= 25
				}
			} else {
				bfile := int8(blackKingSquare.File())
				if File(bfile-1) != file &&
					File(bfile) != file &&
					File(bfile)+1 != file {
					whiteKingSafetyCentiPawns -= 40 // black can pile up
				}
			}
		}
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

	aggressivityFactor := int32(1)
	if !isEndgame {
		aggressivityFactor = 2
	}
	whiteCentipawns += aggressivityFactor * int32(wAttackCounts-bAttackCounts)
	blackCentipawns += aggressivityFactor * int32(bAttackCounts-wAttackCounts)

	whiteCentipawns += aggressivityFactor * int32(2*(whiteAggressivity-blackAggressivity))
	blackCentipawns += aggressivityFactor * int32(2*(blackAggressivity-whiteAggressivity))

	// king safety
	whiteFactor := (blackKnightsCount*n.Weight() + blackBishopsCount*b.Weight() +
		blackRooksCount*r.Weight() + blackQueensCount*q.Weight()) * 2 / p.Weight()

	blackFactor := (whiteKnightsCount*n.Weight() + whiteBishopsCount*b.Weight() +
		whiteRooksCount*r.Weight() + whiteQueensCount*q.Weight()) * 2 / p.Weight()

	blackCentipawns += (blackKingSafetyCentiPawns - blackFactor)
	whiteCentipawns += (whiteKingSafetyCentiPawns - whiteFactor)

	var eval int32
	if turn == White {
		eval = whiteCentipawns - blackCentipawns
	} else {
		eval = blackCentipawns - whiteCentipawns
	}
	EvalTable.Set(hash, eval)
	return eval
}
