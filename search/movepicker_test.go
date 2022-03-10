package search

import (
	"fmt"
	"testing"
	"time"

	. "github.com/amanjpro/zahak/engine"
)

var mp = EmptyMovePicker()

func TestMovepickerNextWithQuietHashmove(t *testing.T) {
	mp := &MovePicker{
		nil,
		NewEngine(nil),
		10,
		MoveList{
			Moves:    []Move{10, 5, 4, 8, 3, 2, 1, 6, 7, 9},
			Scores:   []int32{10000, 500, 400, 800, 300, 200, 100, 600, 700, 900},
			IsScored: true,
			Size:     10,
			Next:     1,
		},
		MoveList{
			Moves:    []Move{20, 15, 14, 18, 13, 12, 11, 16, 17, 19},
			Scores:   []int32{2000, 1500, 1400, 1800, 1300, 1200, 1100, 1600, 1700, -1900},
			IsScored: true,
			Size:     10,
			Next:     0,
		},
		make([]int32, 10),
		1,
		EmptyMove,
		true,
		false,
		EmptyMove,
		EmptyMove,
		1,
		EmptyMove,
		0,
	}

	expectedOrder := []Move{10, 20, 18, 17, 16, 15, 14, 13, 12, 11, 9, 8, 7, 6, 5, 4, 3, 2, 1, 19}

	for i := 0; i < 20; i++ {
		actual := mp.Next()
		expected := expectedOrder[i]
		if actual != Move(expected) {
			t.Error(fmt.Sprintf("Expected %d But got %d\n", expected, int32(actual)))
		}
	}
}

func TestMovepickerNextWithCaptureHashmove(t *testing.T) {
	capture := NewMove(E1, E2, WhitePawn, WhiteKing, NoType, Capture)
	mp := &MovePicker{
		nil,
		NewEngine(nil),
		capture,
		MoveList{
			Moves:    []Move{10, 5, 4, 8, 3, 2, 1, 6, 7, 9},
			Scores:   []int32{1000, 500, 400, 800, 300, 200, 100, 600, 700, 900},
			IsScored: true,
			Size:     10,
			Next:     0,
		},
		MoveList{
			Moves:    []Move{capture, 20, 15, 14, 13, 12, 11, 16, 17, 19},
			Scores:   []int32{18000, 2000, 1500, 1400, 1300, 1200, 1100, 1600, 1700, -1900},
			IsScored: true,
			Size:     10,
			Next:     1,
		},
		make([]int32, 10),
		1,
		EmptyMove,
		true,
		false,
		EmptyMove,
		EmptyMove,
		1,
		EmptyMove,
		0,
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
}

func TestMovepickerNextWithNoHashmove(t *testing.T) {
	mp := &MovePicker{
		nil,
		NewEngine(nil),
		0,
		MoveList{
			Moves:    []Move{10, 5, 4, 8, 3, 2, 1, 6, 7, 9},
			Scores:   []int32{1000, 500, 400, 800, 300, 200, 100, 600, 700, 900},
			IsScored: true,
			Size:     10,
			Next:     0,
		},
		MoveList{
			Moves:    []Move{20, 15, 14, 18, 13, 12, 11, 16, 17, 19},
			Scores:   []int32{2000, 1500, 1400, 1800, 1300, 1200, 1100, 1600, 1700, -1900},
			IsScored: true,
			Size:     10,
			Next:     0,
		},
		make([]int32, 10),
		1,
		EmptyMove,
		false,
		false,
		EmptyMove,
		EmptyMove,
		1,
		EmptyMove,
		0,
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
}

func TestMovePickerNormalSearch(t *testing.T) {
	fen := "rnbqkb1r/ppp2ppp/5n2/3p4/4P3/2N1P3/PPP2PPP/R1BQKBNR w KQkq - 1 2"

	game := FromFen(fen)
	r := NewRunner(1)
	tm, _, _ := NewTimeManager(time.Now(), 400_000, true, 0, 0, false)
	r.TimeManager = tm
	engine := NewEngine(r)
	engine.Position = game.Position()
	engine.ClearForSearch()

	engine.searchHistory.addKillerForTest(NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0), 1)
	engine.searchHistory.addKillerForTest(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1)
	engine.searchHistory.addHistoryForTest(NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0), 1)
	engine.searchHistory.addHistoryForTest(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1) // this is no-op
	mp.RecycleWith(game.Position(), engine, 1, NewMove(A1, B1, WhiteRook, NoPiece, NoType, 0), 0, false)

	moves := []Move{
		NewMove(A1, B1, WhiteRook, NoPiece, NoType, 0),
		NewMove(E4, D5, WhitePawn, BlackPawn, NoType, Capture),
		NewMove(C3, D5, WhiteKnight, BlackPawn, NoType, Capture),
		NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0),
		NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F2, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F3, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G4, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G3, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H4, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H3, WhitePawn, NoPiece, NoType, 0),
		NewMove(E4, E5, WhitePawn, NoPiece, NoType, 0),
		NewMove(G1, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, F3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, H3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B1, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, A4, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B5, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C1, D2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, E2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, D3, WhiteBishop, NoPiece, NoType, 0),
		NewMove(A2, A3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F1, B5, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, A6, WhiteBishop, NoPiece, NoType, 0),
		NewMove(A2, A4, WhitePawn, NoPiece, NoType, 0),
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
}

