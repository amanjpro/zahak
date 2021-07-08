package tuning

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"

	// "strconv"
	"strings"
	// "sync"

	. "github.com/amanjpro/zahak/engine"
	. "github.com/amanjpro/zahak/evaluation"
	. "github.com/amanjpro/zahak/search"
)

type TestPosition struct {
	pos     *Position
	outcome float64
}

var testPositions []TestPosition
var initialGuesses = computeInitialGuesses()
var K_PRECISION = 10
var NUM_PROCESSORS = 16
var initialK = 1.0
var skipParams map[int]bool
var answers = make(chan float64)
var ml = NewMoveList(500)

func initEngines() []*Engine {
	res := make([]*Engine, NUM_PROCESSORS)
	for i := 0; i < NUM_PROCESSORS; i++ {
		res[i] = NewEngine(nil)
	}
	return res
}

var futileIndices map[int]bool = initFutileIndices()

func initFutileIndices() map[int]bool {
	futileIndices := []int{
		0, 1, 2, 3, 4, 5, 6, 7,
		56, 57, 58, 59, 60, 61, 62, 63,
		64, 65, 66, 67, 68, 69, 70, 71,
		120, 121, 122, 123, 124, 125, 126, 127,
	}

	res := make(map[int]bool, len(futileIndices))

	for _, v := range futileIndices {
		res[v] = true
	}

	return res
}

func computeInitialGuesses() []int16 {
	var guesses = make([]int16, 0, 800)
	guesses = append(guesses, EarlyPawnPst[:]...)                    // 0-63
	guesses = append(guesses, LatePawnPst[:]...)                     // 64-127
	guesses = append(guesses, EarlyKnightPst[:]...)                  // 128-191
	guesses = append(guesses, LateKnightPst[:]...)                   // 192-255
	guesses = append(guesses, EarlyBishopPst[:]...)                  // 256-319
	guesses = append(guesses, LateBishopPst[:]...)                   // 320-383
	guesses = append(guesses, EarlyRookPst[:]...)                    // 384-447
	guesses = append(guesses, LateRookPst[:]...)                     // 448-511
	guesses = append(guesses, EarlyQueenPst[:]...)                   // 512-575
	guesses = append(guesses, LateQueenPst[:]...)                    // 576-639
	guesses = append(guesses, EarlyKingPst[:]...)                    // 640-703
	guesses = append(guesses, LateKingPst[:]...)                     // 704-767
	guesses = append(guesses, MiddlegameBackwardPawnPenalty)         // 768
	guesses = append(guesses, EndgameBackwardPawnPenalty)            // 769
	guesses = append(guesses, MiddlegameIsolatedPawnPenalty)         // 770
	guesses = append(guesses, EndgameIsolatedPawnPenalty)            // 771
	guesses = append(guesses, MiddlegameDoublePawnPenalty)           // 772
	guesses = append(guesses, EndgameDoublePawnPenalty)              // 773
	guesses = append(guesses, MiddlegamePassedPawnAward)             // 774
	guesses = append(guesses, EndgamePassedPawnAward)                // 775
	guesses = append(guesses, MiddlegameAdvancedPassedPawnAward)     // 776
	guesses = append(guesses, EndgameAdvancedPassedPawnAward)        // 777
	guesses = append(guesses, MiddlegameCandidatePassedPawnAward)    // 778
	guesses = append(guesses, EndgameCandidatePassedPawnAward)       // 779
	guesses = append(guesses, MiddlegameRookOpenFileAward)           // 780
	guesses = append(guesses, EndgameRookOpenFileAward)              // 781
	guesses = append(guesses, MiddlegameRookSemiOpenFileAward)       // 782
	guesses = append(guesses, EndgameRookSemiOpenFileAward)          // 783
	guesses = append(guesses, MiddlegameVeritcalDoubleRookAward)     // 784
	guesses = append(guesses, EndgameVeritcalDoubleRookAward)        // 785
	guesses = append(guesses, MiddlegameHorizontalDoubleRookAward)   // 786
	guesses = append(guesses, EndgameHorizontalDoubleRookAward)      // 787
	guesses = append(guesses, MiddlegamePawnFactorCoeff)             // 788
	guesses = append(guesses, EndgamePawnFactorCoeff)                // 789
	guesses = append(guesses, MiddlegameMobilityFactorCoeff)         // 790
	guesses = append(guesses, EndgameMobilityFactorCoeff)            // 791
	guesses = append(guesses, MiddlegameAggressivityFactorCoeff)     // 792
	guesses = append(guesses, EndgameAggressivityFactorCoeff)        // 793
	guesses = append(guesses, MiddlegameInnerPawnToKingAttackCoeff)  // 794
	guesses = append(guesses, EndgameInnerPawnToKingAttackCoeff)     // 795
	guesses = append(guesses, MiddlegameOuterPawnToKingAttackCoeff)  // 796
	guesses = append(guesses, EndgameOuterPawnToKingAttackCoeff)     // 797
	guesses = append(guesses, MiddlegameInnerMinorToKingAttackCoeff) // 798
	guesses = append(guesses, EndgameInnerMinorToKingAttackCoeff)    // 799
	guesses = append(guesses, MiddlegameOuterMinorToKingAttackCoeff) // 800
	guesses = append(guesses, EndgameOuterMinorToKingAttackCoeff)    // 801
	guesses = append(guesses, MiddlegameInnerMajorToKingAttackCoeff) // 802
	guesses = append(guesses, EndgameInnerMajorToKingAttackCoeff)    // 803
	guesses = append(guesses, MiddlegameOuterMajorToKingAttackCoeff) // 804
	guesses = append(guesses, EndgameOuterMajorToKingAttackCoeff)    // 805
	guesses = append(guesses, MiddlegamePawnShieldPenalty)           // 806
	guesses = append(guesses, EndgamePawnShieldPenalty)              // 807
	guesses = append(guesses, MiddlegameNotCastlingPenalty)          // 808
	guesses = append(guesses, EndgameNotCastlingPenalty)             // 809
	guesses = append(guesses, MiddlegameKingZoneOpenFilePenalty)     // 810
	guesses = append(guesses, EndgameKingZoneOpenFilePenalty)        // 811
	guesses = append(guesses, MiddlegameKingZoneMissingPawnPenalty)  // 812
	guesses = append(guesses, EndgameKingZoneMissingPawnPenalty)     // 813
	guesses = append(guesses, MiddlegameKnightOutpostAward)          // 814
	guesses = append(guesses, EndgameKnightOutpostAward)             // 815
	guesses = append(guesses, MiddlegameBishopPairAward)             // 816
	guesses = append(guesses, EndgameBishopPairAward)                // 817

	return guesses
}

