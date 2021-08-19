package evaluation

import (
	"math/bits"

	. "github.com/amanjpro/zahak/engine"
)

type Eval struct {
	blackMG int16
	whiteMG int16
	blackEG int16
	whiteEG int16
}

const CHECKMATE_EVAL int16 = 30000
const MAX_NON_CHECKMATE int16 = 25000
const PawnPhase int16 = 0
const KnightPhase int16 = 1
const BishopPhase int16 = 1
const RookPhase int16 = 2
const QueenPhase int16 = 4
const TotalPhase int16 = PawnPhase*16 + KnightPhase*4 + BishopPhase*4 + RookPhase*4 + QueenPhase*2
const HalfPhase = TotalPhase / 2
const Tempo int16 = 5

const BlackKingSideMask = uint64(1<<F8 | 1<<G8 | 1<<H8 | 1<<F7 | 1<<G7 | 1<<H7)
const WhiteKingSideMask = uint64(1<<F1 | 1<<G1 | 1<<H1 | 1<<F2 | 1<<G2 | 1<<H2)
const BlackQueenSideMask = uint64(1<<C8 | 1<<B8 | 1<<A8 | 1<<A7 | 1<<B7 | 1<<C7)
const WhiteQueenSideMask = uint64(1<<C1 | 1<<B1 | 1<<A1 | 1<<A2 | 1<<B2 | 1<<C2)

const BlackAShield = uint64(1<<A7 | 1<<A6)
const BlackBShield = uint64(1<<B7 | 1<<B6)
const BlackCShield = uint64(1<<C7 | 1<<C6)
const BlackFShield = uint64(1<<F7 | 1<<F6)
const BlackGShield = uint64(1<<G7 | 1<<G6)
const BlackHShield = uint64(1<<H7 | 1<<H6)
const WhiteAShield = uint64(1<<A2 | 1<<A3)
const WhiteBShield = uint64(1<<B2 | 1<<B3)
const WhiteCShield = uint64(1<<C2 | 1<<C3)
const WhiteFShield = uint64(1<<F2 | 1<<F3)
const WhiteGShield = uint64(1<<G2 | 1<<G3)
const WhiteHShield = uint64(1<<H2 | 1<<H3)

func PSQT(piece Piece, sq Square, isEndgame bool) int16 {
	if isEndgame {
		switch piece {
		case WhitePawn:
			return LatePawnPst[Flip[int(sq)]]
		case WhiteKnight:
			return LateKnightPst[Flip[int(sq)]]
		case WhiteBishop:
			return LateBishopPst[Flip[int(sq)]]
		case WhiteRook:
			return LateRookPst[Flip[int(sq)]]
		case WhiteQueen:
			return LateQueenPst[Flip[int(sq)]]
		case WhiteKing:
			return LateKingPst[Flip[int(sq)]]
		case BlackPawn:
			return LatePawnPst[int(sq)]
		case BlackKnight:
			return LateKnightPst[int(sq)]
		case BlackBishop:
			return LateBishopPst[int(sq)]
		case BlackRook:
			return LateRookPst[int(sq)]
		case BlackQueen:
			return LateQueenPst[int(sq)]
		case BlackKing:
			return LateKingPst[int(sq)]
		}
	} else {
		switch piece {
		case WhitePawn:
			return EarlyPawnPst[Flip[int(sq)]]
		case WhiteKnight:
			return EarlyKnightPst[Flip[int(sq)]]
		case WhiteBishop:
			return EarlyBishopPst[Flip[int(sq)]]
		case WhiteRook:
			return EarlyRookPst[Flip[int(sq)]]
		case WhiteQueen:
			return EarlyQueenPst[Flip[int(sq)]]
		case WhiteKing:
			return EarlyKingPst[Flip[int(sq)]]
		case BlackPawn:
			return EarlyPawnPst[int(sq)]
		case BlackKnight:
			return EarlyKnightPst[int(sq)]
		case BlackBishop:
			return EarlyBishopPst[int(sq)]
		case BlackRook:
			return EarlyRookPst[int(sq)]
		case BlackQueen:
			return EarlyQueenPst[int(sq)]
		case BlackKing:
			return EarlyKingPst[int(sq)]
		}
	}
	return 0
}

