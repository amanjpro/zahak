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
)

type PieceType int8

const (
	Pawn PieceType = iota
	Knight
	Bishop
	Rook
	Queen
	King
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
	}
	return King
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
