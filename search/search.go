package search

import (
	"fmt"
	"time"

	. "github.com/amanjpro/zahak/book"
	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
)

func (e *Engine) Search(depth int8) {
	e.ClearForSearch()
	e.rootSearch(depth)
}

var p = WhitePawn.Weight()
var r = WhiteRook.Weight()
var b = WhiteBishop.Weight()
var q = WhiteQueen.Weight()

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

			iterStartTime := time.Now()

			e.innerLines[0].Recycle()
			score := e.aspirationWindow(e.score, iterationDepth)
			// score := e.alphaBeta(iterationDepth, 0, -MAX_INT, MAX_INT)

			if e.AbruptStop {
				break
			}
			e.pv.Clone(e.innerLines[0])
			e.score = score
			e.move = e.pv.MoveAt(0)
			e.SendPv(iterationDepth, false)
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

			lastIterationTime := time.Now().Sub(iterStartTime).Milliseconds()

			// We expect the next iteration will take around three times as much as this iteration
			// And, if that means we will exceed the allocated time before we can finish the search,
			// we will be wasting time, this condition tries to avoid this.
			if !e.CanFinishSearch(lastIterationTime) {
				break
			}
		}

	}

	e.SendPv(lastDepth, e.IsBench)
	if e.move == EmptyMove { // we didn't have time to pick a move, pick a random one
		allMoves := e.Position.LegalMoves()
		e.move = allMoves[0]
	}
}

func (e *Engine) aspirationWindow(score int16, iterationDepth int8) int16 {
	if iterationDepth <= 6 {
		alpha := -MAX_INT
		beta := MAX_INT
		score = e.alphaBeta(iterationDepth, 0, alpha, beta)
	} else {
		var alpha, beta int16
		alphaMargin := int16(25)
		betaMargin := int16(25)
		for i := 0; i < 3; i++ {
			if i < 2 {
				alpha = max16(score-alphaMargin, -MAX_INT)
				beta = min16(score+betaMargin, MAX_INT)
			} else {
				alpha = -MAX_INT
				beta = MAX_INT
			}
			score = e.alphaBeta(iterationDepth, 0, alpha, beta)
			if score <= alpha {
				alphaMargin *= 2
			} else if score >= beta {
				betaMargin *= 2
			} else {
				return score
			}
		}
	}
	return score
}

