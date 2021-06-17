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
const BlackQueenSideMask = uint64(1<<C8 | 1<<B8 | 1<<A8 | 1<<C7 | 1<<B7 | 1<<A7)
const WhiteQueenSideMask = uint64(1<<C1 | 1<<B1 | 1<<A1 | 1<<C2 | 1<<B2 | 1<<A2)

const BlackInnerAShield = uint64(1 << A7)
const BlackInnerBShield = uint64(1 << B7)
const BlackInnerCShield = uint64(1 << C7)
const BlackInnerFShield = uint64(1 << F7)
const BlackInnerGShield = uint64(1 << G7)
const BlackInnerHShield = uint64(1 << H7)

const BlackOuterAShield = uint64(1 << A6)
const BlackOuterBShield = uint64(1 << B6)
const BlackOuterCShield = uint64(1 << C6)
const BlackOuterFShield = uint64(1 << F6)
const BlackOuterGShield = uint64(1 << G6)
const BlackOuterHShield = uint64(1 << H6)

const WhiteInnerAShield = uint64(1 << A2)
const WhiteInnerBShield = uint64(1 << B2)
const WhiteInnerCShield = uint64(1 << C2)
const WhiteInnerFShield = uint64(1 << F2)
const WhiteInnerGShield = uint64(1 << G2)
const WhiteInnerHShield = uint64(1 << H2)

const WhiteOuterAShield = uint64(1 << A3)
const WhiteOuterBShield = uint64(1 << B3)
const WhiteOuterCShield = uint64(1 << C3)
const WhiteOuterFShield = uint64(1 << F3)
const WhiteOuterGShield = uint64(1 << G3)
const WhiteOuterHShield = uint64(1 << H3)

// Piece Square Tables
// Middle-game
var EarlyPawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	94, 130, 66, 100, 98, 123, 14, -27,
	-15, -13, 19, 18, 59, 69, 17, -20,
	-24, -6, -4, 14, 16, 11, 2, -30,
	-33, -23, -12, 5, 8, 4, -8, -33,
	-34, -29, -20, -22, -10, -9, 8, -23,
	-42, -24, -33, -34, -31, 11, 12, -32,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var EarlyKnightPst = [64]int16{
	-212, -87, -46, -48, 60, -111, -27, -137,
	-76, -41, 78, 40, 30, 73, 7, -14,
	-39, 70, 45, 64, 97, 139, 68, 51,
	8, 30, 26, 54, 24, 73, 15, 32,
	6, 32, 30, 23, 42, 26, 37, 8,
	-4, 11, 27, 33, 45, 32, 45, 7,
	-4, -30, 13, 21, 20, 36, 15, 13,
	-104, 5, -33, -15, 16, 2, 11, 0,
}

var EarlyBishopPst = [64]int16{
	-30, 23, -81, -42, -19, -35, 8, 8,
	3, 49, 14, 3, 55, 74, 37, -25,
	20, 68, 77, 62, 62, 74, 48, 22,
	30, 37, 41, 69, 59, 51, 32, 21,
	31, 51, 45, 60, 65, 47, 46, 39,
	31, 53, 54, 51, 54, 73, 57, 41,
	42, 62, 57, 43, 54, 62, 80, 44,
	3, 38, 33, 28, 35, 33, 0, 16,
}

var EarlyRookPst = [64]int16{
	0, 14, -11, 22, 24, -12, -1, -6,
	5, 3, 33, 30, 59, 62, -2, 24,
	-34, -7, 0, -6, -26, 34, 45, -14,
	-42, -29, -13, 0, -12, 20, -19, -37,
	-48, -44, -28, -26, -16, -27, 1, -37,
	-52, -32, -32, -30, -20, -11, -14, -32,
	-47, -26, -34, -24, -16, 1, -17, -69,
	-18, -20, -12, -2, 1, -3, -36, -14,
}

var EarlyQueenPst = [64]int16{
	-65, -27, -11, -18, 34, 32, 30, 9,
	-34, -60, -25, -24, -62, 30, -4, 22,
	-16, -23, -15, -38, -10, 27, 6, 22,
	-41, -39, -38, -48, -32, -22, -33, -28,
	-10, -41, -24, -24, -23, -18, -16, -13,
	-24, 5, -11, -3, -6, -2, 5, -2,
	-24, -2, 17, 8, 15, 24, 4, 15,
	15, -6, 6, 22, -5, -8, -13, -32,
}

var EarlyKingPst = [64]int16{
	-54, 92, 84, 41, -45, -15, 42, 37,
	86, 32, 24, 51, 18, -3, -14, -64,
	33, 36, 41, 7, 26, 52, 55, -13,
	-17, -12, 12, -18, -20, -21, -21, -54,
	-47, 9, -34, -74, -70, -54, -60, -73,
	-3, -17, -27, -57, -58, -47, -20, -36,
	11, 15, -13, -56, -38, -17, 9, 16,
	-9, 36, 12, -65, 1, -35, 25, 27,
}

