package perft

import (
	"fmt"
	"os"
	"time"

	. "github.com/amanjpro/zahak/engine"
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

func PerftTree(game Game, depth int, strMoves []string) {
	sum := int64(0)
	game.ParseGameMoves(strMoves)
	if depth > 0 {
		depth -= 1
		cache = make([]map[uint64]int64, depth)
		for i := 0; i < depth; i++ {
			cache[i] = make(map[uint64]int64, 1000_000)
		}
		moves := game.Position().PseudoLegalMoves()
		for _, move := range moves {
			if ep, tg, hc, ok := game.Position().MakeMove(move); ok {
				nodes := bulkyPerft(game.Position(), depth)
				fmt.Printf("%s %d\n", move.ToString(), nodes)
				sum += nodes
				game.Position().UnMakeMove(move, tg, ep, hc)
			}
		}
	}

	fmt.Printf("\n%d\n", sum)
}

func StartPerftTest(slow bool) {
	result := int8(0)
	if !slow {
		result += test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 0,
			PerftNodes{1, 0, 0, 0, 0, 0, 0})
		result += test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 1,
			PerftNodes{20, 0, 0, 0, 0, 0, 0})
		result += test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 2,
			PerftNodes{400, 0, 0, 0, 0, 0, 0})
		result += test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 3,
			PerftNodes{8902, 12, 34, 0, 0, 0, 0})
		result += test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 4,
			PerftNodes{197281, 469, 1576, 0, 0, 0, 8})
		result += test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 5,
			PerftNodes{4865609, 27351 /* 27351 + 6 */, 82719, 258, 0, 0, 347})
		result += testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 1,
			48)
		result += testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 2,
			2039)
		result += testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 3,
			97862)
		result += testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 4,
			4085603)
		result += testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 1, 24)
		result += testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 2, 496)
		result += testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 3, 9483)
		result += testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 4, 182838)

		// initial position
		result += testNodesOnly("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 6, 119060324)

		result += testNodesOnly("rnb1kbnr/pp1pp1pp/1qp2p2/8/Q1P5/N7/PP1PPPPP/1RB1KBNR b Kkq - 2 4", 7, 14794751816)

		// old positions
		result += testNodesOnly("1Q5Q/8/3p1p2/2RpkpR1/2PpppP1/2QPQPB1/8/4K3 b - - 0 1", 7, 290063345)
		result += testNodesOnly("1Q5Q/8/3p1p2/2RpkpR1/2PpppP1/2QPQPB1/8/4K3 b - - 0 1", 8, 17665826996)
		result += testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 5, 193690690)
		result += testNodesOnly("8/PPP4k/8/8/8/8/4Kppp/8 w - - 0 1", 6, 34336777)
		result += testNodesOnly("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1", 7, 178633661)
		result += testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 6, 8031647685)
		result += testNodesOnly("r3r1k1/1pq2pp1/2p2n2/1PNn4/2QN2b1/6P1/3RPP2/2R3KB b - - 0 1", 6, 15097513050)
		result += testNodesOnly("8/p1p1p3/8/1P1P2k1/1K2p1p1/8/3P1P1P/8 w - - 0 1", 7, 118590233)

		// new positions
		result += testNodesOnly("r3k2r/8/8/8/3pPp2/8/8/R3K1RR b KQkq e3 0 1", 6, 485647607)
		result += testNodesOnly("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1", 6, 706045033)
		result += testNodesOnly("8/7p/p5pb/4k3/P1pPn3/8/P5PP/1rB2RK1 b - d3 0 28", 6, 38633283)
		result += testNodesOnly("8/3K4/2p5/p2b2r1/5k2/8/8/1q6 b - - 1 67", 7, 493407574)
		result += testNodesOnly("rnbqkb1r/ppppp1pp/7n/4Pp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 3", 6, 244063299)
		result += testNodesOnly("8/p7/8/1P6/K1k3p1/6P1/7P/8 w - - 0 1", 8, 8103790)
		result += testNodesOnly("r3k2r/p6p/8/B7/1pp1p3/3b4/P6P/R3K2R w KQkq - 0 1", 6, 77054993)
		result += testNodesOnly("8/5p2/8/2k3P1/p3K3/8/1P6/8 b - - 0 1", 8, 64451405)
		result += testNodesOnly("r3k2r/pb3p2/5npp/n2p4/1p1PPB2/6P1/P2N1PBP/R3K2R w KQkq - 0 1", 5, 29179893)
		result += testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 4, 4085603)
		result += testNodesOnly("rnbqkbnr/ppN5/3ppppp/8/B3P3/8/PPPP1PPP/R1BQK1NR b KQkq - 0 1", 7, 1029536265)
		result += testNodesOnly("k7/8/8/8/8/8/3b1P2/4K3 w - - 0 1", 10, 1555125531)
		result += testNodesOnly("r3k2r/1bp2pP1/5n2/1P1Q4/1pPq4/5N2/1B1P2p1/R3K2R b KQkq c3 0 1", 6, 8419356881)
		result += testNodesOnly("rrrrkr1r/rr1rr3/8/8/8/8/8/6KR w k - 0 1", 7, 1941335854)

		// avoid illegal ep (thanks to Steve Maughan):
		result += testNodesOnly("3k4/3p4/8/K1P4r/8/8/8/8 b - - 0 1", 6, 1134888)
		result += testNodesOnly("8/8/8/8/k1p4R/8/3P4/3K4 w - - 0 1", 6, 1134888)

		// avoid illegal ep #2
		result += testNodesOnly("8/8/4k3/8/2p5/8/B2P2K1/8 w - - 0 1", 6, 1015133)
		result += testNodesOnly("8/b2p2k1/8/2P5/8/4K3/8/8 b - - 0 1", 6, 1015133)

		// en passant capture checks opponent:
		result += testNodesOnly("8/8/1k6/2b5/2pP4/8/5K2/8 b - d3 0 1", 6, 1440467)
		result += testNodesOnly("8/5k2/8/2Pp4/2B5/1K6/8/8 w - d6 0 1", 6, 1440467)

		// short castling gives check:
		result += testNodesOnly("5k2/8/8/8/8/8/8/4K2R w K - 0 1", 6, 661072)
		result += testNodesOnly("4k2r/8/8/8/8/8/8/5K2 b k - 0 1", 6, 661072)

		// long castling gives check:
		result += testNodesOnly("3k4/8/8/8/8/8/8/R3K3 w Q - 0 1", 6, 803711)
		result += testNodesOnly("r3k3/8/8/8/8/8/8/3K4 b q - 0 1", 6, 803711)

		// castling (including losing cr due to rook capture):
		result += testNodesOnly("r3k2r/1b4bq/8/8/8/8/7B/R3K2R w KQkq - 0 1", 4, 1274206)
		result += testNodesOnly("r3k2r/7b/8/8/8/8/1B4BQ/R3K2R b KQkq - 0 1", 4, 1274206)

		// castling prevented:
		result += testNodesOnly("r3k2r/8/3Q4/8/8/5q2/8/R3K2R b KQkq - 0 1", 4, 1720476)
		result += testNodesOnly("r3k2r/8/5Q2/8/8/3q4/8/R3K2R w KQkq - 0 1", 4, 1720476)

		// promote out of check:
		result += testNodesOnly("2K2r2/4P3/8/8/8/8/8/3k4 w - - 0 1", 6, 3821001)
		result += testNodesOnly("3K4/8/8/8/8/8/4p3/2k2R2 b - - 0 1", 6, 3821001)

		// discovered check:
		result += testNodesOnly("8/8/1P2K3/8/2n5/1q6/8/5k2 b - - 0 1", 5, 1004658)
		result += testNodesOnly("5K2/8/1Q6/2N5/8/1p2k3/8/8 w - - 0 1", 5, 1004658)

		// promote to give check:
		result += testNodesOnly("4k3/1P6/8/8/8/8/K7/8 w - - 0 1", 6, 217342)
		result += testNodesOnly("8/k7/8/8/8/8/1p6/4K3 b - - 0 1", 6, 217342)

		// underpromote to check:
		result += testNodesOnly("8/P1k5/K7/8/8/8/8/8 w - - 0 1", 6, 92683)
		result += testNodesOnly("8/8/8/8/8/k7/p1K5/8 b - - 0 1", 6, 92683)

		// self stalemate:
		result += testNodesOnly("K1k5/8/P7/8/8/8/8/8 w - - 0 1", 6, 2217)
		result += testNodesOnly("8/8/8/8/8/p7/8/k1K5 b - - 0 1", 6, 2217)

		// stalemate/checkmate:
		result += testNodesOnly("8/k1P5/8/1K6/8/8/8/8 w - - 0 1", 7, 567584)
		result += testNodesOnly("8/8/8/8/1k6/8/K1p5/8 b - - 0 1", 7, 567584)

		// double check:
		result += testNodesOnly("8/8/2k5/5q2/5n2/8/5K2/8 b - - 0 1", 4, 23527)
		result += testNodesOnly("8/5k2/8/5N2/5Q2/2K5/8/8 w - - 0 1", 4, 23527)

		// capture-to-square
		result += testNodesOnly("3k4/8/8/2Pp3r/2K5/8/8/8 w - d6 0 1", 2, 112)
		result += testNodesOnly("1RR4K/3P4/8/8/8/8/3p4/4rr1k w - - 0 1", 6, 419523239)
		result += testNodesOnly("1RR4K/3P4/8/8/8/8/3p4/4rr1k b - - 1 1", 6, 395340738)
		result += testNodesOnly("1RR5/7K/3P4/8/8/3p4/7k/4rr2 w - - 0 1", 6, 310492012)
		result += testNodesOnly("1RR5/7K/3P4/8/8/3p4/7k/4rr2 b - - 1 1", 6, 302653359)

		// somewhat slower than the others
		result += testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 5,
			193690690)
		result += testNodesOnly("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 6,
			8031647685)
		result += testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 5, 3605103)
		result += testNodesOnly("rnb1kbnr/ppp1pppp/8/3p4/1P6/P2P3q/2P1PPP1/RNBQKBNR b KQkq - 0 4", 7, 44950307154)
	} else {
		result += testNodesOnly("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 7, 3195901860)
		result += testNodesOnly("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 6, 71179139)
		result += test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 6,
			PerftNodes{119060324, 809099, 2812008, 5248, 0, 0, 10828})
		result += test("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 9,
			PerftNodes{2439530234167, 36095901903,
				125208536153, 319496827, 1784356000, 17334376, 400191963})
	}

	if result > 0 {
		fmt.Printf("%d perft tests failed... hard luck\n", result)
		os.Exit(1)
	} else {
		fmt.Println("All perft tests passed")
	}
}

