package main

import (
	"github.com/notnil/chess"
)

func eval(position *chess.Position) float64 {
	return (pieceCount(position) +
		files(position) +
		bishopPairs(position) +
		passedPawns(position) +
		mobility(position) +
		castling(position) +
		center(position) +
		backwardPawn(position) +
		doublePawns(position))
}

func pieceCount(position *chess.Position) float64 {
	board := position.Board()
	allPieces := board.SquareMap()
	whiteBishops := 0.0
	whiteKnights := 0.0
	blackBishops := 0.0
	blackKnights := 0.0
	pawns := 0.0
	total := 0.0

	for _, piece := range allPieces {
		switch piece {
		case chess.WhiteQueen:
			total += 9
		case chess.BlackQueen:
			total -= 9
		case chess.WhiteRook:
			total += 5
		case chess.BlackRook:
			total -= 5
		case chess.WhiteBishop:
			whiteBishops += 1
		case chess.BlackBishop:
			blackBishops += 1
		case chess.WhiteKnight:
			whiteKnights += 1
		case chess.BlackKnight:
			blackKnights += 1
		case chess.WhitePawn:
			pawns += 1
			total += 1
		case chess.BlackPawn:
			pawns += 1
			total -= 1
		}
	}
	total += whiteBishops * 3 * (1 + (16-pawns)/64)
	total += blackBishops * -3 * (1 + (16-pawns)/64)
	total += whiteKnights * 3 * (1 - (16-pawns)/64)
	total += blackKnights * -3 * (1 - (16-pawns)/64)
	return total
}

func files(position *chess.Position) float64 {
	return 0.0
}

func bishopPairs(position *chess.Position) float64 {
	return 0.0
}

func passedPawns(position *chess.Position) float64 {
	return 0.0
}

func mobility(position *chess.Position) float64 {
	return 0.0
}

func castling(position *chess.Position) float64 {
	return 0.0
}

func center(position *chess.Position) float64 {
	return 0.0
}

func backwardPawn(position *chess.Position) float64 {
	return 0.0
}

func doublePawns(position *chess.Position) float64 {
	return 0.0
}
