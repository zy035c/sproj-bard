package timestamp

import (
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

func GetTimestamp(c *gin.Context) uint64 {
	// get the timestamp atomically
	lamportCounter := c.MustGet("lamport-counter").(*uint64)
	return atomic.LoadUint64(lamportCounter)
}

func IncrementTimestamp(c *gin.Context) uint64 {
	// increment the timestamp atomically
	lamportCounter := c.MustGet("lamport-counter").(*uint64)
	atomic.AddUint64(lamportCounter, 1)

	return atomic.LoadUint64(lamportCounter)
}

func SetTimestamp(c *gin.Context, timestamp uint64) {
	// set the timestamp atomically
	lamportCounter := c.MustGet("lamport-counter").(*uint64)
	atomic.StoreUint64(lamportCounter, timestamp)
}
