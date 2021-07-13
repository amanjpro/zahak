package uci

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	. "github.com/amanjpro/zahak/book"
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
	. "github.com/amanjpro/zahak/search"
)

const startFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

type UCI struct {
	version     string
	engine      *Engine
	timeManager *TimeManager
	withBook    bool
	bookPath    string
}

func NewUCI(version string, withBook bool, bookPath string) *UCI {
	return &UCI{
		version,
		NewEngine(NewCache(DEFAULT_CACHE_SIZE)),
		nil,
		withBook,
		bookPath,
	}
}

func (uci *UCI) Start() {
	var game Game
	var depth = int8(MAX_DEPTH)
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
				uci.engine.AttachTimeManager(uci.timeManager)
				uci.timeManager = nil
				uci.engine.SendPv(-1)
			case "quit":
				return
			case "eval":
				dir := int16(1)
				if game.Position().Turn() == Black {
					dir = -1
				}
				fmt.Printf("%d\n", dir*Evaluate(game.Position()))
			case "uci":
				fmt.Printf("id name Zahak %s\n", uci.version)
				fmt.Print("id author Amanj\n")
				fmt.Print("option name Ponder type check default false\n")
				fmt.Printf("option name Hash type spin default %d min 1 max %d\n", DEFAULT_CACHE_SIZE, MAX_CACHE_SIZE)
				fmt.Printf("option name Pawnhash type spin default %d min 1 max %d\n", DEFAULT_PAWNHASH_SIZE, MAX_PAWNHASH_SIZE)
				fmt.Printf("option name Book type check default %t\n", uci.withBook)
				fmt.Print("uciok\n")
			case "isready":
				fmt.Print("readyok\n")
			case "isdraw":
				fmt.Print(game.Position().IsDraw(), "\n")
			case "draw":
				fmt.Print(game.Position().Board.Draw(), "\n")
			case "ucinewgame", "position startpos":
				size := uci.engine.TranspositionTable.Size()
				pawnSize := Pawnhash.Size()
				Pawnhash = nil
				uci.engine.TranspositionTable = nil
				runtime.GC()
				uci.engine.TranspositionTable = NewCache(size)
				Pawnhash = NewPawnCache(pawnSize)
				game = FromFen(startFen, true)
			case "stop":
				if uci.engine.Pondering {
					uci.stopPondering()
				} else {
					if uci.engine.TimeManager != nil {
						uci.engine.TimeManager.StopSearchNow = true
					}
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
				} else if strings.HasPrefix(cmd, "setoption name Pawnhash value") {
					options := strings.Fields(cmd)
					mg := options[len(options)-1]
					hashSize, _ := strconv.Atoi(mg)
					Pawnhash = nil
					runtime.GC()
					Pawnhash = NewPawnCache(hashSize)
				} else if strings.HasPrefix(cmd, "setoption name Hash value") {
					options := strings.Fields(cmd)
					mg := options[len(options)-1]
					hashSize, _ := strconv.Atoi(mg)
					uci.engine.TranspositionTable = nil
					runtime.GC()
					uci.engine.TranspositionTable = NewCache(uint32(hashSize))
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
					game = FromFen(startFen, false)
				} else if strings.HasPrefix(cmd, "position fen") {
					uci.stopPondering()
					cmd := strings.Fields(cmd)
					var fen string
					if len(cmd) < 8 {
						fen = fmt.Sprintf("%s %s %s %s %d %d", cmd[2], cmd[3], cmd[4], cmd[5], 0, 1)
					} else {
						fen = fmt.Sprintf("%s %s %s %s %s %s", cmd[2], cmd[3], cmd[4], cmd[5], cmd[6], cmd[7])
					}
					moves := []string{}
					if len(cmd) > 9 {
						moves = cmd[9:]
						game = FromFen(fen, false)
					} else {
						game = FromFen(fen, false)
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
	uci.timeManager = nil
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
			noTC = true
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

	uci.engine.Position = game.Position()
	uci.engine.Ply = ply
	if !noTC {
		if uci.engine.Pondering {
			uci.timeManager = NewTimeManager(time.Now(), int64(timeToThink), perMove, int64(inc), int64(movesToGo))
			uci.engine.InitTimeManager(MAX_TIME, false, 0, 0)
		} else {
			uci.engine.InitTimeManager(int64(timeToThink), perMove, int64(inc), int64(movesToGo))
		}
		uci.engine.Search(depth)
		uci.engine.SendBestMove()
		uci.engine.Pondering = false
	} else {
		uci.engine.InitTimeManager(MAX_TIME, false, 0, 0)
		uci.timeManager = uci.engine.TimeManager
		uci.engine.Search(depth)
		uci.engine.SendBestMove()
		uci.engine.Pondering = false
	}
}

func (uci *UCI) stopPondering() {
	if uci.engine.Pondering {
		uci.engine.TimeManager.StopSearchNow = true
		for !uci.engine.Pondering {
		} // wait until stopped
	}
}
