package search

import (
	"fmt"
	"testing"

	. "github.com/amanjpro/zahak/engine"
)

func TestMovepickerNextAndResetWithQuietHashmove(t *testing.T) {
	mp := &MovePicker{
		nil,
		nil,
		10,
		[]Move{10, 5, 4, 8, 3, 2, 1, 6, 7, 9},
		[]int32{10000, 500, 400, 800, 300, 200, 100, 600, 700, 900},
		[]Move{20, 15, 14, 18, 13, 12, 11, 16, 17, 19},
		[]int32{2000, 1500, 1400, 1800, 1300, 1200, 1100, 1600, 1700, -1900},
		0,
		1,
		0,
		true,
		false,
	}

	expectedOrder := []Move{10, 20, 18, 17, 16, 15, 14, 13, 12, 11, 9, 8, 7, 6, 5, 4, 3, 2, 1, 19}

	for i := 0; i < 20; i++ {
		actual := mp.Next()
		expected := expectedOrder[i]
		if actual != Move(expected) {
			t.Error(fmt.Sprintf("Expected %d But got %d\n", expected, int32(actual)))
		}
	}

	mp.Reset()

	for i := 0; i < 20; i++ {
		actual := mp.Next()
		expected := expectedOrder[i]
		if actual != Move(expected) {
			t.Error(fmt.Sprintf("MovePicker Reset is broken.\nExpected %d But got %d\n", expected, int32(actual)))
		}
	}
}

func TestMovepickerNextAndResetWithCaptureHashmove(t *testing.T) {
	capture := NewMove(E1, E2, WhitePawn, WhiteKing, NoType, Capture)
	mp := &MovePicker{
		nil,
		nil,
		capture,
		[]Move{10, 5, 4, 8, 3, 2, 1, 6, 7, 9},
		[]int32{1000, 500, 400, 800, 300, 200, 100, 600, 700, 900},
		[]Move{capture, 20, 15, 14, 13, 12, 11, 16, 17, 19},
		[]int32{18000, 2000, 1500, 1400, 1300, 1200, 1100, 1600, 1700, -1900},
		0,
		0,
		1,
		true,
		false,
	}

	expectedOrder := []Move{capture, 20, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 19}
	i := 0
	for ; ; i++ {
		actual := mp.Next()
		if actual == EmptyMove {
			break
		}
		expected := expectedOrder[i]
		if actual != Move(expected) {
			t.Error(fmt.Sprintf("Expected %d But got %d\n", expected, int32(actual)))
		}
	}

	if i != 20 {
		t.Error("Wrong number of moves!")
	}

	mp.Reset()

	for i = 0; ; i++ {
		actual := mp.Next()
		if actual == EmptyMove {
			break
		}
		expected := expectedOrder[i]
		if actual != Move(expected) {
			t.Error(fmt.Sprintf("MovePicker Reset is broken.\nExpected %d But got %d\n", expected, int32(actual)))
		}
	}

	if i != 20 {
		t.Error("Wrong number of moves in reset!")
	}
}

func TestMovepickerNextAndResetWithNoHashmove(t *testing.T) {
	mp := &MovePicker{
		nil,
		nil,
		0,
		[]Move{10, 5, 4, 8, 3, 2, 1, 6, 7, 9},
		[]int32{1000, 500, 400, 800, 300, 200, 100, 600, 700, 900},
		[]Move{20, 15, 14, 18, 13, 12, 11, 16, 17, 19},
		[]int32{2000, 1500, 1400, 1800, 1300, 1200, 1100, 1600, 1700, -1900},
		0,
		0,
		0,
		false,
		false,
	}

	expectedOrder := []Move{20, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 19}

	i := 0
	for ; ; i++ {
		actual := mp.Next()
		if actual == EmptyMove {
			break
		}
		expected := expectedOrder[i]
		if actual != Move(expected) {
			t.Error(fmt.Sprintf("Expected %d But got %d\n", expected, int32(actual)))
		}
	}

	if i != 20 {
		t.Error("Wrong number of moves!")
	}
	mp.Reset()

	for i = 0; ; i++ {
		actual := mp.Next()
		if actual == EmptyMove {
			break
		}
		expected := expectedOrder[i]
		if actual != Move(expected) {
			t.Error(fmt.Sprintf("MovePicker Reset is broken.\nExpected %d But got %d\n", expected, int32(actual)))
		}
	}

	if i != 20 {
		t.Error("Wrong number of moves in reset!")
	}
}

