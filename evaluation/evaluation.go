package evaluation

import (
	"math/bits"

	. "github.com/amanjpro/zahak/engine"
)

type Eval struct {
	blackMG int16
	whiteMG int16
	blackEG int16
	whiteEG int16
}

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

const BlackKingSideMask = uint64(1<<F8 | 1<<G8 | 1<<H8 | 1<<F7 | 1<<G7 | 1<<H7)
const WhiteKingSideMask = uint64(1<<F1 | 1<<G1 | 1<<H1 | 1<<F2 | 1<<G2 | 1<<H2)
const BlackQueenSideMask = uint64(1<<C8 | 1<<B8 | 1<<A8 | 1<<A7 | 1<<B7 | 1<<C7)
const WhiteQueenSideMask = uint64(1<<C1 | 1<<B1 | 1<<A1 | 1<<A2 | 1<<B2 | 1<<C2)

const BlackAShield = uint64(1<<A7 | 1<<A6)
const BlackBShield = uint64(1<<B7 | 1<<B6)
const BlackCShield = uint64(1<<C7 | 1<<C6)
const BlackFShield = uint64(1<<F7 | 1<<F6)
const BlackGShield = uint64(1<<G7 | 1<<G6)
const BlackHShield = uint64(1<<H7 | 1<<H6)
const WhiteAShield = uint64(1<<A2 | 1<<A3)
const WhiteBShield = uint64(1<<B2 | 1<<B3)
const WhiteCShield = uint64(1<<C2 | 1<<C3)
const WhiteFShield = uint64(1<<F2 | 1<<F3)
const WhiteGShield = uint64(1<<G2 | 1<<G3)
const WhiteHShield = uint64(1<<H2 | 1<<H3)

// Piece Square Tables
// Middle-game
var EarlyPawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	96, 131, 66, 105, 97, 126, 15, -31,
	-14, -13, 21, 21, 62, 74, 16, -18,
	-24, -6, -3, 16, 15, 12, 1, -29,
	-35, -23, -13, 5, 10, 3, -9, -34,
	-35, -30, -20, -22, -10, -10, 6, -25,
	-43, -25, -33, -35, -30, 9, 11, -33,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var EarlyKnightPst = [64]int16{
	-211, -81, -50, -44, 54, -119, -26, -140,
	-79, -39, 79, 39, 30, 75, 4, -12,
	-39, 71, 48, 67, 95, 139, 69, 55,
	9, 32, 25, 57, 27, 74, 17, 36,
	8, 34, 36, 25, 46, 27, 34, 9,
	-2, 13, 29, 34, 49, 36, 48, 9,
	-4, -23, 14, 24, 24, 40, 19, 17,
	-99, 7, -29, -14, 18, 4, 13, 5,
}

var EarlyBishopPst = [64]int16{
	-25, 24, -83, -45, -22, -36, 11, 12,
	5, 50, 14, 0, 57, 77, 41, -25,
	20, 69, 81, 66, 63, 78, 52, 25,
	32, 39, 42, 72, 59, 53, 35, 24,
	36, 52, 48, 60, 68, 48, 46, 40,
	33, 58, 55, 53, 57, 76, 58, 43,
	46, 65, 59, 47, 57, 66, 83, 47,
	6, 41, 37, 29, 40, 34, 3, 16,
}

var EarlyRookPst = [64]int16{
	-3, 14, -19, 21, 19, -19, 0, -9,
	3, 0, 33, 31, 59, 59, -2, 20,
	-36, -10, -5, -5, -29, 24, 43, -17,
	-43, -32, -15, 2, -16, 14, -19, -38,
	-51, -46, -28, -26, -13, -25, 0, -38,
	-52, -35, -32, -33, -18, -7, -15, -34,
	-50, -26, -34, -24, -15, 3, -15, -73,
	-19, -20, -12, -2, 0, -1, -38, -18,
}

