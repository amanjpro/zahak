package engine

import (
	"encoding/binary"
	"io"
	"math"
	"os"
)

const NetInputSize = 768
const NetHiddenSize = 128
const NetOutputSize = 1
const NetLayers = 1
const MaximumDepth = 128

var CurrentHiddenWeights []float32
var CurrentHiddenBiases []float32
var CurrentOutputWeights []float32
var CurrentOutputBias float32
var CurrentNetworkId uint32

type Updates struct {
	Diff []Update
	Size int
}

func (u *Updates) Add(update Update) {
	u.Diff[u.Size] = update
	u.Size += 1
}

type NetworkState struct {
	HiddenOutputs     [][]float32
	EmptyHiddenOutput []float32
	CurrentHidden     int
	HiddenWeights     []float32
	HiddenBiases      []float32
	OutputWeights     []float32
	OutputBias        float32
}

func NewNetworkState() *NetworkState {
	net := NetworkState{
		HiddenWeights: make([]float32, NetInputSize*NetHiddenSize),
		HiddenBiases:  make([]float32, NetHiddenSize),
		OutputWeights: make([]float32, NetHiddenSize),
		OutputBias:    CurrentOutputBias,
	}

	copy(net.HiddenWeights, CurrentHiddenWeights)
	copy(net.HiddenBiases, CurrentHiddenBiases)
	copy(net.OutputWeights, CurrentOutputWeights)
	net.EmptyHiddenOutput = make([]float32, NetHiddenSize)
	net.HiddenOutputs = make([][]float32, MaximumDepth)
	for i := 0; i < MaximumDepth; i++ {
		net.HiddenOutputs[i] = make([]float32, NetHiddenSize)
	}
	return &net
}

type Change int8

const (
	Remove Change = -1
	Add    Change = 1
)

type Update struct {
	Index int16
	Value Change
}

func calculateNetInputIndex(sq Square, piece Piece) int16 {
	return int16(piece-1)*64 + int16(sq)
}

func (n *NetworkState) RevertHidden() {
	n.CurrentHidden -= 1
}

func (n *NetworkState) UpdateHidden(updates *Updates) {
	n.CurrentHidden += 1
	hiddenOutputs := n.HiddenOutputs[n.CurrentHidden]
	copy(hiddenOutputs, n.HiddenOutputs[n.CurrentHidden-1])

	for i := 0; i < updates.Size; i++ {
		d := updates.Diff[i]
		weights := n.HiddenWeights
		for j := 0; j < len(hiddenOutputs); j++ {
			hiddenOutputs[j] += float32(d.Value) * weights[int(d.Index)*NetHiddenSize+j]
		}
	}
}

func (n *NetworkState) Recalculate(input []int16) {
	n.CurrentHidden = 0
	// apply hidden layer
	hiddenOutputs := n.HiddenOutputs[n.CurrentHidden]
	copy(hiddenOutputs, n.EmptyHiddenOutput)

	for index := 0; index < len(input); index++ {
		i := int(input[index])
		weights := n.HiddenWeights
		for j := 0; j < len(hiddenOutputs); j++ {
			hiddenOutputs[j] += weights[i*NetHiddenSize+j]
		}
	}
	for i := 0; i < len(hiddenOutputs); i++ {
		hiddenOutputs[i] += n.HiddenBiases[i]
	}
}

func (n *NetworkState) QuickFeed() float32 {
	// apply output layer
	output := float32(0)
	hiddenOutputs := n.HiddenOutputs[n.CurrentHidden]
	for i := 0; i < len(n.OutputWeights); i++ {
		output += ReLu(hiddenOutputs[i]) * n.OutputWeights[i]
	}
	return output + n.OutputBias
}

// load a neural network from file
func LoadNetwork(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Read headers
	buf := make([]byte, 4)
	_, err = io.ReadFull(f, buf)
	if err != nil {
		panic(err)
	}
	// if buf[0] != 66 || buf[1] != 90 || buf[2] != 1 || buf[3] != 0 {
	// 	panic("Magic word does not match expected, exiting")
	// }

	_, err = io.ReadFull(f, buf)
	if err != nil {
		panic(err)
	}
	id := binary.LittleEndian.Uint32(buf)

	// Read Topology Header
	buf = make([]byte, 12)
	_, err = io.ReadFull(f, buf)
	if err != nil {
		panic(err)
	}
	// inputs := binary.LittleEndian.Uint32(buf[:4])
	// outputs := binary.LittleEndian.Uint32(buf[4:8])
	// layers := binary.LittleEndian.Uint32(buf[8:])

	// if inputs != NetInputSize || outputs != NetOutputSize || layers != NetLayers {
	// 	panic("Topology is not supported, exiting")
	// }

	buf = make([]byte, 4)
	_, err = io.ReadFull(f, buf)
	if err != nil {
		panic(err)
	}
	// neurons := binary.LittleEndian.Uint32(buf)
	// if neurons != NetHiddenSize {
	// 	panic("Topology is not supported, exiting")
	// }

	CurrentNetworkId = id

	CurrentHiddenWeights = make([]float32, NetHiddenSize*NetInputSize)
	for i := uint32(0); i < NetHiddenSize*NetInputSize; i++ {
		_, err := io.ReadFull(f, buf)
		if err != nil {
			panic(err)
		}
		CurrentHiddenWeights[i] = math.Float32frombits(binary.LittleEndian.Uint32(buf))
	}

	CurrentHiddenBiases = make([]float32, NetHiddenSize)
	for i := uint32(0); i < NetHiddenSize; i++ {
		_, err := io.ReadFull(f, buf)
		if err != nil {
			panic(err)
		}
		CurrentHiddenBiases[i] = math.Float32frombits(binary.LittleEndian.Uint32(buf))
	}

	CurrentOutputWeights = make([]float32, NetHiddenSize)
	for i := uint32(0); i < NetOutputSize*NetHiddenSize; i++ {
		_, err := io.ReadFull(f, buf)
		if err != nil {
			panic(err)
		}
		CurrentOutputWeights[i] = math.Float32frombits(binary.LittleEndian.Uint32(buf))
	}

	_, err = io.ReadFull(f, buf)
	if err != nil {
		panic(err)
	}
	CurrentOutputBias = math.Float32frombits(binary.LittleEndian.Uint32(buf))
}

func ReLu(x float32) float32 {
	if x < 0 {
		return 0
	}
	return x
}
