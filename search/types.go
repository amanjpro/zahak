package search

import (
	"context"
	"errors"
	"fmt"
	"time"

	. "github.com/amanjpro/zahak/engine"
)

type Runner struct {
	nodesVisited int64
	cacheHits    int64
	Engines      []*Engine
	TimeManager  *TimeManager
	DebugMode    bool
	pv           PVLine
	isBookmove   bool
	depth        int8
	move         Move
	score        int16
	Ctx          context.Context
	CancelFunc   context.CancelFunc
}

type Engine struct {
	Position        *Position
	Ply             uint16
	nodesVisited    int64
	cacheHits       int64
	positionMoves   []Move
	searchHistory   MoveHistory
	MovePickers     []*MovePicker
	triedQuietMoves [][]Move
	triedNoisyMoves [][]Move
	pred            Predecessors
	seldepth        int8
	innerLines      []PVLine
	staticEvals     []int16
	TotalTime       float64
	doPruning       bool
	isMainThread    bool
	StartTime       time.Time
	parent          *Runner
	score           int16
	startDepth      int8
	skipMove        Move
	tt              *Cache
	rootMove        Move
	skipHeight      int8
	MovesToSearch   []Move
	TempMovePicker  *MovePicker
	MultiPV         int
	CurrentPV       int
	MultiPVs        []PVLine
	Scores          []int16
	NoMoves         bool
	tbHit           int64
	stop            bool
}

const MaxMultiPV = 120
const MAX_DEPTH int8 = int8(100)

var errTimeout = errors.New("Search timeout")

func (e *Engine) TimeManager() *TimeManager {
	return e.parent.TimeManager
}

func NewRunner(numberOfThreads int) *Runner {
	t := &Runner{}
	engines := make([]*Engine, numberOfThreads)
	for i := 0; i < numberOfThreads; i++ {
		var engine *Engine
		engine = NewEngine(t)
		if i == 0 {
			engine.isMainThread = true
		}
		engines[i] = engine
	}
	t.Engines = engines
	t.pv = NewPVLine(MAX_DEPTH)
	return t
}

func NewEngine(parent *Runner) *Engine {
	innerLines := make([]PVLine, MAX_DEPTH)
	for i := int8(0); i < MAX_DEPTH; i++ {
		line := NewPVLine(MAX_DEPTH)
		innerLines[i] = line
	}
	movePickers := make([]*MovePicker, MAX_DEPTH)
	for i := int8(0); i < MAX_DEPTH; i++ {
		movePickers[i] = EmptyMovePicker()
	}

	multiPVs := make([]PVLine, MaxMultiPV)
	for i := int8(0); i < MAX_DEPTH; i++ {
		line := NewPVLine(MAX_DEPTH)
		multiPVs[i] = line
	}

	return &Engine{
		Position:        nil,
		Ply:             0,
		nodesVisited:    0,
		cacheHits:       0,
		positionMoves:   make([]Move, MAX_DEPTH),
		searchHistory:   MoveHistory{},
		MovePickers:     movePickers,
		triedQuietMoves: make([][]Move, 250),
		triedNoisyMoves: make([][]Move, 250),
		pred:            NewPredecessors(),
		innerLines:      innerLines,
		staticEvals:     make([]int16, MAX_DEPTH),
		StartTime:       time.Now(),
		TotalTime:       0,
		doPruning:       false,
		isMainThread:    false,
		parent:          parent,
		score:           0,
		TempMovePicker:  EmptyMovePicker(),
		skipMove:        EmptyMove,
		rootMove:        EmptyMove,
		skipHeight:      MAX_DEPTH,
		MultiPV:         1,
		MultiPVs:        multiPVs,
		Scores:          make([]int16, MaxMultiPV),
		stop:            false,
	}
}

func (e *Engine) SetStaticEvals(height int, eval int16) {
	e.staticEvals[height] = eval
}

func (t *Runner) AddTimeManager(tm *TimeManager) {
	t.TimeManager = tm
}

func (r *Runner) Ponderhit() {
	r.TimeManager.StartTime = time.Now()
	r.TimeManager.Pondering = false
	fmt.Printf("info nodes %d\n", r.nodesVisited)
}

func (r *Runner) ResetHistory() {
	for i := 0; i < len(r.Engines); i++ {
		r.Engines[i].searchHistory = MoveHistory{}
	}
}

func (r *Runner) ClearForSearch() {
	r.nodesVisited = 0
	r.score = -MAX_INT
	r.depth = 0
	r.isBookmove = false
	r.cacheHits = 0
}

func (e *Engine) ClearForSearch() {

	for i := 0; i < len(e.MultiPVs); i++ {
		e.MultiPVs[i].Recycle()
		e.Scores[i] = 0
	}
	for i := 0; i < len(e.innerLines); i++ {
		e.innerLines[i].Recycle()
		e.staticEvals[i] = 0
	}

	e.stop = false
	e.seldepth = 0
	e.tbHit = 0
	e.score = 0
	e.NoMoves = false
	e.rootMove = EmptyMove

	// e.searchHistory.Reset()

	e.skipMove = EmptyMove
	e.skipHeight = MAX_DEPTH

	for i := 0; i < len(e.triedQuietMoves); i++ {
		if e.triedQuietMoves[i] == nil {
			e.triedQuietMoves[i] = make([]Move, 250) // Number of potential legal moves per position
		}
		if e.triedNoisyMoves[i] == nil {
			e.triedNoisyMoves[i] = make([]Move, 250) // Number of potential legal moves per position
		}
		for j := 0; j < len(e.triedQuietMoves[i]); j++ {
			e.triedQuietMoves[i][j] = EmptyMove
			e.triedNoisyMoves[i][j] = EmptyMove
		}
	}
	e.tt = TranspositionTable

	e.nodesVisited = 0
	e.cacheHits = 0

	e.StartTime = time.Now()

	e.pred.Clear()

	e.Position.Net.Recalculate(e.Position.NetInput())
}

