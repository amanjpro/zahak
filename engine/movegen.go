package engine

import (
	"math/bits"
)

func (p *Position) addAllMoves(allMoves *[]Move, ms ...Move) {
	color := p.Turn()
	for _, m := range ms {
		// make the move
		p.partialMakeMove(m)

		// Does the move puts the moving player in check
		pNotInCheck := !isInCheck(p.Board, color)

		if pNotInCheck {
			if isInCheck(p.Board, p.Turn()) { // We put opponent in check
				m.AddCheckTag()
				*allMoves = append(*allMoves, m)
			} else { // The move does not put us in check
				// do nothing
				*allMoves = append(*allMoves, m)
			}
		}
		p.partialUnMakeMove(m)
	}
}

func (p *Position) addCaptureMoves(allMoves *[]Move, withChecks bool, isChecked bool, ms ...Move) {
	color := p.Turn()
	for _, m := range ms {
		// make the move
		p.partialMakeMove(m)

		// Does the move puts the moving player in check
		pNotInCheck := !isInCheck(p.Board, color)

		if pNotInCheck {
			isCheckMove := isInCheck(p.Board, p.Turn())

			if isCheckMove {
				m.AddCheckTag()
			}

			if withChecks && isCheckMove { // We put opponent in check
				*allMoves = append(*allMoves, m)
			} else if m.IsCapture() { // The move is a capture
				*allMoves = append(*allMoves, m)
			} else if isChecked { // Check replies are also considered
				*allMoves = append(*allMoves, m)
			}
		}
		p.partialUnMakeMove(m)
	}
}

func (p *Position) LegalMoves() []Move {
	allMoves := make([]Move, 0, 256)

	p.generateMoves(&allMoves, false, p.IsInCheck(), false)

	return allMoves
}

func (p *Position) QuiesceneMoves(withChecks bool) []Move {
	allMoves := make([]Move, 0, 256)

	isChecked := p.IsInCheck()

	p.generateMoves(&allMoves, !(withChecks || isChecked), isChecked, true)

	return allMoves
}

