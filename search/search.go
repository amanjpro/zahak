package search

import (
	"fmt"
	"time"

	. "github.com/amanjpro/zahak/cache"
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

type Info struct {
	efpCounter            int
	rfpCounter            int
	fpCounter             int
	razoringCounter       int
	checkExtentionCounter int
	multiCutCounter       int
	nullMoveCounter       int
	lmrCounter            int
	deltaPruningCounter   int
	seeCounter            int
	seeQuiescenceCounter  int
	mainSearchCounter     int
	zwCounter             int
	quiesceCounter        int
	killerCounter         int
	historyCounter        int
}

func (i *Info) Print() {
	fmt.Println("EFP: ", i.efpCounter)
	fmt.Println("RFP: ", i.rfpCounter)
	fmt.Println("FP: ", i.fpCounter)
	fmt.Println("Razoring: ", i.razoringCounter)
	fmt.Println("Check Extension: ", i.checkExtentionCounter)
	fmt.Println("Mult-Cut: ", i.multiCutCounter)
	fmt.Println("Null-Move: ", i.nullMoveCounter)
	fmt.Println("LMR: ", i.lmrCounter)
	fmt.Println("Delta Pruning: ", i.deltaPruningCounter)
	fmt.Println("SEE: ", i.seeCounter)
	fmt.Println("SEE Quiescence: ", i.seeQuiescenceCounter)
	fmt.Println("PV Nodes: ", i.mainSearchCounter)
	fmt.Println("ZW Nodes: ", i.zwCounter)
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
	score          int32
	killerMoves    [][]Move
	searchHistory  [][]int32
	startTime      time.Time
	ThinkTime      int64
	info           Info
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
	}
}

var NoInfo = Info{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func (e *Engine) ShouldStop() bool {
	if e.StopSearchFlag {
		return true
	}
	now := time.Now()
	return now.Sub(e.startTime).Milliseconds() >= e.ThinkTime
}

var EmptyMove = Move{NoSquare, NoSquare, 0, 0}

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
	if !move.HasTag(Capture) {
		e.info.killerCounter += 1
		e.killerMoves[ply][1] = e.killerMoves[ply][0]
		e.killerMoves[ply][0] = move
	}
}

func (e *Engine) MoveHistoryScore(movingPiece Piece, destination Square, ply int8) int32 {
	if e.searchHistory[movingPiece] == nil {
		return 0
	}
	return 60_000 + e.searchHistory[movingPiece][destination]
}

func (e *Engine) AddMoveHistory(move Move, movingPiece Piece, destination Square, ply int8) {
	if !move.HasTag(Capture) {
		e.info.historyCounter += 1
		e.searchHistory[movingPiece][destination] += int32(ply)
	}
}

func (e *Engine) SendBestMove() {
	mv := e.Move()
	fmt.Printf("bestmove %s\n", mv.ToString())
}

func (e *Engine) Move() Move {
	return e.move
}

func (e *Engine) Score() int32 {
	return e.score
}

func (e *Engine) SendPv(depth int8) {
	thinkTime := time.Now().Sub(e.startTime)
	fmt.Printf("info depth %d seldepth %d nps %d tbhits %d hashfull %d nodes %d score cp %d time %d pv %s\n\n",
		depth, e.pv.moveCount, nps(e.nodesVisited, thinkTime.Seconds()),
		e.cacheHits, TranspositionTable.Consumed(), e.nodesVisited, e.score,
		thinkTime.Milliseconds(), e.pv.ToString())
}

func (e *Engine) VisitNode() {
	e.nodesVisited += 1
}

func (e *Engine) CacheHit() {
	e.cacheHits += 1
}

func (e *Engine) Search(position *Position, depth int8, ply uint16) {
	e.ClearForSearch()
	e.rootSearch(position, depth, ply)
}

func (e *Engine) rootSearch(position *Position, depth int8, ply uint16) {

	var previousBestMove Move
	alpha := -MAX_INT
	beta := MAX_INT

	e.move = EmptyMove
	e.score = alpha
	fruitelessIterations := 0

	firstScore := true
	lastDepth := int8(1)
	for iterationDepth := int8(1); iterationDepth <= depth; iterationDepth++ {
		if e.ShouldStop() {
			break
		}
		line := NewPVLine(100)
		score, ok := e.alphaBeta(position, iterationDepth, 0, alpha, beta, ply, line, EmptyMove, true, true, 0)
		if ok && (firstScore || line.moveCount >= e.pv.moveCount) {
			e.pv = line
			e.score = score
			e.move = e.pv.MoveAt(0)
			e.SendPv(iterationDepth)
			firstScore = false
		}
		lastDepth = iterationDepth
		if iterationDepth >= 20 && e.move == previousBestMove {
			fruitelessIterations++
			if fruitelessIterations > 4 {
				break
			}
		} else {
			fruitelessIterations = 0
		}
		if e.score == CHECKMATE_EVAL {
			break
		}
		previousBestMove = e.move
		e.info.Print()
	}

	e.SendPv(lastDepth)
}

