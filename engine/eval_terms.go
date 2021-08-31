package engine

// Piece Square Tables
// Middle-game
var EarlyPawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	67, 69, 81, 83, 66, 66, -19, -59,
	-11, -7, 27, 27, 40, 72, 26, -1,
	-28, -17, -9, -2, 18, 15, -1, -14,
	-35, -28, -15, 2, 5, 1, -12, -21,
	-37, -31, -20, -12, -2, -16, -5, -25,
	-39, -35, -30, -28, -18, -9, -1, -35,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var EarlyKnightPst = [64]int16{
	-174, -71, -32, -49, 34, -73, -1, -112,
	-21, 4, 35, 56, 19, 105, 16, 9,
	13, 55, 49, 53, 95, 116, 62, 43,
	34, 45, 57, 67, 39, 74, 37, 57,
	26, 35, 47, 42, 54, 42, 45, 32,
	2, 29, 39, 47, 57, 44, 49, 21,
	5, 10, 27, 45, 42, 42, 35, 34,
	-44, 10, 5, 19, 20, 36, 13, -13,
}

var EarlyBishopPst = [64]int16{
	-23, -20, -43, -57, -51, -26, -1, -36,
	2, 20, 14, 12, 19, 21, 7, 20,
	38, 45, 41, 50, 39, 72, 47, 53,
	25, 40, 40, 52, 59, 41, 37, 11,
	30, 26, 36, 55, 51, 38, 33, 38,
	31, 50, 52, 49, 48, 58, 54, 44,
	45, 52, 60, 39, 51, 59, 69, 53,
	33, 53, 35, 36, 38, 29, 36, 42,
}

var EarlyRookPst = [64]int16{
	-16, -20, -7, -4, -2, -8, 10, 9,
	-21, -29, -5, 8, 0, 36, 14, 41,
	-35, -11, -20, -19, 7, 30, 57, 11,
	-35, -27, -21, -14, -26, -5, -9, -18,
	-41, -43, -41, -32, -35, -40, -6, -22,
	-42, -37, -36, -31, -28, -19, 4, -13,
	-34, -37, -30, -28, -21, -11, -2, -28,
	-22, -21, -19, -11, -9, -8, -3, -25,
}

var EarlyQueenPst = [64]int16{
	-40, -41, -9, 0, 14, 43, 38, -13,
	-15, -34, -33, -27, -49, 8, -22, 32,
	3, 3, -2, -14, -12, 20, 19, 34,
	-16, -3, -17, -21, -26, -14, -16, -5,
	1, -20, -10, -7, -10, -10, -6, 4,
	-10, 11, 2, 2, 8, 4, 17, 5,
	4, 14, 19, 23, 23, 26, 27, 37,
	19, 5, 12, 21, 15, 7, 13, 11,
}

var EarlyKingPst = [64]int16{
	-70, 140, 133, 78, -17, -14, 44, 17,
	97, 72, 64, 71, 38, 26, 9, -64,
	-7, 55, 59, 22, 42, 50, 56, -36,
	-48, -11, -21, -44, -37, -41, -62, -99,
	-71, -39, -36, -70, -80, -69, -92, -124,
	-25, -12, -35, -53, -47, -43, -23, -46,
	37, 11, -6, -16, -14, -10, 14, 16,
	20, 46, 17, -47, -24, -38, 21, 33,
}

// Endgame
var LatePawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	131, 112, 99, 58, 58, 75, 110, 123,
	73, 64, 34, 4, -3, 12, 45, 50,
	26, 7, -6, -22, -25, -16, -1, 6,
	16, 2, -7, -8, -11, -8, -2, 0,
	7, -2, -9, -11, -6, -10, -13, -8,
	8, -3, -5, -9, 1, -11, -18, -8,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var LateKnightPst = [64]int16{
	-45, -15, -13, -6, -25, -34, -43, -83,
	-32, -11, -22, -21, -24, -46, -25, -46,
	-20, -12, 0, -11, -35, -29, -14, -29,
	-4, 0, 5, -7, 0, 7, 6, -21,
	-17, -3, -1, 11, 8, -2, -5, -18,
	-15, -10, -16, 6, -1, -17, -19, -20,
	-13, -14, -5, -9, -13, -15, -11, -10,
	14, -41, -6, -10, -12, -25, -37, -35,
}

var LateBishopPst = [64]int16{
	-16, 1, 0, -7, -7, -17, -30, -23,
	-31, -27, -19, -13, -29, -30, -29, -30,
	-16, -18, -25, -31, -29, -23, -14, -17,
	-9, -9, -15, -10, -21, -15, -16, -4,
	-19, -17, -14, -15, -23, -21, -20, -35,
	-30, -11, -11, -14, -6, -22, -26, -25,
	-25, -29, -32, -19, -20, -29, -26, -49,
	-23, -28, -36, -19, -15, -20, -33, -42,
}

