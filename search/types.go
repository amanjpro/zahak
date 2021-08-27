package search

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

type Runner struct {
	mu           sync.RWMutex
	Engines      []*Engine
	globalInfo   Info
	nodesVisited int64
	Stop         bool
	TimeManager  *TimeManager
	DebugMode    bool
	cacheHits    int64
	pv           PVLine
	isBookmove   bool
	depth        int8
	move         Move
	score        int16
}

type Info struct {
	fpCounter                  int64
	efpCounter                 int64
	rfpCounter                 int64
	razoringCounter            int64
	checkExtentionCounter      int64
	nullMoveCounter            int64
	lmrCounter                 int64
	lmpCounter                 int64
	deltaPruningCounter        int64
	seeQuiescenceCounter       int64
	seeCounter                 int64
	mainSearchCounter          int64
	zwCounter                  int64
	researchCounter            int64
	quiesceCounter             int64
	killerCounter              int64
	historyCounter             int64
	probCutCounter             int64
	singularExtensionCounter   int64
	historyPruningCounter      int64
	multiCutCounter            int64
	internalIterativeReduction int64
}

func (e *Engine) ShareInfo() {
	atomic.AddInt64(&e.parent.globalInfo.fpCounter, e.info.fpCounter)
	atomic.AddInt64(&e.parent.globalInfo.efpCounter, e.info.efpCounter)
	atomic.AddInt64(&e.parent.globalInfo.rfpCounter, e.info.rfpCounter)
	atomic.AddInt64(&e.parent.globalInfo.razoringCounter, e.info.razoringCounter)
	atomic.AddInt64(&e.parent.globalInfo.checkExtentionCounter, e.info.checkExtentionCounter)
	atomic.AddInt64(&e.parent.globalInfo.nullMoveCounter, e.info.nullMoveCounter)
	atomic.AddInt64(&e.parent.globalInfo.lmrCounter, e.info.lmrCounter)
	atomic.AddInt64(&e.parent.globalInfo.lmpCounter, e.info.lmpCounter)
	atomic.AddInt64(&e.parent.globalInfo.deltaPruningCounter, e.info.deltaPruningCounter)
	atomic.AddInt64(&e.parent.globalInfo.seeQuiescenceCounter, e.info.seeQuiescenceCounter)
	atomic.AddInt64(&e.parent.globalInfo.seeCounter, e.info.seeCounter)
	atomic.AddInt64(&e.parent.globalInfo.mainSearchCounter, e.info.mainSearchCounter)
	atomic.AddInt64(&e.parent.globalInfo.zwCounter, e.info.zwCounter)
	atomic.AddInt64(&e.parent.globalInfo.researchCounter, e.info.researchCounter)
	atomic.AddInt64(&e.parent.globalInfo.quiesceCounter, e.info.quiesceCounter)
	atomic.AddInt64(&e.parent.globalInfo.killerCounter, e.info.killerCounter)
	atomic.AddInt64(&e.parent.globalInfo.historyCounter, e.info.historyCounter)
	atomic.AddInt64(&e.parent.globalInfo.probCutCounter, e.info.probCutCounter)
	atomic.AddInt64(&e.parent.globalInfo.historyPruningCounter, e.info.historyPruningCounter)
	atomic.AddInt64(&e.parent.globalInfo.internalIterativeReduction, e.info.internalIterativeReduction)
	atomic.AddInt64(&e.parent.globalInfo.singularExtensionCounter, e.info.singularExtensionCounter)
	atomic.AddInt64(&e.parent.globalInfo.multiCutCounter, e.info.multiCutCounter)

	atomic.AddInt64(&e.parent.nodesVisited, e.nodesVisited)
	atomic.AddInt64(&e.parent.cacheHits, e.cacheHits)
	e.info = NoInfo
	e.nodesVisited = 0
	e.cacheHits = 0
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
	fmt.Printf("info string Singular Extension: %d\n", i.singularExtensionCounter)
	fmt.Printf("info string Multi-Cut: %d\n", i.multiCutCounter)
	fmt.Printf("info string Internal Iterative Reduction: %d\n", i.internalIterativeReduction)
}

