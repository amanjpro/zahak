//go:build !cgo
// +build !cgo

package fathom

import (
	. "github.com/amanjpro/zahak/engine"
)

const DefaultProbeDepth = 0

var MinProbeDepth int8 = DefaultProbeDepth

const (
	TB_LOSS          uint32 = 0 /* LOSS */
	TB_DRAW          uint32 = 2 /* DRAW */
	TB_WIN           uint32 = 4 /* WIN  */
	TB_RESULT_FAILED uint32 = 0xFFFFFFFF
)

func SetSyzygyPath(path string) {
}

func ClearSyzygy() {
}

func ProbeWDL(pos *Position, depth int8) uint32 {
	return TB_RESULT_FAILED
}

func ProbeDTZ(pos *Position) Move {
	return EmptyMove
}
