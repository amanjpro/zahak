package evaluation

import (
	"math/bits"

	. "github.com/amanjpro/zahak/engine"
)

const CHECKMATE_EVAL int16 = 30000
const MAX_NON_CHECKMATE int16 = 25000

// Piece Square Tables

// Middle game
var EarlyPawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	50, 50, 50, 50, 50, 50, 50, 50,
	10, 10, 20, 30, 30, 20, 10, 10,
	5, 5, 10, 25, 25, 10, 5, 5,
	0, 0, 0, 20, 20, 0, 0, 0,
	5, -5, -10, 0, 0, -10, -5, 5,
	5, 10, 10, -20, -20, 10, 10, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
	//
	// 0, 0, 0, 0, 0, 0, 0, 0,
	// 50, 50, 50, 50, 50, 50, 50, 50,
	// 30, 30, 30, 40, 40, 30, 30, 30,
	// 5, 5, 5, 30, 30, 5, 5, 5,
	// 0, 0, 0, 20, 20, 0, 0, 0,
	// 0, 0, 0, 10, 10, 0, 0, 0,
	// 0, 0, 0, -5, -5, 0, 0, 0,
	// 0, 0, 0, 0, 0, 0, 0, 0,
}

var EarlyKnightPst = [64]int16{
	-40, -25, -25, -25, -25, -25, -25, -40,
	-30, 0, 0, 0, 0, 0, 0, -30,
	-30, 0, 0, 0, 0, 0, 0, -30,
	-30, 0, 0, 15, 15, 0, 0, -30,
	-30, 0, 0, 15, 15, 0, 0, -30,
	-30, 0, 10, 0, 0, 10, 0, -30,
	-30, 0, 0, 5, 5, 0, 0, -30,
	-40, -30, -25, -25, -25, -25, -30, -40,
}

var EarlyBishopPst = [64]int16{
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-10, 0, 5, 0, 0, 5, 0, -10,
	-10, 0, 0, 10, 10, 0, 0, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-10, -20, -20, -20, -20, -20, -20, -10,
}

var EarlyRookPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	10, 10, 10, 10, 10, 10, 10, 10,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 5, 5, 0, 0, 0,
}

var EarlyQueenPst = [64]int16{
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	5, 5, -25, -25, -25, -25, 5, 5,
	5, 5, 10, 15, 15, 10, 5, 5,
}

var EarlyKingPst = [64]int16{
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	20, 25, 20, -15, -15, -15, 25, 20,
}

// Endgame

var LatePawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	80, 80, 80, 80, 80, 80, 80, 80,
	60, 60, 60, 60, 60, 60, 60, 60,
	40, 40, 40, 40, 40, 40, 40, 40,
	10, 10, 10, 10, 10, 10, 10, 10,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var LateKnightPst = [64]int16{
	-30, -20, -10, -20, -20, -20, -30, -30,
	-10, -10, -10, -5, -5, -10, -10, -10,
	-10, -10, 10, 10, -10, -10, -10, -10,
	-10, 5, 10, 10, 10, 10, 10, -10,
	-10, -5, 10, 15, 10, 15, 5, -10,
	-10, -5, 0, 10, 10, 0, -10, -10,
	-25, -20, -10, -5, -5, -20, -20, -25,
	-30, -30, -30, -10, -10, -30, -30, -30,
}

var LateBishopPst = [64]int16{
	-10, -10, -10, -10, -10, -10, -10, -10,
	-10, -5, 5, -10, -5, -10, -5, -10,
	5, -10, 0, 0, 0, 5, 0, 5,
	-5, 10, 10, 10, 10, 10, 5, 0,
	-5, 5, 10, 15, 5, 10, -5, -10,
	-10, -5, 10, 10, 15, 5, -5, -10,
	-10, -15, -5, 0, 5, -5, -10, -15,
	-15, -10, -15, -5, -10, -10, -5, -15,
}

var LateRookPst = [64]int16{
	15, 10, 15, 15, 15, 15, 10, 5,
	10, 10, 10, 10, 5, 5, 10, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
	-5, -5, -5, -5, -5, -5, -5, -5,
	-10, -10, -10, -10, -10, -10, -10, -10,
}

