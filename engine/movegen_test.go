package engine

import (
	"fmt"
	"testing"
)

func TestBishopMoves(t *testing.T) {
	fen := "rnbqkbnr/pPp1pppp/4P3/3pP3/3p4/4B1N1/PP2BPPP/1NRQK2R w Kkq - 0 1"
	g := FromFen(fen, true)
	board := g.position.Board
	moves := make([]Move, 0, 8)
	g.position.bbSlidingMoves(board.whiteBishop, board.whitePieces, board.blackPieces,
		White, WhiteBishop, false, false, false, false, &moves)
	expectedMoves := []Move{
		NewMove(E2, F1, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E2, F3, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E2, G4, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E2, H5, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E2, D3, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E2, C4, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E2, B5, WhiteBishop, NoPiece, NoType, Check),
		NewMove(E2, A6, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E3, D2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E3, F4, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E3, G5, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E3, H6, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E3, D4, WhiteBishop, BlackPawn, NoType, Capture),
	}
	expectedLen := len(expectedMoves)
	if len(moves) != expectedLen || !equalMoves(expectedMoves, moves) {
		fmt.Println("Got:")
		for _, i := range moves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		fmt.Println("Expected:")
		for _, i := range expectedMoves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected different number of moves to be generated%s",
			fmt.Sprintf("\nExpected: %d\nGot: %d\n", expectedLen, len(moves)))
	}
}

func TestRookMoves(t *testing.T) {
	fen := "rnkqbbnr/ppp1pppp/4P3/3pP3/3P4/4B1N1/PP2BPPP/1NRQK2R w Kkq - 0 1"
	g := FromFen(fen, true)
	board := g.position.Board
	moves := make([]Move, 0, 8)
	g.position.bbSlidingMoves(board.whiteRook, board.whitePieces, board.blackPieces,
		White, WhiteRook, false, false, false, false, &moves)
	expectedMoves := []Move{
		NewMove(H1, G1, WhiteRook, NoPiece, NoType, 0),
		NewMove(H1, F1, WhiteRook, NoPiece, NoType, 0),
		NewMove(C1, C2, WhiteRook, NoPiece, NoType, 0),
		NewMove(C1, C3, WhiteRook, NoPiece, NoType, 0),
		NewMove(C1, C4, WhiteRook, NoPiece, NoType, 0),
		NewMove(C1, C5, WhiteRook, NoPiece, NoType, 0),
		NewMove(C1, C6, WhiteRook, NoPiece, NoType, 0),
		NewMove(C1, C7, WhiteRook, BlackPawn, NoType, Capture|Check),
	}
	expectedLen := len(expectedMoves)
	if len(moves) != expectedLen || !equalMoves(expectedMoves, moves) {
		fmt.Println("Got:")
		for _, i := range moves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		fmt.Println("Expected:")
		for _, i := range expectedMoves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected different number of moves to be generated%s",
			fmt.Sprintf("\nExpected: %d\nGot: %d\n", expectedLen, len(moves)))
	}
}

func TestQueenMoves(t *testing.T) {
	fen := "rnbqkbnr/pPp1pppp/4P3/3pP3/3p4/4B1N1/PP2BPPP/1NRQK2R w Kkq - 0 1"
	g := FromFen(fen, true)
	board := g.position.Board
	moves := make([]Move, 0, 8)
	g.position.bbSlidingMoves(board.whiteQueen, board.whitePieces, board.blackPieces,
		White, WhiteQueen, false, false, false, false, &moves)
	expectedMoves := []Move{
		NewMove(D1, D2, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, D3, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, D4, WhiteQueen, BlackPawn, NoType, Capture),
		NewMove(D1, C2, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, B3, WhiteQueen, NoPiece, NoType, 0),
		NewMove(D1, A4, WhiteQueen, NoPiece, NoType, Check),
	}
	expectedLen := len(expectedMoves)
	if len(moves) != expectedLen || !equalMoves(expectedMoves, moves) {
		fmt.Println("Got:")
		for _, i := range moves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		fmt.Println("Expected:")
		for _, i := range expectedMoves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected different number of moves to be generated%s",
			fmt.Sprintf("\nExpected: %d\nGot: %d\n", expectedLen, len(moves)))
	}
}