var EarlyQueenPst = [64]int16{
	-64, -33, -12, -17, 36, 33, 32, 12,
	-35, -62, -24, -24, -65, 27, -5, 20,
	-18, -24, -14, -42, -7, 28, 1, 19,
	-39, -38, -38, -49, -33, -21, -34, -27,
	-10, -41, -24, -24, -24, -18, -17, -14,
	-25, 7, -10, -3, -5, 0, 6, -3,
	-25, -2, 19, 11, 17, 24, 7, 16,
	13, -4, 9, 25, -2, -8, -10, -31,
}

var EarlyKingPst = [64]int16{
	-52, 95, 93, 48, -45, -13, 44, 46,
	97, 33, 25, 63, 19, 4, -21, -61,
	35, 38, 53, 8, 26, 60, 56, -14,
	-21, -8, 10, -21, -19, -20, -22, -56,
	-45, 9, -36, -72, -72, -50, -60, -78,
	-6, -15, -27, -58, -59, -46, -21, -37,
	14, 17, -10, -61, -40, -15, 10, 17,
	-9, 37, 15, -66, -8, -35, 28, 26,
}

// Endgame
var LatePawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	169, 150, 134, 103, 116, 104, 155, 191,
	87, 85, 61, 34, 18, 24, 61, 73,
	23, 4, -7, -28, -19, -12, 1, 9,
	20, 9, -1, -9, -9, -7, -2, 4,
	4, 0, -8, -5, -4, -7, -14, -12,
	15, -1, 5, 3, 11, -8, -13, -10,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var LateKnightPst = [64]int16{
	-44, -50, -15, -39, -44, -33, -74, -99,
	-30, -14, -47, -18, -28, -51, -34, -61,
	-36, -39, -12, -13, -29, -37, -38, -62,
	-24, -9, 7, 9, 15, -6, -2, -29,
	-27, -24, 0, 13, 2, 4, -6, -24,
	-31, -14, -18, 0, -7, -18, -33, -34,
	-42, -27, -19, -15, -12, -30, -31, -54,
	-23, -57, -31, -15, -28, -27, -58, -74,
}

var LateBishopPst = [64]int16{
	-10, -26, -4, -7, -4, -7, -15, -26,
	-11, -16, -3, -15, -15, -22, -15, -13,
	-2, -20, -17, -18, -17, -16, -9, -2,
	-5, 1, 2, -4, -1, -3, -8, 1,
	-12, -9, 2, 5, -10, -2, -14, -12,
	-11, -9, 0, 0, 4, -13, -10, -15,
	-16, -23, -13, -4, -2, -15, -20, -32,
	-19, -9, -18, -4, -10, -12, -4, -13,
}

var LateRookPst = [64]int16{
	8, 2, 15, 3, 5, 13, 6, 7,
	8, 12, 2, 3, -17, -6, 12, 3,
	12, 9, 2, 4, 5, -6, -8, 2,
	15, 9, 15, 0, 7, 6, 2, 15,
	15, 19, 16, 12, 4, 8, -3, 5,
	12, 12, 7, 12, 2, -2, 4, -2,
	10, 4, 12, 12, 1, -3, -4, 15,
	4, 11, 9, 0, -3, -4, 10, -14,
}

var LateQueenPst = [64]int16{
	32, 59, 50, 50, 41, 36, 25, 58,
	5, 45, 51, 63, 91, 41, 57, 39,
	-1, 22, 18, 84, 67, 49, 51, 37,
	43, 44, 44, 71, 77, 60, 92, 68,
	0, 57, 42, 67, 52, 55, 67, 53,
	26, -24, 31, 19, 30, 36, 47, 40,
	-4, -9, -20, 0, 4, -1, -15, -12,
	-19, -17, -7, -30, 21, -13, 1, -28,
}

