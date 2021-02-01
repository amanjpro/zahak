package main

type Bitboard struct {
	blackPawn   uint64
	blackKnight uint64
	blackBishop uint64
	blackRook   uint64
	blackQueen  uint64
	blackKing   uint64
	whitePawn   uint64
	whiteKnight uint64
	whiteBishop uint64
	whiteRook   uint64
	whiteQueen  uint64
	whiteKing   uint64
	whitePieces uint64
	blackPieces uint64
}

func (b *Bitboard) AllPieces() map[Square]Piece {
	allPieces := make(map[Square]Piece, 32)
	for pos := 0; pos < 64; pos++ {
		file := File(pos % 8)
		rank := Rank(pos / 8)
		if b.blackPawn&(1<<pos) != 0 {
			allPieces[SquareOf(file, rank)] = BlackPawn
		} else if b.whitePawn&(1<<pos) != 0 {
			allPieces[SquareOf(file, rank)] = WhitePawn
		} else if b.blackKnight&(1<<pos) != 0 {
			allPieces[SquareOf(file, rank)] = BlackKnight
		} else if b.whiteKnight&(1<<pos) != 0 {
			allPieces[SquareOf(file, rank)] = WhiteKnight
		} else if b.blackBishop&(1<<pos) != 0 {
			allPieces[SquareOf(file, rank)] = BlackBishop
		} else if b.whiteBishop&(1<<pos) != 0 {
			allPieces[SquareOf(file, rank)] = WhiteBishop
		} else if b.blackRook&(1<<pos) != 0 {
			allPieces[SquareOf(file, rank)] = BlackRook
		} else if b.whiteRook&(1<<pos) != 0 {
			allPieces[SquareOf(file, rank)] = WhiteRook
		} else if b.blackQueen&(1<<pos) != 0 {
			allPieces[SquareOf(file, rank)] = BlackQueen
		} else if b.whiteQueen&(1<<pos) != 0 {
			allPieces[SquareOf(file, rank)] = WhiteQueen
		} else if b.blackKing&(1<<pos) != 0 {
			allPieces[SquareOf(file, rank)] = BlackKing
		} else if b.whiteKing&(1<<pos) != 0 {
			allPieces[SquareOf(file, rank)] = WhiteKing
		}
	}
	return allPieces
}

func (b *Bitboard) UpdateSquare(square *Square, piece Piece) {

	pos := square.BitboardIndex()

	// Remove the piece from source square and add it to destination
	switch piece {
	case BlackPawn:
		b.blackPawn |= (1 << pos)
		b.blackPieces |= (1 << pos)
	case BlackKnight:
		b.blackKnight |= (1 << pos)
		b.blackPieces |= (1 << pos)
	case BlackBishop:
		b.blackBishop |= (1 << pos)
		b.blackPieces |= (1 << pos)
	case BlackRook:
		b.blackRook |= (1 << pos)
		b.blackPieces |= (1 << pos)
	case BlackQueen:
		b.blackQueen |= (1 << pos)
		b.blackPieces |= (1 << pos)
	case BlackKing:
		b.blackKing |= (1 << pos)
		b.blackPieces |= (1 << pos)
	case WhitePawn:
		b.whitePawn |= (1 << pos)
		b.whitePieces |= (1 << pos)
	case WhiteKnight:
		b.whiteKnight |= (1 << pos)
		b.whitePieces |= (1 << pos)
	case WhiteBishop:
		b.whiteBishop |= (1 << pos)
		b.whitePieces |= (1 << pos)
	case WhiteRook:
		b.whiteRook |= (1 << pos)
		b.whitePieces |= (1 << pos)
	case WhiteQueen:
		b.whiteQueen |= (1 << pos)
		b.whitePieces |= (1 << pos)
	case WhiteKing:
		b.whiteKing |= (1 << pos)
		b.whitePieces |= (1 << pos)
	}
}

func (b *Bitboard) PieceAt(sq *Square) Piece {
	pos := sq.BitboardIndex()
	return b.PieceAtIndex(pos)
}

func (b *Bitboard) PieceAtIndex(pos int8) Piece {
	if b.blackPawn&(1<<pos) != 0 {
		return BlackPawn
	} else if b.whitePawn&(1<<pos) != 0 {
		return WhitePawn
	} else if b.blackKnight&(1<<pos) != 0 {
		return BlackKnight
	} else if b.whiteKnight&(1<<pos) != 0 {
		return WhiteKnight
	} else if b.blackBishop&(1<<pos) != 0 {
		return BlackBishop
	} else if b.whiteBishop&(1<<pos) != 0 {
		return WhiteBishop
	} else if b.blackRook&(1<<pos) != 0 {
		return BlackRook
	} else if b.whiteRook&(1<<pos) != 0 {
		return WhiteRook
	} else if b.blackQueen&(1<<pos) != 0 {
		return BlackQueen
	} else if b.whiteQueen&(1<<pos) != 0 {
		return WhiteQueen
	} else if b.blackKing&(1<<pos) != 0 {
		return BlackKing
	} else if b.whiteKing&(1<<pos) != 0 {
		return WhiteKing
	}
	return NoPiece
}

