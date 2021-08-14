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
	89, 126, 68, 110, 93, 132, 14, -38,
	-11, -16, 19, 20, 63, 80, 17, -16,
	-24, -9, -9, 15, 14, 13, 0, -26,
	-39, -32, -16, 0, 5, -1, -14, -37,
	-34, -34, -18, -18, -6, -10, 4, -23,
	-44, -30, -34, -27, -26, 11, 10, -30,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var EarlyKnightPst = [64]int16{
	-187, -81, -44, -40, 63, -112, -19, -122,
	-63, -24, 96, 40, 36, 87, 17, 0,
	-26, 84, 59, 72, 103, 148, 84, 66,
	25, 50, 42, 67, 40, 86, 38, 49,
	27, 49, 51, 40, 57, 44, 47, 28,
	19, 30, 43, 49, 62, 52, 67, 25,
	15, -9, 32, 42, 44, 55, 31, 34,
	-86, 26, -12, 3, 36, 19, 29, 23,
}

var EarlyBishopPst = [64]int16{
	-13, 30, -92, -52, -30, -35, 14, 16,
	2, 44, 10, -6, 49, 75, 36, -25,
	19, 64, 74, 57, 58, 73, 46, 22,
	32, 34, 36, 68, 54, 51, 32, 21,
	33, 47, 40, 57, 62, 38, 46, 39,
	32, 57, 54, 48, 55, 74, 58, 41,
	45, 63, 56, 44, 54, 64, 81, 46,
	5, 39, 36, 26, 35, 35, 2, 18,
}

var EarlyRookPst = [64]int16{
	-4, 13, -18, 21, 23, -22, 1, -11,
	2, -8, 30, 29, 56, 58, -4, 20,
	-39, -18, -10, -9, -31, 21, 40, -22,
	-44, -31, -20, -2, -19, 14, -23, -36,
	-53, -51, -35, -30, -19, -32, -6, -43,
	-55, -35, -32, -33, -21, -10, -17, -38,
	-48, -26, -34, -25, -14, 4, -15, -74,
	-21, -20, -12, -3, -1, 1, -38, -20,
}

var EarlyQueenPst = [64]int16{
	-54, -29, -12, -12, 46, 43, 39, 19,
	-29, -56, -23, -25, -71, 29, -9, 29,
	-13, -17, -9, -45, -9, 27, 5, 20,
	-35, -31, -35, -46, -32, -24, -34, -22,
	-7, -40, -21, -23, -24, -18, -19, -16,
	-22, 12, -10, 1, -4, -1, 4, -2,
	-20, 3, 21, 18, 25, 29, 12, 22,
	16, 3, 16, 30, 3, -4, -5, -32,
}

var EarlyKingPst = [64]int16{
	-52, 116, 111, 56, -51, -18, 43, 46,
	108, 49, 26, 79, 23, 16, -14, -72,
	35, 49, 64, 17, 33, 70, 74, -13,
	-15, -4, 20, -18, -22, -21, -19, -65,
	-46, 18, -35, -72, -75, -52, -62, -85,
	-12, -12, -28, -57, -56, -47, -20, -40,
	17, 22, -9, -59, -35, -16, 12, 17,
	-7, 37, 14, -59, -10, -35, 28, 24,
}

// Endgame
var LatePawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	173, 146, 130, 100, 111, 100, 151, 191,
	83, 80, 55, 28, 9, 18, 57, 67,
	20, 1, -11, -31, -21, -16, -1, 8,
	21, 10, -2, -10, -11, -6, -2, 6,
	2, -3, -12, -10, -8, -10, -18, -14,
	14, -4, 2, -1, 2, -13, -18, -13,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var LateKnightPst = [64]int16{
	-35, -42, -13, -35, -40, -27, -70, -89,
	-25, -11, -45, -11, -25, -49, -31, -54,
	-31, -35, -8, -13, -29, -34, -32, -55,
	-21, -4, 10, 4, 6, -3, -1, -26,
	-25, -19, 1, 15, 3, 4, -4, -23,
	-29, -9, -15, 2, -5, -20, -33, -26,
	-39, -20, -17, -13, -13, -26, -27, -52,
	-13, -57, -24, -12, -28, -22, -60, -73,
}