var LateKingPst = [64]int16{
	-72, -51, -34, -29, -4, 17, -2, -15,
	-29, 1, 1, -3, 6, 31, 16, 20,
	3, 5, 5, 6, 4, 30, 29, 15,
	-7, 13, 15, 23, 19, 28, 21, 10,
	-14, -14, 19, 27, 29, 22, 9, 0,
	-20, -8, 9, 22, 24, 18, 3, -3,
	-28, -19, 5, 12, 14, 4, -11, -21,
	-54, -48, -23, 3, -26, -4, -37, -56,
}

var MiddlegameBackwardPawnPenalty int16 = 10
var EndgameBackwardPawnPenalty int16 = 4
var MiddlegameIsolatedPawnPenalty int16 = 11
var EndgameIsolatedPawnPenalty int16 = 6
var MiddlegameDoublePawnPenalty int16 = 2
var EndgameDoublePawnPenalty int16 = 26
var MiddlegamePassedPawnAward int16 = 0
var EndgamePassedPawnAward int16 = 10
var MiddlegameAdvancedPassedPawnAward int16 = 10
var EndgameAdvancedPassedPawnAward int16 = 62
var MiddlegameCandidatePassedPawnAward int16 = 32
var EndgameCandidatePassedPawnAward int16 = 51
var MiddlegameRookOpenFileAward int16 = 45
var EndgameRookOpenFileAward int16 = 0
var MiddlegameRookSemiOpenFileAward int16 = 14
var EndgameRookSemiOpenFileAward int16 = 20
var MiddlegameVeritcalDoubleRookAward int16 = 10
var EndgameVeritcalDoubleRookAward int16 = 10
var MiddlegameHorizontalDoubleRookAward int16 = 27
var EndgameHorizontalDoubleRookAward int16 = 12
var MiddlegamePawnFactorCoeff int16 = 0
var EndgamePawnFactorCoeff int16 = 0
var MiddlegameMobilityFactorCoeff int16 = 6
var EndgameMobilityFactorCoeff int16 = 3
var MiddlegameAggressivityFactorCoeff int16 = 1
var EndgameAggressivityFactorCoeff int16 = 6
var MiddlegameInnerPawnToKingAttackCoeff int16 = 0
var EndgameInnerPawnToKingAttackCoeff int16 = 0
var MiddlegameOuterPawnToKingAttackCoeff int16 = 4
var EndgameOuterPawnToKingAttackCoeff int16 = 0
var MiddlegameInnerMinorToKingAttackCoeff int16 = 18
var EndgameInnerMinorToKingAttackCoeff int16 = 0
var MiddlegameOuterMinorToKingAttackCoeff int16 = 11
var EndgameOuterMinorToKingAttackCoeff int16 = 1
var MiddlegameInnerMajorToKingAttackCoeff int16 = 17
var EndgameInnerMajorToKingAttackCoeff int16 = 0
var MiddlegameOuterMajorToKingAttackCoeff int16 = 8
var EndgameOuterMajorToKingAttackCoeff int16 = 5
var MiddlegamePawnShieldPenalty int16 = 10
var EndgamePawnShieldPenalty int16 = 8
var MiddlegameNotCastlingPenalty int16 = 26
var EndgameNotCastlingPenalty int16 = 5
var MiddlegameKingZoneOpenFilePenalty int16 = 35
var EndgameKingZoneOpenFilePenalty int16 = 0
var MiddlegameKingZoneMissingPawnPenalty int16 = 16
var EndgameKingZoneMissingPawnPenalty int16 = 0

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
	pieceIter := bbBlackPawn
	for pieceIter > 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		blackPawnsCount++
		blackCentipawnsEG += LatePawnPst[index]
		blackCentipawnsMG += EarlyPawnPst[index]
		pieceIter ^= mask
	}

	// PST for white pawns
	pieceIter = bbWhitePawn
	for pieceIter != 0 {
		whitePawnsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask(uint64(index))
		whiteCentipawnsEG += LatePawnPst[flip[index]]
		whiteCentipawnsMG += EarlyPawnPst[flip[index]]
		pieceIter ^= mask
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
		blackCentipawnsMG += EarlyKingPst[index]
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
		whiteCentipawnsMG += EarlyKingPst[flip[index]]
		whiteKingIndex = index
		pieceIter ^= mask
	}

	pawnFactorMG := int16(16-blackPawnsCount-whitePawnsCount) * MiddlegamePawnFactorCoeff
	pawnFactorEG := int16(16-blackPawnsCount-whitePawnsCount) * EndgamePawnFactorCoeff

	blackCentipawnsMG += blackPawnsCount * BlackPawn.Weight()
	blackCentipawnsMG += blackKnightsCount * (BlackKnight.Weight() - pawnFactorMG)
	blackCentipawnsMG += blackBishopsCount * (BlackBishop.Weight())
	blackCentipawnsMG += blackRooksCount * (BlackRook.Weight() + pawnFactorMG)
	blackCentipawnsMG += blackQueensCount * BlackQueen.Weight()

	blackCentipawnsEG += blackPawnsCount * BlackPawn.Weight()
	blackCentipawnsEG += blackKnightsCount * (BlackKnight.Weight() - pawnFactorEG)
	blackCentipawnsEG += blackBishopsCount * (BlackBishop.Weight())
	blackCentipawnsEG += blackRooksCount * (BlackRook.Weight() + pawnFactorEG)
	blackCentipawnsEG += blackQueensCount * BlackQueen.Weight()

	whiteCentipawnsMG += whitePawnsCount * WhitePawn.Weight()
	whiteCentipawnsMG += whiteKnightsCount * (WhiteKnight.Weight() - pawnFactorMG)
	whiteCentipawnsMG += whiteBishopsCount * (WhiteBishop.Weight())
	whiteCentipawnsMG += whiteRooksCount * (WhiteRook.Weight() + pawnFactorMG)
	whiteCentipawnsMG += whiteQueensCount * WhiteQueen.Weight()

	whiteCentipawnsEG += whitePawnsCount * WhitePawn.Weight()
	whiteCentipawnsEG += whiteKnightsCount * (WhiteKnight.Weight() - pawnFactorEG)
	whiteCentipawnsEG += whiteBishopsCount * (WhiteBishop.Weight())
	whiteCentipawnsEG += whiteRooksCount * (WhiteRook.Weight() + pawnFactorEG)
	whiteCentipawnsEG += whiteQueensCount * WhiteQueen.Weight()

	mobilityEval := Mobility(position, blackKingIndex, whiteKingIndex)

	whiteCentipawnsMG += mobilityEval.whiteMG
	whiteCentipawnsEG += mobilityEval.whiteEG
	blackCentipawnsMG += mobilityEval.blackMG
	blackCentipawnsEG += mobilityEval.blackEG

	rookEval := RookFilesEval(bbBlackRook, bbWhiteRook, bbBlackPawn, bbWhitePawn)
	whiteCentipawnsMG += rookEval.whiteMG
	whiteCentipawnsEG += rookEval.whiteEG
	blackCentipawnsMG += rookEval.blackMG
	blackCentipawnsEG += rookEval.blackEG

	pawnStructureEval := PawnStructureEval(position)
	whiteCentipawnsMG += pawnStructureEval.whiteMG
	whiteCentipawnsEG += pawnStructureEval.whiteEG
	blackCentipawnsMG += pawnStructureEval.blackMG
	blackCentipawnsEG += pawnStructureEval.blackEG

	kingSafetyEval := KingSafety(bbBlackKing, bbWhiteKing, bbBlackPawn, bbWhitePawn,
		position.HasTag(BlackCanCastleQueenSide) || position.HasTag(BlackCanCastleKingSide),
		position.HasTag(WhiteCanCastleQueenSide) || position.HasTag(WhiteCanCastleKingSide),
	)
	whiteCentipawnsMG += kingSafetyEval.whiteMG
	whiteCentipawnsEG += kingSafetyEval.whiteEG
	blackCentipawnsMG += kingSafetyEval.blackMG
	blackCentipawnsEG += kingSafetyEval.blackEG

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

