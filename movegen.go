package main

func (p *Position) LegalMoves() []*Move {
	allMoves := make([]*Move, 0, 350)

	color := p.Turn()

	add := func(ms ...*Move) {
		for _, m := range ms {
			// capture the position state
			oldTag := p.tag
			oldEnPassant := p.enPassant

			// make the move
			capturedPiece := p.MakeMove(m)

			// Does the move puts the moving player in check
			pNotInCheck := !isInCheck(p.board, color)

			if pNotInCheck && isInCheck(p.board, p.Turn()) { // We put opponent in check
				m.SetTag(Check)
				allMoves = append(allMoves, m)
			} else if pNotInCheck { // The move does not put us in check
				// do nothing
				allMoves = append(allMoves, m)
			}
			p.UnMakeMove(m, oldTag, oldEnPassant, capturedPiece)
		}
	}

	board := p.board

	taboo := tabooSquares(board, color)

	// If it is double check, only king can move
	if isDoubleCheck(board, color) {
		if color == White {
			bbKingMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
				taboo, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide), add)
		} else if color == Black {
			bbKingMoves(board.blackKing, board.blackPieces, board.whitePieces, board.whiteKing,
				taboo, p.HasTag(BlackCanCastleKingSide), p.HasTag(BlackCanCastleQueenSide), add)
		}
	} else {

		if color == White {
			bbPawnMoves(board.whitePawn, board.whitePieces, board.blackPieces,
				board.blackKing, color, p.enPassant, add)
			bbKnightMoves(board.whiteKnight, board.whitePieces, board.blackPieces,
				board.blackKing, add)
			bbSlidingMoves(board.whiteBishop, board.whitePieces, board.blackPieces, board.blackKing,
				color, bishopAttacks, add)
			bbSlidingMoves(board.whiteRook, board.whitePieces, board.blackPieces, board.blackKing,
				color, rookAttacks, add)
			bbSlidingMoves(board.whiteQueen, board.whitePieces, board.blackPieces, board.blackKing,
				color, queenAttacks, add)
			bbKingMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
				taboo, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide), add)
		} else if color == Black {
			bbPawnMoves(board.blackPawn, board.blackPieces, board.whitePieces,
				board.whiteKing, color, p.enPassant, add)
			bbKnightMoves(board.blackKnight, board.blackPieces, board.whitePieces,
				board.whiteKing, add)
			bbSlidingMoves(board.blackBishop, board.blackPieces, board.whitePieces, board.whiteKing,
				color, bishopAttacks, add)
			bbSlidingMoves(board.blackRook, board.blackPieces, board.whitePieces, board.whiteKing,
				color, rookAttacks, add)
			bbSlidingMoves(board.blackQueen, board.blackPieces, board.whitePieces, board.whiteKing,
				color, queenAttacks, add)
			bbKingMoves(board.blackKing, board.blackPieces, board.whitePieces, board.whiteKing,
				taboo, p.HasTag(BlackCanCastleKingSide), p.HasTag(BlackCanCastleQueenSide), add)
		}
	}
	return allMoves
}

