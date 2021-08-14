package engine

import (
	"math/bits"
)

func (p *Position) IsPseudoLegal(move Move) bool {
	turn := p.Turn()
	src := SquareMask[move.Source()]
	dest := SquareMask[move.Destination()]
	board := p.Board
	var ownPieces, enemyPieces, enemyKing, enemyPawns uint64
	var canQueenSideCastle, canKingSideCastle bool
	var queenSideRookSquare, kingSideRookSquare Square
	if turn == White {
		ownPieces = board.whitePieces
		enemyPieces = board.blackPieces
		enemyKing = board.blackKing
		enemyPawns = board.blackPawn
		canQueenSideCastle = p.HasTag(WhiteCanCastleQueenSide)
		canKingSideCastle = p.HasTag(WhiteCanCastleKingSide)
		queenSideRookSquare = A1
		kingSideRookSquare = H1
	} else {
		ownPieces = board.blackPieces
		enemyPieces = board.whitePieces
		enemyKing = board.whiteKing
		enemyPawns = board.whitePawn
		canQueenSideCastle = p.HasTag(BlackCanCastleQueenSide)
		canKingSideCastle = p.HasTag(BlackCanCastleKingSide)
		queenSideRookSquare = A8
		kingSideRookSquare = H8
	}
	allPieces := enemyPieces | ownPieces

	// we are moving own pieces, and do not capture own pieces
	if ownPieces&src == 0 || ownPieces&dest != 0 {
		return false
	}
	if move.IsEnPassant() {
		ep := findEnPassantCaptureSquare(move)
		if ep != p.EnPassant || enemyPawns&dest == 0 {
			return false
		}
	}

	// if we do capture enemy pieces, we should have Capture tag
	if enemyPieces&dest != 0 && !move.IsCapture() {
		return false
	}

	// the board is consistent with capturing piece
	if move.IsCapture() {
		cp := move.CapturedPiece()
		if cp == WhitePawn && board.whitePawn&dest == 0 {
			return false
		} else if cp == WhiteKnight && board.whiteKnight&dest == 0 {
			return false
		} else if cp == WhiteBishop && board.whiteBishop&dest == 0 {
			return false
		} else if cp == WhiteRook && board.whiteRook&dest == 0 {
			return false
		} else if cp == WhiteQueen && board.whiteQueen&dest == 0 {
			return false
		} else if cp == BlackPawn && board.blackPawn&dest == 0 {
			return false
		} else if cp == BlackKnight && board.blackKnight&dest == 0 {
			return false
		} else if cp == BlackBishop && board.blackBishop&dest == 0 {
			return false
		} else if cp == BlackRook && board.blackRook&dest == 0 {
			return false
		} else if cp == BlackQueen && board.blackQueen&dest == 0 {
			return false
		}
	}

	// We are not capturing enemy king
	if enemyKing&dest != 0 {
		return false
	}

	// Check that sliding pieces (2 push pawn, rook, bishop and queen, are not jumping over anything)
	movingPieceType := move.MovingPiece().Type()
	if movingPieceType == Queen || movingPieceType == Bishop || movingPieceType == Rook || movingPieceType == Pawn {
		squares := squaresInBetween[src][dest]
		if squares&allPieces != 0 {
			return false
		}
	}
	// Check that castling is correct (has castling right, the squares that matter are not in check)
	if move.IsQueenSideCastle() {
		if !canQueenSideCastle && p.IsInCheck() {
			return false
		}
		// check in between squares
		checkFreeSquares := squaresInBetween[src][dest]
		emptySquares := squaresInBetween[src][queenSideRookSquare]
		if emptySquares&allPieces != 0 {
			return false
		}
		taboo := tabooSquares(board, turn)
		if taboo&checkFreeSquares != 0 {
			return false
		}
	} else if move.IsKingSideCastle() {
		if !canKingSideCastle || p.IsInCheck() {
			return false
		}
		// check in between squares
		checkFreeSquares := squaresInBetween[src][dest]
		emptySquares := squaresInBetween[src][kingSideRookSquare]
		if emptySquares&allPieces != 0 {
			return false
		}
		taboo := tabooSquares(board, turn)
		if taboo&checkFreeSquares != 0 {
			return false
		}
	}

	return true
}