type Engine struct {
	Position           *Position
	Ply                uint16
	nodesVisited       int64
	cacheHits          int64
	positionMoves      []Move
	killerMoves        [][]Move
	searchHistory      [][]int32
	MovePickers        []*MovePicker
	triedQuietMoves    [][]Move
	info               Info
	pred               Predecessors
	score              int16
	innerLines         []PVLine
	staticEvals        []int16
	TranspositionTable *Cache
	Pawnhash           *PawnCache
	TotalTime          float64
	doPruning          bool
	isMainThread       bool
	StartTime          time.Time
	parent             *Runner
	startDepth         int8
	skipMove           Move
	skipHeight         int8
	TempMovePicker     *MovePicker
}

var MAX_DEPTH int8 = int8(100)

func (e *Engine) TimeManager() *TimeManager {
	return e.parent.TimeManager
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
			engine = NewEngine(tt, NewPawnCache(ph.Size()), t)
		}
		engines[i] = engine
	}
	t.pv = NewPVLine(MAX_DEPTH)
	t.globalInfo = NoInfo
	t.Engines = engines
	return t
}

func NewEngine(tt *Cache, ph *PawnCache, parent *Runner) *Engine {
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
		TempMovePicker:     EmptyMovePicker(),
		skipMove:           EmptyMove,
		skipHeight:         MAX_DEPTH,
	}
}

func (t *Runner) AddTimeManager(tm *TimeManager) {
	t.TimeManager = tm
}

func (r *Runner) Ponderhit() {
	r.TimeManager.StartTime = time.Now()
	r.TimeManager.Pondering = false
	fmt.Printf("info nodes %d\n", r.nodesVisited)
}

var NoInfo = Info{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func (r *Runner) ClearForSearch() {
	r.nodesVisited = 0
	r.score = -MAX_INT
	r.depth = 0
	r.isBookmove = false
	r.cacheHits = 0
	r.pv.Pop() // pop our move
	r.pv.Pop() // pop our opponent's move
	r.Stop = false
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

	e.skipMove = EmptyMove
	e.skipHeight = MAX_DEPTH

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

	e.info = NoInfo
	e.StartTime = time.Now()

	e.pred.Clear()
}

func (e *Engine) KillerMoveScore(move Move, searchHeight int8) int32 {
	if move == EmptyMove || searchHeight < 0 || e.killerMoves[searchHeight] == nil {
		return 0
	}
	if e.killerMoves[searchHeight][0] == move {
		return 90_000_000
	}
	if e.killerMoves[searchHeight][1] == move {
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
		if e.killerMoves[searchHeight][0] != move {
			e.killerMoves[searchHeight][1] = e.killerMoves[searchHeight][0]
			e.killerMoves[searchHeight][0] = move
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

func (e *Engine) AddKillerMove(move Move, searchHeight int8) {
	if !move.IsCapture() {
		e.info.killerCounter += 1
		e.killerMoves[searchHeight][1] = e.killerMoves[searchHeight][0]
		e.killerMoves[searchHeight][0] = move
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

func (e *Engine) SendPv(pv PVLine, score int16, depth int8) {
	if depth == -1 {
		depth = pv.moveCount
	}
	thinkTime := time.Since(e.StartTime)
	nodesVisited := e.parent.nodesVisited
	nps := int64(float64(nodesVisited) / thinkTime.Seconds())
	fmt.Printf("info depth %d seldepth %d hashfull %d nodes %d nps %d score %s time %d pv %s\n",
		depth, pv.moveCount, e.TranspositionTable.Consumed(),
		nodesVisited, nps, ScoreToCp(score),
		thinkTime.Milliseconds(), pv.ToString())
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
