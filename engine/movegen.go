package engine

import (
	"math/bits"
)

// func (p *Position) addAllMoves(allMoves *[]Move, ms ...Move) {
// 	*allMoves = append(*allMoves, ms...)
// 	// color := p.Turn()
// 	// for _, m := range ms {
// 	// 	// make the move
// 	// 	// p.partialMakeMove(m)
// 	//
// 	// 	// Does the move puts the moving player in check
// 	// 	// pNotInCheck := !isInCheck(p.Board, color)
// 	//
// 	// 	// if pNotInCheck {
// 	// 	// if isInCheck(p.Board, p.Turn()) { // We put opponent in check
// 	// 	// 	m.AddCheckTag()
// 	// 	// 	*allMoves = append(*allMoves, m)
// 	// 	// } else { // The move does not put us in check
// 	// 	// 	// do nothing
// 	// 	// 	*allMoves = append(*allMoves, m)
// 	// 	// }
// 	// 	// }
// 	// 	p.partialUnMakeMove(m)
// 	// }
// }

func (p *Position) LegalMoves() []Move {
	allMoves := make([]Move, 0, 256)

	allMoves = append(allMoves, p.GetCaptureMoves()...)
	allMoves = append(allMoves, p.GetQuietMoves()...)

	return allMoves
}

func (p *Position) GetQuietMoves() []Move {
	allMoves := make([]Move, 0, 150)
	color := p.Turn()
	board := p.Board

	taboo := tabooSquares(board, color)
	if color == White {
		p.pawnQuietMoves(board.whitePawn, board.whitePieces, board.blackPieces,
			color, false, &allMoves)
		p.knightQuietMoves(WhiteKnight, board.whiteKnight, board.whitePieces, board.blackPieces,
			false, &allMoves)
		p.slidingQuietMoves(board.whiteBishop, board.whitePieces, board.blackPieces,
			color, WhiteBishop, false, &allMoves)
		p.slidingQuietMoves(board.whiteRook, board.whitePieces, board.blackPieces,
			color, WhiteRook, false, &allMoves)
		p.slidingQuietMoves(board.whiteQueen, board.whitePieces, board.blackPieces,
			color, WhiteQueen, false, &allMoves)
		p.kingQuietMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
			taboo, color, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide),
			false, &allMoves)
	} else if color == Black {
		p.pawnQuietMoves(board.blackPawn, board.blackPieces, board.whitePieces,
			color, false, &allMoves)
		p.knightQuietMoves(BlackKnight, board.blackKnight, board.blackPieces, board.whitePieces,
			false, &allMoves)
		p.slidingQuietMoves(board.blackBishop, board.blackPieces, board.whitePieces,
			color, BlackBishop, false, &allMoves)
		p.slidingQuietMoves(board.blackRook, board.blackPieces, board.whitePieces,
			color, BlackRook, false, &allMoves)
		p.slidingQuietMoves(board.blackQueen, board.blackPieces, board.whitePieces,
			color, BlackQueen, false, &allMoves)
		p.kingQuietMoves(board.blackKing, board.blackPieces, board.whitePieces, board.whiteKing,
			taboo, color, p.HasTag(BlackCanCastleKingSide), p.HasTag(BlackCanCastleQueenSide),
			false, &allMoves)
	}

	return allMoves
}

func (p *Position) GetCaptureMoves() []Move {
	allMoves := make([]Move, 0, 150)
	color := p.Turn()
	board := p.Board

	taboo := tabooSquares(board, color)

	if color == White {
		p.pawnCaptureMoves(board.whitePawn, board.whitePieces, board.blackPieces,
			color, false, &allMoves)
		p.knightCaptureMoves(WhiteKnight, board.whiteKnight, board.whitePieces, board.blackPieces,
			false, &allMoves)
		p.slidingCaptureMoves(board.whiteBishop, board.whitePieces, board.blackPieces,
			color, WhiteBishop, false, &allMoves)
		p.slidingCaptureMoves(board.whiteRook, board.whitePieces, board.blackPieces,
			color, WhiteRook, false, &allMoves)
		p.slidingCaptureMoves(board.whiteQueen, board.whitePieces, board.blackPieces,
			color, WhiteQueen, false, &allMoves)
		p.kingCaptureMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
			taboo, color, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide),
			false, &allMoves)
	} else if color == Black {
		p.pawnCaptureMoves(board.blackPawn, board.blackPieces, board.whitePieces,
			color, false, &allMoves)
		p.knightCaptureMoves(BlackKnight, board.blackKnight, board.blackPieces, board.whitePieces,
			false, &allMoves)
		p.slidingCaptureMoves(board.blackBishop, board.blackPieces, board.whitePieces,
			color, BlackBishop, false, &allMoves)
		p.slidingCaptureMoves(board.blackRook, board.blackPieces, board.whitePieces,
			color, BlackRook, false, &allMoves)
		p.slidingCaptureMoves(board.blackQueen, board.blackPieces, board.whitePieces,
			color, BlackQueen, false, &allMoves)
		p.kingCaptureMoves(board.blackKing, board.blackPieces, board.whitePieces, board.whiteKing,
			taboo, color, p.HasTag(BlackCanCastleKingSide), p.HasTag(BlackCanCastleQueenSide),
			false, &allMoves)
	}

	return allMoves
}

