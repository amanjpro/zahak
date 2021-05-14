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
	. "github.com/amanjpro/zahak/tuning"
	. "github.com/amanjpro/zahak/uci"
)

var version = "dev"

func main() {
	args := os.Args
	if len(args) > 1 && args[1] == "bench" {
		NewUCI(version, false, "", true).Start()
	} else {
		var perftFlag = flag.Bool("perft", false, "Provide this to run perft tests")
		var slowFlag = flag.Bool("slow", false, "Run all perft tests, even the very slow tests")
		var tuneFlag = flag.Bool("tune", false, "Peform texel tuning for optimal evaluation values")
		var prepareTuningFlag = flag.Bool("prepare-tuning-data", false, "Prepare quiet EPDs for tuning")
		var perftTreeFlag = flag.Bool("perft-tree", false, "Run the engine in prefttree mode")
		var profileFlag = flag.Bool("profile", false, "Run the engine in profiling mode")
		var bookPath = flag.String("book", "", "Path to openning book in PolyGlot (bin) format")
		var epdPath = flag.String("test-positions", "", "Path to EPD positions, used to test the strength of the engine")
		var excludeParams = flag.String("exclude-params", "", "Exclude parameters when tuning, format: 1, 9, 10, 11 or 1, 9-11")
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
		if *prepareTuningFlag && *epdPath != "" {
			PrepareTuningData(*epdPath)
		} else if *tuneFlag && *epdPath != "" {
			paramsToExclude := make(map[int]bool)
			if excludeParams != nil {
				fields := strings.Split(*excludeParams, ",")
				for _, f := range fields {
					fields := strings.Split(f, "-")
					if len(fields) == 1 {
						f = strings.Trim(f, " ")
						i, _ := strconv.Atoi(f)
						paramsToExclude[i] = true
					} else if len(fields) == 2 {
						fst := strings.Trim(fields[0], " ")
						lst := strings.Trim(fields[1], " ")
						i, _ := strconv.Atoi(fst)
						j, _ := strconv.Atoi(lst)
						for ; i <= j; i++ {
							paramsToExclude[i] = true
						}
					}
				}
			}
			Tune(*epdPath, paramsToExclude)
		} else if *epdPath != "" {
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
			NewUCI(version, *bookPath != "", *bookPath, false).Start()
		}
	}
}
