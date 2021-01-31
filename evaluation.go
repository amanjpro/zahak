package main

import (
	"math"

	"github.com/notnil/chess"
)

func eval(position *chess.Position) float64 {
	board := position.Board()
	allPieces := board.SquareMap()
	return evaluate(position, &allPieces)
}

func evaluate(position *chess.Position, allPieces *map[chess.Square]chess.Piece) float64 {
	if position.Status() == chess.Checkmate || position.Status() == chess.Resignation {
		if position.Turn() == chess.Black {
			return math.Inf(1)
		}
		return math.Inf(-1)
	}

	if position.Status() != chess.NoMethod {
		return 0.0
	}

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
		case chess.WhiteKing:
			// This doesn't consider endgame
			if position.CastleRights().CanCastle(chess.White, chess.KingSide) ||
				position.CastleRights().CanCastle(chess.White, chess.QueenSide) ||
				square == chess.A1 || square == chess.A2 ||
				square == chess.B1 || square == chess.B2 ||
				square == chess.F1 || square == chess.F2 ||
				square == chess.G1 || square == chess.G2 ||
				square == chess.H1 || square == chess.H2 {
				centipawn += 1
			}
		case chess.BlackKing:
			// This doesn't consider endgame
			if position.CastleRights().CanCastle(chess.Black, chess.KingSide) ||
				position.CastleRights().CanCastle(chess.Black, chess.QueenSide) ||
				square == chess.A7 || square == chess.A8 ||
				square == chess.B7 || square == chess.B8 ||
				square == chess.F7 || square == chess.F8 ||
				square == chess.G7 || square == chess.G8 ||
				square == chess.H7 || square == chess.H8 {
				centipawn -= 1
			}
		case chess.WhiteQueen:
			centipawn += 9
		case chess.BlackQueen:
			centipawn -= 9
		case chess.WhiteRook:
			white := whitePawnsPerFile[file]
			black := blackPawnsPerFile[file]
			bonus := 0.0
			if white == 0 && black == 0 { // open file
				bonus = 1
			} else if white == 0 { // semi-open file
				bonus = 0.5
			}
			centipawn += 5 + bonus
		case chess.BlackRook:
			white := whitePawnsPerFile[file]
			black := blackPawnsPerFile[file]
			bonus := 0.0
			if white == 0 && black == 0 { // open file
				bonus = 1
			} else if black == 0 { // semi-open file
				bonus = 0.5
			}
			centipawn -= 5 + bonus
		case chess.WhiteBishop:
			whiteBishops += 1
		case chess.BlackBishop:
			blackBishops += 1
		case chess.WhiteKnight:
			whiteKnights += 1
		case chess.BlackKnight:
			blackKnights += 1
		case chess.WhitePawn:
			white := whitePawnsPerFile[file]
			black := blackPawnsPerFile[file]
			bonus := 0.0
			if black == 0 { // passed pawn
				if file != chess.FileH {
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

			if white > 1 {
				bonus -= 0.25
			}
			whitePawns += 1
			centipawn += 1 + bonus
		case chess.BlackPawn:
			white := whitePawnsPerFile[file]
			black := blackPawnsPerFile[file]
			bonus := 0.0
			if white == 0 { // passed pawn
				if file != chess.FileH {
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

			if black > 1 {
				bonus -= 0.25
			}
			blackPawns += 1
			centipawn -= 1 + bonus
		}
	}
	pawns := blackPawns + whitePawns
	centipawn += whiteBishops * 3 * (1 + (16-pawns)/64)
	centipawn -= blackBishops * 3 * (1 + (16-pawns)/64)
	centipawn += whiteKnights * 3 * (1 - (16-pawns)/64)
	centipawn -= blackKnights * 3 * (1 - (16-pawns)/64)

	if whiteBishops >= 2 && blackBishops < 2 {
		centipawn += 0.3 + (8-blackPawns)/64
	}
	if whiteBishops < 2 && blackBishops >= 2 {
		centipawn -= 0.3 + (8-whitePawns)/64
	}
	return centipawn
}

func pawnsPerFile(allPieces *map[chess.Square]chess.Piece) (map[chess.File](int8), map[chess.File](int8)) {
	whites := make(map[chess.File]int8)
	blacks := make(map[chess.File]int8)

	files := [8]chess.File{chess.FileA, chess.FileB, chess.FileC, chess.FileD, chess.FileE, chess.FileF, chess.FileG, chess.FileH}

	for _, file := range files {
		white, black := pawnsInFile(file, allPieces)
		whites[file] = white
		blacks[file] = black
	}

	return whites, blacks
}

func pawnsInFile(file chess.File, allPieces *map[chess.Square]chess.Piece) (int8, int8) {
	ranks := [8]int{0, 1, 2, 3, 4, 5, 6, 7}
	var blackPawn int8 = 0
	var whitePawn int8 = 0
	for _, rank := range ranks {
		square := chess.Square((rank * 8) + int(file))
		piece, ok := (*allPieces)[square]
		if ok {
			if piece == chess.BlackPawn {
				blackPawn += 1
			} else if piece == chess.WhitePawn {
				whitePawn += 1
			}
		}
	}

	return whitePawn, blackPawn
}

func pawnsPerRank(allPieces *map[chess.Square]chess.Piece) (map[chess.Rank](int8), map[chess.Rank](int8)) {
	whites := make(map[chess.Rank]int8)
	blacks := make(map[chess.Rank]int8)

	ranks := [8]chess.Rank{chess.Rank1, chess.Rank2, chess.Rank3, chess.Rank4, chess.Rank5, chess.Rank6, chess.Rank7, chess.Rank8}

	for _, rank := range ranks {
		white, black := pawnsInRank(rank, allPieces)
		whites[rank] = white
		blacks[rank] = black
	}

	return whites, blacks
}

func pawnsInRank(rank chess.Rank, allPieces *map[chess.Square]chess.Piece) (int8, int8) {
	files := [8]int{0, 1, 2, 3, 4, 5, 6, 7}
	var blackPawn int8 = 0
	var whitePawn int8 = 0
	for _, file := range files {
		square := chess.Square((int(rank) * 8) + file)
		piece, ok := (*allPieces)[square]
		if ok {
			if piece == chess.BlackPawn {
				blackPawn += 1
			} else if piece == chess.WhitePawn {
				whitePawn += 1
			}
		}
	}

	return whitePawn, blackPawn
}

func mobility(position *chess.Position, allPieces *map[chess.Square]chess.Piece) float64 {
	return 0.0
}

func castling(position *chess.Position, allPieces *map[chess.Square]chess.Piece) float64 {
	return 0.0
}

func center(position *chess.Position, allPieces *map[chess.Square]chess.Piece) float64 {
	return 0.0
}

func backwardPawn(position *chess.Position, allPieces *map[chess.Square]chess.Piece) float64 {
	return 0.0
}