func TestKingMoves(t *testing.T) {
	fen := "rnbqkbn1/pPp1pppp/4P3/3pP3/3p4/4B1N1/PP1rBPPP/R3K2R w Kkq - 0 1"
	g := FromFen(fen, true)
	p := g.position
	board := g.position.Board
	moves := make([]Move, 0, 8)
	color := White
	taboo := tabooSquares(board, color)
	g.position.bbKingMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
		taboo, color, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide), false, false, false, false, &moves)
	expectedMoves := []Move{
		NewMove(E1, D2, WhiteKing, BlackRook, NoType, Capture),
		NewMove(E1, F1, WhiteKing, NoPiece, NoType, 0),
		NewMove(E1, G1, WhiteKing, NoPiece, NoType, KingSideCastle),
	}
	expectedLen := len(expectedMoves)
	if len(moves) != expectedLen || !equalMoves(expectedMoves, moves) {
		fmt.Println("Got:")
		for _, i := range moves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		fmt.Println("Expected:")
		for _, i := range expectedMoves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected different number of moves to be generated%s",
			fmt.Sprintf("\nExpected: %d\nGot: %d\n", expectedLen, len(moves)))
	}
}

func TestKingCastlingWithOccupiedSquares(t *testing.T) {
	fen := "rnbqkbnr/1p6/p1p3Pp/1B1pp2Q/1P6/B7/P1PP1PPP/RN2K1NR w KQkq - 0 1"
	g := FromFen(fen, true)
	p := g.position
	board := g.position.Board
	moves := make([]Move, 0, 8)
	color := White
	taboo := tabooSquares(board, color)
	g.position.bbKingMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
		taboo, color, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide), false, false, false, false, &moves)
	expectedMoves := []Move{
		NewMove(E1, E2, WhiteKing, NoPiece, NoType, 0),
		NewMove(E1, F1, WhiteKing, NoPiece, NoType, 0),
		NewMove(E1, D1, WhiteKing, NoPiece, NoType, 0),
	}
	expectedLen := len(expectedMoves)
	if len(moves) != expectedLen || !equalMoves(expectedMoves, moves) {
		fmt.Println("Got:")
		for _, i := range moves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		fmt.Println("Expected:")
		for _, i := range expectedMoves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected different number of moves to be generated%s",
			fmt.Sprintf("\nExpected: %d\nGot: %d\n", expectedLen, len(moves)))
	}
}

func TestKingQueenSideCastling(t *testing.T) {
	fen := "rnbqkbnr/1p6/p1p3Pp/1B1pp2Q/1P6/B7/P1PP1PPP/R3K1NR w KQkq - 0 1"
	g := FromFen(fen, true)
	p := g.position
	board := g.position.Board
	moves := make([]Move, 0, 8)
	color := White
	taboo := tabooSquares(board, color)
	g.position.bbKingMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
		taboo, color, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide), false, false, false, false, &moves)
	expectedMoves := []Move{
		NewMove(E1, E2, WhiteKing, NoPiece, NoType, 0),
		NewMove(E1, F1, WhiteKing, NoPiece, NoType, 0),
		NewMove(E1, D1, WhiteKing, NoPiece, NoType, 0),
		NewMove(E1, C1, WhiteKing, NoPiece, NoType, QueenSideCastle),
	}
	expectedLen := len(expectedMoves)
	if len(moves) != expectedLen || !equalMoves(expectedMoves, moves) {
		fmt.Println("Got:")
		for _, i := range moves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		fmt.Println("Expected:")
		for _, i := range expectedMoves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected different number of moves to be generated%s",
			fmt.Sprintf("\nExpected: %d\nGot: %d\n", expectedLen, len(moves)))
	}
}