func (p *Position) generateMoves(allMoves *[]Move, captureOnly bool) {

	color := p.Turn()
	board := p.Board

	taboo := tabooSquares(board, color)

	if color == White {
		p.pawnCaptureMoves(board.whitePawn, board.whitePieces, board.blackPieces,
			color, false, allMoves)
		p.knightCaptureMoves(WhiteKnight, board.whiteKnight, board.whitePieces, board.blackPieces,
			false, allMoves)
		p.slidingCaptureMoves(board.whiteBishop, board.whitePieces, board.blackPieces,
			color, WhiteBishop, false, allMoves)
		p.slidingCaptureMoves(board.whiteRook, board.whitePieces, board.blackPieces,
			color, WhiteRook, false, allMoves)
		p.slidingCaptureMoves(board.whiteQueen, board.whitePieces, board.blackPieces,
			color, WhiteQueen, false, allMoves)
		p.kingCaptureMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
			taboo, color, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide),
			false, allMoves)
		if !captureOnly {
			p.pawnQuietMoves(board.whitePawn, board.whitePieces, board.blackPieces,
				color, false, allMoves)
			p.knightQuietMoves(WhiteKnight, board.whiteKnight, board.whitePieces, board.blackPieces,
				false, allMoves)
			p.slidingQuietMoves(board.whiteBishop, board.whitePieces, board.blackPieces,
				color, WhiteBishop, false, allMoves)
			p.slidingQuietMoves(board.whiteRook, board.whitePieces, board.blackPieces,
				color, WhiteRook, false, allMoves)
			p.slidingQuietMoves(board.whiteQueen, board.whitePieces, board.blackPieces,
				color, WhiteQueen, false, allMoves)
			p.kingQuietMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
				taboo, color, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide),
				false, allMoves)
		}
	} else if color == Black {
		p.pawnCaptureMoves(board.blackPawn, board.blackPieces, board.whitePieces,
			color, false, allMoves)
		p.knightCaptureMoves(BlackKnight, board.blackKnight, board.blackPieces, board.whitePieces,
			false, allMoves)
		p.slidingCaptureMoves(board.blackBishop, board.blackPieces, board.whitePieces,
			color, BlackBishop, false, allMoves)
		p.slidingCaptureMoves(board.blackRook, board.blackPieces, board.whitePieces,
			color, BlackRook, false, allMoves)
		p.slidingCaptureMoves(board.blackQueen, board.blackPieces, board.whitePieces,
			color, BlackQueen, false, allMoves)
		p.kingCaptureMoves(board.blackKing, board.blackPieces, board.whitePieces, board.whiteKing,
			taboo, color, p.HasTag(BlackCanCastleKingSide), p.HasTag(BlackCanCastleQueenSide),
			false, allMoves)
		if !captureOnly {
			p.pawnQuietMoves(board.blackPawn, board.blackPieces, board.whitePieces,
				color, false, allMoves)
			p.knightQuietMoves(BlackKnight, board.blackKnight, board.blackPieces, board.whitePieces,
				false, allMoves)
			p.slidingQuietMoves(board.blackBishop, board.blackPieces, board.whitePieces,
				color, BlackBishop, false, allMoves)
			p.slidingQuietMoves(board.blackRook, board.blackPieces, board.whitePieces,
				color, BlackRook, false, allMoves)
			p.slidingQuietMoves(board.blackQueen, board.blackPieces, board.whitePieces,
				color, BlackQueen, false, allMoves)
			p.kingQuietMoves(board.blackKing, board.blackPieces, board.whitePieces, board.whiteKing,
				taboo, color, p.HasTag(BlackCanCastleKingSide), p.HasTag(BlackCanCastleQueenSide),
				false, allMoves)
		}
	}
}

func (p *Position) checkMove(m Move) bool {
	color := p.Turn()
	// make the move
	p.partialMakeMove(m)

	// Does the move puts the moving player in check
	pNotInCheck := !isInCheck(p.Board, color)
	p.partialUnMakeMove(m)

	return pNotInCheck
}

func (p *Position) HasLegalMoves() bool {
	color := p.Turn()
	board := p.Board

	taboo := tabooSquares(board, color)

	if color == White {
		return p.kingQuietMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
			taboo, color, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide), true, nil) ||
			p.kingCaptureMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
				taboo, color, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide), true, nil) ||
			p.knightCaptureMoves(WhiteKnight, board.whiteKnight, board.whitePieces, board.blackPieces,
				true, nil) ||
			p.knightQuietMoves(WhiteKnight, board.whiteKnight, board.whitePieces, board.blackPieces,
				true, nil) ||
			p.slidingCaptureMoves(board.whiteBishop, board.whitePieces, board.blackPieces,
				color, WhiteBishop, true, nil) ||
			p.slidingQuietMoves(board.whiteBishop, board.whitePieces, board.blackPieces,
				color, WhiteBishop, true, nil) ||
			p.slidingCaptureMoves(board.whiteRook, board.whitePieces, board.blackPieces,
				color, WhiteRook, true, nil) ||
			p.slidingQuietMoves(board.whiteRook, board.whitePieces, board.blackPieces,
				color, WhiteRook, true, nil) ||
			p.slidingCaptureMoves(board.whiteQueen, board.whitePieces, board.blackPieces,
				color, WhiteQueen, true, nil) ||
			p.slidingQuietMoves(board.whiteQueen, board.whitePieces, board.blackPieces,
				color, WhiteQueen, true, nil) ||
			p.pawnQuietMoves(board.whitePawn, board.whitePieces, board.blackPieces,
				color, true, nil) ||
			p.pawnCaptureMoves(board.whitePawn, board.whitePieces, board.blackPieces,
				color, true, nil)

	} else if color == Black {
		return p.kingQuietMoves(board.blackKing, board.blackPieces, board.whitePieces, board.whiteKing,
			taboo, color, p.HasTag(BlackCanCastleKingSide), p.HasTag(BlackCanCastleQueenSide), true, nil) ||
			p.kingCaptureMoves(board.blackKing, board.blackPieces, board.whitePieces, board.whiteKing,
				taboo, color, p.HasTag(BlackCanCastleKingSide), p.HasTag(BlackCanCastleQueenSide), true, nil) ||
			p.knightCaptureMoves(BlackKnight, board.blackKnight, board.blackPieces, board.whitePieces,
				true, nil) ||
			p.knightQuietMoves(BlackKnight, board.blackKnight, board.blackPieces, board.whitePieces,
				true, nil) ||
			p.slidingCaptureMoves(board.blackBishop, board.blackPieces, board.whitePieces,
				color, BlackBishop, true, nil) ||
			p.slidingQuietMoves(board.blackBishop, board.blackPieces, board.whitePieces,
				color, BlackBishop, true, nil) ||
			p.slidingCaptureMoves(board.blackRook, board.blackPieces, board.whitePieces,
				color, BlackRook, true, nil) ||
			p.slidingQuietMoves(board.blackRook, board.blackPieces, board.whitePieces,
				color, BlackRook, true, nil) ||
			p.slidingCaptureMoves(board.blackQueen, board.blackPieces, board.whitePieces,
				color, BlackQueen, true, nil) ||
			p.slidingQuietMoves(board.blackQueen, board.blackPieces, board.whitePieces,
				color, BlackQueen, true, nil) ||
			p.pawnQuietMoves(board.blackPawn, board.blackPieces, board.whitePieces,
				color, true, nil) ||
			p.pawnCaptureMoves(board.blackPawn, board.blackPieces, board.whitePieces,
				color, true, nil)
	}
	return false
}

