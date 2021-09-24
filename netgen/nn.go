//+build ignore

package main

import (
	"fmt"
	"os"

	. "github.com/amanjpro/zahak/engine"
)

var netPath = "dev"

func main() {
	LoadNetwork(netPath)

	f, err := os.Create(fmt.Sprintf("engine%cnn.go", os.PathSeparator))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("package engine\n\n")
	f.WriteString("// Code generated by go generate; DO NOT EDIT.\n\n")
	f.WriteString("func init() {\n")

	f.WriteString("HiddenWeights = []float32 {\n")
	for i := 0; i < len(HiddenWeights); i++ {
		if i%10 == 0 {
			f.WriteString("\n")
		}
		f.WriteString(fmt.Sprintf("%f,", HiddenWeights[i]))
	}
	f.WriteString("\n}\n")

	f.WriteString("HiddenBiases = []float32 {\n")
	for i := 0; i < len(HiddenBiases); i++ {
		if i%10 == 0 {
			f.WriteString("\n")
		}
		f.WriteString(fmt.Sprintf("%f,", HiddenBiases[i]))
	}
	f.WriteString("\n}\n")

	f.WriteString("OutputWeights = []float32 {\n")
	for i := 0; i < len(OutputWeights); i++ {
		if i%10 == 0 {
			f.WriteString("\n")
		}
		f.WriteString(fmt.Sprintf("%f,", OutputWeights[i]))
	}
	f.WriteString("\n}\n")

	f.WriteString(fmt.Sprintf("OutputBias = %f\n", OutputBias))
	f.WriteString(fmt.Sprintf("NetworkId = %d\n", NetworkId))

	f.WriteString("}\n")
}