func (e *Engine) alphaBeta(position *Position, depthLeft int8, searchHeight int8, alpha int32, beta int32, ply uint16, pvline *PVLine, currentMove Move, multiCutFlag bool, nullMove bool, inNullMoveSearch int8) (int32, bool) {
	e.VisitNode()

	outcome := position.Status()
	if outcome == Checkmate {
		return -CHECKMATE_EVAL, true
	} else if outcome == Draw {
		return 0, true
	}

	isRootNode := searchHeight == 0
	isPvNode := alpha != beta-1

	if depthLeft <= 0 {
		return e.quiescence(position, alpha, beta, currentMove, 0, Evaluate(position), searchHeight)
	}

	if isPvNode {
		e.info.mainSearchCounter += 1
	} else {
		e.info.zwCounter += 1
	}

	hash := position.Hash()
	cachedEval, found := TranspositionTable.Get(hash)
	if found && cachedEval.Depth >= depthLeft {
		score := cachedEval.Eval
		if score >= beta && (cachedEval.Type == UpperBound || cachedEval.Type == Exact) {
			e.CacheHit()
			return beta, true
		}
		if score <= alpha && (cachedEval.Type == LowerBound || cachedEval.Type == Exact) {
			e.CacheHit()
			return alpha, true
		}
	}

	if e.ShouldStop() {
		return -MAX_INT, false
	}

	var isInCheck bool
	if currentMove == EmptyMove {
		isInCheck = position.IsInCheck()
	} else {
		isInCheck = currentMove.HasTag(Check)
	}

	if isInCheck {
		e.info.checkExtentionCounter += 1
		depthLeft += 1 // Singular Extension
	}

	eval := Evaluate(position)
	rook := WhiteRook

	// NullMove pruning
	R := int8(3)
	isNullMoveAllowed := !isRootNode && !isPvNode && nullMove && depthLeft >= R+1 && !position.IsEndGame() && !isInCheck

	if isNullMoveAllowed {
		bound := beta
		if inNullMoveSearch == 0 {
			tempo := int32(15)   // TODO: Make it variable with a formula like: 10*(numPGAM > 0) + 10* numPGAM > 15);
			bound = beta - tempo // variable bound
		}
		ep := position.MakeNullMove()
		newBeta := 1 - bound
		line := NewPVLine(100)
		score, ok := e.alphaBeta(position, depthLeft-R, searchHeight+1, newBeta-1, newBeta, ply, line, EmptyMove, !multiCutFlag, false, inNullMoveSearch+1)
		score = -score
		position.UnMakeNullMove(ep)
		if !ok {
			return score, false
		}
		if score >= bound && abs32(score) != CHECKMATE_EVAL {
			e.info.nullMoveCounter += 1
			return beta, true // null move pruning
		}
	}

	// Razoring
	pawn := WhitePawn
	razoringMargin := 3 * pawn.Weight()
	if !isRootNode && !isPvNode && depthLeft < 2 && eval+razoringMargin < beta {
		newEval, ok := e.quiescence(position, alpha, beta, currentMove, 0, eval, searchHeight)
		if !ok {
			return newEval, ok
		}
		if newEval < beta {
			e.info.razoringCounter += 1
			return newEval, true
		}
	}

	reverseFutilityMargin := rook.Weight() // Rook + Pawn
	// Reverse Futility Pruning
	if !isRootNode && !isPvNode && depthLeft < 2 && eval-reverseFutilityMargin >= beta {
		e.info.rfpCounter += 1
		return eval - reverseFutilityMargin, true /* fail soft */
	}

	legalMoves := position.LegalMoves()
	movePicker := NewMovePicker(position, e, legalMoves, searchHeight)
	line := NewPVLine(100)

	// Internal Iterative Deepening
	if depthLeft >= 8 && !movePicker.HasPVMove() && !isInCheck {
		e.alphaBeta(position, depthLeft-7, searchHeight+1, alpha, beta, ply, line, currentMove, false, false, 0)
		if line.moveCount != 0 {
			movePicker.UpgradeToPvMove(line.MoveAt(0))
		}
	}

	// Multi-Cut Pruning
	M := 6
	C := 3
	if !isRootNode && !isPvNode && depthLeft >= R+1 && multiCutFlag && len(legalMoves) > M {
		cutNodeCounter := 0
		for i := 0; i < M; i++ {
			line.Recycle()
			move := movePicker.Next()
			capturedPiece, oldEnPassant, oldTag, hc := position.MakeMove(move)
			newBeta := 1 - beta
			// newBeta := -beta + 1
			score, ok := e.alphaBeta(position, depthLeft-R, searchHeight+1, newBeta-1, newBeta, ply, line, move, !multiCutFlag, true, inNullMoveSearch)
			score = -score
			position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece, hc)
			if !ok {
				return score, ok
			}
			if score >= beta {
				cutNodeCounter++
				if cutNodeCounter == C {
					e.info.multiCutCounter += 1
					return beta, ok // mc-prune
				}
			}
		}
	}

	// Extended Futility Pruning
	reductionsAllowed := !isRootNode && !isPvNode && !isInCheck && depthLeft > 2 && searchHeight > 3

	movePicker.Reset()

	hasSeenExact := false

	// using fail soft with negamax:
	move := movePicker.Next()
	capturedPiece, oldEnPassant, oldTag, hc := position.MakeMove(move)
	line.Recycle()
	bestscore, ok := e.alphaBeta(position, depthLeft-1, searchHeight+1, -beta, -alpha, ply, line, move, !multiCutFlag, true, inNullMoveSearch)
	bestscore = -bestscore
	position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece, hc)
	if !ok {
		return bestscore, ok
	}
	if bestscore > alpha {
		if bestscore >= beta {
			// Those scores are never useful
			if bestscore != -MAX_INT && bestscore != MAX_INT {
				TranspositionTable.Set(hash, CachedEval{hash, bestscore, depthLeft, UpperBound, ply})
			}
			e.AddKillerMove(move, searchHeight)
			return bestscore, true
		}
		alpha = bestscore
		pvline.AddFirst(move)
		pvline.ReplaceLine(line)
		hasSeenExact = true
		e.AddMoveHistory(move, position.Board.PieceAt(move.Source), move.Destination, searchHeight)
	}

	for i := 1; i < len(legalMoves); i++ {
		line.Recycle()
		move := movePicker.Next()
		if isRootNode {
			fmt.Printf("info currmove %s currmovenumber %d\n\n", move.ToString(), i+1)
		}

		LMR := int8(0)
		if reductionsAllowed && move.PromoType != NoType && !move.HasTag(Capture) && !move.HasTag(Check) {
			// Extended Futility Pruning
			gain := Evaluate(position) + rook.Weight()
			isCheckMove := move.HasTag(Check)
			if gain <= alpha && !isCheckMove && move.PromoType == NoType {
				e.info.efpCounter += 1
				continue
			}

			// Late Move Reduction
			if depthLeft > 3 && i >= 2 && !isCheckMove && move.PromoType == NoType {
				if i >= 5 {
					LMR = 2
				} else {
					LMR = 1
				}
			}
		}
		capturedPiece, oldEnPassant, oldTag, hc := position.MakeMove(move)
		score, ok := e.alphaBeta(position, depthLeft-1-LMR, searchHeight+1, -alpha-1, -alpha, ply, line, move, !multiCutFlag, true, inNullMoveSearch)
		score = -score
		if !ok {
			position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece, hc)
			return score, ok
		}
		if score > alpha && score < beta {
			line.Recycle()
			// research with window [alpha;beta]
			score, ok = e.alphaBeta(position, depthLeft-1, searchHeight+1, -beta, -alpha, ply, line, move, !multiCutFlag, true, inNullMoveSearch)
			score = -score
			if !ok {
				position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece, hc)
				return score, ok
			}
			if score > alpha {
				alpha = score
			}
		}
		position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece, hc)

		if score > bestscore { //}||
			// (score == CHECKMATE_EVAL && score >= alpha &&
			// (pvline == nil || pvline.moveCount < line.moveCount+1)) { // shorter checkmate?
			if score >= beta {
				// Those scores are never useful
				if score != -MAX_INT && score != MAX_INT {
					TranspositionTable.Set(hash, CachedEval{hash, score, depthLeft, UpperBound, ply})
				}
				e.AddKillerMove(move, searchHeight)
				return score, ok
			}

			bestscore = score
			// Potential PV move, lets copy it to the current pv-line
			pvline.AddFirst(move)
			pvline.ReplaceLine(line)
			hasSeenExact = true
			e.AddMoveHistory(move, position.Board.PieceAt(move.Source), move.Destination, searchHeight)
		}
	}
	if hasSeenExact {
		TranspositionTable.Set(hash, CachedEval{hash, bestscore, depthLeft, Exact, ply})
	} else {
		TranspositionTable.Set(hash, CachedEval{hash, bestscore, depthLeft, LowerBound, ply})
	}
	return bestscore, true
}

func nps(nodes int64, dur float64) int64 {
	if dur == 0 {
		return 0
	}
	return int64(float64(nodes) / 1000 * dur)
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