// Checks and Pins
func isInCheck(b Bitboard, colorOfKing Color) bool {
	return isKingAttacked(b, colorOfKing)
}

func isKingAttacked(b Bitboard, colorOfKing Color) bool {
	var ownKing, opPawnAttacks, opKnights, opRQ, opBQ uint64
	var squareOfKing Square
	occupiedBB := b.whitePieces | b.blackPieces
	if colorOfKing == White {
		kingIndex := bitScanForward(b.whiteKing)
		ownKing = squareMask[kingIndex]
		squareOfKing = Square(kingIndex)
		opPawnAttacks = wPawnsAble2CaptureAny(ownKing, b.blackPawn)
		opKnights = b.blackKnight
		opRQ = b.blackRook | b.blackQueen
		opBQ = b.blackBishop | b.blackQueen
	} else {
		kingIndex := bitScanForward(b.blackKing)
		ownKing = squareMask[kingIndex]
		squareOfKing = Square(kingIndex)
		opPawnAttacks = bPawnsAble2CaptureAny(ownKing, b.whitePawn)
		opKnights = b.whiteKnight
		opRQ = b.whiteRook | b.whiteQueen
		opBQ = b.whiteBishop | b.whiteQueen
	}
	pawnChecks := opPawnAttacks
	if pawnChecks != 0 {
		return true
	}

	knightChecks := (knightAttacks(ownKing) & opKnights)
	if knightChecks != 0 {
		return true
	}
	// Knights and pawns cannot discover each other
	bishopChecks := (bishopAttacks(squareOfKing, occupiedBB, empty) & opBQ)
	if bishopChecks != 0 {
		return true
	}
	rookChecks := (rookAttacks(squareOfKing, occupiedBB, empty) & opRQ)
	if rookChecks != 0 {
		return true
	}

	return false
}

// func isDoubleCheck(b Bitboard, colorOfKing Color) bool {
// 	return isKingAttacked(b, colorOfKing, true)
// }
//
func tabooSquares(b Bitboard, colorOfKing Color) uint64 {
	var opPawns, opKnights, opR, opB, opQ, opKing, opPieces uint64
	occupiedBB := b.whitePieces | b.blackPieces
	if colorOfKing == White {
		opPawns = bPawnsAble2CaptureAny(b.blackPawn, universal)
		opKnights = b.blackKnight
		opR = b.blackRook
		opB = b.blackBishop
		opQ = b.blackQueen
		opKing = b.blackKing
		opPieces = b.blackPieces
	} else {
		opPawns = wPawnsAble2CaptureAny(b.whitePawn, universal)
		opKnights = b.whiteKnight
		opR = b.whiteRook
		opB = b.whiteBishop
		opQ = b.whiteQueen
		opKing = b.whiteKing
		opPieces = b.whitePieces
	}
	taboo := opPawns | computedKingAttacks[bitScanForward(opKing)]

	for opKnights != 0 {
		sq := bitScanForward(opKnights)
		taboo |= computedKnightAttacks[sq]
		opKnights ^= squareMask[sq]
	}

	for opB != 0 {
		sq := bitScanForward(opB)
		taboo |= bishopAttacks(Square(sq), occupiedBB, opPieces)
		opB ^= squareMask[sq]
	}

	for opR != 0 {
		sq := bitScanForward(opR)
		taboo |= rookAttacks(Square(sq), occupiedBB, opPieces)
		opR ^= squareMask[sq]
	}

	for opQ != 0 {
		sq := bitScanForward(opQ)
		taboo |= queenAttacks(Square(sq), occupiedBB, opPieces)
		opQ ^= squareMask[sq]
	}

	return taboo
}

// Quiet moves

func (p *Position) pawnQuietMoves(bbPawn uint64, ownPieces uint64, otherPieces uint64, color Color,
	isLegalityCheck bool, allMoves *[]Move) bool {
	emptySquares := (otherPieces | ownPieces) ^ universal
	if color == White {
		bbPawn &^= rank7
		for bbPawn != 0 {
			src := bitScanForward(bbPawn)
			srcSq := Square(src)
			pawn := squareMask[src]
			if srcSq.Rank() == Rank2 {
				dbl := wDoublePushTargets(pawn, emptySquares)
				if dbl != 0 {
					dest := Square(bitScanForward(dbl))
					var tag MoveTag = 0
					m := NewMove(srcSq, dest, WhitePawn, NoPiece, NoType, tag)
					if isLegalityCheck && p.checkMove(m) {
						return true
					} else if !isLegalityCheck {
						*allMoves = append(*allMoves, m)
					}
				}
			}
			sngl := wSinglePushTargets(pawn, emptySquares)
			if sngl != 0 {
				dest := Square(bitScanForward(sngl))
				m := NewMove(srcSq, dest, WhitePawn, NoPiece, NoType, 0)
				if isLegalityCheck && p.checkMove(m) {
					return true
				} else if !isLegalityCheck {
					*allMoves = append(*allMoves, m)
				}
			}
			bbPawn ^= pawn
		}
	} else if color == Black {
		bbPawn &^= rank2
		for bbPawn != 0 {
			src := bitScanForward(bbPawn)
			srcSq := Square(src)
			pawn := squareMask[src]
			dbl := bDoublePushTargets(pawn, emptySquares)
			if dbl != 0 {
				dest := Square(bitScanForward(dbl))
				var tag MoveTag = 0
				m := NewMove(srcSq, dest, BlackPawn, NoPiece, NoType, tag)
				if isLegalityCheck && p.checkMove(m) {
					return true
				} else if !isLegalityCheck {
					*allMoves = append(*allMoves, m)
				}
			}
			sngl := bSinglePushTargets(pawn, emptySquares)
			if sngl != 0 {
				dest := Square(bitScanForward(sngl))
				m := NewMove(srcSq, dest, BlackPawn, NoPiece, NoType, 0)
				if isLegalityCheck && p.checkMove(m) {
					return true
				} else if !isLegalityCheck {
					*allMoves = append(*allMoves, m)
				}
			}
			bbPawn ^= pawn
		}
	}

	return false
}