func TestPawnMovesForWhite(t *testing.T) {
	fen := "rnbqkbn1/pPp1pppp/4P3/3pP3/3p4/4B1N1/PP1rBPPP/R3K2R w Kkq d6 0 1"
	g := FromFen(fen, true)
	p := g.position
	board := g.position.Board
	moves := make([]Move, 0, 8)
	color := White
	g.position.bbPawnMoves(board.whitePawn, board.whitePieces, board.blackPieces,
		color, p.EnPassant, false, false, false, false, &moves)
	expectedMoves := []Move{
		NewMove(H2, H4, WhitePawn, NoPiece, NoType, 0),
		NewMove(H2, H3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F3, WhitePawn, NoPiece, NoType, 0),
		NewMove(A2, A4, WhitePawn, NoPiece, NoType, 0),
		NewMove(A2, A3, WhitePawn, NoPiece, NoType, 0),
		NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0),
		NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0),
		NewMove(E5, D6, WhitePawn, BlackPawn, NoType, EnPassant|Capture),
		NewMove(E6, F7, WhitePawn, BlackPawn, NoType, Capture|Check),
		NewMove(B7, A8, WhitePawn, BlackRook, Queen, Capture),
		NewMove(B7, A8, WhitePawn, BlackRook, Rook, Capture),
		NewMove(B7, A8, WhitePawn, BlackRook, Bishop, Capture),
		NewMove(B7, A8, WhitePawn, BlackRook, Knight, Capture),
		NewMove(B7, C8, WhitePawn, BlackBishop, Queen, Capture),
		NewMove(B7, C8, WhitePawn, BlackBishop, Rook, Capture),
		NewMove(B7, C8, WhitePawn, BlackBishop, Bishop, Capture),
		NewMove(B7, C8, WhitePawn, BlackBishop, Knight, Capture),
	}
	expectedLen := len(expectedMoves)
	if len(moves) != expectedLen || !equalMoves(expectedMoves, moves) {
		fmt.Println("Got:")
		for _, i := range moves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		fmt.Println("Expected:")
		for _, i := range expectedMoves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected different number of moves to be generated%s",
			fmt.Sprintf("\nExpected: %d\nGot: %d\n", expectedLen, len(moves)))
	}
}

func TestPawnMovesForBlack(t *testing.T) {
	fen := "rnbqkbnr/ppp3pp/3p1p2/1P4P1/4pP2/N6N/P1PPP2P/R1BQKB1R b KQkq f3 0 1"
	g := FromFen(fen, true)
	p := g.position
	board := g.position.Board
	moves := make([]Move, 0, 8)
	color := Black
	g.position.bbPawnMoves(board.blackPawn, board.blackPieces, board.whitePieces,
		color, p.EnPassant, false, false, false, false, &moves)
	expectedMoves := []Move{
		NewMove(H7, H6, BlackPawn, NoPiece, NoType, 0),
		NewMove(H7, H5, BlackPawn, NoPiece, NoType, 0),
		NewMove(G7, G6, BlackPawn, NoPiece, NoType, 0),
		NewMove(F6, F5, BlackPawn, NoPiece, NoType, 0),
		NewMove(F6, G5, BlackPawn, WhitePawn, NoType, Capture),
		NewMove(E4, E3, BlackPawn, NoPiece, NoType, 0),
		NewMove(E4, F3, BlackPawn, WhitePawn, NoType, EnPassant|Capture),
		NewMove(D6, D5, BlackPawn, NoPiece, NoType, 0),
		NewMove(C7, C6, BlackPawn, NoPiece, NoType, 0),
		NewMove(C7, C5, BlackPawn, NoPiece, NoType, 0),
		NewMove(B7, B6, BlackPawn, NoPiece, NoType, 0),
		NewMove(A7, A6, BlackPawn, NoPiece, NoType, 0),
		NewMove(A7, A5, BlackPawn, NoPiece, NoType, 0),
	}
	expectedLen := len(expectedMoves)
	if len(moves) != expectedLen || !equalMoves(expectedMoves, moves) {
		fmt.Println("Got:")
		for _, i := range moves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		fmt.Println("Expected:")
		for _, i := range expectedMoves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected different number of moves to be generated%s",
			fmt.Sprintf("\nExpected: %d\nGot: %d\n", expectedLen, len(moves)))
	}
}

