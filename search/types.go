package search

import (
	"fmt"
	"sync/atomic"
	"time"

	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

type Runner struct {
	Engines      []*Engine
	globalInfo   Info
	nodesVisited int64
	timeManager  atomic.Value
	DebugMode    bool
}

const MAX_TIME int64 = 9_223_372_036_854_775_807

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
	seeCounter                 int
	mainSearchCounter          int
	zwCounter                  int
	researchCounter            int
	quiesceCounter             int
	killerCounter              int
	historyCounter             int
	probCutCounter             int
	historyPruningCounter      int
	internalIterativeReduction int
}

func (this *Info) add(other Info) {
	this.fpCounter = other.fpCounter
	this.efpCounter = other.efpCounter
	this.rfpCounter = other.rfpCounter
	this.razoringCounter = other.razoringCounter
	this.checkExtentionCounter = other.checkExtentionCounter
	this.nullMoveCounter = other.nullMoveCounter
	this.lmrCounter = other.lmrCounter
	this.lmpCounter = other.lmpCounter
	this.deltaPruningCounter = other.deltaPruningCounter
	this.seeQuiescenceCounter = other.seeQuiescenceCounter
	this.seeCounter = other.seeCounter
	this.mainSearchCounter = other.mainSearchCounter
	this.zwCounter = other.zwCounter
	this.researchCounter = other.researchCounter
	this.quiesceCounter = other.quiesceCounter
	this.killerCounter = other.killerCounter
	this.historyCounter = other.historyCounter
	this.probCutCounter = other.probCutCounter
	this.historyPruningCounter = other.historyCounter
	this.internalIterativeReduction = other.internalIterativeReduction
}

func (t *Runner) PrintGlobalInfo() {
	t.globalInfo = NoInfo

	for i := 0; i < len(t.Engines); i++ {
		t.globalInfo.add(t.Engines[i].info)
	}
	t.globalInfo.Print()
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
	fmt.Printf("info string ProbCut: %d\n", i.probCutCounter)
	fmt.Printf("info string Delta Pruning: %d\n", i.deltaPruningCounter)
	fmt.Printf("info string SEE Quiescence: %d\n", i.seeQuiescenceCounter)
	fmt.Printf("info string SEE: %d\n", i.seeCounter)
	fmt.Printf("info string PV Nodes: %d\n", i.mainSearchCounter)
	fmt.Printf("info string ZW Nodes: %d\n", i.zwCounter)
	fmt.Printf("info string Research: %d\n", i.researchCounter)
	fmt.Printf("info string Quiescence Nodes: %d\n", i.quiesceCounter)
	fmt.Printf("info string Killer Moves: %d\n", i.killerCounter)
	fmt.Printf("info string History Moves: %d\n", i.historyCounter)
	fmt.Printf("info string History Pruning: %d\n", i.historyPruningCounter)
	fmt.Printf("info string Internal Iterative Reduction: %d\n", i.internalIterativeReduction)
}

type Engine struct {
	Position           *Position
	Ply                uint16
	nodesVisited       int64
	cacheHits          int64
	pv                 PVLine
	move               Move
	score              int16
	positionMoves      []Move
	killerMoves        [][]Move
	searchHistory      [][]int32
	MovePickers        []*MovePicker
	triedQuietMoves    [][]Move
	info               Info
	pred               Predecessors
	innerLines         []PVLine
	staticEvals        []int16
	TranspositionTable *Cache
	Pawnhash           *PawnCache
	TotalTime          float64
	doPruning          bool
	isMainThread       bool
	StartTime          time.Time
	parent             *Runner
}

var MAX_DEPTH int8 = int8(100)

func (r *Runner) TimeManager() *TimeManager {
	tm := r.timeManager.Load()
	if tm == nil {
		return nil
	}
	return tm.(*TimeManager)
}

func (e *Engine) TimeManager() *TimeManager {
	return e.parent.TimeManager()
}

func NewRunner(tt *Cache, ph *PawnCache, numberOfThreads int) *Runner {
	t := &Runner{}
	engines := make([]*Engine, numberOfThreads)
	for i := 0; i < numberOfThreads; i++ {
		var engine *Engine
		if i == 0 {
			engine = NewEngine(tt, ph, t)
			engine.isMainThread = true
		} else {
			engine = NewEngine(tt, ph, t)
		}
		engines[i] = engine
	}
	t.globalInfo = NoInfo
	t.Engines = engines
	return t
}

func NewEngine(tt *Cache, ph *PawnCache, parent *Runner) *Engine {
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
		Position:           nil,
		Ply:                0,
		nodesVisited:       0,
		cacheHits:          0,
		pv:                 line,
		move:               EmptyMove,
		score:              0,
		positionMoves:      make([]Move, MAX_DEPTH),
		killerMoves:        make([][]Move, 125), // We assume there will be at most 126 iterations for each move/search
		searchHistory:      make([][]int32, 12), // We have 12 pieces only
		MovePickers:        movePickers,
		triedQuietMoves:    make([][]Move, 250),
		info:               NoInfo,
		pred:               NewPredecessors(),
		innerLines:         innerLines,
		staticEvals:        make([]int16, MAX_DEPTH),
		TranspositionTable: tt,
		Pawnhash:           ph,
		StartTime:          time.Now(),
		TotalTime:          0,
		doPruning:          false,
		isMainThread:       false,
		parent:             parent,
	}
}