func (p *Position) knightQuietMoves(movingPiece Piece, bbPiece uint64, ownPieces uint64, otherPieces uint64,
	isLegalityCheck bool, allMoves *[]Move) bool {
	both := otherPieces | ownPieces
	for bbPiece != 0 {
		src := bitScanForward(bbPiece)
		srcSq := Square(src)
		knight := squareMask[src]
		moves := knightMovesNoCaptures(srcSq, both)
		for moves != 0 {
			sq := bitScanForward(moves)
			dest := Square(sq)
			m := NewMove(srcSq, dest, movingPiece, NoPiece, NoType, 0)
			if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
				return true
			} else if !isLegalityCheck {
				*allMoves = append(*allMoves, m)
			}
			moves ^= squareMask[sq]
		}
		bbPiece ^= knight
	}

	return false
}

func (p *Position) slidingQuietMoves(bbPiece uint64, ownPieces uint64, otherPieces uint64,
	color Color, movingPiece Piece, isLegalityCheck bool, allMoves *[]Move) bool {
	both := otherPieces | ownPieces
	var rayAttacks uint64
	for bbPiece != 0 {
		src := bitScanForward(bbPiece)
		srcSq := Square(src)
		switch movingPiece {
		case WhiteBishop, BlackBishop:
			rayAttacks = bishopAttacks(srcSq, both, ownPieces)
		case WhiteRook, BlackRook:
			rayAttacks = rookAttacks(srcSq, both, ownPieces)
		case WhiteQueen, BlackQueen:
			rayAttacks = queenAttacks(srcSq, both, ownPieces)
		}
		passiveMoves := rayAttacks &^ otherPieces
		for passiveMoves != 0 {
			sq := bitScanForward(passiveMoves)
			dest := Square(sq)
			m := NewMove(srcSq, dest, movingPiece, NoPiece, NoType, 0)
			if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
				return true
			} else if !isLegalityCheck {
				*allMoves = append(*allMoves, m)
			}
			passiveMoves ^= squareMask[sq]
		}

		bbPiece ^= squareMask[src]
	}
	return false
}

func (p *Position) kingQuietMoves(bbPiece uint64, ownPieces uint64, otherPieces uint64, otherKing uint64,
	tabooSquares uint64, color Color, kingSideCastle bool, queenSideCastle bool,
	isLegalityCheck bool, allMoves *[]Move) bool {
	both := (otherPieces | ownPieces)
	var movingPiece = BlackKing
	if color == White {
		movingPiece = WhiteKing
	}
	if bbPiece != 0 {
		src := bitScanForward(bbPiece)
		srcSq := Square(src)
		moves := kingMovesNoCaptures(srcSq, both, tabooSquares)
		for moves != 0 {
			sq := bitScanForward(moves)
			dest := Square(sq)
			m := NewMove(srcSq, dest, movingPiece, NoPiece, NoType, 0)
			if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
				return true
			} else if !isLegalityCheck {
				*allMoves = append(*allMoves, m)
			}

			moves ^= squareMask[sq]
		}

		E := E1
		F := F1
		G := G1
		D := D1
		C := C1
		B := B1
		if color == Black && srcSq.Rank() == Rank8 {
			E = E8
			F = F8
			G = G8
			D = D8
			C = C8
			B = B8
		}

		kingSide := uint64(squareMask[F] | squareMask[G])
		queenSide := uint64(squareMask[D] | squareMask[C])

		if srcSq == E && kingSideCastle &&
			((ownPieces|otherPieces)&kingSide == 0) && // are empty
			(tabooSquares&(kingSide|squareMask[E]) == 0) { // Not in check
			m := NewMove(srcSq, G, movingPiece, NoPiece, NoType, KingSideCastle)
			if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
				return true
			} else if !isLegalityCheck {
				*allMoves = append(*allMoves, m)
			}
		}

		if srcSq == E && queenSideCastle &&
			((ownPieces|otherPieces)&(queenSide|(squareMask[B])) == 0) && // are empty
			(tabooSquares&(queenSide|squareMask[E]) == 0) { // Not in check
			m := NewMove(srcSq, C, movingPiece, NoPiece, NoType, QueenSideCastle)
			if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
				return true
			} else if !isLegalityCheck {
				*allMoves = append(*allMoves, m)
			}

		}
	}

	return false
}

// Capture moves

