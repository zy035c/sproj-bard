package simulation

import (
	"sync/atomic"
)

var last_read atomic.Int64
var MRC_metric atomic.Int64
var RYWC_metric atomic.Int64
var MWC_metric atomic.Int64
var EC_metric atomic.Int64
var CC_metric atomic.Int64

func InitMetrics() {
	last_read.Store(0)
	MRC_metric.Store(0)
	RYWC_metric.Store(0)
	MWC_metric.Store(0)
	EC_metric.Store(0)
	CC_metric.Store(0)
}

func PrintMetrics(vid int, version_chain []int, lastWriteVid int, mid int) (int64, int64, int64, int64, int64) {
	return PrintMRC(vid, version_chain), PrintRYWC(vid, version_chain, lastWriteVid, mid), PrintMWC(vid, version_chain), PrintEC(vid, version_chain), PrintCC(vid, version_chain)
}

func PrintMRC(vid int, version_chain []int) int64 {
	if len(version_chain) >= 2 {
		n := len(version_chain)
		if version_chain[n-1] < version_chain[n-2] {
			MRC_metric.Add(1)
		}
	}
	return MRC_metric.Load()
}

func PrintRYWC(vid int, version_chain []int, lastWriteVid int, mid int) int64 {
	if len(version_chain) == 0 {
		return RYWC_metric.Load()
	}
	vr := version_chain[len(version_chain)-1]

	if lastWriteVid > vr {
		RYWC_metric.Add(1)
	}

	return RYWC_metric.Load()

}

func PrintMWC(vid int, version_chain []int) int64 {
	s := 0
	for i, v := range version_chain {
		if i == 0 {
			continue
		}
		if v < version_chain[i-1] {
			s += 1
		}
	}
	MWC_metric.Store(int64(s))
	return MWC_metric.Load()
}

func PrintEC(vid int, version_chain []int) int64 {
	// if len(version_chain) == 0 {
	// 	return 0
	// }
	// last := version_chain[len(version_chain)-1]
	// if vid == last {
	// 	return EC_metric.Load()
	// }
	// EC_metric.Add(1)
	// return EC_metric.Load()
	return 0
}

func PrintCC(vid int, version_chain []int) int64 {
	// Find the maximum number in the list
	maxNum := 0
	for _, num := range version_chain {
		if num > maxNum {
			maxNum = num
		}
	}

	// Create a frequency map to count occurrences
	freq := make(map[int]int)
	for _, num := range version_chain {
		freq[num]++
	}

	// Calculate the score
	score := 0
	for i := 1; i <= maxNum; i++ {
		count, exists := freq[i]
		if !exists {
			// Missing number
			score++
		} else if count > 1 {
			// Extra occurrences
			score += count - 1
		}
	}
	CC_metric.Add(int64(score))
	return CC_metric.Load()
}
