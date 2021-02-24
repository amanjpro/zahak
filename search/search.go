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
	move           *Move
	score          int32
	killerMoves    [][]*Move
	searchHistory  [][]int32
}

func NewEngine() *Engine {
	return &Engine{
		0,
		0,
		NewPVLine(100),
		false,
		nil,
		0,
		make([][]*Move, 600), // We assume there will be 300 moves at most
		make([][]int32, 32),  // We have 32 pieces only
	}
}

func (e *Engine) KillerMoveScore(move *Move, ply uint16) int32 {
	if e.killerMoves[ply] == nil {
		return 0
	}
	if e.killerMoves[ply][0] != nil && *e.killerMoves[ply][0] == *move {
		return 100_000
	}
	if e.killerMoves[ply][1] != nil && *e.killerMoves[ply][1] == *move {
		return 90_000
	}
	return 0
}

func (e *Engine) AddKillerMove(move *Move, ply uint16) {
	if e.killerMoves[ply] == nil {
		e.killerMoves[ply] = make([]*Move, 2)
	}
	if !move.HasTag(Capture) {
		e.killerMoves[ply][1] = e.killerMoves[ply][0]
		e.killerMoves[ply][0] = move
	}
}

func (e *Engine) MoveHistoryScore(movingPiece Piece, destination Square, ply uint16) int32 {
	if e.searchHistory[movingPiece] == nil {
		return 0
	}
	return 60_000 + e.searchHistory[movingPiece][destination]
}

func (e *Engine) AddMoveHistory(move *Move, movingPiece Piece, destination Square, ply uint16) {
	if e.searchHistory[movingPiece] == nil {
		e.searchHistory[movingPiece] = make([]int32, 64) // Number of Squares
	}
	if !move.HasTag(Capture) {
		e.searchHistory[movingPiece][destination] += int32(ply)
	}
}

func (e *Engine) SendBestMove() {
	fmt.Printf("bestmove %s\n", e.Move().ToString())
}

func (e *Engine) Move() *Move {
	return e.move
}

func (e *Engine) Score() int32 {
	return e.score
}

