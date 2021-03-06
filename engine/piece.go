package engine

type Piece int8

const (
	NoPiece Piece = iota
	WhitePawn
	WhiteKnight
	WhiteBishop
	WhiteRook
	WhiteQueen
	WhiteKing
	BlackPawn
	BlackKnight
	BlackBishop
	BlackRook
	BlackQueen
	BlackKing
)

type PieceType int8

const (
	NoType PieceType = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

type Color int8

const (
	NoColor Color = iota
	White
	Black
)

func (c *Color) Other() Color {
	if *c == White {
		return Black
	} else if *c == Black {
		return White
	}
	return NoColor
}

func (p *Piece) Type() PieceType {
	switch *p {
	case WhitePawn, BlackPawn:
		return Pawn
	case WhiteKnight, BlackKnight:
		return Knight
	case WhiteBishop, BlackBishop:
		return Bishop
	case WhiteRook, BlackRook:
		return Rook
	case WhiteQueen, BlackQueen:
		return Queen
	case WhiteKing, BlackKing:
		return King
	}
	return NoType
}

const MAX_INT = int32(2_147_483_647)

func (p *Piece) Weight() int32 {
	switch *p {
	case WhitePawn, BlackPawn:
		return 100
	case WhiteKnight, BlackKnight:
		return 325
	case WhiteBishop, BlackBishop:
		return 325
	case WhiteRook, BlackRook:
		return 550
	case WhiteQueen, BlackQueen:
		return 1000
	case WhiteKing, BlackKing:
		return MAX_INT
	}
	return 0
}

func (p *Piece) Name() string {
	switch *p {
	case WhitePawn:
		return "P"
	case WhiteKnight:
		return "N"
	case WhiteBishop:
		return "B"
	case WhiteRook:
		return "R"
	case WhiteQueen:
		return "Q"
	case WhiteKing:
		return "K"
	case BlackPawn:
		return "p"
	case BlackKnight:
		return "n"
	case BlackBishop:
		return "b"
	case BlackRook:
		return "r"
	case BlackQueen:
		return "q"
	case BlackKing:
		return "k"
	}
	return "nothing"
}

func pieceFromName(name rune) Piece {
	switch name {
	case 'P':
		return WhitePawn
	case 'N':
		return WhiteKnight
	case 'B':
		return WhiteBishop
	case 'R':
		return WhiteRook
	case 'Q':
		return WhiteQueen
	case 'K':
		return WhiteKing
	case 'p':
		return BlackPawn
	case 'n':
		return BlackKnight
	case 'b':
		return BlackBishop
	case 'r':
		return BlackRook
	case 'q':
		return BlackQueen
	case 'k':
		return BlackKing
	}
	return NoPiece
}

func (p *Piece) Color() Color {
	switch *p {
	case WhitePawn, WhiteKnight, WhiteBishop, WhiteRook, WhiteQueen, WhiteKing:
		return White
	case NoPiece:
		return NoColor
	}
	return Black
}

func GetPiece(pieceType PieceType, color Color) Piece {
	if color == White {
		switch pieceType {
		case King:
			return WhiteKing
		case Queen:
			return WhiteQueen
		case Rook:
			return WhiteRook
		case Bishop:
			return WhiteBishop
		case Knight:
			return WhiteKnight
		case Pawn:
			return WhitePawn
		}
	}
	if color == Black {
		switch pieceType {
		case King:
			return BlackKing
		case Queen:
			return BlackQueen
		case Rook:
			return BlackRook
		case Bishop:
			return BlackBishop
		case Knight:
			return BlackKnight
		case Pawn:
			return BlackPawn
		}
	}
	return NoPiece
}