func RookFilesEval(blackRook uint64, whiteRook uint64, blackPawns uint64, whitePawns uint64) Eval {
	var blackMG, whiteMG, blackEG, whiteEG int16

	blackFiles := FileFill(blackRook)
	whiteFiles := FileFill(whiteRook)

	allPawns := FileFill(blackPawns | whitePawns)

	// open files
	blackRooksNoPawns := blackFiles &^ allPawns
	whiteRooksNoPawns := whiteFiles &^ allPawns

	blackRookOpenFiles := blackRook & blackRooksNoPawns
	whiteRookOpenFiles := whiteRook & whiteRooksNoPawns

	count := int16(bits.OnesCount64(blackRookOpenFiles))
	blackMG += MiddlegameRookOpenFileAward * count
	blackEG += EndgameRookOpenFileAward * count

	count = int16(bits.OnesCount64(whiteRookOpenFiles))
	whiteMG += MiddlegameRookOpenFileAward * count
	whiteEG += EndgameRookOpenFileAward * count

	// semi-open files
	blackRooksNoOwnPawns := blackFiles &^ FileFill(blackPawns)
	whiteRooksNoOwnPawns := whiteFiles &^ FileFill(whitePawns)

	blackRookSemiOpenFiles := (blackRook &^ blackRookOpenFiles) & blackRooksNoOwnPawns
	whiteRookSemiOpenFiles := (whiteRook &^ whiteRookOpenFiles) & whiteRooksNoOwnPawns

	count = int16(bits.OnesCount64(blackRookSemiOpenFiles))
	blackMG += MiddlegameRookSemiOpenFileAward * count
	blackEG += EndgameRookSemiOpenFileAward * count

	count = int16(bits.OnesCount64(whiteRookSemiOpenFiles))
	whiteMG += MiddlegameRookSemiOpenFileAward * count
	whiteEG += EndgameRookSemiOpenFileAward * count

	return Eval{blackMG: blackMG, whiteMG: whiteMG, blackEG: blackEG, whiteEG: whiteEG}
}

