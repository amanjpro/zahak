package uci

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/search"
)

const startFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var timerHalter chan bool

func UCI() {
	var game Game
	var depth = int8(100)
	reader := bufio.NewReader(os.Stdin)
	for true {
		cmd, err := reader.ReadString('\n')
		if err == nil {
			switch cmd {
			case "quit\n":
				return
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
				if timerHalter != nil {
					timerHalter <- true
				}
			default:
				if strings.HasPrefix(cmd, "go") {
					go findMove(game, depth, game.MoveClock(), cmd)
				} else if strings.HasPrefix(cmd, "position startpos moves") {
					moves := strings.Fields(cmd)[3:]
					game = FromFen(startFen, false)
					for _, move := range game.Position().ParseMoves(moves) {
						game.Move(move)
					}
				} else if strings.HasPrefix(cmd, "position startpos") {
					game = FromFen(startFen, true)
				} else if strings.HasPrefix(cmd, "position fen") {
					cmd := strings.Fields(cmd)
					fen := fmt.Sprintf("%s %s %s %s %s %s", cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7])
					moves := []string{}
					if len(cmd) > 9 {
						moves = cmd[9:]
						game = FromFen(fen, false)
					} else {
						game = FromFen(fen, true)
					}
					for _, move := range game.Position().ParseMoves(moves) {
						game.Move(move)
					}
				} else {
					fmt.Println("Didn't understand", cmd)
				}
			}
		}
	}
}

func findMove(game Game, depth int8, ply uint16, cmd string) {
	fields := strings.Fields(cmd)

	pos := game.Position()
	noTC := false
	timeToThink := 0
	inc := 0
	movesToGo := 0
	perMove := false
	for i := 0; i < len(fields); i++ {
		switch fields[i] {
		case "wtime":
			if pos.Turn() == White {
				timeToThink, _ = strconv.Atoi(fields[i+1])
				i++
			}
		case "btime":
			if pos.Turn() == Black {
				timeToThink, _ = strconv.Atoi(fields[i+1])
				i++
			}
		case "winc":
			if pos.Turn() == White {
				inc, _ = strconv.Atoi(fields[i+1])
				i++
			}
		case "binc":
			if pos.Turn() == Black {
				inc, _ = strconv.Atoi(fields[i+1])
				i++
			}
		case "movestogo":
			movesToGo, _ = strconv.Atoi(fields[i+1])
			i++
		case "depth":
			newPly, _ := strconv.Atoi(fields[i+1])
			ply = uint16(newPly)
			i++
		case "movetime":
			timeToThink, _ = strconv.Atoi(fields[i+1])
			perMove = true
			i++
		case "infinite":
			noTC = true
		}
	}

	if !noTC {
		timerHalter = make(chan bool)
		go InitiateTimer(&game, timeToThink, perMove, inc, movesToGo, timerHalter)
		evalMove := Search(game.Position(), depth, ply)
		fmt.Printf("bestmove %s\n", evalMove.Move().ToString())
		timerHalter <- true
	} else {
		evalMove := Search(game.Position(), depth, ply)
		fmt.Printf("bestmove %s\n", evalMove.Move().ToString())
	}
}