func TestUpgradeMoveToHashmoveQuiet(t *testing.T) {
	fen := "rnbqkb1r/ppp2ppp/5n2/3p4/4P3/2N1P3/PPP2PPP/R1BQKBNR w KQkq - 1 2"

	game := FromFen(fen)
	r := NewRunner(1)
	tm, _, _ := NewTimeManager(time.Now(), 400_000, true, 0, 0, false)
	r.TimeManager = tm
	engine := NewEngine(r)

	engine.Position = game.Position()
	engine.ClearForSearch()

	engine.searchHistory.addKillerForTest(NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0), 1)
	engine.searchHistory.addKillerForTest(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1)
	engine.searchHistory.addHistoryForTest(NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0), 1)
	engine.searchHistory.addHistoryForTest(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1) // this is no-op

	mp.RecycleWith(game.Position(), engine, 1, EmptyMove, 0, false)

	moves := []Move{
		NewMove(A1, B1, WhiteRook, NoPiece, NoType, 0),
		NewMove(E4, D5, WhitePawn, BlackPawn, NoType, Capture),
		NewMove(C3, D5, WhiteKnight, BlackPawn, NoType, Capture),
		NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0),
		NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F2, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F3, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G4, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G3, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H4, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H3, WhitePawn, NoPiece, NoType, 0),
		NewMove(E4, E5, WhitePawn, NoPiece, NoType, 0),
		NewMove(G1, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, F3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, H3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B1, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, A4, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B5, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C1, D2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, E2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, D3, WhiteBishop, NoPiece, NoType, 0),
		NewMove(A2, A3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F1, B5, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, A6, WhiteBishop, NoPiece, NoType, 0),
		NewMove(A2, A4, WhitePawn, NoPiece, NoType, 0),
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
}

func TestMovePickerNormalSearchNoHashmove(t *testing.T) {
	fen := "rnbqkb1r/ppp2ppp/5n2/3p4/4P3/2N1P3/PPP2PPP/R1BQKBNR w KQkq - 1 2"

	game := FromFen(fen)
	r := NewRunner(1)
	tm, _, _ := NewTimeManager(time.Now(), 400_000, true, 0, 0, false)
	r.TimeManager = tm
	engine := NewEngine(r)
	engine.Position = game.Position()
	engine.ClearForSearch()

	engine.searchHistory.addKillerForTest(NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0), 1)
	engine.searchHistory.addKillerForTest(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1)
	engine.searchHistory.addHistoryForTest(NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0), 1)
	engine.searchHistory.addHistoryForTest(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1) // this is no-op
	mp.RecycleWith(game.Position(), engine, 1, EmptyMove, 0, false)

	moves := []Move{
		NewMove(E4, D5, WhitePawn, BlackPawn, NoType, Capture),
		NewMove(C3, D5, WhiteKnight, BlackPawn, NoType, Capture),
		NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0),
		NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0),
		NewMove(A2, A3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F3, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G4, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G3, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H4, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H3, WhitePawn, NoPiece, NoType, 0),
		NewMove(E4, E5, WhitePawn, NoPiece, NoType, 0),
		NewMove(G1, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, F3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, H3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B1, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, A4, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B5, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C1, D2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, E2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, D3, WhiteBishop, NoPiece, NoType, 0),
		NewMove(A2, A4, WhitePawn, NoPiece, NoType, 0),
		NewMove(F1, B5, WhiteBishop, NoPiece, NoType, 0),
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
}