// Endgame
var LatePawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	169, 152, 135, 107, 118, 107, 155, 192,
	88, 84, 63, 37, 20, 26, 60, 71,
	23, 6, -7, -25, -18, -12, 1, 10,
	20, 9, -1, -9, -9, -8, -2, 4,
	5, 1, -7, -5, -4, -7, -14, -10,
	16, -1, 7, 3, 9, -7, -11, -9,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var LateKnightPst = [64]int16{
	-46, -47, -17, -40, -45, -37, -75, -99,
	-32, -13, -46, -18, -27, -51, -37, -60,
	-35, -40, -11, -12, -31, -39, -39, -61,
	-25, -9, 5, 9, 14, -6, 0, -27,
	-25, -25, 2, 13, 3, 5, -8, -26,
	-29, -14, -18, -1, -7, -18, -34, -31,
	-42, -24, -18, -15, -12, -28, -30, -54,
	-22, -57, -28, -14, -28, -25, -58, -72,
}

var LateBishopPst = [64]int16{
	-8, -23, -6, -9, -6, -8, -16, -26,
	-11, -15, -1, -15, -14, -21, -15, -13,
	-5, -19, -16, -17, -17, -14, -8, 0,
	-7, 2, 2, -3, -1, -4, -6, 1,
	-11, -10, 3, 4, -8, 1, -13, -12,
	-10, -8, -1, 2, 6, -12, -11, -15,
	-15, -24, -12, -3, -1, -14, -22, -31,
	-18, -9, -18, -5, -10, -12, -3, -14,
}

var LateRookPst = [64]int16{
	10, 3, 13, 4, 6, 13, 7, 8,
	9, 12, 4, 4, -14, -6, 12, 3,
	12, 9, 1, 5, 5, -10, -7, 1,
	14, 8, 15, 1, 4, 4, 4, 14,
	15, 18, 17, 11, 5, 8, -3, 5,
	12, 11, 8, 11, 2, -1, 5, -2,
	9, 5, 12, 12, 0, -2, -4, 14,
	4, 11, 9, 0, -3, -1, 10, -12,
}

var LateQueenPst = [64]int16{
	32, 54, 49, 51, 42, 36, 26, 59,
	5, 45, 51, 64, 87, 40, 54, 37,
	-3, 22, 18, 81, 69, 50, 50, 31,
	42, 45, 44, 67, 77, 57, 92, 71,
	1, 54, 42, 66, 51, 54, 65, 48,
	26, -23, 29, 17, 28, 36, 44, 39,
	-10, -11, -23, -1, 5, -6, -17, -13,
	-25, -17, -7, -31, 20, -12, -3, -30,
}

var LateKingPst = [64]int16{
	-71, -53, -32, -29, -5, 16, -3, -13,
	-31, 1, 1, 0, 4, 30, 14, 20,
	3, 6, 7, 5, 4, 30, 27, 11,
	-9, 13, 15, 22, 20, 28, 20, 9,
	-15, -15, 20, 28, 29, 22, 8, -3,
	-21, -8, 10, 22, 25, 18, 1, -3,
	-26, -17, 8, 13, 16, 5, -11, -20,
	-52, -47, -22, 5, -26, -3, -34, -54,
}

