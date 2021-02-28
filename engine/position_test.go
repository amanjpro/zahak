package engine

import (
	"testing"
)

func TestMakeMove(t *testing.T) {
	game := FromFen("rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R w KQkq d6 0 1", true)
	move := NewMove(F3, G4, WhiteBishop, NoPiece, NoType, 0)
	game.position.MakeMove(move)
	fen := game.Fen()
	expected := "rnbqkbnr/pPp1pppp/4P3/3pP3/4p1B1/6N1/PP3PPP/RNBQK2R b KQkq - 1 1"
	if fen != expected {
		t.Errorf("Move was not generated properly\nGot: %s\n", fen)
		t.Errorf("But expected: %s\n", expected)
	}
}

func TestMakeMoveDoublePushPawn(t *testing.T) {
	game := FromFen("rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R w KQkq d6 0 1", true)
	move := NewMove(H2, H4, WhitePawn, NoPiece, NoType, 0)
	game.position.MakeMove(move)
	fen := game.Fen()
	expected := "rnbqkbnr/pPp1pppp/4P3/3pP3/4p2P/5BN1/PP3PP1/RNBQK2R b KQkq h3 0 1"
	if fen != expected {
		t.Errorf("Move was not generated properly\nGot: %s\n", fen)
		t.Errorf("But expected: %s\n", expected)
	}
}

func TestMakeMoveCapture(t *testing.T) {
	game := FromFen("rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R w KQkq d6 0 1", true)
	move := NewMove(F3, E4, WhiteBishop, BlackPawn, NoType, Capture)
	game.position.MakeMove(move)
	fen := game.Fen()
	expected := "rnbqkbnr/pPp1pppp/4P3/3pP3/4B3/6N1/PP3PPP/RNBQK2R b KQkq - 0 1"
	if fen != expected {
		t.Errorf("Move was not generated properly\nGot: %s\n", fen)
		t.Errorf("But expected: %s\n", expected)
	}
}

func TestMakeMoveCastling(t *testing.T) {
	game := FromFen("rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R w KQkq d6 0 1", true)
	move := NewMove(E1, G1, WhiteKing, NoPiece, NoType, KingSideCastle)
	game.position.MakeMove(move)
	fen := game.Fen()
	expected := "rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQ1RK1 b kq - 1 1"
	if fen != expected {
		t.Errorf("Move was not generated properly\nGot: %s\n", fen)
		t.Errorf("But expected: %s\n", expected)
	}
}

func TestMakeMoveEnPassant(t *testing.T) {
	game := FromFen("rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R w KQkq d6 0 1", true)
	move := NewMove(E5, D6, WhitePawn, BlackPawn, NoType, EnPassant|Capture)
	game.position.MakeMove(move)
	fen := game.Fen()
	expected := "rnbqkbnr/pPp1pppp/3PP3/8/4p3/5BN1/PP3PPP/RNBQK2R b KQkq - 0 1"
	if fen != expected {
		t.Errorf("Move was not generated properly\nGot: %s\n", fen)
		t.Errorf("But expected: %s\n", expected)
	}
}

func TestMakeMovePromotion(t *testing.T) {
	game := FromFen("rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R w KQkq d6 0 1", true)
	move := NewMove(B7, A8, WhitePawn, BlackRook, Queen, Capture)
	game.position.MakeMove(move)
	fen := game.Fen()
	expected := "Qnbqkbnr/p1p1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R b KQk - 0 1"
	if fen != expected {
		t.Errorf("Move was not generated properly\nGot: %s\n", fen)
		t.Errorf("But expected: %s\n", expected)
	}
}

func TestUnMakeMove(t *testing.T) {
	startFen := "rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R w KQkq d6 0 1"
	game := FromFen(startFen, true)
	startHash := game.position.Hash()
	move := NewMove(F3, G4, WhiteBishop, NoPiece, NoType, 0)
	ep, tag, hc := game.position.MakeMove(move)
	game.position.UnMakeMove(move, tag, ep, hc)
	fen := game.Fen()
	if fen != startFen {
		t.Errorf("Move was not undone properly\nGot: %s\n", fen)
		t.Errorf("But expected: %s\n", startFen)
	}
	newHash := game.position.Hash()
	if startHash != newHash {
		t.Errorf("Move was not undone properly\nGot hash: %d\n", newHash)
		t.Errorf("But expected: %d\n", startHash)
	}
}

func TestUnMakeMoveDoublePushPawn(t *testing.T) {
	startFen := "rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R w KQkq d6 0 1"
	game := FromFen(startFen, true)
	startHash := game.position.Hash()
	move := NewMove(H2, H4, WhitePawn, NoPiece, NoType, 0)
	ep, tag, hc := game.position.MakeMove(move)
	game.position.UnMakeMove(move, tag, ep, hc)
	fen := game.Fen()
	if fen != startFen {
		t.Errorf("Move was not undone properly\nGot: %s\n", fen)
		t.Errorf("But expected: %s\n", startFen)
	}
	newHash := game.position.Hash()
	if startHash != newHash {
		t.Errorf("Move was not undone properly\nGot hash: %d\n", newHash)
		t.Errorf("But expected: %d\n", startHash)
	}
}