var LateRookPst = [64]int16{
	9, 13, 14, 14, 5, 22, 5, 8,
	13, 29, 28, 20, 18, 10, 14, -4,
	14, 12, 18, 8, -4, -2, -12, -12,
	18, 16, 15, 7, 3, 0, 5, 3,
	16, 9, 12, 13, 10, 16, -9, -4,
	15, 5, 13, 10, 6, -1, -18, -11,
	11, 9, 13, 11, 7, -4, -9, 2,
	12, 3, 7, -5, -5, -1, -4, 4,
}

var LateQueenPst = [64]int16{
	50, 56, 58, 49, 51, 50, 47, 43,
	20, 68, 87, 80, 119, 77, 87, 56,
	22, 35, 56, 76, 89, 87, 89, 48,
	57, 58, 64, 92, 92, 74, 104, 71,
	25, 71, 48, 73, 65, 65, 68, 58,
	40, 9, 49, 35, 41, 55, 44, 35,
	18, 9, 5, 19, 18, 9, -2, -18,
	2, 9, 16, -2, 15, 12, 15, 17,
}

var LateKingPst = [64]int16{
	-141, -79, -41, -15, 31, 12, -29, -74,
	-17, 20, 38, 25, 41, 47, 46, 14,
	-4, 24, 31, 38, 33, 32, 33, 7,
	-6, 17, 27, 38, 36, 32, 22, 7,
	-6, 1, 18, 30, 26, 14, 9, 2,
	-23, -11, 4, 15, 14, 2, -11, -8,
	-26, -13, -2, -2, -1, -5, -20, -25,
	-38, -49, -31, -6, -28, -7, -38, -56,
}

var MiddlegameBackwardPawnPenalty int16 = 12
var EndgameBackwardPawnPenalty int16 = 7
var MiddlegameIsolatedPawnPenalty int16 = 14
var EndgameIsolatedPawnPenalty int16 = 6
var MiddlegameDoublePawnPenalty int16 = 7
var EndgameDoublePawnPenalty int16 = 23
var MiddlegamePassedPawnAward int16 = 0
var EndgamePassedPawnAward int16 = 13
var MiddlegameAdvancedPassedPawnAward int16 = 17
var EndgameAdvancedPassedPawnAward int16 = 62
var MiddlegameCandidatePassedPawnAward int16 = 53
var EndgameCandidatePassedPawnAward int16 = 52
var MiddlegameRookOpenFileAward int16 = 37
var EndgameRookOpenFileAward int16 = 2
var MiddlegameRookSemiOpenFileAward int16 = 15
var EndgameRookSemiOpenFileAward int16 = 10
var MiddlegameVeritcalDoubleRookAward int16 = 0
var EndgameVeritcalDoubleRookAward int16 = 19
var MiddlegameHorizontalDoubleRookAward int16 = 22
var EndgameHorizontalDoubleRookAward int16 = 7
var MiddlegamePawnFactorCoeff int16 = 1
var EndgamePawnFactorCoeff int16 = 1
var MiddlegamePawnSquareControlCoeff int16 = 1
var EndgamePawnSquareControlCoeff int16 = 7
var MiddlegameMinorMobilityFactorCoeff int16 = 3
var EndgameMinorMobilityFactorCoeff int16 = 3
var MiddlegameMinorAggressivityFactorCoeff int16 = 5
var EndgameMinorAggressivityFactorCoeff int16 = 7
var MiddlegameMajorMobilityFactorCoeff int16 = 3
var EndgameMajorMobilityFactorCoeff int16 = 4
var MiddlegameMajorAggressivityFactorCoeff int16 = 0
var EndgameMajorAggressivityFactorCoeff int16 = 6
var MiddlegameInnerPawnToKingAttackCoeff int16 = 6
var EndgameInnerPawnToKingAttackCoeff int16 = 0
var MiddlegameOuterPawnToKingAttackCoeff int16 = 1
var EndgameOuterPawnToKingAttackCoeff int16 = 7
var MiddlegameInnerMinorToKingAttackCoeff int16 = 20
var EndgameInnerMinorToKingAttackCoeff int16 = 0
var MiddlegameOuterMinorToKingAttackCoeff int16 = 11
var EndgameOuterMinorToKingAttackCoeff int16 = 3
var MiddlegameInnerMajorToKingAttackCoeff int16 = 19
var EndgameInnerMajorToKingAttackCoeff int16 = 0
var MiddlegameOuterMajorToKingAttackCoeff int16 = 6
var EndgameOuterMajorToKingAttackCoeff int16 = 4
var MiddlegamePawnShieldPenalty int16 = 16
var EndgamePawnShieldPenalty int16 = 10
var MiddlegameNotCastlingPenalty int16 = 55
var EndgameNotCastlingPenalty int16 = 0
var MiddlegameKingZoneOpenFilePenalty int16 = 60
var EndgameKingZoneOpenFilePenalty int16 = 0
var MiddlegameKingZoneMissingPawnPenalty int16 = 27
var EndgameKingZoneMissingPawnPenalty int16 = 0
var MiddlegameKnightOutpostAward int16 = 16
var EndgameKnightOutpostAward int16 = 30
var MiddlegameBishopPairAward int16 = 23
var EndgameBishopPairAward int16 = 50

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
	UpdatePSQTs()
}

func UpdatePSQTs() {
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
