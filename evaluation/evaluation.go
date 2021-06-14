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

const BlackKingSideMask = uint64(1<<G8 | 1<<H8 | 1<<G7 | 1<<H7)
const WhiteKingSideMask = uint64(1<<G1 | 1<<H1 | 1<<G2 | 1<<H2)
const BlackQueenSideMask = uint64(1<<C8 | 1<<B8 | 1<<A8 | 1<<A7 | 1<<B7)
const WhiteQueenSideMask = uint64(1<<C1 | 1<<B1 | 1<<A1 | 1<<A2 | 1<<B2)

const BlackAShield = uint64(1<<A7 | 1<<A6)
const BlackBShield = uint64(1<<B7 | 1<<B6)
const BlackCShield = uint64(1 << C7)
const BlackFShield = uint64(1 << F7)
const BlackGShield = uint64(1<<G7 | 1<<G6)
const BlackHShield = uint64(1<<H7 | 1<<H6)
const WhiteAShield = uint64(1<<A2 | 1<<A3)
const WhiteBShield = uint64(1<<B2 | 1<<B3)
const WhiteCShield = uint64(1 << C2)
const WhiteFShield = uint64(1 << F2)
const WhiteGShield = uint64(1<<G2 | 1<<G3)
const WhiteHShield = uint64(1<<H2 | 1<<H3)

// Piece Square Tables
// Middle-game
var EarlyPawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	98, 132, 65, 104, 96, 126, 16, -26,
	-12, -12, 19, 17, 58, 70, 12, -19,
	-24, -6, -4, 14, 15, 10, 1, -32,
	-33, -21, -13, 5, 8, 4, -8, -34,
	-34, -29, -20, -22, -11, -4, 8, -23,
	-42, -25, -32, -34, -31, 11, 12, -32,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var EarlyKnightPst = [64]int16{
	-212, -82, -50, -52, 49, -113, -26, -135,
	-77, -45, 78, 41, 28, 74, 11, -13,
	-41, 71, 45, 64, 94, 137, 72, 53,
	5, 29, 24, 52, 25, 73, 14, 32,
	7, 31, 31, 23, 42, 26, 34, 7,
	-3, 11, 26, 31, 46, 32, 45, 7,
	-4, -25, 13, 20, 20, 36, 11, 11,
	-102, 6, -34, -15, 14, 3, 11, 5,
}

var EarlyBishopPst = [64]int16{
	-26, 20, -79, -42, -21, -32, 12, 4,
	3, 44, 11, 2, 54, 77, 38, -27,
	17, 67, 77, 63, 63, 75, 53, 23,
	28, 35, 43, 67, 56, 50, 31, 23,
	30, 52, 46, 58, 65, 45, 45, 38,
	30, 51, 54, 50, 54, 72, 54, 41,
	43, 61, 57, 43, 54, 60, 79, 44,
	5, 38, 34, 26, 35, 31, 5, 16,
}

var EarlyRookPst = [64]int16{
	-1, 15, -13, 24, 24, -11, -2, -3,
	5, 4, 34, 31, 57, 61, -3, 23,
	-32, -7, 1, -3, -24, 31, 44, -13,
	-44, -30, -13, 5, -14, 20, -18, -35,
	-50, -41, -28, -22, -13, -28, 2, -37,
	-52, -31, -30, -31, -18, -10, -12, -33,
	-48, -25, -33, -24, -14, 2, -17, -71,
	-19, -20, -13, -1, 1, 0, -36, -17,
}

var EarlyQueenPst = [64]int16{
	-66, -29, -9, -18, 34, 33, 28, 11,
	-34, -61, -24, -23, -64, 26, -7, 23,
	-18, -25, -15, -42, -8, 26, 1, 18,
	-41, -40, -36, -49, -33, -21, -35, -27,
	-10, -42, -25, -25, -25, -20, -16, -14,
	-27, 5, -11, -5, -7, -2, 6, -3,
	-25, -3, 18, 8, 14, 21, 5, 14,
	15, -7, 6, 22, -4, -8, -12, -33,
}

