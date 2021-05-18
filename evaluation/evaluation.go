package evaluation

import (
	"math/bits"

	. "github.com/amanjpro/zahak/engine"
)

const CHECKMATE_EVAL int16 = 30000
const MAX_NON_CHECKMATE int16 = 25000
const PawnPhase int16 = 0
const KnightPhase int16 = 1
const BishopPhase int16 = 1
const RookPhase int16 = 2
const QueenPhase int16 = 4
const TotalPhase int16 = PawnPhase*16 + KnightPhase*4 + BishopPhase*4 + RookPhase*4 + QueenPhase*2
const HalfPhase = TotalPhase / 2
const Tempo int16 = 5

// Piece Square Tables
// Middle-game
var EarlyPawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	128, 149, 105, 111, 111, 117, 50, 4,
	-10, 6, 19, 29, 63, 57, 18, -26,
	-19, 6, -5, 14, 14, 1, 9, -35,
	-35, -10, -14, 3, 8, -4, 0, -36,
	-31, -14, -14, -22, -8, -8, 25, -23,
	-43, -10, -32, -38, -29, 14, 27, -34,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var EarlyKnightPst = [64]int16{
	-210, -76, -44, -47, 39, -100, -32, -132,
	-78, -46, 76, 41, 29, 66, 11, -22,
	-52, 64, 47, 65, 91, 128, 84, 37,
	-12, 21, 24, 63, 45, 84, 25, 29,
	-12, 13, 20, 17, 34, 26, 30, -3,
	-21, -5, 16, 18, 26, 23, 31, -13,
	-25, -50, -9, 0, 2, 25, -8, -9,
	-111, -19, -59, -35, -13, -24, -16, -26,
}

var EarlyBishopPst = [64]int16{
	-45, 19, -67, -27, -21, -19, 12, -4,
	-20, 36, -5, -1, 58, 72, 49, -33,
	-6, 52, 65, 65, 51, 69, 56, 19,
	12, 23, 40, 72, 65, 64, 30, 14,
	13, 38, 34, 48, 56, 34, 33, 26,
	23, 38, 36, 37, 33, 50, 39, 28,
	29, 39, 37, 19, 28, 46, 54, 26,
	-13, 21, 6, -3, 6, 7, -15, -4,
}

var EarlyRookPst = [64]int16{
	19, 28, 11, 41, 44, 3, 10, 17,
	25, 30, 54, 50, 68, 61, 22, 37,
	-12, 10, 24, 25, 13, 40, 51, 13,
	-28, -19, 3, 19, 22, 29, -2, -19,
	-35, -34, -10, -4, 7, -12, 11, -20,
	-50, -27, -18, -26, 2, -7, -6, -35,
	-50, -16, -22, -15, -7, 9, -13, -76,
	-22, -15, 1, 12, 14, 0, -34, -22,
}

var EarlyQueenPst = [64]int16{
	-64, -24, -4, -11, 28, 25, 25, 6,
	-44, -53, -17, -9, -30, 35, 24, 28,
	-28, -29, -3, -3, 16, 37, 28, 41,
	-45, -43, -28, -26, -11, 8, -13, -13,
	-19, -41, -20, -22, -14, -14, -8, -9,
	-29, -6, -22, -13, -15, -8, 1, -7,
	-37, -14, 3, -9, -2, 6, -12, -8,
	-1, -23, -18, 2, -22, -30, -30, -46,
}

var EarlyKingPst = [64]int16{
	-35, 49, 46, 23, -29, -1, 35, 14,
	47, 9, 9, 34, 12, -10, -25, -43,
	3, 16, 14, -8, 4, 30, 26, -18,
	-20, -13, -1, -21, -23, -23, -19, -48,
	-50, -9, -37, -63, -59, -42, -44, -60,
	-9, -12, -32, -59, -58, -45, -19, -34,
	-1, 7, -22, -80, -60, -26, 5, 10,
	-13, 45, 14, -70, 3, -40, 34, 25,
}

// Endgame
var LatePawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	189, 182, 158, 141, 148, 142, 178, 213,
	104, 110, 94, 71, 56, 53, 90, 92,
	34, 25, 14, 4, -4, 4, 17, 19,
	15, 10, -4, -10, -11, -11, 2, -2,
	2, 8, -8, 2, 0, -6, -5, -11,
	15, 9, 9, 12, 12, -1, 1, -9,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var LateKnightPst = [64]int16{
	-64, -57, -22, -43, -44, -45, -79, -114,
	-39, -18, -40, -16, -25, -42, -41, -63,
	-36, -34, -4, -2, -19, -26, -40, -56,
	-28, -8, 12, 7, 9, -6, -6, -34,
	-29, -20, 5, 17, 3, 6, -10, -35,
	-38, -15, -13, 1, -5, -15, -36, -36,
	-57, -33, -26, -17, -15, -37, -39, -64,
	-44, -71, -35, -27, -37, -35, -64, -80,
}

