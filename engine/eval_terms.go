package engine

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

var Flip = [64]int16{
	56, 57, 58, 59, 60, 61, 62, 63,
	48, 49, 50, 51, 52, 53, 54, 55,
	40, 41, 42, 43, 44, 45, 46, 47,
	32, 33, 34, 35, 36, 37, 38, 39,
	24, 25, 26, 27, 28, 29, 30, 31,
	16, 17, 18, 19, 20, 21, 22, 23,
	8, 9, 10, 11, 12, 13, 14, 15,
	0, 1, 2, 3, 4, 5, 6, 7,
}

var EarlyPieceSquareTables [12][64]int16
var LatePieceSquareTables [12][64]int16

func init() {
	for j := 0; j < 64; j++ {
		// White pawn
		EarlyPieceSquareTables[WhitePawn-1][j] = EarlyPawnPst[Flip[j]]
		LatePieceSquareTables[WhitePawn-1][j] = LatePawnPst[Flip[j]]
		// White knight
		EarlyPieceSquareTables[WhiteKnight-1][j] = EarlyKnightPst[Flip[j]]
		LatePieceSquareTables[WhiteKnight-1][j] = LateKnightPst[Flip[j]]
		// White bishop
		EarlyPieceSquareTables[WhiteBishop-1][j] = EarlyBishopPst[Flip[j]]
		LatePieceSquareTables[WhiteBishop-1][j] = LateBishopPst[Flip[j]]
		// White rook
		EarlyPieceSquareTables[WhiteRook-1][j] = EarlyRookPst[Flip[j]]
		LatePieceSquareTables[WhiteRook-1][j] = LateRookPst[Flip[j]]
		// White queen
		EarlyPieceSquareTables[WhiteQueen-1][j] = EarlyQueenPst[Flip[j]]
		LatePieceSquareTables[WhiteQueen-1][j] = LateQueenPst[Flip[j]]
		// White king
		EarlyPieceSquareTables[WhiteKing-1][j] = EarlyKingPst[Flip[j]]
		LatePieceSquareTables[WhiteKing-1][j] = LateKingPst[Flip[j]]

		// Black pawn
		EarlyPieceSquareTables[BlackPawn-1][j] = EarlyPawnPst[j]
		LatePieceSquareTables[BlackPawn-1][j] = LatePawnPst[j]
		// Black knight
		EarlyPieceSquareTables[BlackKnight-1][j] = EarlyKnightPst[j]
		LatePieceSquareTables[BlackKnight-1][j] = LateKnightPst[j]
		// Black bishop
		EarlyPieceSquareTables[BlackBishop-1][j] = EarlyBishopPst[j]
		LatePieceSquareTables[BlackBishop-1][j] = LateBishopPst[j]
		// Black rook
		EarlyPieceSquareTables[BlackRook-1][j] = EarlyRookPst[j]
		LatePieceSquareTables[BlackRook-1][j] = LateRookPst[j]
		// Black queen
		EarlyPieceSquareTables[BlackQueen-1][j] = EarlyQueenPst[j]
		LatePieceSquareTables[BlackQueen-1][j] = LateQueenPst[j]
		// Black king
		EarlyPieceSquareTables[BlackKing-1][j] = EarlyKingPst[j]
		LatePieceSquareTables[BlackKing-1][j] = LateKingPst[j]
	}
}
