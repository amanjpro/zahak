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

	bookmove := GetBookMove(position)
	lastDepth := int8(1)

	if bookmove != EmptyMove {
		e.move = bookmove
		e.pv.Recycle()
		e.pv.AddFirst(bookmove)
	}

	if e.move == EmptyMove {
		e.pv.Recycle()
		for iterationDepth := int8(1); iterationDepth <= depth; iterationDepth++ {
			if e.ShouldStop() {
				break
			}
			e.innerLines[0].Recycle()
			score, ok := e.alphaBeta(position, iterationDepth, 0, alpha, beta, ply, EmptyMove, true, true)
			if ok {
				e.pv.Clone(e.innerLines[0])
				e.score = score
				e.move = e.pv.MoveAt(0)
				e.SendPv(iterationDepth)
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

	e.SendPv(lastDepth)
	if e.move == EmptyMove { // we didn't have time to pick a move, pick a random one
		allMoves := position.LegalMoves()
		e.move = allMoves[0]
	}
}

func (e *Engine) alphaBeta(position *Position, depthLeft int8, searchHeight int8, alpha int16, beta int16, ply uint16, currentMove Move, multiCutFlag bool, nullMove bool) (int16, bool) {
	e.VisitNode()

	isRootNode := searchHeight == 0
	isPvNode := alpha != beta-1

	var isInCheck = currentMove.IsCheck()

	// Position is drawn
	if IsRepetition(position, e.pred, currentMove) || position.IsDraw() {
		return 0, true
	}

	if isInCheck && isPvNode {
		e.info.checkExtentionCounter += 1
		depthLeft += 1 // Singular Extension
	}

	if depthLeft <= 0 {
		return e.quiescence(position, alpha, beta, currentMove, 0, Evaluate(position), searchHeight)
	}

	if isPvNode {
		e.info.mainSearchCounter += 1
	} else {
		e.info.zwCounter += 1
	}

	hash := position.Hash()
	nHashMove, nEval, nDepth, nType, found := e.TranspositionTable.Get(hash)
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

	if nHashMove == EmptyMove && !position.HasLegalMoves() {
		if isInCheck {
			return -CHECKMATE_EVAL, true
		} else {
			return 0, true
		}
	}

	if e.ShouldStop() {
		return -MAX_INT, false
	}

	// NullMove pruning
	R := int8(4)
	if depthLeft == 4 {
		R = 3
	}
	isNullMoveAllowed := !isRootNode && !isPvNode && nullMove && depthLeft > R && !position.IsEndGame() && !isInCheck

	if isNullMoveAllowed {
		ep := position.MakeNullMove()
		oldPred := e.pred
		e.pred = NewPredecessors()
		e.innerLines[searchHeight+1].Recycle()
		score, ok := e.alphaBeta(position, depthLeft-R, searchHeight+1, -beta, -beta+1, ply, EmptyMove, !multiCutFlag, false)
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

	var eval int16
	if !isInCheck {
		eval = Evaluate(position)
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

	movePicker := NewMovePicker(position, e, searchHeight, nHashMove, false)

	// Internal Iterative Deepening
	if depthLeft >= 8 && movePicker.HasNoPVMove() {
		e.innerLines[searchHeight].Recycle()
		score, ok := e.alphaBeta(position, depthLeft-7, searchHeight, alpha, beta, ply, currentMove, false, false)
		line := e.innerLines[searchHeight]
		if ok && line.moveCount != 0 && score > alpha && score < beta {
			hashmove := e.innerLines[searchHeight].MoveAt(0)
			movePicker.UpgradeToPvMove(hashmove)
		} else if !ok {
			return score, false
		}
		e.innerLines[searchHeight].Recycle()
	}

	// Multi-Cut Pruning
	M := 6
	C := 3
	R = 4
	if !isRootNode && !isPvNode && depthLeft > R && searchHeight > 3 && multiCutFlag {
		cutNodeCounter := 0
		for i := 0; i < M; i++ {
			move := movePicker.Next()
			if move == EmptyMove {
				break
			}
			oldEnPassant, oldTag, hc := position.MakeMove(move)
			newBeta := 1 - beta
			// newBeta := -beta + 1
			e.pred.Push(position.Hash())
			e.innerLines[searchHeight+1].Recycle()
			score, ok := e.alphaBeta(position, depthLeft-R, searchHeight+1, newBeta-1, newBeta, ply, move, !multiCutFlag, true)
			score = -score
			e.pred.Pop()
			position.UnMakeMove(move, oldTag, oldEnPassant, hc)
			if !ok {
				return score, false
			}
			if score >= beta {
				cutNodeCounter++
				if cutNodeCounter == C {
					e.info.multiCutCounter += 1
					return beta, true // mc-prune
				}
			}
		}
		movePicker.Reset()
	}

	// Extended Futility Pruning
	reductionsAllowed := !isRootNode && !isPvNode && !isInCheck

	hasSeenExact := false

	// using fail soft with negamax:
	move := movePicker.Next()
	oldEnPassant, oldTag, hc := position.MakeMove(move)
	e.pred.Push(position.Hash())
	e.innerLines[searchHeight+1].Recycle()
	bestscore, ok := e.alphaBeta(position, depthLeft-1, searchHeight+1, -beta, -alpha, ply, move, !multiCutFlag, true)
	bestscore = -bestscore
	hashmove := move
	e.pred.Pop()
	position.UnMakeMove(move, oldTag, oldEnPassant, hc)
	if !ok {
		return bestscore, ok
	}
	if bestscore > alpha {
		if bestscore >= beta {
			e.TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, UpperBound, ply)
			e.AddKillerMove(move, searchHeight)
			e.AddMoveHistory(move, move.MovingPiece(), move.Destination(), searchHeight)
			return bestscore, true
		}
		alpha = bestscore
		e.innerLines[searchHeight].AddFirst(move)
		e.innerLines[searchHeight].ReplaceLine(e.innerLines[searchHeight+1])
		hasSeenExact = true
	}

	for i := 1; ; i++ {
		move := movePicker.Next()
		if move == EmptyMove {
			break
		}
		if isRootNode {
			fmt.Printf("info depth %d currmove %s currmovenumber %d\n", depthLeft, move.ToString(), i+1)
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
		e.innerLines[searchHeight+1].Recycle()
		score, ok := e.alphaBeta(position, depthLeft-1-LMR, searchHeight+1, -alpha-1, -alpha, ply, move, !multiCutFlag, true)
		score = -score
		e.pred.Pop()
		if !ok {
			position.UnMakeMove(move, oldTag, oldEnPassant, hc)
			return score, ok
		}
		if score > alpha && score < beta {
			e.info.researchCounter += 1
			// research with window [alpha;beta]
			e.pred.Push(position.Hash())
			e.innerLines[searchHeight+1].Recycle()
			score, ok = e.alphaBeta(position, depthLeft-1, searchHeight+1, -beta, -alpha, ply, move, !multiCutFlag, true)
			score = -score
			e.pred.Pop()
			if !ok {
				position.UnMakeMove(move, oldTag, oldEnPassant, hc)
				return score, ok
			}
			if score > alpha {
				alpha = score
			}
		}
		position.UnMakeMove(move, oldTag, oldEnPassant, hc)

		if score > bestscore {
			if score >= beta {
				e.TranspositionTable.Set(hash, move, score, depthLeft, UpperBound, ply)
				e.AddKillerMove(move, searchHeight)
				e.AddMoveHistory(move, move.MovingPiece(), move.Destination(), searchHeight)
				return score, true
			}

			bestscore = score
			hashmove = move
			// Potential PV move, lets copy it to the current pv-line
			e.innerLines[searchHeight].AddFirst(move)
			e.innerLines[searchHeight].ReplaceLine(e.innerLines[searchHeight+1])
			hasSeenExact = true
		}
	}
	if hasSeenExact {
		e.TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, Exact, ply)
	} else {
		e.TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, LowerBound, ply)
	}
	return bestscore, true
}
