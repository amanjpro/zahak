package engine

// Piece Square Tables
// Middle-game
var EarlyPawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	77, 55, 89, 87, 81, 32, -20, -53,
	-17, -16, 16, 11, 24, 61, 17, -10,
	-32, -24, -13, -8, 12, 8, -7, -18,
	-36, -32, -17, 0, 3, -3, -14, -23,
	-39, -35, -21, -15, -4, -20, -8, -28,
	-42, -38, -32, -30, -20, -15, -4, -39,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var EarlyKnightPst = [64]int16{
	-180, -74, -35, -57, 25, -70, -12, -112,
	-13, 4, 26, 57, 18, 104, 10, 8,
	19, 52, 45, 50, 93, 108, 61, 37,
	32, 40, 55, 63, 37, 68, 34, 53,
	23, 34, 45, 39, 51, 38, 42, 28,
	2, 26, 35, 44, 52, 40, 46, 19,
	2, 9, 26, 42, 38, 38, 33, 29,
	-35, 7, 4, 17, 17, 33, 10, -18,
}

var EarlyBishopPst = [64]int16{
	-29, -42, -22, -64, -59, -27, -3, -50,
	3, 15, 12, 13, 13, 15, 2, 25,
	38, 43, 40, 48, 35, 70, 46, 54,
	23, 37, 39, 50, 54, 39, 35, 10,
	27, 23, 34, 53, 49, 36, 32, 36,
	30, 47, 50, 46, 46, 54, 52, 42,
	42, 50, 57, 37, 49, 57, 67, 52,
	37, 52, 33, 35, 37, 27, 35, 46,
}

var EarlyRookPst = [64]int16{
	-17, -25, -11, -10, -13, -1, 11, 15,
	-24, -28, -8, 5, -8, 33, 20, 48,
	-37, -11, -20, -20, 8, 31, 58, 17,
	-32, -28, -21, -16, -25, -6, -3, -17,
	-42, -44, -41, -31, -35, -39, -5, -23,
	-43, -39, -35, -32, -28, -20, 6, -10,
	-33, -37, -30, -28, -21, -12, -2, -25,
	-23, -21, -19, -12, -9, -9, 2, -24,
}

var EarlyQueenPst = [64]int16{
	-36, -35, -6, 7, 13, 48, 39, -11,
	-12, -27, -31, -22, -43, 3, -24, 34,
	8, 9, -1, -6, -9, 19, 24, 41,
	-10, 1, -14, -14, -21, -9, -10, 1,
	2, -15, -7, -3, -4, -7, -1, 8,
	-9, 11, 4, 4, 10, 6, 21, 10,
	5, 16, 20, 23, 24, 27, 29, 40,
	22, 6, 14, 21, 18, 9, 13, 17,
}

var EarlyKingPst = [64]int16{
	-48, 175, 159, 77, -8, -22, 46, 41,
	68, 79, 66, 47, 14, 28, 12, -53,
	-38, 52, 44, -16, 16, 37, 49, -50,
	-75, -21, -38, -67, -48, -67, -90, -113,
	-76, -63, -31, -66, -76, -74, -98, -141,
	-24, -8, -33, -48, -45, -38, -19, -42,
	34, 8, -8, -6, -6, -11, 14, 16,
	20, 43, 15, -43, -24, -39, 21, 33,
}

// Endgame
var LatePawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	144, 131, 110, 71, 66, 103, 124, 137,
	48, 40, 6, -26, -30, -4, 22, 33,
	26, 9, -7, -18, -23, -14, 1, 6,
	8, -3, -14, -17, -18, -13, -7, -6,
	5, -4, -15, -13, -9, -13, -15, -11,
	6, -6, -8, -11, -4, -11, -20, -9,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var LateKnightPst = [64]int16{
	-48, -18, -15, -5, -24, -40, -44, -90,
	-36, -13, -15, -24, -25, -47, -23, -44,
	-26, -12, 0, -10, -35, -26, -18, -27,
	-6, 0, 3, -8, -3, 8, 5, -22,
	-16, -4, -2, 11, 8, -2, -5, -18,
	-17, -7, -15, 4, 0, -16, -19, -21,
	-12, -16, -8, -10, -13, -15, -12, -5,
	0, -39, -9, -14, -11, -26, -35, -29,
}

var LateBishopPst = [64]int16{
	-16, 6, -6, -7, -4, -22, -30, -17,
	-32, -26, -19, -16, -25, -27, -27, -39,
	-20, -18, -26, -31, -29, -24, -17, -25,
	-11, -9, -17, -12, -21, -17, -19, -7,
	-18, -17, -15, -17, -23, -22, -24, -34,
	-32, -11, -12, -15, -9, -22, -26, -27,
	-27, -30, -32, -21, -22, -30, -27, -56,
	-33, -33, -36, -22, -17, -21, -38, -53,
}