func TestMovePickerNormalSearch(t *testing.T) {
	fen := "rnbqkb1r/ppp2ppp/5n2/3p4/4P3/2N1P3/PPP2PPP/R1BQKBNR w KQkq - 1 2"

	game := FromFen(fen, true)
	engine := NewEngine(NewCache(2))
	engine.ClearForSearch()
	mp := NewMovePicker(game.Position(), engine, 1, NewMove(A1, B1, WhiteRook, NoPiece, NoType, 0), false)

	engine.AddKillerMove(NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0), 1)
	engine.AddKillerMove(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1)
	engine.AddMoveHistory(NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0), WhiteBishop, C4, 1)
	engine.AddMoveHistory(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), WhitePawn, B4, 1) // this is no-op

	moves := []Move{
		NewMove(A1, B1, WhiteRook, NoPiece, NoType, 0),
		NewMove(E4, D5, WhitePawn, BlackPawn, NoType, Capture),
		NewMove(C3, D5, WhiteKnight, BlackPawn, NoType, Capture),
		NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0),
		NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, B5, WhiteBishop, NoPiece, NoType, Check),
		NewMove(F2, F3, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G4, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G3, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H4, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H3, WhitePawn, NoPiece, NoType, 0),
		NewMove(E4, E5, WhitePawn, NoPiece, NoType, 0),
		NewMove(A2, A3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(A2, A4, WhitePawn, NoPiece, NoType, 0),
		NewMove(C3, B1, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, A4, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B5, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, F3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, H3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C1, D2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, E2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, A6, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, D3, WhiteBishop, NoPiece, NoType, 0),
		NewMove(D1, D2, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, E2, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, D3, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, F3, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, D4, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, G4, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, H5, WhiteQueen, NoPiece, NoType, 0),
		NewMove(E1, D2, WhiteKing, NoPiece, NoType, 0),
		NewMove(E1, E2, WhiteKing, NoPiece, NoType, 0),
		NewMove(D1, D5, WhiteQueen, BlackPawn, NoType, Capture),
	}

	i := 0
	for ; ; i++ {
		actual := mp.Next()
		if actual == EmptyMove {
			break
		}
		expected := moves[i]
		if actual != expected {
			t.Error(fmt.Sprintf("Move number %d Expected %s But got %s which has score of %d\n", i+1, expected.ToString(), actual.ToString(), mp.getScore(actual)))
		}
	}

	if i != len(moves) {
		t.Error("Wrong number of moves!")
	}
	mp.Reset()

	for i = 0; ; i++ {
		actual := mp.Next()
		if actual == EmptyMove {
			break
		}
		expected := moves[i]
		if actual != expected {
			t.Error(fmt.Sprintf("Reset is broken.\nMove number %d Expected %s But got %s which has score of %d\n", i+1, expected.ToString(), actual.ToString(), mp.getScore(actual)))
		}
	}

	if i != len(moves) {
		t.Error("Wrong number of moves in reset!")
	}
}

