package main

import (
	"fmt"
	"lamport/models"
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	/**
		Here's how the algorithm works:

	    Each process maintains a logical clock, and whenever it performs a local event, its logical clock is incremented.
	    When a process sends a message, it attaches its logical clock value to the message. Upon receiving the message, the recipient process updates its own logical clock to be greater than the timestamp in the received message.
	    If two events have the same timestamp, the tie is broken by comparing process identifiers.
	*/

	/* init lamport counter */
	var lamportCounter uint64 = 0

	server := gin.Default()

	server.Use(func(c *gin.Context) {
		c.Set("lamport-counter", &lamportCounter)
		// db connection ...
		c.Next()
	})

	server.GET("/get-data", GetOrder)
	server.GET("/get-proc-id", GetProcID)
	server.GET("/insert-data", InsertData)
	server.GET("/sync", RcvMsg)
}

type RequestBody struct {
	ID      uint64 `json:"id"`
	ProcID  int    `json:"proc-id"`
	ReqType string `json:"req-type"`

	Order models.Order `json:"order"`
}

/* router func: get by id */
func GetOrder(c *gin.Context) {
	// get by id
	// parse request by bind to a struct
	// return the struct

	var reqBody RequestBody
	if err := c.Bind(&reqBody); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if reqBody.ReqType != "get" {
		c.JSON(200, gin.H{"error": "invalid request type"})
		return
	}

	pid := os.Getpid()
	if reqBody.ProcID != pid {
		println("I received a request that is not for this proc id")
		return
	}

	timestamp := GetTimestamp(c)
	fmt.Println("I will fetch the db for this proc: ", pid, "at timestamp: ", timestamp)

	db_conn := make_conn(c)

	order, err := models.GetMostRecent(timestamp, db_conn)
	if err != nil {
		c.JSON(404, gin.H{"error": err})
		return
	}

	fmt.Println("I fetched the db for pid: ", *order)
}

func GetProcID(c *gin.Context) {
	// get the proc id of this server
	pid := os.Getpid()
	c.JSON(200, gin.H{"proc-id": pid})
}

func InsertData(c *gin.Context) {
	// insert data
	// parse request by bind to a struct
	// return the struct

	var reqBody RequestBody
	if err := c.Bind(&reqBody); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if reqBody.ReqType != "insert" {
		c.JSON(200, gin.H{"error": "invalid request type"})
		return
	}

	if reqBody.ProcID != os.Getpid() {
		println("I received a request that is not for this proc id")
		return
	}

	fmt.Println("I will insert the db for id: ", reqBody.ID)
}

func ResetLocalDatabase() {
	// this func will drop/create a db
}

func RcvMsg(c *gin.Context) {
	var reqBody RequestBody
	if err := c.Bind(&reqBody); err != nil {
		c.JSON(400, gin.H{"error": "invalid sync msg"})
		return
	}

	if reqBody.ReqType != "sync" {
		c.JSON(200, gin.H{"error": "invalid request type"})
		return
	}

	/**
	 *  Lamport's logical clock algorithm
	 */
	order := reqBody.Order
	local_timestamp := GetTimestamp(c)
	if order.Timestamp > local_timestamp {
		println("I will update my local timestamp to: ", order.Timestamp)
		atomic.StoreUint64(
			c.MustGet("lamport-counter").(*uint64),
			order.Timestamp+1,
		)

		order.Insert(make_conn(c))
	} else {
		// if the timestamp is less than or equal to the local timestamp
		// then we don't need to update the db
	}

	fmt.Println("Sync the db for pid: ", order)
}

func compute() {
	// simulate a time consuming computation

	time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
	// sleep for a random time between 0 and 2000 ms
}

func GetTimestamp(c *gin.Context) uint64 {
	// get the timestamp atomically
	lamportCounter := c.MustGet("lamport-counter").(*uint64)
	return atomic.LoadUint64(lamportCounter)
}
