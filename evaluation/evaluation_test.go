package evaluation

import (
	. "github.com/amanjpro/zahak/engine"
	"testing"
)

func TestMaterialValue(t *testing.T) {
	fen := "rnb2bnr/ppqppkpp/8/2p5/4P3/8/PPPP1PPP/RNB1KBNR w KQ - 0 1"
	game := FromFen(fen, false)

	actual := Evaluate(game.Position())

	if actual >= 0 {
		t.Errorf("Expected: a negative number\nGot: %d\n", actual)
	}

	fen = "3N1k2/N7/1p2ppR1/1P6/P2pP3/3Pb3/1K5n/8 w - - 0 4"

	game = FromFen(fen, false)

	actual = Evaluate(game.Position())

	if actual <= 0 {
		t.Errorf("Expected: a positive number\nGot: %d\n", actual)
	}
}