func TestUpgradeMoveToHashmoveQuiet(t *testing.T) {
	fen := "rnbqkb1r/ppp2ppp/5n2/3p4/4P3/2N1P3/PPP2PPP/R1BQKBNR w KQkq - 1 2"

	game := FromFen(fen, true)
	engine := NewEngine(NewCache(2))
	engine.ClearForSearch()
	mp := NewMovePicker(game.Position(), engine, 1, EmptyMove, false)

	engine.AddKillerMove(NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0), 1)
	engine.AddKillerMove(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1)
	engine.AddMoveHistory(NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0), WhiteBishop, C4, 1)
	engine.AddMoveHistory(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), WhitePawn, B4, 1) // this is no-op

	moves := []Move{
		NewMove(A1, B1, WhiteRook, NoPiece, NoType, 0),
		NewMove(E4, D5, WhitePawn, BlackPawn, NoType, Capture),
		NewMove(C3, D5, WhiteKnight, BlackPawn, NoType, Capture),
		NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0),
		NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, B5, WhiteBishop, NoPiece, NoType, Check),
		NewMove(F2, F3, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G4, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G3, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H4, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H3, WhitePawn, NoPiece, NoType, 0),
		NewMove(E4, E5, WhitePawn, NoPiece, NoType, 0),
		NewMove(A2, A3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(A2, A4, WhitePawn, NoPiece, NoType, 0),
		NewMove(C3, B1, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, A4, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B5, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, F3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, H3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C1, D2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, E2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, A6, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, D3, WhiteBishop, NoPiece, NoType, 0),
		NewMove(D1, D2, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, E2, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, D3, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, F3, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, D4, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, G4, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, H5, WhiteQueen, NoPiece, NoType, 0),
		NewMove(E1, D2, WhiteKing, NoPiece, NoType, 0),
		NewMove(E1, E2, WhiteKing, NoPiece, NoType, 0),
		NewMove(D1, D5, WhiteQueen, BlackPawn, NoType, Capture),
	}

	mp.UpgradeToPvMove(NewMove(A1, B1, WhiteRook, NoPiece, NoType, 0))

	i := 0
	for ; ; i++ {
		actual := mp.Next()
		if actual == EmptyMove {
			break
		}
		expected := moves[i]
		if actual != expected {
			t.Error(fmt.Sprintf("Move number %d Expected %s But got %s which has score of %d\n", i+1, expected.ToString(), actual.ToString(), mp.getScore(actual)))
		}
	}

	if i != len(moves) {
		t.Error("Wrong number of moves!")
	}
	mp.Reset()

	for i = 0; ; i++ {
		actual := mp.Next()
		if actual == EmptyMove {
			break
		}
		expected := moves[i]
		if actual != expected {
			t.Error(fmt.Sprintf("Reset is broken.\nMove number %d Expected %s But got %s which has score of %d\n", i+1, expected.ToString(), actual.ToString(), mp.getScore(actual)))
		}
	}

	if i != len(moves) {
		t.Error("Wrong number of moves in reset!")
	}
}

func TestMovePickerNormalSearchNoHashmove(t *testing.T) {
	fen := "rnbqkb1r/ppp2ppp/5n2/3p4/4P3/2N1P3/PPP2PPP/R1BQKBNR w KQkq - 1 2"

	game := FromFen(fen, true)
	engine := NewEngine(NewCache(2))
	engine.ClearForSearch()
	mp := NewMovePicker(game.Position(), engine, 1, EmptyMove, false)

	engine.AddKillerMove(NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0), 1)
	engine.AddKillerMove(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1)
	engine.AddMoveHistory(NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0), WhiteBishop, C4, 1)
	engine.AddMoveHistory(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), WhitePawn, B4, 1) // this is no-op

	moves := []Move{
		NewMove(E4, D5, WhitePawn, BlackPawn, NoType, Capture),
		NewMove(C3, D5, WhiteKnight, BlackPawn, NoType, Capture),
		NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0),
		NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, B5, WhiteBishop, NoPiece, NoType, Check),
		NewMove(F2, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F3, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G4, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G3, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H4, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H3, WhitePawn, NoPiece, NoType, 0),
		NewMove(E4, E5, WhitePawn, NoPiece, NoType, 0),
		NewMove(A2, A4, WhitePawn, NoPiece, NoType, 0),
		NewMove(A2, A3, WhitePawn, NoPiece, NoType, 0),
		NewMove(G1, H3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B1, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, A4, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B5, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, F3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(F1, D3, WhiteBishop, NoPiece, NoType, 0),
		NewMove(C1, D2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, E2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, A6, WhiteBishop, NoPiece, NoType, 0),
		NewMove(A1, B1, WhiteRook, NoPiece, NoType, 0),
		NewMove(D1, D2, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, E2, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, D3, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, F3, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, D4, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, G4, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, H5, WhiteQueen, NoPiece, NoType, 0),
		NewMove(E1, D2, WhiteKing, NoPiece, NoType, 0),
		NewMove(E1, E2, WhiteKing, NoPiece, NoType, 0),
		NewMove(D1, D5, WhiteQueen, BlackPawn, NoType, Capture),
	}

	i := 0
	for ; ; i++ {
		actual := mp.Next()
		if actual == EmptyMove {
			break
		}
		expected := moves[i]
		if actual != expected {
			t.Error(fmt.Sprintf("Move number %d Expected %s But got %s which has score of %d\n", i+1, expected.ToString(), actual.ToString(), mp.getScore(actual)))
		}
	}

	if i != len(moves) {
		t.Error("Wrong number of moves!")
	}
	mp.Reset()

	for i = 0; ; i++ {
		actual := mp.Next()
		if actual == EmptyMove {
			break
		}
		expected := moves[i]
		if actual != expected {
			t.Error(fmt.Sprintf("Reset is broken.\nMove number %d Expected %s But got %s which has score of %d\n", i+1, expected.ToString(), actual.ToString(), mp.getScore(actual)))
		}
	}

	if i != len(moves) {
		t.Error("Wrong number of moves in reset!")
	}
}

