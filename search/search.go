package search

import (
	"fmt"

	. "github.com/amanjpro/zahak/book"
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func (e *Engine) Search(depth int8) {
	e.ClearForSearch()
	e.rootSearch(depth)
}

func (e *Engine) rootSearch(depth int8) {

	var previousBestMove Move
	// alpha := -MAX_INT
	// beta := MAX_INT

	e.move = EmptyMove
	e.score = -MAX_INT //alpha
	fruitelessIterations := 0

	bookmove := GetBookMove(e.Position)
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
			score, ok := e.aspirationWindow(e.score, iterationDepth)
			// score, ok := e.alphaBeta(iterationDepth, 0, -MAX_INT, MAX_INT, EmptyMove)

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
			if IsCheckmateEval(e.score) {
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
		allMoves := e.Position.LegalMoves()
		e.move = allMoves[0]
	}
}

func (e *Engine) aspirationWindow(score int16, iterationDepth int8) (int16, bool) {
	var ok bool
	if iterationDepth <= 6 {
		alpha := -MAX_INT
		beta := MAX_INT
		score, ok = e.alphaBeta(iterationDepth, 0, alpha, beta, EmptyMove)
	} else {
		var alpha, beta int16
		alphaMargin := int16(25)
		betaMargin := int16(25)
		// counter := 0
		for i := 0; ; i++ {
			if i < 2 {
				alpha = max16(score-alphaMargin, -MAX_INT)
				beta = min16(score+betaMargin, MAX_INT)
			} else {
				alpha = -MAX_INT
				beta = MAX_INT
			}
			score, ok = e.alphaBeta(iterationDepth, 0, alpha, beta, EmptyMove)
			if !ok {
				return score, false
			}
			if score <= alpha {
				alphaMargin *= 2
			} else if score >= beta {
				betaMargin *= 2
			} else {
				return score, true
			}
		}
	}
	return score, ok
}

func (e *Engine) alphaBeta(depthLeft int8, searchHeight int8, alpha int16, beta int16, currentMove Move) (int16, bool) {
	e.VisitNode()

	isRootNode := searchHeight == 0
	isPvNode := alpha != beta-1

	position := e.Position
	var isInCheck = position.IsInCheck()

	// Position is drawn
	if IsRepetition(position, e.pred, currentMove) || position.IsDraw() {
		return 0, true
	}

	if isInCheck {
		e.info.checkExtentionCounter += 1
		depthLeft += 1 // Singular Extension
	}

	if depthLeft <= 0 {
		return e.quiescence(alpha, beta, currentMove, Evaluate(position), searchHeight)
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
			return -CHECKMATE_EVAL + int16(searchHeight), true
		} else {
			return 0, true
		}
	}

	if e.ShouldStop() {
		return -MAX_INT, false
	}

	isEndgame := position.IsEndGame()
	isNullMoveAllowed := !isRootNode && !isPvNode && currentMove != EmptyMove && !isInCheck && !isEndgame

	// NullMove pruning
	R := int8(4)
	if depthLeft == 4 {
		R = 3
	}

	if isNullMoveAllowed && depthLeft > R {
		ep := position.MakeNullMove()
		oldPred := e.pred
		e.pred = NewPredecessors()
		e.innerLines[searchHeight+1].Recycle()
		score, ok := e.alphaBeta(depthLeft-R, searchHeight+1, -beta, -beta+1, EmptyMove)
		score = -score
		e.pred = oldPred
		position.UnMakeNullMove(ep)
		if !ok {
			return score, false
		}
		if score >= beta { //}&& abs16(score) <= CHECKMATE_EVAL {
			e.info.nullMoveCounter += 1
			return score, true // null move pruning
		}
	}

	eval := Evaluate(position)
	e.staticEvals[searchHeight] = eval
	improving := currentMove != EmptyMove ||
		searchHeight > 2 && e.staticEvals[searchHeight] > e.staticEvals[searchHeight-2]

	// Reverse Futility Pruning
	reverseFutilityMargin := WhiteRook.Weight()
	if improving {
		reverseFutilityMargin = int16(depthLeft) * WhitePawn.Weight()
	}
	if isNullMoveAllowed && abs16(beta) < CHECKMATE_EVAL && depthLeft <= 3 && eval-reverseFutilityMargin >= beta {
		e.info.rfpCounter += 1
		return eval - reverseFutilityMargin, true /* fail soft */
	}

	// Razoring
	razoringMargin := WhiteRook.Weight()
	if improving {
		razoringMargin = int16(depthLeft) * WhitePawn.Weight()
	}
	if !isRootNode && !isPvNode && depthLeft <= 3 && eval+razoringMargin < beta {
		newEval, ok := e.quiescence(alpha, beta, currentMove, eval, searchHeight)
		if !ok {
			return newEval, false
		}
		if newEval < beta {
			e.info.razoringCounter += 1
			return newEval, true
		}
	}

	// Internal Iterative Deepening
	if depthLeft >= 8 && nHashMove == EmptyMove {
		e.innerLines[searchHeight].Recycle()
		score, ok := e.alphaBeta(depthLeft-7, searchHeight, alpha, beta, currentMove)
		line := e.innerLines[searchHeight]
		if ok && line.moveCount != 0 { // }&& score > alpha && score < beta {
			hashmove := e.innerLines[searchHeight].MoveAt(0)
			nHashMove = hashmove // movePicker.UpgradeToPvMove(hashmove)
		} else if !ok {
			return score, false
		}
		e.innerLines[searchHeight].Recycle()
	}
	movePicker := e.MovePickers[searchHeight]
	movePicker.RecycleWith(position, e, depthLeft, nHashMove, isEndgame, false)

	// Pruning
	reductionsAllowed := !isRootNode && !isPvNode && !isInCheck

	hasSeenExact := false
	hashmove := EmptyMove
	var bestscore int16

	// using fail soft with negamax:
	for hashmove == EmptyMove {
		move := movePicker.Next()
		if move == EmptyMove {
			break
		}
		if oldEnPassant, oldTag, hc, legal := position.MakeMove(&move); legal {
			hashmove = move
			e.pred.Push(position.Hash())
			e.innerLines[searchHeight+1].Recycle()
			score, ok := e.alphaBeta(depthLeft-1, searchHeight+1, -beta, -alpha, hashmove)
			bestscore = -score
			e.pred.Pop()
			position.UnMakeMove(hashmove, oldTag, oldEnPassant, hc)
			if !ok {
				return bestscore, false
			}
			if bestscore > alpha {
				e.innerLines[searchHeight].AddFirst(hashmove)
				e.innerLines[searchHeight].ReplaceLine(e.innerLines[searchHeight+1])
				if bestscore >= beta {
					e.TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, UpperBound, e.Ply)
					e.AddHistory(hashmove, hashmove.MovingPiece(), hashmove.Destination(), depthLeft)
					return bestscore, true
				}
				alpha = bestscore
				hasSeenExact = true
			}
			e.RemoveMoveHistory(hashmove, hashmove.MovingPiece(), hashmove.Destination(), searchHeight)
		}
	}

	pruningThreashold := int(5 + depthLeft*depthLeft)
	if !improving {
		pruningThreashold /= 2
	}

	for i := 1; ; i++ {
		move := movePicker.Next()
		if move == EmptyMove {
			break
		}
		if oldEnPassant, oldTag, hc, legal := position.MakeMove(&move); legal {
			if isRootNode {
				fmt.Printf("info depth %d currmove %s currmovenumber %d\n", depthLeft, move.ToString(), i+1)
			}

			isCheckMove := move.IsCheck()
			isCaptureMove := move.IsCapture()
			promoType := move.PromoType()

			// Late Move Pruning
			if reductionsAllowed && promoType == NoType && !isCaptureMove && !isCheckMove && depthLeft <= 8 &&
				searchHeight > 5 && i > pruningThreashold && e.KillerMoveScore(move, searchHeight) <= 0 && alpha > -(CHECKMATE_EVAL-int16(MAX_DEPTH)) {
				e.info.lmpCounter += 1
				position.UnMakeMove(move, oldTag, oldEnPassant, hc)
				continue // LMP
			}

			LMR := int8(0)
			// Late Move Reduction
			if reductionsAllowed && promoType == NoType && !isCaptureMove && !isCheckMove && depthLeft > 3 && i > 4 {
				e.info.lmrCounter += 1
				if i > 10 {
					LMR = 2
				} else {
					LMR = 1
				}
			}

			e.pred.Push(position.Hash())
			e.innerLines[searchHeight+1].Recycle()
			score, ok := e.alphaBeta(depthLeft-1-LMR, searchHeight+1, -alpha-1, -alpha, move)
			score = -score
			e.pred.Pop()
			if !ok {
				position.UnMakeMove(move, oldTag, oldEnPassant, hc)
				return score, false
			}
			if score > alpha && score < beta {
				e.info.researchCounter += 1
				// research with window [alpha;beta]
				e.pred.Push(position.Hash())
				e.innerLines[searchHeight+1].Recycle()
				v, ok := e.alphaBeta(depthLeft-1, searchHeight+1, -beta, -alpha, move)
				score = -v
				e.pred.Pop()
				if !ok {
					position.UnMakeMove(move, oldTag, oldEnPassant, hc)
					return score, false
				}
				if score > alpha {
					alpha = score
				}
			}
			position.UnMakeMove(move, oldTag, oldEnPassant, hc)

			if score > bestscore {
				// Potential PV move, lets copy it to the current pv-line
				e.innerLines[searchHeight].AddFirst(move)
				e.innerLines[searchHeight].ReplaceLine(e.innerLines[searchHeight+1])
				if score >= beta {
					e.TranspositionTable.Set(hash, move, score, depthLeft, UpperBound, e.Ply)
					e.AddHistory(hashmove, hashmove.MovingPiece(), hashmove.Destination(), depthLeft)
					return score, true
				}
				bestscore = score
				hashmove = move
				hasSeenExact = true
			}
			e.RemoveMoveHistory(hashmove, hashmove.MovingPiece(), hashmove.Destination(), searchHeight)
		} else {
			i -= 1
		}
	}
	if hasSeenExact {
		e.TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, Exact, e.Ply)
	} else {
		e.TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, LowerBound, e.Ply)
	}
	return bestscore, true
}

func IsCheckmateEval(eval int16) bool {
	absEval := abs16(eval)
	if absEval == MAX_INT {
		return false
	}
	return absEval >= CHECKMATE_EVAL-int16(MAX_DEPTH)
}