func (p *Position) PseudoLegalMoves() []Move {
	ml := NewMoveList(500)

	p.GetCaptureMoves(ml)
	p.GetQuietMoves(ml)

	return ml.Moves[:ml.Size]
}

func (p *Position) GetQuietMoves(ml *MoveList) {
	color := p.Turn()
	board := p.Board

	var taboo uint64

	if color == White &&
		(p.HasTag(WhiteCanCastleKingSide | WhiteCanCastleQueenSide)) {
		taboo = tabooSquares(board, color)
	} else if color == Black &&
		(p.HasTag(BlackCanCastleKingSide | BlackCanCastleQueenSide)) {
		taboo = tabooSquares(board, color)
	}

	p.pawnQuietMoves(color, ml)
	p.knightQuietMoves(color, ml)
	p.slidingQuietMoves(color, Bishop, ml)
	p.slidingQuietMoves(color, Rook, ml)
	p.slidingQuietMoves(color, Queen, ml)
	p.kingQuietMoves(taboo, color, ml)
}

func (p *Position) GetCaptureMoves(ml *MoveList) {
	color := p.Turn()

	p.pawnCaptureMoves(color, ml)
	p.knightCaptureMoves(color, ml)
	p.slidingCaptureMoves(color, Bishop, ml)
	p.slidingCaptureMoves(color, Rook, ml)
	p.slidingCaptureMoves(color, Queen, ml)
	p.kingCaptureMoves(color, ml)
}

// Checks and Pins
func isInCheck(b *Bitboard, colorOfKing Color) bool {
	return isKingAttacked(b, colorOfKing)
}

func isKingAttacked(b *Bitboard, colorOfKing Color) bool {
	var ownKing, opPawnAttacks, opKnights, opRQ, opBQ, opKing uint64
	var squareOfKing Square
	occupiedBB := b.whitePieces | b.blackPieces
	if colorOfKing == White {
		kingIndex := bitScanForward(b.whiteKing)
		ownKing = SquareMask[kingIndex]
		opPawnAttacks = wPawnsAble2CaptureAny(ownKing, b.blackPawn)
		if opPawnAttacks != 0 {
			return true
		}
		squareOfKing = Square(kingIndex)
		opKnights = b.blackKnight
		opRQ = b.blackRook | b.blackQueen
		opBQ = b.blackBishop | b.blackQueen
		opKing = b.blackKing
	} else {
		kingIndex := bitScanForward(b.blackKing)
		ownKing = SquareMask[kingIndex]
		opPawnAttacks = bPawnsAble2CaptureAny(ownKing, b.whitePawn)
		if opPawnAttacks != 0 {
			return true
		}
		squareOfKing = Square(kingIndex)
		opKnights = b.whiteKnight
		opRQ = b.whiteRook | b.whiteQueen
		opBQ = b.whiteBishop | b.whiteQueen
		opKing = b.whiteKing
	}

	knightChecks := (knightAttacks(ownKing) & opKnights)
	if knightChecks != 0 {
		return true
	}

	kingChecks := (kingAttacks(ownKing) & opKing)
	if kingChecks != 0 {
		return true
	}

	bishopChecks := (bishopAttacks(squareOfKing, occupiedBB, empty) & opBQ)

	if bishopChecks != 0 {
		return true
	}

	rookChecks := (rookAttacks(squareOfKing, occupiedBB, empty) & opRQ)
	if rookChecks != 0 {
		return true
	}

	if colorOfKing == White {
		return opPawnAttacks != 0
	}

	return false
}

func tabooSquares(b *Bitboard, colorOfKing Color) uint64 {
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
	taboo := opPawns | (knightAttacks(opKnights)) | kingAttacks(opKing)
	for opB != 0 {
		sq := bitScanForward(opB)
		taboo |= bishopAttacks(Square(sq), occupiedBB, opPieces)
		opB ^= SquareMask[sq]
	}

	for opR != 0 {
		sq := bitScanForward(opR)
		taboo |= rookAttacks(Square(sq), occupiedBB, opPieces)
		opR ^= SquareMask[sq]
	}

	for opQ != 0 {
		sq := bitScanForward(opQ)
		taboo |= queenAttacks(Square(sq), occupiedBB, opPieces)
		opQ ^= SquareMask[sq]
	}

	return taboo
}

