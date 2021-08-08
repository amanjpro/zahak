package search

import (
	"fmt"
	"testing"
	"time"

	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func TestBlackShouldFindEscape(t *testing.T) {
	game := FromFen("3rbbn1/BQ1kp3/2p1q2p/N4p2/8/3P4/P1P2PPP/5RK1 b - - 0 27")
	r := NewRunner(NewCache(DEFAULT_CACHE_SIZE), NewPawnCache(DEFAULT_PAWNHASH_SIZE), 1)
	r.AddTimeManager(NewTimeManager(time.Now(), 400_000, true, 0, 0, false))
	e := r.Engines[0]
	e.Position = game.Position()
	e.Ply = 27
	e.Search(7)
	expected := NewMove(D7, D6, BlackKing, NoPiece, NoType, 0)
	mv := e.Move()
	mvStr := mv.ToString()
	if mv != expected {
		t.Errorf("Unexpected move was played:%s\n", fmt.Sprintf("Expected: %s\nGot: %s\n", expected.ToString(), mvStr))
	}
}

func TestBlackCanFindASimpleTactic(t *testing.T) {
	game := FromFen("3N1k2/N7/1p2ppR1/1P6/P2pP3/3Pb3/2r4n/3K4 b - - 0 1")
	r := NewRunner(NewCache(DEFAULT_CACHE_SIZE), NewPawnCache(DEFAULT_PAWNHASH_SIZE), 1)
	r.AddTimeManager(NewTimeManager(time.Now(), 400_000, true, 0, 0, false))
	e := r.Engines[0]
	e.Position = game.Position()
	e.Ply = 1
	e.Search(7)
	expected := NewMove(C2, D2, BlackRook, NoPiece, NoType, 0)
	mv := e.Move()
	mvStr := mv.ToString()
	if mv != expected {
		t.Errorf("Unexpected move was played:%s\n", fmt.Sprintf("Expected: %s\nGot: %s\n", expected.ToString(), mvStr))
	}
}

func TestBlackCanFindASimpleMaterialGainWithDiscoveredCheck(t *testing.T) {
	game := FromFen("3N1k2/N7/1p2ppR1/1P6/P2pP3/3Pb3/3r3n/2K5 b - - 1 1")
	r := NewRunner(NewCache(DEFAULT_CACHE_SIZE), NewPawnCache(DEFAULT_PAWNHASH_SIZE), 1)
	r.AddTimeManager(NewTimeManager(time.Now(), 400_000, true, 0, 0, false))
	e := r.Engines[0]
	e.Position = game.Position()
	e.Search(7)
	expected := NewMove(D2, G2, BlackRook, NoPiece, NoType, 0)
	mv := e.Move()
	mvStr := mv.ToString()
	if mv != expected {
		t.Errorf("Unexpected move was played:%s\n", fmt.Sprintf("Expected: %s\nGot: %s\n", expected.ToString(), mvStr))
	}
}

func TestWhiteShouldAcceptMaterialLossToAvoidCheckmate(t *testing.T) {
	game := FromFen("3N1k2/N7/1p2ppR1/1P6/P2pP3/3Pb3/3r3n/3K4 w - - 0 1")
	r := NewRunner(NewCache(DEFAULT_CACHE_SIZE), NewPawnCache(DEFAULT_PAWNHASH_SIZE), 1)
	r.AddTimeManager(NewTimeManager(time.Now(), 400_000, true, 0, 0, false))
	e := r.Engines[0]
	e.Position = game.Position()
	e.Search(7)
	expected := NewMove(D1, C1, WhiteKing, NoPiece, NoType, 0)
	mv := e.Move()
	mvStr := mv.ToString()
	if mv != expected {
		t.Errorf("Unexpected move was played:%s\n", fmt.Sprintf("Expected: %s\nGot: %s\n", expected.ToString(), mvStr))
	}
}

func TestSearchOnlyMove(t *testing.T) {
	game := FromFen("rnbqkbnr/ppppp1p1/7p/5P1Q/8/8/PPPP1PPP/RNB1KBNR b KQkq - 0 1")
	r := NewRunner(NewCache(DEFAULT_CACHE_SIZE), NewPawnCache(DEFAULT_PAWNHASH_SIZE), 1)
	r.AddTimeManager(NewTimeManager(time.Now(), 400_000, true, 0, 0, false))
	e := r.Engines[0]
	e.Position = game.Position()
	e.Search(7)
	expected := NewMove(G7, G6, BlackPawn, NoPiece, NoType, 0)
	mv := e.Move()
	score := e.Score()
	mvStr := mv.ToString()
	if mv != expected {
		t.Errorf("Unexpected move was played:%s\n", fmt.Sprintf("Expected: %s\nGot: %s\n", expected.ToString(), mvStr))
	}
	if score != -CHECKMATE_EVAL+2 {
		t.Errorf("Unexpected eval was returned:%s\n", fmt.Sprintf("Expected: %d\nGot: %d\n", -CHECKMATE_EVAL+2, score))
	}
}

func TestWhiteCanFindMateInTwo(t *testing.T) {
	game := FromFen("3N1k2/N7/1p2ppR1/1P6/P2pP3/3Pbn2/3r4/4K3 w - - 2 2")
	r := NewRunner(NewCache(DEFAULT_CACHE_SIZE), NewPawnCache(DEFAULT_PAWNHASH_SIZE), 1)
	r.AddTimeManager(NewTimeManager(time.Now(), 400_000, true, 0, 0, false))
	e := r.Engines[0]
	e.Position = game.Position()
	e.Search(7)
	expected := NewMove(E1, F1, WhiteKing, NoPiece, NoType, 0)
	mv := e.Move()
	mvStr := mv.ToString()
	if mv != expected {
		t.Errorf("Unexpected move was played:%s\n", fmt.Sprintf("Expected: %s\nGot: %s\n", expected.ToString(), mvStr))
	}
	score := e.Score()
	if score != -CHECKMATE_EVAL+2 {
		t.Errorf("Unexpected eval was returned:%s\n", fmt.Sprintf("Expected: %d\nGot: %d\n", -CHECKMATE_EVAL+2, score))
	}
}

func TestNestedMakeUnMake(t *testing.T) {
	fen := "rnb1kbnr/pQpp1ppp/4p3/8/7q/2P5/PP1PPPPP/RNB1KBNR b KQkq - 0 1"
	g := FromFen(fen)
	p := g.Position()
	originalHash := p.Hash()

	m1 := NewMove(G8, E7, BlackKnight, NoPiece, NoType, 0)
	ep1, tg1, hc1, _ := p.MakeMove(m1)

	m2 := NewMove(G2, G3, WhitePawn, NoPiece, NoType, 0)
	ep2, tg2, hc2, _ := p.MakeMove(m2)

	m3 := NewMove(H4, G5, BlackQueen, NoPiece, NoType, 0)
	ep3, tg3, hc3, _ := p.MakeMove(m3)

	m4 := NewMove(G3, G4, WhitePawn, NoPiece, NoType, 0)
	ep4, tg4, hc4, _ := p.MakeMove(m4)

	m5 := NewMove(C8, B7, BlackBishop, WhiteQueen, NoType, Capture)
	ep5, tg5, hc5, _ := p.MakeMove(m5)

	m6 := NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0)
	ep6, tg6, hc6, _ := p.MakeMove(m6)

	actualFen := g.Fen()
	expectedFen := "rn2kb1r/pbppnppp/4p3/6q1/1P4P1/2P5/P2PPP1P/RNB1KBNR b KQkq b3 0 1"
	if actualFen != expectedFen {
		t.Errorf("Unexected Fen after making the moves:\n%s", fmt.Sprintf("Got: %s\nExpected: %s\n", actualFen, expectedFen))
	}

	p.UnMakeMove(m6, tg6, ep6, hc6)
	p.UnMakeMove(m5, tg5, ep5, hc5)
	p.UnMakeMove(m4, tg4, ep4, hc4)
	p.UnMakeMove(m3, tg3, ep3, hc3)
	p.UnMakeMove(m2, tg2, ep2, hc2)
	p.UnMakeMove(m1, tg1, ep1, hc1)

	endHash := p.Hash()
	actualFen = g.Fen()
	expectedFen = fen
	if actualFen != expectedFen {
		t.Errorf("Unexected Fen after unmaking the moves:\n%s", fmt.Sprintf("Got: %s\nExpected: %s\n", actualFen, expectedFen))
	}
	if originalHash != endHash {
		t.Errorf("Nested Make/UnMake broke hashing %s", fmt.Sprintf("Got: %d\nExpected: %d\n", endHash, originalHash))
	}
}