func TestMovePickerNormalSearchCaptureHashmove(t *testing.T) {
	fen := "rnbqkb1r/ppp2ppp/5n2/3p4/4P3/2N1P3/PPP2PPP/R1BQKBNR w KQkq - 1 2"

	game := FromFen(fen)
	r := NewRunner(1)
	tm, _, _ := NewTimeManager(time.Now(), 400_000, true, 0, 0, false)
	r.TimeManager = tm
	engine := NewEngine(r)
	engine.Position = game.Position()
	engine.ClearForSearch()

	engine.searchHistory.addKillerForTest(NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0), 1)
	engine.searchHistory.addKillerForTest(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1)
	engine.searchHistory.addHistoryForTest(NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0), 1)
	engine.searchHistory.addHistoryForTest(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1) // this is no-op
	mp.RecycleWith(game.Position(), engine, 1, NewMove(C3, D5, WhiteKnight, BlackPawn, NoType, Capture), 0, false)

	moves := []Move{
		NewMove(C3, D5, WhiteKnight, BlackPawn, NoType, Capture),
		NewMove(E4, D5, WhitePawn, BlackPawn, NoType, Capture),
		NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0),
		NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0),
		NewMove(A2, A3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F3, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G4, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G3, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H4, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H3, WhitePawn, NoPiece, NoType, 0),
		NewMove(E4, E5, WhitePawn, NoPiece, NoType, 0),
		NewMove(G1, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, F3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, H3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B1, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, A4, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B5, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C1, D2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, E2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, D3, WhiteBishop, NoPiece, NoType, 0),
		NewMove(A2, A4, WhitePawn, NoPiece, NoType, 0),
		NewMove(F1, B5, WhiteBishop, NoPiece, NoType, 0),
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
}

func TestMovePickerNormalSearchUpgradeToHashmoveCapture(t *testing.T) {
	fen := "rnbqkb1r/ppp2ppp/5n2/3p4/4P3/2N1P3/PPP2PPP/R1BQKBNR w KQkq - 1 2"

	game := FromFen(fen)
	r := NewRunner(1)
	tm, _, _ := NewTimeManager(time.Now(), 400_000, true, 0, 0, false)
	r.TimeManager = tm
	engine := NewEngine(r)
	engine.Position = game.Position()
	engine.ClearForSearch()

	engine.searchHistory.addKillerForTest(NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0), 1)
	engine.searchHistory.addKillerForTest(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1)
	engine.searchHistory.addHistoryForTest(NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0), 1)
	engine.searchHistory.addHistoryForTest(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1) // this is no-op
	mp.RecycleWith(game.Position(), engine, 1, EmptyMove, 0, false)
	mp.UpgradeToPvMove(NewMove(C3, D5, WhiteKnight, BlackPawn, NoType, Capture))

	moves := []Move{
		NewMove(C3, D5, WhiteKnight, BlackPawn, NoType, Capture),
		NewMove(E4, D5, WhitePawn, BlackPawn, NoType, Capture),
		NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0),
		NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0),
		NewMove(A2, A3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F3, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G4, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G3, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H4, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H3, WhitePawn, NoPiece, NoType, 0),
		NewMove(E4, E5, WhitePawn, NoPiece, NoType, 0),
		NewMove(G1, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, F3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G1, H3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B1, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, E2, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, A4, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C3, B5, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C1, D2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, E2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, D3, WhiteBishop, NoPiece, NoType, 0),
		NewMove(A2, A4, WhitePawn, NoPiece, NoType, 0),
		NewMove(F1, B5, WhiteBishop, NoPiece, NoType, 0),
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
}

func TestMovePickerQuiescenceSearch(t *testing.T) {
	fen := "rnbqkb1r/ppp2ppp/5n2/3p4/4P3/2N1P3/PPP2PPP/R1BQKBNR w KQkq - 1 2"

	game := FromFen(fen)
	r := NewRunner(1)
	tm, _, _ := NewTimeManager(time.Now(), 400_000, true, 0, 0, false)
	r.TimeManager = tm
	engine := NewEngine(r)
	engine.Position = game.Position()
	engine.ClearForSearch()

	// all these are no-op
	engine.searchHistory.addKillerForTest(NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0), 1)
	engine.searchHistory.addKillerForTest(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1)
	engine.searchHistory.addHistoryForTest(NewMove(F1, C4, WhiteBishop, NoPiece, NoType, 0), 1)
	engine.searchHistory.addHistoryForTest(NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0), 1)
	mp.RecycleWith(game.Position(), engine, 1, EmptyMove, 0, true)

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
}

