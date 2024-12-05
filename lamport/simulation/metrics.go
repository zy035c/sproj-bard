package simulation

import (
	"sync/atomic"
)

var last_read atomic.Int64
var MRC_metric atomic.Int64
var RYWC_metric atomic.Int64
var MWC_metric atomic.Int64

func InitMetrics() {
	last_read.Store(0)
	MRC_metric.Store(0)
	RYWC_metric.Store(0)
	MWC_metric.Store(0)
}

func PrintMetrics(vid int, version_chain []int, lastWriteVid int, mid int) (int64, int64, int64) {
	return PrintMRC(vid, version_chain), PrintRYWC(vid, version_chain, lastWriteVid, mid), PrintMWC(vid, version_chain)
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
