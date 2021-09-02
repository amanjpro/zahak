package engine

// Piece Square Tables
// Middle-game
var EarlyPawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	80, 54, 89, 91, 84, 31, -20, -59,
	-13, -11, 20, 15, 26, 65, 20, -6,
	-30, -20, -11, -7, 12, 9, -8, -19,
	-35, -29, -17, 0, 2, -3, -14, -24,
	-38, -32, -22, -13, -8, -19, -11, -27,
	-40, -35, -32, -33, -23, -19, -7, -39,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var EarlyKnightPst = [64]int16{
	-178, -70, -32, -57, 27, -69, -11, -112,
	-11, 4, 28, 56, 19, 99, 10, 5,
	19, 51, 45, 50, 91, 104, 60, 34,
	34, 41, 55, 62, 35, 69, 32, 54,
	21, 33, 44, 39, 50, 38, 42, 28,
	-1, 24, 34, 41, 52, 38, 45, 18,
	2, 7, 24, 41, 37, 39, 35, 30,
	-36, 5, 2, 17, 17, 33, 9, -18,
}

var EarlyBishopPst = [64]int16{
	-28, -43, -20, -66, -62, -25, -4, -48,
	3, 16, 12, 11, 12, 15, 0, 22,
	39, 43, 41, 48, 37, 69, 45, 52,
	24, 36, 40, 51, 55, 37, 35, 10,
	27, 23, 34, 55, 49, 35, 32, 35,
	29, 47, 50, 45, 47, 54, 52, 43,
	43, 50, 57, 37, 48, 58, 68, 52,
	36, 51, 32, 35, 35, 27, 34, 44,
}

var EarlyRookPst = [64]int16{
	-16, -26, -9, -11, -15, 0, 13, 16,
	-23, -30, -8, 6, -8, 35, 19, 48,
	-36, -10, -19, -18, 10, 31, 60, 18,
	-33, -30, -20, -17, -26, -6, -4, -17,
	-42, -46, -42, -31, -36, -40, -7, -21,
	-43, -37, -37, -31, -28, -22, 5, -11,
	-34, -38, -29, -29, -21, -11, -1, -22,
	-23, -22, -20, -14, -11, -11, 3, -21,
}

var EarlyQueenPst = [64]int16{
	-32, -37, -6, 8, 15, 49, 39, -11,
	-11, -27, -30, -21, -41, 6, -24, 35,
	9, 10, 0, -3, -5, 19, 26, 42,
	-10, 4, -13, -9, -19, -7, -7, 4,
	2, -12, -4, 0, -1, -5, 1, 11,
	-7, 13, 6, 5, 14, 8, 23, 13,
	8, 17, 21, 24, 25, 29, 31, 42,
	22, 7, 13, 22, 19, 11, 16, 19,
}

var EarlyKingPst = [64]int16{
	-50, 179, 160, 77, -4, -21, 51, 44,
	66, 81, 70, 46, 12, 29, 11, -54,
	-41, 51, 42, -15, 18, 38, 50, -51,
	-76, -22, -40, -66, -45, -65, -90, -116,
	-78, -65, -30, -64, -72, -72, -96, -145,
	-26, -8, -31, -44, -41, -33, -16, -43,
	29, 6, -8, -4, -4, -8, 15, 10,
	14, 37, 16, -37, -18, -37, 20, 29,
}

// Endgame
var LatePawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	145, 135, 112, 74, 68, 108, 127, 140,
	48, 38, 5, -25, -31, -5, 21, 30,
	27, 7, -8, -21, -25, -15, 0, 6,
	10, -2, -13, -17, -18, -12, -8, -5,
	5, -5, -13, -15, -9, -12, -14, -10,
	6, -6, -8, -11, -5, -12, -21, -8,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var LateKnightPst = [64]int16{
	-45, -17, -13, -2, -23, -37, -44, -87,
	-37, -13, -15, -21, -25, -43, -22, -43,
	-25, -12, 0, -9, -35, -24, -15, -25,
	-7, -1, 2, -8, -2, 8, 5, -20,
	-15, -3, -2, 12, 8, -2, -4, -18,
	-15, -8, -15, 5, -1, -16, -18, -21,
	-12, -14, -6, -11, -13, -16, -11, -5,
	2, -39, -9, -14, -11, -26, -34, -27,
}

var LateBishopPst = [64]int16{
	-14, 7, -7, -4, -3, -22, -28, -16,
	-31, -26, -18, -15, -22, -25, -26, -41,
	-20, -17, -26, -30, -28, -25, -16, -25,
	-12, -9, -17, -13, -22, -16, -19, -7,
	-17, -16, -14, -17, -24, -22, -22, -34,
	-30, -10, -13, -14, -9, -23, -26, -28,
	-25, -29, -32, -21, -22, -31, -27, -52,
	-31, -33, -37, -23, -16, -22, -36, -52,
}