func (p *Position) pawnCaptureMoves(bbPawn uint64, ownPieces uint64, otherPieces uint64, color Color,
	isLegalityCheck bool, allMoves *[]Move) bool {
	emptySquares := (otherPieces | ownPieces) ^ universal
	enPassant := p.EnPassant
	if color == White {
		for bbPawn != 0 {
			src := bitScanForward(bbPawn)
			srcSq := Square(src)
			pawn := squareMask[src]
			attacks := wPawnsAble2CaptureAny(pawn, otherPieces)
			for attacks != 0 {
				sq := bitScanForward(attacks)
				dest := Square(sq)
				cp := p.Board.PieceAt(dest)
				if dest.Rank() == Rank8 {
					m1 := NewMove(srcSq, dest, WhitePawn, cp, Queen, Capture)
					m2 := NewMove(srcSq, dest, WhitePawn, cp, Rook, Capture)
					m3 := NewMove(srcSq, dest, WhitePawn, cp, Bishop, Capture)
					m4 := NewMove(srcSq, dest, WhitePawn, cp, Knight, Capture)
					if isLegalityCheck && p.checkMove(m1) { // if one is illegal, they all are
						return true
					} else if !isLegalityCheck {
						*allMoves = append(*allMoves, m1, m2, m3, m4)
					}
				} else {
					m := NewMove(srcSq, dest, WhitePawn, cp, NoType, Capture)
					if isLegalityCheck && p.checkMove(m) {
						return true
					} else if !isLegalityCheck {
						*allMoves = append(*allMoves, m)
					}
				}
				attacks ^= squareMask[sq]
			}
			if srcSq.Rank() == Rank5 && enPassant != NoSquare && enPassant.Rank() == Rank6 {
				ep := squareMask[enPassant]
				r := wPawnsAble2CaptureAny(pawn, ep)
				if r != 0 {
					dest := Square(bitScanForward(r))
					var tag MoveTag = Capture | EnPassant
					m := NewMove(srcSq, dest, WhitePawn, BlackPawn, NoType, tag)
					if isLegalityCheck && p.checkMove(m) {
						return true
					} else if !isLegalityCheck {
						*allMoves = append(*allMoves, m)
					}
				}
			}
			sngl := wSinglePushTargets(pawn&rank7, emptySquares)
			if sngl != 0 {
				dest := Square(bitScanForward(sngl))
				m1 := NewMove(srcSq, dest, WhitePawn, NoPiece, Queen, 0)
				m2 := NewMove(srcSq, dest, WhitePawn, NoPiece, Rook, 0)
				m3 := NewMove(srcSq, dest, WhitePawn, NoPiece, Bishop, 0)
				m4 := NewMove(srcSq, dest, WhitePawn, NoPiece, Knight, 0)
				if isLegalityCheck && p.checkMove(m1) { // if one is illegal, they all are illegal
					return true
				} else if !isLegalityCheck {
					*allMoves = append(*allMoves, m1, m2, m3, m4)
				}
			}
			bbPawn ^= pawn
		}
	} else if color == Black {
		for bbPawn != 0 {
			src := bitScanForward(bbPawn)
			srcSq := Square(src)
			pawn := squareMask[src]
			attacks := bPawnsAble2CaptureAny(pawn, otherPieces)
			for attacks != 0 {
				sq := bitScanForward(attacks)
				dest := Square(sq)
				cp := p.Board.PieceAt(dest)
				if dest.Rank() == Rank1 {
					m1 := NewMove(srcSq, dest, BlackPawn, cp, Queen, Capture)
					m2 := NewMove(srcSq, dest, BlackPawn, cp, Rook, Capture)
					m3 := NewMove(srcSq, dest, BlackPawn, cp, Bishop, Capture)
					m4 := NewMove(srcSq, dest, BlackPawn, cp, Knight, Capture)
					if isLegalityCheck && p.checkMove(m1) { // if one is illegal, they all are illegal
						return true
					} else if !isLegalityCheck {
						*allMoves = append(*allMoves, m1, m2, m3, m4)
					}
				} else {
					var tag MoveTag = Capture
					m := NewMove(srcSq, dest, BlackPawn, cp, NoType, tag)
					if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
						return true
					} else if !isLegalityCheck {
						*allMoves = append(*allMoves, m)
					}
				}
				attacks ^= squareMask[sq]
			}
			if srcSq.Rank() == Rank4 && enPassant != NoSquare && enPassant.Rank() == Rank3 {
				ep := squareMask[enPassant]
				r := bPawnsAble2CaptureAny(pawn, ep)
				if r != 0 {
					dest := Square(bitScanForward(r))
					var tag MoveTag = Capture | EnPassant
					m := NewMove(srcSq, dest, BlackPawn, WhitePawn, NoType, tag)
					if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
						return true
					} else if !isLegalityCheck {
						*allMoves = append(*allMoves, m)
					}
				}
			}
			sngl := bSinglePushTargets(pawn&rank2, emptySquares)
			if sngl != 0 {
				dest := Square(bitScanForward(sngl))
				m1 := NewMove(srcSq, dest, BlackPawn, NoPiece, Queen, 0)
				m2 := NewMove(srcSq, dest, BlackPawn, NoPiece, Rook, 0)
				m3 := NewMove(srcSq, dest, BlackPawn, NoPiece, Bishop, 0)
				m4 := NewMove(srcSq, dest, BlackPawn, NoPiece, Knight, 0)
				if isLegalityCheck && p.checkMove(m1) {
					return true
				} else if !isLegalityCheck {
					*allMoves = append(*allMoves, m1, m2, m3, m4)
				}
			}
			bbPawn ^= pawn
		}
	}

	return false
}

func (p *Position) knightCaptureMoves(movingPiece Piece, bbPiece uint64, ownPieces uint64, otherPieces uint64,
	isLegalityCheck bool, allMoves *[]Move) bool {
	for bbPiece != 0 {
		src := bitScanForward(bbPiece)
		srcSq := Square(src)
		knight := squareMask[src]
		captures := knightCaptures(srcSq, otherPieces)
		for captures != 0 {
			sq := bitScanForward(captures)
			dest := Square(sq)
			cp := p.Board.PieceAt(dest)
			m := NewMove(srcSq, dest, movingPiece, cp, NoType, Capture)
			if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
				return true
			} else if !isLegalityCheck {
				*allMoves = append(*allMoves, m)
			}

			captures ^= squareMask[sq]
		}
		bbPiece ^= knight
	}

	return false
}

func (p *Position) slidingCaptureMoves(bbPiece uint64, ownPieces uint64, otherPieces uint64,
	color Color, movingPiece Piece, isLegalityCheck bool, allMoves *[]Move) bool {
	both := otherPieces | ownPieces
	var rayAttacks uint64
	for bbPiece != 0 {
		src := bitScanForward(bbPiece)
		srcSq := Square(src)
		switch movingPiece {
		case WhiteBishop, BlackBishop:
			rayAttacks = bishopAttacks(srcSq, both, ownPieces)
		case WhiteRook, BlackRook:
			rayAttacks = rookAttacks(srcSq, both, ownPieces)
		case WhiteQueen, BlackQueen:
			rayAttacks = queenAttacks(srcSq, both, ownPieces)
		}
		captureMoves := rayAttacks & otherPieces
		for captureMoves != 0 {
			sq := bitScanForward(captureMoves)
			dest := Square(sq)
			cp := p.Board.PieceAt(dest)
			m := NewMove(srcSq, dest, movingPiece, cp, NoType, Capture)
			if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
				return true
			} else if !isLegalityCheck {
				*allMoves = append(*allMoves, m)
			}

			captureMoves ^= squareMask[sq]
		}
		bbPiece ^= squareMask[src]
	}
	return false
}

