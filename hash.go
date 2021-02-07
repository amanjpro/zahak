package main

import (
	"math/rand"
)

var piecesZC [12][64]uint64
var castleRightsZC [4]uint64
var enPassantZC [16]uint64
var whiteTurnZC uint64

func initZobrist() {
	whiteTurnZC = rand.Uint64()
	for i := 0; i < 12; i++ {
		for j := 0; j < 64; j++ {
			piecesZC[i][j] = rand.Uint64()
		}
	}
	for i := 0; i < 4; i++ {
		castleRightsZC[i] = rand.Uint64()
	}
	for i := 0; i < 16; i++ {
		enPassantZC[i] = rand.Uint64()
	}
}

func generateZobristHash(pos *Position) uint64 {
	var hash uint64 = 0
	/* Turn */
	if pos.Turn() == White {
		hash ^= whiteTurnZC
	}

	/* Castle */
	if pos.HasTag(WhiteCanCastleKingSide) {
		hash ^= castleRightsZC[0]
	}
	if pos.HasTag(WhiteCanCastleQueenSide) {
		hash ^= castleRightsZC[1]
	}
	if pos.HasTag(BlackCanCastleKingSide) {
		hash ^= castleRightsZC[2]
	}
	if pos.HasTag(BlackCanCastleQueenSide) {
		hash ^= castleRightsZC[3]
	}

	/* En passant */
	enPassant := pos.enPassant
	if enPassant != NoSquare {
		if pos.Turn() == Black {
			/* Next mov Black -> Current pos White -> White en passant square */
			hash ^= enPassantZC[enPassant-16]
		} else {
			/* Next mov White -> Current pos Black -> Black en passant square */
			hash ^= enPassantZC[enPassant-40+8]
		}
	}

	/* Board */
	board := pos.board
	for sq := A1; sq <= H8; sq++ {
		p := board.PieceAt(sq)
		if p != NoPiece {
			hash ^= piecesZC[int8(p)][sq]
		}
	}

	return hash
}

func init() {
	initZobrist()
}
