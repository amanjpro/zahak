package engine

import (
	"fmt"
	"testing"
)

func TestSimpleStaticExchangeEval(t *testing.T) {
	fen := "1k1r4/1pp4p/p7/4p3/8/P5P1/1PP4P/2K1R3 w - - 0 1"
	game := FromFen(fen, true)
	board := game.position.Board

	actual := board.StaticExchangeEval(E5, BlackPawn, E1, WhiteRook)
	expected := int32(100)

	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}
}

func TestComplicatedStaticExchangeEval(t *testing.T) {
	fen := "1k1r3q/1ppn3p/p4b2/4p3/8/P2N2P1/1PP1R1BP/2K1Q3 w - - 0 1"
	game := FromFen(fen, true)
	board := game.position.Board

	actual := board.StaticExchangeEval(E5, BlackPawn, D3, WhiteKnight)
	expected := int32(-220)

	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}
}

func TestSlidingPiecesStaticExchangeEval(t *testing.T) {
	fen := "k3r3/pp2r3/2b5/3p4/4P3/5P2/PP2R3/K3R2B b - - 0 1"
	game := FromFen(fen, true)
	board := game.position.Board

	actual := board.StaticExchangeEval(E4, WhitePawn, D5, BlackPawn)

	if actual > 0 {
		t.Error(fmt.Sprintf("Expected: a non-positive number\n, Got: %d\n", actual))
	}
}

func TestSlidingPiecesStaticExchangeEvalPositive(t *testing.T) {
	fen := "k3r3/pp2r3/2b5/3p4/4P3/8/PP2R3/K3R2B b - - 0 1"
	game := FromFen(fen, true)
	board := game.position.Board

	actual := board.StaticExchangeEval(E4, WhitePawn, D5, BlackPawn)
	expected := int32(100)

	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}
}