func (p *Position) generateMoves(allMoves *[]Move, capturesOnly bool, positionIsInCheck bool, isQuiescence bool) {

	color := p.Turn()
	board := p.Board

	taboo := tabooSquares(board, color)

	// If it is double check, only king can move
	if positionIsInCheck && isDoubleCheck(board, color) {
		if color == White {
			p.bbKingMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
				taboo, color, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide),
				capturesOnly, positionIsInCheck, false, isQuiescence, allMoves)
		} else if color == Black {
			p.bbKingMoves(board.blackKing, board.blackPieces, board.whitePieces, board.whiteKing,
				taboo, color, p.HasTag(BlackCanCastleKingSide), p.HasTag(BlackCanCastleQueenSide),
				capturesOnly, positionIsInCheck, false, isQuiescence, allMoves)
		}
	} else {

		if color == White {
			p.bbPawnMoves(board.whitePawn, board.whitePieces, board.blackPieces,
				color, p.EnPassant, capturesOnly, positionIsInCheck, false, isQuiescence, allMoves)
			p.bbKnightMoves(WhiteKnight, board.whiteKnight, board.whitePieces, board.blackPieces,
				capturesOnly, positionIsInCheck, false, isQuiescence, allMoves)
			p.bbSlidingMoves(board.whiteBishop, board.whitePieces, board.blackPieces,
				color, WhiteBishop, capturesOnly, positionIsInCheck, false, isQuiescence, allMoves)
			p.bbSlidingMoves(board.whiteRook, board.whitePieces, board.blackPieces,
				color, WhiteRook, capturesOnly, positionIsInCheck, false, isQuiescence, allMoves)
			p.bbSlidingMoves(board.whiteQueen, board.whitePieces, board.blackPieces,
				color, WhiteQueen, capturesOnly, positionIsInCheck, false, isQuiescence, allMoves)
			p.bbKingMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
				taboo, color, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide),
				capturesOnly, positionIsInCheck, false, isQuiescence, allMoves)
		} else if color == Black {
			p.bbPawnMoves(board.blackPawn, board.blackPieces, board.whitePieces,
				color, p.EnPassant, capturesOnly, positionIsInCheck, false, isQuiescence, allMoves)
			p.bbKnightMoves(BlackKnight, board.blackKnight, board.blackPieces, board.whitePieces,
				capturesOnly, positionIsInCheck, false, isQuiescence, allMoves)
			p.bbSlidingMoves(board.blackBishop, board.blackPieces, board.whitePieces,
				color, BlackBishop, capturesOnly, positionIsInCheck, false, isQuiescence, allMoves)
			p.bbSlidingMoves(board.blackRook, board.blackPieces, board.whitePieces,
				color, BlackRook, capturesOnly, positionIsInCheck, false, isQuiescence, allMoves)
			p.bbSlidingMoves(board.blackQueen, board.blackPieces, board.whitePieces,
				color, BlackQueen, capturesOnly, positionIsInCheck, false, isQuiescence, allMoves)
			p.bbKingMoves(board.blackKing, board.blackPieces, board.whitePieces, board.whiteKing,
				taboo, color, p.HasTag(BlackCanCastleKingSide), p.HasTag(BlackCanCastleQueenSide),
				capturesOnly, positionIsInCheck, false, isQuiescence, allMoves)
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

	// If it is double check, only king can move
	if isDoubleCheck(board, color) {
		if color == White {
			return p.bbKingMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
				taboo, color, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide), false, false, true, false, nil)
		} else if color == Black {
			return p.bbKingMoves(board.blackKing, board.blackPieces, board.whitePieces, board.whiteKing,
				taboo, color, p.HasTag(BlackCanCastleKingSide), p.HasTag(BlackCanCastleQueenSide), false, false, true, false, nil)
		}
		return false
	} else {

		if color == White {
			return p.bbPawnMoves(board.whitePawn, board.whitePieces, board.blackPieces,
				color, p.EnPassant, false, false, true, false, nil) ||
				p.bbKnightMoves(WhiteKnight, board.whiteKnight, board.whitePieces, board.blackPieces,
					false, false, true, false, nil) ||
				p.bbSlidingMoves(board.whiteBishop, board.whitePieces, board.blackPieces,
					color, WhiteBishop, false, false, true, false, nil) ||
				p.bbSlidingMoves(board.whiteRook, board.whitePieces, board.blackPieces,
					color, WhiteRook, false, false, true, false, nil) ||
				p.bbSlidingMoves(board.whiteQueen, board.whitePieces, board.blackPieces,
					color, WhiteQueen, false, false, true, false, nil) ||
				p.bbKingMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
					taboo, color, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide), false, false, true, false, nil)
		} else if color == Black {
			return p.bbPawnMoves(board.blackPawn, board.blackPieces, board.whitePieces,
				color, p.EnPassant, false, false, true, false, nil) ||
				p.bbKnightMoves(BlackKnight, board.blackKnight, board.blackPieces, board.whitePieces,
					false, false, true, false, nil) ||
				p.bbSlidingMoves(board.blackBishop, board.blackPieces, board.whitePieces,
					color, BlackBishop, false, false, true, false, nil) ||
				p.bbSlidingMoves(board.blackRook, board.blackPieces, board.whitePieces,
					color, BlackRook, false, false, true, false, nil) ||
				p.bbSlidingMoves(board.blackQueen, board.blackPieces, board.whitePieces,
					color, BlackQueen, false, false, true, false, nil) ||
				p.bbKingMoves(board.blackKing, board.blackPieces, board.whitePieces, board.whiteKing,
					taboo, color, p.HasTag(BlackCanCastleKingSide), p.HasTag(BlackCanCastleQueenSide), false, false, true, false, nil)
		}
	}
	return false
}

// Checks and Pins
func isInCheck(b Bitboard, colorOfKing Color) bool {
	return isKingAttacked(b, colorOfKing, false)
}