var LateBishopPst = [64]int16{
	-15, -29, -14, -14, -9, -18, -21, -31,
	-10, -10, 5, -15, -13, -18, -13, -16,
	1, -13, -6, -9, -6, -2, -4, -1,
	-5, 6, 8, 5, 7, 2, -3, 0,
	-11, -4, 10, 13, 3, 5, -8, -15,
	-17, -8, 4, 8, 12, -2, -9, -18,
	-23, -24, -12, -3, 2, -18, -16, -37,
	-30, -18, -29, -9, -13, -19, -13, -20,
}

var LateRookPst = [64]int16{
	28, 23, 33, 27, 25, 24, 20, 20,
	21, 22, 21, 22, 4, 11, 18, 13,
	21, 22, 17, 18, 13, 5, 4, 8,
	19, 18, 24, 12, 11, 11, 6, 15,
	15, 19, 17, 13, 3, 5, -3, -3,
	8, 11, 4, 11, -2, -3, -1, -6,
	6, 2, 9, 12, 1, -3, -1, 11,
	-4, 9, 9, 5, 0, -4, 9, -22,
}

var LateQueenPst = [64]int16{
	13, 44, 48, 51, 54, 38, 22, 47,
	-7, 30, 44, 61, 82, 50, 44, 24,
	-13, 17, 18, 67, 69, 62, 47, 31,
	23, 38, 40, 59, 74, 55, 73, 50,
	-16, 43, 31, 63, 47, 48, 54, 27,
	3, -26, 24, 14, 16, 28, 25, 19,
	-32, -23, -23, -10, -7, -16, -33, -27,
	-36, -33, -21, -52, -9, -39, -27, -56,
}

