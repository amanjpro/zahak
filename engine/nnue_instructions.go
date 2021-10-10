//go:build !amd64
// +build !amd64

package engine

func (n *NetworkState) QuickFeed() int16 {
	// apply output layer
	output := int16(0)
	hiddenOutputs := n.HiddenOutputs[n.CurrentHidden]
	for i := 0; i < len(n.OutputWeights); i++ {
		output += ReLu(hiddenOutputs[i]) * n.OutputWeights[i]
	}
	output += n.OutputBias * QPrecisionIn
	return output / QPrecisionIn / QPrecisionOut
}

func (n *NetworkState) UpdateHidden(updates *Updates) {
	n.CurrentHidden += 1
	hiddenOutputs := n.HiddenOutputs[n.CurrentHidden]
	copy(hiddenOutputs, n.HiddenOutputs[n.CurrentHidden-1])

	for i := 0; i < updates.Size; i++ {
		weights := n.HiddenWeights
		for j := 0; j < len(hiddenOutputs); j++ {
			hiddenOutputs[j] += int16(updates.Coeffs[i]) * weights[int(updates.Indices[i])*NetHiddenSize+j]
		}
	}
}
