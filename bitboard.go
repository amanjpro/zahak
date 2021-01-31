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

func (b *Bitboard) Move(s1 Square, s2 Square) {

	src := s1.BitboardIndex()
	dest := s2.BitboardIndex()

	// clear destination square
	b.blackPawn &= ^(1 << dest)
	b.blackKnight &= ^(1 << dest)
	b.blackBishop &= ^(1 << dest)
	b.blackRook &= ^(1 << dest)
	b.blackQueen &= ^(1 << dest)
	b.blackKing &= ^(1 << dest)
	b.blackPieces &= ^(1 << dest)
	b.whitePawn &= ^(1 << dest)
	b.whiteKnight &= ^(1 << dest)
	b.whiteBishop &= ^(1 << dest)
	b.whiteRook &= ^(1 << dest)
	b.whiteQueen &= ^(1 << dest)
	b.whiteKing &= ^(1 << dest)
	b.whitePieces &= ^(1 << dest)

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
		b.blackKing &= ^(1 << src)
		b.blackKing |= (1 << dest)
		b.blackPieces |= (1 << dest)
	} else if b.whiteKing&(1<<src) != 0 {
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
