package main

import (
	"fmt"
	"testing"
)

func TestSearchOnlyMove(t *testing.T) {
	game := FromFen("rnbqkbnr/ppppp1p1/7p/5P1Q/8/8/PPPP1PPP/RNB1KBNR b KQkq - 0 1", true)
	fmt.Println(game.position.board.Draw())
	evalMove := search(game.position, 10)
	expected := Move{G7, G6, NoType, 0}
	mv := *evalMove.move
	mvStr := mv.ToString()
	fmt.Println(game.position.board.Draw())
	if mv != expected {
		t.Errorf("Unexpected move was played:%s\n", fmt.Sprintf("Expected: %s\nGot: %s\n", expected.ToString(), mvStr))
	}
	if evalMove.eval != CHECKMATE_EVAL {
		t.Errorf("Unexpected eval was returned:%s\n", fmt.Sprintf("Expected: %f\nGot: %f\n", CHECKMATE_EVAL, evalMove.eval))
	}
}

func TestNestedMakeUnMake(t *testing.T) {
	fen := "rnb1kbnr/pQpp1ppp/4p3/8/7q/2P5/PP1PPPPP/RNB1KBNR b KQkq - 0 1"
	g := FromFen(fen, true)
	p := g.position

	m1 := &Move{G8, E7, NoType, 0}
	ep1 := p.enPassant
	tg1 := p.tag
	cp1 := p.MakeMove(m1)

	m2 := &Move{G2, G3, NoType, 0}
	ep2 := p.enPassant
	tg2 := p.tag
	cp2 := p.MakeMove(m2)

	m3 := &Move{H4, G5, NoType, 0}
	ep3 := p.enPassant
	tg3 := p.tag
	cp3 := p.MakeMove(m3)

	m4 := &Move{G3, G4, NoType, 0}
	ep4 := p.enPassant
	tg4 := p.tag
	cp4 := p.MakeMove(m4)

	m5 := &Move{C8, B7, NoType, Capture}
	ep5 := p.enPassant
	tg5 := p.tag
	cp5 := p.MakeMove(m5)

	m6 := &Move{B2, B4, NoType, 0}
	ep6 := p.enPassant
	tg6 := p.tag
	cp6 := p.MakeMove(m6)

	actualFen := g.Fen()
	expectedFen := "rn2kb1r/pbppnppp/4p3/6q1/1P4P1/2P5/P2PPP1P/RNB1KBNR b KQkq b3 0 1"
	if actualFen != expectedFen {
		t.Errorf("Unexected Fen after making the moves:\n%s", fmt.Sprintf("Got: %s\nExpected: %s\n", actualFen, expectedFen))
	}

	p.UnMakeMove(m6, tg6, ep6, cp6)
	p.UnMakeMove(m5, tg5, ep5, cp5)
	p.UnMakeMove(m4, tg4, ep4, cp4)
	p.UnMakeMove(m3, tg3, ep3, cp3)
	p.UnMakeMove(m2, tg2, ep2, cp2)
	p.UnMakeMove(m1, tg1, ep1, cp1)

	actualFen = g.Fen()
	expectedFen = fen
	if actualFen != expectedFen {
		t.Errorf("Unexected Fen after unmaking the moves:\n%s", fmt.Sprintf("Got: %s\nExpected: %s\n", actualFen, expectedFen))
	}
}