func isKingAttacked(b Bitboard, colorOfKing Color, doubleCheck bool) bool {
	var ownKing, opPawnAttacks, opKnights, opRQ, opBQ uint64
	var squareOfKing Square
	occupiedBB := b.whitePieces | b.blackPieces
	if colorOfKing == White {
		kingIndex := bitScanForward(b.whiteKing)
		ownKing = (1 << kingIndex)
		squareOfKing = Square(kingIndex)
		opPawnAttacks = wPawnsAble2CaptureAny(ownKing, b.blackPawn)
		opKnights = b.blackKnight
		opRQ = b.blackRook | b.blackQueen
		opBQ = b.blackBishop | b.blackQueen
	} else {
		kingIndex := bitScanForward(b.blackKing)
		ownKing = (1 << kingIndex)
		squareOfKing = Square(kingIndex)
		opPawnAttacks = bPawnsAble2CaptureAny(ownKing, b.whitePawn)
		opKnights = b.whiteKnight
		opRQ = b.whiteRook | b.whiteQueen
		opBQ = b.whiteBishop | b.whiteQueen
	}
	acc := 0
	pawnChecks := opPawnAttacks
	if pawnChecks != 0 {
		if !doubleCheck {
			return true
		} else {
			acc++
		}
	}

	knightChecks := (knightAttacks(ownKing) & opKnights)
	if knightChecks != 0 {
		if !doubleCheck || acc == 1 {
			return true
		} else {
			acc += 1
		}
	}
	// Knights and pawns cannot discover each other
	bishopChecks := (bishopAttacks(squareOfKing, occupiedBB, empty) & opBQ)

	if bishopChecks != 0 && !doubleCheck {
		return true
	} else {
		for bishopChecks != 0 {
			sq := bitScanForward(bishopChecks)
			acc += 1
			if (!doubleCheck && acc >= 1) || acc > 1 {
				return true
			}
			bishopChecks ^= (1 << sq)
		}
	}

	rookChecks := (rookAttacks(squareOfKing, occupiedBB, empty) & opRQ)

	if rookChecks != 0 && !doubleCheck {
		return true
	} else {
		for rookChecks != 0 {
			sq := bitScanForward(rookChecks)
			acc += 1
			if (!doubleCheck && acc >= 1) || acc > 1 {
				return true
			}
			rookChecks ^= (1 << sq)
		}
	}

	return false
}

func isDoubleCheck(b Bitboard, colorOfKing Color) bool {
	return isKingAttacked(b, colorOfKing, true)
}

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
	taboo := opPawns | (knightAttacks(opKnights)) | kingAttacks(opKing)
	for opB != 0 {
		sq := bitScanForward(opB)
		taboo |= bishopAttacks(Square(sq), occupiedBB, opPieces)
		opB ^= (1 << sq)
	}

	for opR != 0 {
		sq := bitScanForward(opR)
		taboo |= rookAttacks(Square(sq), occupiedBB, opPieces)
		opR ^= (1 << sq)
	}

	for opQ != 0 {
		sq := bitScanForward(opQ)
		taboo |= queenAttacks(Square(sq), occupiedBB, opPieces)
		opQ ^= (1 << sq)
	}

	return taboo
}

// Pawns