func (p *Position) HasLegalMoves() bool {
	hasMoves := false

	color := p.Turn()

	add := func(ms ...*Move) {
		for _, m := range ms {
			// capture the position state
			oldTag := p.tag
			oldEnPassant := p.enPassant

			// make the move
			capturedPiece := p.MakeMove(m)

			// Does the move puts the moving player in check
			pNotInCheck := !isInCheck(p.board, color)

			if pNotInCheck && isInCheck(p.board, p.Turn()) { // We put opponent in check
				m.SetTag(Check)
				hasMoves = true
				p.UnMakeMove(m, oldTag, oldEnPassant, capturedPiece)
				break
			} else if pNotInCheck { // The move does not put us in check
				// do nothing
				hasMoves = true
				p.UnMakeMove(m, oldTag, oldEnPassant, capturedPiece)
				break
			}
			p.UnMakeMove(m, oldTag, oldEnPassant, capturedPiece)
		}
	}

	board := p.board

	taboo := tabooSquares(board, color)

	// If it is double check, only king can move
	if isDoubleCheck(board, color) {
		if color == White {
			bbKingMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
				taboo, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide), add)
		} else if color == Black {
			bbKingMoves(board.blackKing, board.blackPieces, board.whitePieces, board.whiteKing,
				taboo, p.HasTag(BlackCanCastleKingSide), p.HasTag(BlackCanCastleQueenSide), add)
		}
		return hasMoves
	} else {

		if color == White {
			bbPawnMoves(board.whitePawn, board.whitePieces, board.blackPieces,
				board.blackKing, color, p.enPassant, add)
			if hasMoves {
				return hasMoves
			}
			bbKnightMoves(board.whiteKnight, board.whitePieces, board.blackPieces,
				board.blackKing, add)
			if hasMoves {
				return hasMoves
			}
			bbSlidingMoves(board.whiteBishop, board.whitePieces, board.blackPieces, board.blackKing,
				color, bishopAttacks, add)
			if hasMoves {
				return hasMoves
			}
			bbSlidingMoves(board.whiteRook, board.whitePieces, board.blackPieces, board.blackKing,
				color, rookAttacks, add)
			if hasMoves {
				return hasMoves
			}
			bbSlidingMoves(board.whiteQueen, board.whitePieces, board.blackPieces, board.blackKing,
				color, queenAttacks, add)
			if hasMoves {
				return hasMoves
			}
			bbKingMoves(board.whiteKing, board.whitePieces, board.blackPieces, board.blackKing,
				taboo, p.HasTag(WhiteCanCastleKingSide), p.HasTag(WhiteCanCastleQueenSide), add)
			if hasMoves {
				return hasMoves
			}
		} else if color == Black {
			bbPawnMoves(board.blackPawn, board.blackPieces, board.whitePieces,
				board.whiteKing, color, p.enPassant, add)
			if hasMoves {
				return hasMoves
			}
			bbKnightMoves(board.blackKnight, board.blackPieces, board.whitePieces,
				board.whiteKing, add)
			if hasMoves {
				return hasMoves
			}
			bbSlidingMoves(board.blackBishop, board.blackPieces, board.whitePieces, board.whiteKing,
				color, bishopAttacks, add)
			if hasMoves {
				return hasMoves
			}
			bbSlidingMoves(board.blackRook, board.blackPieces, board.whitePieces, board.whiteKing,
				color, rookAttacks, add)
			if hasMoves {
				return hasMoves
			}
			bbSlidingMoves(board.blackQueen, board.blackPieces, board.whitePieces, board.whiteKing,
				color, queenAttacks, add)
			if hasMoves {
				return hasMoves
			}
			bbKingMoves(board.blackKing, board.blackPieces, board.whitePieces, board.whiteKing,
				taboo, p.HasTag(BlackCanCastleKingSide), p.HasTag(BlackCanCastleQueenSide), add)
			if hasMoves {
				return hasMoves
			}
		}
	}
	return hasMoves
}

// Checks and Pins
func isInCheck(b Bitboard, colorOfKing Color) bool {
	bbKing := b.whiteKing
	if colorOfKing == Black {
		bbKing = b.blackKing
	}
	return tabooSquares(b, colorOfKing)&bbKing != 0
}

func isDoubleCheck(b Bitboard, colorOfKing Color) bool {
	var opPawns, opKnights, opR, opB, opQ, opPieces, ownKing uint64
	occupiedBB := b.whitePieces | b.blackPieces
	if colorOfKing == White {
		opPawns = bPawnsAble2CaptureAny(b.blackPawn, b.whitePieces)
		opKnights = b.blackKnight
		opR = b.blackRook
		opB = b.blackBishop
		opQ = b.blackQueen
		opPieces = b.blackPieces
		ownKing = b.whiteKing
	} else {
		opPawns = wPawnsAble2CaptureAny(b.whitePawn, b.blackPieces)
		opKnights = b.whiteKnight
		opR = b.whiteRook
		opB = b.whiteBishop
		opQ = b.whiteQueen
		opPieces = b.whitePieces
		ownKing = b.blackKing
	}
	checkCounts := 0
	if opPawns&ownKing != 0 {
		checkCounts += 1
	}
	if knightAttacks(opKnights)&ownKing != 0 {
		checkCounts += 1
	}
	for opB != 0 {
		sq := bitScanReverse(opB)
		if bishopAttacks(Square(sq), occupiedBB, opPieces)&ownKing != 0 {
			checkCounts += 1
		}
		opB ^= (1 << sq)
	}

	for opR != 0 {
		sq := bitScanReverse(opR)
		if rookAttacks(Square(sq), occupiedBB, opPieces)&ownKing != 0 {
			checkCounts += 1
		}
		opR ^= (1 << sq)
	}

	for opQ != 0 {
		sq := bitScanReverse(opQ)
		if queenAttacks(Square(sq), occupiedBB, opPieces)&ownKing != 0 {
			checkCounts += 1
		}
		opQ ^= (1 << sq)
	}

	return checkCounts > 1
}

