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
const BlackFShield = uint64(1<<F7 | 1<<F7)
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
	97, 128, 65, 103, 97, 124, 17, -31,
	-15, -12, 19, 20, 59, 71, 13, -20,
	-23, -7, -5, 15, 15, 11, 1, -31,
	-34, -23, -12, 5, 8, 3, -9, -34,
	-35, -30, -20, -22, -10, -9, 6, -25,
	-43, -25, -32, -34, -30, 8, 10, -33,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var EarlyKnightPst = [64]int16{
	-210, -82, -47, -45, 54, -119, -25, -139,
	-79, -41, 78, 40, 30, 75, 6, -13,
	-41, 68, 47, 66, 95, 136, 70, 54,
	7, 31, 23, 54, 26, 72, 16, 34,
	8, 33, 32, 24, 43, 26, 33, 8,
	-2, 12, 27, 34, 47, 33, 46, 8,
	-3, -26, 15, 23, 22, 37, 18, 15,
	-99, 7, -28, -13, 16, 1, 12, 3,
}

var EarlyBishopPst = [64]int16{
	-25, 22, -79, -42, -19, -37, 13, 7,
	3, 48, 13, 1, 56, 76, 41, -26,
	19, 67, 80, 65, 60, 78, 54, 24,
	31, 39, 42, 71, 58, 52, 35, 23,
	36, 51, 47, 59, 68, 48, 46, 40,
	33, 54, 54, 53, 55, 74, 56, 44,
	45, 64, 57, 45, 57, 63, 81, 45,
	5, 39, 36, 29, 39, 34, 2, 17,
}

var EarlyRookPst = [64]int16{
	-1, 14, -17, 22, 22, -17, 2, -7,
	4, 1, 33, 31, 59, 61, -3, 21,
	-35, -9, -3, -3, -28, 25, 43, -16,
	-43, -30, -15, 1, -15, 16, -18, -38,
	-51, -43, -28, -26, -12, -27, 0, -38,
	-52, -35, -31, -33, -19, -7, -13, -34,
	-50, -26, -34, -24, -15, 3, -14, -72,
	-19, -20, -12, -2, 0, -1, -37, -18,
}

var EarlyQueenPst = [64]int16{
	-64, -29, -11, -18, 35, 33, 33, 10,
	-32, -61, -24, -24, -63, 26, -4, 20,
	-18, -24, -13, -40, -8, 29, 2, 17,
	-40, -38, -37, -47, -32, -21, -33, -25,
	-9, -42, -24, -24, -23, -18, -16, -14,
	-24, 6, -9, -3, -6, 0, 7, 0,
	-24, -2, 18, 10, 17, 23, 8, 17,
	13, -6, 8, 23, -1, -8, -9, -31,
}

var EarlyKingPst = [64]int16{
	-54, 94, 90, 44, -45, -15, 43, 42,
	86, 32, 25, 58, 15, -1, -21, -61,
	35, 36, 50, 7, 25, 56, 50, -13,
	-20, -8, 8, -20, -20, -20, -21, -55,
	-46, 10, -37, -71, -71, -50, -59, -77,
	-3, -16, -25, -58, -58, -46, -21, -37,
	13, 16, -10, -61, -40, -15, 10, 17,
	-6, 37, 15, -66, -8, -33, 28, 26,
}

// Endgame
var LatePawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	167, 150, 133, 104, 117, 104, 154, 189,
	86, 84, 62, 35, 18, 25, 60, 71,
	23, 4, -7, -27, -17, -12, 1, 9,
	20, 9, -2, -9, -9, -8, -2, 4,
	4, 0, -9, -5, -4, -7, -15, -12,
	15, -1, 5, 3, 11, -8, -13, -10,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var LateKnightPst = [64]int16{
	-44, -50, -17, -40, -44, -34, -74, -98,
	-29, -13, -45, -18, -26, -51, -36, -61,
	-36, -39, -12, -12, -30, -37, -39, -62,
	-25, -8, 7, 8, 15, -6, -1, -28,
	-27, -24, 0, 13, 2, 4, -5, -25,
	-32, -13, -17, 0, -7, -18, -33, -32,
	-44, -27, -19, -15, -12, -28, -31, -53,
	-22, -57, -29, -16, -28, -25, -56, -72,
}