func (p *Position) kingCaptureMoves(bbPiece uint64, ownPieces uint64, otherPieces uint64, otherKing uint64,
	tabooSquares uint64, color Color, kingSideCastle bool, queenSideCastle bool,
	isLegalityCheck bool, allMoves *[]Move) bool {
	var movingPiece = BlackKing
	if color == White {
		movingPiece = WhiteKing
	}
	if bbPiece != 0 {
		src := bitScanForward(bbPiece)
		srcSq := Square(src)
		captures := kingCaptures(srcSq, otherPieces, tabooSquares)
		for captures != 0 {
			sq := bitScanForward(captures)
			dest := Square(sq)
			cp := p.Board.PieceAt(dest)
			m := NewMove(srcSq, dest, movingPiece, cp, NoType, Capture)
			if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
				return true
			} else if !isLegalityCheck {
				*allMoves = append(*allMoves, m)
			}
			captures ^= squareMask[sq]
		}
	}

	return false
}

// Pawn Pushes

func wSinglePushTargets(wpawns uint64, empty uint64) uint64 {
	return nortOne(wpawns) & empty
}

func wDoublePushTargets(wpawns uint64, empty uint64) uint64 {
	singlePushs := wSinglePushTargets(wpawns, empty)
	return nortOne(singlePushs) & empty & rank4
}

func bSinglePushTargets(bpawns uint64, empty uint64) uint64 {
	return soutOne(bpawns) & empty
}

func bDoublePushTargets(bpawns uint64, empty uint64) uint64 {
	singlePushs := bSinglePushTargets(bpawns, empty)
	return soutOne(singlePushs) & empty & rank5
}

func wPawnAnyAttacks(wpawns uint64) uint64 {
	return noEaOne(wpawns) | noWeOne(wpawns)
}

func bPawnAnyAttacks(bpawns uint64) uint64 {
	return soEaOne(bpawns) | soWeOne(bpawns)
}

func wPawnsAble2CaptureAny(wpawns uint64, bpieces uint64) uint64 {
	return wPawnAnyAttacks(wpawns) & bpieces
}

func bPawnsAble2CaptureAny(bpawns uint64, wpieces uint64) uint64 {
	return bPawnAnyAttacks(bpawns) & wpieces
}

// Sliding pieces

func getPositiveRayAttacks(sq Square, occupied uint64, dir Direction) uint64 {
	positiveAttacks := rayAttacksArray[dir][sq]
	blocker := positiveAttacks & occupied
	if blocker != 0 {
		square := bitScanForward(blocker)
		positiveAttacks ^= rayAttacksArray[dir][square]
	}
	return positiveAttacks
}

func getNegativeRayAttacks(sq Square, occupied uint64, dir Direction) uint64 {
	negativeAttacks := rayAttacksArray[dir][sq]
	blocker := negativeAttacks & occupied
	if blocker != 0 {
		square := bitScanReverse(blocker)
		negativeAttacks ^= rayAttacksArray[dir][square]
	}
	return negativeAttacks
}

// The mighty knight

func knightCheckTag(from Square, otherKing uint64) MoveTag {
	if knightCaptures(from, otherKing) != 0 {
		return Check
	}
	return 0
}

var computedKnightAttacks [64]uint64

func knightMovesNoCaptures(sq Square, other uint64) uint64 {
	attacks := computedKnightAttacks[sq]
	return attacks &^ other
}

func knightCaptures(sq Square, other uint64) uint64 {
	return computedKnightAttacks[sq] & other
}

func knightAttacks(b uint64) uint64 {
	return noNoEa(b) | noEaEa(b) | soEaEa(b) | soSoEa(b) |
		noNoWe(b) | noWeWe(b) | soWeWe(b) | soSoWe(b)
}

// King & Kingslayer
func kingMovesNoCaptures(sq Square, others uint64, tabooSquares uint64) uint64 {
	attacks := computedKingAttacks[sq]
	return attacks &^ (others | tabooSquares)
}

func kingCaptures(sq Square, others uint64, tabooSquares uint64) uint64 {
	attacks := computedKingAttacks[sq]
	return (attacks & others) &^ tabooSquares
}

func kingAttacks(b uint64) uint64 {
	return soutOne(b) | nortOne(b) | eastOne(b) | noEaOne(b) |
		soEaOne(b) | westOne(b) | soWeOne(b) | noWeOne(b)
}

var computedKingAttacks [64]uint64

// Utilites
func bitScanForward(bb uint64) uint8 {
	return uint8(bits.TrailingZeros64(bb))
}

func bitScanReverse(bb uint64) uint8 {
	return uint8(bits.LeadingZeros64(bb) ^ 63)
}

// directions

type Direction uint8

const (
	North Direction = iota
	NorthEast
	East
	SouthEast
	South
	SouthWest
	West
	NorthWest
)

var rayAttacksArray [8][64]uint64

func southRay(b uint64) uint64 {
	res := uint64(0)
	b = soutOne(b)
	for b != 0 {
		res |= b
		b = soutOne(b)
	}
	return res
}

func northRay(b uint64) uint64 {
	res := uint64(0)
	b = nortOne(b)
	for b != 0 {
		res |= b
		b = nortOne(b)
	}
	return res
}

func northEastRay(b uint64) uint64 {
	res := uint64(0)
	b = noEaOne(b)
	for b != 0 {
		res |= b
		b = noEaOne(b)
	}
	return res
}

func northWestRay(b uint64) uint64 {
	res := uint64(0)
	b = noWeOne(b)
	for b != 0 {
		res |= b
		b = noWeOne(b)
	}
	return res
}

func westRay(b uint64) uint64 {
	res := uint64(0)
	b = westOne(b)
	for b != 0 {
		res |= b
		b = westOne(b)
	}
	return res
}

func eastRay(b uint64) uint64 {
	res := uint64(0)
	b = eastOne(b)
	for b != 0 {
		res |= b
		b = eastOne(b)
	}
	return res
}

func southEastRay(b uint64) uint64 {
	res := uint64(0)
	b = soEaOne(b)
	for b != 0 {
		res |= b
		b = soEaOne(b)
	}
	return res
}

func southWestRay(b uint64) uint64 {
	res := uint64(0)
	b = soWeOne(b)
	for b != 0 {
		res |= b
		b = soWeOne(b)
	}
	return res
}

func soutOne(b uint64) uint64 {
	return b >> 8
}
func nortOne(b uint64) uint64 {
	return b << 8
}

