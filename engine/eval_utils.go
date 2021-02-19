package engine

/**
A bunch of bittwiddling functions mainly meant for evaluation and search
This is not about move generation
*/

func (b *Bitboard) attacksTo(occupied uint64, sq Square) uint64 {
	knights := b.blackKnight | b.whiteKnight
	kings := b.blackKing | b.whiteKing
	rooksQueens := b.blackQueen | b.whiteQueen
	bishopsQueens := rooksQueens
	rooksQueens |= b.blackRook | b.whiteRook
	bishopsQueens |= b.blackBishop | b.whiteBishop

	sqMask := uint64(1 << sq)
	return wPawnsAble2CaptureAny(sqMask, b.blackPawn) |
		bPawnsAble2CaptureAny(sqMask, b.whitePawn) |
		(computedKnightAttacks[sq] & knights) |
		(computedKingAttacks[sq] & kings) |
		(bishopAttacks(sq, occupied, empty) & bishopsQueens) |
		(rookAttacks(sq, occupied, empty) & rooksQueens)
}

func (b *Bitboard) getLeastValuablePiece(attacks uint64, color Color) (uint64, Piece) {
	shift := uint8(0)

	if attacks == 0 {
		return 0, NoPiece
	}

	if color == Black {
		shift = uint8(6)
	}

	for piece := uint8(WhitePawn) + shift; piece >= uint8(WhiteKing)+shift; piece-- {
		subset := attacks & b.GetBitboardOf(Piece(piece))
		if subset != 0 {
			return subset & -subset, Piece(piece) // The piece and its location on the board
		}
	}
	return 0, NoPiece
}

func (b *Bitboard) StaticExchangeEval(toSq Square, target Piece, frSq Square, aPiece Piece) int {

	gain := make([]int, 32)
	d := 0

	mayXray := /* b.blackPawn | b.whitePawn | */ b.blackBishop | b.whiteBishop |
		b.blackRook | b.whiteRook | b.blackQueen | b.whiteQueen

	fromSet := uint64(1 << frSq)
	occupied := b.whitePieces | b.blackPieces
	attacks := b.attacksTo(occupied, toSq)

	// Ray Attacks, to update the attack def
	rooksQueens := b.blackQueen | b.whiteQueen
	bishopsQueens := rooksQueens
	rooksQueens |= b.blackRook | b.whiteRook
	bishopsQueens |= b.blackBishop | b.whiteBishop

	gain[d] = target.Weight()

	for fromSet != 0 {
		d++                                   // next depth and side
		gain[d] = aPiece.Weight() - gain[d-1] // speculative store, if defended
		if max(-gain[d-1], gain[d]) < 0 {
			break // pruning does not influence the result
		}
		attacks ^= fromSet       // reset bit in set to traverse
		occupied ^= fromSet      // reset bit in temporary occupancy (for x-Rays)
		bishopsQueens ^= fromSet // reset bit in temporary occupancy for bishops/queens
		rooksQueens ^= fromSet   // reset bit in temporary occupancy for rooks/queens
		if fromSet&mayXray != 0 {
			attacks |= (bishopAttacks(toSq, occupied, empty) & bishopsQueens)
			attacks |= (rookAttacks(toSq, occupied, empty) & rooksQueens)
		}
		color := aPiece.Color()
		fromSet, aPiece = b.getLeastValuablePiece(attacks, color.Other())
	}
	for d--; d > 0; d-- {
		gain[d-1] = -max(-gain[d-1], gain[d])
	}
	return gain[0]
}

func max(x int, y int) int {
	if x > y {
		return x
	}
	return y
}
