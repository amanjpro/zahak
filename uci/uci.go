package uci

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	. "github.com/amanjpro/zahak/book"
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/fathom"
	. "github.com/amanjpro/zahak/search"
)

const startFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var defaultCPU = 1
var minCPU = 1
var maxCPU = runtime.NumCPU()

type UCI struct {
	version     string
	runner      *Runner
	timeManager *TimeManager
	withBook    bool
	bookPath    string
}

func NewUCI(version string, withBook bool, bookPath string) *UCI {
	return &UCI{
		version,
		NewRunner(NewCache(DEFAULT_CACHE_SIZE), 1),
		nil,
		withBook,
		bookPath,
	}
}

func (uci *UCI) Start() {
	var game Game = FromFen(startFen)
	var depth = int8(MAX_DEPTH)
	if uci.withBook {
		InitBook(uci.bookPath)
	}
	reader := bufio.NewReader(os.Stdin)
	defer os.Stdin.Close()

	for true {
		cmd, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		cmd = strings.Trim(cmd, "\n\r")
		if err == nil {
			switch cmd {
			case "debug on":
				uci.runner.DebugMode = true
			case "debug off":
				uci.runner.DebugMode = false
			case "ponderhit":
				uci.runner.Ponderhit()
				uci.timeManager = nil
			case "quit":
				return
			case "eval":
				game.Position().Net.Recalculate(game.Position().NetInput())
				fmt.Printf("%d\n", game.Position().Evaluate())
			case "uci":
				fmt.Printf("id name Zahak %s\n", uci.version)
				fmt.Print("id author Amanj\n")
				fmt.Printf("id EvalFile %d\n", CurrentNetworkId)
				fmt.Print("option name Ponder type check default false\n")
				fmt.Printf("option name Hash type spin default %d min 1 max %d\n", DEFAULT_CACHE_SIZE, MAX_CACHE_SIZE)
				fmt.Printf("option name Book type check default %t\n", uci.withBook)
				fmt.Printf("option name Threads type spin default %d min %d max %d\n", defaultCPU, minCPU, maxCPU)
				fmt.Print("option name EvalFile type string default <empty>\n")
				fmt.Print("option name BookFile type string default <empty>\n")
				fmt.Print("option name SyzygyPath type string default <empty>\n")
				fmt.Printf("option name SyzygyProbeDepth type spin default %d min 0 max 128\n", DefaultProbeDepth)
				fmt.Print("uciok\n")
			case "isready":
				fmt.Print("readyok\n")
			case "isdraw":
				fmt.Print(game.Position().IsDraw(), "\n")
			case "draw":
				fmt.Print(game.Position().Board.Draw(), "\n")
			case "ucinewgame", "position startpos":
				size := uci.runner.Engines[0].TranspositionTable.Size()
				newTT := NewCache(size)
				for i := 0; i < len(uci.runner.Engines); i++ {
					uci.runner.Engines[i].TranspositionTable = nil
					runtime.GC()
					uci.runner.Engines[i].TranspositionTable = newTT
				}
				game = FromFen(startFen)
			case "stop":
				if uci.runner.TimeManager != nil {
					if uci.runner.TimeManager.Pondering {
						uci.stopPondering()
					} else {
						uci.runner.TimeManager.StopSearchNow = true
					}
				}
			default:
				if strings.HasPrefix(cmd, "tb-probe") {
					fmt.Println(ProbeWDL(game.Position(), 0))
					dtz := ProbeDTZ(game.Position())
					fmt.Println(dtz.ToString())
				} else if strings.HasPrefix(cmd, "setoption name SyzygyProbeDepth value") {
					options := strings.Fields(cmd)
					v := options[len(options)-1]
					depth, _ := strconv.Atoi(v)
					MinProbeDepth = int8(depth)
				} else if strings.HasPrefix(cmd, "setoption name SyzygyPath value") {
					path := strings.TrimSpace(strings.ReplaceAll(cmd, "setoption name SyzygyPath value", ""))
					if path == "" || path == "<empty>" {
						ClearSyzygy()
					}
					fmt.Println(path)
					SetSyzygyPath(path)
				} else if strings.HasPrefix(cmd, "setoption name EvalFile value") {

					path := strings.TrimSpace(strings.ReplaceAll(cmd, "setoption name EvalFile value", ""))
					if path == "" || path == "<empty>" {
						fmt.Print("info string no eval file is selected, ignoring\n")
						continue
					}
					LoadNetwork(path)
					fmt.Printf("info string new EvalFile loaded, the id of the new EvalFile is %d\n", CurrentNetworkId)
				} else if strings.HasPrefix(cmd, "setoption name BookFile value") {
					path := strings.TrimSpace(strings.ReplaceAll(cmd, "setoption name BookFile value", ""))
					if path == "" || path == "<empty>" {
						fmt.Print("info string no eval file is selected, ignoring\n")
						continue
					}
					uci.bookPath = path
					if uci.withBook {
						InitBook(uci.bookPath)
					}
				} else if strings.HasPrefix(cmd, "setoption name Ponder value") {
					continue
				} else if strings.HasPrefix(cmd, "setoption name Book value ") {
					options := strings.Fields(cmd)
					opt := options[len(options)-1]
					if opt == "false" {
						ResetBook()
					} else if !IsBoookLoaded() && opt == "true" { // if it is loaded, no need to reload
						uci.withBook = true
						InitBook(uci.bookPath)
					}
				} else if strings.HasPrefix(cmd, "setoption name Threads value") {
					options := strings.Fields(cmd)
					v := options[len(options)-1]
					cpu, _ := strconv.Atoi(v)
					uci.runner = NewRunner(uci.runner.Engines[0].TranspositionTable, cpu)
				} else if strings.HasPrefix(cmd, "setoption name Hash value") {
					options := strings.Fields(cmd)
					mg := options[len(options)-1]
					hashSize, _ := strconv.Atoi(mg)
					newTT := NewCache(uint32(hashSize))
					for i := 0; i < len(uci.runner.Engines); i++ {
						uci.runner.Engines[i].TranspositionTable = nil
						runtime.GC()
						uci.runner.Engines[i].TranspositionTable = newTT
					}
				} else if strings.HasPrefix(cmd, "go") {
					uci.findMove(game, depth, game.MoveClock(), cmd)
				} else if strings.HasPrefix(cmd, "position startpos moves") {
					uci.stopPondering()
					moves := strings.Fields(cmd)[3:]
					game = FromFen(startFen)
					for _, move := range game.Position().ParseMoves(moves) {
						game.Move(move)
					}
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
						game = FromFen(fen)
					} else {
						game = FromFen(fen)
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
	pondering := false
	for i := 0; i < len(fields); i++ {
		switch fields[i] {
		case "ponder":
			pondering = true
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

	for i := 0; i < len(uci.runner.Engines); i++ {
		uci.runner.Engines[i].Position = game.Position().Copy()
		uci.runner.Engines[i].Ply = ply
	}

	if timeToThink == 0 && inc == 0 {
		noTC = true
	}
	if !noTC {
		if pondering {
			tm := NewTimeManager(time.Now(), int64(timeToThink), perMove, int64(inc), int64(movesToGo), pondering)
			uci.timeManager = tm
			uci.runner.AddTimeManager(tm)
		} else {
			tm := NewTimeManager(time.Now(), int64(timeToThink), perMove, int64(inc), int64(movesToGo), pondering)
			uci.runner.AddTimeManager(tm)
		}
		go uci.runner.Search(depth)
	} else {
		tm := NewTimeManager(time.Now(), MAX_TIME, false, 0, 0, pondering)
		uci.runner.AddTimeManager(tm)
		uci.timeManager = tm
		go uci.runner.Search(depth)
	}
}

func (uci *UCI) stopPondering() {
	if uci.runner.TimeManager != nil && uci.runner.TimeManager.Pondering {
		uci.runner.TimeManager.Pondering = false
		uci.runner.TimeManager.StopSearchNow = true
		for uci.runner.TimeManager.Pondering {
		} // wait until stopped
	}
}