func Evaluate(position *Position, pawnhash *PawnCache) int16 {
	var drawDivider int16 = 0
	// position.MaterialAndPSQT()
	board := position.Board
	turn := position.Turn()

	// fetch material balance
	whitePawnsCount := position.MaterialsOnBoard[WhitePawn-1]
	whiteKnightsCount := position.MaterialsOnBoard[WhiteKnight-1]
	whiteBishopsCount := position.MaterialsOnBoard[WhiteBishop-1]
	whiteRooksCount := position.MaterialsOnBoard[WhiteRook-1]
	whiteQueensCount := position.MaterialsOnBoard[WhiteQueen-1]

	blackPawnsCount := position.MaterialsOnBoard[BlackPawn-1]
	blackKnightsCount := position.MaterialsOnBoard[BlackKnight-1]
	blackBishopsCount := position.MaterialsOnBoard[BlackBishop-1]
	blackRooksCount := position.MaterialsOnBoard[BlackRook-1]
	blackQueensCount := position.MaterialsOnBoard[BlackQueen-1]

	// Fetch PSQT Rewards/Penalty
	whiteCentipawnsMG := position.WhiteMiddlegamePSQT
	whiteCentipawnsEG := position.WhiteEndgamePSQT
	blackCentipawnsMG := position.BlackMiddlegamePSQT
	blackCentipawnsEG := position.BlackEndgamePSQT

	allPiecesCount :=
		whitePawnsCount +
			blackPawnsCount +
			whiteKnightsCount +
			blackKnightsCount +
			whiteBishopsCount +
			blackBishopsCount +
			whiteRooksCount +
			blackRooksCount +
			whiteQueensCount +
			blackQueensCount

	// KPK endgame, ask the bitbase
	if allPiecesCount == 1 && (whitePawnsCount == 1 || blackPawnsCount == 1) {
		if whitePawnsCount == 1 {
			return KpkProbe(board, White, turn)
		}
		return KpkProbe(board, Black, turn)
	}

	whites := board.GetWhitePieces()
	blacks := board.GetBlackPieces()
	all := whites | blacks

	bbBlackKing := board.GetBitboardOf(BlackKing)
	bbWhiteKing := board.GetBitboardOf(WhiteKing)

	bbBlackRook := board.GetBitboardOf(BlackRook)
	bbWhiteRook := board.GetBitboardOf(WhiteRook)

	bbBlackPawn := board.GetBitboardOf(BlackPawn)
	bbWhitePawn := board.GetBitboardOf(WhitePawn)

	whiteKingIndex := bits.TrailingZeros64(bbWhiteKing)
	blackKingIndex := bits.TrailingZeros64(bbBlackKing)

	// Double Rooks
	if blackRooksCount > 1 {
		sq := Square(bits.TrailingZeros64(bbBlackRook))
		if board.IsVerticalDoubleRook(sq, bbBlackRook, all) {
			// double-rook vertical
			blackCentipawnsEG += EndgameVeritcalDoubleRookAward
			blackCentipawnsMG += MiddlegameVeritcalDoubleRookAward
		} else if board.IsHorizontalDoubleRook(sq, bbBlackRook, all) {
			// double-rook horizontal
			blackCentipawnsMG += MiddlegameHorizontalDoubleRookAward
			blackCentipawnsEG += EndgameHorizontalDoubleRookAward
		}
	}

	if whiteRooksCount > 1 {
		sq := Square(bits.TrailingZeros64(bbWhiteRook))
		if board.IsVerticalDoubleRook(sq, bbWhiteRook, all) {
			// double-rook vertical
			whiteCentipawnsEG += EndgameVeritcalDoubleRookAward
			whiteCentipawnsMG += MiddlegameVeritcalDoubleRookAward
		} else if board.IsHorizontalDoubleRook(sq, bbWhiteRook, all) {
			// double-rook horizontal
			whiteCentipawnsMG += MiddlegameHorizontalDoubleRookAward
			whiteCentipawnsEG += EndgameHorizontalDoubleRookAward
		}
	}

	// Draw scenarios
	{

		if (allPiecesCount == 2 && whiteRooksCount == 1 && (blackKnightsCount == 1 || blackBishopsCount == 1)) ||
			(allPiecesCount == 2 && blackRooksCount == 1 && (whiteKnightsCount == 1 || whiteBishopsCount == 1)) ||
			(allPiecesCount == 2 && (blackKnightsCount == 1 || blackBishopsCount == 1) && whitePawnsCount == 1) ||
			(allPiecesCount == 2 && (whiteKnightsCount == 1 || whiteBishopsCount == 1) && blackPawnsCount == 1) ||
			(allPiecesCount == 3 && blackRooksCount == 1 && whiteRooksCount == 1 && (whiteKnightsCount == 1 || blackKnightsCount == 1 || blackBishopsCount == 1 || whiteBishopsCount == 1)) {
			drawDivider = 3
		}
	}

	pawnFactorMG := int16(16-blackPawnsCount-whitePawnsCount) * MiddlegamePawnFactorCoeff
	pawnFactorEG := int16(16-blackPawnsCount-whitePawnsCount) * EndgamePawnFactorCoeff

	blackCentipawnsMG += blackPawnsCount * BlackPawn.Weight()
	blackCentipawnsMG += blackKnightsCount * (BlackKnight.Weight() - pawnFactorMG)
	blackCentipawnsMG += blackBishopsCount * (BlackBishop.Weight())
	blackCentipawnsMG += blackRooksCount * (BlackRook.Weight() + pawnFactorMG)
	blackCentipawnsMG += blackQueensCount * BlackQueen.Weight()

	blackCentipawnsEG += blackPawnsCount * BlackPawn.Weight()
	blackCentipawnsEG += blackKnightsCount * (BlackKnight.Weight() - pawnFactorEG)
	blackCentipawnsEG += blackBishopsCount * (BlackBishop.Weight())
	blackCentipawnsEG += blackRooksCount * (BlackRook.Weight() + pawnFactorEG)
	blackCentipawnsEG += blackQueensCount * BlackQueen.Weight()

	whiteCentipawnsMG += whitePawnsCount * WhitePawn.Weight()
	whiteCentipawnsMG += whiteKnightsCount * (WhiteKnight.Weight() - pawnFactorMG)
	whiteCentipawnsMG += whiteBishopsCount * (WhiteBishop.Weight())
	whiteCentipawnsMG += whiteRooksCount * (WhiteRook.Weight() + pawnFactorMG)
	whiteCentipawnsMG += whiteQueensCount * WhiteQueen.Weight()

	whiteCentipawnsEG += whitePawnsCount * WhitePawn.Weight()
	whiteCentipawnsEG += whiteKnightsCount * (WhiteKnight.Weight() - pawnFactorEG)
	whiteCentipawnsEG += whiteBishopsCount * (WhiteBishop.Weight())
	whiteCentipawnsEG += whiteRooksCount * (WhiteRook.Weight() + pawnFactorEG)
	whiteCentipawnsEG += whiteQueensCount * WhiteQueen.Weight()

	// Bishop Pair
	if whiteBishopsCount >= 2 {
		whiteCentipawnsMG += MiddlegameBishopPairAward
		whiteCentipawnsEG += EndgameBishopPairAward
	}
	if blackBishopsCount >= 2 {
		blackCentipawnsMG += MiddlegameBishopPairAward
		blackCentipawnsEG += EndgameBishopPairAward
	}

	mobilityEval := Mobility(position, blackKingIndex, whiteKingIndex)

	whiteCentipawnsMG += mobilityEval.whiteMG
	whiteCentipawnsEG += mobilityEval.whiteEG
	blackCentipawnsMG += mobilityEval.blackMG
	blackCentipawnsEG += mobilityEval.blackEG

	rookEval := RookFilesEval(bbBlackRook, bbWhiteRook, bbBlackPawn, bbWhitePawn)
	whiteCentipawnsMG += rookEval.whiteMG
	whiteCentipawnsEG += rookEval.whiteEG
	blackCentipawnsMG += rookEval.blackMG
	blackCentipawnsEG += rookEval.blackEG

	pawnMG, pawnEG := CachedPawnStructureEval(position, pawnhash)

	kingSafetyEval := KingSafety(bbBlackKing, bbWhiteKing, bbBlackPawn, bbWhitePawn,
		position.HasTag(BlackCanCastleQueenSide) || position.HasTag(BlackCanCastleKingSide),
		position.HasTag(WhiteCanCastleQueenSide) || position.HasTag(WhiteCanCastleKingSide),
	)
	whiteCentipawnsMG += kingSafetyEval.whiteMG
	whiteCentipawnsEG += kingSafetyEval.whiteEG
	blackCentipawnsMG += kingSafetyEval.blackMG
	blackCentipawnsEG += kingSafetyEval.blackEG

	knightOutpostEval := KnightOutpostEval(position)
	whiteCentipawnsMG += knightOutpostEval.whiteMG
	whiteCentipawnsEG += knightOutpostEval.whiteEG
	blackCentipawnsMG += knightOutpostEval.blackMG
	blackCentipawnsEG += knightOutpostEval.blackEG

	phase := TotalPhase -
		whitePawnsCount*PawnPhase -
		blackPawnsCount*PawnPhase -
		whiteKnightsCount*KnightPhase -
		blackKnightsCount*KnightPhase -
		whiteBishopsCount*BishopPhase -
		blackBishopsCount*BishopPhase -
		whiteRooksCount*RookPhase -
		blackRooksCount*RookPhase -
		whiteQueensCount*QueenPhase -
		blackQueensCount*QueenPhase

	phase = (phase*256 + HalfPhase) / TotalPhase

	var evalEG, evalMG int16

	if turn == White {
		evalEG = whiteCentipawnsEG - blackCentipawnsEG + pawnEG
		evalMG = whiteCentipawnsMG - blackCentipawnsMG + pawnMG
	} else {
		evalEG = blackCentipawnsEG - whiteCentipawnsEG - pawnEG
		evalMG = blackCentipawnsMG - whiteCentipawnsMG - pawnMG
	}

	// The following formula overflows if I do not convert to int32 first
	// then I have to convert back to int16, as the function return requires
	// and that is also safe, due to the division
	mg := int32(evalMG)
	eg := int32(evalEG)
	phs := int32(phase)
	taperedEval := int16(((mg * (256 - phs)) + eg*phs) / 256)
	return toEval(taperedEval+Tempo) >> drawDivider
}