func tabooSquares(b Bitboard, colorOfKing Color) uint64 {
	var opPawns, opKnights, opR, opB, opQ, opKing, opPieces uint64
	occupiedBB := b.whitePieces | b.blackPieces
	if colorOfKing == White {
		opPawns = bPawnsAble2CaptureAny(b.blackPawn, b.whitePieces)
		opKnights = b.blackKnight
		opR = b.blackRook
		opB = b.blackBishop
		opQ = b.blackQueen
		opKing = b.blackKing
		opPieces = b.blackPieces
	} else {
		opPawns = wPawnsAble2CaptureAny(b.whitePawn, b.blackPieces)
		opKnights = b.whiteKnight
		opR = b.whiteRook
		opB = b.whiteBishop
		opQ = b.whiteQueen
		opKing = b.whiteKing
		opPieces = b.whitePieces
	}
	taboo := opPawns | (knightAttacks(opKnights)) | kingMovesNoCaptures(opKing, occupiedBB, 0)
	for opB != 0 {
		sq := bitScanReverse(opB)
		taboo |= bishopAttacks(Square(sq), occupiedBB, opPieces)
		opB ^= (1 << sq)
	}

	for opR != 0 {
		sq := bitScanReverse(opR)
		taboo |= rookAttacks(Square(sq), occupiedBB, opPieces)
		opR ^= (1 << sq)
	}

	for opQ != 0 {
		sq := bitScanReverse(opQ)
		taboo |= queenAttacks(Square(sq), occupiedBB, opPieces)
		opQ ^= (1 << sq)
	}

	return taboo
}

// Pawns

func bbPawnMoves(bbPawn uint64, ownPieces uint64, otherPieces uint64, otherKing uint64,
	color Color, enPassant Square, add func(m ...*Move)) {
	emptySquares := (otherPieces | ownPieces) ^ universal
	if color == White {
		for bbPawn != 0 {
			src := bitScanReverse(bbPawn)
			srcSq := Square(src)
			pawn := uint64(1 << src)
			dbl := wDblPushTargets(pawn, emptySquares)
			if dbl != 0 {
				dest := Square(bitScanReverse(dbl))
				var tag MoveTag = 0
				add(&Move{srcSq, dest, NoType, tag})
			}
			sngl := wSinglePushTargets(pawn, emptySquares)
			if sngl != 0 {
				dest := Square(bitScanReverse(sngl))
				if dest.Rank() == Rank8 {
					add(
						&Move{srcSq, dest, Queen, 0},
						&Move{srcSq, dest, Rook, 0},
						&Move{srcSq, dest, Bishop, 0},
						&Move{srcSq, dest, Knight, 0})
				} else {
					var tag MoveTag = 0
					add(&Move{srcSq, dest, NoType, tag})
				}
			}
			for _, sq := range getIndicesOfOnes(wPawnsAble2CaptureAny(pawn, otherPieces)) {
				dest := Square(sq)
				if dest.Rank() == Rank8 {
					add(
						&Move{srcSq, dest, Queen, Capture},
						&Move{srcSq, dest, Rook, Capture},
						&Move{srcSq, dest, Bishop, Capture},
						&Move{srcSq, dest, Knight, Capture})
				} else {
					var tag MoveTag = Capture
					add(&Move{srcSq, dest, NoType, tag})
				}
			}
			if srcSq.Rank() == Rank5 && enPassant != NoSquare && enPassant.Rank() == Rank6 {
				ep := uint64(1 << enPassant)
				r := wPawnsAble2CaptureAny(pawn, ep)
				if r != 0 {
					dest := Square(bitScanReverse(r))
					var tag MoveTag = Capture | EnPassant
					add(&Move{srcSq, dest, NoType, tag})
				}
			}
			bbPawn ^= pawn
		}
	} else if color == Black {
		for bbPawn != 0 {
			src := bitScanReverse(bbPawn)
			srcSq := Square(src)
			pawn := uint64(1 << src)
			dbl := bDoublePushTargets(pawn, emptySquares)
			if dbl != 0 {
				dest := Square(bitScanReverse(dbl))
				var tag MoveTag = 0
				add(&Move{srcSq, dest, NoType, tag})
			}
			sngl := bSinglePushTargets(pawn, emptySquares)
			if sngl != 0 {
				dest := Square(bitScanReverse(sngl))
				if dest.Rank() == Rank1 {
					add(
						&Move{srcSq, dest, Queen, 0},
						&Move{srcSq, dest, Rook, 0},
						&Move{srcSq, dest, Bishop, 0},
						&Move{srcSq, dest, Knight, 0})
				} else {
					add(&Move{srcSq, dest, NoType, 0})
				}
			}
			for _, sq := range getIndicesOfOnes(bPawnsAble2CaptureAny(pawn, otherPieces)) {
				dest := Square(sq)
				if dest.Rank() == Rank1 {
					add(
						&Move{srcSq, dest, Queen, Capture},
						&Move{srcSq, dest, Rook, Capture},
						&Move{srcSq, dest, Bishop, Capture},
						&Move{srcSq, dest, Knight, Capture})
				} else {
					var tag MoveTag = Capture
					add(&Move{srcSq, dest, NoType, tag})
				}
			}
			if srcSq.Rank() == Rank4 && enPassant != NoSquare && enPassant.Rank() == Rank3 {
				ep := uint64(1 << enPassant)
				r := bPawnsAble2CaptureAny(pawn, ep)
				if r != 0 {
					dest := Square(bitScanReverse(r))
					var tag MoveTag = Capture | EnPassant
					add(&Move{srcSq, dest, NoType, tag})
				}
			}
			bbPawn ^= pawn
		}
	}
}

