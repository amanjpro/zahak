package engine

import (
	"fmt"
	"testing"
)

func TestNewMove(t *testing.T) {
	move := NewMove(E1, E2, WhiteKing, BlackRook, Queen, Check|Capture|KingSideCastle|QueenSideCastle|EnPassant)

	if move.Source() != E1 {
		t.Errorf("NewMove cannot assign source properly, %s", fmt.Sprintf("\nExpected: %b\nGot: %b", E1, move.Source()))
	}
	if move.Destination() != E2 {
		t.Errorf("NewMove cannot assign destination properly, %s", fmt.Sprintf("\nExpected: %b\nGot: %b", E2, move.Destination()))
	}
	if move.MovingPiece() != WhiteKing {
		t.Errorf("NewMove cannot assign moving piece properly, %s", fmt.Sprintf("\nExpected: %b\nGot: %b", WhiteKing, move.MovingPiece()))
	}
	if move.CapturedPiece() != BlackRook {
		t.Errorf("NewMove cannot assign captured piece properly, %s", fmt.Sprintf("\nExpected: %b\nGot: %b", BlackRook, move.CapturedPiece()))
	}
	if move.PromoType() != Queen {
		t.Errorf("NewMove cannot assign promo type properly, %s", fmt.Sprintf("\nExpected: %b\nGot: %b", Queen, move.PromoType()))
	}
	if !move.IsCheck() {
		t.Error("NewMove doesn't set check flag properly, Expected true, got false")
	}
	if !move.IsCapture() {
		t.Error("NewMove doesn't set capture flag properly, Expected true, got false")
	}
	if !move.IsEnPassant() {
		t.Error("NewMove doesn't set enpassant flag properly, Expected true, got false")
	}
	if !move.IsQueenSideCastle() {
		t.Error("NewMove doesn't set queen-side castle flag properly, Expected true, got false")
	}
	if !move.IsKingSideCastle() {
		t.Error("NewMove doesn't set king-side castle flag properly, Expected true, got false")
	}
}

func TestEnPassantFlag(t *testing.T) {
	move := NewMove(E1, E2, WhiteKing, BlackRook, Queen, EnPassant)

	if move.IsCheck() {
		t.Error("NewMove doesn't set check flag properly, Expected false, got true")
	}
	if move.IsCapture() {
		t.Error("NewMove doesn't set capture flag properly, Expected false, got true")
	}
	if !move.IsEnPassant() {
		t.Error("NewMove doesn't set enpassant flag properly, Expected true, got false")
	}
	if move.IsQueenSideCastle() {
		t.Error("NewMove doesn't set queen-side castle flag properly, Expected false, got true")
	}
	if move.IsKingSideCastle() {
		t.Error("NewMove doesn't set king-side castle flag properly, Expected false, got true")
	}
}

func TestCheckFlag(t *testing.T) {
	move := NewMove(E1, E2, WhiteKing, BlackRook, Queen, Check)

	if !move.IsCheck() {
		t.Error("NewMove doesn't set check flag properly, Expected true, got false")
	}
	if move.IsCapture() {
		t.Error("NewMove doesn't set capture flag properly, Expected false, got true")
	}
	if move.IsEnPassant() {
		t.Error("NewMove doesn't set enpassant flag properly, Expected false, got true")
	}
	if move.IsQueenSideCastle() {
		t.Error("NewMove doesn't set queen-side castle flag properly, Expected false, got true")
	}
	if move.IsKingSideCastle() {
		t.Error("NewMove doesn't set king-side castle flag properly, Expected false, got true")
	}
}

func TestCaptureFlag(t *testing.T) {
	move := NewMove(E1, E2, WhiteKing, BlackRook, Queen, Capture)

	if move.IsCheck() {
		t.Error("NewMove doesn't set check flag properly, Expected false, got true")
	}
	if !move.IsCapture() {
		t.Error("NewMove doesn't set capture flag properly, Expected true, got false")
	}
	if move.IsEnPassant() {
		t.Error("NewMove doesn't set enpassant flag properly, Expected false, got true")
	}
	if move.IsQueenSideCastle() {
		t.Error("NewMove doesn't set queen-side castle flag properly, Expected false, got true")
	}
	if move.IsKingSideCastle() {
		t.Error("NewMove doesn't set king-side castle flag properly, Expected false, got true")
	}
}

func TestQeenSideCastleFlag(t *testing.T) {
	move := NewMove(E1, E2, WhiteKing, BlackRook, Queen, QueenSideCastle)

	if move.IsCheck() {
		t.Error("NewMove doesn't set check flag properly, Expected false, got true")
	}
	if move.IsCapture() {
		t.Error("NewMove doesn't set capture flag properly, Expected false, got true")
	}
	if move.IsEnPassant() {
		t.Error("NewMove doesn't set enpassant flag properly, Expected false, got true")
	}
	if !move.IsQueenSideCastle() {
		t.Error("NewMove doesn't set queen-side castle flag properly, Expected true, got false")
	}
	if move.IsKingSideCastle() {
		t.Error("NewMove doesn't set king-side castle flag properly, Expected false, got true")
	}
}

func TestKingSideCastleFlag(t *testing.T) {
	move := NewMove(E1, E2, WhiteKing, BlackRook, Queen, KingSideCastle)

	if move.IsCheck() {
		t.Error("NewMove doesn't set check flag properly, Expected false, got true")
	}
	if move.IsCapture() {
		t.Error("NewMove doesn't set capture flag properly, Expected false, got true")
	}
	if move.IsEnPassant() {
		t.Error("NewMove doesn't set enpassant flag properly, Expected false, got true")
	}
	if move.IsQueenSideCastle() {
		t.Error("NewMove doesn't set queen-side castle flag properly, Expected false, got true")
	}
	if !move.IsKingSideCastle() {
		t.Error("NewMove doesn't set king-side castle flag properly, Expected true, got false")
	}
}

func TestAddCheckTag(t *testing.T) {
	move := NewMove(E1, E2, WhiteKing, BlackRook, Queen, 0)

	move.AddCheckTag()

	if !move.IsCheck() {
		t.Error("NewMove doesn't set check flag properly, Expected true, got false")
	}
	if move.IsCapture() {
		t.Error("NewMove doesn't set capture flag properly, Expected false, got true")
	}
	if move.IsEnPassant() {
		t.Error("NewMove doesn't set enpassant flag properly, Expected false, got true")
	}
	if move.IsQueenSideCastle() {
		t.Error("NewMove doesn't set queen-side castle flag properly, Expected false, got true")
	}
	if move.IsKingSideCastle() {
		t.Error("NewMove doesn't set king-side castle flag properly, Expected false, got true")
	}
}

func TestNewMoveWithPromo(t *testing.T) {
	m1 := NewMove(E2, F1, BlackPawn, WhiteKnight, Queen, Capture)
	m2 := NewMove(E2, F1, BlackPawn, WhiteKnight, Knight, Capture)

	if m1 == m2 {
		t.Error("The two moves are wrongfully deemed equal")
	}
	if m1.PromoType() != Queen {
		t.Error("The promo type of m1 is not stored correctly")
	}
	if m2.PromoType() != Knight {
		t.Error("The promo type of m2 is not stored correctly")
	}
}

func TestNoPromoType(t *testing.T) {
	m1 := NewMove(E2, F1, BlackPawn, WhiteKnight, NoType, Capture)

	if m1.PromoType() != NoType {
		t.Error("NoType promo type is not supported")
	}
}