func updateEvalParams(guesses []int16) {
	for i := 0; i < 64; i++ {
		EarlyPawnPst[i] = guesses[i+0*64]
		LatePawnPst[i] = guesses[i+1*64]
		EarlyKnightPst[i] = guesses[i+2*64]
		LateKnightPst[i] = guesses[i+3*64]
		EarlyBishopPst[i] = guesses[i+4*64]
		LateBishopPst[i] = guesses[i+5*64]
		EarlyRookPst[i] = guesses[i+6*64]
		LateRookPst[i] = guesses[i+7*64]
		EarlyQueenPst[i] = guesses[i+8*64]
		LateQueenPst[i] = guesses[i+9*64]
		EarlyKingPst[i] = guesses[i+10*64]
		LateKingPst[i] = guesses[i+11*64]
	}
	MiddlegameBackwardPawnPenalty = guesses[768]
	EndgameBackwardPawnPenalty = guesses[769]
	MiddlegameIsolatedPawnPenalty = guesses[770]
	EndgameIsolatedPawnPenalty = guesses[771]
	MiddlegameDoublePawnPenalty = guesses[772]
	EndgameDoublePawnPenalty = guesses[773]
	MiddlegamePassedPawnAward = guesses[774]
	EndgamePassedPawnAward = guesses[775]
	MiddlegameAdvancedPassedPawnAward = guesses[776]
	EndgameAdvancedPassedPawnAward = guesses[777]
	MiddlegameCandidatePassedPawnAward = guesses[778]
	EndgameCandidatePassedPawnAward = guesses[779]
	MiddlegameRookOpenFileAward = guesses[780]
	EndgameRookOpenFileAward = guesses[781]
	MiddlegameRookSemiOpenFileAward = guesses[782]
	EndgameRookSemiOpenFileAward = guesses[783]
	MiddlegameVeritcalDoubleRookAward = guesses[784]
	EndgameVeritcalDoubleRookAward = guesses[785]
	MiddlegameHorizontalDoubleRookAward = guesses[786]
	EndgameHorizontalDoubleRookAward = guesses[787]
	MiddlegamePawnFactorCoeff = guesses[788]
	EndgamePawnFactorCoeff = guesses[789]
	MiddlegameMobilityFactorCoeff = guesses[790]
	EndgameMobilityFactorCoeff = guesses[791]
	MiddlegameAggressivityFactorCoeff = guesses[792]
	EndgameAggressivityFactorCoeff = guesses[793]
	MiddlegameInnerPawnToKingAttackCoeff = guesses[794]
	EndgameInnerPawnToKingAttackCoeff = guesses[795]
	MiddlegameOuterPawnToKingAttackCoeff = guesses[796]
	EndgameOuterPawnToKingAttackCoeff = guesses[797]
	MiddlegameInnerMinorToKingAttackCoeff = guesses[798]
	EndgameInnerMinorToKingAttackCoeff = guesses[799]
	MiddlegameOuterMinorToKingAttackCoeff = guesses[800]
	EndgameOuterMinorToKingAttackCoeff = guesses[801]
	MiddlegameInnerMajorToKingAttackCoeff = guesses[802]
	EndgameInnerMajorToKingAttackCoeff = guesses[803]
	MiddlegameOuterMajorToKingAttackCoeff = guesses[804]
	EndgameOuterMajorToKingAttackCoeff = guesses[805]
	MiddlegamePawnShieldPenalty = guesses[806]
	EndgamePawnShieldPenalty = guesses[807]
	MiddlegameNotCastlingPenalty = guesses[808]
	EndgameNotCastlingPenalty = guesses[809]
	MiddlegameKingZoneOpenFilePenalty = guesses[810]
	EndgameKingZoneOpenFilePenalty = guesses[811]
	MiddlegameKingZoneMissingPawnPenalty = guesses[812]
	EndgameKingZoneMissingPawnPenalty = guesses[813]
	MiddlegameKnightOutpostAward = guesses[814]
	EndgameKnightOutpostAward = guesses[815]
	MiddlegameBishopPairAward = guesses[816]
	EndgameBishopPairAward = guesses[817]
}