func PawnStructureEval(p *Position) Eval {
	var blackMG, whiteMG, blackEG, whiteEG int16

	// passed pawns
	countP, countS := p.CountPassedPawns(Black)
	blackMG += MiddlegamePassedPawnAward * countP
	blackEG += EndgamePassedPawnAward * countP

	blackMG += MiddlegameAdvancedPassedPawnAward * countS
	blackEG += EndgameAdvancedPassedPawnAward * countS

	countP, countS = p.CountPassedPawns(White)
	whiteMG += MiddlegamePassedPawnAward * countP
	whiteEG += EndgamePassedPawnAward * countP

	whiteMG += MiddlegameAdvancedPassedPawnAward * countS
	whiteEG += EndgameAdvancedPassedPawnAward * countS

	// candidate passed pawns
	count := p.CountCandidatePawns(Black)
	blackMG += MiddlegameCandidatePassedPawnAward * count
	blackEG += EndgameCandidatePassedPawnAward * count

	count = p.CountCandidatePawns(White)
	whiteMG += MiddlegameCandidatePassedPawnAward * count
	whiteEG += EndgameCandidatePassedPawnAward * count

	// backward pawns
	count = p.CountBackwardPawns(Black)
	blackMG -= MiddlegameBackwardPawnPenalty * count
	blackEG -= EndgameBackwardPawnPenalty * count

	count = p.CountBackwardPawns(White)
	whiteMG -= MiddlegameBackwardPawnPenalty * count
	whiteEG -= EndgameBackwardPawnPenalty * count

	// isolated pawns
	count = p.CountIsolatedPawns(Black)
	blackMG -= MiddlegameIsolatedPawnPenalty * count
	blackEG -= EndgameIsolatedPawnPenalty * count

	count = p.CountIsolatedPawns(White)
	whiteMG -= MiddlegameIsolatedPawnPenalty * count
	whiteEG -= EndgameIsolatedPawnPenalty * count

	// double pawns
	count = p.CountDoublePawns(Black)
	blackMG -= MiddlegameDoublePawnPenalty * count
	blackEG -= EndgameDoublePawnPenalty * count

	count = p.CountDoublePawns(White)
	whiteMG -= MiddlegameDoublePawnPenalty * count
	whiteEG -= EndgameDoublePawnPenalty * count

	return Eval{blackMG: blackMG, whiteMG: whiteMG, blackEG: blackEG, whiteEG: whiteEG}
}

