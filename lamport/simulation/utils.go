package simulation

import (
	"lamport/chart"
	"lamport/goclock"
	"math"
	"strconv"
)

func GenerateCorrectVersionChain(epoch int) []int {
	if epoch < 0 {
		return []int{}
	}
	chain := []int{}
	for i := epoch - 1; i > 0; i-- {
		chain = append(chain, i)
	}
	return chain
}

func CalculateScore(versionChain []string, epoch int) int {
	standardChain := GenerateCorrectVersionChain(epoch)
	roundList := make([]int, len(versionChain))
	for i, x := range versionChain {
		num, _ := strconv.Atoi(x[5:]) // remove chars 'epoch'
		roundList[i] = num
	}
	// fmt.Println("standardChain", standardChain)
	// fmt.Println("roundList", roundList)
	return DamerauLevenshteinScore(roundList, standardChain)
}

func CalcScore(epoch int, m Machine[string, int]) float64 {
	vchain := m.GetVersionChainData()
	// fmt.Println("- vchain =", vchain)
	score := CalculateScore(vchain, epoch)
	// fmt.Println("- Score =", score)
	var good_chain_rate float64 = float64(score) / float64(epoch)
	return good_chain_rate
}

func PlotScore(epoch int, ms []goclock.Machine[string, int], time float64, dyn_chart *chart.DynamicChart) float64 {
	var score_sum float64 = 0
	for _, m := range ms {
		score := CalcScore(epoch, m)
		score_sum += score
	}
	// fmt.Println("- Score_sum ", score_sum)
	score := score_sum / float64(len(ms))
	go dyn_chart.SendDataPoint(chart.DataPoint{X: float64(epoch), Y: score})
	return score
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