func KnightOutpostEval(p *Position) Eval {
	var blackMG, whiteMG, blackEG, whiteEG int16
	blackOutposts := p.CountKnightOutposts(Black)
	whiteOutposts := p.CountKnightOutposts(White)

	blackMG = MiddlegameKnightOutpostAward * blackOutposts
	blackEG = EndgameKnightOutpostAward * blackOutposts
	whiteMG = MiddlegameKnightOutpostAward * whiteOutposts
	whiteEG = EndgameKnightOutpostAward * whiteOutposts

	return Eval{blackMG: blackMG, whiteMG: whiteMG, blackEG: blackEG, whiteEG: whiteEG}
}

func RookFilesEval(blackRook uint64, whiteRook uint64, blackPawns uint64, whitePawns uint64) Eval {
	var blackMG, whiteMG, blackEG, whiteEG int16

	blackFiles := FileFill(blackRook)
	whiteFiles := FileFill(whiteRook)

	allPawns := FileFill(blackPawns | whitePawns)

	// open files
	blackRooksNoPawns := blackFiles &^ allPawns
	whiteRooksNoPawns := whiteFiles &^ allPawns

	blackRookOpenFiles := blackRook & blackRooksNoPawns
	whiteRookOpenFiles := whiteRook & whiteRooksNoPawns

	count := int16(bits.OnesCount64(blackRookOpenFiles))
	blackMG += MiddlegameRookOpenFileAward * count
	blackEG += EndgameRookOpenFileAward * count

	count = int16(bits.OnesCount64(whiteRookOpenFiles))
	whiteMG += MiddlegameRookOpenFileAward * count
	whiteEG += EndgameRookOpenFileAward * count

	// semi-open files
	blackRooksNoOwnPawns := blackFiles &^ FileFill(blackPawns)
	whiteRooksNoOwnPawns := whiteFiles &^ FileFill(whitePawns)

	blackRookSemiOpenFiles := (blackRook &^ blackRookOpenFiles) & blackRooksNoOwnPawns
	whiteRookSemiOpenFiles := (whiteRook &^ whiteRookOpenFiles) & whiteRooksNoOwnPawns

	count = int16(bits.OnesCount64(blackRookSemiOpenFiles))
	blackMG += MiddlegameRookSemiOpenFileAward * count
	blackEG += EndgameRookSemiOpenFileAward * count

	count = int16(bits.OnesCount64(whiteRookSemiOpenFiles))
	whiteMG += MiddlegameRookSemiOpenFileAward * count
	whiteEG += EndgameRookSemiOpenFileAward * count

	return Eval{blackMG: blackMG, whiteMG: whiteMG, blackEG: blackEG, whiteEG: whiteEG}
}

