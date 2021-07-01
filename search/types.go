package search

import (
	"fmt"
	"time"

	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

type Info struct {
	fpCounter                  int
	efpCounter                 int
	rfpCounter                 int
	razoringCounter            int
	checkExtentionCounter      int
	nullMoveCounter            int
	lmrCounter                 int
	lmpCounter                 int
	deltaPruningCounter        int
	seeQuiescenceCounter       int
	mainSearchCounter          int
	zwCounter                  int
	researchCounter            int
	quiesceCounter             int
	killerCounter              int
	historyCounter             int
	internalIterativeReduction int
}

func (i *Info) Print() {
	fmt.Printf("info string LMP: %d\n", i.lmpCounter)
	fmt.Printf("info string FP: %d\n", i.fpCounter)
	fmt.Printf("info string EFP: %d\n", i.efpCounter)
	fmt.Printf("info string RFP: %d\n", i.rfpCounter)
	fmt.Printf("info string Razoring: %d\n", i.razoringCounter)
	fmt.Printf("info string Check Extension: %d\n", i.checkExtentionCounter)
	fmt.Printf("info string Null-Move: %d\n", i.nullMoveCounter)
	fmt.Printf("info string LMR: %d\n", i.lmrCounter)
	fmt.Printf("info string Delta Pruning: %d\n", i.deltaPruningCounter)
	fmt.Printf("info string SEE Quiescence: %d\n", i.seeQuiescenceCounter)
	fmt.Printf("info string PV Nodes: %d\n", i.mainSearchCounter)
	fmt.Printf("info string ZW Nodes: %d\n", i.zwCounter)
	fmt.Printf("info string Research: %d\n", i.researchCounter)
	fmt.Printf("info string Quiescence Nodes: %d\n", i.quiesceCounter)
	fmt.Printf("info string Killer Moves: %d\n", i.killerCounter)
	fmt.Printf("info string History Moves: %d\n", i.historyCounter)
	fmt.Printf("info string Internal Iterative Reduction: %d\n", i.internalIterativeReduction)
}

type Engine struct {
	Position            *Position
	Ply                 uint16
	nodesVisited        int64
	cacheHits           int64
	pv                  PVLine
	StopSearchFlag      bool
	AbruptStop          bool
	move                Move
	score               int16
	positionMoves       []Move
	killerMoves         [][]Move
	searchHistory       [][]int32
	MovePickers         []*MovePicker
	StartTime           time.Time
	ThinkTime           int64
	info                Info
	pred                Predecessors
	innerLines          []PVLine
	staticEvals         []int16
	TranspositionTable  *Cache
	DebugMode           bool
	Pondering           bool
	TotalTime           float64
	nodesSinceTimeCheck int
}

var MAX_DEPTH int8 = int8(100)

func NewEngine(tt *Cache) *Engine {
	line := NewPVLine(MAX_DEPTH)
	innerLines := make([]PVLine, MAX_DEPTH)
	for i := int8(0); i < MAX_DEPTH; i++ {
		line := NewPVLine(MAX_DEPTH)
		innerLines[i] = line
	}
	movePickers := make([]*MovePicker, MAX_DEPTH)
	for i := int8(0); i < MAX_DEPTH; i++ {
		movePickers[i] = EmptyMovePicker()
	}
	return &Engine{
		nil,
		0,
		0,
		0,
		line,
		false,
		false,
		EmptyMove,
		0,
		make([]Move, MAX_DEPTH),
		make([][]Move, 125), // We assume there will be at most 126 iterations for each move/search
		make([][]int32, 12), // We have 12 pieces only
		movePickers,
		time.Now(),
		0,
		NoInfo,
		NewPredecessors(),
		innerLines,
		make([]int16, MAX_DEPTH),
		tt,
		false,
		false,
		0,
		0,
	}
}

var NoInfo = Info{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func (e *Engine) ShouldStop() bool {
	if e.nodesSinceTimeCheck < 2000 {
		e.nodesSinceTimeCheck += 1
		e.AbruptStop = e.AbruptStop || e.StopSearchFlag
		return e.AbruptStop
	}
	e.nodesSinceTimeCheck = 0
	e.AbruptStop = e.AbruptStop || e.StopSearchFlag || time.Since(e.StartTime).Milliseconds() >= e.ThinkTime
	return e.AbruptStop
}

func (e *Engine) CanFinishSearch(lastIterationTime int64) bool {
	return time.Since(e.StartTime).Milliseconds()+3*lastIterationTime < e.ThinkTime
}

func (e *Engine) ClearForSearch() {
	for i := 0; i < len(e.innerLines); i++ {
		e.innerLines[i].Recycle()
		e.staticEvals[i] = 0
	}
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
	e.AbruptStop = false
	e.nodesVisited = 0
	e.cacheHits = 0
	e.pv.Pop() // pop our move
	e.pv.Pop() // pop our opponent's move

	e.info = NoInfo

	e.pred.Clear()

	e.StartTime = time.Now()
}

func (e *Engine) NodesVisited() int64 {
	return e.nodesVisited
}

func (e *Engine) KillerMoveScore(move Move, ply int8) int32 {
	if move == EmptyMove || ply < 0 || e.killerMoves[ply] == nil {
		return 0
	}
	if e.killerMoves[ply][0] == move {
		return 100_000
	}
	if e.killerMoves[ply][1] == move {
		return 90_000
	}
	return 0
}

func (e *Engine) AddHistory(move Move, movingPiece Piece, destination Square, ply int8) {
	if ply >= 0 && !move.IsCapture() {
		e.info.killerCounter += 1
		if e.killerMoves[ply][0] != move {
			e.killerMoves[ply][1] = e.killerMoves[ply][0]
			e.killerMoves[ply][0] = move
		}

		if ply <= 1 {
			return
		}

		e.info.historyCounter += 1
		e.searchHistory[movingPiece-1][destination] += int32(ply * ply)
	}
}

func (e *Engine) RemoveMoveHistory(move Move, movingPiece Piece, destination Square, ply int8) {
	if ply >= 0 && !move.IsCapture() && e.searchHistory[movingPiece-1][destination] != 0 {
		value := e.searchHistory[movingPiece-1][destination] - int32(ply*ply)
		if value < 0 {
			value = 0
		}
		e.searchHistory[movingPiece-1][destination] = value
	}
}

func (e *Engine) AddKillerMove(move Move, ply int8) {
	if !move.IsCapture() {
		e.info.killerCounter += 1
		e.killerMoves[ply][1] = e.killerMoves[ply][0]
		e.killerMoves[ply][0] = move
	}
}

func (e *Engine) MoveHistoryScore(movingPiece Piece, destination Square, ply int8) int32 {
	if ply < 0 || e.searchHistory[movingPiece-1][destination] == 0 {
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
	if e.pv.moveCount >= 2 {
		fmt.Printf("bestmove %s ponder %s\n", mv.ToString(), e.pv.MoveAt(1).ToString())
	} else {
		fmt.Printf("bestmove %s\n", mv.ToString())
	}
}

func (e *Engine) Move() Move {
	return e.move
}

func (e *Engine) Score() int16 {
	return e.score
}

func (e *Engine) SendPv(depth int8) {
	if depth == -1 {
		depth = e.pv.moveCount
	}
	thinkTime := time.Since(e.StartTime)
	nps := int64(float64(e.nodesVisited) / thinkTime.Seconds())
	fmt.Printf("info depth %d seldepth %d tbhits %d hashfull %d nodes %d nps %d score %s time %d pv %s\n",
		depth, e.pv.moveCount, e.cacheHits, e.TranspositionTable.Consumed(),
		e.nodesVisited, nps, ScoreToCp(e.score),
		thinkTime.Milliseconds(), e.pv.ToString())
	e.TotalTime = thinkTime.Seconds()
}

func ScoreToCp(score int16) string {
	if IsCheckmateEval(score) {
		if score < 0 {
			return fmt.Sprintf("mate -%d", (CHECKMATE_EVAL+score)/2)
		} else {
			return fmt.Sprintf("mate +%d", (CHECKMATE_EVAL-score)/2)
		}
	}
	return fmt.Sprintf("cp %d", score)
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
	return Predecessors{make([]uint64, MAX_DEPTH), -1}
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
	if currentMove == EmptyMove || p.HalfMoveClock == 0 {
		return false
	}

	current := p.Hash()
	for i := pred.maxIndex - 1; i >= 0; i-- {
		var candidate = pred.line[i]
		if current == candidate {
			return true
		}
	}

	previouslySeen := p.Positions[current]

	if previouslySeen >= 2 {
		return true
	}

	return false
}

func IsPromoting(move Move) bool {
	switch move.MovingPiece() {
	case WhitePawn:
		return move.Destination().Rank() > 5
	case BlackPawn:
		return move.Destination().Rank() < 4
	default:
		return false
	}
}

func abs16(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

func max16(x int16, y int16) int16 {
	if x < y {
		return y
	}
	return x
}

func min16(x int16, y int16) int16 {
	if x >= y {
		return y
	}
	return x
}

func min8(x int8, y int8) int8 {
	if x >= y {
		return y
	}
	return x
}
