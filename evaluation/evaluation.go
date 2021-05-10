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

// Middle game
var EarlyPawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	50, 50, 50, 50, 50, 50, 50, 50,
	10, 10, 20, 30, 30, 20, 10, 10,
	5, 5, 10, 25, 25, 10, 5, 5,
	0, 0, 0, 20, 20, 0, 0, 0,
	5, -5, -10, 0, 0, -10, -5, 5,
	5, 10, 10, -20, -20, 10, 10, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var EarlyKnightPst = [64]int16{
	-40, -25, -25, -25, -25, -25, -25, -40,
	-30, 0, 0, 0, 0, 0, 0, -30,
	-30, 0, 0, 0, 0, 0, 0, -30,
	-30, 0, 0, 15, 15, 0, 0, -30,
	-30, 0, 0, 15, 15, 0, 0, -30,
	-30, 0, 10, 0, 0, 10, 0, -30,
	-30, 0, 0, 5, 5, 0, 0, -30,
	-40, -30, -25, -25, -25, -25, -30, -40,
}

var EarlyBishopPst = [64]int16{
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-10, 0, 5, 0, 0, 5, 0, -10,
	-10, 0, 0, 10, 10, 0, 0, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-10, -20, -20, -20, -20, -20, -20, -10,
}

var EarlyRookPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	10, 10, 10, 10, 10, 10, 10, 10,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 5, 5, 0, 0, 0,
}

var EarlyQueenPst = [64]int16{
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	5, 5, -25, -25, -25, -25, 5, 5,
	5, 5, 10, 15, 15, 10, 5, 5,
}

var EarlyKingPst = [64]int16{
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	20, 25, 20, -15, -15, -15, 25, 20,
}

// Endgame

var LatePawnPst = [64]int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	80, 80, 80, 80, 80, 80, 80, 80,
	60, 60, 60, 60, 60, 60, 60, 60,
	40, 40, 40, 40, 40, 40, 40, 40,
	10, 10, 10, 10, 10, 10, 10, 10,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var LateKnightPst = [64]int16{
	-30, -20, -10, -20, -20, -20, -30, -30,
	-10, -10, -10, -5, -5, -10, -10, -10,
	-10, -10, 10, 10, -10, -10, -10, -10,
	-10, 5, 10, 10, 10, 10, 10, -10,
	-10, -5, 10, 15, 10, 15, 5, -10,
	-10, -5, 0, 10, 10, 0, -10, -10,
	-25, -20, -10, -5, -5, -20, -20, -25,
	-30, -30, -30, -10, -10, -30, -30, -30,
}

var LateBishopPst = [64]int16{
	-10, -10, -10, -10, -10, -10, -10, -10,
	-10, -5, 5, -10, -5, -10, -5, -10,
	5, -10, 0, 0, 0, 5, 0, 5,
	-5, 10, 10, 10, 10, 10, 5, 0,
	-5, 5, 10, 15, 5, 10, -5, -10,
	-10, -5, 10, 10, 15, 5, -5, -10,
	-10, -15, -5, 0, 5, -5, -10, -15,
	-15, -10, -15, -5, -10, -10, -5, -15,
}

var LateRookPst = [64]int16{
	15, 10, 15, 15, 15, 15, 10, 5,
	10, 10, 10, 10, 5, 5, 10, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	0, 0, 0, 0, 0, 0, 0, 0,
	-5, -5, -5, -5, -5, -5, -5, -5,
	-10, -10, -10, -10, -10, -10, -10, -10,
}

var LateQueenPst = [64]int16{
	-10, 20, 20, 25, 25, 20, 10, 20,
	-15, 20, 30, 40, 40, 20, 20, 0,
	-20, 5, 10, 30, 30, 30, 5, -20,
	5, 20, 20, 30, 30, 20, 20, 5,
	-15, 25, 20, 30, 30, 20, 25, -15,
	-15, -25, 10, 5, 10, 15, 10, 5,
	-20, -20, -30, -15, -15, -20, -20, -20,
	-30, -30, -20, -30, -5, -20, -20, -20,
}