func kingSafetyPenalty(color Color, side PieceType, ownPawn uint64, allPawn uint64) (int16, int16) {
	var mg, eg int16
	var a_shield, b_shield, c_shield, f_shield, g_shield, h_shield uint64
	if color == White {
		a_shield = WhiteAShield
		b_shield = WhiteBShield
		c_shield = WhiteCShield
		f_shield = WhiteFShield
		g_shield = WhiteGShield
		h_shield = WhiteHShield
	} else {
		a_shield = BlackAShield
		b_shield = BlackBShield
		c_shield = BlackCShield
		f_shield = BlackFShield
		g_shield = BlackGShield
		h_shield = BlackHShield
	}
	if side == King {
		if H_FileFill&allPawn == 0 { // no pawns, super bad
			mg += MiddlegameKingZoneOpenFilePenalty
			eg += EndgameKingZoneOpenFilePenalty
		} else if H_FileFill&ownPawn == 0 { // semi-open file, bad
			mg += MiddlegameKingZoneMissingPawnPenalty
			eg += EndgameKingZoneMissingPawnPenalty
		} else if h_shield&ownPawn == 0 {
			mg += MiddlegamePawnShieldPenalty
			eg += EndgamePawnShieldPenalty
		}

		if G_FileFill&allPawn == 0 { // no pawns, super bad
			mg += MiddlegameKingZoneOpenFilePenalty
			eg += EndgameKingZoneOpenFilePenalty
		} else if G_FileFill&ownPawn == 0 { // semi-open file, bad
			mg += MiddlegameKingZoneMissingPawnPenalty
			eg += EndgameKingZoneMissingPawnPenalty
		} else if g_shield&ownPawn == 0 {
			mg += MiddlegamePawnShieldPenalty
			eg += EndgamePawnShieldPenalty
		}

		if F_FileFill&allPawn == 0 { // no pawns, super bad
			mg += MiddlegameKingZoneOpenFilePenalty
			eg += EndgameKingZoneOpenFilePenalty
		} else if F_FileFill&ownPawn == 0 { // semi-open file, bad
			mg += MiddlegameKingZoneMissingPawnPenalty
			eg += EndgameKingZoneMissingPawnPenalty
		} else if f_shield&ownPawn == 0 {
			mg += MiddlegamePawnShieldPenalty
			eg += EndgamePawnShieldPenalty
		}
	} else {
		if C_FileFill&allPawn == 0 { // no pawns, super bad
			mg += MiddlegameKingZoneOpenFilePenalty
			eg += EndgameKingZoneOpenFilePenalty
		} else if C_FileFill&ownPawn == 0 { // semi-open file, bad
			mg += MiddlegameKingZoneMissingPawnPenalty
			eg += EndgameKingZoneMissingPawnPenalty
		} else if c_shield&ownPawn == 0 {
			mg += MiddlegamePawnShieldPenalty
			eg += EndgamePawnShieldPenalty
		}

		if B_FileFill&allPawn == 0 { // no pawns, super bad
			mg += MiddlegameKingZoneOpenFilePenalty
			eg += EndgameKingZoneOpenFilePenalty
		} else if B_FileFill&ownPawn == 0 { // semi-open file, bad
			mg += MiddlegameKingZoneMissingPawnPenalty
			eg += EndgameKingZoneMissingPawnPenalty
		} else if b_shield&ownPawn == 0 {
			mg += MiddlegamePawnShieldPenalty
			eg += EndgamePawnShieldPenalty
		}

		if A_FileFill&allPawn == 0 { // no pawns, super bad
			mg += MiddlegameKingZoneOpenFilePenalty
			eg += EndgameKingZoneOpenFilePenalty
		} else if A_FileFill&ownPawn == 0 { // semi-open file, bad
			mg += MiddlegameKingZoneMissingPawnPenalty
			eg += EndgameKingZoneMissingPawnPenalty
		} else if a_shield&ownPawn == 0 {
			mg += MiddlegamePawnShieldPenalty
			eg += EndgamePawnShieldPenalty
		}
	}

	return mg, eg
}