var EarlyKingPst = [64]int16{
	-52, 90, 84, 41, -45, -21, 40, 36,
	79, 31, 17, 51, 15, -5, -16, -64,
	37, 35, 40, 9, 21, 49, 51, -18,
	-19, -8, 13, -20, -20, -20, -20, -54,
	-51, 6, -35, -73, -71, -52, -55, -70,
	-2, -15, -26, -57, -58, -47, -21, -36,
	12, 13, -14, -64, -42, -15, 9, 16,
	-10, 38, 16, -67, -8, -31, 28, 26,
}

// Endgame
var LatePawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	170, 154, 136, 108, 119, 109, 156, 192,
	88, 85, 63, 37, 21, 27, 63, 74,
	24, 8, -4, -23, -17, -11, 3, 11,
	20, 9, -1, -10, -9, -8, -2, 4,
	5, 2, -6, -3, -4, -6, -13, -9,
	17, 0, 7, 3, 10, -6, -11, -8,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var LateKnightPst = [64]int16{
	-45, -50, -16, -38, -43, -36, -74, -102,
	-33, -13, -47, -18, -27, -50, -38, -60,
	-34, -42, -10, -13, -29, -37, -39, -61,
	-23, -9, 7, 9, 16, -7, -1, -26,
	-27, -23, 1, 14, 3, 3, -7, -24,
	-32, -14, -18, 1, -8, -19, -33, -32,
	-45, -27, -19, -15, -11, -29, -30, -52,
	-22, -56, -27, -14, -28, -27, -58, -76,
}

var LateBishopPst = [64]int16{
	-11, -25, -7, -9, -6, -8, -16, -23,
	-11, -14, -1, -16, -15, -21, -16, -14,
	-3, -18, -17, -17, -15, -12, -8, -2,
	-6, 2, 1, -2, 1, -2, -4, -1,
	-10, -10, 4, 5, -8, 2, -12, -11,
	-10, -7, 1, 1, 5, -11, -8, -13,
	-16, -23, -13, -3, -1, -12, -22, -32,
	-19, -7, -17, -4, -9, -11, -7, -15,
}

var LateRookPst = [64]int16{
	9, 5, 15, 4, 6, 13, 8, 7,
	9, 11, 3, 4, -16, -7, 11, 3,
	12, 8, 2, 5, 2, -9, -9, 1,
	15, 8, 15, 0, 6, 4, 4, 14,
	15, 17, 16, 11, 3, 6, -4, 4,
	12, 12, 8, 11, 2, -1, 3, -3,
	9, 4, 11, 13, 0, -3, -3, 14,
	4, 11, 9, 0, -3, -2, 10, -14,
}

var LateQueenPst = [64]int16{
	31, 53, 49, 50, 44, 35, 24, 57,
	3, 44, 50, 61, 91, 41, 53, 33,
	-4, 22, 19, 83, 66, 48, 51, 30,
	41, 45, 44, 67, 77, 59, 88, 67,
	-2, 56, 41, 64, 52, 54, 63, 49,
	25, -22, 30, 17, 28, 37, 43, 38,
	-9, -12, -24, 0, 3, -6, -19, -11,
	-24, -17, -9, -32, 18, -15, -1, -31,
}

var LateKingPst = [64]int16{
	-73, -52, -34, -28, -4, 17, -1, -12,
	-27, 3, 3, 1, 6, 32, 15, 21,
	4, 7, 9, 6, 6, 32, 31, 15,
	-6, 14, 16, 23, 21, 28, 21, 10,
	-14, -12, 21, 29, 30, 23, 8, -3,
	-20, -7, 10, 23, 25, 19, 3, -3,
	-28, -18, 4, 15, 16, 3, -11, -20,
	-52, -46, -23, 4, -23, -3, -37, -56,
}

