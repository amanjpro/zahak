package cmd

import (
	"fmt"
)

type Move struct {
	source      Square
	destination Square
	promoType   PieceType
	moveTag     MoveTag
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

func (m *Move) SetTag(tag MoveTag)      { m.moveTag |= tag }
func (m *Move) ClearTag(tag MoveTag)    { m.moveTag &^= tag }
func (m *Move) ToggleTag(tag MoveTag)   { m.moveTag ^= tag }
func (m *Move) HasTag(tag MoveTag) bool { return m.moveTag&tag != 0 }

func (m *Move) ToString() string {
	notation := fmt.Sprintf("%s%s", m.source.Name(), m.destination.Name())
	if m.promoType != NoType {
		// color doesn't matter here, I picked black as it prints lower case letters
		piece := getPiece(m.promoType, Black)
		notation = fmt.Sprintf("%s%s", notation, piece.Name())
	}
	return notation
}