func toEvalParams(guesses []float64) []int16 {
	params := make([]int16, len(guesses))
	for i, v := range guesses {
		params[i] = int16(v)
	}
	return params
}

func printPst(pst []int16, varname string) {
	acc := ""
	for i := 0; i < 64; i++ {
		if i%8 == 0 {
			acc = fmt.Sprintf("%s\n", acc)
		}
		acc = fmt.Sprintf("%s %d,", acc, pst[i])
	}
	fmt.Printf("var %s = [64]int16 { %s \n}\n\n", varname, acc)
}

func printOptimalGuesses(guesses []int16) {
	fmt.Println("// Middle-game")
	printPst(guesses[0:64], "EarlyPawnPst")
	printPst(guesses[128:192], "EarlyKnightPst")
	printPst(guesses[256:320], "EarlyBishopPst")
	printPst(guesses[384:448], "EarlyRookPst")
	printPst(guesses[512:576], "EarlyQueenPst")
	printPst(guesses[640:704], "EarlyKingPst")
	fmt.Println("// Endgame")
	printPst(guesses[64:128], "LatePawnPst")
	printPst(guesses[192:256], "LateKnightPst")
	printPst(guesses[320:384], "LateBishopPst")
	printPst(guesses[448:512], "LateRookPst")
	printPst(guesses[576:640], "LateQueenPst")
	printPst(guesses[704:768], "LateKingPst")

	fmt.Printf("var MiddlegameBackwardPawnPenalty int16 = %d\n", guesses[768])
	fmt.Printf("var EndgameBackwardPawnPenalty int16 = %d\n", guesses[769])
	fmt.Printf("var MiddlegameIsolatedPawnPenalty int16 = %d\n", guesses[770])
	fmt.Printf("var EndgameIsolatedPawnPenalty int16 = %d\n", guesses[771])
	fmt.Printf("var MiddlegameDoublePawnPenalty int16 = %d\n", guesses[772])
	fmt.Printf("var EndgameDoublePawnPenalty int16 = %d\n", guesses[773])
	fmt.Printf("var MiddlegamePassedPawnAward int16 = %d\n", guesses[774])
	fmt.Printf("var EndgamePassedPawnAward int16 = %d\n", guesses[775])
	fmt.Printf("var MiddlegameAdvancedPassedPawnAward int16 = %d\n", guesses[776])
	fmt.Printf("var EndgameAdvancedPassedPawnAward int16 = %d\n", guesses[777])
	fmt.Printf("var MiddlegameCandidatePassedPawnAward int16 = %d\n", guesses[778])
	fmt.Printf("var EndgameCandidatePassedPawnAward int16 = %d\n", guesses[779])
	fmt.Printf("var MiddlegameRookOpenFileAward int16 = %d\n", guesses[780])
	fmt.Printf("var EndgameRookOpenFileAward int16 = %d\n", guesses[781])
	fmt.Printf("var MiddlegameRookSemiOpenFileAward int16 = %d\n", guesses[782])
	fmt.Printf("var EndgameRookSemiOpenFileAward int16 = %d\n", guesses[783])
	fmt.Printf("var MiddlegameVeritcalDoubleRookAward int16 = %d\n", guesses[784])
	fmt.Printf("var EndgameVeritcalDoubleRookAward int16 = %d\n", guesses[785])
	fmt.Printf("var MiddlegameHorizontalDoubleRookAward int16 = %d\n", guesses[786])
	fmt.Printf("var EndgameHorizontalDoubleRookAward int16 = %d\n", guesses[787])
	fmt.Printf("var MiddlegamePawnFactorCoeff int16 = %d\n", guesses[788])
	fmt.Printf("var EndgamePawnFactorCoeff int16 = %d\n", guesses[789])
	fmt.Printf("var MiddlegameMobilityFactorCoeff int16 = %d\n", guesses[790])
	fmt.Printf("var EndgameMobilityFactorCoeff int16 = %d\n", guesses[791])
	fmt.Printf("var MiddlegameAggressivityFactorCoeff int16 = %d\n", guesses[792])
	fmt.Printf("var EndgameAggressivityFactorCoeff int16 = %d\n", guesses[793])
	fmt.Printf("var MiddlegameInnerPawnToKingAttackCoeff int16 = %d\n", guesses[794])
	fmt.Printf("var EndgameInnerPawnToKingAttackCoeff int16 = %d\n", guesses[795])
	fmt.Printf("var MiddlegameOuterPawnToKingAttackCoeff int16 = %d\n", guesses[796])
	fmt.Printf("var EndgameOuterPawnToKingAttackCoeff int16 = %d\n", guesses[797])
	fmt.Printf("var MiddlegameInnerMinorToKingAttackCoeff int16 = %d\n", guesses[798])
	fmt.Printf("var EndgameInnerMinorToKingAttackCoeff int16 = %d\n", guesses[799])
	fmt.Printf("var MiddlegameOuterMinorToKingAttackCoeff int16 = %d\n", guesses[800])
	fmt.Printf("var EndgameOuterMinorToKingAttackCoeff int16 = %d\n", guesses[801])
	fmt.Printf("var MiddlegameInnerMajorToKingAttackCoeff int16 = %d\n", guesses[802])
	fmt.Printf("var EndgameInnerMajorToKingAttackCoeff int16 = %d\n", guesses[803])
	fmt.Printf("var MiddlegameOuterMajorToKingAttackCoeff int16 = %d\n", guesses[804])
	fmt.Printf("var EndgameOuterMajorToKingAttackCoeff int16 = %d\n", guesses[805])
	fmt.Printf("var MiddlegamePawnShieldPenalty int16 = %d\n", guesses[806])
	fmt.Printf("var EndgamePawnShieldPenalty int16 = %d\n", guesses[807])
	fmt.Printf("var MiddlegameNotCastlingPenalty int16 = %d\n", guesses[808])
	fmt.Printf("var EndgameNotCastlingPenalty int16 = %d\n", guesses[809])
	fmt.Printf("var MiddlegameKingZoneOpenFilePenalty int16 = %d\n", guesses[810])
	fmt.Printf("var EndgameKingZoneOpenFilePenalty int16 = %d\n", guesses[811])
	fmt.Printf("var MiddlegameKingZoneMissingPawnPenalty int16 = %d\n", guesses[812])
	fmt.Printf("var EndgameKingZoneMissingPawnPenalty int16 = %d\n", guesses[813])
	fmt.Printf("var MiddlegameKnightOutpostAward int16 = %d\n", guesses[814])
	fmt.Printf("var EndgameKnightOutpostAward int16 = %d\n", guesses[815])
	fmt.Printf("var MiddlegameBishopPairAward int16 = %d\n", guesses[816])
	fmt.Printf("var EndgameBishopPairAward int16 = %d\n", guesses[817])

	// fmt.Printf("var MiddlegameCastlingAward int16 = %d\n", guesses[792])
	fmt.Println("===================================================")
}

