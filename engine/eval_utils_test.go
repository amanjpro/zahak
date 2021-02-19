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
	expected := 100

	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}
}

func TestComplicatedStaticExchangeEval(t *testing.T) {
	fen := "1k1r3q/1ppn3p/p4b2/4p3/8/P2N2P1/1PP1R1BP/2K1Q3 w - - 0 1"
	game := FromFen(fen, true)
	board := game.position.Board

	actual := board.StaticExchangeEval(E5, BlackPawn, D3, WhiteKnight)
	expected := -200

	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}
}
