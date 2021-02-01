package main

import (
	"math"
)

type Piece int8

const (
	WhitePawn Piece = iota
	BlackPawn
	WhiteKnight
	BlackKnight
	WhiteBishop
	BlackBishop
	WhiteRook
	BlackRook
	WhiteQueen
	BlackQueen
	WhiteKing
	BlackKing
	NoPiece
)

type PieceType int8

const (
	Pawn PieceType = iota
	Knight
	Bishop
	Rook
	Queen
	King
	NoType
)

type Color int8

const (
	White Color = iota
	Black
)

func (p *Piece) Type() PieceType {
	switch *p {
	case WhitePawn | BlackPawn:
		return Pawn
	case WhiteKnight | BlackKnight:
		return Knight
	case WhiteBishop | BlackBishop:
		return Bishop
	case WhiteRook | BlackRook:
		return Rook
	case WhiteQueen | BlackQueen:
		return Queen
	case WhiteKing | BlackKing:
		return King
	}
	return NoType
}

func (p *Piece) Weight() float64 {
	switch *p {
	case WhitePawn | BlackPawn:
		return 1
	case WhiteKnight | BlackKnight:
		return 3
	case WhiteBishop | BlackBishop:
		return 3
	case WhiteRook | BlackRook:
		return 5
	case WhiteQueen | BlackQueen:
		return 9
	}
	return math.Inf(1)
}

func (p *Piece) Name() string {
	switch *p {
	case WhitePawn:
		return "P"
	case WhiteKnight:
		return "N"
	case WhiteBishop:
		return "B"
	case WhiteRook:
		return "R"
	case WhiteQueen:
		return "Q"
	case WhiteKing:
		return "K"
	}

	switch *p {
	case BlackPawn:
		return "p"
	case BlackKnight:
		return "n"
	case BlackBishop:
		return "b"
	case BlackRook:
		return "r"
	case BlackQueen:
		return "q"
	}
	return "k"
}

func (p *Piece) Color() Color {
	switch *p {
	case WhitePawn | WhiteKnight | WhiteBishop | WhiteRook | WhiteQueen | WhiteKing:
		return White
	}
	return Black
}