func CachedPawnStructureEval(p *Position, pawnhash *PawnCache) (int16, int16) {
	hash := p.Pawnhash()
	mg, eg, ok := pawnhash.Get(hash)

	if ok {
		return mg, eg
	}

	eval := PawnStructureEval(p)
	mg = eval.whiteMG - eval.blackMG
	eg = eval.whiteEG - eval.blackEG
	pawnhash.Set(hash, mg, eg)

	return mg, eg
}

func PawnStructureEval(p *Position) Eval {
	var blackMG, whiteMG, blackEG, whiteEG int16

	// passed pawns
	countP, countS := p.CountPassedPawns(Black)
	blackMG += MiddlegamePassedPawnAward * countP
	blackEG += EndgamePassedPawnAward * countP

	blackMG += MiddlegameAdvancedPassedPawnAward * countS
	blackEG += EndgameAdvancedPassedPawnAward * countS

	countP, countS = p.CountPassedPawns(White)
	whiteMG += MiddlegamePassedPawnAward * countP
	whiteEG += EndgamePassedPawnAward * countP

	whiteMG += MiddlegameAdvancedPassedPawnAward * countS
	whiteEG += EndgameAdvancedPassedPawnAward * countS

	// candidate passed pawns
	count := p.CountCandidatePawns(Black)
	blackMG += MiddlegameCandidatePassedPawnAward * count
	blackEG += EndgameCandidatePassedPawnAward * count

	count = p.CountCandidatePawns(White)
	whiteMG += MiddlegameCandidatePassedPawnAward * count
	whiteEG += EndgameCandidatePassedPawnAward * count

	// backward pawns
	count = p.CountBackwardPawns(Black)
	blackMG -= MiddlegameBackwardPawnPenalty * count
	blackEG -= EndgameBackwardPawnPenalty * count

	count = p.CountBackwardPawns(White)
	whiteMG -= MiddlegameBackwardPawnPenalty * count
	whiteEG -= EndgameBackwardPawnPenalty * count

	// isolated pawns
	count = p.CountIsolatedPawns(Black)
	blackMG -= MiddlegameIsolatedPawnPenalty * count
	blackEG -= EndgameIsolatedPawnPenalty * count

	count = p.CountIsolatedPawns(White)
	whiteMG -= MiddlegameIsolatedPawnPenalty * count
	whiteEG -= EndgameIsolatedPawnPenalty * count

	// double pawns
	count = p.CountDoublePawns(Black)
	blackMG -= MiddlegameDoublePawnPenalty * count
	blackEG -= EndgameDoublePawnPenalty * count

	count = p.CountDoublePawns(White)
	whiteMG -= MiddlegameDoublePawnPenalty * count
	whiteEG -= EndgameDoublePawnPenalty * count

	return Eval{blackMG: blackMG, whiteMG: whiteMG, blackEG: blackEG, whiteEG: whiteEG}
}

