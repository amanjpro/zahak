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
			hash ^= piecesZC[int8(p)][sq]
		}
	}

	return hash
}

func updateHashForNullMove(pos *Position) {
	if pos.hash == 0 {
		pos.Hash()
		return
	}
	/* Turn */
	pos.hash ^= whiteTurnZC

}

// capture square is provided for the case of enpassant
func updateHash(pos *Position, move *Move, movingPiece Piece, capturedPiece Piece,
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
			hash ^= piecesZC[int8(WhiteRook)][H1]
			hash ^= piecesZC[int8(WhiteRook)][F1]
		}
		if move.HasTag(QueenSideCastle) {
			hash ^= piecesZC[int8(WhiteRook)][A1]
			hash ^= piecesZC[int8(WhiteRook)][D1]
		}
	} else if move.Source == E8 { // Black
		if move.HasTag(KingSideCastle) {
			hash ^= piecesZC[int8(BlackRook)][H8]
			hash ^= piecesZC[int8(BlackRook)][F8]
		}
		if move.HasTag(QueenSideCastle) {
			hash ^= piecesZC[int8(BlackRook)][A8]
			hash ^= piecesZC[int8(BlackRook)][D8]
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
	hash ^= piecesZC[int8(movingPiece)][move.Source]
	if promoPiece != NoPiece {
		hash ^= piecesZC[int8(promoPiece)][move.Destination]
	} else {
		hash ^= piecesZC[int8(movingPiece)][move.Destination]
	}

	if capturedPiece != NoPiece {
		hash ^= piecesZC[int8(capturedPiece)][captureSquare]
	}

	pos.hash = hash
}
