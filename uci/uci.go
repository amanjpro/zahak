package uci

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/amanjpro/zahak/book"
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/search"
)

const startFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

type UCI struct {
	engine    *Engine
	thinkTime int64
	withBook  bool
	bookPath  string
}

func NewUCI(withBook bool, bookPath string) *UCI {
	return &UCI{
		NewEngine(),
		0,
		withBook,
		bookPath,
	}
}

func (uci *UCI) Start() {
	var game Game
	var depth = int8(100)
	if uci.withBook {
		InitBook(uci.bookPath)
	}
	reader := bufio.NewReader(os.Stdin)
	for true {
		cmd, err := reader.ReadString('\n')
		cmd = strings.Trim(cmd, "\n\r")
		if err == nil {
			switch cmd {
			case "debug on":
				uci.engine.DebugMode = true
			case "debug off":
				uci.engine.DebugMode = false
			case "ponderhit":
				uci.engine.StartTime = time.Now()
				uci.engine.Pondering = false
				uci.engine.ThinkTime = uci.thinkTime
				uci.thinkTime = 0
				uci.engine.SendPv(-1)
			case "quit":
				return
			case "uci":
				fmt.Print("id name Zahak\n")
				fmt.Print("id author Amanj\n")
				fmt.Print("option name Ponder type check default false\n")
				fmt.Printf("option name Hash type spin default %d min 1 max %d\n", DEFAULT_CACHE_SIZE, MAX_CACHE_SIZE)
				fmt.Printf("option name Book type check default %t\n", uci.withBook)
				fmt.Print("uciok\n")
			case "isready":
				fmt.Print("readyok\n")
			case "ucinewgame":
				game = FromFen(startFen, true)
			case "stop":
				if uci.engine.Pondering {
					uci.stopPondering()
				} else {
					uci.engine.StopSearchFlag = true
				}
			default:
				if strings.HasPrefix(cmd, "setoption name Ponder value") {
					continue
				} else if strings.HasPrefix(cmd, "setoption name Book value ") {
					options := strings.Fields(cmd)
					opt := options[len(options)-1]
					if opt == "false" {
						ResetBook()
					} else if !IsBoookLoaded() && opt == "true" { // if it is loaded, no need to reload
						InitBook(uci.bookPath)
					}
				} else if strings.HasPrefix(cmd, "setoption name Hash value") {
					options := strings.Fields(cmd)
					mg := options[len(options)-1]
					hashSize, _ := strconv.Atoi(mg)
					NewCache(uint32(hashSize))
				} else if strings.HasPrefix(cmd, "go") {
					go uci.findMove(game, depth, game.MoveClock(), cmd)
				} else if strings.HasPrefix(cmd, "position startpos moves") {
					uci.stopPondering()
					moves := strings.Fields(cmd)[3:]
					game = FromFen(startFen, false)
					for _, move := range game.Position().ParseMoves(moves) {
						game.Move(move)
					}
				} else if strings.HasPrefix(cmd, "position startpos") {
					uci.stopPondering()
					game = FromFen(startFen, true)
				} else if strings.HasPrefix(cmd, "position fen") {
					uci.stopPondering()
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

func (uci *UCI) findMove(game Game, depth int8, ply uint16, cmd string) {
	uci.thinkTime = 0
	fields := strings.Fields(cmd)

	pos := game.Position()
	noTC := false
	timeToThink := 0
	inc := 0
	movesToGo := 0
	perMove := false
	uci.engine.Pondering = false
	for i := 0; i < len(fields); i++ {
		switch fields[i] {
		case "ponder":
			uci.engine.Pondering = true
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
			depth = int8(newPly)
			i++
		case "movetime":
			timeToThink, _ = strconv.Atoi(fields[i+1])
			perMove = true
			i++
		case "infinite":
			noTC = true
		}
	}

	var MAX_TIME int64 = 9_223_372_036_854_775_807

	if !noTC {
		if uci.engine.Pondering {
			uci.thinkTime = uci.engine.InitiateTimer(&game, timeToThink, perMove, inc, movesToGo)
			uci.engine.ThinkTime = MAX_TIME
		} else {
			uci.engine.ThinkTime = uci.engine.InitiateTimer(&game, timeToThink, perMove, inc, movesToGo)
		}
		uci.engine.Search(game.Position(), depth, ply)
		uci.engine.SendBestMove()
		uci.engine.Pondering = false
	} else {
		uci.engine.ThinkTime = MAX_TIME
		uci.thinkTime = uci.engine.ThinkTime
		uci.engine.Search(game.Position(), depth, ply)
		uci.engine.SendBestMove()
		uci.engine.Pondering = false
	}
}

func (uci *UCI) stopPondering() {
	if uci.engine.Pondering {
		uci.engine.StopSearchFlag = true
		for !uci.engine.Pondering {
		} // wait until stopped
	}
}