var LateQueenPst = [64]int16{
	-10, 20, 20, 25, 25, 20, 10, 20,
	-15, 20, 30, 40, 40, 20, 20, 0,
	-20, 5, 10, 30, 30, 30, 5, -20,
	5, 20, 20, 30, 30, 20, 20, 5,
	-15, 25, 20, 30, 30, 20, 25, -15,
	-15, -25, 10, 5, 10, 15, 10, 5,
	-20, -20, -30, -15, -15, -20, -20, -20,
	-30, -30, -20, -30, -5, -20, -20, -20,
}

var LateKingPst = [64]int16{
	-50, -50, -50, -50, -50, -50, -50, -50,
	-15, 15, 15, 15, 15, 15, 15, -15,
	10, 15, 20, 15, 20, 20, 15, 10,
	-10, 20, 20, 20, 20, 20, 20, -10,
	-15, -5, 20, 20, 20, 20, -5, -15,
	-15, -5, 10, 20, 20, 10, -5, -15,
	-20, -10, 5, 15, 15, 5, -10, -20,
	-40, -40, -20, -10, -10, -20, -40, -40,
}

var flip = [64]int16{
	56, 57, 58, 59, 60, 61, 62, 63,
	48, 49, 50, 51, 52, 53, 54, 55,
	40, 41, 42, 43, 44, 45, 46, 47,
	32, 33, 34, 35, 36, 37, 38, 39,
	24, 25, 26, 27, 28, 29, 30, 31,
	16, 17, 18, 19, 20, 21, 22, 23,
	8, 9, 10, 11, 12, 13, 14, 15,
	0, 1, 2, 3, 4, 5, 6, 7,
}

var BackwardPawnAward = int16(25)
var IsolatedPawnAward = int16(25)
var DoublePawnAward = int16(25)
var EndgamePassedPawnAward = int16(50)
var MiddlegamePassedPawnAward = int16(20)
var EndgameCandidatePassedPawnAward = int16(25)
var MiddlegameCandidatePassedPawnAward = int16(10)
var RookOpenFileAward = int16(25)
var RookSemiOpenFileAward = int16(15)
var VeritcalDoubleRookAward = int16(50)
var HorizontalDoubleRookAward = int16(30)
var PawnFactorCoeff = int16(2)
var AggressivityFactorCoeff = int16(1)
var MiddlegameAggressivityFactorCoeff = int16(2)
var CastlingAward = int16(10)

func PSQT(piece Piece, sq Square, isEndgame bool) int16 {
	switch piece {
	case WhitePawn:
		return EarlyPawnPst[flip[int(sq)]]
	case WhiteKnight:
		return EarlyKnightPst[flip[int(sq)]]
	case WhiteBishop:
		return EarlyBishopPst[flip[int(sq)]]
	case WhiteRook:
		return EarlyRookPst[flip[int(sq)]]
	case WhiteQueen:
		return EarlyQueenPst[flip[int(sq)]]
	case WhiteKing:
		if isEndgame {
			return LateKingPst[flip[int(sq)]]
		}
		return EarlyKingPst[flip[int(sq)]]
	case BlackPawn:
		return EarlyPawnPst[int(sq)]
	case BlackKnight:
		return EarlyKnightPst[int(sq)]
	case BlackBishop:
		return EarlyBishopPst[int(sq)]
	case BlackRook:
		return EarlyRookPst[int(sq)]
	case BlackQueen:
		return EarlyQueenPst[int(sq)]
	case BlackKing:
		if isEndgame {
			return LateKingPst[int(sq)]
		}
		return EarlyKingPst[int(sq)]
	}
	return 0
}

