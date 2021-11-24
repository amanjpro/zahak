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

const MaxSkillLevels = 7

type UCI struct {
	version     string
	runner      *Runner
	timeManager *TimeManager
	withBook    bool
	bookPath    string
	multiPV     int
}

func NewUCI(version string, withBook bool, bookPath string) *UCI {
	return &UCI{
		version,
		NewRunner(1),
		nil,
		withBook,
		bookPath,
		1,
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
				fmt.Printf("option name OwnBook type check default %t\n", uci.withBook)
				fmt.Printf("option name Threads type spin default %d min %d max %d\n", defaultCPU, minCPU, maxCPU)
				fmt.Print("option name EvalFile type string default <empty>\n")
				fmt.Print("option name BookFile type string default <empty>\n")
				fmt.Print("option name SyzygyPath type string default <empty>\n")
				fmt.Printf("option name SyzygyProbeDepth type spin default %d min 0 max 128\n", DefaultProbeDepth)
				fmt.Printf("option name MultiPV type spin default 1 min 1 max %d\n", MaxMultiPV)
				fmt.Printf("option name Skill Level type spin default %d min 1 max %d\n", MaxSkillLevels, MaxSkillLevels)
				fmt.Print("uciok\n")
			case "isready":
				fmt.Print("readyok\n")
			case "isdraw":
				fmt.Print(game.Position().IsDraw(), "\n")
			case "draw":
				fmt.Print(game.Position().Board.Draw(), "\n")
			case "ucinewgame", "position startpos":
				game = FromFen(startFen)
				uci.runner.ResetHistory()
			case "fen":
				fmt.Println(game.Position().Fen())
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
					ClearSyzygy()
					SetSyzygyPath(path)
				} else if strings.HasPrefix(cmd, "setoption name EvalFile value") {
					path := strings.TrimSpace(strings.ReplaceAll(cmd, "setoption name EvalFile value", ""))
					if path == "" || path == "<empty>" {
						fmt.Print("info string no eval file is selected, ignoring\n")
						continue
					}
					err := LoadNetwork(path)
					if err != nil {
						fmt.Printf("info string an error happened when loading the network %s\n", err)
					} else {
						fmt.Printf("info string new EvalFile loaded, the id of the new EvalFile is %d\n", CurrentNetworkId)
					}
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
				} else if strings.HasPrefix(cmd, "setoption name MultiPV value ") {
					options := strings.Fields(cmd)
					v := options[len(options)-1]
					multiPV, _ := strconv.Atoi(v)
					uci.runner.Engines[0].MultiPV = multiPV
					uci.multiPV = multiPV
				} else if strings.HasPrefix(cmd, "setoption name OwnBook value ") {
					options := strings.Fields(cmd)
					opt := options[len(options)-1]
					if opt == "false" {
						ResetBook()
					} else if !IsBoookLoaded() && opt == "true" { // if it is loaded, no need to reload
						uci.withBook = true
						if uci.bookPath != "" && uci.bookPath != "<empty>" {
							InitBook(uci.bookPath)
						}
					}
				} else if strings.HasPrefix(cmd, "setoption name Skill Level value") {
					options := strings.Fields(cmd)
					v := options[len(options)-1]
					level, _ := strconv.Atoi(v)
					switchToLevel(level)
				} else if strings.HasPrefix(cmd, "setoption name Threads value") {
					options := strings.Fields(cmd)
					v := options[len(options)-1]
					cpu, _ := strconv.Atoi(v)
					uci.runner = NewRunner(cpu)
					uci.runner.Engines[0].MultiPV = uci.multiPV
				} else if strings.HasPrefix(cmd, "setoption name Hash value") {
					options := strings.Fields(cmd)
					mg := options[len(options)-1]
					hashSize, _ := strconv.Atoi(mg)
					TranspositionTable = nil
					runtime.GC()
					TranspositionTable = NewCache(hashSize)
				} else if strings.HasPrefix(cmd, "go") {
					uci.findMove(game, depth, game.MoveClock(), cmd)
				} else if strings.HasPrefix(cmd, "position startpos moves") {
					uci.stopPondering()
					moves := strings.Fields(cmd)[3:]
					game = FromFen(startFen)
					for _, move := range game.Position().ParseGameMoves(moves) {
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
					for _, move := range game.Position().ParseGameMoves(moves) {
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
	nodes := -1
	mateIn := -1

	pos := game.Position()
	noTC := false
	timeToThink := 0
	inc := 0
	movesToGo := 0
	perMove := false
	pondering := false
	var movesToSearch []Move
	for i := 0; i < len(fields); i++ {
		switch fields[i] {
		case "searchmoves":
			movesToSearch = pos.ParseSearchMoves(fields[i+1:])
			i += len(movesToSearch)
		case "ponder":
			pondering = true
		case "wtime":
			if pos.Turn() == White {
				timeToThink, _ = strconv.Atoi(fields[i+1])
				i++
			}
		case "nodes":
			nodes, _ = strconv.Atoi(fields[i+1])
			i++
		case "mate":
			mateIn, _ = strconv.Atoi(fields[i+1])
			i++
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
		uci.runner.Engines[i].MovesToSearch = movesToSearch
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
		go uci.runner.Search(depth, 2*int16(mateIn), int64(nodes))
	} else {
		tm := NewTimeManager(time.Now(), MAX_TIME, false, 0, 0, pondering)
		uci.runner.AddTimeManager(tm)
		uci.timeManager = tm
		go uci.runner.Search(depth, 2*int16(mateIn), int64(nodes))
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

func switchToLevel(level int) {
	switch level {
	case 1:
		Skills1Init()
		NetHiddenSize = Skills1NetHiddenSize
		CurrentHiddenWeights = Skills1HiddenWeights
		CurrentHiddenBiases = Skills1HiddenBiases
		CurrentOutputWeights = Skills1OutputWeights
		CurrentOutputBias = Skills1OutputBias
		CurrentNetworkId = Skills1NetworkId
	case 2:
		Skills2Init()
		NetHiddenSize = Skills2NetHiddenSize
		CurrentHiddenWeights = Skills2HiddenWeights
		CurrentHiddenBiases = Skills2HiddenBiases
		CurrentOutputWeights = Skills2OutputWeights
		CurrentOutputBias = Skills2OutputBias
		CurrentNetworkId = Skills2NetworkId
	case 3:
		Skills3Init()
		NetHiddenSize = Skills3NetHiddenSize
		CurrentHiddenWeights = Skills3HiddenWeights
		CurrentHiddenBiases = Skills3HiddenBiases
		CurrentOutputWeights = Skills3OutputWeights
		CurrentOutputBias = Skills3OutputBias
		CurrentNetworkId = Skills3NetworkId
	case 4:
		Skills4Init()
		NetHiddenSize = Skills4NetHiddenSize
		CurrentHiddenWeights = Skills4HiddenWeights
		CurrentHiddenBiases = Skills4HiddenBiases
		CurrentOutputWeights = Skills4OutputWeights
		CurrentOutputBias = Skills4OutputBias
		CurrentNetworkId = Skills4NetworkId
	case 5:
		Skills5Init()
		NetHiddenSize = Skills5NetHiddenSize
		CurrentHiddenWeights = Skills5HiddenWeights
		CurrentHiddenBiases = Skills5HiddenBiases
		CurrentOutputWeights = Skills5OutputWeights
		CurrentOutputBias = Skills5OutputBias
		CurrentNetworkId = Skills5NetworkId
	case 6:
		Skills6Init()
		NetHiddenSize = Skills6NetHiddenSize
		CurrentHiddenWeights = Skills6HiddenWeights
		CurrentHiddenBiases = Skills6HiddenBiases
		CurrentOutputWeights = Skills6OutputWeights
		CurrentOutputBias = Skills6OutputBias
		CurrentNetworkId = Skills6NetworkId
	case MaxSkillLevels:
		NetHiddenSize = DefaultNetHiddenSize
		CurrentHiddenWeights = DefaultHiddenWeights
		CurrentHiddenBiases = DefaultHiddenBiases
		CurrentOutputWeights = DefaultOutputWeights
		CurrentOutputBias = DefaultOutputBias
		CurrentNetworkId = DefaultNetworkId
	default:
		fmt.Printf("info string unsupported skills level, only values between 1 and %d is supported\n", MaxSkillLevels)
	}
}
