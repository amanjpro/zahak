package evaluation

import (
	"math/bits"

	. "github.com/amanjpro/zahak/engine"
)

func Evaluate(position *Position) int16 {
	// board := position.Board
	// allPieces := board.AllPieces()
	return middlegameEval(position)
}

const CHECKMATE_EVAL int16 = 3100
const DIVIDER int16 = 800

// Piece Square Tables
var pawnPst = []int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 15, 15, 0, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 0, 5, 5, 0, 0, 0,
	0, 0, 0, -25, -25, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var knightPst = []int16{
	-40, -25, -25, -25, -25, -25, -25, -40,
	-30, 0, 0, 0, 0, 0, 0, -30,
	-30, 0, 0, 0, 0, 0, 0, -30,
	-30, 0, 0, 15, 15, 0, 0, -30,
	-30, 0, 0, 15, 15, 0, 0, -30,
	-30, 0, 10, 0, 0, 10, 0, -30,
	-30, 0, 0, 5, 5, 0, 0, -30,
	-40, -30, -25, -25, -25, -25, -30, -40,
}

var bishopPst = []int16{
	-10, 0, 0, 0, 0, 0, 0, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-10, 0, 5, 0, 0, 5, 0, -10,
	-10, 0, 0, 10, 10, 0, 0, -10,
	-10, 0, 5, 10, 10, 5, 0, -10,
	-10, 0, 5, 0, 0, 5, 0, -10,
	-10, 5, 0, 0, 0, 0, 5, -10,
	-10, -20, -20, -20, -20, -20, -20, -10,
}

var rookPst = []int16{
	0, 0, 0, 0, 0, 0, 0, 0,
	10, 10, 10, 10, 10, 10, 10, 10,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 5, 5, 0, 0, 0,
}

var queenPst = []int16{
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	5, 5, 10, 10, 10, 15, 5, 5,
	5, 5, 10, 15, 15, 10, 5, 5,
}

var kingPst = []int16{
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	-25, -25, -25, -25, -25, -25, -25, -25,
	10, 15, 10, -15, -15, 15, 15, 10,
}

var flip = []int16{
	56, 57, 58, 59, 60, 61, 62, 63,
	48, 49, 50, 51, 52, 53, 54, 55,
	40, 41, 42, 43, 44, 45, 46, 47,
	32, 33, 34, 35, 36, 37, 38, 39,
	24, 25, 26, 27, 28, 29, 30, 31,
	16, 17, 18, 19, 20, 21, 22, 23,
	8, 9, 10, 11, 12, 13, 14, 15,
	0, 1, 2, 3, 4, 5, 6, 7,
}