var MiddlegameBackwardPawnPenalty int16 = 9
var EndgameBackwardPawnPenalty int16 = 1
var MiddlegameIsolatedPawnPenalty int16 = 11
var EndgameIsolatedPawnPenalty int16 = 7
var MiddlegameDoublePawnPenalty int16 = 3
var EndgameDoublePawnPenalty int16 = 27
var MiddlegamePassedPawnAward int16 = 0
var EndgamePassedPawnAward int16 = 9
var MiddlegameAdvancedPassedPawnAward int16 = 9
var EndgameAdvancedPassedPawnAward int16 = 57
var MiddlegameCandidatePassedPawnAward int16 = 34
var EndgameCandidatePassedPawnAward int16 = 47
var MiddlegameRookOpenFileAward int16 = 45
var EndgameRookOpenFileAward int16 = 0
var MiddlegameRookSemiOpenFileAward int16 = 14
var EndgameRookSemiOpenFileAward int16 = 20
var MiddlegameVeritcalDoubleRookAward int16 = 10
var EndgameVeritcalDoubleRookAward int16 = 10
var MiddlegameHorizontalDoubleRookAward int16 = 28
var EndgameHorizontalDoubleRookAward int16 = 8
var MiddlegamePawnFactorCoeff int16 = 0
var EndgamePawnFactorCoeff int16 = 0
var MiddlegameMobilityFactorCoeff int16 = 6
var EndgameMobilityFactorCoeff int16 = 3
var MiddlegameAggressivityFactorCoeff int16 = 1
var EndgameAggressivityFactorCoeff int16 = 5
var MiddlegameInnerPawnToKingAttackCoeff int16 = 0
var EndgameInnerPawnToKingAttackCoeff int16 = 0
var MiddlegameOuterPawnToKingAttackCoeff int16 = 3
var EndgameOuterPawnToKingAttackCoeff int16 = 0
var MiddlegameInnerMinorToKingAttackCoeff int16 = 17
var EndgameInnerMinorToKingAttackCoeff int16 = 0
var MiddlegameOuterMinorToKingAttackCoeff int16 = 10
var EndgameOuterMinorToKingAttackCoeff int16 = 1
var MiddlegameInnerMajorToKingAttackCoeff int16 = 17
var EndgameInnerMajorToKingAttackCoeff int16 = 0
var MiddlegameOuterMajorToKingAttackCoeff int16 = 7
var EndgameOuterMajorToKingAttackCoeff int16 = 5
var MiddlegamePawnShieldPenalty int16 = 10
var EndgamePawnShieldPenalty int16 = 0
var MiddlegameNotCastlingPenalty int16 = 21
var EndgameNotCastlingPenalty int16 = 2
var MiddlegameKingZoneOpenFilePenalty int16 = 24
var EndgameKingZoneOpenFilePenalty int16 = 0

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