// Quiet moves

func (p *Position) pawnQuietMoves(color Color,
	ml *MoveList) {
	var bbPawn, ownPieces, otherPieces uint64
	if color == White {
		bbPawn = p.Board.whitePawn
		ownPieces = p.Board.whitePieces
		otherPieces = p.Board.blackPieces
	} else {
		bbPawn = p.Board.blackPawn
		ownPieces = p.Board.blackPieces
		otherPieces = p.Board.whitePieces
	}
	emptySquares := (otherPieces | ownPieces) ^ universal
	if color == White {
		bbPawn &^= rank7
		for bbPawn != 0 {
			src := bitScanForward(bbPawn)
			srcSq := Square(src)
			pawn := SquareMask[src]
			if srcSq.Rank() == Rank2 {
				dbl := wDoublePushTargets(pawn, emptySquares)
				if dbl != 0 {
					dest := Square(bitScanForward(dbl))
					var tag MoveTag = 0
					m := NewMove(srcSq, dest, WhitePawn, NoPiece, NoType, tag)
					ml.Add(m)
				}
			}
			sngl := wSinglePushTargets(pawn, emptySquares)
			if sngl != 0 {
				dest := Square(bitScanForward(sngl))
				m := NewMove(srcSq, dest, WhitePawn, NoPiece, NoType, 0)
				ml.Add(m)
			}
			bbPawn ^= pawn
		}
	} else if color == Black {
		bbPawn &^= rank2
		for bbPawn != 0 {
			src := bitScanForward(bbPawn)
			srcSq := Square(src)
			pawn := SquareMask[src]
			dbl := bDoublePushTargets(pawn, emptySquares)
			if dbl != 0 {
				dest := Square(bitScanForward(dbl))
				var tag MoveTag = 0
				m := NewMove(srcSq, dest, BlackPawn, NoPiece, NoType, tag)
				ml.Add(m)
			}
			sngl := bSinglePushTargets(pawn, emptySquares)
			if sngl != 0 {
				dest := Square(bitScanForward(sngl))
				m := NewMove(srcSq, dest, BlackPawn, NoPiece, NoType, 0)
				ml.Add(m)
			}
			bbPawn ^= pawn
		}
	}
}

func (p *Position) knightQuietMoves(color Color, ml *MoveList) {
	var movingPiece Piece
	var bbPiece, ownPieces, otherPieces uint64
	if color == White {
		movingPiece = WhiteKnight
		bbPiece = p.Board.whiteKnight
		ownPieces = p.Board.whitePieces
		otherPieces = p.Board.blackPieces
	} else {
		movingPiece = BlackKnight
		bbPiece = p.Board.blackKnight
		ownPieces = p.Board.blackPieces
		otherPieces = p.Board.whitePieces
	}
	both := otherPieces | ownPieces
	for bbPiece != 0 {
		src := bitScanForward(bbPiece)
		srcSq := Square(src)
		knight := SquareMask[src]
		moves := knightMovesNoCaptures(srcSq, both)
		for moves != 0 {
			sq := bitScanForward(moves)
			dest := Square(sq)
			m := NewMove(srcSq, dest, movingPiece, NoPiece, NoType, 0)
			ml.Add(m)
			moves ^= SquareMask[sq]
		}
		bbPiece ^= knight
	}
}

func (p *Position) slidingQuietMoves(color Color, pieceType PieceType, ml *MoveList) {
	var ownPieces, otherPieces uint64
	movingPiece := GetPiece(pieceType, color)
	bbPiece := p.Board.GetBitboardOf(movingPiece)
	if color == White {
		ownPieces = p.Board.whitePieces
		otherPieces = p.Board.blackPieces
	} else {
		ownPieces = p.Board.blackPieces
		otherPieces = p.Board.whitePieces
	}
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
			ml.Add(m)
			passiveMoves ^= SquareMask[sq]
		}

		bbPiece ^= SquareMask[src]
	}
}

