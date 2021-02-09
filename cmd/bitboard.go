package cmd

import (
	"fmt"
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

func (b *Bitboard) AllPieces() map[Square]Piece {
	allPieces := make(map[Square]Piece, 32)
	for sq := A1; sq <= H8 && len(allPieces) <= 32; sq++ {
		if b.blackPawn&(1<<sq) != 0 {
			allPieces[sq] = BlackPawn
		} else if b.whitePawn&(1<<sq) != 0 {
			allPieces[sq] = WhitePawn
		} else if b.blackKnight&(1<<sq) != 0 {
			allPieces[sq] = BlackKnight
		} else if b.whiteKnight&(1<<sq) != 0 {
			allPieces[sq] = WhiteKnight
		} else if b.blackBishop&(1<<sq) != 0 {
			allPieces[sq] = BlackBishop
		} else if b.whiteBishop&(1<<sq) != 0 {
			allPieces[sq] = WhiteBishop
		} else if b.blackRook&(1<<sq) != 0 {
			allPieces[sq] = BlackRook
		} else if b.whiteRook&(1<<sq) != 0 {
			allPieces[sq] = WhiteRook
		} else if b.blackQueen&(1<<sq) != 0 {
			allPieces[sq] = BlackQueen
		} else if b.whiteQueen&(1<<sq) != 0 {
			allPieces[sq] = WhiteQueen
		} else if b.blackKing&(1<<sq) != 0 {
			allPieces[sq] = BlackKing
		} else if b.whiteKing&(1<<sq) != 0 {
			allPieces[sq] = WhiteKing
		}
	}
	return allPieces
}

func (b *Bitboard) UpdateSquare(sq Square, piece Piece) {
	// Remove the piece from source square and add it to destination
	b.Clear(sq)
	switch piece {
	case BlackPawn:
		b.blackPawn |= (1 << sq)
	case BlackKnight:
		b.blackKnight |= (1 << sq)
	case BlackBishop:
		b.blackBishop |= (1 << sq)
	case BlackRook:
		b.blackRook |= (1 << sq)
	case BlackQueen:
		b.blackQueen |= (1 << sq)
	case BlackKing:
		b.blackKing |= (1 << sq)
	case WhitePawn:
		b.whitePawn |= (1 << sq)
	case WhiteKnight:
		b.whiteKnight |= (1 << sq)
	case WhiteBishop:
		b.whiteBishop |= (1 << sq)
	case WhiteRook:
		b.whiteRook |= (1 << sq)
	case WhiteQueen:
		b.whiteQueen |= (1 << sq)
	case WhiteKing:
		b.whiteKing |= (1 << sq)
	}

	b.blackPieces = b.blackPawn | b.blackKnight | b.blackBishop | b.blackRook | b.blackQueen | b.blackKing
	b.whitePieces = b.whitePawn | b.whiteKnight | b.whiteBishop | b.whiteRook | b.whiteQueen | b.whiteKing
}

func (b *Bitboard) PieceAt(sq Square) Piece {
	if sq == NoSquare {
		return NoPiece
	} else if b.blackPawn&(1<<sq) != 0 {
		return BlackPawn
	} else if b.whitePawn&(1<<sq) != 0 {
		return WhitePawn
	} else if b.blackKnight&(1<<sq) != 0 {
		return BlackKnight
	} else if b.whiteKnight&(1<<sq) != 0 {
		return WhiteKnight
	} else if b.blackBishop&(1<<sq) != 0 {
		return BlackBishop
	} else if b.whiteBishop&(1<<sq) != 0 {
		return WhiteBishop
	} else if b.blackRook&(1<<sq) != 0 {
		return BlackRook
	} else if b.whiteRook&(1<<sq) != 0 {
		return WhiteRook
	} else if b.blackQueen&(1<<sq) != 0 {
		return BlackQueen
	} else if b.whiteQueen&(1<<sq) != 0 {
		return WhiteQueen
	} else if b.blackKing&(1<<sq) != 0 {
		return BlackKing
	} else if b.whiteKing&(1<<sq) != 0 {
		return WhiteKing
	}
	return NoPiece
}

func (b *Bitboard) Clear(square Square) {

	b.blackPawn &^= (1 << square)
	b.blackKnight &^= (1 << square)
	b.blackBishop &^= (1 << square)
	b.blackRook &^= (1 << square)
	b.blackQueen &^= (1 << square)
	b.blackKing &^= (1 << square)
	b.blackPieces &^= (1 << square)
	b.whitePawn &^= (1 << square)
	b.whiteKnight &^= (1 << square)
	b.whiteBishop &^= (1 << square)
	b.whiteRook &^= (1 << square)
	b.whiteQueen &^= (1 << square)
	b.whiteKing &^= (1 << square)
	b.whitePieces &^= (1 << square)
}

func (b *Bitboard) Move(src Square, dest Square) {

	// clear destination square
	b.Clear(dest)

	// Remove the piece from source square and add it to destination
	if b.blackPawn&(1<<src) != 0 {
		b.blackPawn &^= (1 << src)
		b.blackPawn |= (1 << dest)
	} else if b.whitePawn&(1<<src) != 0 {
		b.whitePawn &^= (1 << src)
		b.whitePawn |= (1 << dest)
	} else if b.blackKnight&(1<<src) != 0 {
		b.blackKnight &^= (1 << src)
		b.blackKnight |= (1 << dest)
	} else if b.whiteKnight&(1<<src) != 0 {
		b.whiteKnight &^= (1 << src)
		b.whiteKnight |= (1 << dest)
	} else if b.blackBishop&(1<<src) != 0 {
		b.blackBishop &^= (1 << src)
		b.blackBishop |= (1 << dest)
	} else if b.whiteBishop&(1<<src) != 0 {
		b.whiteBishop &^= (1 << src)
		b.whiteBishop |= (1 << dest)
	} else if b.blackRook&(1<<src) != 0 {
		b.blackRook &^= (1 << src)
		b.blackRook |= (1 << dest)
	} else if b.whiteRook&(1<<src) != 0 {
		b.whiteRook &^= (1 << src)
		b.whiteRook |= (1 << dest)
	} else if b.blackQueen&(1<<src) != 0 {
		b.blackQueen &^= (1 << src)
		b.blackQueen |= (1 << dest)
	} else if b.whiteQueen&(1<<src) != 0 {
		b.whiteQueen &^= (1 << src)
		b.whiteQueen |= (1 << dest)
	} else if b.blackKing&(1<<src) != 0 {
		b.blackKing &^= (1 << src)
		b.blackKing |= (1 << dest)
		// Is it a castle?
		if src == E8 && dest == G8 {
			b.Move(H8, F8)
		} else if src == E8 && dest == C8 {
			b.Move(A8, D8)
		}
	} else if b.whiteKing&(1<<src) != 0 {
		b.whiteKing &^= (1 << src)
		b.whiteKing |= (1 << dest)
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

// Draw returns visual representation of the board useful for debugging.
func (b *Bitboard) Draw() string {
	pieceUnicodes := []string{"♔", "♕", "♖", "♗", "♘", "♙", "♚", "♛", "♜", "♝", "♞", "♟"}
	s := "\n A B C D E F G H\n"
	for r := 7; r >= 0; r-- {
		s += fmt.Sprint(Rank(r + 1))
		for f := 0; f < len(files); f++ {
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
