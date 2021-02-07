package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const startFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func uci() {
	var game Game
	var depth = int8(5)
	reader := bufio.NewReader(os.Stdin)
	for true {
		cmd, err := reader.ReadString('\n')
		if err == nil {
			switch cmd {
			case "uci\n":
				fmt.Println("id name Zahak\n")
				fmt.Println("id author Amanj\n")
				fmt.Println("uciok\n")
			case "isready\n":
				fmt.Println("readyok\n")
			case "ucinewgame\n":
				game = FromFen(startFen)
			case "stop\n":
				STOP_SEARCH_GLOBALLY = true
			default:
				if strings.HasPrefix(cmd, "go") {
					go findMove(game.position, depth)
				} else if strings.HasPrefix(cmd, "position startpos moves") {
					moves := strings.Fields(cmd)
					game = FromFen(startFen)
					for i := 3; i < len(moves); i++ {
						move := moves[i]
						source := NameToSquareMap[string([]byte{move[0], move[1]})]
						destination := NameToSquareMap[string([]byte{move[2], move[3]})]
						promo := NoType
						if len(move) == 5 {
							p := pieceFromName(rune(move[4]))
							promo = p.Type()
						}

						var tag MoveTag = 0
						if destination.File() != source.File() {
							pos := game.position
							board := pos.board
							movingPiece := board.PieceAt(source)
							destPiece := board.PieceAt(destination)
							if movingPiece.Type() == Pawn && destPiece == NoPiece {
								tag = EnPassant
							}
						}

						game.Move(Move{source, destination, promo, tag})
					}
				} else if strings.HasPrefix(cmd, "position startpos") {
					game = FromFen(startFen)
				} else if strings.HasPrefix(cmd, "position") {
					cmd := strings.Fields(cmd)
					fen := fmt.Sprintf("%s %s %s %s %s %s", cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6])
					game = FromFen(fen)
					for i := 8; i < len(cmd); i++ {
						move := cmd[i]
						source := NameToSquareMap[fmt.Sprintf("%x%x", move[0], move[1])]
						destination := NameToSquareMap[fmt.Sprintf("%x%x", move[2], move[3])]
						promo := NoType
						if len(move) == 5 {
							p := pieceFromName(rune(move[4]))
							promo = p.Type()
						}

						var tag MoveTag = 0
						if destination.File() != source.File() {
							pos := game.position
							board := pos.board
							movingPiece := board.PieceAt(source)
							destPiece := board.PieceAt(destination)
							if movingPiece.Type() == Pawn && destPiece == NoPiece {
								tag = EnPassant
							}
						}

						game.Move(Move{source, destination, promo, tag})
					}

				} else {
					fmt.Println("Didn't understand", cmd)
				}
			}
		}
	}
}

func findMove(pos *Position, depth int8) {
	evalMove := search(pos, depth)
	pos.MakeMove(evalMove.move)
	fmt.Printf("bestmove %s\n", evalMove.move.ToString())
}