func (e *Engine) SendPv(thinkTime time.Duration) {
	fmt.Printf("info depth %d nps %d tbhits %d hashfull %d nodes %d score cp %d time %d pv %s\n\n",
		e.pv.moveCount, nps(e.nodesVisited, thinkTime.Seconds()),
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
	e.StopSearchFlag = false
	e.nodesVisited = 0
	e.cacheHits = 0
	e.rootSearch(position, depth, ply)
	e.pv.Pop() // pop our move
	e.pv.Pop() // pop our opponent's move
}

func (e *Engine) rootSearch(position *Position, depth int8, ply uint16) {

	// Collect evaluation for moves per iteration to help us order moves for the
	// next iteration
	legalMoves := position.LegalMoves()
	iterationEvals := make([]int32, len(legalMoves))

	var previousBestMove *Move
	alpha := -MAX_INT
	beta := MAX_INT

	e.move = nil
	e.score = alpha
	fruitelessIterations := 0

	start := time.Now()
END_LOOP:
	for iterationDepth := int8(1); iterationDepth <= depth; iterationDepth++ {
		if e.StopSearchFlag {
			break END_LOOP
		}
		currentBestScore := -MAX_INT
		orderedMoves := orderIterationMoves(&IterationMoves{legalMoves, e, iterationEvals})
		line := NewPVLine(iterationDepth + 1)
		for index, move := range orderedMoves {
			if e.StopSearchFlag {
				break END_LOOP
			}
			fmt.Printf("info currmove %s currmovenumber %d\n\n", move.ToString(), index+1)
			sendPv := false
			cp, ep, tg, hc := position.MakeMove(move)
			score := -MAX_INT
			score = -e.alphaBeta(position, iterationDepth, 1, -beta, -alpha, ply, line)
			// This only works, because checkmate eval is clearly distinguished from
			// maximum/minimum beta/alpha
			if score > beta {
				beta = score
			}
			if score == CHECKMATE_EVAL {
				alpha = score
				iterationEvals[index] = score
				if score > currentBestScore || e.pv == nil || e.pv.moveCount > line.moveCount+1 {
					currentBestScore = score
					sendPv = true
					e.pv.AddFirst(move)
					e.pv.ReplaceLine(line)
					e.move = move
					e.score = currentBestScore
				}
			} else if score > alpha && score < beta { // no very hard alpha-beta cutoff
				iterationEvals[index] = score
				alpha = score
				if score > currentBestScore {
					currentBestScore = score
					sendPv = true
					e.pv.AddFirst(move)
					e.pv.ReplaceLine(line)
					e.move = move
					e.score = currentBestScore
				}
			} else {
				iterationEvals[index] = -MAX_INT // if it is, then too bad, that is a bad move
			}
			position.UnMakeMove(move, tg, ep, cp, hc)

			timeSpent := time.Now().Sub(start)
			if sendPv {
				e.SendPv(timeSpent)
			}
		}

		if iterationDepth >= 10 && *e.move == *previousBestMove {
			fruitelessIterations++
			if fruitelessIterations > 4 {
				break
			}
		} else {
			fruitelessIterations = 0
		}
		timeSpent := time.Now().Sub(start)
		e.SendPv(timeSpent)
		if e.score == CHECKMATE_EVAL {
			break
		}
		previousBestMove = e.move
		alpha = -MAX_INT
		beta = MAX_INT
	}

	timeSpent := time.Now().Sub(start)
	e.SendPv(timeSpent)
}

func (e *Engine) alphaBeta(position *Position, depthLeft int8, searchHeight int8, alpha int32, beta int32, ply uint16, pvline *PVLine) int32 {
	e.VisitNode()
	outcome := position.Status()
	if outcome == Checkmate {
		return -CHECKMATE_EVAL
	} else if outcome == Draw {
		return 0
	}

	if e.StopSearchFlag {
		return alpha
	}

	if depthLeft == 0 {
		return e.quiescence(position, alpha, beta, 0, Evaluate(position), searchHeight)
	}

	hash := position.Hash()
	cachedEval, found := TranspositionTable.Get(hash)
	if found && cachedEval.Depth >= depthLeft {
		score := cachedEval.Eval
		if score >= beta && (cachedEval.Type == UpperBound || cachedEval.Type == Exact) {
			e.CacheHit()
			return beta
		}
		if score <= alpha && (cachedEval.Type == LowerBound || cachedEval.Type == Exact) {
			e.CacheHit()
			return alpha
		}
	}

	searchPv := true

	legalMoves := position.LegalMoves()
	movePicker := NewMovePicker(position, e, legalMoves, searchHeight, ply+uint16(searchHeight))

	for i := 0; i < len(legalMoves); i++ {
		move := movePicker.Next()
		capturedPiece, oldEnPassant, oldTag, hc := position.MakeMove(move)
		line := NewPVLine(depthLeft - 1)
		score := -MAX_INT
		if searchPv {
			score = -e.alphaBeta(position, depthLeft-1, searchHeight+1, -beta, -alpha, ply, line)
		} else {
			score = -e.zeroWindowSearch(position, depthLeft-1, searchHeight+1, -alpha, ply, true, true)
			if score > alpha { // in fail-soft ... && score < beta ) is common
				score = -e.alphaBeta(position, depthLeft-1, searchHeight+1, -beta, -alpha, ply, line) // re-search
			}
		}
		position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece, hc)
		if score >= beta {
			// Those scores are never useufl
			if score != -MAX_INT && score != MAX_INT {
				TranspositionTable.Set(hash, &CachedEval{hash, score, depthLeft, UpperBound, ply})
			}
			e.AddKillerMove(move, uint16(searchHeight)+ply)
			return beta
		}
		if score > alpha ||
			(score == CHECKMATE_EVAL && score >= alpha &&
				(pvline == nil || pvline.moveCount < line.moveCount+1)) { // shorter checkmate?
			alpha = score
			// Potential PV move, lets copy it to the current pv-line
			pvline.AddFirst(move)
			pvline.ReplaceLine(line)
			e.AddMoveHistory(move, position.Board.PieceAt(move.Source), move.Destination, uint16(searchHeight)+ply)
			searchPv = false
		}
	}
	if !searchPv {
		TranspositionTable.Set(hash, &CachedEval{hash, alpha, depthLeft, LowerBound, ply})
	} else {
		TranspositionTable.Set(hash, &CachedEval{hash, alpha, depthLeft, Exact, ply})
	}
	return alpha
}