func TestMovePickerNormalSearchCaptureHashmove(t *testing.T) {
	fen := "rnbqkb1r/ppp2ppp/5n2/3p4/4P3/2N1P3/PPP2PPP/R1BQKBNR w KQkq - 1 2"

	game := FromFen(fen, true)
	engine := NewEngine(NewCache(2))
	engine.ClearForSearch()
	mp := NewMovePicker(game.Position(), engine, 1, NewMove(C3, D5, WhiteKnight, BlackPawn, NoType, Capture), false)

	engine.AddKillerMove(NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0), 1)
	engine.AddKillerMove(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1)
	engine.AddMoveHistory(NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0), WhiteBishop, C4, 1)
	engine.AddMoveHistory(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), WhitePawn, B4, 1) // this is no-op

	moves := []Move{
		NewMove(C3, D5, WhiteKnight, BlackPawn, NoType, Capture),
		NewMove(E4, D5, WhitePawn, BlackPawn, NoType, Capture),
		NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0),
		NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, B5, WhiteBishop, NoPiece, NoType, Check),
		NewMove(F2, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F3, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G4, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G3, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H4, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H3, WhitePawn, NoPiece, NoType, 0),
		NewMove(E4, E5, WhitePawn, NoPiece, NoType, 0),
		NewMove(A2, A4, WhitePawn, NoPiece, NoType, 0),
		NewMove(A2, A3, WhitePawn, NoPiece, NoType, 0),
		NewMove(G1, H3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B1, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, A4, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B5, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, F3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(F1, D3, WhiteBishop, NoPiece, NoType, 0),
		NewMove(C1, D2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, E2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, A6, WhiteBishop, NoPiece, NoType, 0),
		NewMove(A1, B1, WhiteRook, NoPiece, NoType, 0),
		NewMove(D1, D2, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, E2, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, D3, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, F3, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, D4, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, G4, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, H5, WhiteQueen, NoPiece, NoType, 0),
		NewMove(E1, D2, WhiteKing, NoPiece, NoType, 0),
		NewMove(E1, E2, WhiteKing, NoPiece, NoType, 0),
		NewMove(D1, D5, WhiteQueen, BlackPawn, NoType, Capture),
	}

	i := 0

	for ; ; i++ {
		actual := mp.Next()
		if actual == EmptyMove {
			break
		}
		expected := moves[i]
		if actual != expected {
			t.Error(fmt.Sprintf("Move number %d Expected %s But got %s which has score of %d\n", i+1, expected.ToString(), actual.ToString(), mp.getScore(actual)))
		}
	}

	if i != len(moves) {
		t.Error("Wrong number of moves!")
	}
	mp.Reset()

	for i = 0; ; i++ {
		actual := mp.Next()
		if actual == EmptyMove {
			break
		}
		expected := moves[i]
		if actual != expected {
			t.Error(fmt.Sprintf("Reset is broken.\nMove number %d Expected %s But got %s which has score of %d\n", i+1, expected.ToString(), actual.ToString(), mp.getScore(actual)))
		}
	}

	if i != len(moves) {
		t.Error("Wrong number of moves in reset!")
	}
}