func localOptimize(initialGuess []int16, K float64) []int16 {
	nParams := len(initialGuess)
	bestE := meanSquareError(testPositions, initialGuess, K)
	bestParValues := append([]int16{}, initialGuess...)
	improved := true
	for improved {
		improved = false
		for pi := 0; pi < nParams; pi++ {
			if _, ok := skipParams[pi]; ok {
				continue
			}
			if _, ok := futileIndices[pi]; ok {
				continue
			}
			bestParValues[pi] += 1
			newE := meanSquareError(testPositions, bestParValues, K)
			if newE < bestE {
				bestE = newE
				fmt.Println("Best parameters so far")
				printOptimalGuesses(bestParValues)
				improved = true
			} else {
				bestParValues[pi] -= 2
				if pi >= 768 && bestParValues[pi] < 0 {
					bestParValues[pi] += 1
					continue
				}
				newE = meanSquareError(testPositions, bestParValues, K)
				if newE < bestE {
					bestE = newE
					fmt.Println("Best parameters so far")
					printOptimalGuesses(bestParValues)
					improved = true
				} else {
					bestParValues[pi] += 1 // reset the guess
				}
			}
			fmt.Println("Current best E", bestE)
		}
	}
	return bestParValues
}

func findK() float64 {
	start := 0.0
	var end float64 = 10
	step := 1.0
	curr := start
	var err float64
	best := meanSquareError(testPositions, initialGuesses, start)

	for i := 0; i < K_PRECISION; i++ {

		// Find the minimum within [start, end] using the current step
		curr = start - step
		for curr < end {
			curr = curr + step
			err = meanSquareError(testPositions, initialGuesses, curr)
			if err <= best {
				best = err
				start = curr
			}
		}

		fmt.Printf("K so far is %f, iteration %d\n", start, i)

		// Adjust the search space
		end = start + step
		start = start - step
		step = step / 10.0
	}

	return start
}