func TestKnightMoves(t *testing.T) {
	fen := "rnbqkbn1/pPp1pppp/4P3/1N1pP3/3p4/4B1N1/PP1rBPPP/R3K2R w Kkq d6 0 1"
	g := FromFen(fen, true)
	p := g.position
	b := p.Board
	moves := make([]Move, 0, 8)
	g.position.bbKnightMoves(WhiteKnight, b.whiteKnight, b.whitePieces, b.blackPieces, false, false, false, false, &moves)
	expectedMoves := []Move{
		NewMove(G3, F1, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G3, E4, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G3, F5, WhiteKnight, NoPiece, NoType, 0),
		NewMove(G3, H5, WhiteKnight, NoPiece, NoType, 0),
		NewMove(B5, A7, WhiteKnight, BlackPawn, NoType, Capture),
		NewMove(B5, A3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(B5, C7, WhiteKnight, BlackPawn, NoType, Capture|Check),
		NewMove(B5, C3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(B5, D4, WhiteKnight, BlackPawn, NoType, Capture),
		NewMove(B5, D6, WhiteKnight, NoPiece, NoType, Check),
	}
	expectedLen := len(expectedMoves)
	if len(moves) != expectedLen || !equalMoves(expectedMoves, moves) {
		fmt.Println("Got:")
		for _, i := range moves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		fmt.Println("Expected:")
		for _, i := range expectedMoves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected different number of moves to be generated%s",
			fmt.Sprintf("\nExpected: %d\nGot: %d\n", expectedLen, len(moves)))
	}
}

func TestCastleAndDiscoveredChecks(t *testing.T) {
	fen := "rnbq1bn1/pPp1pppp/4P3/3pP3/3p4/4B1N1/PP1rBPPP/k3K2R w Kkq - 0 1"
	g := FromFen(fen, true)
	p := g.position
	legalMoves := p.LegalMoves()
	move := NewMove(E1, G1, WhiteKing, NoPiece, NoType, Check|KingSideCastle)
	if !containsMove(legalMoves, move) {
		fmt.Println("Got:")
		for _, i := range legalMoves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected to see %s", fmt.Sprintf("%s %d", move.ToString(), move.Tag()))
	}
	move = NewMove(E1, D2, WhiteKing, BlackRook, NoType, Check|Capture)
	if !containsMove(legalMoves, move) {
		fmt.Println("Got:")
		for _, i := range legalMoves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected to see %s", fmt.Sprintf("%s %d", move.ToString(), move.Tag()))
	}
}

func TestCastleAndPawnAttack(t *testing.T) {
	fen := "r3k2r/p1ppqpb1/1n2pnp1/1b1PN3/1p2P3/P1N2Q2/1PPBBPpP/1R2K2R w Kkq - 0 1"
	g := FromFen(fen, true)
	p := g.position
	board := g.position.Board
	moves := make([]Move, 0, 8)
	color := White
	taboo := tabooSquares(board, color)
	g.position.bbKingMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
		taboo, color, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide), false, false, false, false, &moves)
	expectedMoves := []Move{
		NewMove(E1, D1, WhiteKing, NoPiece, NoType, 0),
	}
	expectedLen := len(expectedMoves)
	if !equalMoves(expectedMoves, moves) {
		fmt.Println(g.position.Board.Draw())
		fmt.Println("Got:")
		for _, i := range moves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		fmt.Println("Expected:")
		for _, i := range expectedMoves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected different number of moves to be generated%s",
			fmt.Sprintf("\nExpected: %d\nGot: %d\n", expectedLen, len(moves)))
	}
}