func noEaOne(b uint64) uint64 {
	return (b << 9) & notAFile
}
func soEaOne(b uint64) uint64 {
	return (b >> 7) & notAFile
}
func westOne(b uint64) uint64 {
	return (b >> 1) & notHFile
}
func soWeOne(b uint64) uint64 {
	return (b >> 9) & notHFile
}
func noWeOne(b uint64) uint64 {
	return (b << 7) & notHFile
}

func eastOne(b uint64) uint64 {
	return (b << 1) & notAFile
}

func noNoEa(b uint64) uint64 {
	return (b << 17) & notAFile
}
func noEaEa(b uint64) uint64 {
	return (b << 10) & notABFile
}
func soEaEa(b uint64) uint64 {
	return (b >> 6) & notABFile
}
func soSoEa(b uint64) uint64 {
	return (b >> 15) & notAFile
}
func noNoWe(b uint64) uint64 {
	return (b << 15) & notHFile
}
func noWeWe(b uint64) uint64 {
	return (b << 6) & notGHFile
}
func soWeWe(b uint64) uint64 {
	return (b >> 10) & notGHFile
}
func soSoWe(b uint64) uint64 {
	return (b >> 17) & notHFile
}

const empty = uint64(0)
const universal = uint64(0xffffffffffffffff)
const notAFile = uint64(0xfefefefefefefefe) // ~0x0101010101010101
const notBFile = uint64(0xfdfdfdfdfdfdfdfd)
const notGFile = uint64(0xbfbfbfbfbfbfbfbf)
const notHFile = uint64(0x7f7f7f7f7f7f7f7f) // ~0x8080808080808080
const notABFile = notAFile & notBFile
const notGHFile = notGFile & notHFile
const rank2 = uint64(0x000000000000FF00)
const rank4 = uint64(0x00000000FF000000)
const rank5 = uint64(0x000000FF00000000)
const rank7 = uint64(0x00FF000000000000)

// I took those from CounterGo, which in turn takes them from Chess Programming Wiki
func bishopAttacks(sq Square, occ uint64, ownPieces uint64) uint64 {
	from := int(sq)
	return bishopAttacksArray[from][((bishopMask[from]&occ)*bishopMult[from])>>bishopShift] &^ ownPieces
}

func rookAttacks(sq Square, occ uint64, ownPieces uint64) uint64 {
	from := int(sq)
	return rookAttacksArray[from][((rookMask[from]&occ)*rookMult[from])>>rookShift] &^ ownPieces
}

func queenAttacks(from Square, occ uint64, ownPieces uint64) uint64 {
	return bishopAttacks(from, occ, ownPieces) | rookAttacks(from, occ, ownPieces)
}

const (
	bishopShift = 55
	rookShift   = 52
)

var rookMult = [...]uint64{
	0x0080001020400080, 0x0040001000200040, 0x0080081000200080, 0x0080040800100080,
	0x0080020400080080, 0x0080010200040080, 0x0080008001000200, 0x0080002040800100,
	0x0000800020400080, 0x0000400020005000, 0x0000801000200080, 0x0000800800100080,
	0x0000800400080080, 0x0000800200040080, 0x0000800100020080, 0x0000800040800100,
	0x0000208000400080, 0x0000404000201000, 0x0000808010002000, 0x0000808008001000,
	0x0000808004000800, 0x0000808002000400, 0x0000010100020004, 0x0000020000408104,
	0x0000208080004000, 0x0000200040005000, 0x0000100080200080, 0x0000080080100080,
	0x0000040080080080, 0x0000020080040080, 0x0000010080800200, 0x0000800080004100,
	0x0000204000800080, 0x0000200040401000, 0x0000100080802000, 0x0000080080801000,
	0x0000040080800800, 0x0000020080800400, 0x0000020001010004, 0x0000800040800100,
	0x0000204000808000, 0x0000200040008080, 0x0000100020008080, 0x0000080010008080,
	0x0000040008008080, 0x0000020004008080, 0x0000010002008080, 0x0000004081020004,
	0x0000204000800080, 0x0000200040008080, 0x0000100020008080, 0x0000080010008080,
	0x0000040008008080, 0x0000020004008080, 0x0000800100020080, 0x0000800041000080,
	0x00FFFCDDFCED714A, 0x007FFCDDFCED714A, 0x003FFFCDFFD88096, 0x0000040810002101,
	0x0001000204080011, 0x0001000204000801, 0x0001000082000401, 0x0001FFFAABFAD1A2,
}

var rookMask = [...]uint64{
	0x000101010101017E, 0x000202020202027C, 0x000404040404047A, 0x0008080808080876,
	0x001010101010106E, 0x002020202020205E, 0x004040404040403E, 0x008080808080807E,
	0x0001010101017E00, 0x0002020202027C00, 0x0004040404047A00, 0x0008080808087600,
	0x0010101010106E00, 0x0020202020205E00, 0x0040404040403E00, 0x0080808080807E00,
	0x00010101017E0100, 0x00020202027C0200, 0x00040404047A0400, 0x0008080808760800,
	0x00101010106E1000, 0x00202020205E2000, 0x00404040403E4000, 0x00808080807E8000,
	0x000101017E010100, 0x000202027C020200, 0x000404047A040400, 0x0008080876080800,
	0x001010106E101000, 0x002020205E202000, 0x004040403E404000, 0x008080807E808000,
	0x0001017E01010100, 0x0002027C02020200, 0x0004047A04040400, 0x0008087608080800,
	0x0010106E10101000, 0x0020205E20202000, 0x0040403E40404000, 0x0080807E80808000,
	0x00017E0101010100, 0x00027C0202020200, 0x00047A0404040400, 0x0008760808080800,
	0x00106E1010101000, 0x00205E2020202000, 0x00403E4040404000, 0x00807E8080808000,
	0x007E010101010100, 0x007C020202020200, 0x007A040404040400, 0x0076080808080800,
	0x006E101010101000, 0x005E202020202000, 0x003E404040404000, 0x007E808080808000,
	0x7E01010101010100, 0x7C02020202020200, 0x7A04040404040400, 0x7608080808080800,
	0x6E10101010101000, 0x5E20202020202000, 0x3E40404040404000, 0x7E80808080808000,
}

