package evaluation

import (
	"fmt"
	"testing"

	. "github.com/amanjpro/zahak/engine"
)

func TestMaterialValue(t *testing.T) {
	fen := "rnb2bnr/ppqppkpp/8/2p5/4P3/8/PPPP1PPP/RNB1KBNR w KQ - 0 1"
	game := FromFen(fen, false)

	actual := Evaluate(game.Position())

	if actual >= 0 {
		t.Errorf("Expected: a negative number\nGot: %d\n", actual)
	}

	fen = "3N1k2/N7/1p2ppR1/1P6/P2pP3/3Pb3/1K5n/8 w - - 0 4"

	game = FromFen(fen, false)

	actual = Evaluate(game.Position())

	if actual <= 0 {
		t.Errorf("Expected: a positive number\nGot: %d\n", actual)
	}

	fen = "2k2b1r/ppp1pppp/4b3/1P6/2P3P1/3BKP1P/7B/1R4N1 b - - 0 23"
	game = FromFen(fen, false)

	actual = Evaluate(game.Position())

	if actual >= 0 {
		t.Errorf("Expected: a negative number\nGot: %d\n", actual)
	}
}

func TestIsBackwardsPawn(t *testing.T) {
	game := FromFen("k7/5p2/4p1p1/8/8/4P1P1/5P2/K7 w - - 0 1", false)
	board := game.Position().Board

	actual := board.IsBackwardPawn(uint64(1<<int(E6)), board.GetBitboardOf(BlackPawn), Black)
	expected := false

	if actual != expected {
		t.Error("Non-Backward Pawn - Black: is falsely identified")
	}

	actual = board.IsBackwardPawn(uint64(1<<int(F7)), board.GetBitboardOf(BlackPawn), Black)
	expected = true

	if actual != expected {
		t.Error("Non-Backward Pawn - Black: is not identified")
	}

	actual = board.IsBackwardPawn(uint64(1<<int(G3)), board.GetBitboardOf(WhitePawn), White)
	expected = false

	if actual != expected {
		t.Error("Non-Backward Pawn - White: is falsely identified")
	}

	actual = board.IsBackwardPawn(uint64(1<<int(F2)), board.GetBitboardOf(WhitePawn), White)
	expected = true

	if actual != expected {
		t.Error("Non-Backward Pawn - White: is not identified")
	}
}

func TestPawnStructureEval(t *testing.T) {
	fen := "k7/4pp2/6p1/8/8/4P1P1/5P2/K7 w - - 0 1"
	game := FromFen(fen, false)

	actual := Evaluate(game.Position())
	expected := int32(-17)

	if actual != expected {
		err := fmt.Sprintf("Backward Pawn - White:\nExpected: %d\nGot: %d\n", expected, actual)
		t.Errorf(err)
	}

	fen = "k7/5p2/4p1p1/8/8/6P1/4PP2/K7 b - - 0 1"
	game = FromFen(fen, false)

	actual = Evaluate(game.Position())
	expected = int32(-17)

	if actual != expected {
		err := fmt.Sprintf("Backward Pawn - Black:\nExpected: %d\nGot: %d\n", expected, actual)
		t.Errorf(err)
	}
}
