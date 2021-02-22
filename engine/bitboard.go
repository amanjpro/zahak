package engine

import (
	"fmt"
	"math/bits"
)

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

func (b *Bitboard) GetWhitePieces() uint64 {
	return b.whitePieces
}

func (b *Bitboard) GetBlackPieces() uint64 {
	return b.blackPieces
}

func (b *Bitboard) GetBitboardOf(piece Piece) uint64 {
	switch piece {
	case BlackPawn:
		return b.blackPawn
	case BlackKnight:
		return b.blackKnight
	case BlackBishop:
		return b.blackBishop
	case BlackRook:
		return b.blackRook
	case BlackQueen:
		return b.blackQueen
	case BlackKing:
		return b.blackKing
	case WhitePawn:
		return b.whitePawn
	case WhiteKnight:
		return b.whiteKnight
	case WhiteBishop:
		return b.whiteBishop
	case WhiteRook:
		return b.whiteRook
	case WhiteQueen:
		return b.whiteQueen
	case WhiteKing:
		return b.whiteKing
	}
	return 0
}

func (b *Bitboard) AllPieces() map[Square]Piece {
	allPieces := make(map[Square]Piece, 32)
	allBits := b.whitePieces | b.blackPieces
	for allBits != 0 {
		index := bitScanForward(allBits)
		mask := uint64(1 << index)
		sq := Square(index)
		if b.blackPawn&(mask) != 0 {
			allPieces[sq] = BlackPawn
		} else if b.whitePawn&(mask) != 0 {
			allPieces[sq] = WhitePawn
		} else if b.blackKnight&(mask) != 0 {
			allPieces[sq] = BlackKnight
		} else if b.whiteKnight&(mask) != 0 {
			allPieces[sq] = WhiteKnight
		} else if b.blackBishop&(mask) != 0 {
			allPieces[sq] = BlackBishop
		} else if b.whiteBishop&(mask) != 0 {
			allPieces[sq] = WhiteBishop
		} else if b.blackRook&(mask) != 0 {
			allPieces[sq] = BlackRook
		} else if b.whiteRook&(mask) != 0 {
			allPieces[sq] = WhiteRook
		} else if b.blackQueen&(mask) != 0 {
			allPieces[sq] = BlackQueen
		} else if b.whiteQueen&(mask) != 0 {
			allPieces[sq] = WhiteQueen
		} else if b.blackKing&(mask) != 0 {
			allPieces[sq] = BlackKing
		} else if b.whiteKing&(mask) != 0 {
			allPieces[sq] = WhiteKing
		}
		allBits ^= mask
	}
	return allPieces
}

func (b *Bitboard) UpdateSquare(sq Square, piece Piece) {
	// Remove the piece from source square and add it to destination
	b.Clear(sq)
	mask := uint64(1 << sq)
	switch piece {
	case BlackPawn:
		b.blackPawn |= mask
	case BlackKnight:
		b.blackKnight |= mask
	case BlackBishop:
		b.blackBishop |= mask
	case BlackRook:
		b.blackRook |= mask
	case BlackQueen:
		b.blackQueen |= mask
	case BlackKing:
		b.blackKing |= mask
	case WhitePawn:
		b.whitePawn |= mask
	case WhiteKnight:
		b.whiteKnight |= mask
	case WhiteBishop:
		b.whiteBishop |= mask
	case WhiteRook:
		b.whiteRook |= mask
	case WhiteQueen:
		b.whiteQueen |= mask
	case WhiteKing:
		b.whiteKing |= mask
	}

	b.blackPieces = b.blackPawn | b.blackKnight | b.blackBishop | b.blackRook | b.blackQueen | b.blackKing
	b.whitePieces = b.whitePawn | b.whiteKnight | b.whiteBishop | b.whiteRook | b.whiteQueen | b.whiteKing
}