var LateBishopPst = [64]int16{
	-22, -35, -14, -16, -14, -18, -24, -34,
	-17, -24, -12, -21, -22, -30, -24, -20,
	-11, -24, -25, -26, -25, -22, -15, -9,
	-17, -5, -6, -10, -8, -12, -17, -8,
	-21, -16, -6, -4, -19, -10, -24, -21,
	-20, -15, -6, -8, -6, -21, -16, -22,
	-26, -31, -20, -14, -13, -21, -30, -44,
	-26, -17, -26, -13, -17, -22, -11, -25,
}

var LateRookPst = [64]int16{
	11, 4, 15, 5, 7, 14, 8, 8,
	9, 14, 3, 3, -16, -7, 11, 4,
	13, 10, 3, 5, 3, -8, -9, 1,
	14, 7, 15, -1, 4, 2, 1, 13,
	15, 18, 16, 10, 3, 5, -3, 3,
	13, 10, 5, 9, 1, -6, 3, -1,
	10, 4, 11, 13, 0, -5, -4, 15,
	6, 10, 9, 1, -2, -1, 10, -13,
}

var LateQueenPst = [64]int16{
	36, 65, 59, 56, 42, 35, 27, 61,
	13, 50, 58, 73, 101, 45, 68, 42,
	6, 24, 18, 89, 71, 51, 58, 44,
	48, 52, 46, 74, 81, 65, 102, 74,
	4, 62, 49, 70, 58, 57, 71, 58,
	32, -20, 37, 26, 36, 45, 56, 48,
	3, -1, -14, 5, 8, 5, -9, -9,
	-12, -13, -5, -23, 23, -2, 3, -15,
}

var LateKingPst = [64]int16{
	-76, -61, -41, -34, -6, 16, -6, -18,
	-38, -4, -1, -7, 2, 25, 13, 19,
	-1, 1, 1, 3, 1, 25, 23, 9,
	-13, 10, 12, 20, 18, 25, 17, 8,
	-17, -17, 18, 27, 28, 20, 6, -3,
	-21, -9, 9, 22, 23, 16, 0, -4,
	-33, -20, 5, 13, 12, 5, -12, -24,
	-57, -47, -23, 1, -24, -4, -37, -57,
}

var MiddlegameBackwardPawnPenalty int16 = 10
var EndgameBackwardPawnPenalty int16 = 4
var MiddlegameIsolatedPawnPenalty int16 = 15
var EndgameIsolatedPawnPenalty int16 = 6
var MiddlegameDoublePawnPenalty int16 = 2
var EndgameDoublePawnPenalty int16 = 25
var MiddlegamePassedPawnAward int16 = 0
var EndgamePassedPawnAward int16 = 10
var MiddlegameAdvancedPassedPawnAward int16 = 11
var EndgameAdvancedPassedPawnAward int16 = 65
var MiddlegameCandidatePassedPawnAward int16 = 40
var EndgameCandidatePassedPawnAward int16 = 51
var MiddlegameRookOpenFileAward int16 = 47
var EndgameRookOpenFileAward int16 = 0
var MiddlegameRookSemiOpenFileAward int16 = 13
var EndgameRookSemiOpenFileAward int16 = 19
var MiddlegameVeritcalDoubleRookAward int16 = 11
var EndgameVeritcalDoubleRookAward int16 = 11
var MiddlegameHorizontalDoubleRookAward int16 = 28
var EndgameHorizontalDoubleRookAward int16 = 12
var MiddlegamePawnFactorCoeff int16 = 0
var EndgamePawnFactorCoeff int16 = 1
var MiddlegamePawnSquareControlCoeff int16 = 6
var EndgamePawnSquareControlCoeff int16 = 4
var MiddlegameMinorMobilityFactorCoeff int16 = 5
var EndgameMinorMobilityFactorCoeff int16 = 1
var MiddlegameMinorAggressivityFactorCoeff int16 = 4
var EndgameMinorAggressivityFactorCoeff int16 = 3
var MiddlegameMajorMobilityFactorCoeff int16 = 3
var EndgameMajorMobilityFactorCoeff int16 = 3
var MiddlegameMajorAggressivityFactorCoeff int16 = 0
var EndgameMajorAggressivityFactorCoeff int16 = 5
var MiddlegameInnerPawnToKingAttackCoeff int16 = 2
var EndgameInnerPawnToKingAttackCoeff int16 = 0
var MiddlegameOuterPawnToKingAttackCoeff int16 = 4
var EndgameOuterPawnToKingAttackCoeff int16 = 1
var MiddlegameInnerMinorToKingAttackCoeff int16 = 17
var EndgameInnerMinorToKingAttackCoeff int16 = 0
var MiddlegameOuterMinorToKingAttackCoeff int16 = 10
var EndgameOuterMinorToKingAttackCoeff int16 = 2
var MiddlegameInnerMajorToKingAttackCoeff int16 = 15
var EndgameInnerMajorToKingAttackCoeff int16 = 0
var MiddlegameOuterMajorToKingAttackCoeff int16 = 11
var EndgameOuterMajorToKingAttackCoeff int16 = 3
var MiddlegamePawnShieldPenalty int16 = 8
var EndgamePawnShieldPenalty int16 = 10
var MiddlegameNotCastlingPenalty int16 = 33
var EndgameNotCastlingPenalty int16 = 6
var MiddlegameKingZoneOpenFilePenalty int16 = 38
var EndgameKingZoneOpenFilePenalty int16 = 0
var MiddlegameKingZoneMissingPawnPenalty int16 = 15
var EndgameKingZoneMissingPawnPenalty int16 = 0
var MiddlegameKnightOutpostAward int16 = 17
var EndgameKnightOutpostAward int16 = 23
var MiddlegameBishopPairAward int16 = 28
var EndgameBishopPairAward int16 = 44

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