func (e *Engine) multiPVSkipRootMove(move Move) bool {
	found := false
	for i := 0; i < e.CurrentPV; i++ {
		if e.MultiPVs[i].moveCount >= 1 && e.MultiPVs[i].line[0] == move {
			found = true
			break
		}
	}
	return found
}

func (e *Engine) NoteMove(move Move, quietMovesCounter int, noisyMovesCounter int, height int8) {
	if height < 0 {
		return
	}
	if noisyMovesCounter >= 0 && move.PromoType() != NoType || move.IsCapture() {
		e.triedNoisyMoves[height][noisyMovesCounter] = move
	} else if quietMovesCounter >= 0 {
		e.triedQuietMoves[height][quietMovesCounter] = move
	}
}

func (r *Runner) SendBestMove() {
	mv := r.Move()
	pv := r.pv
	if pv.moveCount >= 2 {
		fmt.Printf("bestmove %s ponder %s\n", mv.ToString(), pv.MoveAt(1).ToString())
	} else {
		fmt.Printf("bestmove %s\n", mv.ToString())
	}
}

func (r *Runner) Move() Move {
	return r.move
}

func (r *Runner) Score() int16 {
	return r.score
}

func (e *Engine) SendMultiPv(pv PVLine, score int16, depth int8) {
	if depth == -1 {
		depth = pv.moveCount
	}
	thinkTime := time.Since(e.StartTime)
	nodesVisited := int64(0) //e.parent.nodesVisited
	tbHits := int64(0)       //e.parent.nodesVisited
	seldepth := int8(0)
	for i := 0; i < len(e.parent.Engines); i++ {
		e := e.parent.Engines[i]
		nodesVisited += e.nodesVisited
		tbHits += e.tbHit
		if e.seldepth > seldepth {
			seldepth = e.seldepth
		}
	}
	nps := int64(float64(nodesVisited) / thinkTime.Seconds())
	fmt.Printf("info depth %d seldepth %d hashfull %d tbhits %d nodes %d nps %d score %s time %d multipv 1 pv %s\n",
		depth, seldepth, TranspositionTable.Consumed(), tbHits,
		nodesVisited, nps, ScoreToCp(score),
		thinkTime.Milliseconds(), pv.ToString())

	for i := 1; i < e.MultiPV; i++ {
		if e.MultiPVs[i].moveCount >= 1 {
			fmt.Printf("info depth %d seldepth %d hashfull %d tbhits %d nodes %d nps %d score %s time %d multipv %d pv %s\n",
				depth, seldepth, TranspositionTable.Consumed(), tbHits,
				nodesVisited, nps, ScoreToCp(e.Scores[i]),
				thinkTime.Milliseconds(), i+1, e.MultiPVs[i].ToString())
		}
	}

	e.TotalTime = thinkTime.Seconds()
}

func (e *Engine) SendPv(pv PVLine, score int16, depth int8) {
	if depth == -1 {
		depth = pv.moveCount
	}
	thinkTime := time.Since(e.StartTime)
	nodesVisited := int64(0) //e.parent.nodesVisited
	tbHits := int64(0)       //e.parent.nodesVisited
	seldepth := int8(0)
	for i := 0; i < len(e.parent.Engines); i++ {
		e := e.parent.Engines[i]
		nodesVisited += e.nodesVisited
		tbHits += e.tbHit
		if e.seldepth > seldepth {
			seldepth = e.seldepth
		}
	}
	nps := int64(float64(nodesVisited) / thinkTime.Seconds())
	fmt.Printf("info depth %d seldepth %d hashfull %d tbhits %d nodes %d nps %d score %s time %d pv %s\n",
		depth, seldepth, TranspositionTable.Consumed(), tbHits,
		nodesVisited, nps, ScoreToCp(score),
		thinkTime.Milliseconds(), pv.ToString())
	e.TotalTime = thinkTime.Seconds()
}

func ScoreToCp(score int16) string {
	if isCheckmateEval(score) {
		if score < 0 {
			return fmt.Sprintf("mate -%d", (CHECKMATE_EVAL+score)/2)
		} else {
			return fmt.Sprintf("mate +%d", (CHECKMATE_EVAL-score+1)/2)
		}
	}
	return fmt.Sprintf("cp %d", score)
}

func (e *Engine) VisitNode(searchHeight int8) {
	e.nodesVisited += 1
	if searchHeight > e.seldepth {
		e.seldepth = searchHeight
	}
	if (e.nodesVisited % 511) == 0 {
		select {
		case <-e.parent.Ctx.Done():
			panic(errTimeout)
		default:
		}
	}
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

var ttWin = CHECKMATE_EVAL - 2*int16(MAX_DEPTH)
var ttLoss = -ttWin

func evalToTT(value int16, searchHeight int8) int16 {

	if value >= ttWin {
		return value + int16(searchHeight)
	}

	if value <= ttLoss {
		return value - int16(searchHeight)
	}

	return value
}

func evalFromTT(value int16, searchHeight int8) int16 {
	if value >= ttWin {
		return value - int16(searchHeight)
	}

	if value <= ttLoss {
		return value + int16(searchHeight)
	}

	return value
}

func (e *Engine) mustSkip(move Move) bool {
	mts := e.MovesToSearch
	notInSearchMoves := len(mts) != 0
	for i := 0; i < len(mts); i++ {
		if mts[i] == move {
			notInSearchMoves = false
			break
		}
	}
	return notInSearchMoves || e.multiPVSkipRootMove(move)
}