func (p *Position) bbPawnMoves(bbPawn uint64, ownPieces uint64, otherPieces uint64, color Color, enPassant Square,
	capturesOnly bool, isPosInCheck bool, isLegalityCheck bool, isQuiescence bool, allMoves *[]Move) bool {
	emptySquares := (otherPieces | ownPieces) ^ universal
	if color == White {
		for bbPawn != 0 {
			src := bitScanForward(bbPawn)
			srcSq := Square(src)
			pawn := uint64(1 << src)
			if !capturesOnly {
				if srcSq.Rank() == Rank2 {
					dbl := wDoublePushTargets(pawn, emptySquares)
					if dbl != 0 {
						dest := Square(bitScanForward(dbl))
						var tag MoveTag = 0
						m := NewMove(srcSq, dest, WhitePawn, NoPiece, NoType, tag)
						if isLegalityCheck && p.checkMove(m) {
							return true
						} else if !isLegalityCheck {
							if isQuiescence {
								p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m)
							} else {
								p.addAllMoves(allMoves, m)
							}
						}
					}
				}
				sngl := wSinglePushTargets(pawn, emptySquares)
				if sngl != 0 {
					dest := Square(bitScanForward(sngl))
					if dest.Rank() == Rank8 {
						m1 := NewMove(srcSq, dest, WhitePawn, NoPiece, Queen, 0)
						m2 := NewMove(srcSq, dest, WhitePawn, NoPiece, Rook, 0)
						m3 := NewMove(srcSq, dest, WhitePawn, NoPiece, Bishop, 0)
						m4 := NewMove(srcSq, dest, WhitePawn, NoPiece, Knight, 0)
						if isLegalityCheck && p.checkMove(m1) { // if one is illegal, they all are illegal
							return true
						} else if !isLegalityCheck {
							if isQuiescence {
								p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m1, m2, m3, m4)
							} else {
								p.addAllMoves(allMoves, m1, m2, m3, m4)
							}
						}
					} else {
						m := NewMove(srcSq, dest, WhitePawn, NoPiece, NoType, 0)
						if isLegalityCheck && p.checkMove(m) {
							return true
						} else if !isLegalityCheck {
							if isQuiescence {
								p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m)
							} else {
								p.addAllMoves(allMoves, m)
							}
						}
					}
				}
			}
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
						if isQuiescence {
							p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m1, m2, m3, m4)
						} else {
							p.addAllMoves(allMoves, m1, m2, m3, m4)
						}
					}
				} else {
					m := NewMove(srcSq, dest, WhitePawn, cp, NoType, Capture)
					if isLegalityCheck && p.checkMove(m) {
						return true
					} else if !isLegalityCheck {
						if isQuiescence {
							p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m)
						} else {
							p.addAllMoves(allMoves, m)
						}
					}
				}
				attacks ^= (1 << sq)
			}
			if srcSq.Rank() == Rank5 && enPassant != NoSquare && enPassant.Rank() == Rank6 {
				ep := uint64(1 << enPassant)
				r := wPawnsAble2CaptureAny(pawn, ep)
				if r != 0 {
					dest := Square(bitScanForward(r))
					var tag MoveTag = Capture | EnPassant
					m := NewMove(srcSq, dest, WhitePawn, BlackPawn, NoType, tag)
					if isLegalityCheck && p.checkMove(m) {
						return true
					} else if !isLegalityCheck {
						if isQuiescence {
							p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m)
						} else {
							p.addAllMoves(allMoves, m)
						}
					}
				}
			}
			bbPawn ^= pawn
		}
	} else if color == Black {
		for bbPawn != 0 {
			src := bitScanForward(bbPawn)
			srcSq := Square(src)
			pawn := uint64(1 << src)
			if !capturesOnly {
				dbl := bDoublePushTargets(pawn, emptySquares)
				if dbl != 0 {
					dest := Square(bitScanForward(dbl))
					var tag MoveTag = 0
					m := NewMove(srcSq, dest, BlackPawn, NoPiece, NoType, tag)
					if isLegalityCheck && p.checkMove(m) {
						return true
					} else if !isLegalityCheck {
						if isQuiescence {
							p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m)
						} else {
							p.addAllMoves(allMoves, m)
						}
					}
				}
				sngl := bSinglePushTargets(pawn, emptySquares)
				if sngl != 0 {
					dest := Square(bitScanForward(sngl))
					if dest.Rank() == Rank1 {
						m1 := NewMove(srcSq, dest, BlackPawn, NoPiece, Queen, 0)
						m2 := NewMove(srcSq, dest, BlackPawn, NoPiece, Rook, 0)
						m3 := NewMove(srcSq, dest, BlackPawn, NoPiece, Bishop, 0)
						m4 := NewMove(srcSq, dest, BlackPawn, NoPiece, Knight, 0)
						if isLegalityCheck && p.checkMove(m1) { // if one is illegal, they all are illegal
							return true
						} else if !isLegalityCheck {
							if isQuiescence {
								p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m1, m2, m3, m4)
							} else {
								p.addAllMoves(allMoves, m1, m2, m3, m4)
							}
						}
					} else {
						m := NewMove(srcSq, dest, BlackPawn, NoPiece, NoType, 0)
						if isLegalityCheck && p.checkMove(m) {
							return true
						} else if !isLegalityCheck {
							if isQuiescence {
								p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m)
							} else {
								p.addAllMoves(allMoves, m)
							}
						}
					}
				}
			}
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
						if isQuiescence {
							p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m1, m2, m3, m4)
						} else {
							p.addAllMoves(allMoves, m1, m2, m3, m4)
						}
					}
				} else {
					var tag MoveTag = Capture
					m := NewMove(srcSq, dest, BlackPawn, cp, NoType, tag)
					if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
						return true
					} else if !isLegalityCheck {
						if isQuiescence {
							p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m)
						} else {
							p.addAllMoves(allMoves, m)
						}
					}
				}
				attacks ^= (1 << sq)
			}
			if srcSq.Rank() == Rank4 && enPassant != NoSquare && enPassant.Rank() == Rank3 {
				ep := uint64(1 << enPassant)
				r := bPawnsAble2CaptureAny(pawn, ep)
				if r != 0 {
					dest := Square(bitScanForward(r))
					var tag MoveTag = Capture | EnPassant
					m := NewMove(srcSq, dest, BlackPawn, WhitePawn, NoType, tag)
					if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
						return true
					} else if !isLegalityCheck {
						if isQuiescence {
							p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m)
						} else {
							p.addAllMoves(allMoves, m)
						}
					}
				}
			}
			bbPawn ^= pawn
		}
	}

	return false
}

