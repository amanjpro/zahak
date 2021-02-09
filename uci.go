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
	var depth = int8(7)
	reader := bufio.NewReader(os.Stdin)
	for true {
		cmd, err := reader.ReadString('\n')
		if err == nil {
			switch cmd {
			case "quit\n":
				os.Exit(0)
			case "uci\n":
				fmt.Print("id name Zahak\n\n")
				fmt.Print("id author Amanj\n\n")
				fmt.Print("uciok\n\n")
			case "isready\n":
				fmt.Print("readyok\n\n")
			case "ucinewgame\n":
				game = FromFen(startFen, true)
			case "stop\n":
				STOP_SEARCH_GLOBALLY = true
			default:
				if strings.HasPrefix(cmd, "go") {
					go findMove(game.position, depth)
				} else if strings.HasPrefix(cmd, "position startpos moves") {
					moves := strings.Fields(cmd)[3:]
					game = FromFen(startFen, true)
					for _, move := range game.position.ParseMoves(moves) {
						game.Move(move)
					}
				} else if strings.HasPrefix(cmd, "position startpos") {
					game = FromFen(startFen, false)
				} else if strings.HasPrefix(cmd, "position") {
					cmd := strings.Fields(cmd)
					moves := cmd[8:]
					fen := fmt.Sprintf("%s %s %s %s %s %s", cmd[1], cmd[2], cmd[3], cmd[4], cmd[5], cmd[6])
					game = FromFen(fen, false)
					for _, move := range game.position.ParseMoves(moves) {
						game.Move(move)
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
