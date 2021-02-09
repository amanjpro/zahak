package main

import (
	"fmt"
	"os"
)

type PerftNodes struct {
	nodes      int64
	checks     int64
	captures   int64
	enPassants int64
	castles    int64
	promotions int64
	checkmates int64
}

func StartPerftTest() {
	result := testStartingPosDepth0() &&
		testStartingPosDepth1() &&
		testStartingPosDepth2() &&
		testStartingPosDepth3() &&
		testStartingPosDepth4() &&
		testStartingPosDepth5() &&
		testStartingPosDepth6() &&
		testStartingPosDepth9() &&
		testMiddleGameDepth1() &&
		testMiddleGameDepth2() &&
		testMiddleGameDepth3() &&
		testMiddleGameDepth4() &&
		testMiddleGameDepth5() &&
		testMiddleGameDepth6() &&
		testEndGamePromotionDepth1() &&
		testEndGamePromotionDepth2() &&
		testEndGamePromotionDepth3() &&
		testEndGamePromotionDepth4() &&
		testEndGamePromotionDepth5() &&
		testEndGamePromotionDepth6()
	if !result {
		os.Exit(1)
	}
}

func testStartingPosDepth0() bool {
	return test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 0, PerftNodes{1, 0, 0, 0, 0, 0, 0})
}

func testStartingPosDepth1() bool {
	return test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 1, PerftNodes{20, 0, 0, 0, 0, 0, 0})
}

func testStartingPosDepth2() bool {
	return test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 2, PerftNodes{400, 0, 0, 0, 0, 0, 0})
}

func testStartingPosDepth3() bool {
	return test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 3, PerftNodes{8902, 12, 34, 0, 0, 0, 0})
}

func testStartingPosDepth4() bool {
	return test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 4, PerftNodes{197281, 469, 1576, 0, 0, 0, 8})
}

func testStartingPosDepth5() bool {
	return test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 5, PerftNodes{4865609, 27351 /* 27351 + 6 */, 82725, 258, 0, 0, 347})
}

func testStartingPosDepth6() bool {
	return test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 6, PerftNodes{119060324, 809099 /* 329 + 46? */, 2812008, 5248, 0, 0, 10828})
}

func testStartingPosDepth9() bool {
	return test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 9, PerftNodes{2439530234167, 36095901903, /* 37101713	+ 5547231	 ? */
		125208536153, 319496827, 1784356000, 17334376, 400191963})
}

func testMiddleGameDepth1() bool {
	return testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 1,
		PerftNodes{48, 0, 0, 0, 0, 0, 0})
}

func testMiddleGameDepth2() bool {
	return testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 2,
		PerftNodes{2039, 0, 0, 0, 0, 0, 0})
}

func testMiddleGameDepth3() bool {
	return testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 3,
		PerftNodes{97862, 0, 0, 0, 0, 0, 0})
}

func testMiddleGameDepth4() bool {
	return testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 4,
		PerftNodes{4085603, 0, 0, 0, 0, 0, 0})
}

func testMiddleGameDepth5() bool {
	return testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 5,
		PerftNodes{193690690, 0, 0, 0, 0, 0, 0})
}

func testMiddleGameDepth6() bool {
	return testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 6,
		PerftNodes{8031647685, 0, 0, 0, 0, 0, 0})
}

func testEndGamePromotionDepth1() bool {
	return testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 1,
		PerftNodes{24, 0, 0, 0, 0, 0, 0})
}

func testEndGamePromotionDepth2() bool {
	return testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 2,
		PerftNodes{496, 0, 0, 0, 0, 0, 0})
}

func testEndGamePromotionDepth3() bool {
	return testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 3,
		PerftNodes{9483, 0, 0, 0, 0, 0, 0})
}

func testEndGamePromotionDepth4() bool {
	return testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 4,
		PerftNodes{182838, 0, 0, 0, 0, 0, 0})
}

func testEndGamePromotionDepth5() bool {
	return testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 5,
		PerftNodes{3605103, 0, 0, 0, 0, 0, 0})
}

func testEndGamePromotionDepth6() bool {
	return testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 6,
		PerftNodes{71179139, 0, 0, 0, 0, 0, 0})
}

func test(fen string, depth int8, expected PerftNodes) bool {
	g := FromFen(fen, true)
	actual := PerftNodes{0, 0, 0, 0, 0, 0, 0}
	perft(g.position, depth, NoType, 0, &actual)
	if actual != expected {
		fmt.Printf("test failed\nExpected: %d\nGot: %d\n", expected, actual)
		return false
	}
	return true
}

func testNodesOnly(fen string, depth int8, expected PerftNodes) bool {
	g := FromFen(fen, true)
	actual := PerftNodes{0, 0, 0, 0, 0, 0, 0}
	perft(g.position, depth, NoType, 0, &actual)
	if actual.nodes != expected.nodes {
		fmt.Printf("test failed\nExpected: %d\nGot: %d\n", expected.nodes, actual.nodes)
		return false
	}
	return true
}

func perft(p *Position, depth int8, lastPromo PieceType, lastTag MoveTag, acc *PerftNodes) {
	isCheckmate := p.Status() == Checkmate
	if depth == 0 || isCheckmate {
		acc.nodes += 1
		if isCheckmate {
			acc.checkmates += 1
		}
		if lastTag&Check != 0 {
			acc.checks += 1
		}

		if lastTag&Capture != 0 && lastTag&EnPassant == 0 {
			acc.captures += 1
		}
		if lastTag&EnPassant != 0 {
			acc.enPassants += 1
		}
		if lastTag&QueenSideCastle != 0 || lastTag&KingSideCastle != 0 {
			acc.castles += 1
		}
		if lastPromo != NoType {
			acc.promotions += 1
		}
		return
	}

	moves := p.LegalMoves()

	for _, move := range moves {
		tag := p.tag
		ep := p.enPassant
		cp := p.MakeMove(move)
		perft(p, depth-1, move.promoType, move.moveTag, acc)
		p.UnMakeMove(move, tag, ep, cp)
	}
}