var LateRookPst = [64]int16{
	13, 17, 19, 20, 13, 22, 7, 9,
	17, 31, 31, 24, 25, 12, 12, -5,
	18, 14, 21, 11, -4, 0, -11, -12,
	16, 18, 15, 9, 6, 2, 3, 4,
	17, 10, 13, 12, 11, 16, -8, -4,
	13, 4, 12, 9, 6, -1, -21, -18,
	6, 7, 10, 10, 7, -3, -9, -3,
	12, 2, 7, -4, -5, -1, -10, 4,
}

var LateQueenPst = [64]int16{
	53, 63, 65, 54, 63, 55, 52, 49,
	29, 72, 97, 88, 125, 99, 105, 62,
	26, 40, 72, 77, 99, 100, 100, 45,
	58, 59, 72, 88, 99, 82, 107, 68,
	32, 69, 55, 73, 65, 69, 72, 59,
	44, 17, 53, 40, 44, 56, 41, 29,
	24, 13, 11, 24, 23, 13, 1, -23,
	7, 13, 21, 5, 15, 16, 18, 12,
}

var LateKingPst = [64]int16{
	-147, -79, -43, -12, 27, 15, -27, -87,
	-12, 22, 39, 29, 41, 45, 46, 12,
	3, 27, 34, 43, 35, 32, 33, 10,
	-2, 17, 27, 38, 34, 32, 26, 9,
	-2, 7, 16, 27, 23, 13, 10, 6,
	-23, -12, 2, 14, 13, 0, -11, -8,
	-27, -12, -2, -2, -1, -6, -21, -26,
	-43, -51, -31, -6, -27, -8, -39, -57,
}

var MiddlegamePassedPawnBonus = [6]int16{
	3, 0, 0, 23, 34, 19,
}

var EndgamePassedPawnBonus = [6]int16{
	0, 2, 27, 49, 87, 34,
}

var MiddlegameBackwardPawnPenalty int16 = 12
var EndgameBackwardPawnPenalty int16 = 6
var MiddlegameIsolatedPawnPenalty int16 = 16
var EndgameIsolatedPawnPenalty int16 = 4
var MiddlegameDoublePawnPenalty int16 = 5
var EndgameDoublePawnPenalty int16 = 23
var MiddlegameCandidatePassedPawnAward int16 = 55
var EndgameCandidatePassedPawnAward int16 = 45
var MiddlegameRookOpenFileAward int16 = 37
var EndgameRookOpenFileAward int16 = 2
var MiddlegameRookSemiOpenFileAward int16 = 15
var EndgameRookSemiOpenFileAward int16 = 7
var MiddlegameVeritcalDoubleRookAward int16 = 0
var EndgameVeritcalDoubleRookAward int16 = 17
var MiddlegameHorizontalDoubleRookAward int16 = 20
var EndgameHorizontalDoubleRookAward int16 = 9
var MiddlegamePawnFactorCoeff int16 = 1
var EndgamePawnFactorCoeff int16 = 1
var MiddlegamePawnSquareControlCoeff int16 = 0
var EndgamePawnSquareControlCoeff int16 = 10
var MiddlegameMinorMobilityFactorCoeff int16 = 3
var EndgameMinorMobilityFactorCoeff int16 = 3
var MiddlegameMinorAggressivityFactorCoeff int16 = 5
var EndgameMinorAggressivityFactorCoeff int16 = 6
var MiddlegameMajorMobilityFactorCoeff int16 = 3
var EndgameMajorMobilityFactorCoeff int16 = 4
var MiddlegameMajorAggressivityFactorCoeff int16 = 0
var EndgameMajorAggressivityFactorCoeff int16 = 5
var MiddlegameInnerPawnToKingAttackCoeff int16 = 9
var EndgameInnerPawnToKingAttackCoeff int16 = 0
var MiddlegameOuterPawnToKingAttackCoeff int16 = 1
var EndgameOuterPawnToKingAttackCoeff int16 = 9
var MiddlegameInnerMinorToKingAttackCoeff int16 = 19
var EndgameInnerMinorToKingAttackCoeff int16 = 0
var MiddlegameOuterMinorToKingAttackCoeff int16 = 11
var EndgameOuterMinorToKingAttackCoeff int16 = 3
var MiddlegameInnerMajorToKingAttackCoeff int16 = 18
var EndgameInnerMajorToKingAttackCoeff int16 = 0
var MiddlegameOuterMajorToKingAttackCoeff int16 = 5
var EndgameOuterMajorToKingAttackCoeff int16 = 3
var MiddlegamePawnShieldPenalty int16 = 17
var EndgamePawnShieldPenalty int16 = 8
var MiddlegameNotCastlingPenalty int16 = 62
var EndgameNotCastlingPenalty int16 = 0
var MiddlegameKingZoneOpenFilePenalty int16 = 59
var EndgameKingZoneOpenFilePenalty int16 = 0
var MiddlegameKingZoneMissingPawnPenalty int16 = 28
var EndgameKingZoneMissingPawnPenalty int16 = 0
var MiddlegameKnightOutpostAward int16 = 15
var EndgameKnightOutpostAward int16 = 32
var MiddlegameBishopPairAward int16 = 20
var EndgameBishopPairAward int16 = 53

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