var LateKingPst = [64]int16{
	-76, -41, -25, -25, -11, 11, -1, -14,
	-15, 18, 11, 10, 14, 41, 27, 17,
	12, 21, 23, 16, 18, 43, 46, 14,
	-9, 24, 25, 29, 28, 37, 31, 7,
	-19, -3, 26, 34, 36, 27, 15, -8,
	-23, -4, 14, 27, 30, 22, 9, -8,
	-31, -14, 7, 20, 21, 7, -4, -22,
	-66, -49, -29, -8, -36, -12, -38, -60,
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

var MiddlegameBackwardPawnPenalty int16 = 3
var EndgameBackwardPawnPenalty int16 = 3

var MiddlegameIsolatedPawnPenalty int16 = 9
var EndgameIsolatedPawnPenalty int16 = 4

var MiddlegameDoublePawnPenalty int16 = 5
var EndgameDoublePawnPenalty int16 = 28

var MiddlegamePassedPawnAward int16 = 0
var EndgamePassedPawnAward int16 = 10

// With tuning, those seemed useless... keeping them around for now
var MiddlegameCandidatePassedPawnAward int16 = 0
var EndgameCandidatePassedPawnAward int16 = 0

var MiddlegameRookOpenFileAward int16 = 40
var EndgameRookOpenFileAward int16 = 5

var MiddlegameRookSemiOpenFileAward int16 = 14
var EndgameRookSemiOpenFileAward int16 = 15

// Somehow tuning thinks that horizontal double rook is better than vertical
var MiddlegameVeritcalDoubleRookAward int16 = 4
var EndgameVeritcalDoubleRookAward int16 = 0

var MiddlegameHorizontalDoubleRookAward int16 = 19
var EndgameHorizontalDoubleRookAward int16 = 0

// tuning doesn't like pawn coeff
var MiddlegamePawnFactorCoeff = int16(0)
var EndgamePawnFactorCoeff = int16(0)

var MiddlegameMobilityFactorCoeff int16 = 4
var EndgameMobilityFactorCoeff int16 = 2
var MiddlegameAggressivityFactorCoeff int16 = 0
var EndgameAggressivityFactorCoeff int16 = 3
var MiddlegameInnerPawnToKingAttackCoeff int16 = 0
var EndgameInnerPawnToKingAttackCoeff int16 = 0
var MiddlegameOuterPawnToKingAttackCoeff int16 = 4
var EndgameOuterPawnToKingAttackCoeff int16 = 0
var MiddlegameInnerMinorToKingAttackCoeff int16 = 9
var EndgameInnerMinorToKingAttackCoeff int16 = 0
var MiddlegameOuterMinorToKingAttackCoeff int16 = 5
var EndgameOuterMinorToKingAttackCoeff int16 = 0
var MiddlegameInnerMajorToKingAttackCoeff int16 = 11
var EndgameInnerMajorToKingAttackCoeff int16 = 0
var MiddlegameOuterMajorToKingAttackCoeff int16 = 3
var EndgameOuterMajorToKingAttackCoeff int16 = 3

//
// var MiddlegameMobilityFactorCoeff int16 = 4
// var EndgameMobilityFactorCoeff int16 = 1
// var MiddlegameAggressivityFactorCoeff int16 = 0
// var EndgameAggressivityFactorCoeff int16 = 2
// var MiddlegameInnerPawnToKingAttackCoeff int16 = 0
// var EndgameInnerPawnToKingAttackCoeff int16 = 0
// var MiddlegameOuterPawnToKingAttackCoeff int16 = 2
// var EndgameOuterPawnToKingAttackCoeff int16 = 1
// var MiddlegameInnerMinorToKingAttackCoeff int16 = 5
// var EndgameInnerMinorToKingAttackCoeff int16 = 3
// var MiddlegameOuterMinorToKingAttackCoeff int16 = 3
// var EndgameOuterMinorToKingAttackCoeff int16 = 2
// var MiddlegameInnerMajorToKingAttackCoeff int16 = 6
// var EndgameInnerMajorToKingAttackCoeff int16 = 5
// var MiddlegameOuterMajorToKingAttackCoeff int16 = 3
// var EndgameOuterMajorToKingAttackCoeff int16 = 2

//
//
// var MiddlegameMobilityFactorCoeff int16 = 2
// var EndgameMobilityFactorCoeff int16 = 1
//
// // tuning doesn't like aggressivity
// var MiddlegameAggressivityFactorCoeff int16 = 0
// var EndgameAggressivityFactorCoeff int16 = 0
//
// var MiddlegameInnerPawnToKingAttackCoeff = int16(0)
// var EndgameInnerPawnToKingAttackCoeff = int16(0)
//
// var MiddlegameOuterPawnToKingAttackCoeff = int16(0)
// var EndgameOuterPawnToKingAttackCoeff = int16(0)
//
// var MiddlegameInnerMinorToKingAttackCoeff = int16(0)
// var EndgameInnerMinorToKingAttackCoeff = int16(0)
//
// var MiddlegameOuterMinorToKingAttackCoeff = int16(0)
// var EndgameOuterMinorToKingAttackCoeff = int16(0)
//
// var MiddlegameInnerMajorToKingAttackCoeff = int16(0)
// var EndgameInnerMajorToKingAttackCoeff = int16(0)
//
// var MiddlegameOuterMajorToKingAttackCoeff = int16(0)
// var EndgameOuterMajorToKingAttackCoeff = int16(0)
//
func PSQT(piece Piece, sq Square, isEndgame bool) int16 {
	if isEndgame {
		switch piece {
		case WhitePawn:
			return LatePawnPst[flip[int(sq)]]
		case WhiteKnight:
			return LateKnightPst[flip[int(sq)]]
		case WhiteBishop:
			return LateBishopPst[flip[int(sq)]]
		case WhiteRook:
			return LateRookPst[flip[int(sq)]]
		case WhiteQueen:
			return LateQueenPst[flip[int(sq)]]
		case WhiteKing:
			return LateKingPst[flip[int(sq)]]
		case BlackPawn:
			return LatePawnPst[int(sq)]
		case BlackKnight:
			return LateKnightPst[int(sq)]
		case BlackBishop:
			return LateBishopPst[int(sq)]
		case BlackRook:
			return LateRookPst[int(sq)]
		case BlackQueen:
			return LateQueenPst[int(sq)]
		case BlackKing:
			return LateKingPst[int(sq)]
		}
	} else {
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
			return EarlyKingPst[int(sq)]
		}
	}
	return 0
}

