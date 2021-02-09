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

func PerftTree(game Game, depth int, moves []*Move) {
	sum := int64(0)
	for _, move := range moves {
		game.position.MakeMove(move)
	}
	moves = game.position.LegalMoves()
	depth -= 1
	for _, move := range moves {
		if move.HasTag(EnPassant) {
			fmt.Println(move.HasTag(EnPassant))
		}
		tg := game.position.tag
		ep := game.position.enPassant
		cp := game.position.MakeMove(move)
		acc := PerftNodes{0, 0, 0, 0, 0, 0, 0}
		perft(game.position, depth, NoType, 0, &acc)
		fmt.Printf("%s %d\n", move.ToString(), acc.nodes)
		sum += acc.nodes
		game.position.UnMakeMove(move, tg, ep, cp)
	}

	fmt.Printf("\n%d\n", sum)
}

func StartPerftTest() {
	result := testStartingPosDepth0()
	result += testStartingPosDepth1()
	result += testStartingPosDepth2()
	result += testStartingPosDepth3()
	result += testStartingPosDepth4()
	// result += testStartingPosDepth5()
	// resu+= && testStartingPosDepth6()
	// resu+= && testStartingPosDepth9()
	result += testMiddleGameDepth1()
	result += testMiddleGameDepth2()
	result += testMiddleGameDepth3()
	result += testMiddleGameDepth4()
	// result += testMiddleGameDepth5()
	// result += testMiddleGameDepth6()
	result += testEndGamePromotionDepth1()
	result += testEndGamePromotionDepth2()
	result += testEndGamePromotionDepth3()
	result += testEndGamePromotionDepth4()
	result += testEndGamePromotionDepth5()
	result += testEndGamePromotionDepth6()
	if result > 0 {
		fmt.Printf("%d perft tests failed... hard luck\n", result)
		os.Exit(1)
	} else {
		fmt.Println("All perft tests passed")
	}
}

func testStartingPosDepth0() int8 {
	return test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 0, PerftNodes{1, 0, 0, 0, 0, 0, 0})
}

func testStartingPosDepth1() int8 {
	return test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 1, PerftNodes{20, 0, 0, 0, 0, 0, 0})
}

func testStartingPosDepth2() int8 {
	return test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 2, PerftNodes{400, 0, 0, 0, 0, 0, 0})
}

func testStartingPosDepth3() int8 {
	return test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 3, PerftNodes{8902, 12, 34, 0, 0, 0, 0})
}

func testStartingPosDepth4() int8 {
	return test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 4, PerftNodes{197281, 469, 1576, 0, 0, 0, 8})
}

func testStartingPosDepth5() int8 {
	return test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 5, PerftNodes{4865609, 27351 /* 27351 + 6 */, 82725, 258, 0, 0, 347})
}

func testStartingPosDepth6() int8 {
	return test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 6, PerftNodes{119060324, 809099 /* 329 + 46? */, 2812008, 5248, 0, 0, 10828})
}

func testStartingPosDepth9() int8 {
	return test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 9, PerftNodes{2439530234167, 36095901903, /* 37101713	+ 5547231	 ? */
		125208536153, 319496827, 1784356000, 17334376, 400191963})
}

func testMiddleGameDepth1() int8 {
	return testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 1,
		PerftNodes{48, 0, 0, 0, 0, 0, 0})
}

func testMiddleGameDepth2() int8 {
	return testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 2,
		PerftNodes{2039, 0, 0, 0, 0, 0, 0})
}

func testMiddleGameDepth3() int8 {
	return testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 3,
		PerftNodes{97862, 0, 0, 0, 0, 0, 0})
}

func testMiddleGameDepth4() int8 {
	return testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 4,
		PerftNodes{4085603, 0, 0, 0, 0, 0, 0})
}

func testMiddleGameDepth5() int8 {
	return testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 5,
		PerftNodes{193690690, 0, 0, 0, 0, 0, 0})
}

func testMiddleGameDepth6() int8 {
	return testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 6,
		PerftNodes{8031647685, 0, 0, 0, 0, 0, 0})
}

func testEndGamePromotionDepth1() int8 {
	return testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 1,
		PerftNodes{24, 0, 0, 0, 0, 0, 0})
}

func testEndGamePromotionDepth2() int8 {
	return testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 2,
		PerftNodes{496, 0, 0, 0, 0, 0, 0})
}

func testEndGamePromotionDepth3() int8 {
	return testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 3,
		PerftNodes{9483, 0, 0, 0, 0, 0, 0})
}

func testEndGamePromotionDepth4() int8 {
	return testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 4,
		PerftNodes{182838, 0, 0, 0, 0, 0, 0})
}

func testEndGamePromotionDepth5() int8 {
	return testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 5,
		PerftNodes{3605103, 0, 0, 0, 0, 0, 0})
}

func testEndGamePromotionDepth6() int8 {
	return testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 6,
		PerftNodes{71179139, 0, 0, 0, 0, 0, 0})
}

func test(fen string, depth int, expected PerftNodes) int8 {
	fmt.Printf("Running perft for %s depth %d\n", fen, depth)
	g := FromFen(fen, true)
	actual := PerftNodes{0, 0, 0, 0, 0, 0, 0}
	perft(g.position, depth, NoType, 0, &actual)
	if actual != expected {
		fmt.Printf("test failed\nExpected: %d\nGot: %d\n", expected, actual)
		return 1
	}
	fmt.Println("Passed")
	return 0
}

func testNodesOnly(fen string, depth int, expected PerftNodes) int8 {
	fmt.Printf("Running perft for %s depth %d\n", fen, depth)
	g := FromFen(fen, true)
	actual := PerftNodes{0, 0, 0, 0, 0, 0, 0}
	perft(g.position, depth, NoType, 0, &actual)
	if actual.nodes != expected.nodes {
		fmt.Printf("test failed\nExpected: %d\nGot: %d\n", expected.nodes, actual.nodes)
		return 1
	}
	fmt.Println("Passed")
	return 0
}

func perft(p *Position, depth int, lastPromo PieceType, lastTag MoveTag, acc *PerftNodes) {
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
