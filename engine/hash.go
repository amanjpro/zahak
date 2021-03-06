package engine

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
	enPassant := pos.EnPassant
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
	board := pos.Board
	for sq := A1; sq <= H8; sq++ {
		p := board.PieceAt(sq)
		if p != NoPiece {
			hash ^= piecesZC[int8(p)-1][sq]
		}
	}

	return hash
}

func updateHashForNullMove(pos *Position, newEnPassant Square, oldEnPassant Square) {
	if pos.hash == 0 {
		pos.Hash()
		return
	}
	var hash uint64 = pos.hash
	/* Turn */
	hash ^= whiteTurnZC

	turn := pos.Turn()
	/* En passant */
	if newEnPassant != NoSquare {
		if turn == Black {
			/* Next mov Black -> Current pos White -> White en passant square */
			hash ^= enPassantZC[newEnPassant-16]
		} else {
			/* Next mov White -> Current pos Black -> Black en passant square */
			hash ^= enPassantZC[newEnPassant-40+8]
		}
	}

	if oldEnPassant != NoSquare {
		if turn == Black {
			/* Previous mov Black -> Current pos White -> Black en passant square */
			hash ^= enPassantZC[oldEnPassant-40+8]
		} else {
			/* Previous mov White -> Current pos Black -> White en passant square */
			hash ^= enPassantZC[oldEnPassant-16]
		}
	}

	pos.hash = hash
}

// capture square is provided for the case of enpassant
func updateHash(pos *Position, move Move, movingPiece Piece, capturedPiece Piece,
	captureSquare Square, newEnPassant Square, oldEnPassant Square, promoPiece Piece,
	oldPositionTag PositionTag) {
	var hash uint64 = pos.hash
	if hash == 0 {
		pos.Hash()
		return
	}
	/* Turn */
	hash ^= whiteTurnZC
	turn := pos.Turn()

	/* Castle */
	if move.Source == E1 { // White
		if move.HasTag(KingSideCastle) {
			hash ^= piecesZC[int8(WhiteRook)-1][H1]
			hash ^= piecesZC[int8(WhiteRook)-1][F1]
		}
		if move.HasTag(QueenSideCastle) {
			hash ^= piecesZC[int8(WhiteRook)-1][A1]
			hash ^= piecesZC[int8(WhiteRook)-1][D1]
		}
	} else if move.Source == E8 { // Black
		if move.HasTag(KingSideCastle) {
			hash ^= piecesZC[int8(BlackRook)-1][H8]
			hash ^= piecesZC[int8(BlackRook)-1][F8]
		}
		if move.HasTag(QueenSideCastle) {
			hash ^= piecesZC[int8(BlackRook)-1][A8]
			hash ^= piecesZC[int8(BlackRook)-1][D8]
		}
	}

	if oldPositionTag&WhiteCanCastleKingSide != pos.Tag&WhiteCanCastleKingSide {
		hash ^= castleRightsZC[0]
	}
	if oldPositionTag&WhiteCanCastleQueenSide != pos.Tag&WhiteCanCastleQueenSide {
		hash ^= castleRightsZC[1]
	}
	if oldPositionTag&BlackCanCastleKingSide != pos.Tag&BlackCanCastleKingSide {
		hash ^= castleRightsZC[2]
	}
	if oldPositionTag&BlackCanCastleQueenSide != pos.Tag&BlackCanCastleQueenSide {
		hash ^= castleRightsZC[3]
	}

	/* En passant */
	if newEnPassant != NoSquare {
		if turn == Black {
			/* Next mov Black -> Current pos White -> White en passant square */
			hash ^= enPassantZC[newEnPassant-16]
		} else {
			/* Next mov White -> Current pos Black -> Black en passant square */
			hash ^= enPassantZC[newEnPassant-40+8]
		}
	}

	if oldEnPassant != NoSquare {
		if turn == Black {
			/* Previous mov Black -> Current pos White -> Black en passant square */
			hash ^= enPassantZC[oldEnPassant-40+8]
		} else {
			/* Previous mov White -> Current pos Black -> White en passant square */
			hash ^= enPassantZC[oldEnPassant-16]
		}
	}

	/* Board */
	hash ^= piecesZC[int8(movingPiece)-1][move.Source]
	if promoPiece != NoPiece {
		hash ^= piecesZC[int8(promoPiece)-1][move.Destination]
	} else {
		hash ^= piecesZC[int8(movingPiece)-1][move.Destination]
	}

	if capturedPiece != NoPiece {
		hash ^= piecesZC[int8(capturedPiece)-1][captureSquare]
	}

	pos.hash = hash
}