func TestMovePickerNormalSearchWithPromotionNoHashmove(t *testing.T) {
	fen := "1k4n1/7P/8/6K1/8/5P2/8/8 w - - 0 1"

	game := FromFen(fen)
	r := NewRunner(1)
	tm, _, _ := NewTimeManager(time.Now(), 400_000, true, 0, 0, false)
	r.TimeManager = tm
	engine := NewEngine(r)
	engine.Position = game.Position()
	engine.ClearForSearch()
	mp.RecycleWith(game.Position(), engine, 1, EmptyMove, 0, false)

	moves := []Move{
		NewMove(H7, G8, WhitePawn, BlackKnight, Queen, Capture),
		NewMove(H7, H8, WhitePawn, NoPiece, Queen, 0),
		NewMove(H7, G8, WhitePawn, BlackKnight, Rook, Capture),
		NewMove(H7, G8, WhitePawn, BlackKnight, Bishop, Capture),
		NewMove(H7, G8, WhitePawn, BlackKnight, Knight, Capture),
		NewMove(H7, H8, WhitePawn, NoPiece, Rook, 0),
		NewMove(H7, H8, WhitePawn, NoPiece, Bishop, 0),
		NewMove(H7, H8, WhitePawn, NoPiece, Knight, 0),
		NewMove(F3, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(G5, F4, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, G4, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, H4, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, F5, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, H5, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, F6, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, G6, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, H6, WhiteKing, NoPiece, NoType, 0),
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
}

func TestMovePickerNormalSearchWithPromotionPromotionQuietHashmove(t *testing.T) {
	fen := "1k4n1/7P/8/6K1/8/5P2/8/8 w - - 0 1"

	game := FromFen(fen)
	r := NewRunner(1)
	tm, _, _ := NewTimeManager(time.Now(), 400_000, true, 0, 0, false)
	r.TimeManager = tm
	engine := NewEngine(r)
	engine.Position = game.Position()
	engine.ClearForSearch()
	mp.RecycleWith(game.Position(), engine, 1, NewMove(H7, H8, WhitePawn, NoPiece, Knight, 0), 0, false)

	moves := []Move{
		NewMove(H7, H8, WhitePawn, NoPiece, Knight, 0),
		NewMove(H7, G8, WhitePawn, BlackKnight, Queen, Capture),
		NewMove(H7, H8, WhitePawn, NoPiece, Queen, 0),
		NewMove(H7, G8, WhitePawn, BlackKnight, Rook, Capture),
		NewMove(H7, G8, WhitePawn, BlackKnight, Bishop, Capture),
		NewMove(H7, G8, WhitePawn, BlackKnight, Knight, Capture),
		NewMove(H7, H8, WhitePawn, NoPiece, Rook, 0),
		NewMove(H7, H8, WhitePawn, NoPiece, Bishop, 0),
		NewMove(F3, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(G5, F4, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, G4, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, H4, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, F5, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, H5, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, F6, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, G6, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, H6, WhiteKing, NoPiece, NoType, 0),
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
}

func TestMovePickerNormalSearchWithPromotionPromotionCaptureHashmove(t *testing.T) {
	fen := "1k4n1/7P/8/6K1/8/5P2/8/8 w - - 0 1"

	game := FromFen(fen)
	r := NewRunner(1)
	tm, _, _ := NewTimeManager(time.Now(), 400_000, true, 0, 0, false)
	r.TimeManager = tm
	engine := NewEngine(r)
	engine.Position = game.Position()
	engine.ClearForSearch()
	mp.RecycleWith(game.Position(), engine, 1, NewMove(H7, G8, WhitePawn, BlackKnight, Knight, Capture), 0, false)

	moves := []Move{
		NewMove(H7, G8, WhitePawn, BlackKnight, Knight, Capture),
		NewMove(H7, G8, WhitePawn, BlackKnight, Queen, Capture),
		NewMove(H7, H8, WhitePawn, NoPiece, Queen, 0),
		NewMove(H7, G8, WhitePawn, BlackKnight, Rook, Capture),
		NewMove(H7, G8, WhitePawn, BlackKnight, Bishop, Capture),
		NewMove(H7, H8, WhitePawn, NoPiece, Rook, 0),
		NewMove(H7, H8, WhitePawn, NoPiece, Bishop, 0),
		NewMove(H7, H8, WhitePawn, NoPiece, Knight, 0),
		NewMove(F3, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(G5, F4, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, G4, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, H4, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, F5, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, H5, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, F6, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, G6, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, H6, WhiteKing, NoPiece, NoType, 0),
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
}

func TestMovePickerNormalSearchWithPromotionUpgradeToPromotionQuietHashmove(t *testing.T) {
	fen := "1k4n1/7P/8/6K1/8/5P2/8/8 w - - 0 1"

	game := FromFen(fen)
	r := NewRunner(1)
	tm, _, _ := NewTimeManager(time.Now(), 400_000, true, 0, 0, false)
	r.TimeManager = tm
	engine := NewEngine(r)
	engine.Position = game.Position()
	engine.ClearForSearch()
	mp.RecycleWith(game.Position(), engine, 1, EmptyMove, 0, false)
	mp.UpgradeToPvMove(NewMove(H7, H8, WhitePawn, NoPiece, Knight, 0))

	moves := []Move{
		NewMove(H7, H8, WhitePawn, NoPiece, Knight, 0),
		NewMove(H7, G8, WhitePawn, BlackKnight, Queen, Capture),
		NewMove(H7, H8, WhitePawn, NoPiece, Queen, 0),
		NewMove(H7, G8, WhitePawn, BlackKnight, Rook, Capture),
		NewMove(H7, G8, WhitePawn, BlackKnight, Bishop, Capture),
		NewMove(H7, G8, WhitePawn, BlackKnight, Knight, Capture),
		NewMove(H7, H8, WhitePawn, NoPiece, Rook, 0),
		NewMove(H7, H8, WhitePawn, NoPiece, Bishop, 0),
		NewMove(F3, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(G5, F4, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, G4, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, H4, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, F5, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, H5, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, F6, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, G6, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, H6, WhiteKing, NoPiece, NoType, 0),
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
}

func TestMovePickerNormalSearchWithPromotionUpgradeToPromotionCaptureHashmove(t *testing.T) {
	fen := "1k4n1/7P/8/6K1/8/5P2/8/8 w - - 0 1"

	game := FromFen(fen)
	r := NewRunner(1)
	tm, _, _ := NewTimeManager(time.Now(), 400_000, true, 0, 0, false)
	r.TimeManager = tm
	engine := NewEngine(r)
	engine.Position = game.Position()
	engine.ClearForSearch()
	mp.RecycleWith(game.Position(), engine, 1, EmptyMove, 0, false)
	mp.UpgradeToPvMove(NewMove(H7, G8, WhitePawn, BlackKnight, Knight, Capture))

	moves := []Move{
		NewMove(H7, G8, WhitePawn, BlackKnight, Knight, Capture),
		NewMove(H7, G8, WhitePawn, BlackKnight, Queen, Capture),
		NewMove(H7, H8, WhitePawn, NoPiece, Queen, 0),
		NewMove(H7, G8, WhitePawn, BlackKnight, Rook, Capture),
		NewMove(H7, G8, WhitePawn, BlackKnight, Bishop, Capture),
		NewMove(H7, H8, WhitePawn, NoPiece, Rook, 0),
		NewMove(H7, H8, WhitePawn, NoPiece, Bishop, 0),
		NewMove(H7, H8, WhitePawn, NoPiece, Knight, 0),
		NewMove(F3, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(G5, F4, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, G4, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, H4, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, F5, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, H5, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, F6, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, G6, WhiteKing, NoPiece, NoType, 0),
		NewMove(G5, H6, WhiteKing, NoPiece, NoType, 0),
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
}

func (mp *MovePicker) getScore(m Move) int32 {
	if mp.hashmove == m {
		return 900_000_000
	}
	if m.IsCapture() {
		for i, s := range mp.captureMoveList.Scores {
			if mp.captureMoveList.Moves[i] == m {
				return s
			}
		}
	}
	for i, s := range mp.quietMoveList.Scores {
		if mp.quietMoveList.Moves[i] == m {
			return s
		}
	}
	return -900_000_000
}

func (m *MoveHistory) addKillerForTest(move Move, h int8) {
	m.killers[h][0], m.killers[h][1] = move, m.killers[h][0]
}

func (m *MoveHistory) addHistoryForTest(move Move, depthLeft int8) {
	m.history[move.MovingPiece()-1][move.Destination()] = int32(depthLeft) * int32(depthLeft)
}