var bishopMult = [...]uint64{
	0x0002020202020200, 0x0002020202020000, 0x0004010202000000, 0x0004040080000000,
	0x0001104000000000, 0x0000821040000000, 0x0000410410400000, 0x0000104104104000,
	0x0000040404040400, 0x0000020202020200, 0x0000040102020000, 0x0000040400800000,
	0x0000011040000000, 0x0000008210400000, 0x0000004104104000, 0x0000002082082000,
	0x0004000808080800, 0x0002000404040400, 0x0001000202020200, 0x0000800802004000,
	0x0000800400A00000, 0x0000200100884000, 0x0000400082082000, 0x0000200041041000,
	0x0002080010101000, 0x0001040008080800, 0x0000208004010400, 0x0000404004010200,
	0x0000840000802000, 0x0000404002011000, 0x0000808001041000, 0x0000404000820800,
	0x0001041000202000, 0x0000820800101000, 0x0000104400080800, 0x0000020080080080,
	0x0000404040040100, 0x0000808100020100, 0x0001010100020800, 0x0000808080010400,
	0x0000820820004000, 0x0000410410002000, 0x0000082088001000, 0x0000002011000800,
	0x0000080100400400, 0x0001010101000200, 0x0002020202000400, 0x0001010101000200,
	0x0000410410400000, 0x0000208208200000, 0x0000002084100000, 0x0000000020880000,
	0x0000001002020000, 0x0000040408020000, 0x0004040404040000, 0x0002020202020000,
	0x0000104104104000, 0x0000002082082000, 0x0000000020841000, 0x0000000000208800,
	0x0000000010020200, 0x0000000404080200, 0x0000040404040400, 0x0002020202020200,
}

var bishopMask = [...]uint64{
	0x0040201008040200, 0x0000402010080400, 0x0000004020100A00, 0x0000000040221400,
	0x0000000002442800, 0x0000000204085000, 0x0000020408102000, 0x0002040810204000,
	0x0020100804020000, 0x0040201008040000, 0x00004020100A0000, 0x0000004022140000,
	0x0000000244280000, 0x0000020408500000, 0x0002040810200000, 0x0004081020400000,
	0x0010080402000200, 0x0020100804000400, 0x004020100A000A00, 0x0000402214001400,
	0x0000024428002800, 0x0002040850005000, 0x0004081020002000, 0x0008102040004000,
	0x0008040200020400, 0x0010080400040800, 0x0020100A000A1000, 0x0040221400142200,
	0x0002442800284400, 0x0004085000500800, 0x0008102000201000, 0x0010204000402000,
	0x0004020002040800, 0x0008040004081000, 0x00100A000A102000, 0x0022140014224000,
	0x0044280028440200, 0x0008500050080400, 0x0010200020100800, 0x0020400040201000,
	0x0002000204081000, 0x0004000408102000, 0x000A000A10204000, 0x0014001422400000,
	0x0028002844020000, 0x0050005008040200, 0x0020002010080400, 0x0040004020100800,
	0x0000020408102000, 0x0000040810204000, 0x00000A1020400000, 0x0000142240000000,
	0x0000284402000000, 0x0000500804020000, 0x0000201008040200, 0x0000402010080400,
	0x0002040810204000, 0x0004081020400000, 0x000A102040000000, 0x0014224000000000,
	0x0028440200000000, 0x0050080402000000, 0x0020100804020000, 0x0040201008040200,
}

func magicify(b uint64, index int) uint64 {
	var bitmask uint64
	var count = bits.OnesCount64(b)

	for i, our := 0, b; i < count; i++ {
		their := ((our - 1) & our) ^ our
		our &= our - 1
		if (1<<uint(i))&index != 0 {
			bitmask |= their
		}
	}

	return bitmask
}

func computeSlideAttacks(f int, occ uint64, fs []func(sq uint64) uint64) uint64 {
	var result uint64
	for _, shift := range fs {
		var x = shift(squareMask[f])
		for x != 0 {
			result |= x
			if (x & occ) != 0 {
				break
			}
			x = shift(x)
		}
	}
	return result
}

var rookAttacksArray [64][1 << 12]uint64
var bishopAttacksArray [64][1 << 9]uint64
var squareMask = initSquareMask()

func SquareMask(sq uint64) uint64 {
	return squareMask[sq]
}

func initSquareMask() [64]uint64 {
	var sqm [64]uint64
	for sq := 0; sq < 64; sq++ {
		var b = uint64(1 << sq)
		sqm[sq] = b
	}
	return sqm
}
func init() {
	var rookShifts = [...]func(uint64) uint64{nortOne, westOne, soutOne, eastOne}
	var bishopShifts = [...]func(uint64) uint64{noWeOne, noEaOne, soWeOne, soEaOne}

	for sq := 0; sq < 64; sq++ {

		// Needs to retire, when we get rid of horizontal and vertical double rooks
		// That is hopefully, when the more efficient in-between mask is created
		rayAttacksArray[North][sq] = northRay(squareMask[sq])
		rayAttacksArray[NorthEast][sq] = northEastRay(squareMask[sq])
		rayAttacksArray[East][sq] = eastRay(squareMask[sq])
		rayAttacksArray[SouthEast][sq] = southEastRay(squareMask[sq])
		rayAttacksArray[South][sq] = southRay(squareMask[sq])
		rayAttacksArray[SouthWest][sq] = southWestRay(squareMask[sq])
		rayAttacksArray[West][sq] = westRay(squareMask[sq])
		rayAttacksArray[NorthWest][sq] = northWestRay(squareMask[sq])

		// Knights
		computedKnightAttacks[sq] = knightAttacks(squareMask[sq])

		// Kings
		computedKingAttacks[sq] = kingAttacks(squareMask[sq])

		// Rooks.
		var mask = rookMask[sq]
		var count = 1 << uint(bits.OnesCount64(mask))
		for i := 0; i < count; i++ {

			var occ = magicify(mask, i)
			var attacks = computeSlideAttacks(sq, occ, rookShifts[:])
			rookAttacksArray[sq][((rookMask[sq]&occ)*rookMult[sq])>>rookShift] = attacks
		}

		// Bishops.
		mask = bishopMask[sq]
		count = 1 << uint(bits.OnesCount64(mask))
		for i := 0; i < count; i++ {
			var occ = magicify(mask, i)
			var attacks = computeSlideAttacks(sq, occ, bishopShifts[:])
			bishopAttacksArray[sq][((bishopMask[sq]&occ)*bishopMult[sq])>>bishopShift] = attacks
		}
	}
}
