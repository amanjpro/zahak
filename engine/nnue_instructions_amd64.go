//go:build amd64
// +build amd64

package engine

import "unsafe"

//go:noescape
func _update_hidden(previous_outputs, update_indices, update_coeffs, update_size, weights, outputs, outputs_len unsafe.Pointer)

//go:noescape
func _quick_feed(hidden_outputs, hidden_outputs_len, weights, weights_len, res unsafe.Pointer)

func (n *NetworkState) UpdateHidden(updates *Updates) {
	n.CurrentHidden += 1

	p1 := unsafe.Pointer(&n.HiddenOutputs[n.CurrentHidden-1][0])
	p2 := unsafe.Pointer(&updates.Indices[0])
	p3 := unsafe.Pointer(&updates.Coeffs[0])
	p4 := unsafe.Pointer(uintptr(updates.Size))
	p5 := unsafe.Pointer(&n.HiddenWeights[0])
	p6 := unsafe.Pointer(&n.HiddenOutputs[n.CurrentHidden][0])
	p7 := unsafe.Pointer(uintptr(NetHiddenSize))

	_update_hidden(p1, p2, p3, p4, p5, p6, p7)
}

func (n *NetworkState) QuickFeed() int16 {
	p1 := unsafe.Pointer(&n.HiddenOutputs[n.CurrentHidden][0])
	p2 := unsafe.Pointer(uintptr(NetHiddenSize))
	p3 := unsafe.Pointer(&n.OutputWeights[0])
	p4 := unsafe.Pointer(uintptr(NetHiddenSize))
	var res int32

	_quick_feed(p1, p2, p3, p4, unsafe.Pointer(&res))
	res += n.OutputBias
	return int16(res) / QPrecision / QPrecision
}
