package strength

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/search"
)

type EPDEntry struct {
	fen       string
	bestmoves []string
	badmoves  []string
	id        string
}

func readEPDs(path string) []EPDEntry {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	epds := make([]EPDEntry, 0, 100)

	scanner := bufio.NewScanner(file)
	var id string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		fields := strings.Fields(line)
		fen := fmt.Sprint(strings.Join(fields[:4], " "), " 0 1")
		bestmoves := make([]string, 0, 10)
		badmoves := make([]string, 0, 10)
		for index := 4; index < len(fields); index += 1 {
			field := fields[index]
			if field == "am" {
				for index += 1; index < len(fields); index += 1 {
					mv := fields[index]
					badmoves = append(badmoves, strings.Trim(mv, ";"))
					if strings.HasSuffix(mv, ";") {
						break
					}
				}
			} else if field == "bm" {
				for index += 1; index < len(fields); index += 1 {
					mv := fields[index]
					bestmoves = append(bestmoves, strings.Trim(mv, ";"))
					if strings.HasSuffix(mv, ";") {
						break
					}
				}
			} else if field == "id" {
				id = strings.Trim(fields[index+1], "\";")
				index += 1
			}
		}
		epds = append(epds, EPDEntry{fen, bestmoves, badmoves, id})
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return epds
}

func RunTestPositions(path string) {
	epds := readEPDs(path)
	success := 0
	for _, epd := range epds {
		old := os.Stdout // keep backup of the real stdout
		_, w, _ := os.Pipe()
		os.Stdout = w

		game := FromFen(epd.fen)
		TranspositionTable = NewCache(DEFAULT_CACHE_SIZE)
		r := NewRunner(1)
		r.AddTimeManager(NewTimeManager(time.Now(), 15000, true, 0, 0, false))
		pos := game.Position()
		r.Engines[0].Position = pos
		r.Search(MAX_DEPTH, -2, -1)
		mv := pos.MoveToPGN(r.Move())

		// back to normal state
		w.Close()
		os.Stdout = old
		if contains(epd.badmoves, mv) {
			fmt.Printf("EPD-id: %s found a very bad move, it found %s, best moves were %s\n", epd.id, mv, epd.bestmoves)
			success -= 1
		} else if !contains(epd.bestmoves, mv) && len(epd.bestmoves) != 0 {
			fmt.Printf("EPD-id: %s failed to find the best move, but avoided the worst move, it found %s, best moves were %s\n", epd.id, mv, epd.bestmoves)
		} else {
			success += 1
		}
	}

	fmt.Printf("Score was %d out of %d\n", success, len(epds))
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
