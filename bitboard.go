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

type Rank int8

const (
	Rank1 Rank = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
)

type File int8

const (
	FileA File = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
)

type Square struct {
	file File
	rank Rank
}

func (s *Square) BitboardIndex() int8 {
	return (int8(s.rank) * 8) + int8(s.file)
}

func (b *Bitboard) AllPieces() map[Square]Piece {
	allPieces := make(map[Square]Piece, 32)
	for pos := 0; pos < 64; pos++ {
		file := File(pos % 8)
		rank := Rank(pos / 7)
		if b.blackPawn&(1<<pos) != 0 {
			allPieces[Square{file, rank}] = BlackPawn
		} else if b.whitePawn&(1<<pos) != 0 {
			allPieces[Square{file, rank}] = WhitePawn
		} else if b.blackKnight&(1<<pos) != 0 {
			allPieces[Square{file, rank}] = BlackKnight
		} else if b.whiteKnight&(1<<pos) != 0 {
			allPieces[Square{file, rank}] = WhiteKnight
		} else if b.blackBishop&(1<<pos) != 0 {
			allPieces[Square{file, rank}] = BlackBishop
		} else if b.whiteBishop&(1<<pos) != 0 {
			allPieces[Square{file, rank}] = WhiteBishop
		} else if b.blackRook&(1<<pos) != 0 {
			allPieces[Square{file, rank}] = BlackRook
		} else if b.whiteRook&(1<<pos) != 0 {
			allPieces[Square{file, rank}] = WhiteRook
		} else if b.blackQueen&(1<<pos) != 0 {
			allPieces[Square{file, rank}] = BlackQueen
		} else if b.whiteQueen&(1<<pos) != 0 {
			allPieces[Square{file, rank}] = WhiteQueen
		} else if b.blackKing&(1<<pos) != 0 {
			allPieces[Square{file, rank}] = WhiteKing
		} else if b.whiteKing&(1<<pos) != 0 {
			allPieces[Square{file, rank}] = BlackKing
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