func middlegameEval(position *Position) int16 {
	board := position.Board
	p := BlackPawn
	n := BlackKnight
	b := BlackBishop
	r := BlackRook
	q := BlackQueen

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

	blackPawnsCount := int16(bits.OnesCount64(bbBlackPawn))
	blackKnightsCount := int16(bits.OnesCount64(bbBlackKnight))
	blackBishopsCount := int16(bits.OnesCount64(bbBlackBishop))
	blackRooksCount := int16(bits.OnesCount64(bbBlackRook))
	blackQueensCount := int16(bits.OnesCount64(bbBlackQueen))

	whitePawnsCount := int16(bits.OnesCount64(bbWhitePawn))
	whiteKnightsCount := int16(bits.OnesCount64(bbWhiteKnight))
	whiteBishopsCount := int16(bits.OnesCount64(bbWhiteBishop))
	whiteRooksCount := int16(bits.OnesCount64(bbWhiteRook))
	whiteQueensCount := int16(bits.OnesCount64(bbWhiteQueen))

	blackCentipawns := blackPawnsCount * p.Weight()
	blackCentipawns += blackKnightsCount * n.Weight()
	blackCentipawns += blackBishopsCount * b.Weight()
	blackCentipawns += blackRooksCount * r.Weight()
	blackCentipawns += blackQueensCount * q.Weight()

	whiteCentipawns := whitePawnsCount * p.Weight()
	whiteCentipawns += whiteKnightsCount * n.Weight()
	whiteCentipawns += whiteBishopsCount * b.Weight()
	whiteCentipawns += whiteRooksCount * r.Weight()
	whiteCentipawns += whiteQueensCount * q.Weight()

	// pawns := blackPawnsCount + whitePawnsCount
	// This is not correct
	// whiteCentipawns += whiteBishopsCount * b.Weight() * int16(1+(16-pawns)/64)
	// whiteCentipawns += whiteKnightsCount * n.Weight() * int16(1-(16-pawns)/64)
	// blackCentipawns += blackBishopsCount * b.Weight() * int16(1+(16-pawns)/64)
	// blackCentipawns += blackKnightsCount * n.Weight() * int16(1-(16-pawns)/64)
	//
	// 2 Bishops vs 2 Knights
	if whiteBishopsCount >= 2 && blackBishopsCount < 2 {
		whiteCentipawns += 25
	}
	if whiteBishopsCount < 2 && blackBishopsCount >= 2 {
		blackCentipawns += 25
	}

	// whites := board.GetWhitePieces()
	// blacks := board.GetBlackPieces()
	// all := whites | blacks

	// PST for black pieces
	pieceIter := bbBlackPawn
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		blackCentipawns += pawnPst[flip[index]]
		pieceIter ^= mask
	}

	pieceIter = bbBlackKnight
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		blackCentipawns += knightPst[flip[index]]
		pieceIter ^= mask
	}

	pieceIter = bbBlackBishop
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		blackCentipawns += bishopPst[flip[index]]
		pieceIter ^= mask
	}

	pieceIter = bbBlackRook
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		blackCentipawns += rookPst[flip[index]]
		pieceIter ^= mask
	}

	pieceIter = bbBlackQueen
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		blackCentipawns += queenPst[flip[index]]
		pieceIter ^= mask
	}

	pieceIter = bbBlackKing
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		blackCentipawns += kingPst[flip[index]]
		pieceIter ^= mask
	}

	// PST for white pieces
	pieceIter = bbWhitePawn
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		whiteCentipawns += pawnPst[index]
		pieceIter ^= mask
	}

	pieceIter = bbWhiteKnight
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		whiteCentipawns += knightPst[index]
		pieceIter ^= mask
	}

	pieceIter = bbWhiteBishop
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		whiteCentipawns += bishopPst[index]
		pieceIter ^= mask
	}

	pieceIter = bbWhiteRook
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		whiteCentipawns += rookPst[index]
		pieceIter ^= mask
	}

	pieceIter = bbWhiteQueen
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		whiteCentipawns += queenPst[index]
		pieceIter ^= mask
	}

	pieceIter = bbWhiteKing
	for pieceIter != 0 {
		index := bits.TrailingZeros64(pieceIter)
		mask := uint64(1 << index)
		whiteCentipawns += kingPst[index]
		pieceIter ^= mask
	}

	if position.Turn() == White {
		return whiteCentipawns - blackCentipawns
	} else {
		return blackCentipawns - whiteCentipawns
	}
}