func TestUnMakeMoveCapture(t *testing.T) {
	startFen := "rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R w KQkq d6 0 1"
	game := FromFen(startFen, true)
	startHash := game.position.Hash()
	move := NewMove(F3, E4, WhiteBishop, BlackPawn, NoType, Capture)
	ep, tag, hc := game.position.MakeMove(move)
	game.position.UnMakeMove(move, tag, ep, hc)
	fen := game.Fen()
	if fen != startFen {
		t.Errorf("Move was not undone properly\nGot: %s\n", fen)
		t.Errorf("But expected: %s\n", startFen)
	}
	newHash := game.position.Hash()
	if startHash != newHash {
		t.Errorf("Move was not undone properly\nGot hash: %d\n", newHash)
		t.Errorf("But expected: %d\n", startHash)
	}
}

func TestUnMakeMoveCastling(t *testing.T) {
	startFen := "rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R w KQkq d6 0 1"
	game := FromFen(startFen, true)
	startHash := game.position.Hash()
	move := NewMove(E1, G1, WhiteKing, NoPiece, NoType, KingSideCastle)
	ep, tag, hc := game.position.MakeMove(move)
	game.position.UnMakeMove(move, tag, ep, hc)
	fen := game.Fen()
	if fen != startFen {
		t.Errorf("Move was not undone properly\nGot: %s\n", fen)
		t.Errorf("But expected: %s\n", startFen)
	}
	newHash := game.position.Hash()
	if startHash != newHash {
		t.Errorf("Move was not undone properly\nGot hash: %d\n", newHash)
		t.Errorf("But expected: %d\n", startHash)
	}
}

func TestUnMakeMoveEnPassant(t *testing.T) {
	startFen := "rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R w KQkq d6 0 1"
	game := FromFen(startFen, true)
	startHash := game.position.Hash()
	move := NewMove(E5, D6, WhitePawn, BlackPawn, NoType, EnPassant|Capture)
	ep, tag, hc := game.position.MakeMove(move)
	game.position.UnMakeMove(move, tag, ep, hc)
	fen := game.Fen()
	if fen != startFen {
		t.Errorf("Move was not undone properly\nGot: %s\n", fen)
		t.Errorf("But expected: %s\n", startFen)
	}
	newHash := game.position.Hash()
	if startHash != newHash {
		t.Errorf("Move was not undone properly\nGot hash: %d\n", newHash)
		t.Errorf("But expected: %d\n", startHash)
	}
}

func TestUnMakeMovePromotion(t *testing.T) {
	startFen := "rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R w KQkq d6 0 1"
	game := FromFen(startFen, true)
	startHash := game.position.Hash()
	move := NewMove(B7, A8, WhitePawn, BlackRook, Queen, Capture)
	ep, tag, hc := game.position.MakeMove(move)
	game.position.UnMakeMove(move, tag, ep, hc)
	fen := game.Fen()
	if fen != startFen {
		t.Errorf("Move was not undone properly\nGot: %s\n", fen)
		t.Errorf("But expected: %s\n", startFen)
	}
	newHash := game.position.Hash()
	if startHash != newHash {
		t.Errorf("Move was not undone properly\nGot hash: %d\n", newHash)
		t.Errorf("But expected: %d\n", startHash)
	}
}

func TestThreeFoldRepetition(t *testing.T) {
	startFen := "k7/3RR3/8/8/8/1q6/8/K1RRRR2 b - - 0 1"
	game := FromFen(startFen, true)

	m1 := Move{B3, A3, NoType, Check}
	m2 := Move{A1, B1, NoType, 0}
	m3 := Move{A3, B3, NoType, Check}
	m4 := Move{B1, A1, NoType, 0}
	m5 := Move{B3, A3, NoType, Check}
	m6 := Move{A1, B1, NoType, 0}
	m7 := Move{A3, B3, NoType, Check}
	m8 := Move{B1, A1, NoType, 0}
	m9 := Move{B3, A3, NoType, Check}

	game.position.MakeMove(m1)
	game.position.MakeMove(m2)
	game.position.MakeMove(m3)
	game.position.MakeMove(m4)
	game.position.MakeMove(m5)
	game.position.MakeMove(m6)
	game.position.MakeMove(m7)
	game.position.MakeMove(m8)
	game.position.MakeMove(m9)

	if game.Status() != Draw {
		t.Errorf("Expected Draw, but got: %d", game.Status())
	}
}
