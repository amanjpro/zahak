package engine

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func (b *Bitboard) Fen() string {
	fen := ""
	for i := len(Ranks) - 1; i >= 0; i-- {
		rank := Ranks[i]
		empty := 0
		for j := 0; j < len(Files); j++ {
			file := Files[j]
			sq := SquareOf(file, rank)
			piece := b.PieceAt(sq)
			if piece == NoPiece {
				empty += 1
			} else {
				if empty != 0 {
					fen = fmt.Sprintf("%s%d%s", fen, empty, piece.Name())
					empty = 0
				} else {
					fen = fmt.Sprintf("%s%s", fen, piece.Name())
				}
			}
		}
		if empty != 0 {
			fen = fmt.Sprintf("%s%d", fen, empty)
		}
		if rank != Rank1 {
			fen = fmt.Sprintf("%s/", fen)
		}
	}
	return fen
}

func (p *Position) Fen() string {
	fen := p.Board.Fen()
	if p.Turn() == Black {
		fen = fmt.Sprintf("%s b ", fen)
	} else {
		fen = fmt.Sprintf("%s w ", fen)
	}

	nocastle := true
	if p.HasTag(WhiteCanCastleKingSide) {
		fen = fmt.Sprintf("%sK", fen)
		nocastle = false
	}
	if p.HasTag(WhiteCanCastleQueenSide) {
		fen = fmt.Sprintf("%sQ", fen)
		nocastle = false
	}
	if p.HasTag(BlackCanCastleKingSide) {
		fen = fmt.Sprintf("%sk", fen)
		nocastle = false
	}
	if p.HasTag(BlackCanCastleQueenSide) {
		fen = fmt.Sprintf("%sq", fen)
		nocastle = false
	}
	if nocastle {
		fen = fmt.Sprintf("%s-", fen)
	}
	if p.EnPassant != NoSquare {
		fen = fmt.Sprintf("%s %s", fen, p.EnPassant.Name())
	} else {
		fen = fmt.Sprintf("%s -", fen)
	}
	return fen
}

func (g *Game) Fen() string {
	fen := fmt.Sprintf("%s %d %d", g.position.Fen(), g.halfMoveClock, g.numberOfMoves)
	return fen
}

func bitboardFromFen(fen string) Bitboard {
	board := Bitboard{}
	ranks := []Square{A8, A7, A6, A5, A4, A3, A2, A1}
	rank := 0
	bitboardIndex := A8
	for _, ch := range fen {
		if ch == ' ' || rank >= len(ranks) {
			break // end of the board
		} else if unicode.IsDigit(ch) {
			n, _ := strconv.Atoi(string(ch))
			bitboardIndex += Square(n)
		} else if ch == '/' && bitboardIndex%8 == 0 {
			rank++
			bitboardIndex = ranks[rank]
			continue
		} else if p := pieceFromName(ch); p != NoPiece {
			board.UpdateSquare(bitboardIndex, p)
			bitboardIndex++
		} else {
			panic(fmt.Sprintf("Invalid FEN notation %s, bitboardIndex == %d, parsing %s",
				fen, bitboardIndex, string(ch)))
		}
	}
	return board
}

func positionFromFen(fen string) Position {
	parts := strings.Fields(fen)
	if len(parts) != 6 {
		panic(fmt.Sprintf("Invalid FEN notation %s, there should be 6 parts", fen))
	}
	p := Position{
		bitboardFromFen(fen),
		NoSquare,
		0,
	}

	if parts[1] == "b" {
		p.SetTag(BlackToMove)
	} else {
		p.SetTag(WhiteToMove)
	}

	for i, ch := range parts[2] {
		if ch == 'K' {
			p.SetTag(WhiteCanCastleKingSide)
		} else if ch == 'Q' {
			p.SetTag(WhiteCanCastleQueenSide)
		} else if ch == 'k' {
			p.SetTag(BlackCanCastleKingSide)
		} else if ch == 'q' {
			p.SetTag(BlackCanCastleQueenSide)
		} else if ch == '-' && i == len(parts[2])-1 {
			break
		} else {
			panic(fmt.Sprintf("Invalid FEN notation %s, castling part is not correct %s", fen, parts[2]))
		}
	}

	sq, ok := NameToSquareMap[parts[3]]
	rank := sq.Rank()
	if !ok && parts[3] != "-" {
		panic(fmt.Sprintf("Invalid FEN notation %s, en-passant part is not correct '%s'", fen, parts[3]))
	} else if ok && rank != Rank3 && rank != Rank6 {
		panic(fmt.Sprintf("Invalid FEN notation %s, en-passant part is not on the right rank %s", fen, parts[3]))
	} else if rank == Rank3 && p.Turn() == White {
		panic(fmt.Sprintf("Invalid FEN notation %s, en-passant part is not on the right rank %s", fen, parts[3]))
	} else if rank == Rank5 && p.Turn() == Black {
		panic(fmt.Sprintf("Invalid FEN notation %s, en-passant part is not on the right rank %s", fen, parts[3]))
	} else if ok {
		p.EnPassant = sq
	}
	return p
}

func FromFen(fen string, clearCache bool) Game {
	parts := strings.Fields(fen)
	if len(parts) != 6 {
		panic(fmt.Sprintf("Invalid FEN notation %s, there should be 6 parts", fen))
	}
	p := positionFromFen(fen)

	halfMoveClock, e1 := strconv.Atoi(parts[4])
	moveCount, e2 := strconv.Atoi(parts[5])
	if e1 != nil {
		panic(fmt.Sprintf("Invalid FEN notation %s, half move clock is not set correctly %s", fen, parts[4]))
	}
	if e2 != nil {
		panic(fmt.Sprintf("Invalid FEN notation %s, move count is not set correctly %s", fen, parts[5]))
	}

	return NewGame(
		&p,
		*p.copy(),
		[]*Move{},
		make(map[uint64]int8, 200),
		uint16(moveCount),
		uint16(halfMoveClock),
		clearCache,
	)
}