func evaluate(position *Position, allPieces map[Square]Piece) int16 {

	whiteBishops := int16(0)
	whiteKnights := int16(0)
	blackBishops := int16(0)
	blackKnights := int16(0)
	blackPawns := int16(0)
	whitePawns := int16(0)

	blackCentipawn := int16(0)
	whiteCentipawn := int16(0)

	whitePawnsPerFile, blackPawnsPerFile := pawnsPerFile(allPieces)
	//whitePawnsPerRank, blackPawnsPerRank := pawnsPerRank(allPieces)

	for square, piece := range allPieces {
		file := square.File()
		rank := square.Rank()
		switch piece {
		case WhiteKing:
			// This doesn't consider endgame
			if position.HasTag(WhiteCanCastleKingSide) ||
				position.HasTag(WhiteCanCastleQueenSide) ||
				square == A1 || square == A2 ||
				square == B1 || square == B2 ||
				square == C1 || square == C2 ||
				square == G1 || square == G2 ||
				square == H1 || square == H2 {
				whiteCentipawn += 10
			}
		case BlackKing:
			// This doesn't consider endgame
			if position.HasTag(BlackCanCastleKingSide) ||
				position.HasTag(BlackCanCastleQueenSide) ||
				square == A7 || square == A8 ||
				square == B7 || square == B8 ||
				square == C7 || square == C8 ||
				square == G7 || square == G8 ||
				square == H7 || square == H8 {
				blackCentipawn += 10
			}
		case WhiteQueen:
			whiteCentipawn += piece.Weight()
			whiteCentipawn += 5 * piece.Weight() * int16(rank+1) / DIVIDER
		case BlackQueen:
			blackCentipawn += piece.Weight()
			blackCentipawn += 5 * piece.Weight() * (9 - int16(rank+1)) / DIVIDER
		case WhiteRook:
			white := whitePawnsPerFile[file]
			black := blackPawnsPerFile[file]
			bonus := int16(0)
			if white == 0 && black == 0 { // open file
				bonus = 1
			} else if white == 0 { // semi-open file
				bonus = 5
			}
			whiteCentipawn += piece.Weight() + bonus
			whiteCentipawn += 5 * piece.Weight() * int16(rank+1) / DIVIDER
		case BlackRook:
			white := whitePawnsPerFile[file]
			black := blackPawnsPerFile[file]
			bonus := int16(0)
			if white == 0 && black == 0 { // open file
				bonus = 1
			} else if black == 0 { // semi-open file
				bonus = 5
			}
			blackCentipawn += piece.Weight() + bonus
			blackCentipawn += 5 * piece.Weight() * (9 - int16(rank+1)) / DIVIDER
		case WhiteBishop:
			whiteBishops += 1
			whiteCentipawn += 5 * piece.Weight() * int16(rank+1) / DIVIDER
		case BlackBishop:
			blackBishops += 1
			blackCentipawn += 5 * piece.Weight() * (9 - int16(rank+1)) / DIVIDER
		case WhiteKnight:
			whiteKnights += 1
			whiteCentipawn += 5 * piece.Weight() * int16(rank+1) / DIVIDER
		case BlackKnight:
			blackKnights += 1
			blackCentipawn += 5 * piece.Weight() * (9 - int16(rank+1)) / DIVIDER
		case WhitePawn:
			white := whitePawnsPerFile[file]
			black := blackPawnsPerFile[file]
			bonus := int16(0)
			if black == 0 { // passed pawn
				if file != FileH {
					white := whitePawnsPerFile[file+1]
					black := blackPawnsPerFile[file+1]
					if white >= 1 && black == 0 { // supported passed pawn
						bonus = 10 * ((int16(rank+1) * 9) / DIVIDER) * int16(32-len(allPieces)) / 32
					} else if white >= 1 { // semi-supported passed pawn
						bonus = 5 * ((int16(rank+1) * 9) / DIVIDER) * int16(32-len(allPieces)) / 32
					} else {
						bonus = 2 * ((int16(rank+1) * 9) / DIVIDER) * int16(32-len(allPieces)) / 32
					}
				} else {
					bonus = 2 * ((int16(rank+1) * 9) / DIVIDER) * int16(32-len(allPieces)) / 32
				}
			}

			// backward pawn
			if rank != Rank7 && file != FileA && file != FileH {
				rPiece, rOk := allPieces[SquareOf(File(file-1), Rank(rank+1))]
				lPiece, lOk := allPieces[SquareOf(File(file+1), Rank(rank+1))]
				if rOk && lOk && rPiece == WhitePawn && lPiece == WhitePawn {
					whiteCentipawn -= 3
				}
			}

			// fawn pawn
			if rank == Rank6 {
				fPiece, fOk := allPieces[SquareOf(file, Rank(rank+1))] // pawn in front
				if fOk && fPiece == BlackPawn {
					if file == FileH {
						rPiece, rOk := allPieces[SquareOf(File(file-1), rank)] // neighbour pawn
						if rOk && rPiece == BlackPawn {
							whiteCentipawn += 5
						}
					} else if file == FileA {
						lPiece, lOk := allPieces[SquareOf(File(file+1), rank)] // another neighbour pawn
						if lOk && lPiece == BlackPawn {
							whiteCentipawn += 5
						}
					} else {
						rPiece, rOk := allPieces[SquareOf(File(file-1), rank)] // neighbour pawn
						lPiece, lOk := allPieces[SquareOf(File(file+1), rank)] // another neighbour pawn
						if (lOk && lPiece == BlackPawn) ||
							(rOk && rPiece == BlackPawn) {
							whiteCentipawn += 5
						}
					}
				}
			}

			// double pawns
			if white > 1 {
				bonus -= 3
			}
			whitePawns += 1
			whiteCentipawn += piece.Weight() + bonus
			whiteCentipawn += 2 * piece.Weight() * (int16(rank + 1)) / DIVIDER
		case BlackPawn:
			white := whitePawnsPerFile[file]
			black := blackPawnsPerFile[file]
			bonus := int16(0)
			if white == 0 { // passed pawn
				if file != FileH {
					white := whitePawnsPerFile[file+1]
					black := blackPawnsPerFile[file+1]
					if black >= 1 && white == 0 { // supported passed pawn
						bonus = 10 * ((9 - int16(rank+1)*9) / DIVIDER) * int16(32-len(allPieces)) / 32
					} else if black >= 1 { // semi-supported passed pawn
						bonus = 5 * ((9 - int16(rank+1)*9) / DIVIDER) * int16(32-len(allPieces)) / 32
					} else {
						bonus = 3 * ((9 - int16(rank+1)*9) / DIVIDER) * int16(32-len(allPieces)) / 32
					}
				} else {
					bonus = 3 * ((9 - int16(rank+1)*9) / DIVIDER) * int16(32-len(allPieces)) / 32
				}
			}

			// backward pawn
			if rank != Rank2 && file != FileA && file != FileH {
				rPiece, rOk := allPieces[SquareOf(File(file-1), Rank(rank-1))]
				lPiece, lOk := allPieces[SquareOf(File(file+1), Rank(rank-1))]
				if rOk && lOk && rPiece == BlackPawn && lPiece == BlackPawn {
					blackCentipawn -= 3
				}
			}

			// fawn pawn
			if rank == Rank2 {
				fPiece, fOk := allPieces[SquareOf(file, Rank(rank-1))] // pawn in front
				if fOk && fPiece == WhitePawn {
					if file == FileH {
						rPiece, rOk := allPieces[SquareOf(File(file-1), rank)] // neighbour pawn
						if rOk && rPiece == WhitePawn {
							blackCentipawn += 5
						}
					} else if file == FileA {
						lPiece, lOk := allPieces[SquareOf(File(file+1), rank)] // another neighbour pawn
						if lOk && lPiece == WhitePawn {
							blackCentipawn += 5
						}
					} else {
						rPiece, rOk := allPieces[SquareOf(File(file-1), rank)] // neighbour pawn
						lPiece, lOk := allPieces[SquareOf(File(file+1), rank)] // another neighbour pawn
						if (lOk && lPiece == WhitePawn) ||
							(rOk && rPiece == WhitePawn) {
							blackCentipawn += 5
						}
					}
				}
			}

			// double pawns
			if black > 1 {
				bonus -= 3
			}
			blackPawns += 1
			blackCentipawn += piece.Weight() + bonus
			blackCentipawn += 2 * piece.Weight() * (9 - int16(rank+1)) / DIVIDER
		}
	}
	pawns := blackPawns + whitePawns
	N := WhiteKnight
	B := WhiteBishop
	whiteCentipawn += whiteBishops * B.Weight() * int16(1+(16-pawns)/64)
	whiteCentipawn += whiteKnights * N.Weight() * int16(1-(16-pawns)/64)
	blackCentipawn += blackBishops * B.Weight() * int16(1+(16-pawns)/64)
	blackCentipawn += blackKnights * N.Weight() * int16(1-(16-pawns)/64)

	if whiteBishops >= 2 && blackBishops < 2 {
		whiteCentipawn += 3 + (8-blackPawns)/64
	}
	if whiteBishops < 2 && blackBishops >= 2 {
		blackCentipawn += 3 + (8-whitePawns)/64
	}
	if position.Turn() == White {
		return whiteCentipawn - blackCentipawn
	} else {
		return blackCentipawn - whiteCentipawn
	}
}