var MiddlegameBackwardPawnPenalty int16 = 10
var EndgameBackwardPawnPenalty int16 = 2
var MiddlegameIsolatedPawnPenalty int16 = 12
var EndgameIsolatedPawnPenalty int16 = 6
var MiddlegameDoublePawnPenalty int16 = 1
var EndgameDoublePawnPenalty int16 = 26
var MiddlegamePassedPawnAward int16 = 0
var EndgamePassedPawnAward int16 = 9
var MiddlegameAdvancedPassedPawnAward int16 = 10
var EndgameAdvancedPassedPawnAward int16 = 58
var MiddlegameCandidatePassedPawnAward int16 = 32
var EndgameCandidatePassedPawnAward int16 = 50
var MiddlegameRookOpenFileAward int16 = 46
var EndgameRookOpenFileAward int16 = 0
var MiddlegameRookSemiOpenFileAward int16 = 13
var EndgameRookSemiOpenFileAward int16 = 20
var MiddlegameVeritcalDoubleRookAward int16 = 8
var EndgameVeritcalDoubleRookAward int16 = 12
var MiddlegameHorizontalDoubleRookAward int16 = 26
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
var MiddlegameInnerMinorToKingAttackCoeff int16 = 18
var EndgameInnerMinorToKingAttackCoeff int16 = 0
var MiddlegameOuterMinorToKingAttackCoeff int16 = 10
var EndgameOuterMinorToKingAttackCoeff int16 = 1
var MiddlegameInnerMajorToKingAttackCoeff int16 = 16
var EndgameInnerMajorToKingAttackCoeff int16 = 0
var MiddlegameOuterMajorToKingAttackCoeff int16 = 7
var EndgameOuterMajorToKingAttackCoeff int16 = 5
var MiddlegameInnerPawnShieldReward int16 = 8
var EndgameInnerPawnShieldReward int16 = 0
var MiddlegameOuterPawnShieldReward int16 = 8
var EndgameOuterPawnShieldReward int16 = 0
var MiddlegameNotCastlingPenalty int16 = 24
var EndgameNotCastlingPenalty int16 = 4
var MiddlegameKingZoneMissingOwnPawnPenalty int16 = 19
var EndgameKingZoneMissingOwnPawnPenalty int16 = 0
var MiddlegameKingZoneOpponentPawnMissingPenalty int16 = 18
var EndgameKingZoneOpponentPawnMissingPenalty int16 = 1

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
// var MiddlegameInnerPawnShieldReward int16 = 10
// var EndgameInnerPawnShieldReward int16 = 1
// var MiddlegameOuterPawnShieldReward int16 = 11
// var EndgameOuterPawnShieldReward int16 = 2
// var MiddlegameNotCastlingPenalty int16 = 21
// var EndgameNotCastlingPenalty int16 = 2
// var MiddlegameKingZoneMissingOwnPawnPenalty int16 = 24
// var EndgameKingZoneMissingOwnPawnPenalty int16 = 0
// var MiddlegameKingZoneOpponentPawnMissingPenalty int16 = 24
// var EndgameKingZoneOpponentPawnMissingPenalty int16 = 0

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

