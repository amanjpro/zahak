package cmd

import (
// "bufio"
// "fmt"
// "math/bits"
// "os"
// "strings"
//
// "github.com/notnil/chess"
)

func main() {
	// for i := 0; i < 8; i++ {
	// 	for j := 0; j < 8; j++ {
	// 		sq := SquareOf(File(j), Rank(i))
	// 		fmt.Print(sq.Name(), " ")
	// 	}
	// 	fmt.Println()
	// }
	// b := StartingBoard()
	// for sq, piece := range b.AllPieces() {
	// 	fmt.Println(sq.rank, sq.file)
	// 	fmt.Printf("In square %s there is %s\n", sq.Name(), piece.Name())
	// }
	// fmt.Println("HERE FEN IS: ", b.Fen())
	// b2 := bitboardFromFen(b.Fen())
	// fmt.Println(b2.Fen())
	uci()
	// g := FromFen("rnbqkbnr/pPp1pppp/4P3/3pP3/4p3/5BN1/PP3PPP/RNBQK2R w KQkq d6 0 1")
	// g := FromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	// for g.Status() == Unknown {
	// 	evalMove := search(g.position, 5)
	// 	g.Move(evalMove.move)
	// 	fmt.Println("Played: ", evalMove.move)
	// 	fmt.Println("Current eval is: ", evalMove.eval)
	// 	fmt.Println("Current tree is: ")
	// 	for _, mv := range evalMove.line {
	// 		fmt.Printf("%s ", mv.ToString())
	// 	}
	// 	// Position currently
	// 	fmt.Println(g.position.board.Draw())
	// }
	// for _, mv := range g.moves {
	// 	fmt.Printf("%s ", mv.ToString())
	// }
	// fmt.Println(g.position.board.Draw())
	// fmt.Println(g.halfMoveClock, g.positions, g.Status())
	// fmt.Println("HERE", g.position.enPassant)
	// gFen := g.Fen()
	// fmt.Println(gFen)
	// g2 := FromFen(gFen)
	// fmt.Println(g2.Fen())
	// fmt.Println(evalMove.eval)
	// fmt.Println(evalMove)

	// for i, m := range g.position.ValidMoves() {
	// 	fmt.Printf("%d %s\n", i, m.ToString())
	// }
	//
	// fmt.Println(bits.OnesCount64(empty))
	// fmt.Println(bits.OnesCount64(universal))
	//
	// fmt.Println(getIndicesOfOnes(g2.position.board.whitePawn))
	// fmt.Println(g.position.LegalMoves())
	//
	// fmt.Println("NOW", g.position.board.blackPieces, g.position.board.whitePieces)
	// for i, m := range g.position.LegalMoves() {
	// 	fmt.Printf("%d %s %d\n", i, m.ToString(), m.moveTag)
	// }
	// fmt.Println("NOW", g.position.board.blackPieces, g.position.board.whitePieces)
	// for i, m2 := range g.position.LegalMoves() {
	// 	fmt.Printf("%d %s %d\n", i, m2.ToString(), m2.moveTag)
	// }
	// fmt.Println(g2.position.board.Draw())
	//
	// game := chess.NewGame()
	// game.MoveStr("e4")
	// game.MoveStr("d5")
	// game.MoveStr("exd5")
	// game.MoveStr("Qxd5")
	// game.MoveStr("Nc3")
	// game.MoveStr("Qe5+")
	// game.MoveStr("Be2")
	// game.MoveStr("Bg4")
	// game.MoveStr("d4")
	// game.MoveStr("Bxe2")
	// game.MoveStr("Ngxe2")
	// game.MoveStr("Qd6")
	// game.MoveStr("Bf4")
	// game.MoveStr("Qd8")
	//
	// // game.MoveStr("e4")
	// // game.MoveStr("g6")
	// // game.MoveStr("d4")
	// // game.MoveStr("c6")
	// // game.MoveStr("Nf3")
	// // game.MoveStr("f6")
	// // game.MoveStr("Bc4")
	// // game.MoveStr("d5")
	// // game.MoveStr("exd5")
	// // game.MoveStr("cxd5")
	// // game.MoveStr("Bb5+")
	// // game.MoveStr("Bd7")
	// // game.MoveStr("Nc3")
	// // game.MoveStr("Bxb5")
	// // game.MoveStr("Nxb5")
	// // game.MoveStr("Bh6")
	// // game.MoveStr("O-O")
	// // game.MoveStr("Bf8")
	// // game.MoveStr("Bf4")
	// // game.MoveStr("Na6")
	// // game.MoveStr("c4")
	// // game.MoveStr("Qa5")
	// // game.MoveStr("Bd2")
	// // game.MoveStr("Qb6")
	// // game.MoveStr("Qa4")
	// // game.MoveStr("Kf7")
	// // game.MoveStr("Ba5")
	// // game.MoveStr("Qc6")
	// // game.MoveStr("Rac1")
	// // game.MoveStr("Bh6")
	// // game.MoveStr("Bd2")
	// // game.MoveStr("Bxd2")
	// // game.MoveStr("Nxd2")
	// // game.MoveStr("Nc7")
	// // game.MoveStr("Qb3")
	// // game.MoveStr("Nxb5")
	// // game.MoveStr("cxd5")
	// // game.MoveStr("Qd7")
	// // game.MoveStr("Rfe1")
	// // game.MoveStr("Nd6")
	// // game.MoveStr("Nc4")
	// // game.MoveStr("Rd8")
	// // game.MoveStr("Ne5+")
	// // game.MoveStr("fxe5")
	// // game.MoveStr("dxe5")
	// // game.MoveStr("Qb5")
	// // game.MoveStr("Qf3")
	// // game.MoveStr("Nf5")
	// // game.MoveStr("d6")
	// // game.MoveStr("Ke8")
	// // game.MoveStr("a4")
	// // game.MoveStr("Qxb2")
	// // game.MoveStr("Rb1")
	// // game.MoveStr("Qa2")
	// // game.MoveStr("Rxb7")
	// // game.MoveStr("Nxd6")
	// // game.MoveStr("exd6")
	// // Ke8 26.a4 Qxb2 27.Rb1 Qa2 28.Rxb7 Nxd6 29.exd6 e6 30.Qf7
	// // game.MoveStr("e5")
	// // game.MoveStr("h6")
	// // game.MoveStr("d4")
	// // game.MoveStr("h5")
	// // game.MoveStr("d5")
	// // game.MoveStr("Nb8")
	// // game.MoveStr("Bf4")
	// // game.MoveStr("Na6")
	// // game.MoveStr("Nf3")
	// // game.MoveStr("g6")
	// // game.MoveStr("Ng5")
	// // game.MoveStr("Nb8")
	// // game.MoveStr("Nc3")
	// // game.MoveStr("Na6")
	// // game.MoveStr("Qd3")
	// // game.MoveStr("Rh6")
	// // game.MoveStr("O-O-O")
	// // game.MoveStr("Rh8")
	// // game.MoveStr("Nxf7")
	// // game.MoveStr("Kxf7")
	// // game.MoveStr("e6")
	// // game.MoveStr("Kg7")
	// // game.MoveStr("Qd4+")
	// // .Qd4+ Kh7 14.Bd3 Bg7 15.Bxg6+ Kxg6 16.Qe4+ Kf6 17.Be5+ Kg5 18.h4+ Kh6 19.Bf4# 1-0
	//
	// // generate moves until game is over
	// // game.MoveStr("e4")
	// // game.MoveStr("c6")
	// // game.MoveStr("d4")
	// // game.MoveStr("Qa5+")
	// // game.MoveStr("Nc3")
	// // game.MoveStr("Qb4")
	// // game.MoveStr("Nf3")
	// // game.MoveStr("e6")
	// // game.MoveStr("d5")
	// // game.MoveStr("exd5")
	// // game.MoveStr("exd5")
	// // game.MoveStr("Be7")
	// // game.MoveStr("Bd2")
	// // game.MoveStr("c5")
	// // game.MoveStr("Be3")
	// // game.MoveStr("Na6")
	// // game.MoveStr("d6")
	// // game.MoveStr("Bf6")
	// // game.MoveStr("Bxa6")
	// // game.MoveStr("bxa6")
	// // game.MoveStr("O-O")
	// // game.MoveStr("Qxb2")
	// // game.MoveStr("Nd5")
	// // game.MoveStr("Qxa1")
	// // game.MoveStr("Nxf6+")
	// // game.MoveStr("Qxf6")
	// // game.MoveStr("Bg5")
	// // game.MoveStr("Qe6")
	// // game.MoveStr("Re1")
	// // game.MoveStr("Qxe1+")
	// // game.MoveStr("Qxe1+")
	// // game.MoveStr("Kf8")
	// // game.MoveStr("Qa5")
	// // game.MoveStr("f6")
	// // game.MoveStr("Qd8+")
	// // game.MoveStr("Kf7")
	// // game.MoveStr("Be3")
	// // game.MoveStr("c4")
	// // game.MoveStr("Nd4")
	// // game.MoveStr("Rb8")
	// fmt.Println(game.Position().Board().Draw())
	// for game.Outcome() == chess.NoOutcome {
	// 	// select a random move
	// 	if game.Position().Turn() == chess.White {
	// 		reader := bufio.NewReader(os.Stdin)
	// 		fmt.Print("Enter your next move: ")
	// 		move, _ := reader.ReadString('\n')
	// 		err := game.MoveStr(strings.TrimSpace(move))
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		fmt.Println("Current eval is: ", eval(game.Position()))
	// 	} else {
	// 		evalMove := search(game.Position(), 5)
	// 		game.Move(evalMove.move)
	// 		fmt.Println("Played: ", evalMove.move)
	// 		fmt.Println("Current eval is: ", evalMove.eval)
	// 		fmt.Println("Current tree is: ")
	// 		for _, mv := range evalMove.line {
	// 			fmt.Printf("%s ", mv.String())
	// 		}
	// 	}
	// 	// Position currently
	// 	fmt.Println(game.Position().Board().Draw())
	// }
	// fmt.Println(game.Position().Turn())
	// fmt.Printf("Game completed. %s by %s.\n", game.Outcome(), game.Method())
	// fmt.Println(game.String())
	// /*
	// 	Output:
	//
	// 	 A B C D E F G H
	// 	8- - - - - - - -
	// 	7- - - - - - ♚ -
	// 	6- - - - ♗ - - -
	// 	5- - - - - - - -
	// 	4- - - - - - - -
	// 	3♔ - - - - - - -
	// 	2- - - - - - - -
	// 	1- - - - - - - -
	//
	// 	Game completed. 1/2-1/2 by InsufficientMaterial.
	//
	// 	1.Nc3 b6 2.a4 e6 3.d4 Bb7 ...
	// */
}