var LateKingPst = [64]int16{
	-50, -50, -50, -50, -50, -50, -50, -50,
	-15, 15, 15, 15, 15, 15, 15, -15,
	10, 15, 20, 15, 20, 20, 15, 10,
	-10, 20, 20, 20, 20, 20, 20, -10,
	-15, -5, 20, 20, 20, 20, -5, -15,
	-15, -5, 10, 20, 20, 10, -5, -15,
	-20, -10, 5, 15, 15, 5, -10, -20,
	-40, -40, -20, -10, -10, -20, -40, -40,
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

var MiddlegameBackwardPawnAward = int16(25)
var EndgameBackwardPawnAward = int16(25)

var MiddlegameIsolatedPawnAward = int16(25)
var EndgameIsolatedPawnAward = int16(25)

var MiddlegameDoublePawnAward = int16(25)
var EndgameDoublePawnAward = int16(25)

var MiddlegamePassedPawnAward = int16(20)
var EndgamePassedPawnAward = int16(50)

var MiddlegameCandidatePassedPawnAward = int16(10)
var EndgameCandidatePassedPawnAward = int16(25)

var MiddlegameRookOpenFileAward = int16(25)
var EndgameRookOpenFileAward = int16(25)

var MiddlegameRookSemiOpenFileAward = int16(15)
var EndgameRookSemiOpenFileAward = int16(15)

var MiddlegameVeritcalDoubleRookAward = int16(50)
var EndgameVeritcalDoubleRookAward = int16(50)

var MiddlegameHorizontalDoubleRookAward = int16(30)
var EndgameHorizontalDoubleRookAward = int16(30)

var MiddlegamePawnFactorCoeff = int16(2)
var EndgamePawnFactorCoeff = int16(2)

var EndgameAggressivityFactorCoeff = int16(1)
var MiddlegameAggressivityFactorCoeff = int16(2)

var MiddlegameCastlingAward = int16(10)

func PSQT(piece Piece, sq Square, isEndgame bool) int16 {
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
		if isEndgame {
			return LateKingPst[flip[int(sq)]]
		}
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
		if isEndgame {
			return LateKingPst[int(sq)]
		}
		return EarlyKingPst[int(sq)]
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
			blackCentipawnsMG -= MiddlegameBackwardPawnAward
			blackCentipawnsEG -= EndgameBackwardPawnAward
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
			whiteCentipawnsMG -= MiddlegameBackwardPawnAward
			whiteCentipawnsEG -= EndgameBackwardPawnAward
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
		// isoLated pawn penalty - white
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
				whiteCentipawnsMG -= MiddlegameIsolatedPawnAward
				whiteCentipawnsEG -= EndgameIsolatedPawnAward
			}
		}

		// isoLated pawn penalty - black
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
				blackCentipawnsMG -= MiddlegameIsolatedPawnAward
				blackCentipawnsEG -= EndgameIsolatedPawnAward
			}
		}

		// double pawn penalty - black
		if blackPawnsPerFile[i] > 1 {
			blackCentipawnsMG -= MiddlegameDoublePawnAward
			blackCentipawnsEG -= EndgameDoublePawnAward
		}
		// double pawn penalty - white
		if whitePawnsPerFile[i] > 1 {
			whiteCentipawnsMG -= MiddlegameDoublePawnAward
			whiteCentipawnsEG -= EndgameDoublePawnAward
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
		if award <= 0 {
			if !position.HasTag(BlackCanCastleKingSide) {
				award -= MiddlegameCastlingAward
			} else if !position.HasTag(BlackCanCastleQueenSide) {
				award -= MiddlegameCastlingAward
			}
		}
		blackCentipawnsMG += award

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
		if award <= 0 {
			if !position.HasTag(WhiteCanCastleKingSide) {
				award -= MiddlegameCastlingAward
			} else if !position.HasTag(WhiteCanCastleQueenSide) {
				award -= MiddlegameCastlingAward
			}
		}
		whiteCentipawnsMG += award

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
	whiteAttacks := board.AllAttacksOn(Black) // get the squares that are taboo for black (white's reach)
	blackAttacks := board.AllAttacksOn(White) // get the squares that are taboo for whtie (black's reach)
	wAttackCounts := bits.OnesCount64(whiteAttacks)
	bAttackCounts := bits.OnesCount64(blackAttacks)

	whiteAggressivity := bits.OnesCount64(whiteAttacks >> 32) // keep hi-bits only (black's half)
	blackAggressivity := bits.OnesCount64(blackAttacks << 32) // keep lo-bits only (white's half)

	whiteCentipawnsMG += MiddlegameAggressivityFactorCoeff * int16(wAttackCounts-bAttackCounts)
	whiteCentipawnsEG += EndgameAggressivityFactorCoeff * int16(wAttackCounts-bAttackCounts)

	blackCentipawnsMG += MiddlegameAggressivityFactorCoeff * int16(bAttackCounts-wAttackCounts)
	blackCentipawnsEG += EndgameAggressivityFactorCoeff * int16(bAttackCounts-wAttackCounts)

	whiteCentipawnsMG += MiddlegameAggressivityFactorCoeff * int16(2*(whiteAggressivity-blackAggressivity))
	whiteCentipawnsEG += EndgameAggressivityFactorCoeff * int16(2*(whiteAggressivity-blackAggressivity))

	blackCentipawnsMG += MiddlegameAggressivityFactorCoeff * int16(2*(blackAggressivity-whiteAggressivity))
	blackCentipawnsEG += EndgameAggressivityFactorCoeff * int16(2*(blackAggressivity-whiteAggressivity))

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
