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
const CostEvalWeight float32 = 0.5
const CostWDLWeight float32 = 1.0 - CostEvalWeight
const SigmoidScale float32 = 2.5 / 1024
const MaximumDepth = 128

var HiddenWeights []float32
var HiddenBiases []float32
var OutputWeights []float32
var OutputBias float32
var NetworkId uint32

type Updates struct {
	Diff [10]Update
	Size int
}

func (u *Updates) Add(update Update) {
	u.Diff[u.Size] = update
	u.Size += 1
}

type NetworkState struct {
	HiddenOutputs [][]float32
	CurrentHidden int
	Output        float32
}

func NewNetworkState() *NetworkState {
	net := NetworkState{}
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

func (n *NetworkState) Recalculate(input []int16) {
	n.CurrentHidden = 0
	n.FeedInput(input)
}

func calculateNetInputIndex(sq Square, piece Piece) int16 {
	return int16(piece)*64 + int16(sq)
}

func (n *NetworkState) RevertHidden() {
	n.CurrentHidden -= 1
}

func (n *NetworkState) UpdateHidden(updates *Updates) {
	n.CurrentHidden += 1
	hiddenOutputs := n.HiddenOutputs[n.CurrentHidden]
	pastHiddenOutputs := n.HiddenOutputs[n.CurrentHidden-1]
	for i := 0; i < NetHiddenSize; i++ {
		hiddenOutputs[i] = pastHiddenOutputs[i]
	}

	for i := 0; i < updates.Size; i++ {
		d := updates.Diff[i]
		for j := 0; j < NetHiddenSize; j++ {
			hiddenOutputs[j] += float32(d.Value) * HiddenWeights[int(d.Index)*NetHiddenSize+j]
		}
	}
}

func (n *NetworkState) FeedInput(input []int16) {

	// apply hidden layer
	hiddenOutputs := n.HiddenOutputs[n.CurrentHidden]
	for i := 0; i < NetHiddenSize; i++ {
		hiddenOutputs[i] = 0
	}
	for index := 0; index < len(input); index++ {
		i := int(input[index])
		for j := 0; j < NetHiddenSize; j++ {
			hiddenOutputs[j] += HiddenWeights[i*NetHiddenSize+j]
		}
	}
	for i := 0; i < NetHiddenSize; i++ {
		hiddenOutputs[i] = hiddenOutputs[i] + HiddenBiases[i]
	}

	n.QuickFeed()
}

func (n *NetworkState) QuickFeed() {
	// apply output layer
	output := float32(0)
	hiddenOutputs := n.HiddenOutputs[n.CurrentHidden]
	for i := 0; i < NetHiddenSize; i++ {
		output += ReLu(hiddenOutputs[i]) * OutputWeights[i]
	}
	output = output + OutputBias

	n.Output = output
}

func (n *NetworkState) Evaluate(input []int16) float32 {
	n.FeedInput(input)
	return n.Output
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
	if buf[0] != 66 || buf[1] != 90 || buf[2] != 1 || buf[3] != 0 {
		panic("Magic word does not match expected, exiting")
	}

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
	inputs := binary.LittleEndian.Uint32(buf[:4])
	outputs := binary.LittleEndian.Uint32(buf[4:8])
	layers := binary.LittleEndian.Uint32(buf[8:])

	if inputs != NetInputSize || outputs != NetOutputSize || layers != NetLayers {
		panic("Topology is not supported, exiting")
	}

	buf = make([]byte, 4)
	_, err = io.ReadFull(f, buf)
	if err != nil {
		panic(err)
	}
	neurons := binary.LittleEndian.Uint32(buf)
	if neurons != NetHiddenSize {
		panic("Topology is not supported, exiting")
	}

	NetworkId = id

	HiddenWeights = make([]float32, NetHiddenSize*NetInputSize)
	for j := 0; j < len(HiddenWeights); j++ {
		_, err := io.ReadFull(f, buf)
		if err != nil {
			panic(err)
		}
		HiddenWeights[j] = math.Float32frombits(binary.LittleEndian.Uint32(buf))
	}

	HiddenBiases = make([]float32, NetHiddenSize)
	for j := 0; j < len(HiddenBiases); j++ {
		_, err := io.ReadFull(f, buf)
		if err != nil {
			panic(err)
		}
		HiddenBiases[j] = math.Float32frombits(binary.LittleEndian.Uint32(buf))
	}

	OutputWeights = make([]float32, NetHiddenSize*NetOutputSize)
	for j := 0; j < len(OutputWeights); j++ {
		_, err := io.ReadFull(f, buf)
		if err != nil {
			panic(err)
		}
		OutputWeights[j] = math.Float32frombits(binary.LittleEndian.Uint32(buf))
	}

	_, err = io.ReadFull(f, buf)
	if err != nil {
		panic(err)
	}
	OutputBias = math.Float32frombits(binary.LittleEndian.Uint32(buf))
}

func ReLu(x float32) float32 {
	if x < 0 {
		return 0
	}
	return x
}
