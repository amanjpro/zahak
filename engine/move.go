package engine

import (
	"fmt"
)

type Move uint32

/*
Move is represented as follows:
Given 32 bit int: 00000 000 0000 0000 0000 000000 000000

	- the lowest 6 bits represented the Source square
	- the next 6 bits represented the Destination square
	- the next 4 bits represent The moving piece
	- the next 4 bits represent The captured piece
	- the next 3 bits represent The promotion type
	- then next 5 bits are for the move tags
	   - first bit is king-side-castle
	   - second bit is queen-side-castle
	   - third bit is capture
	   - fourth bit is enpassant
	   - fifth bit is check

	6+6+4+4+3+5 = 28 bits, that leaves us 4 more bits in case
	more tags were needed
*/

const EmptyMove = Move(0)

func NewMove(from Square, to Square, movingPiece Piece, capturedPiece Piece, promoType PieceType, tag MoveTag) Move {
	s := uint32(from)                // the first 6 bits, offset = 0 = 0
	d := uint32(to) << 6             // the second 6 bits, offset = 0 + 6 = 6
	m := uint32(movingPiece) << 12   // next 4 bits, offset = 0 + 6 + 6 = 12
	c := uint32(capturedPiece) << 16 // next 4 bits, offset = 0 + 6 + 6 + 4 = 16
	p := uint32(promoType) << 20     // next 4 bits, offset = 0 + 6 + 6 + 4 + 4 = 20
	t := uint32(tag) << 23           // reminder, offset = 0 + 6 + 6 + 4 + 4 + 3 = 23
	mv := Move(s | d | m | c | p | t)
	return mv
}

type MoveTag uint8

const (
	// KingSideCastle indicates that the move is a king side castle.
	KingSideCastle MoveTag = 1 << iota
	// QueenSideCastle indicates that the move is a queen side castle.
	QueenSideCastle
	// Capture indicates that the move captures a piece.
	Capture
	// EnPassant indicates that the move captures via en passant.
	EnPassant
	// Check indicates that the move puts the opposing player in check.
	Check
)

func (m Move) Source() Square {
	return Square(uint32(m) & 0x3F)
}

func (m Move) Destination() Square {
	return Square(uint32(m) & 0xFC0 >> 6)
}

func (m Move) MovingPiece() Piece {
	return Piece(uint32(m) & 0xF000 >> 12)
}

func (m Move) CapturedPiece() Piece {
	return Piece(uint32(m) & 0xF0000 >> 16)
}

func (m Move) PromoType() PieceType {
	return PieceType(uint32(m) & 0x700000 >> 20)
}

func (m Move) Tag() MoveTag {
	return MoveTag(uint32(m >> 23))
}

func (m Move) IsKingSideCastle() bool {
	return uint32(m)&0x800000 != 0
}

func (m Move) IsQueenSideCastle() bool {
	return uint32(m)&0x1000000 != 0
}

func (m Move) IsCapture() bool {
	return uint32(m)&0x2000000 != 0
}

func (m Move) IsEnPassant() bool {
	return uint32(m)&0x4000000 != 0
}

func (m Move) IsCheck() bool {
	return uint32(m)&0x8000000 != 0
}

func (m *Move) AddCheckTag() {
	*m = Move(uint32(*m) | 0x8000000)
}

func (m Move) ToString() string {
	notation := fmt.Sprintf("%s%s", m.Source().Name(), m.Destination().Name())
	if m.PromoType() != NoType {
		// color doesn't matter here, I picked black as it prints lower case letters
		piece := GetPiece(m.PromoType(), Black)
		notation = fmt.Sprintf("%s%s", notation, piece.Name())
	}
	return notation
}