func kingSafetyPenalty(color Color, side PieceType, ownPawn uint64, opPawn uint64) (int16, int16) {
	var mg, eg int16
	var a_inner_shield, b_inner_shield, c_inner_shield, f_inner_shield, g_inner_shield, h_inner_shield,
		a_outer_shield, b_outer_shield, c_outer_shield, f_outer_shield, g_outer_shield, h_outer_shield uint64

	if color == White {
		a_inner_shield = WhiteInnerAShield
		b_inner_shield = WhiteInnerBShield
		c_inner_shield = WhiteInnerCShield
		f_inner_shield = WhiteInnerFShield
		g_inner_shield = WhiteInnerGShield
		h_inner_shield = WhiteInnerHShield

		a_outer_shield = WhiteOuterAShield
		b_outer_shield = WhiteOuterBShield
		c_outer_shield = WhiteOuterCShield
		f_outer_shield = WhiteOuterFShield
		g_outer_shield = WhiteOuterGShield
		h_outer_shield = WhiteOuterHShield

	} else {
		a_inner_shield = BlackInnerAShield
		b_inner_shield = BlackInnerBShield
		c_inner_shield = BlackInnerCShield
		f_inner_shield = BlackInnerFShield
		g_inner_shield = BlackInnerGShield
		h_inner_shield = BlackInnerHShield

		a_outer_shield = BlackOuterAShield
		b_outer_shield = BlackOuterBShield
		c_outer_shield = BlackOuterCShield
		f_outer_shield = BlackOuterFShield
		g_outer_shield = BlackOuterGShield
		h_outer_shield = BlackOuterHShield

	}
	if side == King {
		if H_FileFill&opPawn == 0 {
			mg -= MiddlegameKingZoneOpponentPawnMissingPenalty
			eg -= EndgameKingZoneOpponentPawnMissingPenalty
		}

		if H_FileFill&ownPawn == 0 {
			mg -= MiddlegameKingZoneMissingOwnPawnPenalty
			eg -= EndgameKingZoneMissingOwnPawnPenalty
		} else {
			if h_inner_shield&ownPawn != 0 {
				mg += MiddlegameInnerPawnShieldReward
				eg += EndgameInnerPawnShieldReward
			} else if h_outer_shield&ownPawn != 0 {
				mg += MiddlegameOuterPawnShieldReward
				eg += EndgameOuterPawnShieldReward
			}
		}

		if G_FileFill&opPawn == 0 {
			mg -= MiddlegameKingZoneOpponentPawnMissingPenalty
			eg -= EndgameKingZoneOpponentPawnMissingPenalty
		}

		if G_FileFill&ownPawn == 0 {
			mg -= MiddlegameKingZoneMissingOwnPawnPenalty
			eg -= EndgameKingZoneMissingOwnPawnPenalty
		} else {
			if g_inner_shield&ownPawn != 0 {
				mg += MiddlegameInnerPawnShieldReward
				eg += EndgameInnerPawnShieldReward
			} else if g_outer_shield&ownPawn != 0 {
				mg += MiddlegameOuterPawnShieldReward
				eg += EndgameOuterPawnShieldReward
			}
		}

		if F_FileFill&opPawn == 0 {
			mg -= MiddlegameKingZoneOpponentPawnMissingPenalty
			eg -= EndgameKingZoneOpponentPawnMissingPenalty
		}

		if F_FileFill&ownPawn == 0 {
			mg -= MiddlegameKingZoneMissingOwnPawnPenalty
			eg -= EndgameKingZoneMissingOwnPawnPenalty
		} else {

			if f_inner_shield&ownPawn != 0 {
				mg += MiddlegameInnerPawnShieldReward
				eg += EndgameInnerPawnShieldReward
			} else if f_outer_shield&ownPawn != 0 {
				mg += MiddlegameOuterPawnShieldReward
				eg += EndgameOuterPawnShieldReward
			}
		}
	} else {
		if C_FileFill&opPawn == 0 {
			mg -= MiddlegameKingZoneOpponentPawnMissingPenalty
			eg -= EndgameKingZoneOpponentPawnMissingPenalty
		}

		if C_FileFill&ownPawn == 0 {
			mg -= MiddlegameKingZoneMissingOwnPawnPenalty
			eg -= EndgameKingZoneMissingOwnPawnPenalty
		} else {
			if c_inner_shield&ownPawn != 0 {
				mg += MiddlegameInnerPawnShieldReward
				eg += EndgameInnerPawnShieldReward
			} else if c_outer_shield&ownPawn != 0 {
				mg += MiddlegameOuterPawnShieldReward
				eg += EndgameOuterPawnShieldReward
			}
		}

		if B_FileFill&opPawn == 0 {
			mg -= MiddlegameKingZoneOpponentPawnMissingPenalty
			eg -= EndgameKingZoneOpponentPawnMissingPenalty
		}

		if B_FileFill&ownPawn == 0 {
			mg -= MiddlegameKingZoneMissingOwnPawnPenalty
			eg -= EndgameKingZoneMissingOwnPawnPenalty
		} else {
			if b_inner_shield&ownPawn != 0 {
				mg += MiddlegameInnerPawnShieldReward
				eg += EndgameInnerPawnShieldReward
			} else if b_outer_shield&ownPawn != 0 {
				mg += MiddlegameOuterPawnShieldReward
				eg += EndgameOuterPawnShieldReward
			}
		}

		if A_FileFill&opPawn == 0 {
			mg -= MiddlegameKingZoneOpponentPawnMissingPenalty
			eg -= EndgameKingZoneOpponentPawnMissingPenalty
		}

		if A_FileFill&ownPawn == 0 {
			mg -= MiddlegameKingZoneMissingOwnPawnPenalty
			eg -= EndgameKingZoneMissingOwnPawnPenalty
		} else {
			if a_inner_shield&ownPawn != 0 {
				mg += MiddlegameInnerPawnShieldReward
				eg += EndgameInnerPawnShieldReward
			} else if a_outer_shield&ownPawn != 0 {
				mg += MiddlegameOuterPawnShieldReward
				eg += EndgameOuterPawnShieldReward
			}
		}
	}

	return mg, eg
}

func KingSafety(blackKing uint64, whiteKing uint64, blackPawn uint64,
	whitePawn uint64, blackCastle bool, whiteCastle bool) Eval {
	var whiteCentipawnsMG, whiteCentipawnsEG, blackCentipawnsMG, blackCentipawnsEG int16

	if blackKing&BlackKingSideMask != 0 {
		blackCastle = true
		// Missing pawn shield
		mg, eg := kingSafetyPenalty(Black, King, blackPawn, whitePawn)
		blackCentipawnsMG += mg
		blackCentipawnsEG += eg
	} else if blackKing&BlackQueenSideMask != 0 {
		blackCastle = true
		// Missing pawn shield
		mg, eg := kingSafetyPenalty(Black, Queen, blackPawn, whitePawn)
		blackCentipawnsMG += mg
		blackCentipawnsEG += eg
	}

	if whiteKing&WhiteKingSideMask != 0 {
		whiteCastle = true
		// Missing pawn shield
		mg, eg := kingSafetyPenalty(White, King, whitePawn, blackPawn)
		whiteCentipawnsMG += mg
		whiteCentipawnsEG += eg
	} else if whiteKing&WhiteQueenSideMask != 0 {
		whiteCastle = true
		// Missing pawn shield
		mg, eg := kingSafetyPenalty(White, Queen, whitePawn, blackPawn)
		whiteCentipawnsMG += mg
		whiteCentipawnsEG += eg
	}

	if !whiteCastle {
		whiteCentipawnsMG -= MiddlegameNotCastlingPenalty
		whiteCentipawnsEG -= EndgameNotCastlingPenalty
	}

	if !blackCastle {
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
