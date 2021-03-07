package engine

import (
	"fmt"
	"testing"
)

func TestPGNParsing(t *testing.T) {
	fen := "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q2/PPPBBPpP/R3K2R b KQkq - 0 1"
	game := FromFen(fen, true)
	actual := game.position.ParseMoves([]string{"g2h1q", "e2f1", "   ", "\n\t", "h8h2"})
	expected := []Move{
		NewMove(G2, H1, BlackPawn, WhiteRook, Queen, Capture|Check),
		NewMove(E2, F1, WhiteBishop, NoPiece, NoType, 0),
		NewMove(H8, H2, BlackRook, WhitePawn, NoType, Capture),
	}
	if !equalMoves(expected, actual) {
		fmt.Println("Got:")
		for _, i := range expected {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		fmt.Println("Expected:")
		for _, i := range actual {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected different number of moves to be generated%s",
			fmt.Sprintf("\nExpected: %d\nGot: %d\n", len(expected), len(actual)))
	}
}
