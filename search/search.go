package search

import (
	"fmt"

	. "github.com/amanjpro/zahak/book"
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

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
	line := NewPVLine(100)

	bookmove := GetBookMove(position)
	lastDepth := int8(1)

	if bookmove != EmptyMove {
		e.move = bookmove
		line.AddFirst(bookmove)
		e.pv = line
	}

	firstScore := true
	if e.move == EmptyMove {
		for iterationDepth := int8(1); iterationDepth <= depth; iterationDepth++ {
			if e.ShouldStop() {
				break
			}
			line.Recycle()
			score, ok := e.alphaBeta(position, iterationDepth, 0, alpha, beta, ply, line, EmptyMove, true, true)
			if ok && (firstScore || line.moveCount >= e.pv.moveCount) {
				e.pv = line
				e.score = score
				e.move = e.pv.MoveAt(0)
				e.SendPv(iterationDepth)
				firstScore = false
			}
			lastDepth = iterationDepth
			if !e.Pondering && iterationDepth >= 35 && e.move == previousBestMove {
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
			e.pred.Clear()
			if !e.Pondering && e.DebugMode {
				e.info.Print()
			}
		}
	}

	if e.Pondering {
		e.SendPv(lastDepth)
	}
}

func (e *Engine) alphaBeta(position *Position, depthLeft int8, searchHeight int8, alpha int16, beta int16, ply uint16, pvline *PVLine, currentMove Move, multiCutFlag bool, nullMove bool) (int16, bool) {
	e.VisitNode()

	isRootNode := searchHeight == 0
	isPvNode := alpha != beta-1

	var isInCheck = currentMove.IsCheck()

	if IsRepetition(position, e.pred, currentMove) {
		return 0, true
	}

	outcome := position.Status()
	if outcome == Checkmate {
		return -CHECKMATE_EVAL, true
	} else if outcome == Draw {
		return 0, true
	}

	if depthLeft == 0 {
		return e.quiescence(position, alpha, beta, currentMove, 0, Evaluate(position), searchHeight)
	}

	if isPvNode {
		e.info.mainSearchCounter += 1
	} else {
		e.info.zwCounter += 1
	}

	hash := position.Hash()
	nHashMove, nEval, nDepth, nType, found := TranspositionTable.Get(hash)
	if found && nDepth >= depthLeft {
		if nEval >= beta && (nType == UpperBound || nType == Exact) {
			e.CacheHit()
			return beta, true
		}
		if nEval <= alpha && (nType == LowerBound || nType == Exact) {
			e.CacheHit()
			return alpha, true
		}
	}

	if e.ShouldStop() {
		return -MAX_INT, false
	}

	var eval int16
	if !isInCheck {
		eval = Evaluate(position)
	}

	// NullMove pruning
	R := int8(4)
	if depthLeft == 4 {
		R = 3
	}
	isNullMoveAllowed := !isRootNode && !isPvNode && nullMove && depthLeft > R && !position.IsEndGame() && !isInCheck

	line := NewPVLine(100)
	if isNullMoveAllowed {
		ep := position.MakeNullMove()
		oldPred := e.pred
		e.pred = NewPredecessors()
		score, ok := e.alphaBeta(position, depthLeft-R, searchHeight+1, -beta, -beta+1, ply, line, EmptyMove, !multiCutFlag, false)
		score = -score
		e.pred = oldPred
		position.UnMakeNullMove(ep)
		if !ok {
			return score, false
		}
		if score >= beta && abs16(score) < CHECKMATE_EVAL {
			e.info.nullMoveCounter += 1
			return beta, true // null move pruning
		}
	}

	// Reverse Futility Pruning
	reverseFutilityMargin := WhiteRook.Weight()
	if !isRootNode && !isPvNode && !isInCheck && depthLeft == 2 && eval-reverseFutilityMargin >= beta {
		e.info.rfpCounter += 1
		return eval - reverseFutilityMargin, true /* fail soft */
	}

	// Razoring
	razoringMargin := 3 * WhitePawn.Weight()
	if depthLeft == 1 {
		razoringMargin = 2 * WhitePawn.Weight()
	}
	if !isRootNode && !isPvNode && depthLeft <= 2 && eval+razoringMargin < beta {
		newEval, ok := e.quiescence(position, alpha, beta, currentMove, 0, eval, searchHeight)
		if !ok {
			return newEval, ok
		}
		if newEval < beta {
			e.info.razoringCounter += 1
			return newEval, true
		}
	}

	legalMoves := position.LegalMoves()

	movePicker := NewMovePicker(position, e, legalMoves, searchHeight, nHashMove)
	line.Recycle()

	// Internal Iterative Deepening
	if depthLeft >= 8 && !movePicker.HasPVMove() {
		e.alphaBeta(position, depthLeft-7, searchHeight+1, alpha, beta, ply, line, currentMove, false, false)
		if line.moveCount != 0 {
			movePicker.UpgradeToPvMove(line.MoveAt(0))
		}
	}

	// Multi-Cut Pruning
	M := 6
	C := 3
	R = 4
	if !isRootNode && !isPvNode && depthLeft > R && searchHeight > 3 && multiCutFlag && len(legalMoves) > M {
		cutNodeCounter := 0
		for i := 0; i < M; i++ {
			line.Recycle()
			move := movePicker.Next()
			oldEnPassant, oldTag, hc := position.MakeMove(move)
			newBeta := 1 - beta
			// newBeta := -beta + 1
			e.pred.Push(position.Hash())
			score, ok := e.alphaBeta(position, depthLeft-R, searchHeight+1, newBeta-1, newBeta, ply, line, move, !multiCutFlag, true)
			score = -score
			e.pred.Pop()
			position.UnMakeMove(move, oldTag, oldEnPassant, hc)
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

	if isInCheck && isPvNode {
		e.info.checkExtentionCounter += 1
		depthLeft += 1 // Singular Extension
	}

	// Extended Futility Pruning
	reductionsAllowed := !isRootNode && !isPvNode && !isInCheck

	movePicker.Reset()

	hasSeenExact := false

	// using fail soft with negamax:
	move := movePicker.Next()
	oldEnPassant, oldTag, hc := position.MakeMove(move)
	line.Recycle()
	e.pred.Push(position.Hash())
	bestscore, ok := e.alphaBeta(position, depthLeft-1, searchHeight+1, -beta, -alpha, ply, line, move, !multiCutFlag, true)
	bestscore = -bestscore
	hashmove := move
	e.pred.Pop()
	position.UnMakeMove(move, oldTag, oldEnPassant, hc)
	if !ok {
		return bestscore, ok
	}
	if bestscore > alpha {
		if bestscore >= beta {
			TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, UpperBound, ply)
			e.AddKillerMove(move, searchHeight)
			return bestscore, true
		}
		alpha = bestscore
		pvline.AddFirst(move)
		pvline.ReplaceLine(line)
		hasSeenExact = true
		e.AddMoveHistory(move, move.MovingPiece(), move.Destination(), searchHeight)
	}

	for i := 1; i < len(legalMoves); i++ {
		line.Recycle()
		move := movePicker.Next()
		if isRootNode {
			fmt.Printf("info depth %d currmove %s currmovenumber %d\n\n", depthLeft, move.ToString(), i+1)
		}

		LMR := int8(0)

		isCheckMove := move.IsCheck()
		isCaptureMove := move.IsCapture()
		promoType := move.PromoType()

		// Extended Futility Pruning
		if reductionsAllowed && !isCheckMove && depthLeft <= 2 && !isCaptureMove &&
			alpha != abs16(CHECKMATE_EVAL) && beta != abs16(CHECKMATE_EVAL) &&
			promoType == NoType {
			margin := BlackBishop.Weight()
			if depthLeft == 2 {
				margin = WhiteRook.Weight()
			}
			if eval+margin <= alpha {
				e.info.efpCounter += 1
				continue
			}
		}

		// Late Move Reduction
		if reductionsAllowed && promoType == NoType && !isCaptureMove && !isCheckMove && depthLeft > 3 && i > 4 {
			e.info.lmrCounter += 1
			LMR = 1
		}
		oldEnPassant, oldTag, hc := position.MakeMove(move)
		e.pred.Push(position.Hash())
		score, ok := e.alphaBeta(position, depthLeft-1-LMR, searchHeight+1, -alpha-1, -alpha, ply, line, move, !multiCutFlag, true)
		score = -score
		e.pred.Pop()
		if !ok {
			position.UnMakeMove(move, oldTag, oldEnPassant, hc)
			return score, ok
		}
		if score > alpha && score < beta {
			line.Recycle()
			e.info.researchCounter += 1
			// research with window [alpha;beta]
			e.pred.Push(position.Hash())
			score, ok = e.alphaBeta(position, depthLeft-1, searchHeight+1, -beta, -alpha, ply, line, move, !multiCutFlag, true)
			score = -score
			e.pred.Pop()
			if !ok {
				position.UnMakeMove(move, oldTag, oldEnPassant, hc)
				return score, ok
			}
			if score > alpha {
				e.AddMoveHistory(move, move.MovingPiece(), move.Destination(), searchHeight)
				alpha = score
			}
		}
		position.UnMakeMove(move, oldTag, oldEnPassant, hc)

		if score > bestscore {
			if score >= beta {
				TranspositionTable.Set(hash, move, score, depthLeft, UpperBound, ply)
				e.AddKillerMove(move, searchHeight)
				return score, ok
			}

			bestscore = score
			hashmove = move
			// Potential PV move, lets copy it to the current pv-line
			pvline.AddFirst(move)
			pvline.ReplaceLine(line)
			hasSeenExact = true
		}
	}
	if hasSeenExact {
		TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, Exact, ply)
	} else {
		TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, LowerBound, ply)
	}
	return bestscore, true
}