func Evaluate(position *Position, pawnhash *PawnCache) int16 {
	var drawDivider int16 = 0
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

	blackPawnsCount = int16(bits.OnesCount64(bbBlackPawn))
	whitePawnsCount = int16(bits.OnesCount64(bbWhitePawn))

	var whiteKingIndex, blackKingIndex int

	// PST for other black pieces
	pieceIter := bbBlackKnight
	for pieceIter != 0 {
		blackKnightsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask[index]
		blackCentipawnsEG += LateKnightPst[index]
		blackCentipawnsMG += EarlyKnightPst[index]
		pieceIter ^= mask
	}

	pieceIter = bbBlackBishop
	for pieceIter != 0 {
		blackBishopsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask[index]
		blackCentipawnsEG += LateBishopPst[index]
		blackCentipawnsMG += EarlyBishopPst[index]
		pieceIter ^= mask
	}

	pieceIter = bbBlackRook
	for pieceIter != 0 {
		blackRooksCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask[index]
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
		mask := SquareMask[index]
		blackCentipawnsEG += LateQueenPst[index]
		blackCentipawnsMG += EarlyQueenPst[index]
		pieceIter ^= mask
	}

	pieceIter = bbBlackKing
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask[index]
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
		mask := SquareMask[index]
		whiteCentipawnsEG += LateKnightPst[flip[index]]
		whiteCentipawnsMG += EarlyKnightPst[flip[index]]
		pieceIter ^= mask
	}

	pieceIter = bbWhiteBishop
	for pieceIter != 0 {
		whiteBishopsCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask[index]
		whiteCentipawnsEG += LateBishopPst[flip[index]]
		whiteCentipawnsMG += EarlyBishopPst[flip[index]]
		pieceIter ^= mask
	}

	pieceIter = bbWhiteRook
	for pieceIter != 0 {
		whiteRooksCount++
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask[index]
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
		mask := SquareMask[index]
		whiteCentipawnsEG += LateQueenPst[flip[index]]
		whiteCentipawnsMG += EarlyQueenPst[flip[index]]
		pieceIter ^= mask
	}

	pieceIter = bbWhiteKing
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask[index]
		whiteCentipawnsEG += LateKingPst[flip[index]]
		whiteCentipawnsMG += EarlyKingPst[flip[index]]
		whiteKingIndex = index
		pieceIter ^= mask
	}

	// Draw scenarios
	{
		allPiecesCount :=
			whitePawnsCount +
				blackPawnsCount +
				whiteKnightsCount +
				blackKnightsCount +
				whiteBishopsCount +
				blackBishopsCount +
				whiteRooksCount +
				blackRooksCount +
				whiteQueensCount +
				blackQueensCount

		if (allPiecesCount == 2 && whiteRooksCount == 1 && (blackKnightsCount == 1 || blackBishopsCount == 1)) ||
			(allPiecesCount == 2 && blackRooksCount == 1 && (whiteKnightsCount == 1 || whiteBishopsCount == 1)) ||
			(allPiecesCount == 2 && (blackKnightsCount == 1 || blackBishopsCount == 1) && whitePawnsCount == 1) ||
			(allPiecesCount == 2 && (whiteKnightsCount == 1 || whiteBishopsCount == 1) && blackPawnsCount == 1) ||
			(allPiecesCount == 3 && blackRooksCount == 1 && whiteRooksCount == 1 && (whiteKnightsCount == 1 || blackKnightsCount == 1 || blackBishopsCount == 1 || whiteBishopsCount == 1)) {
			drawDivider = 3
		}
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

	// Bishop Pair
	if whiteBishopsCount >= 2 {
		whiteCentipawnsMG += MiddlegameBishopPairAward
		whiteCentipawnsEG += EndgameBishopPairAward
	}
	if blackBishopsCount >= 2 {
		blackCentipawnsMG += MiddlegameBishopPairAward
		blackCentipawnsEG += EndgameBishopPairAward
	}

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

	pawnMG, pawnEG := CachedPawnStructureEval(position, pawnhash)

	kingSafetyEval := KingSafety(bbBlackKing, bbWhiteKing, bbBlackPawn, bbWhitePawn,
		position.HasTag(BlackCanCastleQueenSide) || position.HasTag(BlackCanCastleKingSide),
		position.HasTag(WhiteCanCastleQueenSide) || position.HasTag(WhiteCanCastleKingSide),
	)
	whiteCentipawnsMG += kingSafetyEval.whiteMG
	whiteCentipawnsEG += kingSafetyEval.whiteEG
	blackCentipawnsMG += kingSafetyEval.blackMG
	blackCentipawnsEG += kingSafetyEval.blackEG

	knightOutpostEval := KnightOutpostEval(position)
	whiteCentipawnsMG += knightOutpostEval.whiteMG
	whiteCentipawnsEG += knightOutpostEval.whiteEG
	blackCentipawnsMG += knightOutpostEval.blackMG
	blackCentipawnsEG += knightOutpostEval.blackEG

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
		evalEG = whiteCentipawnsEG - blackCentipawnsEG + pawnEG
		evalMG = whiteCentipawnsMG - blackCentipawnsMG + pawnMG
	} else {
		evalEG = blackCentipawnsEG - whiteCentipawnsEG - pawnEG
		evalMG = blackCentipawnsMG - whiteCentipawnsMG - pawnMG
	}

	// The following formula overflows if I do not convert to int32 first
	// then I have to convert back to int16, as the function return requires
	// and that is also safe, due to the division
	mg := int32(evalMG)
	eg := int32(evalEG)
	phs := int32(phase)
	taperedEval := int16(((mg * (256 - phs)) + eg*phs) / 256)
	return toEval(taperedEval+Tempo) >> drawDivider
}

