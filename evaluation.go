package main

// import (
// 	"math"
// )

func eval(position *Position) float64 {
	board := position.board
	allPieces := board.AllPieces()
	return evaluate(position, &allPieces)
}

func evaluate(position *Position, allPieces *map[Square]Piece) float64 {
	// FIXME
	// if position.Status() == chess.Checkmate || position.Status() == chess.Resignation {
	// 	if position.Turn() == chess.Black {
	// 		return math.Inf(1)
	// 	}
	// 	return math.Inf(-1)
	// }
	//
	// if position.Status() != chess.NoMethod {
	// 	return 0.0
	// }

	whiteBishops := 0.0
	whiteKnights := 0.0
	blackBishops := 0.0
	blackKnights := 0.0
	blackPawns := 0.0
	whitePawns := 0.0
	centipawn := 0.0

	whitePawnsPerFile, blackPawnsPerFile := pawnsPerFile(allPieces)
	//whitePawnsPerRank, blackPawnsPerRank := pawnsPerRank(allPieces)

	for square, piece := range *allPieces {
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
				centipawn += 1
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
				centipawn -= 1
			}
		case WhiteQueen:
			centipawn += piece.Weight()
			centipawn += 0.025 * piece.Weight() * float64(rank) / 8
		case BlackQueen:
			centipawn -= piece.Weight()
			centipawn -= 0.025 * piece.Weight() * (9 - float64(rank)) / 8
		case WhiteRook:
			white := whitePawnsPerFile[file]
			black := blackPawnsPerFile[file]
			bonus := 0.0
			if white == 0 && black == 0 { // open file
				bonus = 1
			} else if white == 0 { // semi-open file
				bonus = 0.5
			}
			centipawn += piece.Weight() + bonus
			centipawn += 0.025 * piece.Weight() * (float64(rank)) / 8
		case BlackRook:
			white := whitePawnsPerFile[file]
			black := blackPawnsPerFile[file]
			bonus := 0.0
			if white == 0 && black == 0 { // open file
				bonus = 1
			} else if black == 0 { // semi-open file
				bonus = 0.5
			}
			centipawn -= piece.Weight() + bonus
			centipawn -= 0.025 * piece.Weight() * (9 - float64(rank)) / 8
		case WhiteBishop:
			whiteBishops += 1
			centipawn += 0.025 * piece.Weight() * (float64(rank)) / 8
		case BlackBishop:
			blackBishops += 1
			centipawn -= 0.025 * piece.Weight() * (9 - float64(rank)) / 8
		case WhiteKnight:
			whiteKnights += 1
			centipawn += 0.025 * piece.Weight() * (float64(rank)) / 8
		case BlackKnight:
			blackKnights += 1
			centipawn -= 0.025 * piece.Weight() * (9 - float64(rank)) / 8
		case WhitePawn:
			white := whitePawnsPerFile[file]
			black := blackPawnsPerFile[file]
			bonus := 0.0
			if black == 0 { // passed pawn
				if file != FileH {
					white := whitePawnsPerFile[file+1]
					black := blackPawnsPerFile[file+1]
					if white >= 1 && black == 0 { // supported passed pawn
						bonus = 1 * ((float64(rank) * 9) / 8) * (32 - float64(len(*allPieces))) / 32
					} else if white >= 1 { // semi-supported passed pawn
						bonus = 0.5 * ((float64(rank) * 9) / 8) * (32 - float64(len(*allPieces))) / 32
					} else {
						bonus = 0.25 * ((float64(rank) * 9) / 8) * (32 - float64(len(*allPieces))) / 32
					}
				} else {
					bonus = 0.25 * ((float64(rank) * 9) / 8) * (32 - float64(len(*allPieces))) / 32
				}
			}

			// backward pawn
			if rank != Rank7 && file != FileA && file != FileH {
				rPiece, rOk := (*allPieces)[SquareOf(File(file-1), Rank(rank+1))]
				lPiece, lOk := (*allPieces)[SquareOf(File(file+1), Rank(rank+1))]
				if rOk && lOk && rPiece == WhitePawn && lPiece == WhitePawn {
					centipawn -= 0.25
				}
			}

			// fawn pawn
			if rank == Rank6 {
				fPiece, fOk := (*allPieces)[SquareOf(file, Rank(rank+1))] // pawn in front
				if fOk && fPiece == BlackPawn {
					if file == FileH {
						rPiece, rOk := (*allPieces)[SquareOf(File(file-1), rank)] // neighbour pawn
						if rOk && rPiece == BlackPawn {
							centipawn += 0.25
						}
					} else if file == FileA {
						lPiece, lOk := (*allPieces)[SquareOf(File(file+1), rank)] // another neighbour pawn
						if lOk && lPiece == BlackPawn {
							centipawn += 0.25
						}
					} else {
						rPiece, rOk := (*allPieces)[SquareOf(File(file-1), rank)] // neighbour pawn
						lPiece, lOk := (*allPieces)[SquareOf(File(file+1), rank)] // another neighbour pawn
						if (lOk && lPiece == BlackPawn) ||
							(rOk && rPiece == BlackPawn) {
							centipawn += 0.25
						}
					}
				}
			}

			// double pawns
			if white > 1 {
				bonus -= 0.25
			}
			whitePawns += 1
			centipawn += piece.Weight() + bonus
			centipawn += 0.125 * piece.Weight() * (float64(rank)) / 8
		case BlackPawn:
			white := whitePawnsPerFile[file]
			black := blackPawnsPerFile[file]
			bonus := 0.0
			if white == 0 { // passed pawn
				if file != FileH {
					white := whitePawnsPerFile[file+1]
					black := blackPawnsPerFile[file+1]
					if black >= 1 && white == 0 { // supported passed pawn
						bonus = 1 * ((9 - float64(rank)*9) / 8) * (32 - float64(len(*allPieces))) / 32
					} else if black >= 1 { // semi-supported passed pawn
						bonus = 0.5 * ((9 - float64(rank)*9) / 8) * (32 - float64(len(*allPieces))) / 32
					} else {
						bonus = 0.25 * ((9 - float64(rank)*9) / 8) * (32 - float64(len(*allPieces))) / 32
					}
				} else {
					bonus = 0.25 * ((9 - float64(rank)*9) / 8) * (32 - float64(len(*allPieces))) / 32
				}
			}

			// backward pawn
			if rank != Rank2 && file != FileA && file != FileH {
				rPiece, rOk := (*allPieces)[SquareOf(File(file-1), Rank(rank-1))]
				lPiece, lOk := (*allPieces)[SquareOf(File(file+1), Rank(rank-1))]
				if rOk && lOk && rPiece == BlackPawn && lPiece == BlackPawn {
					centipawn += 0.25
				}
			}

			// fawn pawn
			if rank == Rank2 {
				fPiece, fOk := (*allPieces)[SquareOf(file, Rank(rank-1))] // pawn in front
				if fOk && fPiece == WhitePawn {
					if file == FileH {
						rPiece, rOk := (*allPieces)[SquareOf(File(file-1), rank)] // neighbour pawn
						if rOk && rPiece == WhitePawn {
							centipawn -= 0.25
						}
					} else if file == FileA {
						lPiece, lOk := (*allPieces)[SquareOf(File(file+1), rank)] // another neighbour pawn
						if lOk && lPiece == WhitePawn {
							centipawn -= 0.25
						}
					} else {
						rPiece, rOk := (*allPieces)[SquareOf(File(file-1), rank)] // neighbour pawn
						lPiece, lOk := (*allPieces)[SquareOf(File(file+1), rank)] // another neighbour pawn
						if (lOk && lPiece == WhitePawn) ||
							(rOk && rPiece == WhitePawn) {
							centipawn -= 0.25
						}
					}
				}
			}

			// double pawns
			if black > 1 {
				bonus -= 0.25
			}
			blackPawns += 1
			centipawn -= piece.Weight() + bonus
			centipawn += 0.125 * piece.Weight() * (9 - float64(rank)) / 8
		}
	}
	pawns := blackPawns + whitePawns
	N := WhiteKnight
	B := WhiteBishop
	centipawn += whiteBishops * B.Weight() * (1 + (16-pawns)/64)
	centipawn -= blackBishops * B.Weight() * (1 + (16-pawns)/64)
	centipawn += whiteKnights * N.Weight() * (1 - (16-pawns)/64)
	centipawn -= blackKnights * N.Weight() * (1 - (16-pawns)/64)

	if whiteBishops >= 2 && blackBishops < 2 {
		centipawn += 0.3 + (8-blackPawns)/64
	}
	if whiteBishops < 2 && blackBishops >= 2 {
		centipawn -= 0.3 + (8-whitePawns)/64
	}
	return centipawn
}

