package main

import (
	"github.com/notnil/chess"
)

// func eval(position *chess.Position) float64 {
// 	board := position.Board()
// 	allPieces := board.SquareMap()
// 	return (pieceCount(position, &allPieces) +
// 		files(position, &allPieces) +
// 		bishopPairs(position, &allPieces) +
// 		passedPawns(position, &allPieces) +
// 		mobility(position, &allPieces) +
// 		castling(position, &allPieces) +
// 		center(position, &allPieces) +
// 		backwardPawn(position, &allPieces) +
// 		doublePawns(position, &allPieces))
// }

func eval(position *chess.Position) float64 {
	board := position.Board()
	allPieces := board.SquareMap()
	return evaluate(position, &allPieces)
}

func evaluate(position *chess.Position, allPieces *map[chess.Square]chess.Piece) float64 {
	whiteBishops := 0.0
	whiteKnights := 0.0
	blackBishops := 0.0
	blackKnights := 0.0
	pawns := 0.0
	centipawn := 0.0

	whitePawnsPerFile, blackPawnsPerFile := pawnsPerFile(allPieces)

	for square, piece := range *allPieces {
		file := square.File()
		rank := square.Rank()
		switch piece {
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
			centipawn -= 5 + bonus
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
						bonus = 1 + (float64(rank)*9)/8
					} else if white >= 1 { // semi-supported passed pawn
						bonus = 0.5 + (float64(rank)*9)/8
					} else {
						bonus = 0.25 + (float64(rank)*9)/8
					}
				} else {
					bonus = 0.25 + (float64(rank)*9)/8
				}
			}

			if white > 1 {
				bonus -= 0.25
			}
			pawns += 1
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
						bonus = 1 + (9-float64(rank)*9)/8
					} else if black >= 1 { // semi-supported passed pawn
						bonus = 0.5 + (9-float64(rank)*9)/8
					} else {
						bonus = 0.25 + (9-float64(rank)*9)/8
					}
				} else {
					bonus = 0.25 + (9-float64(rank)*9)/8
				}
			}

			if black > 1 {
				bonus -= 0.25
			}
			pawns += 1
			centipawn -= 1 + bonus
		}
	}
	centipawn += whiteBishops * 3 * (1 + (16-pawns)/64)
	centipawn += blackBishops * -3 * (1 + (16-pawns)/64)
	centipawn += whiteKnights * 3 * (1 - (16-pawns)/64)
	centipawn += blackKnights * -3 * (1 - (16-pawns)/64)
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

func files(position *chess.Position, allPieces *map[chess.Square]chess.Piece) float64 {
	return 0.0
}

func bishopPairs(position *chess.Position, allPieces *map[chess.Square]chess.Piece) float64 {
	return 0.0
}

func passedPawns(position *chess.Position, allPieces *map[chess.Square]chess.Piece) float64 {
	return 0.0
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

func doublePawns(position *chess.Position, allPieces *map[chess.Square]chess.Piece) float64 {
	return 0.0
}
