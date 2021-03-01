package search

import (
	"fmt"
	"time"

	. "github.com/amanjpro/zahak/cache"
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

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
	}
}

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
		line := NewPVLine(iterationDepth + 1)
		score, ok := e.alphaBeta(position, iterationDepth, 0, alpha, beta, ply, line, true, true, 0)
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
	}

	e.SendPv(lastDepth)
}

func (e *Engine) alphaBeta(position *Position, depthLeft int8, searchHeight int8, alpha int32, beta int32, ply uint16, pvline *PVLine,
	multiCutFlag bool, nullMove bool, inNullMoveSearch int8) (int32, bool) {
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
		return e.quiescence(position, alpha, beta, 0, Evaluate(position), searchHeight), true
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

	legalMoves := position.LegalMoves()

	movePicker := NewMovePicker(position, e, legalMoves, searchHeight)

	if e.ShouldStop() {
		return -MAX_INT, false
	}

	isInCheck := false
	if !isRootNode && !isPvNode { // only compute it if it is not pv-node
		isInCheck = position.IsInCheck()
	}

	if isInCheck {
		depthLeft += 1 // Singular Extension
	}

	eval := Evaluate(position)

	// NullMove pruning
	R := int8(3)
	if searchHeight > 6 {
		R = 2
	}
	isNullMoveAllowed := !isRootNode && !isPvNode && nullMove && depthLeft >= R+2 && !position.IsEndGame() && !isInCheck && eval >= beta

	if isNullMoveAllowed {
		bound := beta
		if inNullMoveSearch == 0 {
			tempo := int32(15)   // TODO: Make it variable with a formula like: 10*(numPGAM > 0) + 10* numPGAM > 15);
			bound = beta - tempo // variable bound
		}
		ep := position.MakeNullMove()
		newBeta := 1 - bound
		line := NewPVLine(depthLeft - 1 - R)
		score, ok := e.alphaBeta(position, depthLeft-R-1, searchHeight+1, newBeta-1, newBeta, ply, line, !multiCutFlag, false, inNullMoveSearch+1)
		score = -score
		position.UnMakeNullMove(ep)
		if !ok {
			return score, false
		}
		if score >= bound {
			return beta, true // null move pruning
		}
	}

	// Multi-Cut Pruning
	M := 6
	C := 3
	if !isRootNode && !isPvNode && depthLeft >= R+2 && multiCutFlag && len(legalMoves) > M {
		cutNodeCounter := 0
		for i := 0; i < M; i++ {
			move := movePicker.Next()
			capturedPiece, oldEnPassant, oldTag, hc := position.MakeMove(move)
			line := NewPVLine(depthLeft - 1 - R)
			newBeta := 1 - beta
			score, ok := e.alphaBeta(position, depthLeft-1-R, searchHeight+1, newBeta-1, newBeta, ply, line, !multiCutFlag, true, inNullMoveSearch)
			score = -score
			position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece, hc)
			if !ok {
				return score, ok
			}
			if score >= beta {
				cutNodeCounter++
				if cutNodeCounter == C {
					return beta, ok // mc-prune
				}
			}
		}
	}

	rook := WhiteRook
	pawn := WhitePawn
	margin := pawn.Weight() + rook.Weight() // Rook + Pawn
	futility := eval + rook.Weight()

	// Razoring
	if !isRootNode && !isPvNode && depthLeft < 2 && eval+margin < beta-1 {
		return e.quiescence(position, alpha, beta, 0, eval, searchHeight), true
	}

	// Reverse Futility Pruning
	if !isRootNode && !isPvNode && depthLeft < 5 && eval-margin >= beta {
		return eval - margin, true /* fail soft */
	}

	// Extended Futility Pruning
	reductionsAllowed := !isRootNode && !isPvNode && !isInCheck

	movePicker.Reset()

	hasSeenExact := false

	// using fail soft with negamax:
	bestscore := -MAX_INT
	move := movePicker.Next()
	capturedPiece, oldEnPassant, oldTag, hc := position.MakeMove(move)
	line := NewPVLine(depthLeft - 1)
	score, ok := e.alphaBeta(position, depthLeft-1, searchHeight+1, -beta, -alpha, ply, line, !multiCutFlag, true, inNullMoveSearch)
	bestscore = -score
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
		if reductionsAllowed && searchHeight >= 6 && depthLeft == 2 {

			// Extended Futility Pruning
			gain := Evaluate(position) + futility
			isCheckMove := move.HasTag(Check)
			if gain <= alpha && !isCheckMove && move.PromoType == NoType {
				continue
			}

			// Late Move Reduction
			if i >= 5 && !isCheckMove && move.PromoType == NoType {
				LMR = 1
			}
		}
		capturedPiece, oldEnPassant, oldTag, hc := position.MakeMove(move)
		score, ok := e.alphaBeta(position, depthLeft-1-LMR, searchHeight+1, -alpha-1, -alpha, ply, line, !multiCutFlag, true, inNullMoveSearch)
		score = -score
		if !ok {
			position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece, hc)
			return score, ok
		}
		if score > alpha && score < beta {
			line.Recycle()
			// research with window [alpha;beta]
			score, ok = e.alphaBeta(position, depthLeft-1-LMR, searchHeight+1, -beta, -alpha, ply, line, !multiCutFlag, true, inNullMoveSearch)
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