func (p *Position) kingQuietMoves(tabooSquares uint64, color Color, ml *MoveList) {
	var bbPiece, ownPieces, otherPieces uint64
	var kingSideCastle, queenSideCastle bool
	var movingPiece Piece
	if color == White {
		bbPiece = p.Board.whiteKing
		ownPieces = p.Board.whitePieces
		otherPieces = p.Board.blackPieces
		kingSideCastle = p.HasTag(WhiteCanCastleKingSide)
		queenSideCastle = p.HasTag(WhiteCanCastleQueenSide)
		movingPiece = WhiteKing
	} else {
		bbPiece = p.Board.blackKing
		ownPieces = p.Board.blackPieces
		otherPieces = p.Board.whitePieces
		kingSideCastle = p.HasTag(BlackCanCastleKingSide)
		queenSideCastle = p.HasTag(BlackCanCastleQueenSide)
		movingPiece = BlackKing
	}
	both := (otherPieces | ownPieces)
	if bbPiece != 0 {
		src := bitScanForward(bbPiece)
		srcSq := Square(src)
		moves := kingMovesNoCaptures(srcSq, both, tabooSquares)
		for moves != 0 {
			sq := bitScanForward(moves)
			dest := Square(sq)
			m := NewMove(srcSq, dest, movingPiece, NoPiece, NoType, 0)
			ml.Add(m)

			moves ^= SquareMask[sq]
		}

		kingSide := whiteKingSideCastle
		queenSide := whiteQueenSideCastle
		kingCastleMove := whiteKingCastleMove
		queenCastleMove := whiteQueenCastleMove
		E := E1
		B := B1
		if color == Black {
			kingCastleMove = blackKingCastleMove
			queenCastleMove = blackQueenCastleMove
			kingSide = blackKingSideCastle
			queenSide = blackQueenSideCastle
			E = E8
			B = B8
		}

		if kingSideCastle &&
			((ownPieces|otherPieces)&kingSide == 0) && // are empty
			(tabooSquares&(kingSide|SquareMask[E]) == 0) { // Not in check
			m := kingCastleMove
			ml.Add(m)
		}

		if srcSq == E && queenSideCastle &&
			((ownPieces|otherPieces)&(queenSide|(SquareMask[B])) == 0) && // are empty
			(tabooSquares&(queenSide|SquareMask[E]) == 0) { // Not in check
			m := queenCastleMove
			ml.Add(m)

		}
	}
}

// Capture moves

