package search

import (
	"fmt"
	"time"

	. "github.com/amanjpro/zahak/cache"
	. "github.com/amanjpro/zahak/engine"
)

type Info struct {
	efpCounter            int
	rfpCounter            int
	razoringCounter       int
	checkExtentionCounter int
	multiCutCounter       int
	nullMoveCounter       int
	lmrCounter            int
	deltaPruningCounter   int
	seeQuiescenceCounter  int
	mainSearchCounter     int
	zwCounter             int
	researchCounter       int
	quiesceCounter        int
	killerCounter         int
	historyCounter        int
}

func (i *Info) Print() {
	fmt.Println("EFP: ", i.efpCounter)
	fmt.Println("RFP: ", i.rfpCounter)
	fmt.Println("Razoring: ", i.razoringCounter)
	fmt.Println("Check Extension: ", i.checkExtentionCounter)
	fmt.Println("Mult-Cut: ", i.multiCutCounter)
	fmt.Println("Null-Move: ", i.nullMoveCounter)
	fmt.Println("LMR: ", i.lmrCounter)
	fmt.Println("Delta Pruning: ", i.deltaPruningCounter)
	fmt.Println("SEE Quiescence: ", i.seeQuiescenceCounter)
	fmt.Println("PV Nodes: ", i.mainSearchCounter)
	fmt.Println("ZW Nodes: ", i.zwCounter)
	fmt.Println("Research: ", i.researchCounter)
	fmt.Println("Quiescence Nodes: ", i.quiesceCounter)
	fmt.Println("Killer Moves: ", i.killerCounter)
	fmt.Println("History Moves: ", i.historyCounter)
}

type Engine struct {
	nodesVisited   int64
	cacheHits      int64
	pv             *PVLine
	StopSearchFlag bool
	move           Move
	score          int16
	killerMoves    [][]Move
	searchHistory  [][]int32
	startTime      time.Time
	ThinkTime      int64
	info           Info
	pred           Predecessors
}

func NewEngine() *Engine {
	return &Engine{
		0,
		0,
		NewPVLine(100),
		false,
		EmptyMove,
		0,
		make([][]Move, 125), // We assume there will be at most 126 iterations for each move/search
		make([][]int32, 12), // We have 12 pieces only
		time.Now(),
		0,
		NoInfo,
		NewPredecessors(),
	}
}

var NoInfo = Info{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func (e *Engine) ShouldStop() bool {
	if e.StopSearchFlag {
		return true
	}
	now := time.Now()
	return now.Sub(e.startTime).Milliseconds() >= e.ThinkTime
}

func (e *Engine) ClearForSearch() {
	for i := 0; i < len(e.killerMoves); i++ {
		if e.killerMoves[i] == nil {
			e.killerMoves[i] = make([]Move, 2)
		}
		for j := 0; j < len(e.killerMoves[i]); j++ {
			e.killerMoves[i][j] = EmptyMove
		}
	}

	for i := 0; i < len(e.searchHistory); i++ {
		if e.searchHistory[i] == nil {
			e.searchHistory[i] = make([]int32, 64) // Number of Squares
		}
		for j := 0; j < len(e.searchHistory[i]); j++ {
			e.searchHistory[i][j] = 0
		}
	}

	e.StopSearchFlag = false
	e.nodesVisited = 0
	e.cacheHits = 0
	e.pv.Pop() // pop our move
	e.pv.Pop() // pop our opponent's move

	e.info = NoInfo

	e.pred.Clear()

	e.startTime = time.Now()
}

func (e *Engine) KillerMoveScore(move Move, ply int8) int32 {
	if e.killerMoves[ply] == nil {
		return 0
	}
	if e.killerMoves[ply][0] != EmptyMove && e.killerMoves[ply][0] == move {
		return 100_000
	}
	if e.killerMoves[ply][1] != EmptyMove && e.killerMoves[ply][1] == move {
		return 90_000
	}
	return 0
}

func (e *Engine) AddKillerMove(move Move, ply int8) {
	if !move.IsCapture() {
		e.info.killerCounter += 1
		e.killerMoves[ply][1] = e.killerMoves[ply][0]
		e.killerMoves[ply][0] = move
	}
}

func (e *Engine) MoveHistoryScore(movingPiece Piece, destination Square, ply int8) int32 {
	if e.searchHistory[movingPiece-1] == nil {
		return 0
	}
	return 60_000 + e.searchHistory[movingPiece-1][destination]
}

func (e *Engine) AddMoveHistory(move Move, movingPiece Piece, destination Square, ply int8) {
	if !move.IsCapture() {
		e.info.historyCounter += 1
		e.searchHistory[movingPiece-1][destination] += int32(ply)
	}
}

func (e *Engine) SendBestMove() {
	mv := e.Move()
	fmt.Printf("bestmove %s\n", mv.ToString())
}

func (e *Engine) Move() Move {
	return e.move
}

func (e *Engine) Score() int16 {
	return e.score
}

func (e *Engine) SendPv(depth int8) {
	thinkTime := time.Now().Sub(e.startTime)
	fmt.Printf("info depth %d seldepth %d tbhits %d hashfull %d nodes %d score cp %d time %d pv %s\n\n",
		depth, e.pv.moveCount, e.cacheHits, TranspositionTable.Consumed(), e.nodesVisited, e.score,
		thinkTime.Milliseconds(), e.pv.ToString())
}

func (e *Engine) VisitNode() {
	e.nodesVisited += 1
}

func (e *Engine) CacheHit() {
	e.cacheHits += 1
}

type Predecessors struct {
	line     []uint64
	maxIndex int
}

func NewPredecessors() Predecessors {
	return Predecessors{make([]uint64, 100), -1}
}

func (p *Predecessors) Push(hash uint64) {
	p.maxIndex += 1
	p.line[p.maxIndex] = hash
}

func (p *Predecessors) Clear() {
	p.maxIndex = -1
}

func (p *Predecessors) Pop() {
	if p.maxIndex < 0 {
		return
	}
	p.maxIndex -= 1
}

func IsRepetition(p *Position, pred Predecessors, currentMove Move) bool {
	current := p.Hash()
	previouslySeen := p.Positions[current]

	if currentMove == EmptyMove || p.HalfMoveClock == 0 {
		return false
	}

	if previouslySeen >= 3 {
		return true
	}

	for i := pred.maxIndex - 1; i >= 0; i-- {
		var candidate = pred.line[i]
		if current == candidate {
			if previouslySeen > 0 {
				return true
			} else {
				previouslySeen += 1
			}
		}
	}

	if previouslySeen >= 2 {
		return true
	}
	return false
}

func abs16(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}