// Sliding moves, for rooks, queens and bishops
func bbSlidingMoves(bbPiece uint64, ownPieces uint64, otherPieces uint64, otherKing uint64,
	color Color, attacks func(sq Square, occ uint64, own uint64) uint64, add func(m ...*Move)) {
	both := otherPieces | ownPieces
	for bbPiece != 0 {
		src := bitScanReverse(bbPiece)
		srcSq := Square(src)
		rayAttacks := attacks(srcSq, both, ownPieces)
		passiveMoves := rayAttacks &^ otherPieces
		captureMoves := rayAttacks & otherPieces
		for passiveMoves != 0 {
			sq := bitScanReverse(passiveMoves)
			dest := Square(sq)
			add(&Move{srcSq, dest, NoType, 0})
			passiveMoves ^= (1 << sq)
		}
		for captureMoves != 0 {
			sq := bitScanReverse(captureMoves)
			dest := Square(sq)
			add(&Move{srcSq, dest, NoType, Capture})
			captureMoves ^= (1 << sq)
		}
		bbPiece ^= (1 << src)
	}
}

// Knights
func bbKnightMoves(bbPiece uint64, ownPieces uint64, otherPieces uint64,
	otherKing uint64, add func(m ...*Move)) {
	both := otherPieces | ownPieces
	for bbPiece != 0 {
		src := bitScanReverse(bbPiece)
		srcSq := Square(src)
		knight := uint64(1 << src)
		moves := knightMovesNoCaptures(knight, both)
		for moves != 0 {
			sq := bitScanReverse(moves)
			dest := Square(sq)
			add(&Move{srcSq, dest, NoType, 0})
			moves ^= (1 << sq)
		}
		captures := knightCaptures(knight, otherPieces)
		for captures != 0 {
			sq := bitScanReverse(captures)
			dest := Square(sq)
			add(&Move{srcSq, dest, NoType, Capture})
			captures ^= (1 << sq)
		}
		bbPiece ^= knight
	}
}

// Kings
func bbKingMoves(bbPiece uint64, ownPieces uint64, otherPieces uint64, otherKing uint64,
	tabooSquares uint64, kingSideCastle bool, queenSideCastle bool,
	add func(m ...*Move)) {
	both := (otherPieces | ownPieces)
	if bbPiece != 0 {
		src := bitScanReverse(bbPiece)
		srcSq := Square(src)
		king := uint64(1 << src)
		moves := kingMovesNoCaptures(king, both, tabooSquares)
		for moves != 0 {
			sq := bitScanReverse(moves)
			dest := Square(sq)
			add(&Move{srcSq, dest, NoType, 0})
			moves ^= (1 << sq)
		}
		captures := kingCaptures(king, otherPieces, tabooSquares)
		for captures != 0 {
			sq := bitScanReverse(captures)
			dest := Square(sq)
			add(&Move{srcSq, dest, NoType, Capture})
			captures ^= (1 << sq)
		}

		E := E1
		F := F1
		G := G1
		D := D1
		C := C1
		B := B1
		if srcSq.Rank() == Rank8 {
			E = E8
			F = F8
			G = G8
			D = D8
			C = C8
			B = B8
		}

		kingSide := uint64(1<<F | 1<<G)
		queenSide := uint64(1<<D | 1<<C | 1<<B)

		if kingSideCastle &&
			(ownPieces&kingSide == 0) && // are empty
			(tabooSquares&(kingSide|1<<E) == 0) { // Not in check
			add(&Move{srcSq, G, NoType, KingSideCastle})
		}

		if queenSideCastle &&
			(ownPieces&queenSide == 0) && // are empty
			(tabooSquares&(queenSide|1<<E) == 0) { // Not in check
			add(&Move{srcSq, G, NoType, QueenSideCastle})
		}
	}
}

