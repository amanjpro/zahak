package engine

import (
	"fmt"
	"testing"
)

func TestAllPieces(t *testing.T) {
	fen := "r1bq1bnr/pppp1p1p/n3p3/2k3p1/2P3P1/7N/PPQPPP1P/RNB1KBR1 w KQkq - 0 1"
	g := FromFen(fen, true)
	expected := map[Square]Piece{
		A8: BlackRook,
		C8: BlackBishop,
		D8: BlackQueen,
		F8: BlackBishop,
		G8: BlackKnight,
		H8: BlackRook,
		A7: BlackPawn,
		B7: BlackPawn,
		C7: BlackPawn,
		D7: BlackPawn,
		E6: BlackPawn,
		F7: BlackPawn,
		H7: BlackPawn,
		G5: BlackPawn,
		A6: BlackKnight,
		C5: BlackKing,
		A1: WhiteRook,
		B1: WhiteKnight,
		C1: WhiteBishop,
		E1: WhiteKing,
		F1: WhiteBishop,
		G1: WhiteRook,
		A2: WhitePawn,
		B2: WhitePawn,
		C4: WhitePawn,
		D2: WhitePawn,
		E2: WhitePawn,
		F2: WhitePawn,
		G4: WhitePawn,
		H2: WhitePawn,
		H3: WhiteKnight,
		C2: WhiteQueen,
	}
	actual := g.position.Board.AllPieces()
	if !equalMaps(actual, expected) {
		err := fmt.Sprintf("Got: %x\nExpected%x", actual, expected)
		t.Errorf("Unexpected return by AllPieces: %s", err)
	}

	m := &Move{H3, G5, NoType, Capture}
	cp, ep, ot := g.position.MakeMove(m)
	g.position.UnMakeMove(m, ot, ep, cp)

	actual = g.position.Board.AllPieces()
	if !equalMaps(actual, expected) {
		err := fmt.Sprintf("Got: %x\nExpected%x", actual, expected)
		t.Errorf("Knight make/unmake move broke all pieces: %s", err)
	}

	m = &Move{G1, G3, NoType, 0}
	cp, ep, ot = g.position.MakeMove(m)
	g.position.UnMakeMove(m, ot, ep, cp)

	actual = g.position.Board.AllPieces()
	if !equalMaps(actual, expected) {
		err := fmt.Sprintf("Got: %x\nExpected%x", actual, expected)
		t.Errorf("Rook make/unmake move broke all pieces: %s", err)
	}
}

func equalMaps(ps1 map[Square]Piece, ps2 map[Square]Piece) bool {
	if len(ps1) != len(ps2) {
		return false
	}
	for sq1, p1 := range ps1 {
		p2, ok := ps2[sq1]
		if !ok || p1 != p2 {
			return false
		}
	}
	return true
}