func TestLegalMoves(t *testing.T) {
	fen := "rn1q1bn1/pPp1pppp/4P3/1N1pP2Q/3p3b/4B3/PP1rBPPP/k3K2R w Kkq d6 0 1"
	g := FromFen(fen, true)
	p := g.position
	legalMoves := p.LegalMoves()
	expectedMoves := []Move{
		NewMove(H1, G1, WhiteRook, NoPiece, NoType, 0),
		NewMove(H1, F1, WhiteRook, NoPiece, NoType, 0),
		NewMove(E1, F1, WhiteKing, NoPiece, NoType, 0),
		NewMove(E1, G1, WhiteKing, NoPiece, NoType, Check|KingSideCastle),
		NewMove(E1, D2, WhiteKing, BlackRook, NoType, Check|Capture),
		NewMove(H2, H3, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G3, WhitePawn, NoPiece, NoType, 0),
		NewMove(G2, G4, WhitePawn, NoPiece, NoType, 0),
		NewMove(E2, F1, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E2, D1, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E2, F3, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E2, G4, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E2, D3, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E2, C4, WhiteBishop, NoPiece, NoType, 0),
		NewMove(B2, B3, WhitePawn, NoPiece, NoType, 0),
		NewMove(B2, B4, WhitePawn, NoPiece, NoType, 0),
		NewMove(A2, A3, WhitePawn, NoPiece, NoType, 0),
		NewMove(A2, A4, WhitePawn, NoPiece, NoType, 0),
		NewMove(E3, D4, WhiteBishop, BlackPawn, NoType, Capture),
		NewMove(E3, F4, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E3, G5, WhiteBishop, NoPiece, NoType, 0),
		NewMove(E3, H6, WhiteBishop, NoPiece, NoType, 0),
		NewMove(H5, H7, WhiteQueen, BlackPawn, NoType, Capture),
		NewMove(H5, H6, WhiteQueen, NoPiece, NoType, 0),
		NewMove(H5, F7, WhiteQueen, BlackPawn, NoType, Capture),
		NewMove(H5, G6, WhiteQueen, NoPiece, NoType, 0),
		NewMove(H5, G5, WhiteQueen, NoPiece, NoType, 0),
		NewMove(H5, F5, WhiteQueen, NoPiece, NoType, 0),
		NewMove(H5, G4, WhiteQueen, NoPiece, NoType, 0),
		NewMove(H5, F3, WhiteQueen, NoPiece, NoType, 0),
		NewMove(H5, H4, WhiteQueen, BlackBishop, NoType, Capture),
		NewMove(E5, D6, WhitePawn, BlackPawn, NoType, Capture|EnPassant),
		NewMove(B5, A3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(B5, C3, WhiteKnight, NoPiece, NoType, 0),
		NewMove(B5, A7, WhiteKnight, BlackPawn, NoType, Capture),
		NewMove(B5, C7, WhiteKnight, BlackPawn, NoType, Capture),
		NewMove(B5, D4, WhiteKnight, BlackPawn, NoType, Capture),
		NewMove(B5, D6, WhiteKnight, NoPiece, NoType, 0),
		NewMove(B5, D6, WhiteKnight, NoPiece, NoType, 0),
		NewMove(E6, F7, WhitePawn, BlackPawn, NoType, Capture),
		NewMove(B7, A8, WhitePawn, BlackRook, Queen, Capture),
		NewMove(B7, A8, WhitePawn, BlackRook, Rook, Capture),
		NewMove(B7, A8, WhitePawn, BlackRook, Bishop, Capture),
		NewMove(B7, A8, WhitePawn, BlackRook, Knight, Capture),
	}
	expectedLen := len(expectedMoves)
	if expectedLen != len(legalMoves) || !equalMoves(expectedMoves, legalMoves) {
		fmt.Println("Got:")
		for _, i := range legalMoves {
			fmt.Println(i.ToString(), i.MovingPiece(), i.CapturedPiece(), i.PromoType(), i.Tag())
		}
		fmt.Println("Expected:")
		for _, i := range expectedMoves {
			fmt.Println(i.ToString(), i.MovingPiece(), i.CapturedPiece(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected different number of moves to be generated%s",
			fmt.Sprintf("\nExpected: %d\nGot: %d\n", expectedLen, len(legalMoves)))
	}
}

func TestDoubleCheckResponses(t *testing.T) {
	fen := "5Q2/8/1q5P/8/6k1/5R2/6P1/2r3K1 w - - 0 1"
	g := FromFen(fen, true)
	p := g.position
	legalMoves := p.LegalMoves()
	expectedMoves := []Move{
		NewMove(G1, H2, WhiteKing, NoPiece, NoType, 0),
	}
	if !p.IsInCheck() {
		t.Errorf("Position is wrongfully considered not check for: %s", fen)
	}
	if !isDoubleCheck(p.Board, White) {
		t.Errorf("Position is wrongfully considered not double-check for: %s", fen)
	}
	if p.Status() != Unknown {
		t.Errorf("Position is wrongfully considered ended: %b", p.Status())
	}
	expectedLen := len(expectedMoves)
	if expectedLen != len(legalMoves) || !equalMoves(expectedMoves, legalMoves) {
		fmt.Println("Got:")
		for _, i := range legalMoves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		fmt.Println("Expected:")
		for _, i := range expectedMoves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected different number of moves to be generated%s",
			fmt.Sprintf("\nExpected: %d\nGot: %d\n", expectedLen, len(legalMoves)))
	}
}

func TestHasLegalMovesCheckmate(t *testing.T) {
	fen := "5Q2/8/1q5P/8/6k1/5R2/6PR/2r3K1 w - - 0 1"
	g := FromFen(fen, true)
	p := g.position
	hasMoves := p.HasLegalMoves()
	if hasMoves {
		for _, i := range p.LegalMoves() {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Position is wrongfully considered playable, %p", p.LegalMoves())
	}
}

func TestHasLegalMovesDraw(t *testing.T) {
	fen := "7k/5Q2/6K1/8/8/8/8/8 b - - 0 1"
	g := FromFen(fen, true)
	p := g.position
	hasMoves := p.HasLegalMoves()
	if hasMoves {
		t.Errorf("Position is wrongfully considered playable")
	}
}

func TestHasLegalMoves(t *testing.T) {
	fen := "5Q2/8/1q5P/8/6k1/5R2/6P1/2r3K1 w - - 0 1"
	g := FromFen(fen, true)
	p := g.position
	legalMoves1 := p.LegalMoves()
	hasMoves := p.HasLegalMoves()
	legalMoves2 := p.LegalMoves()
	if !hasMoves || !equalMoves(legalMoves1, legalMoves2) {
		fmt.Println("First call to LegalMoves")
		for _, i := range legalMoves1 {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		fmt.Println("Second call to LegalMoves")
		for _, i := range legalMoves2 {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Position is wrongfully considered lost, %p", p.LegalMoves())
	}
}

func TestLegalMovesInOpenning(t *testing.T) {
	fen := "rnbqkbnr/ppp3pp/3ppp2/1P6/6P1/N6N/P1PPPP1P/R1BQKB1R w KQkq - 0 1"
	g := FromFen(fen, true)
	p := g.position
	legalMoves := p.LegalMoves()
	expectedMoves := []Move{
		NewMove(H1, G1, WhiteRook, NoPiece, NoType, 0),
		NewMove(G4, G5, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F3, WhitePawn, NoPiece, NoType, 0),
		NewMove(F2, F4, WhitePawn, NoPiece, NoType, 0),
		NewMove(E2, E3, WhitePawn, NoPiece, NoType, 0),
		NewMove(E2, E4, WhitePawn, NoPiece, NoType, 0),
		NewMove(D2, D3, WhitePawn, NoPiece, NoType, 0),
		NewMove(D2, D4, WhitePawn, NoPiece, NoType, 0),
		NewMove(C2, C3, WhitePawn, NoPiece, NoType, 0),
		NewMove(C2, C4, WhitePawn, NoPiece, NoType, 0),
		NewMove(B5, B6, WhitePawn, NoPiece, NoType, 0),
		NewMove(A1, B1, WhiteRook, NoPiece, NoType, 0),
		NewMove(A3, C4, WhiteKnight, NoPiece, NoType, 0),
		NewMove(A3, B1, WhiteKnight, NoPiece, NoType, 0),
		NewMove(C1, B2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(F1, G2, WhiteBishop, NoPiece, NoType, 0),
		NewMove(H3, G5, WhiteKnight, NoPiece, NoType, 0),
		NewMove(H3, F4, WhiteKnight, NoPiece, NoType, 0),
		NewMove(H3, G1, WhiteKnight, NoPiece, NoType, 0),
	}
	expectedLen := len(expectedMoves)
	if expectedLen != len(legalMoves) || !equalMoves(expectedMoves, legalMoves) {
		fmt.Println("Got:")
		for _, i := range legalMoves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		fmt.Println("Expected:")
		for _, i := range expectedMoves {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected different number of moves to be generated%s",
			fmt.Sprintf("\nExpected: %d\nGot: %d\n", expectedLen, len(legalMoves)))
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

func containsMove(moves1 []Move, move Move) bool {
	exists := false
	for _, m := range moves1 {
		if m == move {
			exists = true
			break
		}
	}
	return exists
}
