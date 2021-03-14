package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"

	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/perft"
	. "github.com/amanjpro/zahak/strength"
	. "github.com/amanjpro/zahak/uci"
)

func main() {
	var perftFlag = flag.Bool("perft", false, "Provide this to run perft tests")
	var slowFlag = flag.Bool("slow", false, "Run all perft tests, even the very slow tests")
	var perftTreeFlag = flag.Bool("perft-tree", false, "Run the engine in prefttree mode")
	var profileFlag = flag.Bool("profile", false, "Run the engine in profiling mode")
	var bookPath = flag.String("book", "", "Path to openning book in PolyGlot (bin) format")
	var epdPath = flag.String("test-positions", "", "Path to EPD positions, used to test the strength of the engine")
	flag.Parse()
	if *profileFlag {
		cpu, err := os.Create("zahak-engine-cpu-profile")
		if err != nil {
			fmt.Println("could not create CPU profile: ", err)
			os.Exit(1)
		}
		if err := pprof.StartCPUProfile(cpu); err != nil {
			fmt.Println("could not start CPU profile: ", err)
			os.Exit(1)
		}

		mem, _ := os.Create("zahak-engine-memory-profile")
		runtime.GC()
		pprof.WriteHeapProfile(mem)
		defer cpu.Close()
		defer pprof.StopCPUProfile()
		defer mem.Close() // error handling omitted for example
	}
	if *epdPath != "" {
		RunTestPositions(*epdPath)
	} else if *perftFlag {
		StartPerftTest(*slowFlag)
	} else if *perftTreeFlag {
		depth, _ := strconv.Atoi(flag.Arg(0))
		fen := flag.Arg(1)
		game := FromFen(fen, true)
		moves := []Move{}
		if len(flag.Args()) > 2 {
			game.Position().ParseMoves(strings.Fields(flag.Args()[2]))
		}
		PerftTree(game, depth, moves)
	} else {
		NewUCI(*bookPath != "", *bookPath).Start()
	}
}
