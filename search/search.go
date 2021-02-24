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

	var previousBestMove *Move
	alpha := -MAX_INT
	beta := MAX_INT

	e.move = nil
	e.score = alpha
	fruitelessIterations := 0

	start := time.Now()
	for iterationDepth := int8(1); iterationDepth <= depth; iterationDepth++ {
		if e.StopSearchFlag {
			break
		}
		line := NewPVLine(iterationDepth + 1)
		e.score = e.alphaBeta(position, iterationDepth, 0, alpha, beta, ply, line, true, true)
		e.pv = line
		e.move = e.pv.MoveAt(0)
		timeSpent := time.Now().Sub(start)
		e.SendPv(timeSpent)
		if iterationDepth >= 10 && *e.move == *previousBestMove {
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

	timeSpent := time.Now().Sub(start)
	e.SendPv(timeSpent)
}

func (e *Engine) alphaBeta(position *Position, depthLeft int8, searchHeight int8, alpha int32, beta int32, ply uint16, pvline *PVLine,
	multiCutFlag bool, nullMove bool) int32 {
	e.VisitNode()
	outcome := position.Status()
	if outcome == Checkmate {
		return -CHECKMATE_EVAL
	} else if outcome == Draw {
		return 0
	}

	isRootNode := searchHeight == 0
	isPvNode := alpha == beta-1

	if e.StopSearchFlag {
		return beta
	}

	if depthLeft <= 0 {
		return e.quiescence(position, alpha, beta, 0, Evaluate(position), searchHeight)
	}

	// NullMove pruning
	isNullMoveAllowed := !isRootNode && !isPvNode && nullMove && depthLeft >= 5 && !position.IsEndGame() && !position.IsInCheck()
	R := int8(3)
	if searchHeight > 6 {
		R = 2
	}

	if isNullMoveAllowed {
		tempo := int32(15)    // TODO: Make it variable with a formula like: 10*(numPGAM > 0) + 10* numPGAM > 15);
		bound := beta - tempo // variable bound
		position.NullMove()
		newBeta := 1 - bound
		line := NewPVLine(depthLeft - 1 - R)
		score := -e.alphaBeta(position, depthLeft-R-1, searchHeight+1, newBeta-1, newBeta, ply, line, !multiCutFlag, !nullMove)
		position.NullMove()
		if score >= bound {
			return beta // null move pruning
		}
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
			score := -e.alphaBeta(position, depthLeft-1-R, searchHeight+1, newBeta-1, newBeta, ply, line, !multiCutFlag, nullMove)
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
	if !isRootNode && !isPvNode && depthLeft < 2 && eval+margin < beta-1 {
		return e.quiescence(position, alpha, beta, 0, eval, searchHeight)
	}

	// Reverse Futility Pruning
	if !isRootNode && !isPvNode && depthLeft < 5 && eval-margin >= beta {
		return eval - margin /* fail soft */
	}

	// Extended Futility Pruning
	reductionsAllowed := !isRootNode || !isPvNode || position.IsInCheck()
	lastRank := Rank7
	if position.Turn() == Black {
		lastRank = Rank2
	}

	movePicker.Reset()

	for i := 0; i < len(legalMoves); i++ {
		move := movePicker.Next()
		if isRootNode {
			fmt.Printf("info currmove %s currmovenumber %d\n\n", move.ToString(), i+1)
		}

		LMR := int8(0)
		if !reductionsAllowed && searchHeight >= 6 && depthLeft == 2 {
			board := position.Board
			movingPiece := board.PieceAt(move.Source)
			isPromoting := (movingPiece.Type() == Pawn && move.Destination.Rank() == lastRank)

			// Extended Futility Pruning
			gain := Evaluate(position) - eval
			if !isRootNode && !isPvNode && !move.HasTag(Check) && futility+gain <= beta-1 &&
				move.PromoType == NoType && !isPromoting {
				continue
			}

			// Late Move Reduction
			if !isRootNode && !isPvNode && i >= 5 && !move.HasTag(Check) && move.PromoType == NoType && !isPromoting {
				LMR = 1
			}
		}
		capturedPiece, oldEnPassant, oldTag, hc := position.MakeMove(move)
		line := NewPVLine(depthLeft - 1 - LMR)
		score := -MAX_INT
		if searchPv {
			score = -e.alphaBeta(position, depthLeft-1-LMR, searchHeight+1, -beta, -alpha, ply, line, multiCutFlag, nullMove)
		} else {
			score = -e.alphaBeta(position, depthLeft-1-LMR, searchHeight+1, -alpha-1, -alpha, ply, line, multiCutFlag, nullMove)
			if score > alpha {
				score = -e.alphaBeta(position, depthLeft-1-LMR, searchHeight+1, -beta, -alpha, ply, line, multiCutFlag, nullMove) // re-search
			}
		}
		position.UnMakeMove(move, oldTag, oldEnPassant, capturedPiece, hc)
		if score >= beta {
			// Those scores are never useful
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
			searchPv = false
			e.AddMoveHistory(move, position.Board.PieceAt(move.Source), move.Destination, uint16(searchHeight)+ply)
		}
	}
	if searchPv {
		TranspositionTable.Set(hash, &CachedEval{hash, alpha, depthLeft, LowerBound, ply})
	} else {
		TranspositionTable.Set(hash, &CachedEval{hash, alpha, depthLeft, Exact, ply})
	}
	return alpha
}

func nps(nodes int64, dur float64) int64 {
	if dur == 0 {
		return 0
	}
	return int64(float64(nodes) / 1000 * dur)
}
