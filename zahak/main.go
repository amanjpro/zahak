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
	. "github.com/amanjpro/zahak/search"
	. "github.com/amanjpro/zahak/strength"
	. "github.com/amanjpro/zahak/tuning"
	. "github.com/amanjpro/zahak/uci"
)

var version = "dev"

func main() {
	args := os.Args
	if len(args) > 1 && args[1] == "bench" {
		RunBenchmark()
	} else {
		var genEpdFlag = flag.Bool("gen-epds", false, "Generate open EPDs for self-play")
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
		if *genEpdFlag {
			GenerateEpds()
		} else if *prepareTuningFlag && *epdPath != "" {
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
			game := FromFen(fen)
			moves := []Move{}
			if len(flag.Args()) > 2 {
				game.Position().ParseMoves(strings.Fields(flag.Args()[2]))
			}
			PerftTree(game, depth, moves)
		} else {
			fmt.Print(`

ZZZZZZZZZZZZZZZZZZZ                 hhhhhhh                               kkkkkkkk
Z:::::::::::::::::Z                 h:::::h                               k::::::k
Z:::::::::::::::::Z                 h:::::h                               k::::::k
Z:::ZZZZZZZZ:::::Z                  h:::::h                               k::::::k
ZZZZZ     Z:::::Z    aaaaaaaaaaaaa   h::::h hhhhh         aaaaaaaaaaaaa    k:::::k    kkkkkkk
        Z:::::Z      a::::::::::::a  h::::hh:::::hhh      a::::::::::::a   k:::::k   k:::::k
       Z:::::Z       aaaaaaaaa:::::a h::::::::::::::hh    aaaaaaaaa:::::a  k:::::k  k:::::k
      Z:::::Z                 a::::a h:::::::hhh::::::h            a::::a  k:::::k k:::::k
     Z:::::Z           aaaaaaa:::::a h::::::h   h::::::h    aaaaaaa:::::a  k::::::k:::::k
    Z:::::Z          aa::::::::::::a h:::::h     h:::::h  aa::::::::::::a  k:::::::::::k
   Z:::::Z          a::::aaaa::::::a h:::::h     h:::::h a::::aaaa::::::a  k:::::::::::k
ZZZ:::::Z     ZZZZZa::::a    a:::::a h:::::h     h:::::ha::::a    a:::::a  k::::::k:::::k
Z::::::ZZZZZZZZ:::Za::::a    a:::::a h:::::h     h:::::ha::::a    a:::::a k::::::k k:::::k
Z:::::::::::::::::Za:::::aaaa::::::a h:::::h     h:::::ha:::::aaaa::::::a k::::::k  k:::::k
Z:::::::::::::::::Z a::::::::::aa:::ah:::::h     h:::::h a::::::::::aa:::ak::::::k   k:::::k
ZZZZZZZZZZZZZZZZZZZ  aaaaaaaaaa  aaaahhhhhhh     hhhhhhh  aaaaaaaaaa  aaaakkkkkkkk    kkkkkkk


`)
			NewUCI(version, *bookPath != "", *bookPath).Start()
		}
	}
}