func Evaluate(position *Position) int16 {
	board := position.Board
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

	blackPawnsCount := int16(0)
	blackKnightsCount := int16(0)
	blackBishopsCount := int16(0)
	blackRooksCount := int16(0)
	blackQueensCount := int16(0)

	whitePawnsCount := int16(0)
	whiteKnightsCount := int16(0)
	whiteBishopsCount := int16(0)
	whiteRooksCount := int16(0)
	whiteQueensCount := int16(0)

	blackCentipawns := int16(0)
	whiteCentipawns := int16(0)
	whites := board.GetWhitePieces()
	blacks := board.GetBlackPieces()
	all := whites | blacks

	// PST for black pawns
	blackPawnsPerFile := [8]int8{0, 0, 0, 0, 0, 0, 0, 0}
	blackLeastAdvancedPawnsPerFile := [8]Rank{Rank1, Rank1, Rank1, Rank1, Rank1, Rank1, Rank1, Rank1}
	blackMostAdvancedPawnsPerFile := [8]Rank{Rank8, Rank8, Rank8, Rank8, Rank8, Rank8, Rank8, Rank8}
	pieceIter := bbBlackPawn
	for pieceIter > 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		blackPawnsCount++
		// backwards pawn
		if board.IsBackwardPawn(mask, bbBlackPawn, Black) {
			blackCentipawns -= BackwardPawnAward
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
			blackCentipawns += LatePawnPst[index]
		} else {
			blackCentipawns += EarlyPawnPst[index]
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
		mask := SquareMask(uint64(index))
		// backwards pawn
		if board.IsBackwardPawn(mask, bbWhitePawn, White) {
			whiteCentipawns -= BackwardPawnAward
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
			whiteCentipawns += LatePawnPst[flip[index]]
		} else {
			whiteCentipawns += EarlyPawnPst[flip[index]]
		}
		pieceIter ^= mask
	}

	for i := 0; i < 8; i++ {
		// isoLated pawn penalty - white
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
				whiteCentipawns -= IsolatedPawnAward
			}
		}

		// isoLated pawn penalty - black
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
				blackCentipawns -= IsolatedPawnAward
			}
		}

		// double pawn penalty - black
		if blackPawnsPerFile[i] > 1 {
			blackCentipawns -= DoublePawnAward
		}
		// double pawn penalty - white
		if whitePawnsPerFile[i] > 1 {
			whiteCentipawns -= DoublePawnAward
		}
		// passed and candidate passed pawn award
		rank := whiteMostAdvancedPawnsPerFile[i]
		if rank != Rank1 {
			if blackLeastAdvancedPawnsPerFile[i] == Rank8 || blackLeastAdvancedPawnsPerFile[i] < rank { // candidate
				if i == 0 {
					if blackLeastAdvancedPawnsPerFile[i+1] == Rank8 || blackLeastAdvancedPawnsPerFile[i+1] < rank { // passed pawn
						if isEndgame {
							whiteCentipawns += EndgamePassedPawnAward //passed pawn
						} else {
							whiteCentipawns += MiddlegamePassedPawnAward //passed pawn
						}
					} else {
						if isEndgame {
							whiteCentipawns += EndgameCandidatePassedPawnAward // candidate passed pawn
						} else {
							whiteCentipawns += MiddlegameCandidatePassedPawnAward
						}
					}
				} else if i == 7 {
					if blackLeastAdvancedPawnsPerFile[i-1] == Rank8 || blackLeastAdvancedPawnsPerFile[i-1] < rank { // passed pawn
						if isEndgame {
							whiteCentipawns += EndgamePassedPawnAward //passed pawn
						} else {
							whiteCentipawns += MiddlegamePassedPawnAward
						}
					} else {
						if isEndgame {
							whiteCentipawns += EndgameCandidatePassedPawnAward // candidate passed pawn
						} else {
							whiteCentipawns += MiddlegameCandidatePassedPawnAward
						}
					}
				} else {
					if (blackLeastAdvancedPawnsPerFile[i-1] == Rank8 || blackLeastAdvancedPawnsPerFile[i-1] < rank) &&
						(blackLeastAdvancedPawnsPerFile[i+1] == Rank8 || blackLeastAdvancedPawnsPerFile[i+1] < rank) { // passed pawn
						if isEndgame {
							whiteCentipawns += EndgamePassedPawnAward //passed pawn
						} else {
							whiteCentipawns += MiddlegamePassedPawnAward //passed pawn
						}
					} else {
						if isEndgame {
							whiteCentipawns += EndgameCandidatePassedPawnAward // candidate passed pawn
						} else {
							whiteCentipawns += MiddlegameCandidatePassedPawnAward // candidate passed pawn
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
							blackCentipawns += EndgamePassedPawnAward //passed pawn
						} else {
							blackCentipawns += MiddlegamePassedPawnAward //passed pawn
						}
					} else {
						if isEndgame {
							blackCentipawns += EndgameCandidatePassedPawnAward // candidate passed pawn
						} else {
							blackCentipawns += MiddlegameCandidatePassedPawnAward // candidate passed pawn
						}
					}
				} else if i == 7 {
					if whiteLeastAdvancedPawnsPerFile[i-1] == Rank1 || whiteLeastAdvancedPawnsPerFile[i-1] > rank { // passed pawn
						if isEndgame {
							blackCentipawns += EndgamePassedPawnAward //passed pawn
						} else {
							blackCentipawns += MiddlegamePassedPawnAward //passed pawn
						}
					} else {
						if isEndgame {
							blackCentipawns += EndgameCandidatePassedPawnAward // candidate passed pawn
						} else {
							blackCentipawns += MiddlegameCandidatePassedPawnAward // candidate passed pawn
						}
					}
				} else {
					if (whiteLeastAdvancedPawnsPerFile[i-1] == Rank1 || whiteLeastAdvancedPawnsPerFile[i-1] > rank) &&
						(whiteLeastAdvancedPawnsPerFile[i+1] == Rank1 || whiteLeastAdvancedPawnsPerFile[i+1] > rank) { // passed pawn
						if isEndgame {
							blackCentipawns += EndgamePassedPawnAward //passed pawn
						} else {
							blackCentipawns += MiddlegamePassedPawnAward //passed pawn
						}
					} else {
						if isEndgame {
							blackCentipawns += EndgameCandidatePassedPawnAward // candidate passed pawn
						} else {
							blackCentipawns += MiddlegameCandidatePassedPawnAward // candidate passed pawn
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
		mask := SquareMask(uint64(index))
		if isEndgame {
			blackCentipawns += LateKnightPst[index]
		} else {
			blackCentipawns += EarlyKnightPst[index]
		}
		pieceIter ^= mask
	}

	pieceIter = bbBlackBishop
	for pieceIter != 0 {
		blackBishopsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		if isEndgame {
			blackCentipawns += LateBishopPst[index]
		} else {
			blackCentipawns += EarlyBishopPst[index]
		}
		pieceIter ^= mask
	}

	pieceIter = bbBlackRook
	for pieceIter != 0 {
		blackRooksCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		file := Square(index).File()
		if blackPawnsPerFile[file] == 0 {
			if whitePawnsPerFile[file] == 0 { // open file
				blackCentipawns += RookOpenFileAward
			} else { // semi-open file
				blackCentipawns += RookSemiOpenFileAward
			}
		}
		sq := Square(index)
		if blackRooksCount == 1 {
			if board.IsVerticalDoubleRook(sq, bbBlackRook, all) {
				// double-rook vertical
				blackCentipawns += VeritcalDoubleRookAward
			} else if board.IsHorizontalDoubleRook(sq, bbBlackRook, all) {
				// double-rook horizontal
				blackCentipawns += HorizontalDoubleRookAward
			}
		}
		if isEndgame {
			blackCentipawns += LateRookPst[index]
		} else {
			blackCentipawns += EarlyRookPst[index]
		}
		pieceIter ^= mask
	}

	pieceIter = bbBlackQueen
	for pieceIter != 0 {
		blackQueensCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		if isEndgame {
			blackCentipawns += LateQueenPst[index]
		} else {
			blackCentipawns += EarlyQueenPst[index]
		}
		pieceIter ^= mask
	}

	pieceIter = bbBlackKing
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		if isEndgame {
			blackCentipawns += LateKingPst[index]
		} else {
			award := EarlyKingPst[index]
			if award <= 0 {
				if !position.HasTag(BlackCanCastleKingSide) {
					award -= CastlingAward
				} else if !position.HasTag(BlackCanCastleQueenSide) {
					award -= CastlingAward
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
		mask := SquareMask(uint64(index))
		if isEndgame {
			whiteCentipawns += LateKnightPst[flip[index]]
		} else {
			whiteCentipawns += EarlyKnightPst[flip[index]]
		}
		pieceIter ^= mask
	}

	pieceIter = bbWhiteBishop
	for pieceIter != 0 {
		whiteBishopsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		if isEndgame {
			whiteCentipawns += LateBishopPst[flip[index]]
		} else {
			whiteCentipawns += EarlyBishopPst[flip[index]]
		}
		pieceIter ^= mask
	}

	pieceIter = bbWhiteRook
	for pieceIter != 0 {
		whiteRooksCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		file := Square(index).File()
		if whitePawnsPerFile[file] == 0 {
			if blackPawnsPerFile[file] == 0 { // open file
				whiteCentipawns += RookOpenFileAward
			} else { // semi-open file
				whiteCentipawns += RookSemiOpenFileAward
			}
		}
		sq := Square(index)
		if whiteRooksCount == 1 {
			if board.IsVerticalDoubleRook(sq, bbWhiteRook, all) {
				// double-rook vertical
				whiteCentipawns += VeritcalDoubleRookAward
			} else if board.IsHorizontalDoubleRook(sq, bbWhiteRook, all) {
				// double-rook horizontal
				whiteCentipawns += HorizontalDoubleRookAward
			}
		}
		if isEndgame {
			whiteCentipawns += LateRookPst[flip[index]]
		} else {
			whiteCentipawns += EarlyRookPst[flip[index]]
		}
		pieceIter ^= mask
	}

	pieceIter = bbWhiteQueen
	for pieceIter != 0 {
		whiteQueensCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		if isEndgame {
			whiteCentipawns += LateQueenPst[flip[index]]
		} else {
			whiteCentipawns += EarlyQueenPst[flip[index]]
		}
		pieceIter ^= mask
	}

	pieceIter = bbWhiteKing
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		if isEndgame {
			whiteCentipawns += LateKingPst[flip[index]]
		} else {
			award := EarlyKingPst[flip[index]]
			if award <= 0 {
				if !position.HasTag(WhiteCanCastleKingSide) {
					award -= CastlingAward
				} else if !position.HasTag(WhiteCanCastleQueenSide) {
					award -= CastlingAward
				}
			}
			whiteCentipawns += award
		}

		pieceIter ^= mask
	}

	pawnFactor := int16(16-blackPawnsCount-whitePawnsCount) * PawnFactorCoeff

	blackCentipawns += blackPawnsCount * BlackPawn.Weight()
	blackCentipawns += blackKnightsCount * (BlackKnight.Weight() - pawnFactor)
	blackCentipawns += blackBishopsCount * (BlackBishop.Weight() + pawnFactor)
	blackCentipawns += blackRooksCount * (BlackRook.Weight() + pawnFactor)
	blackCentipawns += blackQueensCount * BlackQueen.Weight()

	whiteCentipawns += whitePawnsCount * WhitePawn.Weight()
	whiteCentipawns += whiteKnightsCount * (WhiteKnight.Weight() - pawnFactor)
	whiteCentipawns += whiteBishopsCount * (WhiteBishop.Weight() + pawnFactor)
	whiteCentipawns += whiteRooksCount * (WhiteRook.Weight() + pawnFactor)
	whiteCentipawns += whiteQueensCount * WhiteQueen.Weight()

	// mobility and attacks
	whiteAttacks := board.AllAttacks(Black) // get the squares that are taboo for black (white's reach)
	blackAttacks := board.AllAttacks(White) // get the squares that are taboo for whtie (black's reach)
	wAttackCounts := bits.OnesCount64(whiteAttacks)
	bAttackCounts := bits.OnesCount64(blackAttacks)

	whiteAggressivity := bits.OnesCount64(whiteAttacks >> 32) // keep hi-bits only (black's half)
	blackAggressivity := bits.OnesCount64(blackAttacks << 32) // keep lo-bits only (white's half)

	aggressivityFactor := AggressivityFactorCoeff
	if !isEndgame {
		aggressivityFactor = MiddlegameAggressivityFactorCoeff
	}

	whiteCentipawns += aggressivityFactor * int16(wAttackCounts-bAttackCounts)
	blackCentipawns += aggressivityFactor * int16(bAttackCounts-wAttackCounts)

	whiteCentipawns += aggressivityFactor * int16(2*(whiteAggressivity-blackAggressivity))
	blackCentipawns += aggressivityFactor * int16(2*(blackAggressivity-whiteAggressivity))

	tempo := int16(5)
	if turn == White {
		return toEval(whiteCentipawns - blackCentipawns + tempo)
	} else {
		return toEval(blackCentipawns - whiteCentipawns + tempo)
	}
}

func toEval(eval int16) int16 {
	if eval >= CHECKMATE_EVAL {
		return MAX_NON_CHECKMATE
	} else if eval <= -CHECKMATE_EVAL {
		return -MAX_NON_CHECKMATE
	}
	return eval

}