// Sliding moves, for rooks, queens and bishops
func (p *Position) bbSlidingMoves(bbPiece uint64, ownPieces uint64, otherPieces uint64,
	color Color, movingPiece Piece,
	capturesOnly bool, isPosInCheck bool, isLegalityCheck bool, isQuiescence bool, allMoves *[]Move) bool {
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
		if !capturesOnly {
			passiveMoves := rayAttacks &^ otherPieces
			for passiveMoves != 0 {
				sq := bitScanForward(passiveMoves)
				dest := Square(sq)
				m := NewMove(srcSq, dest, movingPiece, NoPiece, NoType, 0)
				if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
					return true
				} else if !isLegalityCheck {
					if isQuiescence {
						p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m)
					} else {
						p.addAllMoves(allMoves, m)
					}
				}
				passiveMoves ^= (1 << sq)
			}
		}
		for captureMoves != 0 {
			sq := bitScanForward(captureMoves)
			dest := Square(sq)
			cp := p.Board.PieceAt(dest)
			m := NewMove(srcSq, dest, movingPiece, cp, NoType, Capture)
			if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
				return true
			} else if !isLegalityCheck {
				if isQuiescence {
					p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m)
				} else {
					p.addAllMoves(allMoves, m)
				}
			}

			captureMoves ^= (1 << sq)
		}
		bbPiece ^= (1 << src)
	}
	return false
}