func Evaluate(position *Position) int16 {
	board := position.Board
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

	blackCentipawnsMG := int16(0)
	blackCentipawnsEG := int16(0)

	whiteCentipawnsMG := int16(0)
	whiteCentipawnsEG := int16(0)

	whites := board.GetWhitePieces()
	blacks := board.GetBlackPieces()
	all := whites | blacks

	var whiteKingIndex, blackKingIndex int

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
			blackCentipawnsMG -= MiddlegameBackwardPawnPenalty
			blackCentipawnsEG -= EndgameBackwardPawnPenalty
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
		blackCentipawnsEG += LatePawnPst[index]
		blackCentipawnsMG += EarlyPawnPst[index]
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
			whiteCentipawnsMG -= MiddlegameBackwardPawnPenalty
			whiteCentipawnsEG -= EndgameBackwardPawnPenalty
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
		whiteCentipawnsEG += LatePawnPst[flip[index]]
		whiteCentipawnsMG += EarlyPawnPst[flip[index]]
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
				whiteCentipawnsMG -= MiddlegameIsolatedPawnPenalty
				whiteCentipawnsEG -= EndgameIsolatedPawnPenalty
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
				blackCentipawnsMG -= MiddlegameIsolatedPawnPenalty
				blackCentipawnsEG -= EndgameIsolatedPawnPenalty
			}
		}

		// double pawn penalty - black
		if blackPawnsPerFile[i] > 1 {
			blackCentipawnsMG -= MiddlegameDoublePawnPenalty
			blackCentipawnsEG -= EndgameDoublePawnPenalty
		}
		// double pawn penalty - white
		if whitePawnsPerFile[i] > 1 {
			whiteCentipawnsMG -= MiddlegameDoublePawnPenalty
			whiteCentipawnsEG -= EndgameDoublePawnPenalty
		}
		// passed and candidate passed pawn award
		rank := whiteMostAdvancedPawnsPerFile[i]
		if rank != Rank1 {
			if blackLeastAdvancedPawnsPerFile[i] == Rank8 || blackLeastAdvancedPawnsPerFile[i] < rank { // candidate
				if i == 0 {
					if blackLeastAdvancedPawnsPerFile[i+1] == Rank8 || blackLeastAdvancedPawnsPerFile[i+1] < rank { // passed pawn
						whiteCentipawnsEG += EndgamePassedPawnAward    //passed pawn
						whiteCentipawnsMG += MiddlegamePassedPawnAward //passed pawn
					} else {
						whiteCentipawnsEG += EndgameCandidatePassedPawnAward // candidate passed pawn
						whiteCentipawnsMG += MiddlegameCandidatePassedPawnAward
					}
				} else if i == 7 {
					if blackLeastAdvancedPawnsPerFile[i-1] == Rank8 || blackLeastAdvancedPawnsPerFile[i-1] < rank { // passed pawn
						whiteCentipawnsEG += EndgamePassedPawnAward //passed pawn
						whiteCentipawnsMG += MiddlegamePassedPawnAward
					} else {
						whiteCentipawnsEG += EndgameCandidatePassedPawnAward // candidate passed pawn
						whiteCentipawnsMG += MiddlegameCandidatePassedPawnAward
					}
				} else {
					if (blackLeastAdvancedPawnsPerFile[i-1] == Rank8 || blackLeastAdvancedPawnsPerFile[i-1] < rank) &&
						(blackLeastAdvancedPawnsPerFile[i+1] == Rank8 || blackLeastAdvancedPawnsPerFile[i+1] < rank) { // passed pawn
						whiteCentipawnsEG += EndgamePassedPawnAward    //passed pawn
						whiteCentipawnsMG += MiddlegamePassedPawnAward //passed pawn
					} else {
						whiteCentipawnsEG += EndgameCandidatePassedPawnAward    // candidate passed pawn
						whiteCentipawnsMG += MiddlegameCandidatePassedPawnAward // candidate passed pawn
					}
				}
			}
		}

		rank = blackMostAdvancedPawnsPerFile[i]
		if rank != Rank8 {
			if whiteLeastAdvancedPawnsPerFile[i] == Rank1 || whiteLeastAdvancedPawnsPerFile[i] > rank { // candidate
				if i == 0 {
					if whiteLeastAdvancedPawnsPerFile[i+1] == Rank1 || whiteLeastAdvancedPawnsPerFile[i+1] > rank { // passed pawn
						blackCentipawnsEG += EndgamePassedPawnAward    //passed pawn
						blackCentipawnsMG += MiddlegamePassedPawnAward //passed pawn
					} else {
						blackCentipawnsEG += EndgameCandidatePassedPawnAward    // candidate passed pawn
						blackCentipawnsMG += MiddlegameCandidatePassedPawnAward // candidate passed pawn
					}
				} else if i == 7 {
					if whiteLeastAdvancedPawnsPerFile[i-1] == Rank1 || whiteLeastAdvancedPawnsPerFile[i-1] > rank { // passed pawn
						blackCentipawnsEG += EndgamePassedPawnAward    //passed pawn
						blackCentipawnsMG += MiddlegamePassedPawnAward //passed pawn
					} else {
						blackCentipawnsEG += EndgameCandidatePassedPawnAward    // candidate passed pawn
						blackCentipawnsMG += MiddlegameCandidatePassedPawnAward // candidate passed pawn
					}
				} else {
					if (whiteLeastAdvancedPawnsPerFile[i-1] == Rank1 || whiteLeastAdvancedPawnsPerFile[i-1] > rank) &&
						(whiteLeastAdvancedPawnsPerFile[i+1] == Rank1 || whiteLeastAdvancedPawnsPerFile[i+1] > rank) { // passed pawn
						blackCentipawnsEG += EndgamePassedPawnAward    //passed pawn
						blackCentipawnsMG += MiddlegamePassedPawnAward //passed pawn
					} else {
						blackCentipawnsEG += EndgameCandidatePassedPawnAward    // candidate passed pawn
						blackCentipawnsMG += MiddlegameCandidatePassedPawnAward // candidate passed pawn
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
		blackCentipawnsEG += LateKnightPst[index]
		blackCentipawnsMG += EarlyKnightPst[index]
		pieceIter ^= mask
	}

	pieceIter = bbBlackBishop
	for pieceIter != 0 {
		blackBishopsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		blackCentipawnsEG += LateBishopPst[index]
		blackCentipawnsMG += EarlyBishopPst[index]
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
				blackCentipawnsMG += MiddlegameRookOpenFileAward
				blackCentipawnsEG += EndgameRookOpenFileAward
			} else { // semi-open file
				blackCentipawnsMG += MiddlegameRookSemiOpenFileAward
				blackCentipawnsEG += EndgameRookSemiOpenFileAward
			}
		}
		sq := Square(index)
		if blackRooksCount == 1 {
			if board.IsVerticalDoubleRook(sq, bbBlackRook, all) {
				// double-rook vertical
				blackCentipawnsEG += EndgameVeritcalDoubleRookAward
				blackCentipawnsMG += MiddlegameVeritcalDoubleRookAward
			} else if board.IsHorizontalDoubleRook(sq, bbBlackRook, all) {
				// double-rook horizontal
				blackCentipawnsMG += MiddlegameHorizontalDoubleRookAward
				blackCentipawnsEG += EndgameHorizontalDoubleRookAward
			}
		}
		blackCentipawnsEG += LateRookPst[index]
		blackCentipawnsMG += EarlyRookPst[index]
		pieceIter ^= mask
	}

	pieceIter = bbBlackQueen
	for pieceIter != 0 {
		blackQueensCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		blackCentipawnsEG += LateQueenPst[index]
		blackCentipawnsMG += EarlyQueenPst[index]
		pieceIter ^= mask
	}

	pieceIter = bbBlackKing
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		blackCentipawnsEG += LateKingPst[index]
		award := EarlyKingPst[index]
		// if award <= 0 {
		// 	if !position.HasTag(BlackCanCastleKingSide) {
		// 		award -= MiddlegameCastlingAward
		// 	} else if !position.HasTag(BlackCanCastleQueenSide) {
		// 		award -= MiddlegameCastlingAward
		// 	}
		// }
		blackCentipawnsMG += award

		blackKingIndex = index

		pieceIter ^= mask
	}

	// PST for other white pieces
	pieceIter = bbWhiteKnight
	for pieceIter != 0 {
		whiteKnightsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		whiteCentipawnsEG += LateKnightPst[flip[index]]
		whiteCentipawnsMG += EarlyKnightPst[flip[index]]
		pieceIter ^= mask
	}

	pieceIter = bbWhiteBishop
	for pieceIter != 0 {
		whiteBishopsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		whiteCentipawnsEG += LateBishopPst[flip[index]]
		whiteCentipawnsMG += EarlyBishopPst[flip[index]]
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
				whiteCentipawnsMG += MiddlegameRookOpenFileAward
				whiteCentipawnsEG += EndgameRookOpenFileAward
			} else { // semi-open file
				whiteCentipawnsMG += MiddlegameRookSemiOpenFileAward
				whiteCentipawnsEG += EndgameRookSemiOpenFileAward
			}
		}
		sq := Square(index)
		if whiteRooksCount == 1 {
			if board.IsVerticalDoubleRook(sq, bbWhiteRook, all) {
				// double-rook vertical
				whiteCentipawnsMG += MiddlegameVeritcalDoubleRookAward
				whiteCentipawnsEG += EndgameVeritcalDoubleRookAward
			} else if board.IsHorizontalDoubleRook(sq, bbWhiteRook, all) {
				// double-rook horizontal
				whiteCentipawnsMG += MiddlegameHorizontalDoubleRookAward
				whiteCentipawnsEG += EndgameHorizontalDoubleRookAward
			}
		}
		whiteCentipawnsEG += LateRookPst[flip[index]]
		whiteCentipawnsMG += EarlyRookPst[flip[index]]
		pieceIter ^= mask
	}

	pieceIter = bbWhiteQueen
	for pieceIter != 0 {
		whiteQueensCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		whiteCentipawnsEG += LateQueenPst[flip[index]]
		whiteCentipawnsMG += EarlyQueenPst[flip[index]]
		pieceIter ^= mask
	}

	pieceIter = bbWhiteKing
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		whiteCentipawnsEG += LateKingPst[flip[index]]
		award := EarlyKingPst[flip[index]]
		// if award <= 0 {
		// 	if !position.HasTag(WhiteCanCastleKingSide) {
		// 		award -= MiddlegameCastlingAward
		// 	} else if !position.HasTag(WhiteCanCastleQueenSide) {
		// 		award -= MiddlegameCastlingAward
		// 	}
		// }
		whiteCentipawnsMG += award

		whiteKingIndex = index

		pieceIter ^= mask
	}

	pawnFactorMG := int16(16-blackPawnsCount-whitePawnsCount) * MiddlegamePawnFactorCoeff
	pawnFactorEG := int16(16-blackPawnsCount-whitePawnsCount) * EndgamePawnFactorCoeff

	blackCentipawnsMG += blackPawnsCount * BlackPawn.Weight()
	blackCentipawnsMG += blackKnightsCount * (BlackKnight.Weight() - pawnFactorMG)
	blackCentipawnsMG += blackBishopsCount * (BlackBishop.Weight() + pawnFactorMG)
	blackCentipawnsMG += blackRooksCount * (BlackRook.Weight() + pawnFactorMG)
	blackCentipawnsMG += blackQueensCount * BlackQueen.Weight()

	blackCentipawnsEG += blackPawnsCount * BlackPawn.Weight()
	blackCentipawnsEG += blackKnightsCount * (BlackKnight.Weight() - pawnFactorEG)
	blackCentipawnsEG += blackBishopsCount * (BlackBishop.Weight() + pawnFactorEG)
	blackCentipawnsEG += blackRooksCount * (BlackRook.Weight() + pawnFactorEG)
	blackCentipawnsEG += blackQueensCount * BlackQueen.Weight()

	whiteCentipawnsMG += whitePawnsCount * WhitePawn.Weight()
	whiteCentipawnsMG += whiteKnightsCount * (WhiteKnight.Weight() - pawnFactorMG)
	whiteCentipawnsMG += whiteBishopsCount * (WhiteBishop.Weight() + pawnFactorMG)
	whiteCentipawnsMG += whiteRooksCount * (WhiteRook.Weight() + pawnFactorMG)
	whiteCentipawnsMG += whiteQueensCount * WhiteQueen.Weight()

	whiteCentipawnsEG += whitePawnsCount * WhitePawn.Weight()
	whiteCentipawnsEG += whiteKnightsCount * (WhiteKnight.Weight() - pawnFactorEG)
	whiteCentipawnsEG += whiteBishopsCount * (WhiteBishop.Weight() + pawnFactorEG)
	whiteCentipawnsEG += whiteRooksCount * (WhiteRook.Weight() + pawnFactorEG)
	whiteCentipawnsEG += whiteQueensCount * WhiteQueen.Weight()

	// mobility and attacks
	whitePawnAttacks, whiteMinorAttacks, whiteOtherAttacks := board.AllAttacks(White) // get the squares that are attacked by white
	blackPawnAttacks, blackMinorAttacks, blackOtherAttacks := board.AllAttacks(Black) // get the squares that are attacked by black

	blackKingZone := SquareInnerRingMask[blackKingIndex] | SquareOuterRingMask[blackKingIndex]
	whiteKingZone := SquareInnerRingMask[whiteKingIndex] | SquareOuterRingMask[whiteKingIndex]

	// king attacks are considered later
	whiteAttacks := (whitePawnAttacks | whiteMinorAttacks | whiteOtherAttacks) &^ blackKingZone
	blackAttacks := (blackPawnAttacks | blackMinorAttacks | blackOtherAttacks) &^ whiteKingZone

	wQuietAttacks := bits.OnesCount64(whiteAttacks << 32) // keep hi-bits only
	bQuietAttacks := bits.OnesCount64(blackAttacks >> 32) // keep lo-bits only

	whiteAggressivity := bits.OnesCount64(whiteAttacks >> 32) // keep hi-bits only
	blackAggressivity := bits.OnesCount64(blackAttacks << 32) // keep lo-bits only

	whiteCentipawnsMG += MiddlegameMobilityFactorCoeff * int16(wQuietAttacks)
	whiteCentipawnsEG += EndgameMobilityFactorCoeff * int16(wQuietAttacks)

	blackCentipawnsMG += MiddlegameMobilityFactorCoeff * int16(bQuietAttacks)
	blackCentipawnsEG += EndgameMobilityFactorCoeff * int16(bQuietAttacks)

	whiteCentipawnsMG += MiddlegameAggressivityFactorCoeff * int16(whiteAggressivity)
	whiteCentipawnsEG += EndgameAggressivityFactorCoeff * int16(whiteAggressivity)

	blackCentipawnsMG += MiddlegameAggressivityFactorCoeff * int16(blackAggressivity)
	blackCentipawnsEG += EndgameAggressivityFactorCoeff * int16(blackAggressivity)

	// whiteCentipawns += 2 * int16(wQuietAttacks)
	// blackCentipawns += 2 * int16(bQuietAttacks)
	//
	// whiteCentipawns += 4 * int16(whiteAggressivity)
	// blackCentipawns += 4 * int16(blackAggressivity)
	//
	// if !isEndgame {
	//

	whiteCentipawnsMG +=
		MiddlegameInnerPawnToKingAttackCoeff*int16(bits.OnesCount64(whitePawnAttacks&SquareInnerRingMask[blackKingIndex])) +
			MiddlegameOuterPawnToKingAttackCoeff*int16(bits.OnesCount64(whitePawnAttacks&SquareOuterRingMask[blackKingIndex])) +
			MiddlegameInnerMinorToKingAttackCoeff*int16(bits.OnesCount64(whiteMinorAttacks&SquareInnerRingMask[blackKingIndex])) +
			MiddlegameOuterMinorToKingAttackCoeff*int16(bits.OnesCount64(whiteMinorAttacks&SquareOuterRingMask[blackKingIndex])) +
			MiddlegameInnerMajorToKingAttackCoeff*int16(bits.OnesCount64(whiteOtherAttacks&SquareInnerRingMask[blackKingIndex])) +
			MiddlegameOuterMajorToKingAttackCoeff*int16(bits.OnesCount64(whiteOtherAttacks&SquareOuterRingMask[blackKingIndex]))

	whiteCentipawnsEG +=
		EndgameInnerPawnToKingAttackCoeff*int16(bits.OnesCount64(whitePawnAttacks&SquareInnerRingMask[blackKingIndex])) +
			EndgameOuterPawnToKingAttackCoeff*int16(bits.OnesCount64(whitePawnAttacks&SquareOuterRingMask[blackKingIndex])) +
			EndgameInnerMinorToKingAttackCoeff*int16(bits.OnesCount64(whiteMinorAttacks&SquareInnerRingMask[blackKingIndex])) +
			EndgameOuterMinorToKingAttackCoeff*int16(bits.OnesCount64(whiteMinorAttacks&SquareOuterRingMask[blackKingIndex])) +
			EndgameInnerMajorToKingAttackCoeff*int16(bits.OnesCount64(whiteOtherAttacks&SquareInnerRingMask[blackKingIndex])) +
			EndgameOuterMajorToKingAttackCoeff*int16(bits.OnesCount64(whiteOtherAttacks&SquareOuterRingMask[blackKingIndex]))

	blackCentipawnsMG +=
		MiddlegameInnerPawnToKingAttackCoeff*int16(bits.OnesCount64(blackPawnAttacks&SquareInnerRingMask[whiteKingIndex])) +
			MiddlegameOuterPawnToKingAttackCoeff*int16(bits.OnesCount64(blackPawnAttacks&SquareOuterRingMask[whiteKingIndex])) +
			MiddlegameInnerMinorToKingAttackCoeff*int16(bits.OnesCount64(blackMinorAttacks&SquareInnerRingMask[whiteKingIndex])) +
			MiddlegameOuterMinorToKingAttackCoeff*int16(bits.OnesCount64(blackMinorAttacks&SquareOuterRingMask[whiteKingIndex])) +
			MiddlegameInnerMajorToKingAttackCoeff*int16(bits.OnesCount64(blackOtherAttacks&SquareInnerRingMask[whiteKingIndex])) +
			MiddlegameOuterMajorToKingAttackCoeff*int16(bits.OnesCount64(blackOtherAttacks&SquareOuterRingMask[whiteKingIndex]))

	blackCentipawnsEG +=
		EndgameInnerPawnToKingAttackCoeff*int16(bits.OnesCount64(blackPawnAttacks&SquareInnerRingMask[whiteKingIndex])) +
			EndgameOuterPawnToKingAttackCoeff*int16(bits.OnesCount64(blackPawnAttacks&SquareOuterRingMask[whiteKingIndex])) +
			EndgameInnerMinorToKingAttackCoeff*int16(bits.OnesCount64(blackMinorAttacks&SquareInnerRingMask[whiteKingIndex])) +
			EndgameOuterMinorToKingAttackCoeff*int16(bits.OnesCount64(blackMinorAttacks&SquareOuterRingMask[whiteKingIndex])) +
			EndgameInnerMajorToKingAttackCoeff*int16(bits.OnesCount64(blackOtherAttacks&SquareInnerRingMask[whiteKingIndex])) +
			EndgameOuterMajorToKingAttackCoeff*int16(bits.OnesCount64(blackOtherAttacks&SquareOuterRingMask[whiteKingIndex]))

		//
	// blackAttacksToKing :=
	// 	7*bits.OnesCount64(blackPawnAttacks&SquareInnerRingMask[whiteKingIndex]) +
	// 		7*bits.OnesCount64(blackPawnAttacks&SquareOuterRingMask[whiteKingIndex]) +
	// 		7*bits.OnesCount64(blackMinorAttacks&SquareInnerRingMask[whiteKingIndex]) +
	// 		6*bits.OnesCount64(blackMinorAttacks&SquareOuterRingMask[whiteKingIndex]) +
	// 		5*bits.OnesCount64(blackOtherAttacks&SquareInnerRingMask[whiteKingIndex]) +
	// 		4*bits.OnesCount64(blackOtherAttacks&SquareOuterRingMask[whiteKingIndex])
	//
	// blackCentipawns += int16(blackAttacksToKing)
	// whiteCentipawns += int16(whiteAttacksToKing)
	//
	// blackCentipawns += blackKingSafetyCentiPawns
	// whiteCentipawns += whiteKingSafetyCentiPawns
	// }

	// whiteAttacks := board.AllAttacksOn(Black) // get the squares that are taboo for black (white's reach)
	// blackAttacks := board.AllAttacksOn(White) // get the squares that are taboo for whtie (black's reach)
	// wAttackCounts := bits.OnesCount64(whiteAttacks)
	// bAttackCounts := bits.OnesCount64(blackAttacks)
	//
	// whiteAggressivity := bits.OnesCount64(whiteAttacks >> 32) // keep hi-bits only (black's half)
	// blackAggressivity := bits.OnesCount64(blackAttacks << 32) // keep lo-bits only (white's half)
	//
	// whiteCentipawnsMG += MiddlegameMobilityFactorCoeff * int16(wAttackCounts)
	// whiteCentipawnsEG += EndgameMobilityFactorCoeff * int16(wAttackCounts)
	//
	// blackCentipawnsMG += MiddlegameMobilityFactorCoeff * int16(bAttackCounts)
	// blackCentipawnsEG += EndgameMobilityFactorCoeff * int16(bAttackCounts)
	//
	// whiteCentipawnsMG += MiddlegameAggressivityFactorCoeff * int16(whiteAggressivity)
	// whiteCentipawnsEG += EndgameAggressivityFactorCoeff * int16(whiteAggressivity)
	//
	// blackCentipawnsMG += MiddlegameAggressivityFactorCoeff * int16(blackAggressivity)
	// blackCentipawnsEG += EndgameAggressivityFactorCoeff * int16(blackAggressivity)
	//
	phase := TotalPhase -
		whitePawnsCount*PawnPhase -
		blackPawnsCount*PawnPhase -
		whiteKnightsCount*KnightPhase -
		blackKnightsCount*KnightPhase -
		whiteBishopsCount*BishopPhase -
		blackBishopsCount*BishopPhase -
		whiteRooksCount*RookPhase -
		blackRooksCount*RookPhase -
		whiteQueensCount*QueenPhase -
		blackQueensCount*QueenPhase

	phase = (phase*256 + HalfPhase) / TotalPhase

	var evalEG, evalMG int16

	if turn == White {
		evalEG = whiteCentipawnsEG - blackCentipawnsEG
		evalMG = whiteCentipawnsMG - blackCentipawnsMG
	} else {
		evalEG = blackCentipawnsEG - whiteCentipawnsEG
		evalMG = blackCentipawnsMG - whiteCentipawnsMG
	}

	// The following formula overflows if I do not convert to int32 first
	// then I have to convert back to int16, as the function return requires
	// and that is also safe, due to the division
	mg := int32(evalMG)
	eg := int32(evalEG)
	phs := int32(phase)
	taperedEval := int16(((mg * (256 - phs)) + eg*phs) / 256)
	return toEval(taperedEval + Tempo)
}

func toEval(eval int16) int16 {
	if eval >= CHECKMATE_EVAL {
		return MAX_NON_CHECKMATE
	} else if eval <= -CHECKMATE_EVAL {
		return -MAX_NON_CHECKMATE
	}
	return eval

}