// Pawn Pushes

func wSinglePushTargets(wpawns uint64, empty uint64) uint64 {
	return nortOne(wpawns) & empty
}

func wDblPushTargets(wpawns uint64, empty uint64) uint64 {
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

func knightCheckTag(from uint64, otherKing uint64) MoveTag {
	if knightCaptures(from, otherKing) != 0 {
		return Check
	}
	return 0
}

func knightMovesNoCaptures(b uint64, other uint64) uint64 {
	attacks := knightAttacks(b)
	return attacks &^ other
}

func knightCaptures(b uint64, other uint64) uint64 {
	return knightAttacks(b) & other
}

func knightAttacks(b uint64) uint64 {
	return noNoEa(b) | noEaEa(b) | soEaEa(b) | soSoEa(b) |
		noNoWe(b) | noWeWe(b) | soWeWe(b) | soSoWe(b)
}

// King & Kingslayer
func kingMovesNoCaptures(b uint64, others uint64, tabooSquares uint64) uint64 {
	attacks := kingAttacks(b)
	return (attacks ^ others) & attacks & (tabooSquares ^ universal)
}

func kingCaptures(b uint64, others uint64, tabooSquares uint64) uint64 {
	return kingAttacks(b) & others & (tabooSquares ^ universal)
}

func kingAttacks(b uint64) uint64 {
	return soutOne(b) | nortOne(b) | eastOne(b) | noEaOne(b) |
		soEaOne(b) | westOne(b) | soWeOne(b) | noWeOne(b)
}

// Utilites
func getIndicesOfOnes(bb uint64) []uint8 {
	indices := make([]uint8, 0, 8)
	for bb != 0 {
		index := bitScanReverse(bb)
		bb ^= (1 << index)
		indices = append(indices, index)
	}
	return indices
}

var forwardIndex = [64]uint8{
	0, 1, 48, 2, 57, 49, 28, 3,
	61, 58, 50, 42, 38, 29, 17, 4,
	62, 55, 59, 36, 53, 51, 43, 22,
	45, 39, 33, 30, 24, 18, 12, 5,
	63, 47, 56, 27, 60, 41, 37, 16,
	54, 35, 52, 21, 44, 32, 23, 11,
	46, 26, 40, 15, 34, 20, 31, 10,
	25, 14, 19, 9, 13, 8, 7, 6,
}

/**
 * bitScanForward
 * @author Martin Läuter (1997)
 *         Charles E. Leiserson
 *         Harald Prokop
 *         Keith H. Randall
 * "Using de Bruijn Sequences to Index a 1 in a Computer Word"
 * @param bb bitboard to scan
 * @precondition bb != 0
 * @return index (0..63) of least significant one bit
 */
func bitScanForward(bb uint64) uint8 {
	const debruijn64 = uint64(0x03f79d71b4cb0a89)
	return forwardIndex[((bb&-bb)*debruijn64)>>58]
}

var reverseIndex = [64]uint8{
	0, 47, 1, 56, 48, 27, 2, 60,
	57, 49, 41, 37, 28, 16, 3, 61,
	54, 58, 35, 52, 50, 42, 21, 44,
	38, 32, 29, 23, 17, 11, 4, 62,
	46, 55, 26, 59, 40, 36, 15, 53,
	34, 51, 20, 43, 31, 22, 10, 45,
	25, 39, 14, 33, 19, 30, 9, 24,
	13, 18, 8, 12, 7, 6, 5, 63,
}

/**
 * bitScanReverse
 * @authors Kim Walisch, Mark Dickinson
 * @param bb bitboard to scan
 * @precondition bb != 0
 * @return index (0..63) of most significant one bit
 */
func bitScanReverse(bb uint64) uint8 {
	const debruijn64 = uint64(0x03f79d71b4cb0a89)
	bb |= bb >> 1
	bb |= bb >> 2
	bb |= bb >> 4
	bb |= bb >> 8
	bb |= bb >> 16
	bb |= bb >> 32
	return reverseIndex[(bb*debruijn64)>>58]
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