func pawnsPerFile(allPieces map[Square]Piece) (map[File](int8), map[File](int8)) {
	whites := make(map[File]int8)
	blacks := make(map[File]int8)

	for _, file := range Files {
		white, black := pawnsInFile(file, allPieces)
		whites[file] = white
		blacks[file] = black
	}

	return whites, blacks
}

func pawnsInFile(file File, allPieces map[Square]Piece) (int8, int8) {
	var blackPawn int8 = 0
	var whitePawn int8 = 0
	for _, rank := range Ranks {
		square := SquareOf(file, rank)
		piece, ok := allPieces[square]
		if ok {
			if piece == BlackPawn {
				blackPawn += 1
			} else if piece == WhitePawn {
				whitePawn += 1
			}
		}
	}

	return whitePawn, blackPawn
}

func pawnsPerRank(allPieces map[Square]Piece) (map[Rank](int8), map[Rank](int8)) {
	whites := make(map[Rank]int8)
	blacks := make(map[Rank]int8)

	for _, rank := range Ranks {
		white, black := pawnsInRank(rank, allPieces)
		whites[rank] = white
		blacks[rank] = black
	}

	return whites, blacks
}

func pawnsInRank(rank Rank, allPieces map[Square]Piece) (int8, int8) {
	var blackPawn int8 = 0
	var whitePawn int8 = 0
	for _, file := range Files {
		square := SquareOf(file, rank)
		piece, ok := allPieces[square]
		if ok {
			if piece == BlackPawn {
				blackPawn += 1
			} else if piece == WhitePawn {
				whitePawn += 1
			}
		}
	}

	return whitePawn, blackPawn
}

// TODO: Implement me
func center(position *Position, allPieces map[Square]Piece) int {
	return 0.0
}

// TODO: Implement me
func mobility(position *Position) int {
	return 0.0
}