var LateBishopPst = [64]int16{
	-10, -25, -6, -9, -6, -7, -15, -24,
	-10, -16, -1, -15, -14, -22, -16, -13,
	-4, -19, -18, -17, -16, -15, -10, -2,
	-6, 1, 2, -3, 0, -3, -8, 1,
	-12, -10, 3, 5, -9, -2, -13, -13,
	-11, -9, 0, 1, 4, -13, -9, -15,
	-17, -24, -12, -3, -2, -14, -19, -31,
	-17, -9, -18, -4, -10, -12, -3, -13,
}

var LateRookPst = [64]int16{
	9, 3, 16, 4, 6, 14, 8, 8,
	8, 12, 2, 3, -17, -6, 12, 3,
	12, 9, 3, 4, 5, -6, -8, 2,
	15, 8, 14, 0, 5, 5, 2, 14,
	16, 17, 16, 11, 3, 6, -3, 5,
	12, 12, 7, 11, 2, -2, 3, -2,
	9, 4, 11, 12, 1, -4, -4, 15,
	4, 11, 9, 0, -3, -4, 10, -14,
}

var LateQueenPst = [64]int16{
	32, 58, 50, 50, 42, 36, 26, 58,
	5, 44, 50, 64, 90, 41, 57, 36,
	-1, 24, 19, 84, 66, 48, 51, 34,
	43, 44, 44, 69, 76, 60, 91, 68,
	0, 56, 43, 66, 51, 55, 66, 50,
	26, -24, 30, 17, 30, 36, 46, 37,
	-7, -10, -21, 0, 4, -3, -17, -12,
	-22, -17, -8, -30, 20, -14, 1, -28,
}

var LateKingPst = [64]int16{
	-72, -54, -35, -30, -5, 14, -3, -16,
	-30, 2, -1, -2, 6, 30, 15, 20,
	3, 7, 6, 6, 4, 30, 30, 13,
	-8, 12, 15, 23, 20, 28, 20, 9,
	-14, -13, 20, 28, 29, 23, 8, 0,
	-21, -8, 9, 23, 25, 18, 2, -3,
	-30, -20, 6, 14, 15, 5, -11, -21,
	-53, -48, -23, 4, -26, -4, -37, -56,
}

var MiddlegameBackwardPawnPenalty int16 = 9
var EndgameBackwardPawnPenalty int16 = 3
var MiddlegameIsolatedPawnPenalty int16 = 11
var EndgameIsolatedPawnPenalty int16 = 6
var MiddlegameDoublePawnPenalty int16 = 2
var EndgameDoublePawnPenalty int16 = 26
var MiddlegamePassedPawnAward int16 = 0
var EndgamePassedPawnAward int16 = 9
var MiddlegameAdvancedPassedPawnAward int16 = 9
var EndgameAdvancedPassedPawnAward int16 = 59
var MiddlegameCandidatePassedPawnAward int16 = 32
var EndgameCandidatePassedPawnAward int16 = 50
var MiddlegameRookOpenFileAward int16 = 45
var EndgameRookOpenFileAward int16 = 0
var MiddlegameRookSemiOpenFileAward int16 = 14
var EndgameRookSemiOpenFileAward int16 = 20
var MiddlegameVeritcalDoubleRookAward int16 = 10
var EndgameVeritcalDoubleRookAward int16 = 9
var MiddlegameHorizontalDoubleRookAward int16 = 26
var EndgameHorizontalDoubleRookAward int16 = 10
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
var MiddlegameInnerMajorToKingAttackCoeff int16 = 16
var EndgameInnerMajorToKingAttackCoeff int16 = 0
var MiddlegameOuterMajorToKingAttackCoeff int16 = 7
var EndgameOuterMajorToKingAttackCoeff int16 = 5
var MiddlegamePawnShieldPenalty int16 = 10
var EndgamePawnShieldPenalty int16 = 6
var MiddlegameNotCastlingPenalty int16 = 25
var EndgameNotCastlingPenalty int16 = 5
var MiddlegameKingZoneOpenFilePenalty int16 = 33
var EndgameKingZoneOpenFilePenalty int16 = 0
var MiddlegameKingZoneMissingPawnPenalty int16 = 18
var EndgameKingZoneMissingPawnPenalty int16 = 0