func kingSafetyPenalty(color Color, side PieceType, ownPawn uint64, allPawn uint64) (int16, int16) {
	var mg, eg int16
	var a_shield, b_shield, c_shield, f_shield, g_shield, h_shield uint64
	if color == White {
		a_shield = WhiteAShield
		b_shield = WhiteBShield
		c_shield = WhiteCShield
		f_shield = WhiteFShield
		g_shield = WhiteGShield
		h_shield = WhiteHShield
	} else {
		a_shield = BlackAShield
		b_shield = BlackBShield
		c_shield = BlackCShield
		f_shield = BlackFShield
		g_shield = BlackGShield
		h_shield = BlackHShield
	}
	if side == King {
		if H_FileFill&allPawn == 0 { // no pawns, super bad
			mg += MiddlegameKingZoneOpenFilePenalty
			eg += EndgameKingZoneOpenFilePenalty
		} else if H_FileFill&ownPawn == 0 { // semi-open file, bad
			mg += MiddlegameKingZoneMissingPawnPenalty
			eg += EndgameKingZoneMissingPawnPenalty
		} else if h_shield&ownPawn == 0 {
			mg += MiddlegamePawnShieldPenalty
			eg += EndgamePawnShieldPenalty
		}

		if G_FileFill&allPawn == 0 { // no pawns, super bad
			mg += MiddlegameKingZoneOpenFilePenalty
			eg += EndgameKingZoneOpenFilePenalty
		} else if G_FileFill&ownPawn == 0 { // semi-open file, bad
			mg += MiddlegameKingZoneMissingPawnPenalty
			eg += EndgameKingZoneMissingPawnPenalty
		} else if g_shield&ownPawn == 0 {
			mg += MiddlegamePawnShieldPenalty
			eg += EndgamePawnShieldPenalty
		}

		if F_FileFill&allPawn == 0 { // no pawns, super bad
			mg += MiddlegameKingZoneOpenFilePenalty
			eg += EndgameKingZoneOpenFilePenalty
		} else if F_FileFill&ownPawn == 0 { // semi-open file, bad
			mg += MiddlegameKingZoneMissingPawnPenalty
			eg += EndgameKingZoneMissingPawnPenalty
		} else if f_shield&ownPawn == 0 {
			mg += MiddlegamePawnShieldPenalty
			eg += EndgamePawnShieldPenalty
		}
	} else {
		if C_FileFill&allPawn == 0 { // no pawns, super bad
			mg += MiddlegameKingZoneOpenFilePenalty
			eg += EndgameKingZoneOpenFilePenalty
		} else if C_FileFill&ownPawn == 0 { // semi-open file, bad
			mg += MiddlegameKingZoneMissingPawnPenalty
			eg += EndgameKingZoneMissingPawnPenalty
		} else if c_shield&ownPawn == 0 {
			mg += MiddlegamePawnShieldPenalty
			eg += EndgamePawnShieldPenalty
		}

		if B_FileFill&allPawn == 0 { // no pawns, super bad
			mg += MiddlegameKingZoneOpenFilePenalty
			eg += EndgameKingZoneOpenFilePenalty
		} else if B_FileFill&ownPawn == 0 { // semi-open file, bad
			mg += MiddlegameKingZoneMissingPawnPenalty
			eg += EndgameKingZoneMissingPawnPenalty
		} else if b_shield&ownPawn == 0 {
			mg += MiddlegamePawnShieldPenalty
			eg += EndgamePawnShieldPenalty
		}

		if A_FileFill&allPawn == 0 { // no pawns, super bad
			mg += MiddlegameKingZoneOpenFilePenalty
			eg += EndgameKingZoneOpenFilePenalty
		} else if A_FileFill&ownPawn == 0 { // semi-open file, bad
			mg += MiddlegameKingZoneMissingPawnPenalty
			eg += EndgameKingZoneMissingPawnPenalty
		} else if a_shield&ownPawn == 0 {
			mg += MiddlegamePawnShieldPenalty
			eg += EndgamePawnShieldPenalty
		}
	}

	return mg, eg
}