// Knights
func (p *Position) bbKnightMoves(movingPiece Piece, bbPiece uint64, ownPieces uint64, otherPieces uint64,
	capturesOnly bool, isPosInCheck bool, isLegalityCheck bool, isQuiescence bool, allMoves *[]Move) bool {
	both := otherPieces | ownPieces
	for bbPiece != 0 {
		src := bitScanForward(bbPiece)
		srcSq := Square(src)
		knight := uint64(1 << src)
		if !capturesOnly {
			moves := knightMovesNoCaptures(srcSq, both)
			for moves != 0 {
				sq := bitScanForward(moves)
				dest := Square(sq)
				m := NewMove(srcSq, dest, movingPiece, NoPiece, NoType, 0)
				if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
					return true
				} else if !isLegalityCheck {
					if isQuiescence {
						p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m)
					} else {
						p.addAllMoves(allMoves, m)
					}
				}

				moves ^= (1 << sq)
			}
		}
		captures := knightCaptures(srcSq, otherPieces)
		for captures != 0 {
			sq := bitScanForward(captures)
			dest := Square(sq)
			cp := p.Board.PieceAt(dest)
			m := NewMove(srcSq, dest, movingPiece, cp, NoType, Capture)
			if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
				return true
			} else if !isLegalityCheck {
				if isQuiescence {
					p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m)
				} else {
					p.addAllMoves(allMoves, m)
				}
			}

			captures ^= (1 << sq)
		}
		bbPiece ^= knight
	}

	return false
}