func KnightOutpostEval(p *Position) Eval {
	var blackMG, whiteMG, blackEG, whiteEG int16
	blackOutposts := p.CountKnightOutposts(Black)
	whiteOutposts := p.CountKnightOutposts(White)

	blackMG = MiddlegameKnightOutpostAward * blackOutposts
	blackEG = EndgameKnightOutpostAward * blackOutposts
	whiteMG = MiddlegameKnightOutpostAward * whiteOutposts
	whiteEG = EndgameKnightOutpostAward * whiteOutposts

	return Eval{blackMG: blackMG, whiteMG: whiteMG, blackEG: blackEG, whiteEG: whiteEG}
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

func CachedPawnStructureEval(p *Position, pawnhash *PawnCache) (int16, int16) {
	hash := p.Pawnhash()
	mg, eg, ok := pawnhash.Get(hash)

	if ok {
		return mg, eg
	}

	eval := PawnStructureEval(p)
	mg = eval.whiteMG - eval.blackMG
	eg = eval.whiteEG - eval.blackEG
	pawnhash.Set(hash, mg, eg)

	return mg, eg
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

	// PST for black pawns
	pieceIter := p.Board.GetBitboardOf(BlackPawn)
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask[index]
		blackEG += LatePawnPst[index]
		blackMG += EarlyPawnPst[index]
		pieceIter ^= mask
	}

	// PST for white pawns
	pieceIter = p.Board.GetBitboardOf(WhitePawn)
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := SquareMask[index]
		whiteEG += LatePawnPst[flip[index]]
		whiteMG += EarlyPawnPst[flip[index]]
		pieceIter ^= mask
	}

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
	whitePawnAttacks, whiteMinorAttacks, whiteMajorAttacks := board.AllAttacks(White) // get the squares that are attacked by white
	blackPawnAttacks, blackMinorAttacks, blackMajorAttacks := board.AllAttacks(Black) // get the squares that are attacked by black

	blackKingZone := SquareInnerRingMask[blackKingIndex] | SquareOuterRingMask[blackKingIndex]
	whiteKingZone := SquareInnerRingMask[whiteKingIndex] | SquareOuterRingMask[whiteKingIndex]

	// Pawn controlled squares
	wPawnAttacks := int16(bits.OnesCount64(whitePawnAttacks &^ blackKingZone))
	bPawnAttacks := int16(bits.OnesCount64(blackPawnAttacks &^ whiteKingZone))

	whiteCentipawnsMG += MiddlegamePawnSquareControlCoeff * wPawnAttacks
	whiteCentipawnsEG += EndgamePawnSquareControlCoeff * wPawnAttacks

	blackCentipawnsMG += MiddlegamePawnSquareControlCoeff * bPawnAttacks
	blackCentipawnsEG += EndgamePawnSquareControlCoeff * bPawnAttacks

	// // Minor mobility
	wMinorAttacksNoKingZone := whiteMinorAttacks &^ blackKingZone
	bMinorAttacksNoKingZone := blackMinorAttacks &^ whiteKingZone

	wMinorQuietAttacks := int16(bits.OnesCount64(wMinorAttacksNoKingZone << 32)) // keep hi-bits only
	bMinorQuietAttacks := int16(bits.OnesCount64(bMinorAttacksNoKingZone >> 32)) // keep lo-bits only

	wMinorAggressivity := int16(bits.OnesCount64(wMinorAttacksNoKingZone >> 32)) // keep hi-bits only
	bMinorAggressivity := int16(bits.OnesCount64(bMinorAttacksNoKingZone << 32)) // keep lo-bits only

	whiteCentipawnsMG += MiddlegameMinorMobilityFactorCoeff * wMinorQuietAttacks
	whiteCentipawnsEG += EndgameMinorMobilityFactorCoeff * wMinorQuietAttacks

	blackCentipawnsMG += MiddlegameMinorMobilityFactorCoeff * bMinorQuietAttacks
	blackCentipawnsEG += EndgameMinorMobilityFactorCoeff * bMinorQuietAttacks

	whiteCentipawnsMG += MiddlegameMinorAggressivityFactorCoeff * wMinorAggressivity
	whiteCentipawnsEG += EndgameMinorAggressivityFactorCoeff * wMinorAggressivity

	blackCentipawnsMG += MiddlegameMinorAggressivityFactorCoeff * bMinorAggressivity
	blackCentipawnsEG += EndgameMinorAggressivityFactorCoeff * bMinorAggressivity

	// Major mobility
	wMajorAttacksNoKingZone := whiteMajorAttacks &^ blackKingZone
	bMajorAttacksNoKingZone := blackMajorAttacks &^ whiteKingZone

	wMajorQuietAttacks := int16(bits.OnesCount64(wMajorAttacksNoKingZone << 32)) // keep hi-bits only
	bMajorQuietAttacks := int16(bits.OnesCount64(bMajorAttacksNoKingZone >> 32)) // keep lo-bits only

	wMajorAggressivity := int16(bits.OnesCount64(wMajorAttacksNoKingZone >> 32)) // keep hi-bits only
	bMajorAggressivity := int16(bits.OnesCount64(bMajorAttacksNoKingZone << 32)) // keep lo-bits only

	whiteCentipawnsMG += MiddlegameMajorMobilityFactorCoeff * wMajorQuietAttacks
	whiteCentipawnsEG += EndgameMajorMobilityFactorCoeff * wMajorQuietAttacks

	blackCentipawnsMG += MiddlegameMajorMobilityFactorCoeff * bMajorQuietAttacks
	blackCentipawnsEG += EndgameMajorMobilityFactorCoeff * bMajorQuietAttacks

	whiteCentipawnsMG += MiddlegameMajorAggressivityFactorCoeff * wMajorAggressivity
	whiteCentipawnsEG += EndgameMajorAggressivityFactorCoeff * wMajorAggressivity

	blackCentipawnsMG += MiddlegameMajorAggressivityFactorCoeff * bMajorAggressivity
	blackCentipawnsEG += EndgameMajorAggressivityFactorCoeff * bMajorAggressivity

	// King attacks
	whiteCentipawnsMG +=
		MiddlegameInnerPawnToKingAttackCoeff*int16(bits.OnesCount64(whitePawnAttacks&SquareInnerRingMask[blackKingIndex])) +
			MiddlegameOuterPawnToKingAttackCoeff*int16(bits.OnesCount64(whitePawnAttacks&SquareOuterRingMask[blackKingIndex])) +
			MiddlegameInnerMinorToKingAttackCoeff*int16(bits.OnesCount64(whiteMinorAttacks&SquareInnerRingMask[blackKingIndex])) +
			MiddlegameOuterMinorToKingAttackCoeff*int16(bits.OnesCount64(whiteMinorAttacks&SquareOuterRingMask[blackKingIndex])) +
			MiddlegameInnerMajorToKingAttackCoeff*int16(bits.OnesCount64(whiteMajorAttacks&SquareInnerRingMask[blackKingIndex])) +
			MiddlegameOuterMajorToKingAttackCoeff*int16(bits.OnesCount64(whiteMajorAttacks&SquareOuterRingMask[blackKingIndex]))

	whiteCentipawnsEG +=
		EndgameInnerPawnToKingAttackCoeff*int16(bits.OnesCount64(whitePawnAttacks&SquareInnerRingMask[blackKingIndex])) +
			EndgameOuterPawnToKingAttackCoeff*int16(bits.OnesCount64(whitePawnAttacks&SquareOuterRingMask[blackKingIndex])) +
			EndgameInnerMinorToKingAttackCoeff*int16(bits.OnesCount64(whiteMinorAttacks&SquareInnerRingMask[blackKingIndex])) +
			EndgameOuterMinorToKingAttackCoeff*int16(bits.OnesCount64(whiteMinorAttacks&SquareOuterRingMask[blackKingIndex])) +
			EndgameInnerMajorToKingAttackCoeff*int16(bits.OnesCount64(whiteMajorAttacks&SquareInnerRingMask[blackKingIndex])) +
			EndgameOuterMajorToKingAttackCoeff*int16(bits.OnesCount64(whiteMajorAttacks&SquareOuterRingMask[blackKingIndex]))

	blackCentipawnsMG +=
		MiddlegameInnerPawnToKingAttackCoeff*int16(bits.OnesCount64(blackPawnAttacks&SquareInnerRingMask[whiteKingIndex])) +
			MiddlegameOuterPawnToKingAttackCoeff*int16(bits.OnesCount64(blackPawnAttacks&SquareOuterRingMask[whiteKingIndex])) +
			MiddlegameInnerMinorToKingAttackCoeff*int16(bits.OnesCount64(blackMinorAttacks&SquareInnerRingMask[whiteKingIndex])) +
			MiddlegameOuterMinorToKingAttackCoeff*int16(bits.OnesCount64(blackMinorAttacks&SquareOuterRingMask[whiteKingIndex])) +
			MiddlegameInnerMajorToKingAttackCoeff*int16(bits.OnesCount64(blackMajorAttacks&SquareInnerRingMask[whiteKingIndex])) +
			MiddlegameOuterMajorToKingAttackCoeff*int16(bits.OnesCount64(blackMajorAttacks&SquareOuterRingMask[whiteKingIndex]))

	blackCentipawnsEG +=
		EndgameInnerPawnToKingAttackCoeff*int16(bits.OnesCount64(blackPawnAttacks&SquareInnerRingMask[whiteKingIndex])) +
			EndgameOuterPawnToKingAttackCoeff*int16(bits.OnesCount64(blackPawnAttacks&SquareOuterRingMask[whiteKingIndex])) +
			EndgameInnerMinorToKingAttackCoeff*int16(bits.OnesCount64(blackMinorAttacks&SquareInnerRingMask[whiteKingIndex])) +
			EndgameOuterMinorToKingAttackCoeff*int16(bits.OnesCount64(blackMinorAttacks&SquareOuterRingMask[whiteKingIndex])) +
			EndgameInnerMajorToKingAttackCoeff*int16(bits.OnesCount64(blackMajorAttacks&SquareInnerRingMask[whiteKingIndex])) +
			EndgameOuterMajorToKingAttackCoeff*int16(bits.OnesCount64(blackMajorAttacks&SquareOuterRingMask[whiteKingIndex]))

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