func (b *Bitboard) PieceAt(sq Square) Piece {
	mask := uint64(1 << sq)
	if sq == NoSquare {
		return NoPiece
	}
	if b.blackPieces&mask != 0 {
		if b.blackPawn&mask != 0 {
			return BlackPawn
		} else if b.blackKnight&mask != 0 {
			return BlackKnight
		} else if b.blackBishop&mask != 0 {
			return BlackBishop
		} else if b.blackRook&mask != 0 {
			return BlackRook
		} else if b.blackQueen&mask != 0 {
			return BlackQueen
		} else if b.blackKing&mask != 0 {
			return BlackKing
		}
	}

	// It is not black? then it is white
	if b.whitePawn&mask != 0 {
		return WhitePawn
	} else if b.whiteKnight&mask != 0 {
		return WhiteKnight
	} else if b.whiteBishop&mask != 0 {
		return WhiteBishop
	} else if b.whiteRook&mask != 0 {
		return WhiteRook
	} else if b.whiteQueen&mask != 0 {
		return WhiteQueen
	} else if b.whiteKing&mask != 0 {
		return WhiteKing
	}
	return NoPiece
}

func (b *Bitboard) Clear(square Square) {

	mask := uint64(1 << square)
	b.blackPawn &^= mask
	b.blackKnight &^= mask
	b.blackBishop &^= mask
	b.blackRook &^= mask
	b.blackQueen &^= mask
	b.blackKing &^= mask
	b.blackPieces &^= mask
	b.whitePawn &^= mask
	b.whiteKnight &^= mask
	b.whiteBishop &^= mask
	b.whiteRook &^= mask
	b.whiteQueen &^= mask
	b.whiteKing &^= mask
	b.whitePieces &^= mask
}

func (b *Bitboard) Move(src Square, dest Square) {

	// clear destination square
	b.Clear(dest)
	maskSrc := uint64(1 << src)
	maskDest := uint64(1 << dest)

	// Remove the piece from source square and add it to destination
	// is black?
	if b.blackPieces&maskSrc != 0 {
		if b.blackPawn&maskSrc != 0 {
			b.blackPawn &^= maskSrc
			b.blackPawn |= maskDest
		} else if b.blackKnight&maskSrc != 0 {
			b.blackKnight &^= maskSrc
			b.blackKnight |= maskDest
		} else if b.blackBishop&maskSrc != 0 {
			b.blackBishop &^= maskSrc
			b.blackBishop |= maskDest
		} else if b.blackRook&maskSrc != 0 {
			b.blackRook &^= maskSrc
			b.blackRook |= maskDest
		} else if b.blackQueen&maskSrc != 0 {
			b.blackQueen &^= maskSrc
			b.blackQueen |= maskDest
		} else if b.blackKing&maskSrc != 0 {
			b.blackKing &^= maskSrc
			b.blackKing |= maskDest
			// Is it a castle?
			if src == E8 && dest == G8 {
				b.Move(H8, F8)
			} else if src == E8 && dest == C8 {
				b.Move(A8, D8)
			}
		}
	}
	// Then it is white
	if b.whitePawn&maskSrc != 0 {
		b.whitePawn &^= maskSrc
		b.whitePawn |= maskDest
	} else if b.whiteKnight&maskSrc != 0 {
		b.whiteKnight &^= maskSrc
		b.whiteKnight |= maskDest
	} else if b.whiteBishop&maskSrc != 0 {
		b.whiteBishop &^= maskSrc
		b.whiteBishop |= maskDest
	} else if b.whiteRook&maskSrc != 0 {
		b.whiteRook &^= maskSrc
		b.whiteRook |= maskDest
	} else if b.whiteQueen&maskSrc != 0 {
		b.whiteQueen &^= maskSrc
		b.whiteQueen |= maskDest
	} else if b.whiteKing&maskSrc != 0 {
		b.whiteKing &^= maskSrc
		b.whiteKing |= maskDest
		// Is it a castle?
		if src == E1 && dest == G1 {
			b.Move(H1, F1)
		} else if src == E1 && dest == C1 {
			b.Move(A1, D1)
		}
	}
	b.blackPieces = b.blackPawn | b.blackKnight | b.blackBishop | b.blackRook | b.blackQueen | b.blackKing
	b.whitePieces = b.whitePawn | b.whiteKnight | b.whiteBishop | b.whiteRook | b.whiteQueen | b.whiteKing
}