// Kings
func (p *Position) bbKingMoves(bbPiece uint64, ownPieces uint64, otherPieces uint64, otherKing uint64,
	tabooSquares uint64, color Color, kingSideCastle bool, queenSideCastle bool,
	capturesOnly bool, isPosInCheck bool, isLegalityCheck bool, isQuiescence bool, allMoves *[]Move) bool {
	both := (otherPieces | ownPieces)
	var movingPiece = BlackKing
	if color == White {
		movingPiece = WhiteKing
	}
	if bbPiece != 0 {
		src := bitScanForward(bbPiece)
		srcSq := Square(src)
		moves := kingMovesNoCaptures(srcSq, both, tabooSquares)
		if !capturesOnly {
			for moves != 0 {
				sq := bitScanForward(moves)
				dest := Square(sq)
				m := NewMove(srcSq, dest, movingPiece, NoPiece, NoType, 0)
				if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
					return true
				} else if !isLegalityCheck {
					if isQuiescence {
						p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m)
					} else {
						p.addAllMoves(allMoves, m)
					}
				}

				moves ^= (1 << sq)
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

			kingSide := uint64(1<<F | 1<<G)
			queenSide := uint64(1<<D | 1<<C)

			if srcSq == E && kingSideCastle &&
				((ownPieces|otherPieces)&kingSide == 0) && // are empty
				(tabooSquares&(kingSide|1<<E) == 0) { // Not in check
				m := NewMove(srcSq, G, movingPiece, NoPiece, NoType, KingSideCastle)
				if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
					return true
				} else if !isLegalityCheck {
					if isQuiescence {
						p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m)
					} else {
						p.addAllMoves(allMoves, m)
					}
				}
			}

			if srcSq == E && queenSideCastle &&
				((ownPieces|otherPieces)&(queenSide|(1<<B)) == 0) && // are empty
				(tabooSquares&(queenSide|1<<E) == 0) { // Not in check
				m := NewMove(srcSq, C, movingPiece, NoPiece, NoType, QueenSideCastle)
				if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
					return true
				} else if !isLegalityCheck {
					if isQuiescence {
						p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m)
					} else {
						p.addAllMoves(allMoves, m)
					}
				}

			}
		}
		captures := kingCaptures(srcSq, otherPieces, tabooSquares)
		for captures != 0 {
			sq := bitScanForward(captures)
			dest := Square(sq)
			cp := p.Board.PieceAt(dest)
			m := NewMove(srcSq, dest, movingPiece, cp, NoType, Capture)
			if isLegalityCheck && p.checkMove(m) { // if one is illegal, they all are illegal
				return true
			} else if !isLegalityCheck {
				if isQuiescence {
					p.addCaptureMoves(allMoves, !capturesOnly, isPosInCheck, m)
				} else {
					p.addAllMoves(allMoves, m)
				}
			}

			captures ^= (1 << sq)
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

func rookAttacks(sq Square, occ uint64, ownPieces uint64) uint64 {
	allAttacks := getPositiveRayAttacks(sq, occ, North) |
		getPositiveRayAttacks(sq, occ, East) |
		getNegativeRayAttacks(sq, occ, South) |
		getNegativeRayAttacks(sq, occ, West)

	return allAttacks &^ ownPieces

}

func bishopAttacks(sq Square, occ uint64, ownPieces uint64) uint64 {
	allAttacks := getPositiveRayAttacks(sq, occ, NorthEast) |
		getPositiveRayAttacks(sq, occ, NorthWest) |
		getNegativeRayAttacks(sq, occ, SouthEast) |
		getNegativeRayAttacks(sq, occ, SouthWest)

	return allAttacks &^ ownPieces
}

func queenAttacks(sq Square, occ uint64, ownPieces uint64) uint64 {
	return rookAttacks(sq, occ, ownPieces) | bishopAttacks(sq, occ, ownPieces)
}

func slidingCheckTag(from Square, occ uint64, ownPieces uint64, otherKing uint64,
	attacks func(sq Square, occ uint64, own uint64) uint64) MoveTag {
	if attacks(from, occ, ownPieces)&otherKing != 0 {
		return Check
	}
	return 0
}

// The mighty knight

func knightCheckTag(from Square, otherKing uint64) MoveTag {
	if knightCaptures(from, otherKing) != 0 {
		return Check
	}
	return 0
}

var computedKnightAttacks = initializeKnightAttacks()

func initializeKnightAttacks() [64]uint64 {
	var attacks = [64]uint64{}
	for sq := uint64(0); sq < 64; sq++ {
		attacks[sq] = knightAttacks(1 << sq)
	}
	return attacks
}

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

var computedKingAttacks = initializeKingAttacks()

func initializeKingAttacks() [64]uint64 {
	var attacks = [64]uint64{}
	for sq := uint64(0); sq < 64; sq++ {
		attacks[sq] = kingAttacks(1 << sq)
	}
	return attacks
}

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

var rayAttacksArray = initializeRayAttacks()

func initializeRayAttacks() [8][64]uint64 {
	var rayAttacks = [8][64]uint64{}
	for sq := uint64(0); sq < 64; sq++ {
		rayAttacks[North][sq] = northRay(1 << sq)
		rayAttacks[NorthEast][sq] = northEastRay(1 << sq)
		rayAttacks[East][sq] = eastRay(1 << sq)
		rayAttacks[SouthEast][sq] = southEastRay(1 << sq)
		rayAttacks[South][sq] = southRay(1 << sq)
		rayAttacks[SouthWest][sq] = southWestRay(1 << sq)
		rayAttacks[West][sq] = westRay(1 << sq)
		rayAttacks[NorthWest][sq] = northWestRay(1 << sq)
	}
	return rayAttacks
}

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

// func noNoEa(b uint64) uint64 {
// 	return (b & notHFile) << 17
// }
// func noEaEa(b uint64) uint64 {
// 	return (b & notGHFile) << 10
// }
// func soEaEa(b uint64) uint64 {
// 	return (b & notGHFile) >> 6
// }
// func soSoEa(b uint64) uint64 {
// 	return (b & notHFile) >> 15
// }
// func noNoWe(b uint64) uint64 {
// 	return (b & notAFile) << 15
// }
// func noWeWe(b uint64) uint64 {
// 	return (b & notABFile) << 6
// }
// func soWeWe(b uint64) uint64 {
// 	return (b & notABFile) >> 10
// }
// func soSoWe(b uint64) uint64 {
// 	return (b & notAFile) >> 17
// }

const empty = uint64(0)
const universal = uint64(0xffffffffffffffff)
const notAFile = uint64(0xfefefefefefefefe) // ~0x0101010101010101
const notBFile = uint64(0xfdfdfdfdfdfdfdfd)
const notGFile = uint64(0xbfbfbfbfbfbfbfbf)
const notHFile = uint64(0x7f7f7f7f7f7f7f7f) // ~0x8080808080808080
const notABFile = notAFile & notBFile
const notGHFile = notGFile & notHFile
const rank4 = uint64(0x00000000FF000000)
const rank5 = uint64(0x000000FF00000000)