func (b *Bitboard) Clear(square *Square) {

	pos := square.BitboardIndex()

	b.blackPawn &= ^(1 << pos)
	b.blackKnight &= ^(1 << pos)
	b.blackBishop &= ^(1 << pos)
	b.blackRook &= ^(1 << pos)
	b.blackQueen &= ^(1 << pos)
	b.blackKing &= ^(1 << pos)
	b.blackPieces &= ^(1 << pos)
	b.whitePawn &= ^(1 << pos)
	b.whiteKnight &= ^(1 << pos)
	b.whiteBishop &= ^(1 << pos)
	b.whiteRook &= ^(1 << pos)
	b.whiteQueen &= ^(1 << pos)
	b.whiteKing &= ^(1 << pos)
	b.whitePieces &= ^(1 << pos)
}

func (b *Bitboard) Move(s1 *Square, s2 *Square) {
	src := s1.BitboardIndex()
	dest := s2.BitboardIndex()

	// clear destination square
	b.Clear(s2)

	// Remove the piece from source square and add it to destination
	if b.blackPawn&(1<<src) != 0 {
		b.blackPawn &= ^(1 << src)
		b.blackPawn |= (1 << dest)
		b.blackPieces |= (1 << dest)
	} else if b.whitePawn&(1<<src) != 0 {
		b.whitePawn &= ^(1 << src)
		b.whitePawn |= (1 << dest)
		b.whitePieces |= (1 << dest)
	} else if b.blackKnight&(1<<src) != 0 {
		b.blackKnight &= ^(1 << src)
		b.blackKnight |= (1 << dest)
		b.blackPieces |= (1 << dest)
	} else if b.whiteKnight&(1<<src) != 0 {
		b.whiteKnight &= ^(1 << src)
		b.whiteKnight |= (1 << dest)
		b.whitePieces |= (1 << dest)
	} else if b.blackBishop&(1<<src) != 0 {
		b.blackBishop &= ^(1 << src)
		b.blackBishop |= (1 << dest)
		b.blackPieces |= (1 << dest)
	} else if b.whiteBishop&(1<<src) != 0 {
		b.whiteBishop &= ^(1 << src)
		b.whiteBishop |= (1 << dest)
		b.whitePieces |= (1 << dest)
	} else if b.blackRook&(1<<src) != 0 {
		b.blackRook &= ^(1 << src)
		b.blackRook |= (1 << dest)
		b.blackPieces |= (1 << dest)
	} else if b.whiteRook&(1<<src) != 0 {
		b.whiteRook &= ^(1 << src)
		b.whiteRook |= (1 << dest)
		b.whitePieces |= (1 << dest)
	} else if b.blackQueen&(1<<src) != 0 {
		b.blackQueen &= ^(1 << src)
		b.blackQueen |= (1 << dest)
		b.blackPieces |= (1 << dest)
	} else if b.whiteQueen&(1<<src) != 0 {
		b.whiteQueen &= ^(1 << src)
		b.whiteQueen |= (1 << dest)
		b.whitePieces |= (1 << dest)
	} else if b.blackKing&(1<<src) != 0 {
		// Is it a castle?
		if *s1 == E8 && *s2 == G8 {
			b.Move(&H8, &F8)
		} else if *s1 == E8 && *s2 == C8 {
			b.Move(&A8, &D8)
		}
		b.blackKing &= ^(1 << src)
		b.blackKing |= (1 << dest)
		b.blackPieces |= (1 << dest)
	} else if b.whiteKing&(1<<src) != 0 {
		// Is it a castle?
		if *s1 == E1 && *s2 == G1 {
			b.Move(&H1, &F1)
		} else if *s1 == E1 && *s2 == C1 {
			b.Move(&A1, &D1)
		}
		b.whiteKing &= ^(1 << src)
		b.whiteKing |= (1 << dest)
		b.whitePieces |= (1 << dest)
	}
}

func StartingBoard() Bitboard {
	bitboard := Bitboard{}
	for pos := 0; pos < 16; pos++ {
		bitboard.whitePieces |= (1 << pos)
	}

	for pos := 48; pos < 64; pos++ {
		bitboard.blackPieces |= (1 << pos)
	}

	for pos := 8; pos < 16; pos++ {
		bitboard.whitePawn |= (1 << pos)
	}

	for pos := 48; pos < 56; pos++ {
		bitboard.blackPawn |= (1 << pos)
	}

	bitboard.whiteRook |= (1 << 0)
	bitboard.whiteRook |= (1 << 7)
	bitboard.whiteKnight |= (1 << 1)
	bitboard.whiteKnight |= (1 << 6)
	bitboard.whiteBishop |= (1 << 2)
	bitboard.whiteBishop |= (1 << 5)
	bitboard.whiteQueen |= (1 << 3)
	bitboard.whiteKing |= (1 << 4)

	bitboard.blackRook |= (1 << 56)
	bitboard.blackRook |= (1 << 63)
	bitboard.blackKnight |= (1 << 57)
	bitboard.blackKnight |= (1 << 62)
	bitboard.blackBishop |= (1 << 58)
	bitboard.blackBishop |= (1 << 61)
	bitboard.blackQueen |= (1 << 59)
	bitboard.blackKing |= (1 << 60)

	return bitboard
}
