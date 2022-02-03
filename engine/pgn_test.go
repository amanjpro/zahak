package engine

import (
	"fmt"
	"testing"
)

func TestPGNSearchmovesParsing(t *testing.T) {
	fen := "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q2/PPPBBPpP/R3K2R b KQkq - 0 1"
	game := FromFen(fen)
	actual := game.position.ParseSearchMoves([]string{"g2h1q", "b4c3", "   ", "\n\t", "wtime"})
	expected := []Move{
		NewMove(G2, H1, BlackPawn, WhiteRook, Queen, Capture),
		NewMove(B4, C3, BlackPawn, WhiteKnight, NoType, Capture),
	}
	if !equalMoves(expected, actual) {
		fmt.Println("Expected:")
		for _, i := range expected {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		fmt.Println("Got:")
		for _, i := range actual {
			fmt.Println(i.ToString(), i.PromoType(), i.Tag())
		}
		t.Errorf("Expected different number of moves to be generated%s",
			fmt.Sprintf("\nExpected: %d\nGot: %d\n", len(expected), len(actual)))
	}
}