func KingSafety(blackKing uint64, whiteKing uint64, blackPawn uint64,
	whitePawn uint64, blackCastleFlag bool, whiteCastleFlag bool) Eval {
	var whiteCentipawnsMG, whiteCentipawnsEG, blackCentipawnsMG, blackCentipawnsEG int16
	allPawn := whitePawn | blackPawn

	if blackKing&BlackKingSideMask != 0 {
		blackCastleFlag = true
		// Missing pawn shield
		mg, eg := kingSafetyPenalty(Black, King, blackPawn, allPawn)
		blackCentipawnsMG -= mg
		blackCentipawnsEG -= eg
	} else if blackKing&BlackQueenSideMask != 0 {
		blackCastleFlag = true
		// Missing pawn shield
		mg, eg := kingSafetyPenalty(Black, Queen, blackPawn, allPawn)
		blackCentipawnsMG -= mg
		blackCentipawnsEG -= eg
	}

	if whiteKing&WhiteKingSideMask != 0 {
		whiteCastleFlag = true
		// Missing pawn shield
		mg, eg := kingSafetyPenalty(White, King, whitePawn, allPawn)
		whiteCentipawnsMG -= mg
		whiteCentipawnsEG -= eg
	} else if whiteKing&WhiteQueenSideMask != 0 {
		whiteCastleFlag = true
		// Missing pawn shield
		mg, eg := kingSafetyPenalty(White, Queen, whitePawn, allPawn)
		whiteCentipawnsMG -= mg
		whiteCentipawnsEG -= eg
	}

	if !whiteCastleFlag {
		whiteCentipawnsMG -= MiddlegameNotCastlingPenalty
		whiteCentipawnsEG -= EndgameNotCastlingPenalty
	}

	if !blackCastleFlag {
		blackCentipawnsMG -= MiddlegameNotCastlingPenalty
		blackCentipawnsEG -= EndgameNotCastlingPenalty
	}
	return Eval{blackMG: blackCentipawnsMG, whiteMG: whiteCentipawnsMG, blackEG: blackCentipawnsEG, whiteEG: whiteCentipawnsEG}
}

