package evaluation

import (
	. "github.com/amanjpro/zahak/engine"
)

func Evaluate(position *Position) int16 {
	board := position.Board
	allPieces := board.AllPieces()
	return evaluate(position, allPieces)
}

const CHECKMATE_EVAL int16 = 3100
const DIVIDER int16 = 800

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
