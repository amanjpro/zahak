package engine

import (
	"fmt"
	"math/bits"
	"testing"
)

func TestSimpleStaticExchangeEval(t *testing.T) {
	fen := "1k1r4/1pp4p/p7/4p3/8/P5P1/1PP4P/2K1R3 w - - 0 1"
	game := FromFen(fen, true)
	board := game.position.Board

	actual := board.StaticExchangeEval(E5, BlackPawn, E1, WhiteRook)
	expected := int16(100)

	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}
}

func TestComplicatedStaticExchangeEval(t *testing.T) {
	fen := "1k1r3q/1ppn3p/p4b2/4p3/8/P2N2P1/1PP1R1BP/2K1Q3 w - - 0 1"
	game := FromFen(fen, true)
	board := game.position.Board

	actual := board.StaticExchangeEval(E5, BlackPawn, D3, WhiteKnight)
	expected := int16(-200)

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
	expected := int16(100)

	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}
}

func TestAllAttacksWhite(t *testing.T) {
	fen := "3k4/8/8/8/8/2Q5/8/2K5 w - - 0 1"
	game := FromFen(fen, true)
	board := game.position.Board

	actualP, actualM, actualO := board.AllAttacks(White)
	actual := bits.OnesCount64(actualP | actualM | actualO)
	expected := 26

	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}
}

func TestAllAttacksBlack(t *testing.T) {
	fen := "3k4/8/8/8/8/2Q5/8/2K5 w - - 0 1"
	game := FromFen(fen, true)
	board := game.position.Board

	actualP, actualM, actualO := board.AllAttacks(Black)
	actual := bits.OnesCount64(actualP | actualM | actualO)
	expected := 5

	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}
}

func TestBackwardPawns(t *testing.T) {
	fen := "rnbqkbnr/5ppp/p2p4/P2P1P2/1P2P3/8/6PP/RNBQKBNR w KQkq - 0 1"
	game := FromFen(fen, true)

	expected := int16(2)
	actual := game.Position().CountBackwardPawns(White)
	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}

	expected = int16(2)
	actual = game.Position().CountBackwardPawns(Black)
	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}
}

func TestCandidatePawns(t *testing.T) {
	fen := "7k/p7/8/PP6/5ppp/8/5P1P/K7 w - - 0 1"
	game := FromFen(fen, true)

	expected := int16(1)
	actual := game.Position().CountCandidatePawns(White)
	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}

	expected = int16(1)
	actual = game.Position().CountCandidatePawns(Black)
	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}
}

func TestDoublePawns(t *testing.T) {
	fen := "7k/1p6/8/PP3p1p/P4p1p/PP6/8/K7 w - - 0 1"
	game := FromFen(fen, true)

	expected := int16(3)
	actual := game.Position().CountDoublePawns(White)
	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}

	expected = int16(2)
	actual = game.Position().CountDoublePawns(Black)
	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}
}

func TestIsolatedPawns(t *testing.T) {
	fen := "7k/8/2p1P3/P2p1p1p/5p2/6P1/7P/K7 w - - 0 1"
	game := FromFen(fen, true)

	expected := int16(2)
	actual := game.Position().CountIsolatedPawns(White)
	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}

	expected = int16(3)
	actual = game.Position().CountIsolatedPawns(Black)
	if actual != expected {
		t.Error(fmt.Sprintf("Expected: %d\n, Got: %d\n", expected, actual))
	}
}

func TestPassedPawns(t *testing.T) {
	fen := "7k/8/1Pp1P3/3p1p1p/P4p2/1p4P1/7P/K7 w - - 0 1"
	game := FromFen(fen, true)

	expectedF, expectedS := int16(1), int16(2)
	actualF, actualS := game.Position().CountPassedPawns(White)
	if actualF != expectedF || actualS != expectedS {
		t.Error(fmt.Sprintf("Expected: (%d, %d)\n, Got: (%d, %d)\n", expectedF, expectedS, actualF, actualS))
	}
	expectedF, expectedS = int16(2), int16(1)
	actualF, actualS = game.Position().CountPassedPawns(Black)
	if actualF != expectedF || actualS != expectedS {
		t.Error(fmt.Sprintf("Expected: (%d, %d)\n, Got: (%d, %d)\n", expectedF, expectedS, actualF, actualS))
	}
}

func TestCountKnightOutpostsWhiteSixthRank(t *testing.T) {
	fen := "r4rk1/p4ppp/1pNp2n1/1Pp5/4P3/8/1PP2PPP/2KRR3 w - - 0 1"
	game := FromFen(fen, true)

	expected := int16(1)
	actual := game.Position().CountKnightOutposts(White)
	if actual != expected {
		t.Error(fmt.Sprintf("White\nExpected: %d\n, Got: %d\n", expected, actual))
	}

	expected = int16(0)
	actual = game.Position().CountKnightOutposts(Black)
	if actual != expected {
		t.Error(fmt.Sprintf("Black\nExpected: %d\n, Got: %d\n", expected, actual))
	}
}

func TestCountKnightOutpostsWhiteFifthRank(t *testing.T) {
	fen := "r4rk1/p4ppp/1p1p2n1/1PpN4/4P3/8/1PP2PPP/2KRR3 w - - 0 1"
	game := FromFen(fen, true)

	expected := int16(1)
	actual := game.Position().CountKnightOutposts(White)
	if actual != expected {
		t.Error(fmt.Sprintf("White\nExpected: %d\n, Got: %d\n", expected, actual))
	}

	expected = int16(0)
	actual = game.Position().CountKnightOutposts(Black)
	if actual != expected {
		t.Error(fmt.Sprintf("Black\nExpected: %d\n, Got: %d\n", expected, actual))
	}
}

func TestCountKnightOutpostsBlackThirdRank(t *testing.T) {
	fen := "r4rk1/p5pp/1p1p4/1P2p3/1p2P1P1/1Pn2N2/2P2P1P/2KRR3 w - - 0 1"
	game := FromFen(fen, true)

	expected := int16(0)
	actual := game.Position().CountKnightOutposts(White)
	if actual != expected {
		t.Error(fmt.Sprintf("White\nExpected: %d\n, Got: %d\n", expected, actual))
	}

	expected = int16(1)
	actual = game.Position().CountKnightOutposts(Black)
	if actual != expected {
		t.Error(fmt.Sprintf("Black\nExpected: %d\n, Got: %d\n", expected, actual))
	}
}

func TestCountKnightOutpostsBlackFourthRank(t *testing.T) {
	fen := "r4rk1/p5pp/1p1p4/1Pp1p3/4PnP1/2N5/1PP2P1P/2KRR3 w - - 0 1"
	game := FromFen(fen, true)

	expected := int16(0)
	actual := game.Position().CountKnightOutposts(White)
	if actual != expected {
		t.Error(fmt.Sprintf("White\nExpected: %d\n, Got: %d\n", expected, actual))
	}

	expected = int16(1)
	actual = game.Position().CountKnightOutposts(Black)
	if actual != expected {
		t.Error(fmt.Sprintf("Black\nExpected: %d\n, Got: %d\n", expected, actual))
	}
}