func (p *Position) pawnCaptureMoves(color Color,
	ml *MoveList) {
	var bbPawn, ownPieces, otherPieces, otherKing uint64
	if color == White {
		bbPawn = p.Board.whitePawn
		ownPieces = p.Board.whitePieces
		otherKing = p.Board.blackKing
		otherPieces = p.Board.blackPieces ^ otherKing
	} else {
		bbPawn = p.Board.blackPawn
		ownPieces = p.Board.blackPieces
		otherKing = p.Board.whiteKing
		otherPieces = p.Board.whitePieces ^ otherKing
	}
	emptySquares := (otherPieces | otherKing | ownPieces) ^ universal
	enPassant := p.EnPassant
	if color == White {
		for bbPawn != 0 {
			src := bitScanForward(bbPawn)
			srcSq := Square(src)
			pawn := SquareMask[src]
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
					ml.AddFour(m1, m2, m3, m4)
				} else {
					m := NewMove(srcSq, dest, WhitePawn, cp, NoType, Capture)
					ml.Add(m)
				}
				attacks ^= SquareMask[sq]
			}
			if srcSq.Rank() == Rank5 && enPassant != NoSquare && enPassant.Rank() == Rank6 {
				ep := SquareMask[enPassant]
				r := wPawnsAble2CaptureAny(pawn, ep)
				if r != 0 {
					dest := Square(bitScanForward(r))
					var tag MoveTag = Capture | EnPassant
					m := NewMove(srcSq, dest, WhitePawn, BlackPawn, NoType, tag)
					ml.Add(m)
				}
			}
			if srcSq.Rank() == Rank7 {
				sngl := wSinglePushTargets(pawn, emptySquares)
				if sngl != 0 {
					sq := bitScanForward(sngl)
					dest := Square(sq)
					m1 := NewMove(srcSq, dest, WhitePawn, NoPiece, Queen, 0)
					m2 := NewMove(srcSq, dest, WhitePawn, NoPiece, Rook, 0)
					m3 := NewMove(srcSq, dest, WhitePawn, NoPiece, Bishop, 0)
					m4 := NewMove(srcSq, dest, WhitePawn, NoPiece, Knight, 0)
					ml.AddFour(m1, m2, m3, m4)
				}
			}
			bbPawn ^= pawn
		}
	} else if color == Black {
		for bbPawn != 0 {
			src := bitScanForward(bbPawn)
			srcSq := Square(src)
			pawn := SquareMask[src]
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
					ml.AddFour(m1, m2, m3, m4)
				} else {
					var tag MoveTag = Capture
					m := NewMove(srcSq, dest, BlackPawn, cp, NoType, tag)
					ml.Add(m)
				}
				attacks ^= SquareMask[sq]
			}
			if srcSq.Rank() == Rank4 && enPassant != NoSquare && enPassant.Rank() == Rank3 {
				ep := SquareMask[enPassant]
				r := bPawnsAble2CaptureAny(pawn, ep)
				if r != 0 {
					dest := Square(bitScanForward(r))
					var tag MoveTag = Capture | EnPassant
					m := NewMove(srcSq, dest, BlackPawn, WhitePawn, NoType, tag)
					ml.Add(m)
				}
			}
			if srcSq.Rank() == Rank2 {
				sngl := bSinglePushTargets(pawn, emptySquares)
				if sngl != 0 {
					dest := Square(bitScanForward(sngl))
					m1 := NewMove(srcSq, dest, BlackPawn, NoPiece, Queen, 0)
					m2 := NewMove(srcSq, dest, BlackPawn, NoPiece, Rook, 0)
					m3 := NewMove(srcSq, dest, BlackPawn, NoPiece, Bishop, 0)
					m4 := NewMove(srcSq, dest, BlackPawn, NoPiece, Knight, 0)
					ml.AddFour(m1, m2, m3, m4)
				}
			}
			bbPawn ^= pawn
		}
	}
}

func (p *Position) knightCaptureMoves(color Color, ml *MoveList) {
	var movingPiece Piece
	var bbPiece, otherPieces uint64
	if color == White {
		movingPiece = WhiteKnight
		bbPiece = p.Board.whiteKnight
		otherPieces = p.Board.blackPieces ^ p.Board.blackKing
	} else {
		movingPiece = BlackKnight
		bbPiece = p.Board.blackKnight
		otherPieces = p.Board.whitePieces ^ p.Board.whiteKing
	}
	for bbPiece != 0 {
		src := bitScanForward(bbPiece)
		srcSq := Square(src)
		knight := SquareMask[src]
		captures := knightCaptures(srcSq, otherPieces)
		for captures != 0 {
			sq := bitScanForward(captures)
			dest := Square(sq)
			cp := p.Board.PieceAt(dest)
			m := NewMove(srcSq, dest, movingPiece, cp, NoType, Capture)
			ml.Add(m)

			captures ^= SquareMask[sq]
		}
		bbPiece ^= knight
	}
}

func (p *Position) slidingCaptureMoves(color Color, pieceType PieceType, ml *MoveList) {
	var ownPieces, otherPieces uint64
	movingPiece := GetPiece(pieceType, color)
	bbPiece := p.Board.GetBitboardOf(movingPiece)
	if color == White {
		ownPieces = p.Board.whitePieces
		otherPieces = p.Board.blackPieces ^ p.Board.blackKing
	} else {
		ownPieces = p.Board.blackPieces
		otherPieces = p.Board.whitePieces ^ p.Board.whiteKing
	}
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
			ml.Add(m)

			captureMoves ^= SquareMask[sq]
		}
		bbPiece ^= SquareMask[src]
	}
}