var LateRookPst = [64]int16{
	14, 19, 19, 22, 16, 23, 8, 10,
	19, 33, 34, 25, 27, 13, 15, -2,
	19, 15, 21, 12, -3, 1, -11, -12,
	19, 21, 17, 11, 7, 4, 4, 4,
	17, 11, 15, 14, 14, 17, -5, -1,
	16, 5, 14, 10, 7, 1, -18, -16,
	10, 10, 12, 12, 8, -2, -9, -3,
	12, 2, 7, -2, -4, 2, -11, 5,
}

var LateQueenPst = [64]int16{
	54, 66, 68, 56, 66, 56, 54, 49,
	29, 73, 100, 89, 127, 102, 104, 62,
	26, 40, 74, 78, 103, 105, 98, 43,
	58, 57, 74, 90, 100, 83, 107, 68,
	37, 69, 56, 75, 67, 72, 73, 57,
	42, 18, 55, 47, 45, 58, 39, 28,
	25, 16, 16, 27, 25, 17, 5, -26,
	11, 16, 24, 10, 18, 16, 17, 12,
}

var LateKingPst = [64]int16{
	-146, -81, -41, -10, 29, 17, -26, -87,
	-12, 22, 42, 30, 44, 46, 48, 14,
	5, 28, 37, 46, 38, 34, 35, 12,
	0, 19, 30, 40, 36, 34, 28, 10,
	-2, 8, 17, 29, 23, 13, 10, 7,
	-21, -10, 5, 14, 13, 1, -11, -9,
	-28, -13, -2, -3, -1, -7, -23, -30,
	-46, -51, -30, -3, -21, -7, -40, -61,
}

var MiddlegamePassedPawnBonus = [6]int16{
	0, 0, 0, 20, 30, 15,
}

var EndgamePassedPawnBonus = [6]int16{
	0, 3, 28, 54, 91, 37,
}

var MiddlegameBackwardPawnPenalty int16 = 12
var EndgameBackwardPawnPenalty int16 = 7
var MiddlegameIsolatedPawnPenalty int16 = 15
var EndgameIsolatedPawnPenalty int16 = 5
var MiddlegameDoublePawnPenalty int16 = 3
var EndgameDoublePawnPenalty int16 = 25
var MiddlegameCandidatePassedPawnAward int16 = 59
var EndgameCandidatePassedPawnAward int16 = 46
var MiddlegameRookOpenFileAward int16 = 38
var EndgameRookOpenFileAward int16 = 2
var MiddlegameRookSemiOpenFileAward int16 = 14
var EndgameRookSemiOpenFileAward int16 = 9
var MiddlegameVeritcalDoubleRookAward int16 = 0
var EndgameVeritcalDoubleRookAward int16 = 18
var MiddlegameHorizontalDoubleRookAward int16 = 21
var EndgameHorizontalDoubleRookAward int16 = 10
var MiddlegamePawnFactorCoeff int16 = 1
var EndgamePawnFactorCoeff int16 = 1
var MiddlegamePawnSquareControlCoeff int16 = 0
var EndgamePawnSquareControlCoeff int16 = 10
var MiddlegameMinorMobilityFactorCoeff int16 = 3
var EndgameMinorMobilityFactorCoeff int16 = 4
var MiddlegameMinorAggressivityFactorCoeff int16 = 5
var EndgameMinorAggressivityFactorCoeff int16 = 7
var MiddlegameMajorMobilityFactorCoeff int16 = 3
var EndgameMajorMobilityFactorCoeff int16 = 4
var MiddlegameMajorAggressivityFactorCoeff int16 = 0
var EndgameMajorAggressivityFactorCoeff int16 = 5
var MiddlegameInnerPawnToKingAttackCoeff int16 = 7
var EndgameInnerPawnToKingAttackCoeff int16 = 0
var MiddlegameOuterPawnToKingAttackCoeff int16 = 2
var EndgameOuterPawnToKingAttackCoeff int16 = 8
var MiddlegameInnerMinorToKingAttackCoeff int16 = 20
var EndgameInnerMinorToKingAttackCoeff int16 = 0
var MiddlegameOuterMinorToKingAttackCoeff int16 = 11
var EndgameOuterMinorToKingAttackCoeff int16 = 2
var MiddlegameInnerMajorToKingAttackCoeff int16 = 18
var EndgameInnerMajorToKingAttackCoeff int16 = 0
var MiddlegameOuterMajorToKingAttackCoeff int16 = 5
var EndgameOuterMajorToKingAttackCoeff int16 = 3
var MiddlegamePawnShieldPenalty int16 = 14
var EndgamePawnShieldPenalty int16 = 7
var MiddlegameNotCastlingPenalty int16 = 59
var EndgameNotCastlingPenalty int16 = 0
var MiddlegameKingZoneOpenFilePenalty int16 = 54
var EndgameKingZoneOpenFilePenalty int16 = 0
var MiddlegameKingZoneMissingPawnPenalty int16 = 23
var EndgameKingZoneMissingPawnPenalty int16 = 0
var MiddlegameKnightOutpostAward int16 = 16
var EndgameKnightOutpostAward int16 = 32
var MiddlegameBishopPairAward int16 = 21
var EndgameBishopPairAward int16 = 53
var MiddlegameKingVirtualMobilityPenalty int16 = 3
var EndgameKingVirtualMobilityPenalty int16 = 0

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