func pawnsPerFile(allPieces *map[Square]Piece) (map[File](int8), map[File](int8)) {
	whites := make(map[File]int8)
	blacks := make(map[File]int8)

	for _, file := range files {
		white, black := pawnsInFile(file, allPieces)
		whites[file] = white
		blacks[file] = black
	}

	return whites, blacks
}

func pawnsInFile(file File, allPieces *map[Square]Piece) (int8, int8) {
	var blackPawn int8 = 0
	var whitePawn int8 = 0
	for _, rank := range ranks {
		square := SquareOf(file, rank)
		piece, ok := (*allPieces)[square]
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

func pawnsPerRank(allPieces *map[Square]Piece) (map[Rank](int8), map[Rank](int8)) {
	whites := make(map[Rank]int8)
	blacks := make(map[Rank]int8)

	for _, rank := range ranks {
		white, black := pawnsInRank(rank, allPieces)
		whites[rank] = white
		blacks[rank] = black
	}

	return whites, blacks
}

func pawnsInRank(rank Rank, allPieces *map[Square]Piece) (int8, int8) {
	var blackPawn int8 = 0
	var whitePawn int8 = 0
	for _, file := range files {
		square := SquareOf(file, rank)
		piece, ok := (*allPieces)[square]
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

func center(position *Position, allPieces *map[Square]Piece) float64 {
	return 0.0
}