func linearEvaluation(pos *Position) int16 {
	eval := Evaluate(pos)
	if pos.Turn() == Black {
		return -eval
	}
	return eval
}

func processor(testPositions []TestPosition, start int, end int, K float64) {
	var acc float64 = 0
	for i := start; i < end; i++ {
		eval := linearEvaluation(testPositions[i].pos)
		acc += math.Pow(testPositions[i].outcome-sigmoid(eval, K), 2)
	}

	answers <- acc
}

func meanSquareError(testPositions []TestPosition, guess []int16, K float64) float64 {
	acc := float64(0)
	updateEvalParams(guess)

	batchSize := len(testPositions) / NUM_PROCESSORS
	for i := 0; i < NUM_PROCESSORS; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if i == NUM_PROCESSORS-1 {
			end = len(testPositions)
		}
		go processor(testPositions, start, end, K)
	}

	for i := 0; i < NUM_PROCESSORS; i++ {
		ans := <-answers
		acc += ans
	}
	return acc / float64(len(testPositions))
}

func sigmoid(eval int16, K float64) float64 {
	return (1.0 / (1.0 + math.Pow(10, -(K*float64(eval))/400.0)))
}

func loadPositions(path string, actionFn func(string)) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	testPositions = make([]TestPosition, 0, 14_000_000)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		actionFn(line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func parseLine(line string) (string, float64) {
	fields := strings.Fields(line)
	fen := strings.Join(fields[:4], " ")
	fen = fmt.Sprintf("%s 0 1", fen)

	outcomeStr := strings.Trim(fields[5], "\";")
	var outcome float64
	if outcomeStr == "1/2-1/2" {
		outcome = 0.5
	} else if outcomeStr == "1-0" {
		outcome = 1.0
	} else if outcomeStr == "0-1" {
		outcome = 0.0
	} else {
		panic(fmt.Sprintf("Unexpected output %s", outcomeStr))
	}
	//
	// fen := strings.Trim(strings.Join(fields[:6], " "), ";")
	//
	// outcomeStr := strings.Replace(fields[8], "pgn=", "", -1)
	// outcome, e := strconv.ParseFloat(outcomeStr, 64)
	// if e != nil {
	// 	panic(e)
	// }
	// if fields[1] == "b" && outcomeStr == "1.0" {
	// 	outcome = 0
	// } else if fields[1] == "b" && outcomeStr == "0.0" {
	// 	outcome = 1
	// }

	return fen, outcome

}

func PrepareTuningData(path string) {
	loadPositions(path, func(line string) {
		fen, _ := parseLine(line)
		game := FromFen(fen, true)
		pos := game.Position()
		ml.Size = 0
		ml.Next = 0
		if pos.IsInCheck() {
			return
		}

		pos.GetCaptureMoves(ml)
		if ml.Size > 0 {
			return
		}
	})
}

func Tune(path string, toExclude map[int]bool) {
	skipParams = toExclude
	loadPositions(path, func(line string) {
		fen, outcome := parseLine(line)
		game := FromFen(fen, true)
		pos := game.Position()
		tp := TestPosition{pos, outcome}
		testPositions = append(testPositions, tp)
	})

	fmt.Printf("%d positions loaded\n", len(testPositions))
	K := findK()
	fmt.Printf("Optimal K is %f\n", K)
	optimalGuesses := localOptimize(initialGuesses, K)
	// tuningVars := make([]Parameter, len(initialGuesses))
	// for i, v := range initialGuesses {
	// 	tuningVars[i] = NewParameter(v)
	// 	if _, ok := futileIndices[i]; ok {
	// 		tuningVars[i].MaxValue = 0
	// 		tuningVars[i].MinValue = 0
	// 	} else if i >= 768 {
	// 		tuningVars[i].MinValue = 0
	// 	}
	// }
	// optimalGuesses := spsaTuning(tuningVars, 1000, K)
	fmt.Println("Optimal Parameters have been found!!")
	fmt.Println("===================================================")
	printOptimalGuesses(optimalGuesses)
	close(answers)
}

type Parameter struct {
	MaxValue float64
	MinValue float64

	OriginalValue int16

	C_END float64
	R_END float64
	C     float64
	A     float64
}

func NewParameter(variable int16) Parameter {
	return Parameter{
		MaxValue:      300.0,
		MinValue:      -300.0,
		OriginalValue: variable,
		C_END:         4.0,
		R_END:         0.002,
		C:             0.0,
		A:             0.0,
	}
}

type DataPoint struct {
	Iteration int

	ThetaError float64

	ThetaPlusError  float64
	ThetaMinusError float64
}

func WriteDataPoint(data []DataPoint) {
	fmt.Println("Iteration,Error,+Error,-Error")
	for _, dp := range data {
		fmt.Printf("%d,%f,%f,%f\n", dp.Iteration, dp.ThetaError, dp.ThetaPlusError, dp.ThetaMinusError)
	}
}
func changedError(newValues []float64, K float64) float64 {

	params := make([]int16, len(newValues))
	// Step 1. Convert newValues to int16, and send it to meanSquareError (which also updates the current variables)
	for i := 0; i < len(newValues); i++ {
		params[i] = int16(newValues[i])
	}

	// Step 2. Compute the new error and return
	return meanSquareError(testPositions, params, K)
}

// This is taken almost verbatim from Loki, many thanks to the author
// Niels who helped me a lot to undrestand it.
// https://github.com/BimmerBass/Loki
// I am currently not using it, because with my current settings, the local optimize did
// an amazing job, but I intend to come back to this and make it work.
func spsaTuning(tuningVars []Parameter, iterations int, K float64) []int16 {
	// Step 2. Set up the vector of parameter values, we call it theta here.
	data := make([]DataPoint, 0, iterations+1)    // This creates a zero size array (hence the second parameter), but which can grow upto iterations+1
	theta := make([]float64, len(initialGuesses)) // this creates an array of size len(initializeGuesses)
	for p := 0; p < len(theta); p++ {
		theta[p] = float64(tuningVars[p].OriginalValue)
	}

	// Step 4. Calculate SPSA constants

	BIG_A := 0.1 * float64(iterations)
	alpha := 0.602
	gamma := 0.101

	for p := 0; p < len(theta); p++ {
		tuningVars[p].C = tuningVars[p].C_END * math.Pow(float64(iterations), gamma)

		a_end := tuningVars[p].R_END * math.Pow(tuningVars[p].C, 2.0)

		tuningVars[p].A = a_end * math.Pow(BIG_A+float64(iterations), alpha)
	}

	// Step 5. Run the tuning with the given number of iterations
	var theta_plus = make([]float64, len(tuningVars))
	var theta_minus = make([]float64, len(tuningVars))

	// Initialize perturbation vector.
	var delta = make([]float64, len(tuningVars))

	/*
		Zero-initialization of step-size and perturbation vectors
	*/
	var an = make([]float64, len(theta))
	var cn = make([]float64, len(theta))

	data = append(data, DataPoint{Iteration: 0, ThetaError: changedError(theta, K), ThetaPlusError: 0, ThetaMinusError: 0})
	for n := 0; n < iterations; n++ {

		for p := 0; p < len(tuningVars); p++ {
			an[p] = tuningVars[p].A / (math.Pow(BIG_A+float64(n)+1, alpha))

			cn[p] = tuningVars[p].C / (math.Pow(float64(n)+1, gamma))
		}

		// I dont' need to clear theta_plus, theta_minus and data
		// becase I use delta[p] syntax, instead of push, which overrides
		// the current values

		// Step 5C. Determine delta and thus theta_plus and theta_minus
		for p := 0; p < len(tuningVars); p++ {
			d := randemacher()

			delta[p] = d // Add the scores to the delta vector.

			// Step 5C.1. Compute theta_plus and theta_minus from these values

			theta_plus[p] = math.Round(float64(theta[p]) + float64(d)*cn[p])
			theta_minus[p] = math.Round(float64(theta[p]) - float64(d)*cn[p])
		}

		// Step 5D. Compute the error of theta_plus and theta_minus respectively.
		thetaPlusErr := float64(len(testPositions)) * changedError(theta_plus, K)
		thetaMinusErr := float64(len(testPositions)) * changedError(theta_minus, K)

		// Step 5E. Compute the gradient for each variable and adjust accordingly.
		for p := 0; p < len(tuningVars); p++ {
			g_hat := (thetaPlusErr - thetaMinusErr) / (2.0 * cn[p] * delta[p])

			// Now adjust the theta values based on the gradient.
			theta[p] -= math.Round(an[p] * g_hat)

			// Lastly, we'll have to make sure that we don't go out of the designated bounds
			theta[p] = math.Max(tuningVars[p].MinValue, math.Min(tuningVars[p].MaxValue, theta[p]))
		}

		for p := 0; p < len(tuningVars); p++ {
			fmt.Printf("Iteration: %d: Tuned: [ %d ]: %f, (Original: [ %d ])\n", n, p, theta[p], tuningVars[p].OriginalValue)
		}

		// Step 5A. Compute the current error. This is only used for outputting the progress.
		err := changedError(theta, K)
		dp := DataPoint{Iteration: n + 1, ThetaError: err, ThetaPlusError: thetaPlusErr, ThetaMinusError: thetaMinusErr}
		data = append(data, dp)
		fmt.Printf("Iteration data point: %d,%f,%f,%f\n", dp.Iteration, dp.ThetaError, dp.ThetaPlusError, dp.ThetaMinusError)
		fmt.Println("----------------------------------")
	}

	params := make([]int16, len(theta))
	for i, p := range theta {
		params[i] = int16(p)
	}

	WriteDataPoint(data)
	return params
}

func min(x, y int16) int16 {
	if x < y {
		return x
	}
	return y
}

func max(x, y int16) int16 {
	if x > y {
		return x
	}
	return y
}

var rnd = rand.New(rand.NewSource(99))

func randemacher() float64 {
	if rnd.Intn(2) < 1 {
		return 1
	}
	return -1
}