func KingSafety(blackKing uint64, whiteKing uint64, blackPawn uint64,
	whitePawn uint64, blackCastleFlag bool, whiteCastleFlag bool) Eval {
	var whiteCentipawnsMG, whiteCentipawnsEG, blackCentipawnsMG, blackCentipawnsEG int16
	allPawn := whitePawn | blackPawn

	if blackKing&BlackKingSideMask != 0 {
		blackCastleFlag = true
		// Missing pawn shield
		mg, eg := kingSafetyPenalty(Black, King, blackPawn, allPawn)
		blackCentipawnsMG -= mg
		blackCentipawnsEG -= eg
	} else if blackKing&BlackQueenSideMask != 0 {
		blackCastleFlag = true
		// Missing pawn shield
		mg, eg := kingSafetyPenalty(Black, Queen, blackPawn, allPawn)
		blackCentipawnsMG -= mg
		blackCentipawnsEG -= eg
	}

	if whiteKing&WhiteKingSideMask != 0 {
		whiteCastleFlag = true
		// Missing pawn shield
		mg, eg := kingSafetyPenalty(White, King, whitePawn, allPawn)
		whiteCentipawnsMG -= mg
		whiteCentipawnsEG -= eg
	} else if whiteKing&WhiteQueenSideMask != 0 {
		whiteCastleFlag = true
		// Missing pawn shield
		mg, eg := kingSafetyPenalty(White, Queen, whitePawn, allPawn)
		whiteCentipawnsMG -= mg
		whiteCentipawnsEG -= eg
	}

	if !whiteCastleFlag {
		whiteCentipawnsMG -= MiddlegameNotCastlingPenalty
		whiteCentipawnsEG -= EndgameNotCastlingPenalty
	}

	if !blackCastleFlag {
		blackCentipawnsMG -= MiddlegameNotCastlingPenalty
		blackCentipawnsEG -= EndgameNotCastlingPenalty
	}
	return Eval{blackMG: blackCentipawnsMG, whiteMG: whiteCentipawnsMG, blackEG: blackCentipawnsEG, whiteEG: whiteCentipawnsEG}
}

func Mobility(p *Position, blackKingIndex int, whiteKingIndex int) Eval {
	board := p.Board
	var whiteCentipawnsMG, whiteCentipawnsEG, blackCentipawnsMG, blackCentipawnsEG int16

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

	return Eval{blackMG: blackCentipawnsMG, whiteMG: whiteCentipawnsMG, blackEG: blackCentipawnsEG, whiteEG: whiteCentipawnsEG}
}

func toEval(eval int16) int16 {
	if eval >= CHECKMATE_EVAL {
		return MAX_NON_CHECKMATE
	} else if eval <= -CHECKMATE_EVAL {
		return -MAX_NON_CHECKMATE
	}
	return eval
}

func min16(x int16, y int16) int16 {
	if x < y {
		return x
	}
	return y
}