func TestMovePickerNormalSearchUpgradeToHashmoveCapture(t *testing.T) {
	fen := "rnbqkb1r/ppp2ppp/5n2/3p4/4P3/2N1P3/PPP2PPP/R1BQKBNR w KQkq - 1 2"

	game := FromFen(fen, true)
	engine := NewEngine(NewCache(2))
	engine.ClearForSearch()
	mp := NewMovePicker(game.Position(), engine, 1, EmptyMove, false)

	mp.UpgradeToPvMove(NewMove(C3, D5, WhiteKnight, BlackPawn, NoType, Capture))

	engine.AddKillerMove(NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0), 1)
	engine.AddKillerMove(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1)
	engine.AddMoveHistory(NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0), WhiteBishop, C4, 1)
	engine.AddMoveHistory(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), WhitePawn, B4, 1) // this is no-op

	moves := []Move{
		NewMove(C3, D5, WhiteKnight, BlackPawn, NoType, Capture),
		NewMove(E4, D5, WhitePawn, BlackPawn, NoType, Capture),
		NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0),
		NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, B5, WhiteBishop, NoPiece, NoType, Check),
		NewMove(F2, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F3, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G4, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G3, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H4, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H3, WhitePawn, NoPiece, NoType, 0),
		NewMove(E4, E5, WhitePawn, NoPiece, NoType, 0),
		NewMove(A2, A4, WhitePawn, NoPiece, NoType, 0),
		NewMove(A2, A3, WhitePawn, NoPiece, NoType, 0),
		NewMove(G1, H3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B1, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, A4, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B5, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, F3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(F1, D3, WhiteBishop, NoPiece, NoType, 0),
		NewMove(C1, D2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, E2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, A6, WhiteBishop, NoPiece, NoType, 0),
		NewMove(A1, B1, WhiteRook, NoPiece, NoType, 0),
		NewMove(D1, D2, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, E2, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, D3, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, F3, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, D4, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, G4, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, H5, WhiteQueen, NoPiece, NoType, 0),
		NewMove(E1, D2, WhiteKing, NoPiece, NoType, 0),
		NewMove(E1, E2, WhiteKing, NoPiece, NoType, 0),
		NewMove(D1, D5, WhiteQueen, BlackPawn, NoType, Capture),
	}

	i := 0

	for ; ; i++ {
		actual := mp.Next()
		if actual == EmptyMove {
			break
		}
		expected := moves[i]
		if actual != expected {
			t.Error(fmt.Sprintf("Move number %d Expected %s But got %s which has score of %d\n", i+1, expected.ToString(), actual.ToString(), mp.getScore(actual)))
		}
	}

	if i != len(moves) {
		t.Error("Wrong number of moves!")
	}
	mp.Reset()

	for i = 0; ; i++ {
		actual := mp.Next()
		if actual == EmptyMove {
			break
		}
		expected := moves[i]
		if actual != expected {
			t.Error(fmt.Sprintf("Reset is broken.\nMove number %d Expected %s But got %s which has score of %d\n", i+1, expected.ToString(), actual.ToString(), mp.getScore(actual)))
		}
	}

	if i != len(moves) {
		t.Error("Wrong number of moves in reset!")
	}
}

func TestMovePickerQuiescenceSearch(t *testing.T) {
	fen := "rnbqkb1r/ppp2ppp/5n2/3p4/4P3/2N1P3/PPP2PPP/R1BQKBNR w KQkq - 1 2"

	game := FromFen(fen, true)
	engine := NewEngine(NewCache(2))
	engine.ClearForSearch()
	mp := NewMovePicker(game.Position(), engine, 1, EmptyMove, true)

	// all these are no-op
	engine.AddKillerMove(NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0), 1)
	engine.AddKillerMove(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1)
	engine.AddMoveHistory(NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0), WhiteBishop, C4, 1)
	engine.AddMoveHistory(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), WhitePawn, B4, 1)

	moves := []Move{
		NewMove(E4, D5, WhitePawn, BlackPawn, NoType, Capture),
		NewMove(C3, D5, WhiteKnight, BlackPawn, NoType, Capture),
		NewMove(D1, D5, WhiteQueen, BlackPawn, NoType, Capture),
	}

	i := 0
	for ; ; i++ {
		actual := mp.Next()
		if actual == EmptyMove {
			break
		}
		expected := moves[i]
		if actual != expected {
			t.Error(fmt.Sprintf("Move number %d Expected %s But got %s which has score of %d\n", i+1, expected.ToString(), actual.ToString(), mp.getScore(actual)))
		}
	}

	if i != len(moves) {
		t.Error("Wrong number of moves!")
	}
	mp.Reset()

	for i = 0; ; i++ {
		actual := mp.Next()
		if actual == EmptyMove {
			break
		}
		expected := moves[i]
		if actual != expected {
			t.Error(fmt.Sprintf("Reset is broken.\nMove number %d Expected %s But got %s which has score of %d\n", i+1, expected.ToString(), actual.ToString(), mp.getScore(actual)))
		}
	}

	if i != len(moves) {
		t.Error("Wrong number of moves in reset!")
	}
}

func (mp *MovePicker) getScore(m Move) int32 {
	if mp.hashmove == m {
		return 900_000_000
	}
	if m.IsCapture() {
		for i, s := range mp.captureScores {
			if mp.captureMoves[i] == m {
				return s
			}
		}
	}
	for i, s := range mp.quietScores {
		if mp.quietMoves[i] == m {
			return s
		}
	}
	return -900_000_000
}