// // Middle-game
// var EarlyPawnPst = [64]int16{
// 	0, 0, 0, 0, 0, 0, 0, 0,
// 	102, 130, 69, 104, 96, 123, 19, -23,
// 	-14, -11, 17, 14, 56, 68, 12, -20,
// 	-23, -5, -3, 15, 16, 10, 1, -33,
// 	-32, -20, -12, 6, 9, 2, -8, -34,
// 	-33, -28, -19, -23, -11, -6, 10, -22,
// 	-41, -24, -32, -35, -31, 13, 13, -30,
// 	0, 0, 0, 0, 0, 0, 0, 0,
// }
//
// var EarlyKnightPst = [64]int16{
// 	-212, -82, -49, -49, 52, -115, -27, -137,
// 	-75, -45, 77, 40, 29, 74, 9, -14,
// 	-39, 70, 46, 63, 95, 137, 72, 53,
// 	4, 28, 24, 51, 23, 72, 13, 31,
// 	5, 30, 31, 22, 40, 24, 33, 6,
// 	-5, 10, 25, 30, 44, 31, 42, 5,
// 	-7, -32, 10, 18, 18, 35, 10, 11,
// 	-103, 4, -35, -18, 13, 0, 9, -2,
// }
//
// var EarlyBishopPst = [64]int16{
// 	-28, 22, -79, -41, -16, -34, 11, 4,
// 	2, 48, 12, 2, 57, 76, 40, -28,
// 	15, 66, 77, 63, 62, 75, 53, 21,
// 	27, 35, 42, 66, 56, 49, 29, 21,
// 	31, 51, 46, 57, 64, 45, 44, 36,
// 	28, 50, 52, 49, 52, 71, 53, 39,
// 	41, 60, 55, 41, 52, 59, 77, 43,
// 	0, 35, 32, 24, 34, 29, 2, 16,
// }
//
// var EarlyRookPst = [64]int16{
// 	0, 16, -10, 24, 23, -11, 0, -4,
// 	7, 5, 35, 32, 58, 61, 1, 24,
// 	-31, -5, 2, 0, -21, 33, 44, -11,
// 	-41, -29, -12, 6, -11, 20, -16, -35,
// 	-49, -42, -27, -24, -12, -29, 3, -36,
// 	-52, -31, -31, -31, -19, -10, -13, -33,
// 	-48, -25, -33, -24, -14, 1, -16, -71,
// 	-18, -19, -12, -1, 1, 0, -36, -16,
// }
//
// var EarlyQueenPst = [64]int16{
// 	-65, -27, -9, -17, 33, 34, 29, 6,
// 	-32, -59, -23, -24, -58, 28, -4, 22,
// 	-18, -24, -14, -38, -7, 27, 4, 21,
// 	-39, -38, -35, -47, -33, -20, -34, -26,
// 	-10, -42, -26, -25, -23, -21, -15, -14,
// 	-27, 5, -13, -6, -7, -3, 3, -2,
// 	-25, -4, 15, 6, 13, 21, 1, 13,
// 	14, -7, 5, 20, -6, -11, -15, -34,
// }
//
// var EarlyKingPst = [64]int16{
// 	-46, 85, 80, 41, -43, -15, 42, 37,
// 	79, 31, 20, 49, 13, -4, -16, -60,
// 	31, 31, 37, 8, 22, 48, 46, -15,
// 	-22, -10, 10, -19, -20, -21, -20, -53,
// 	-47, 6, -34, -72, -69, -50, -55, -70,
// 	-4, -15, -26, -57, -58, -47, -21, -36,
// 	8, 13, -15, -67, -45, -18, 8, 15,
// 	-14, 38, 16, -71, -1, -34, 29, 25,
// }
//
// // Endgame
// var LatePawnPst = [64]int16{
// 	0, 0, 0, 0, 0, 0, 0, 0,
// 	174, 158, 138, 111, 123, 113, 158, 194,
// 	91, 89, 70, 43, 26, 28, 67, 75,
// 	26, 9, -3, -21, -15, -9, 4, 13,
// 	21, 9, 0, -10, -9, -8, -2, 4,
// 	6, 2, -6, -2, -3, -6, -13, -8,
// 	18, 1, 9, 5, 11, -4, -9, -7,
// 	0, 0, 0, 0, 0, 0, 0, 0,
// }
//
// var LateKnightPst = [64]int16{
// 	-46, -50, -17, -39, -44, -36, -75, -100,
// 	-33, -12, -45, -17, -26, -50, -37, -61,
// 	-36, -39, -10, -11, -29, -37, -39, -61,
// 	-25, -8, 7, 10, 16, -7, -1, -26,
// 	-26, -24, 2, 13, 3, 5, -7, -25,
// 	-32, -13, -17, 0, -7, -18, -33, -32,
// 	-45, -25, -20, -14, -11, -29, -30, -52,
// 	-22, -56, -27, -14, -28, -25, -57, -72,
// }
//
// var LateBishopPst = [64]int16{
// 	-10, -25, -7, -10, -7, -9, -15, -22,
// 	-10, -15, -1, -15, -16, -22, -14, -13,
// 	-3, -18, -16, -15, -15, -12, -8, 1,
// 	-5, 2, 1, -2, 1, 0, -4, 0,
// 	-10, -9, 3, 5, -8, 1, -12, -12,
// 	-11, -7, 1, 2, 5, -11, -8, -14,
// 	-15, -23, -13, -3, 0, -13, -21, -31,
// 	-17, -8, -17, -4, -8, -11, -5, -13,
// }
//
// var LateRookPst = [64]int16{
// 	10, 5, 14, 5, 7, 13, 9, 8,
// 	9, 12, 3, 5, -14, -6, 12, 4,
// 	12, 9, 2, 5, 5, -8, -7, 0,
// 	15, 9, 16, 1, 5, 4, 4, 14,
// 	16, 18, 16, 11, 2, 6, -4, 5,
// 	12, 12, 8, 11, 2, -2, 3, -3,
// 	10, 5, 11, 13, 0, -3, -4, 15,
// 	4, 12, 9, 0, -3, -2, 10, -13,
// }
//
// var LateQueenPst = [64]int16{
// 	31, 53, 49, 52, 43, 36, 23, 56,
// 	3, 45, 50, 63, 86, 40, 52, 34,
// 	-4, 25, 20, 80, 65, 49, 49, 28,
// 	42, 45, 44, 67, 78, 55, 88, 64,
// 	-2, 54, 41, 66, 50, 53, 62, 47,
// 	24, -24, 28, 15, 27, 34, 41, 36,
// 	-12, -14, -24, -2, 2, -8, -22, -14,
// 	-24, -18, -11, -32, 17, -16, -4, -32,
// }
//
// var LateKingPst = [64]int16{
// 	-73, -52, -31, -28, -4, 17, -2, -14,
// 	-28, 5, 3, 1, 9, 33, 17, 21,
// 	5, 9, 12, 7, 8, 34, 33, 16,
// 	-7, 15, 18, 24, 21, 30, 22, 11,
// 	-14, -11, 21, 29, 31, 23, 9, -2,
// 	-20, -7, 10, 23, 26, 19, 3, -3,
// 	-31, -20, 5, 16, 17, 4, -11, -21,
// 	-58, -47, -24, 3, -24, -5, -38, -57,
// }
//
// var MiddlegameBackwardPawnPenalty int16 = 9
// var EndgameBackwardPawnPenalty int16 = 1
// var MiddlegameIsolatedPawnPenalty int16 = 11
// var EndgameIsolatedPawnPenalty int16 = 7
// var MiddlegameDoublePawnPenalty int16 = 3
// var EndgameDoublePawnPenalty int16 = 28
// var MiddlegamePassedPawnAward int16 = 0
// var EndgamePassedPawnAward int16 = 10
// var MiddlegameAdvancedPassedPawnAward int16 = 8
// var EndgameAdvancedPassedPawnAward int16 = 56
// var MiddlegameCandidatePassedPawnAward int16 = 33
// var EndgameCandidatePassedPawnAward int16 = 47
// var MiddlegameRookOpenFileAward int16 = 45
// var EndgameRookOpenFileAward int16 = 0
// var MiddlegameRookSemiOpenFileAward int16 = 13
// var EndgameRookSemiOpenFileAward int16 = 21
// var MiddlegameVeritcalDoubleRookAward int16 = 9
// var EndgameVeritcalDoubleRookAward int16 = 10
// var MiddlegameHorizontalDoubleRookAward int16 = 28
// var EndgameHorizontalDoubleRookAward int16 = 8
// var MiddlegamePawnFactorCoeff int16 = 0
// var EndgamePawnFactorCoeff int16 = 0
// var MiddlegameMobilityFactorCoeff int16 = 6
// var EndgameMobilityFactorCoeff int16 = 3
// var MiddlegameAggressivityFactorCoeff int16 = 0
// var EndgameAggressivityFactorCoeff int16 = 5
// var MiddlegameInnerPawnToKingAttackCoeff int16 = 0
// var EndgameInnerPawnToKingAttackCoeff int16 = 0
// var MiddlegameOuterPawnToKingAttackCoeff int16 = 2
// var EndgameOuterPawnToKingAttackCoeff int16 = 0
// var MiddlegameInnerMinorToKingAttackCoeff int16 = 17
// var EndgameInnerMinorToKingAttackCoeff int16 = 0
// var MiddlegameOuterMinorToKingAttackCoeff int16 = 10
// var EndgameOuterMinorToKingAttackCoeff int16 = 1
// var MiddlegameInnerMajorToKingAttackCoeff int16 = 17
// var EndgameInnerMajorToKingAttackCoeff int16 = 0
// var MiddlegameOuterMajorToKingAttackCoeff int16 = 7
// var EndgameOuterMajorToKingAttackCoeff int16 = 5
// var MiddlegamePawnShieldPenalty int16 = 10
// var EndgamePawnShieldPenalty int16 = 0
// var MiddlegameNotCastlingPenalty int16 = 21
// var EndgameNotCastlingPenalty int16 = 2
// var MiddlegameKingZoneOpenFilePenalty int16 = 24
// var EndgameKingZoneOpenFilePenalty int16 = 0
// var MiddlegameKingZoneMissingPawnPenalty int16 = 24
// var EndgameKingZoneMissingPawnPenalty int16 = 0

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
	whitePawn uint64, blackCanCastle bool, whiteCanCastle bool) Eval {
	var whiteCentipawnsMG, whiteCentipawnsEG, blackCentipawnsMG, blackCentipawnsEG int16
	allPawn := whitePawn | blackPawn

	blackHasCastled := false
	whiteHasCastled := false

	if blackKing&BlackKingSideMask != 0 {
		blackHasCastled = true
		// Missing pawn shield
		mg, eg := kingSafetyPenalty(Black, King, blackPawn, allPawn)
		blackCentipawnsMG -= mg
		blackCentipawnsEG -= eg
	} else if blackKing&BlackQueenSideMask != 0 {
		blackHasCastled = true
		// Missing pawn shield
		mg, eg := kingSafetyPenalty(Black, Queen, blackPawn, allPawn)
		blackCentipawnsMG -= mg
		blackCentipawnsEG -= eg
	}

	if whiteKing&WhiteKingSideMask != 0 {
		whiteHasCastled = true
		// Missing pawn shield
		mg, eg := kingSafetyPenalty(White, King, whitePawn, allPawn)
		whiteCentipawnsMG -= mg
		whiteCentipawnsEG -= eg
	} else if whiteKing&WhiteQueenSideMask != 0 {
		whiteHasCastled = true
		// Missing pawn shield
		mg, eg := kingSafetyPenalty(White, Queen, whitePawn, allPawn)
		whiteCentipawnsMG -= mg
		whiteCentipawnsEG -= eg
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