func (e *Engine) alphaBeta(depthLeft int8, searchHeight int8, alpha int16, beta int16) int16 {
	e.VisitNode()

	isRootNode := searchHeight == 0
	isPvNode := alpha != beta-1

	position := e.Position
	var isInCheck = position.IsInCheck()

	currentMove := e.positionMoves[searchHeight]
	// Position is drawn
	if IsRepetition(position, e.pred, currentMove) || position.IsDraw() {
		return 0
	}

	if isInCheck {
		e.info.checkExtentionCounter += 1
		depthLeft += 1 // Singular Extension
	}

	if depthLeft <= 0 {
		e.staticEvals[searchHeight] = Evaluate(position)
		return e.quiescence(alpha, beta, searchHeight)
	}

	if isPvNode {
		e.info.mainSearchCounter += 1
	} else {
		e.info.zwCounter += 1
	}

	hash := position.Hash()
	nHashMove, nEval, nDepth, nType, found := e.TranspositionTable.Get(hash)
	if !isPvNode && found && nDepth >= depthLeft {
		if nEval >= beta && nType == LowerBound {
			e.CacheHit()
			return nEval
		}
		if nEval <= alpha && nType == UpperBound {
			e.CacheHit()
			return nEval
		}
	}

	// if nHashMove == EmptyMove && !position.HasLegalMoves() {
	// 	if isInCheck {
	// 		return -CHECKMATE_EVAL + int16(searchHeight), true
	// 	} else {
	// 		return 0, true
	// 	}
	// }
	//
	if e.ShouldStop() {
		return -MAX_INT
	}

	eval := Evaluate(position)
	e.staticEvals[searchHeight] = eval
	improving := currentMove == EmptyMove ||
		(searchHeight > 2 && e.staticEvals[searchHeight] > e.staticEvals[searchHeight-2])

	// Razoring
	razoringMargin := r
	// if improving {
	// 	razoringMargin += int16(depthLeft) * p
	// }
	if !isRootNode && !isPvNode && currentMove != EmptyMove && !isInCheck && depthLeft <= 3 && eval+razoringMargin < beta {
		e.staticEvals[searchHeight] = eval
		newEval := e.quiescence(alpha, beta, searchHeight)
		if newEval < beta {
			e.info.razoringCounter += 1
			return newEval
		}
	}

	isNullMoveAllowed := !isRootNode && !isPvNode && currentMove != EmptyMove && !isInCheck && !position.IsEndGame()
	// Reverse Futility Pruning
	reverseFutilityMargin := int16(depthLeft) * (b - p)
	// if improving {
	// 	reverseFutilityMargin += int16(depthLeft) * p
	// }
	if isNullMoveAllowed && depthLeft < 7 && eval-reverseFutilityMargin >= beta {
		e.info.rfpCounter += 1
		return eval - reverseFutilityMargin /* fail soft */
	}

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
		e.positionMoves[searchHeight+1] = EmptyMove
		score := -e.alphaBeta(depthLeft-R, searchHeight+1, -beta, -beta+1)
		e.pred = oldPred
		position.UnMakeNullMove(ep)
		if score >= beta { //}&& abs16(score) <= CHECKMATE_EVAL {
			e.info.nullMoveCounter += 1
			// if abs16(score) <= CHECKMATE_EVAL {
			return score
			// }
			// return beta, true // null move pruning
		}
	}

	// Internal Iterative Deepening
	if depthLeft >= 8 && nHashMove == EmptyMove {
		e.innerLines[searchHeight].Recycle()
		score := e.alphaBeta(depthLeft-7, searchHeight, alpha, beta)
		if e.AbruptStop {
			return score
		}
		line := e.innerLines[searchHeight]
		if line.moveCount != 0 { // }&& score > alpha && score < beta {
			hashmove := e.innerLines[searchHeight].MoveAt(0)
			nHashMove = hashmove // movePicker.UpgradeToPvMove(hashmove)
		}
		e.innerLines[searchHeight].Recycle()
	}

	movePicker := e.MovePickers[searchHeight]
	movePicker.RecycleWith(position, e, depthLeft, nHashMove, false)

	// Pruning
	reductionsAllowed := !isRootNode && !isPvNode && !isInCheck

	hasSeenExact := false

	// using fail soft with negamax:
	hashmove := movePicker.Next()
	if hashmove == EmptyMove {
		if isInCheck {
			return -CHECKMATE_EVAL + int16(searchHeight)
		} else {
			return 0
		}
	}
	oldEnPassant, oldTag, hc := position.MakeMove(hashmove)
	e.pred.Push(position.Hash())
	e.innerLines[searchHeight+1].Recycle()
	e.positionMoves[searchHeight+1] = hashmove
	bestscore := -e.alphaBeta(depthLeft-1, searchHeight+1, -beta, -alpha)
	e.pred.Pop()
	position.UnMakeMove(hashmove, oldTag, oldEnPassant, hc)
	if bestscore > alpha {
		e.innerLines[searchHeight].AddFirst(hashmove)
		e.innerLines[searchHeight].ReplaceLine(e.innerLines[searchHeight+1])
		if bestscore >= beta {
			e.TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, LowerBound, e.Ply)
			// e.AddKillerMove(hashmove, searchHeight)
			e.AddHistory(hashmove, hashmove.MovingPiece(), hashmove.Destination(), depthLeft)
			return bestscore
		}
		alpha = bestscore
		hasSeenExact = true
	}
	e.RemoveMoveHistory(hashmove, hashmove.MovingPiece(), hashmove.Destination(), depthLeft)

	pruningThreashold := int(5 + depthLeft*depthLeft)
	if !improving {
		pruningThreashold /= 2
	}

	for i := 1; ; i++ {
		move := movePicker.Next()
		if move == EmptyMove {
			break
		}
		if isRootNode {
			fmt.Printf("info depth %d currmove %s currmovenumber %d\n", depthLeft, move.ToString(), i+1)
		}

		isCheckMove := move.IsCheck()
		isCaptureMove := move.IsCapture()
		promoType := move.PromoType()
		notPromoting := !IsPromoting(move)

		// Late Move Pruning
		if reductionsAllowed && notPromoting && !isCaptureMove && !isCheckMove && depthLeft <= 8 &&
			searchHeight > 5 && i > pruningThreashold && e.KillerMoveScore(move, searchHeight) <= 0 && abs16(bestscore) < CHECKMATE_EVAL {
			e.info.lmpCounter += 1
			continue // LMP
		}

		LMR := int8(0)
		// Late Move Reduction
		if reductionsAllowed && promoType == NoType && !isCaptureMove && !isCheckMove && depthLeft > 3 && i > 4 {
			e.info.lmrCounter += 1
			if i >= 8 && notPromoting {
				LMR = 2
			} else {
				LMR = 1
			}
		}

		oldEnPassant, oldTag, hc := position.MakeMove(move)
		e.pred.Push(position.Hash())
		e.innerLines[searchHeight+1].Recycle()
		e.positionMoves[searchHeight+1] = move
		score := -e.alphaBeta(depthLeft-1-LMR, searchHeight+1, -alpha-1, -alpha)
		e.pred.Pop()
		if score > alpha && score < beta {
			e.info.researchCounter += 1
			// research with window [alpha;beta]
			e.pred.Push(position.Hash())
			e.innerLines[searchHeight+1].Recycle()
			score = -e.alphaBeta(depthLeft-1, searchHeight+1, -beta, -alpha)
			e.pred.Pop()
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
				e.TranspositionTable.Set(hash, move, score, depthLeft, LowerBound, e.Ply)
				// e.AddKillerMove(move, searchHeight)
				e.AddHistory(move, move.MovingPiece(), move.Destination(), depthLeft)
				return score
			}
			bestscore = score
			hashmove = move
			hasSeenExact = true
		}
		e.RemoveMoveHistory(move, move.MovingPiece(), move.Destination(), depthLeft)
	}
	if hasSeenExact {
		e.TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, Exact, e.Ply)
	} else {
		e.TranspositionTable.Set(hash, hashmove, bestscore, depthLeft, UpperBound, e.Ply)
	}
	return bestscore
}

func IsCheckmateEval(eval int16) bool {
	absEval := abs16(eval)
	if absEval == MAX_INT {
		return false
	}
	return absEval >= CHECKMATE_EVAL-int16(MAX_DEPTH)
}