func (p *Position) kingCaptureMoves(color Color, ml *MoveList) {
	var bbPiece, otherPieces uint64
	var movingPiece Piece
	if color == White {
		bbPiece = p.Board.whiteKing
		otherPieces = p.Board.blackPieces ^ p.Board.blackKing
		movingPiece = WhiteKing
	} else {
		bbPiece = p.Board.blackKing
		otherPieces = p.Board.whitePieces ^ p.Board.whiteKing
		movingPiece = BlackKing
	}
	if bbPiece != 0 {
		src := bitScanForward(bbPiece)
		srcSq := Square(src)
		captures := kingCaptures(srcSq, otherPieces)
		for captures != 0 {
			sq := bitScanForward(captures)
			dest := Square(sq)
			cp := p.Board.PieceAt(dest)
			m := NewMove(srcSq, dest, movingPiece, cp, NoType, Capture)
			ml.Add(m)
			captures ^= SquareMask[sq]
		}
	}
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

func kingCaptures(sq Square, others uint64) uint64 {
	attacks := computedKingAttacks[sq]
	return (attacks & others)
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
const blackQueenSideCastle = uint64(1<<D8 | 1<<C8)
const whiteQueenSideCastle = uint64(1<<D1 | 1<<C1)
const blackKingSideCastle = uint64(1<<F8 | 1<<G8)
const whiteKingSideCastle = uint64(1<<F1 | 1<<G1)

var whiteKingCastleMove = NewMove(E1, G1, WhiteKing, NoPiece, NoType, KingSideCastle)
var whiteQueenCastleMove = NewMove(E1, C1, WhiteKing, NoPiece, NoType, QueenSideCastle)
var blackKingCastleMove = NewMove(E8, G8, BlackKing, NoPiece, NoType, KingSideCastle)
var blackQueenCastleMove = NewMove(E8, C8, BlackKing, NoPiece, NoType, QueenSideCastle)

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
		var x = shift(SquareMask[f])
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

var squaresInBetween [64][64]uint64
var rookAttacksArray [64][1 << 12]uint64
var bishopAttacksArray [64][1 << 9]uint64
var SquareMask = initSquareMask()

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
		rayAttacksArray[North][sq] = northRay(SquareMask[sq])
		rayAttacksArray[NorthEast][sq] = northEastRay(SquareMask[sq])
		rayAttacksArray[East][sq] = eastRay(SquareMask[sq])
		rayAttacksArray[SouthEast][sq] = southEastRay(SquareMask[sq])
		rayAttacksArray[South][sq] = southRay(SquareMask[sq])
		rayAttacksArray[SouthWest][sq] = southWestRay(SquareMask[sq])
		rayAttacksArray[West][sq] = westRay(SquareMask[sq])
		rayAttacksArray[NorthWest][sq] = northWestRay(SquareMask[sq])

		// Knights
		computedKnightAttacks[sq] = knightAttacks(SquareMask[sq])

		// Kings
		computedKingAttacks[sq] = kingAttacks(SquareMask[sq])

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

		for s1 := 0; s1 < 64; s1++ {
			for s2 := 0; s2 < 64; s2++ {
				squaresInBetween[s1][s2] = 0
				if (queenAttacks(Square(s1), 0, 0) & SquareMask[s2]) != 0 {
					var delta = ((s2 - s1) / squareDistance(s1, s2))
					for s := s1 + delta; s != s2; s += delta {
						squaresInBetween[s1][s2] |= SquareMask[s]
					}
				}
			}
		}
	}
}

func fileDistance(sq1 int, sq2 int) int {
	diff := int(Square(sq1).File() - Square(sq2).File())
	if diff < 0 {
		return -diff
	}
	return diff
}

func rankDistance(sq1 int, sq2 int) int {
	diff := int(Square(sq1).Rank() - Square(sq2).Rank())
	if diff < 0 {
		return -diff
	}
	return diff
}

func squareDistance(sq1, sq2 int) int {
	fileDist := fileDistance(sq1, sq2)
	rankDist := rankDistance(sq1, sq2)
	if fileDist > rankDist {
		return fileDist
	}
	return rankDist
}
