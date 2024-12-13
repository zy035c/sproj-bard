package simulation

import (
	"fmt"
	"lamport/chart"
	"lamport/goclock"
	"lamport/utils"
	"math"
	"strconv"
	"strings"
)

func GenerateCorrectVersionChain(vid int) []int {
	if vid < 0 {
		return []int{}
	}
	chain := []int{}
	for i := 1; i < vid; i++ {
		chain = append(chain, i)
	}
	return chain
}

func VerChainStrToInt(chain []string) []int {
	roundList := make([]int, len(chain))
	for i, x := range chain {
		num, _ := strconv.Atoi(x[3:]) // remove chars 'Vid'
		roundList[i] = num
	}
	return roundList
}

func CalcFlawRate(vid int, machineChain []int) float64 {
	// fmt.Println("- vchain =", machineChain)
	standardChain := GenerateCorrectVersionChain(vid)
	score := DamerauLevenshteinScore(machineChain, standardChain)
	// fmt.Println("- Score =", score)
	var flaw_rate float64 = float64(score) / float64(vid)
	return flaw_rate
}

func PlotFlawMetric(vid int, ms []goclock.Machine[string, int], time float64, dyn_chart *chart.DynamicChart) {
	var score_sum float64 = 0
	for _, m := range ms {
		flaw_rate := CalcFlawRate(vid, VerChainStrToInt(m.GetVersionChainData()))
		if flaw_rate > 0 {
			fmt.Println(">> Flaw detected: ", strings.Join(m.GetVersionChainData(), "->"))
		}
		score_sum += flaw_rate
	}
	// fmt.Println("- Score_sum ", score_sum)
	score := score_sum / float64(len(ms))
	go dyn_chart.SendDataPoint(chart.DataPoint{X: float64(vid), Y: score})
	// fmt.Println("- Metric Flaw Rate: ", score)
}

func RandomSampleNodeVersionChain(ms []goclock.Machine[string, int], rid int, version_chain []string) {
	// fmt.Println("- Sampling Node", rid, strings.Join(version_chain, "->"))
}

func PrintCycleMetric(epoch int, ms []goclock.Machine[string, int]) {
	chains := make([][]int, 1024)
	for _, m := range ms {
		chains = append(chains, VerChainStrToInt(m.GetVersionChainData()))
	}
	chains = append(chains, GenerateCorrectVersionChain(epoch))
	nCycle := CalcDependGraphCycle(chains, epoch)
	fmt.Println("- Metric Dependency Cycle N: ", nCycle)
}

// min returns the minimum value among the given integers
func min(a, b, c, d int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	if d < m {
		m = d
	}
	return m
}

// DamerauLevenshteinScore computes the Damerau-Levenshtein distance between two strings
func DamerauLevenshteinScore(s1, s2 []int) int {
	lenS1 := len(s1)
	lenS2 := len(s2)

	// Create a 2D slice to store distances
	distance := make([][]int, lenS1+1)
	for i := range distance {
		distance[i] = make([]int, lenS2+1)
	}

	// Initialize the distance matrix
	for i := 0; i <= lenS1; i++ {
		distance[i][0] = i
	}
	for j := 0; j <= lenS2; j++ {
		distance[0][j] = j
	}

	// Compute the distance
	for i := 1; i <= lenS1; i++ {
		for j := 1; j <= lenS2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			// Calculate minimum distance considering insert, delete, and replace
			distance[i][j] = min(
				distance[i-1][j]+1,      // Deletion
				distance[i][j-1]+1,      // Insertion
				distance[i-1][j-1]+cost, // Substitution
				math.MaxInt,             // Initialize with a high value; will be updated for transposition
			)

			// Check for transposition (swap of adjacent characters)
			if i > 1 && j > 1 && s1[i-1] == s2[j-2] && s1[i-2] == s2[j-1] {
				distance[i][j] = min(
					distance[i][j],
					distance[i-2][j-2]+cost, // Transposition
					math.MaxInt,
					math.MaxInt,
				)
			}
		}
	}

	// Return the final computed distance
	return distance[lenS1][lenS2]
}

func CalcDependGraphCycle(chains [][]int, maxEpoch int) int {
	/*
		* Input:
			chains is a slice of slices
				index: Machine Id: int
				value: The Version Chain stored on that machine

		*
		* Each chain starts with 0, with the last number being the latest version that node considers
		* Now build a graph with nodes each represents a Version Id
		* We have two sources for partial orders (directed edge):
		1. On a single version chain, e.g. chains[3] = [0, 1, 3, 2, 4, 6]
		We conclude that these nodes are in order:
			0 BEFORE 1 BEFORE 3 BEFORE 2 BEFORE 4 BEFORE 6
			We add directed edge if not exists.

		2. The number has an ensured being monotonically increased:
			0 BEFORE 1 BEFORE 2 BEFORE 3 BEFORE 4 BEFORE (5) BEFORE 6
			We add directed edge if not exists.

		Iterate over each chain to add edges, add monotonic order edges.
		Calculate the total number of directed cycle in the graph.
		Also write a test function.
	*/

	adj := make(map[int][]int)

	for _, chain := range chains {

		for i := 0; i < len(chain)-1; i++ {
			from := chain[i]
			to := chain[i+1]
			if _, exists := adj[from]; !exists {
				adj[from] = []int{}
			}
			adj[from] = utils.Append_if_not_exist(to, adj[from])
		}
	}

	return utils.CountCycles(adj)
}

// \de\addplot[color=red, mark=square*, thick] coordinates {(1, 100) (10, 200) (100, 400) (1000, 800) (10000, 900)};
// \addplot[color=blue, mark=*, thick] coordinates {(1, 150) (10, 250) (100, 500) (1000, 600) (10000, 850)};
// \addplot[color=Emerald, mark=triangle*, thick] coordinates {(1, 200) (10, 300) (100, 450) (1000, 700) (10000, 950)};
// \addplot[color=BurntOrange, mark=x, thick] coordinates {(1, 50) (10, 100) (100, 300) (1000, 500) (10000, 750)};