func (e *Engine) zeroWindowSearch(position *Position, depthLeft int8, searchHeight int8, beta int32, ply uint16,
	multiCutFlag bool, nullMove bool) int32 {
	e.VisitNode()

	if e.StopSearchFlag {
		return beta - 1
	}

	if depthLeft <= 0 {
		return e.quiescence(position, beta-1, beta, 0, Evaluate(position), searchHeight)
	}

	// NullMove pruning
	isNullMoveAllowed := nullMove && depthLeft >= 5 && !position.IsEndGame() && !position.IsInCheck()
	R := int8(3)
	if searchHeight > 6 {
		R = 2
	}

	if isNullMoveAllowed {
		tempo := int32(15)    // TODO: Make it variable with a formula like: 10*(numPGAM > 0) + 10* numPGAM > 15);
		bound := beta - tempo // variable bound
		position.NullMove()
		score := -e.zeroWindowSearch(position, depthLeft-R-1, searchHeight+1, 1-bound, ply, !multiCutFlag, !nullMove)
		position.NullMove()
		if score >= bound {
			return beta // null move pruning
		}
	}

	hash := position.Hash()
	cachedEval, found := TranspositionTable.Get(hash)
	if found &&
		cachedEval.Depth >= depthLeft {
		score := cachedEval.Eval
		if score >= beta && (cachedEval.Type != UpperBound || cachedEval.Type == Exact) {
			e.CacheHit()
			return beta
		}
		if score <= beta-1 && (cachedEval.Type != LowerBound || cachedEval.Type == Exact) {
			e.CacheHit()
			return beta - 1
		}
	}
	// Multi-Cut Pruning
	M := 6
	C := 3

	legalMoves := position.LegalMoves()
	movePicker := NewMovePicker(position, e, legalMoves, searchHeight, ply+uint16(searchHeight))

	if depthLeft >= R && multiCutFlag && len(legalMoves) > M {
		cutNodeCounter := 0
		for i := 0; i < M; i++ {
			move := movePicker.Next()
			capturedPiece, oldEnPassant, oldTag, hc := position.MakeMove(move)
			score := -e.zeroWindowSearch(position, depthLeft-1-R, searchHeight+1, 1-beta, ply, !multiCutFlag, nullMove)
			position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece, hc)
			if score >= beta {
				cutNodeCounter++
				if cutNodeCounter == C {
					return beta // mc-prune
				}
			}
		}
	}

	eval := Evaluate(position)

	rook := WhiteRook
	pawn := WhitePawn
	margin := pawn.Weight() + rook.Weight() // Rook + Pawn
	futility := eval + rook.Weight()

	// Razoring
	if depthLeft < 2 && eval+margin < beta-1 {
		return e.quiescence(position, beta-1, beta, 0, eval, searchHeight)
	}

	// Reverse Futility Pruning
	if depthLeft < 5 && eval-margin >= beta {
		return eval - margin /* fail soft */
	}

	// Extended Futility Pruning
	isInCheck := position.IsInCheck()
	lastRank := Rank7
	if position.Turn() == Black {
		lastRank = Rank2
	}

	movePicker.Reset()

	for i := 0; i < len(legalMoves); i++ {
		move := movePicker.Next()
		LMR := int8(0)
		if !isInCheck && searchHeight >= 6 && depthLeft == 2 {
			board := position.Board
			movingPiece := board.PieceAt(move.Source)
			isPromoting := (movingPiece.Type() == Pawn && move.Destination.Rank() == lastRank)

			// Extended Futility Pruning
			gain := Evaluate(position) - eval
			if !move.HasTag(Check) && futility+gain <= beta-1 &&
				move.PromoType == NoType && !isPromoting {
				continue
			}

			// Late Move Reduction
			if i >= 5 && !move.HasTag(Check) && move.PromoType == NoType && !isPromoting {
				LMR = 1
			}
		}

		capturedPiece, oldEnPassant, oldTag, hc := position.MakeMove(move)
		score := -e.zeroWindowSearch(position, depthLeft-1-LMR, searchHeight+1, 1-beta, ply, !multiCutFlag, nullMove)

		position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece, hc)
		if score >= beta {
			// Those scores are never useufl
			if score != -MAX_INT && score != MAX_INT {
				TranspositionTable.Set(hash, &CachedEval{hash, score, depthLeft, UpperBound, ply})
			}
			return beta // fail-hard beta-cutoff
		}
	}
	TranspositionTable.Set(hash, &CachedEval{hash, beta - 1, depthLeft, LowerBound, ply})
	return beta - 1 // fail-hard, return alpha
}

func nps(nodes int64, dur float64) int64 {
	if dur == 0 {
		return 0
	}
	return int64(float64(nodes) / 1000 * dur)
}