func test(fen string, depth int, expected PerftNodes) int8 {
	fmt.Printf("Running perft for %s depth %d\n", fen, depth)
	g := FromFen(fen)
	actual := PerftNodes{0, 0, 0, 0, 0, 0, 0}
	perft(g.Position(), depth, EmptyMove, &actual)
	if actual != expected {
		fmt.Printf("test failed\nExpected: %d\nGot: %d\n", expected, actual)
		return 1
	}
	fmt.Println("Passed")
	return 0
}

var cache []map[uint64]int64

func testNodesOnly(fen string, depth int, expected int64) int8 {
	fmt.Printf("Running perft for %s depth %d\n", fen, depth)
	g := FromFen(fen)
	cache = make([]map[uint64]int64, depth)
	for i := 0; i < depth; i++ {
		cache[i] = make(map[uint64]int64, 1000_000)
	}
	start := time.Now()
	actual := bulkyPerft(g.Position(), depth)
	end := time.Now()
	fmt.Printf("mnsp is %f\n", float64(actual)/(1000_000*(end.Sub(start).Seconds())))
	if actual != expected {
		fmt.Printf("test failed\nExpected: %d\nGot: %d\n", expected, actual)
		return 1
	}
	fmt.Println("Passed")
	return 0
}

func perft(p *Position, depth int, currentMove Move, acc *PerftNodes) {
	if depth == 0 {
		isCheck := p.IsInCheck()
		moves := p.PseudoLegalMoves()
		hasLegalMoves := false
		for _, move := range moves {
			if ep, tag, hc, ok := p.MakeMove(move); ok {
				hasLegalMoves = true
				p.UnMakeMove(move, tag, ep, hc)
				break
			}
		}
		isCheckmate := isCheck && !hasLegalMoves
		acc.nodes += 1
		if isCheckmate {
			acc.checkmates += 1
		}
		if isCheck {
			acc.checks += 1
		}

		if currentMove.IsCapture() {
			acc.captures += 1
		}
		if currentMove.IsEnPassant() {
			acc.enPassants += 1
		}
		if currentMove.IsQueenSideCastle() || currentMove.IsKingSideCastle() {
			acc.castles += 1
		}
		if currentMove.PromoType() != NoType {
			acc.promotions += 1
		}
		return
	}

	moves := p.PseudoLegalMoves()

	for _, move := range moves {
		if ep, tag, hc, ok := p.MakeMove(move); ok {
			perft(p, depth-1, move, acc)
			p.UnMakeMove(move, tag, ep, hc)
		}
	}
}

func bulkyPerft(p *Position, depth int) int64 {
	nodes := int64(0)

	if depth == 0 {
		return 1
	}

	moves := p.PseudoLegalMoves()
	if depth == 1 {
		count := int64(0)
		for _, move := range moves {
			if ep, tag, hc, ok := p.MakeMove(move); ok {
				count += 1
				p.UnMakeMove(move, tag, ep, hc)
			}
		}
		return count
	}

	for _, move := range moves {
		if ep, tag, hc, ok := p.MakeMove(move); ok {
			hash := p.Hash()
			n, ok := cache[depth-1][hash]
			if ok {
				nodes += n
			} else {
				n := bulkyPerft(p, depth-1)
				cache[depth-1][hash] = n
				nodes += n
			}
			p.UnMakeMove(move, tag, ep, hc)
		}
	}
	return nodes
}