func (t *Runner) AddTimeManager(tm *TimeManager) {
	t.timeManager.Store(tm)
}

func (t *Runner) Ponderhit() {
	tm := t.TimeManager()
	tm.StartTime = time.Now()
	tm.UpdatePondering(false)
	t.timeManager.Store(tm)
	fmt.Printf("info nodes %d\n", t.nodesVisited)
}

var NoInfo = Info{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

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

	for i := 0; i < len(e.triedQuietMoves); i++ {
		if e.triedQuietMoves[i] == nil {
			e.triedQuietMoves[i] = make([]Move, 250) // Number of potential legal moves per position
		}
		for j := 0; j < len(e.triedQuietMoves[i]); j++ {
			e.triedQuietMoves[i][j] = EmptyMove
		}
	}

	e.nodesVisited = 0
	e.cacheHits = 0
	e.pv.Pop() // pop our move
	e.pv.Pop() // pop our opponent's move

	e.info = NoInfo
	e.StartTime = time.Now()

	e.pred.Clear()
}

func (e *Engine) NodesVisited() int64 {
	return e.nodesVisited
}

func (e *Engine) KillerMoveScore(move Move, depthLeft int8) int32 {
	if move == EmptyMove || depthLeft < 0 || e.killerMoves[depthLeft] == nil {
		return 0
	}
	if e.killerMoves[depthLeft][0] == move {
		return 90_000_000
	}
	if e.killerMoves[depthLeft][1] == move {
		return 80_000_000
	}
	return 0
}

func historyBonus(current int32, bonus int32) int32 {
	return current + 32*bonus - current*abs32(bonus)/512
}

func (e *Engine) AddHistory(move Move, movingPiece Piece, destination Square, depthLeft int8, searchHeight int8, quietMovesCounter int) {
	if depthLeft >= 0 && move.PromoType() == NoType && !move.IsCapture() {
		e.info.killerCounter += 1
		if e.killerMoves[depthLeft][0] != move {
			e.killerMoves[depthLeft][1] = e.killerMoves[depthLeft][0]
			e.killerMoves[depthLeft][0] = move
		}

		if depthLeft <= 1 {
			return
		}

		e.RemoveMoveHistory(move, quietMovesCounter, depthLeft, searchHeight)
		e.info.historyCounter += 1
		e.searchHistory[movingPiece-1][destination] = historyBonus(e.searchHistory[movingPiece-1][destination], int32(depthLeft*depthLeft))
	}
}

func (e *Engine) NoteMove(move Move, quietMovesCounter int, height int8) {
	if quietMovesCounter < 0 || height < 0 || move.PromoType() != NoType || move.IsCapture() {
		return
	}
	e.triedQuietMoves[height][quietMovesCounter] = move
}

func (e *Engine) RemoveMoveHistory(killerMove Move, quietMovesCounter int, depthLeft int8, searchHeight int8) {
	if searchHeight < 0 || depthLeft < 0 {
		return
	}
	triedMoves := e.triedQuietMoves[searchHeight]
	for i := 0; i <= quietMovesCounter; i++ {
		move := triedMoves[i]
		destination := move.Destination()
		movingPiece := move.MovingPiece()
		if move != killerMove && move.PromoType() == NoType && !move.IsCapture() /* && e.searchHistory[movingPiece-1][destination] != 0 */ {
			// value := e.searchHistory[movingPiece-1][destination] - int32(depthLeft*depthLeft)
			// e.searchHistory[movingPiece-1][destination] = value
			e.searchHistory[movingPiece-1][destination] = historyBonus(e.searchHistory[movingPiece-1][destination], -int32(depthLeft*depthLeft))
		}
	}
}

func (e *Engine) AddKillerMove(move Move, depthLeft int8) {
	if !move.IsCapture() {
		e.info.killerCounter += 1
		e.killerMoves[depthLeft][1] = e.killerMoves[depthLeft][0]
		e.killerMoves[depthLeft][0] = move
	}
}

func (e *Engine) MoveHistoryScore(movingPiece Piece, destination Square, depthLeft int8) int32 {
	if depthLeft < 0 || e.searchHistory[movingPiece-1][destination] == 0 {
		return 0
	}
	return e.searchHistory[movingPiece-1][destination]
}

func (e *Engine) AddMoveHistory(move Move, movingPiece Piece, destination Square, depthLeft int8) {
	if move.PromoType() == NoType && !move.IsCapture() {
		e.info.historyCounter += 1
		e.searchHistory[movingPiece-1][destination] += int32(depthLeft)
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
	fmt.Printf("info depth %d seldepth %d hashfull %d nodes %d nps %d score %s time %d pv %s\n",
		depth, e.pv.moveCount, e.TranspositionTable.Consumed(),
		e.nodesVisited, nps, ScoreToCp(e.score),
		thinkTime.Milliseconds(), e.pv.ToString())
	e.TotalTime = thinkTime.Seconds()
}

func ScoreToCp(score int16) string {
	if isCheckmateEval(score) {
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
	for i := pred.maxIndex - 2; i >= 0; i -= 2 {
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

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
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

func min(x int, y int) int {
	if x >= y {
		return y
	}
	return x
}

func max8(x int8, y int8) int8 {
	if x <= y {
		return y
	}
	return x
}

func isCheckmateEval(eval int16) bool {
	absEval := abs16(eval)
	if absEval == MAX_INT {
		return false
	}
	return absEval >= CHECKMATE_EVAL-int16(MAX_DEPTH)
}