func StartingBoard() Bitboard {
	bitboard := Bitboard{}
	bitboard.UpdateSquare(A2, WhitePawn)
	bitboard.UpdateSquare(B2, WhitePawn)
	bitboard.UpdateSquare(C2, WhitePawn)
	bitboard.UpdateSquare(D2, WhitePawn)
	bitboard.UpdateSquare(E2, WhitePawn)
	bitboard.UpdateSquare(F2, WhitePawn)
	bitboard.UpdateSquare(G2, WhitePawn)
	bitboard.UpdateSquare(H2, WhitePawn)

	bitboard.UpdateSquare(A7, BlackPawn)
	bitboard.UpdateSquare(B7, BlackPawn)
	bitboard.UpdateSquare(C7, BlackPawn)
	bitboard.UpdateSquare(D7, BlackPawn)
	bitboard.UpdateSquare(E7, BlackPawn)
	bitboard.UpdateSquare(F7, BlackPawn)
	bitboard.UpdateSquare(G7, BlackPawn)
	bitboard.UpdateSquare(H7, BlackPawn)

	bitboard.UpdateSquare(A1, WhiteRook)
	bitboard.UpdateSquare(B1, WhiteKnight)
	bitboard.UpdateSquare(C1, WhiteBishop)
	bitboard.UpdateSquare(D1, WhiteQueen)
	bitboard.UpdateSquare(E1, WhiteKing)
	bitboard.UpdateSquare(F1, WhiteBishop)
	bitboard.UpdateSquare(G1, WhiteKnight)
	bitboard.UpdateSquare(H1, WhiteRook)

	bitboard.UpdateSquare(A8, BlackRook)
	bitboard.UpdateSquare(B8, BlackKnight)
	bitboard.UpdateSquare(C8, BlackBishop)
	bitboard.UpdateSquare(D8, BlackQueen)
	bitboard.UpdateSquare(E8, BlackKing)
	bitboard.UpdateSquare(F8, BlackBishop)
	bitboard.UpdateSquare(G8, BlackKnight)
	bitboard.UpdateSquare(H8, BlackRook)

	return bitboard
}

func (b *Bitboard) CountPieces() int {
	return bits.OnesCount64(b.whitePieces | b.whitePieces)
}

// Draw returns visual representation of the board useful for debugging.
func (b *Bitboard) Draw() string {
	pieceUnicodes := []string{"♔", "♕", "♖", "♗", "♘", "♙", "♚", "♛", "♜", "♝", "♞", "♟"}
	s := "\n A B C D E F G H\n"
	for r := 7; r >= 0; r-- {
		s += fmt.Sprint(Rank(r + 1))
		for f := 0; f < len(Files); f++ {
			p := b.PieceAt(SquareOf(File(f), Rank(r)))
			if p == NoPiece {
				s += "-"
			} else {
				s += pieceUnicodes[int(p)]
			}
			s += " "
		}
		s += "\n"
	}
	return s
}

func (b *Bitboard) copy() *Bitboard {
	return &Bitboard{
		b.blackPawn,
		b.blackKnight,
		b.blackBishop,
		b.blackRook,
		b.blackQueen,
		b.blackKing,
		b.whitePawn,
		b.whiteKnight,
		b.whiteBishop,
		b.whiteRook,
		b.whiteQueen,
		b.whiteKing,
		b.whitePieces,
		b.blackPieces,
	}
}
