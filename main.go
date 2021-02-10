package main

import (
	"flag"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
)

func main() {
	var perftFlag = flag.Bool("perft", false, "Provide this to run perft tests")
	var slowFlag = flag.Bool("slow", false, "Run all perft tests, even the very slow tests")
	var perftTreeFlag = flag.Bool("perft-tree", false, "Run the engine in prefttree mode")
	var profileFlag = flag.Bool("profile", false, "Run the engine in profiling mode")
	flag.Parse()
	if *profileFlag {
		cpu, _ := os.Create("zahak-engine-cpu-profile")
		mem, _ := os.Create("zahak-engine-memory-profile")
		pprof.StartCPUProfile(cpu)
		runtime.GC()
		defer pprof.StopCPUProfile()
		defer mem.Close() // error handling omitted for example
	}
	if *perftFlag {
		StartPerftTest(*slowFlag)
	} else if *perftTreeFlag {
		depth, _ := strconv.Atoi(flag.Arg(0))
		fen := flag.Arg(1)
		game := FromFen(fen, true)
		moves := game.position.ParseMoves(flag.Args()[2:])
		PerftTree(game, depth, moves)
	} else {
		uci()
	}
}