func Mobility(p *Position, blackKingIndex int, whiteKingIndex int) Eval {
	board := p.Board
	var whiteCentipawnsMG, whiteCentipawnsEG, blackCentipawnsMG, blackCentipawnsEG int16

	// mobility and attacks
	whitePawnAttacks, whiteMinorAttacks, whiteMajorAttacks := board.AllAttacks(White) // get the squares that are attacked by white
	blackPawnAttacks, blackMinorAttacks, blackMajorAttacks := board.AllAttacks(Black) // get the squares that are attacked by black

	blackKingZone := SquareInnerRingMask[blackKingIndex] | SquareOuterRingMask[blackKingIndex]
	whiteKingZone := SquareInnerRingMask[whiteKingIndex] | SquareOuterRingMask[whiteKingIndex]

	// Pawn controlled squares
	wPawnAttacks := int16(bits.OnesCount64(whitePawnAttacks &^ blackKingZone))
	bPawnAttacks := int16(bits.OnesCount64(blackPawnAttacks &^ whiteKingZone))

	whiteCentipawnsMG += MiddlegamePawnSquareControlCoeff * wPawnAttacks
	whiteCentipawnsEG += EndgamePawnSquareControlCoeff * wPawnAttacks

	blackCentipawnsMG += MiddlegamePawnSquareControlCoeff * bPawnAttacks
	blackCentipawnsEG += EndgamePawnSquareControlCoeff * bPawnAttacks

	// // Minor mobility
	wMinorAttacksNoKingZone := whiteMinorAttacks &^ blackKingZone
	bMinorAttacksNoKingZone := blackMinorAttacks &^ whiteKingZone

	wMinorQuietAttacks := int16(bits.OnesCount64(wMinorAttacksNoKingZone << 32)) // keep hi-bits only
	bMinorQuietAttacks := int16(bits.OnesCount64(bMinorAttacksNoKingZone >> 32)) // keep lo-bits only

	wMinorAggressivity := int16(bits.OnesCount64(wMinorAttacksNoKingZone >> 32)) // keep hi-bits only
	bMinorAggressivity := int16(bits.OnesCount64(bMinorAttacksNoKingZone << 32)) // keep lo-bits only

	whiteCentipawnsMG += MiddlegameMinorMobilityFactorCoeff * wMinorQuietAttacks
	whiteCentipawnsEG += EndgameMinorMobilityFactorCoeff * wMinorQuietAttacks

	blackCentipawnsMG += MiddlegameMinorMobilityFactorCoeff * bMinorQuietAttacks
	blackCentipawnsEG += EndgameMinorMobilityFactorCoeff * bMinorQuietAttacks

	whiteCentipawnsMG += MiddlegameMinorAggressivityFactorCoeff * wMinorAggressivity
	whiteCentipawnsEG += EndgameMinorAggressivityFactorCoeff * wMinorAggressivity

	blackCentipawnsMG += MiddlegameMinorAggressivityFactorCoeff * bMinorAggressivity
	blackCentipawnsEG += EndgameMinorAggressivityFactorCoeff * bMinorAggressivity

	// Major mobility
	wMajorAttacksNoKingZone := whiteMajorAttacks &^ blackKingZone
	bMajorAttacksNoKingZone := blackMajorAttacks &^ whiteKingZone

	wMajorQuietAttacks := int16(bits.OnesCount64(wMajorAttacksNoKingZone << 32)) // keep hi-bits only
	bMajorQuietAttacks := int16(bits.OnesCount64(bMajorAttacksNoKingZone >> 32)) // keep lo-bits only

	wMajorAggressivity := int16(bits.OnesCount64(wMajorAttacksNoKingZone >> 32)) // keep hi-bits only
	bMajorAggressivity := int16(bits.OnesCount64(bMajorAttacksNoKingZone << 32)) // keep lo-bits only

	whiteCentipawnsMG += MiddlegameMajorMobilityFactorCoeff * wMajorQuietAttacks
	whiteCentipawnsEG += EndgameMajorMobilityFactorCoeff * wMajorQuietAttacks

	blackCentipawnsMG += MiddlegameMajorMobilityFactorCoeff * bMajorQuietAttacks
	blackCentipawnsEG += EndgameMajorMobilityFactorCoeff * bMajorQuietAttacks

	whiteCentipawnsMG += MiddlegameMajorAggressivityFactorCoeff * wMajorAggressivity
	whiteCentipawnsEG += EndgameMajorAggressivityFactorCoeff * wMajorAggressivity

	blackCentipawnsMG += MiddlegameMajorAggressivityFactorCoeff * bMajorAggressivity
	blackCentipawnsEG += EndgameMajorAggressivityFactorCoeff * bMajorAggressivity

	// King attacks
	whiteCentipawnsMG +=
		MiddlegameInnerPawnToKingAttackCoeff*int16(bits.OnesCount64(whitePawnAttacks&SquareInnerRingMask[blackKingIndex])) +
			MiddlegameOuterPawnToKingAttackCoeff*int16(bits.OnesCount64(whitePawnAttacks&SquareOuterRingMask[blackKingIndex])) +
			MiddlegameInnerMinorToKingAttackCoeff*int16(bits.OnesCount64(whiteMinorAttacks&SquareInnerRingMask[blackKingIndex])) +
			MiddlegameOuterMinorToKingAttackCoeff*int16(bits.OnesCount64(whiteMinorAttacks&SquareOuterRingMask[blackKingIndex])) +
			MiddlegameInnerMajorToKingAttackCoeff*int16(bits.OnesCount64(whiteMajorAttacks&SquareInnerRingMask[blackKingIndex])) +
			MiddlegameOuterMajorToKingAttackCoeff*int16(bits.OnesCount64(whiteMajorAttacks&SquareOuterRingMask[blackKingIndex]))

	whiteCentipawnsEG +=
		EndgameInnerPawnToKingAttackCoeff*int16(bits.OnesCount64(whitePawnAttacks&SquareInnerRingMask[blackKingIndex])) +
			EndgameOuterPawnToKingAttackCoeff*int16(bits.OnesCount64(whitePawnAttacks&SquareOuterRingMask[blackKingIndex])) +
			EndgameInnerMinorToKingAttackCoeff*int16(bits.OnesCount64(whiteMinorAttacks&SquareInnerRingMask[blackKingIndex])) +
			EndgameOuterMinorToKingAttackCoeff*int16(bits.OnesCount64(whiteMinorAttacks&SquareOuterRingMask[blackKingIndex])) +
			EndgameInnerMajorToKingAttackCoeff*int16(bits.OnesCount64(whiteMajorAttacks&SquareInnerRingMask[blackKingIndex])) +
			EndgameOuterMajorToKingAttackCoeff*int16(bits.OnesCount64(whiteMajorAttacks&SquareOuterRingMask[blackKingIndex]))

	blackCentipawnsMG +=
		MiddlegameInnerPawnToKingAttackCoeff*int16(bits.OnesCount64(blackPawnAttacks&SquareInnerRingMask[whiteKingIndex])) +
			MiddlegameOuterPawnToKingAttackCoeff*int16(bits.OnesCount64(blackPawnAttacks&SquareOuterRingMask[whiteKingIndex])) +
			MiddlegameInnerMinorToKingAttackCoeff*int16(bits.OnesCount64(blackMinorAttacks&SquareInnerRingMask[whiteKingIndex])) +
			MiddlegameOuterMinorToKingAttackCoeff*int16(bits.OnesCount64(blackMinorAttacks&SquareOuterRingMask[whiteKingIndex])) +
			MiddlegameInnerMajorToKingAttackCoeff*int16(bits.OnesCount64(blackMajorAttacks&SquareInnerRingMask[whiteKingIndex])) +
			MiddlegameOuterMajorToKingAttackCoeff*int16(bits.OnesCount64(blackMajorAttacks&SquareOuterRingMask[whiteKingIndex]))

	blackCentipawnsEG +=
		EndgameInnerPawnToKingAttackCoeff*int16(bits.OnesCount64(blackPawnAttacks&SquareInnerRingMask[whiteKingIndex])) +
			EndgameOuterPawnToKingAttackCoeff*int16(bits.OnesCount64(blackPawnAttacks&SquareOuterRingMask[whiteKingIndex])) +
			EndgameInnerMinorToKingAttackCoeff*int16(bits.OnesCount64(blackMinorAttacks&SquareInnerRingMask[whiteKingIndex])) +
			EndgameOuterMinorToKingAttackCoeff*int16(bits.OnesCount64(blackMinorAttacks&SquareOuterRingMask[whiteKingIndex])) +
			EndgameInnerMajorToKingAttackCoeff*int16(bits.OnesCount64(blackMajorAttacks&SquareInnerRingMask[whiteKingIndex])) +
			EndgameOuterMajorToKingAttackCoeff*int16(bits.OnesCount64(blackMajorAttacks&SquareOuterRingMask[whiteKingIndex]))

	return Eval{blackMG: blackCentipawnsMG, whiteMG: whiteCentipawnsMG, blackEG: blackCentipawnsEG, whiteEG: whiteCentipawnsEG}
}

func toEval(eval int16) int16 {
	if eval >= CHECKMATE_EVAL {
		return MAX_NON_CHECKMATE
	} else if eval <= -CHECKMATE_EVAL {
		return -MAX_NON_CHECKMATE
	}
	return eval
}

func min16(x int16, y int16) int16 {
	if x < y {
		return x
	}
	return y
}