func KingSafety(blackKing uint64, whiteKing uint64, blackPawn uint64,
	whitePawn uint64, blackCanCastle bool, whiteCanCastle bool) Eval {
	var whiteCentipawnsMG, whiteCentipawnsEG, blackCentipawnsMG, blackCentipawnsEG int16
	allPawns := whitePawn | blackPawn

	blackHasCastled := false
	whiteHasCastled := false

	queenSideOpenFiles := 3 - (min16(int16(bits.OnesCount64(A_FileFill&allPawns)), 1) +
		min16(int16(bits.OnesCount64(B_FileFill&allPawns)), 1) +
		min16(int16(bits.OnesCount64(C_FileFill&allPawns)), 1))

	kingSideOpenFiles := 3 - (min16(int16(bits.OnesCount64(H_FileFill&allPawns)), 1) +
		min16(int16(bits.OnesCount64(G_FileFill&allPawns)), 1) +
		min16(int16(bits.OnesCount64(F_FileFill&allPawns)), 1))

	if blackKing&BlackKingSideMask != 0 {
		blackHasCastled = true
		// Missing pawn shield
		if BlackHShield&blackPawn == 0 {
			blackCentipawnsMG -= MiddlegamePawnShieldPenalty
			blackCentipawnsEG -= EndgamePawnShieldPenalty
		}
		if BlackGShield&blackPawn == 0 {
			blackCentipawnsMG -= MiddlegamePawnShieldPenalty
			blackCentipawnsEG -= EndgamePawnShieldPenalty
		}
		if BlackFShield&blackPawn == 0 {
			blackCentipawnsMG -= MiddlegamePawnShieldPenalty
			blackCentipawnsEG -= EndgamePawnShieldPenalty
		}

		blackCentipawnsMG -= kingSideOpenFiles * MiddlegameKingZoneOpenFilePenalty
		blackCentipawnsEG -= kingSideOpenFiles * EndgameKingZoneOpenFilePenalty
	} else if blackKing&BlackQueenSideMask != 0 {
		blackHasCastled = true
		// Missing pawn shield
		if BlackAShield&blackPawn == 0 {
			blackCentipawnsMG -= MiddlegamePawnShieldPenalty
			blackCentipawnsEG -= EndgamePawnShieldPenalty
		}
		if BlackBShield&blackPawn == 0 {
			blackCentipawnsMG -= MiddlegamePawnShieldPenalty
			blackCentipawnsEG -= EndgamePawnShieldPenalty
		}
		if BlackCShield&blackPawn == 0 {
			blackCentipawnsMG -= MiddlegamePawnShieldPenalty
			blackCentipawnsEG -= EndgamePawnShieldPenalty
		}

		blackCentipawnsMG -= queenSideOpenFiles * MiddlegameKingZoneOpenFilePenalty
		blackCentipawnsEG -= queenSideOpenFiles * EndgameKingZoneOpenFilePenalty
	}

	if whiteKing&WhiteKingSideMask != 0 {
		whiteHasCastled = true
		// Missing pawn shield
		if WhiteHShield&whitePawn == 0 {
			whiteCentipawnsMG -= MiddlegamePawnShieldPenalty
			whiteCentipawnsEG -= EndgamePawnShieldPenalty
		}
		if WhiteGShield&whitePawn == 0 {
			whiteCentipawnsMG -= MiddlegamePawnShieldPenalty
			whiteCentipawnsEG -= EndgamePawnShieldPenalty
		}
		if WhiteFShield&whitePawn == 0 {
			whiteCentipawnsMG -= MiddlegamePawnShieldPenalty
			whiteCentipawnsEG -= EndgamePawnShieldPenalty
		}

		whiteCentipawnsMG -= kingSideOpenFiles * MiddlegameKingZoneOpenFilePenalty
		whiteCentipawnsEG -= kingSideOpenFiles * EndgameKingZoneOpenFilePenalty
	} else if whiteKing&WhiteQueenSideMask != 0 {
		whiteHasCastled = true
		// Missing pawn shield
		if WhiteAShield&whitePawn == 0 {
			whiteCentipawnsMG -= MiddlegamePawnShieldPenalty
			whiteCentipawnsEG -= EndgamePawnShieldPenalty
		}
		if WhiteBShield&whitePawn == 0 {
			whiteCentipawnsMG -= MiddlegamePawnShieldPenalty
			whiteCentipawnsEG -= EndgamePawnShieldPenalty
		}
		if WhiteCShield&whitePawn == 0 {
			whiteCentipawnsMG -= MiddlegamePawnShieldPenalty
			whiteCentipawnsEG -= EndgamePawnShieldPenalty
		}

		whiteCentipawnsMG -= queenSideOpenFiles * MiddlegameKingZoneOpenFilePenalty
		whiteCentipawnsEG -= queenSideOpenFiles * EndgameKingZoneOpenFilePenalty
	}

	if !whiteHasCastled && !whiteCanCastle {
		whiteCentipawnsMG -= MiddlegameNotCastlingPenalty
		whiteCentipawnsEG -= EndgameNotCastlingPenalty
	}

	if !blackHasCastled && !blackCanCastle {
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