func TestReubenFineBasicChessEndingsPosition70(t *testing.T) {
	fen := "8/k7/3p4/p2P1p2/P2P1P2/8/8/K7 w - - 0 1"
	game := FromFen(fen)
	r := NewRunner(NewCache(DEFAULT_CACHE_SIZE), NewPawnCache(DEFAULT_PAWNHASH_SIZE), 1)
	r.AddTimeManager(NewTimeManager(time.Now(), 400_000, true, 0, 0, false))
	e := r.Engines[0]
	e.Position = game.Position()
	e.Search(25)
	expected := NewMove(A1, B1, WhiteKing, NoPiece, NoType, 0)
	mv := e.Move()
	mvStr := mv.ToString()
	if mv != expected {
		t.Errorf("Unexpected move was played:%s\n", fmt.Sprintf("Expected: %s\nGot: %s\n", expected.ToString(), mvStr))
	}
}

func TestSearchFindsThreeFoldRepetitionToAvoidMate(t *testing.T) {
	fen := "k7/3RR3/8/8/8/1q6/8/K1RRRR2 b - - 0 1"
	game := FromFen(fen)
	r := NewRunner(NewCache(DEFAULT_CACHE_SIZE), NewPawnCache(DEFAULT_PAWNHASH_SIZE), 1)
	r.AddTimeManager(NewTimeManager(time.Now(), 400_000, true, 0, 0, false))
	e := r.Engines[0]
	e.Position = game.Position()
	e.Search(13)
	expected := []Move{
		NewMove(B3, A3, BlackQueen, NoPiece, NoType, 0),
		NewMove(A1, B1, WhiteKing, NoPiece, NoType, 0),
		NewMove(A3, B3, BlackQueen, NoPiece, NoType, 0),
		NewMove(B1, A1, WhiteKing, NoPiece, NoType, 0),
		NewMove(B3, A3, BlackQueen, NoPiece, NoType, 0)}
	actual := e.pv.line
	if equalMoves(expected, actual) {
		actualString := e.pv.ToString()
		e.pv.line = expected
		expectedString := e.pv.ToString()
		t.Errorf("Unexpected move was played:%s\n", fmt.Sprintf("Expected: %s\nGot: %s\n", expectedString, actualString))
	}
}

func equalMoves(moves1 []Move, moves2 []Move) bool {
	if len(moves1) != len(moves2) {
		return false
	}
	for _, m1 := range moves1 {
		exists := false
		for _, m2 := range moves2 {
			if m1 == m2 {
				exists = true
				break
			}
		}
		if !exists {
			fmt.Println("Missing", m1.ToString(), m1.Tag())
			return false
		}
	}
	return true
}
