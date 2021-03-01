package search

import (
	"fmt"
	"testing"

	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func TestBlackShouldFindEscape(t *testing.T) {
	game := FromFen("3rbbn1/BQ1kp3/2p1q2p/N4p2/8/3P4/P1P2PPP/5RK1 b - - 0 27", true)
	e := NewEngine()
	e.ThinkTime = 400_000
	e.Search(game.Position(), 7, 27)
	expected := Move{D7, D6, NoType, 0}
	mv := e.Move()
	mvStr := mv.ToString()
	if mv != expected {
		t.Errorf("Unexpected move was played:%s\n", fmt.Sprintf("Expected: %s\nGot: %s\n", expected.ToString(), mvStr))
	}
}

func TestBlackCanFindASimpleTactic(t *testing.T) {
	game := FromFen("3N1k2/N7/1p2ppR1/1P6/P2pP3/3Pb3/2r4n/3K4 b - - 0 1", true)
	e := NewEngine()
	e.ThinkTime = 400_000
	e.Search(game.Position(), 7, 1)
	expected := Move{C2, D2, NoType, Check}
	mv := e.Move()
	mvStr := mv.ToString()
	if mv != expected {
		t.Errorf("Unexpected move was played:%s\n", fmt.Sprintf("Expected: %s\nGot: %s\n", expected.ToString(), mvStr))
	}
}

func TestBlackCanFindASimpleMaterialGainWithDiscoveredCheck(t *testing.T) {
	game := FromFen("3N1k2/N7/1p2ppR1/1P6/P2pP3/3Pb3/3r3n/2K5 b - - 1 1", true)
	e := NewEngine()
	e.ThinkTime = 400_000
	e.Search(game.Position(), 7, 1)
	expected := Move{D2, G2, NoType, Check}
	mv := e.Move()
	mvStr := mv.ToString()
	if mv != expected {
		t.Errorf("Unexpected move was played:%s\n", fmt.Sprintf("Expected: %s\nGot: %s\n", expected.ToString(), mvStr))
	}
}

func TestWhiteShouldAcceptMaterialLossToAvoidCheckmate(t *testing.T) {
	game := FromFen("3N1k2/N7/1p2ppR1/1P6/P2pP3/3Pb3/3r3n/3K4 w - - 0 1", true)
	e := NewEngine()
	e.ThinkTime = 400_000
	e.Search(game.Position(), 7, 1)
	expected := Move{D1, C1, NoType, 0}
	mv := e.Move()
	mvStr := mv.ToString()
	if mv != expected {
		t.Errorf("Unexpected move was played:%s\n", fmt.Sprintf("Expected: %s\nGot: %s\n", expected.ToString(), mvStr))
	}
}

func TestSearchOnlyMove(t *testing.T) {
	game := FromFen("rnbqkbnr/ppppp1p1/7p/5P1Q/8/8/PPPP1PPP/RNB1KBNR b KQkq - 0 1", true)
	e := NewEngine()
	e.ThinkTime = 400_000
	e.Search(game.Position(), 7, 1)
	expected := Move{G7, G6, NoType, 0}
	mv := e.Move()
	score := e.Score()
	mvStr := mv.ToString()
	if mv != expected {
		t.Errorf("Unexpected move was played:%s\n", fmt.Sprintf("Expected: %s\nGot: %s\n", expected.ToString(), mvStr))
	}
	if score != -CHECKMATE_EVAL {
		t.Errorf("Unexpected eval was returned:%s\n", fmt.Sprintf("Expected: %d\nGot: %d\n", -CHECKMATE_EVAL, score))
	}
}

func TestWhiteCanFindMateInTwo(t *testing.T) {
	game := FromFen("3N1k2/N7/1p2ppR1/1P6/P2pP3/3Pbn2/3r4/4K3 w - - 2 2", true)
	e := NewEngine()
	e.ThinkTime = 400_000
	e.Search(game.Position(), 7, 1)
	expected := Move{E1, F1, NoType, 0}
	mv := e.Move()
	mvStr := mv.ToString()
	if mv != expected {
		t.Errorf("Unexpected move was played:%s\n", fmt.Sprintf("Expected: %s\nGot: %s\n", expected.ToString(), mvStr))
	}
	score := e.Score()
	if score != -CHECKMATE_EVAL {
		t.Errorf("Unexpected eval was returned:%s\n", fmt.Sprintf("Expected: %d\nGot: %d\n", -CHECKMATE_EVAL, score))
	}
}

func TestNestedMakeUnMake(t *testing.T) {
	fen := "rnb1kbnr/pQpp1ppp/4p3/8/7q/2P5/PP1PPPPP/RNB1KBNR b KQkq - 0 1"
	g := FromFen(fen, true)
	p := g.Position()
	originalHash := p.Hash()

	m1 := Move{G8, E7, NoType, 0}
	cp1, ep1, tg1, hc1 := p.MakeMove(m1)

	m2 := Move{G2, G3, NoType, 0}
	cp2, ep2, tg2, hc2 := p.MakeMove(m2)

	m3 := Move{H4, G5, NoType, 0}
	cp3, ep3, tg3, hc3 := p.MakeMove(m3)

	m4 := Move{G3, G4, NoType, 0}
	cp4, ep4, tg4, hc4 := p.MakeMove(m4)

	m5 := Move{C8, B7, NoType, Capture}
	cp5, ep5, tg5, hc5 := p.MakeMove(m5)

	m6 := Move{B2, B4, NoType, 0}
	cp6, ep6, tg6, hc6 := p.MakeMove(m6)

	actualFen := g.Fen()
	expectedFen := "rn2kb1r/pbppnppp/4p3/6q1/1P4P1/2P5/P2PPP1P/RNB1KBNR b KQkq b3 0 1"
	if actualFen != expectedFen {
		t.Errorf("Unexected Fen after making the moves:\n%s", fmt.Sprintf("Got: %s\nExpected: %s\n", actualFen, expectedFen))
	}

	p.UnMakeMove(m6, tg6, ep6, cp6, hc6)
	p.UnMakeMove(m5, tg5, ep5, cp5, hc5)
	p.UnMakeMove(m4, tg4, ep4, cp4, hc4)
	p.UnMakeMove(m3, tg3, ep3, cp3, hc3)
	p.UnMakeMove(m2, tg2, ep2, cp2, hc2)
	p.UnMakeMove(m1, tg1, ep1, cp1, hc1)

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

func TestSearchFindsThreeFoldRepetitionToAvoidMate(t *testing.T) {
	fen := "k7/3RR3/8/8/8/1q6/8/K1RRRR2 b - - 0 1"
	game := FromFen(fen, true)
	e := NewEngine()
	e.ThinkTime = 400_000
	e.Search(game.Position(), 13, 1)
	expected := []Move{
		Move{B3, A3, NoType, Check},
		Move{A1, B1, NoType, 0},
		Move{A3, B3, NoType, Check},
		Move{B1, A1, NoType, 0},
		Move{B3, A3, NoType, Check},
		Move{A1, B1, NoType, 0},
		Move{A3, B3, NoType, Check},
		Move{B1, A1, NoType, 0},
		Move{B3, A3, NoType, Check},
		Move{A1, B1, NoType, 0},
		Move{A3, B3, NoType, Check},
		Move{B1, A1, NoType, 0}}
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
			fmt.Println("Missing", m1.ToString(), m1.Tag)
			return false
		}
	}
	return true
}
