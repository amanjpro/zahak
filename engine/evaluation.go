package engine

const Rank2Fill = uint64(1<<A2 | 1<<B2 | 1<<C2 | 1<<D2 | 1<<E2 | 1<<F2 | 1<<G2 | 1<<H2)
const Rank7Fill = uint64(1<<A7 | 1<<B7 | 1<<C7 | 1<<D7 | 1<<E7 | 1<<F7 | 1<<G7 | 1<<H7)

func (p *Position) Evaluate() int16 {
	output := p.Net.QuickFeed()
	var eval int
	if p.Turn() == Black {
		eval = -toEval(output)
	} else {
		eval = toEval(output)
	}
	return int16(eval * (128 + p.phase) / (128 + maxPhase))
}

func toEval(eval float32) int {
	if eval >= MAX_NON_CHECKMATE {
		return int(MAX_NON_CHECKMATE)
	} else if eval <= MIN_NON_CHECKMATE {
		return int(MIN_NON_CHECKMATE)
	}
	return int(eval)
}

func abs16(x int16) int16 {
	if x >= 0 {
		return x
	}
	return -x
}
