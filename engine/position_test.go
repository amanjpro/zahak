package engine

import (
	"testing"
)

func TestMakeMove(t *testing.T) {
	game := FromFen("rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R w KQkq d6 0 1", true)
	move := &Move{F3, G4, NoType, 0}
	game.position.MakeMove(move)
	fen := game.Fen()
	expected := "rnbqkbnr/pPp1pppp/4P3/3pP3/4p1B1/6N1/PP3PPP/RNBQK2R b KQkq - 0 1"
	if fen != expected {
		t.Errorf("Move was not generated properly\nGot: %s\n", fen)
		t.Errorf("But expected: %s\n", expected)
	}
}

func TestMakeMoveDoublePushPawn(t *testing.T) {
	game := FromFen("rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R w KQkq d6 0 1", true)
	move := &Move{H2, H4, NoType, 0}
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
	move := &Move{F3, E4, NoType, Capture}
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
	move := &Move{E1, G1, NoType, KingSideCastle}
	game.position.MakeMove(move)
	fen := game.Fen()
	expected := "rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQ1RK1 b kq - 0 1"
	if fen != expected {
		t.Errorf("Move was not generated properly\nGot: %s\n", fen)
		t.Errorf("But expected: %s\n", expected)
	}
}

func TestMakeMoveEnPassant(t *testing.T) {
	game := FromFen("rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R w KQkq d6 0 1", true)
	move := &Move{E5, D6, NoType, EnPassant | Capture}
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
	move := &Move{B7, A8, Queen, Capture}
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
	move := &Move{F3, G4, NoType, 0}
	cp, ep, tag, hc := game.position.MakeMove(move)
	game.position.UnMakeMove(move, tag, ep, cp, hc)
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
	move := &Move{H2, H4, NoType, 0}
	cp, ep, tag, hc := game.position.MakeMove(move)
	game.position.UnMakeMove(move, tag, ep, cp, hc)
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
	move := &Move{F3, E4, NoType, Capture}
	cp, ep, tag, hc := game.position.MakeMove(move)
	game.position.UnMakeMove(move, tag, ep, cp, hc)
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
	move := &Move{E1, G1, NoType, KingSideCastle}
	cp, ep, tag, hc := game.position.MakeMove(move)
	game.position.UnMakeMove(move, tag, ep, cp, hc)
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
	move := &Move{E5, D6, Pawn, EnPassant | Capture}
	cp, ep, tag, hc := game.position.MakeMove(move)
	game.position.UnMakeMove(move, tag, ep, cp, hc)
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
	move := &Move{B7, A8, Queen, Capture}
	cp, ep, tag, hc := game.position.MakeMove(move)
	game.position.UnMakeMove(move, tag, ep, cp, hc)
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
