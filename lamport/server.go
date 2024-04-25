package main

import (
	"bytes"
	"container/heap"
	"encoding/json"
	"flag"
	"fmt"
	"lamport/controller"
	"lamport/models"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"lamport/mypq"
	"lamport/timestamp"

	"github.com/gin-gonic/gin"
)

// read command options
// -port followed by a port number
// -ps followed by a list of port numbers of other processes
// example: go run server.go -port=9090 -ps=8000,9000,10000

func main() {

	var port int
	var ps []string

	// Read port
	if gPort := os.Getenv("GPORT"); gPort != "" {
		fmt.Println("Using GPORT environment variable:", gPort)
		port = parseIntOrPanic(gPort)
	} else {
		// if no env var, read from cmd option
		portNum := flag.Int("port", 8080, "port number")
		flag.Parse()
		port = *portNum
	}

	// Read port list
	if gPs := os.Getenv("GPLIST"); gPs != "" {
		fmt.Println("Using GPLIST environment variable:", gPs)
		ps = strings.Split(gPs, ",")
		for _, p := range ps {
			fmt.Println("Other process port:", p)
		}
	} else {
		// if no env var, read from cmd option
		ps_ := flag.String("ps", "", "list of port numbers of other processes")
		flag.Parse()
		if *ps_ != "" {
			ps = strings.Split(*ps_, ",")
			for _, p := range ps {
				fmt.Println("Other process port:", p)
			}
		}
	}

	server := gin.Default()

	/* init lamport counter */
	var lamportCounter uint64 = 0
	fake_db := make(map[uint64]models.Order)
	taskQueue := mypq.NewPriorityQueue()

	server.Use(func(c *gin.Context) {

		c.Set("lamport-counter", &lamportCounter)

		c.Set("db_conn", nil) // not in use

		/* As of now, we don't store data in db. Instead, we use a hashmap. */
		c.Set("fake_db", &fake_db)

		c.Set("task_queue", taskQueue)

		c.Set("port_list", ps)

		// c.Next()
	})

	go func() {
		taskQueue.LoopAndPoll(&lamportCounter)
	}()

	server.GET("/get-data", GetOrder)
	server.GET("/get-proc-id", GetProcID)
	server.POST("/insert-data", InsertData)
	server.POST("/sync", RcvMsg)

	server.GET("/get-port", func(c *gin.Context) {
		c.JSON(200, gin.H{"port": port})
	})

	println("-- My proc id: ", os.Getpid())

	server.Run(fmt.Sprintf(":%d", port))
}

type RequestBody struct {
	ID      uint64 `json:"id"`
	Port    int    `json:"port"`
	ReqType string `json:"req_type" binding:"required"`

	Order models.Order `json:"order"`
}

/* router func: get by id */
func GetOrder(c *gin.Context) {
	// get by current timestamp
	// parse request by bind to a struct
	// return the struct

	fmt.Println("[Router] Will fetch the db for this proc:", os.Getpid(), "ts=", timestamp.GetTimestamp(c))
	order, err := controller.GetMostRecent(c, timestamp.GetTimestamp(c))
	fmt.Println(err)
	if err != nil {
		c.JSON(404, gin.H{"error": err})
		return
	}

	fmt.Println("[Router] Fetched data", order)

	c.JSON(200, gin.H{"order": order})
}

func GetProcID(c *gin.Context) {
	// get the proc id of this server
	pid := os.Getpid()
	c.JSON(200, gin.H{"proc-id": pid})
}

func InsertData(c *gin.Context) {
	// insert (or said update) data
	// parse request by bind to a struct
	// return the struct

	var reqBody RequestBody
	if err := c.Bind(&reqBody); err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if reqBody.ReqType != "insert" {
		c.JSON(200, gin.H{"error": "invalid request type"})
		return
	}

	computeOrDelay()
	task := func() {
		computeOrDelay()

		reqBody.Order.ProcID = os.Getpid()
		_, err := controller.InsertOrUpdate(c, &reqBody.Order, timestamp.IncrementTimestamp(c))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("[Queue Task] Insertion ok. Current ts=", timestamp.GetTimestamp(c))
		reqBody.Order.Timestamp = timestamp.GetTimestamp(c) // ?
		SendSyncMsg(&reqBody.Order, GetPortList(c))
		fmt.Println("[Queue Task] Send sync msg ok.")
	}

	// fmt.Println("[Queue Task] Will insert to db. Priority:", timestamp.GetTimestamp(c))
	taskQueue := GetTaskQueue(c)
	heap.Push(taskQueue, &mypq.Item{Value: task, Priority: timestamp.GetTimestamp(c)})

	c.JSON(200, gin.H{"status": "ok"})
}

func RcvMsg(c *gin.Context) {
	var reqBody RequestBody
	if err := c.Bind(&reqBody); err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"error": "invalid sync msg"})
		return
	}

	if reqBody.ReqType != "sync" {
		fmt.Println("Invalid request type")
		c.JSON(200, gin.H{"error": "invalid request type"})
		return
	}

	order := reqBody.Order
	fmt.Println("[Router] Received sync msg from proc: ", order.ProcID, " ts=", order.Timestamp)
	computeOrDelay()

	/**
	 *  Lamport's logical clock algorithm
		Here's how the algorithm works:

	    Each process maintains a logical clock, and whenever it performs a local event, its logical clock is incremented.
	    When a process sends a message, it attaches its logical clock value to the message. Upon receiving the message,
		the recipient process updates its own logical clock to be greater than the timestamp in the received message.
	    If two events have the same timestamp, the tie is broken by comparing process identifiers.
	*/
	task := func() {
		computeOrDelay()

		order := reqBody.Order
		cur_timestamp := timestamp.GetTimestamp(c)

		if order.Timestamp == cur_timestamp {
			fmt.Println("[Queue Task] Tie at ts.")
			fmt.Println("[Queue Task] Compare proc id. Incoming pid:", order.ProcID, "Local pid:", controller.GetLastProcID(c, cur_timestamp))
			if order.ProcID > controller.GetLastProcID(c, cur_timestamp) {
				fmt.Println("Greater incoming pid, will accept entry")
				controller.InsertOrUpdate(c, &order, order.Timestamp)
				fmt.Println("[Queue Task] Update entry to", order, "Priority:", order.Timestamp)
			} else {
				fmt.Println("[Queue Task] Outdated sync msg and less pid, abort.", "ts=", order.Timestamp, "local ts=", cur_timestamp)
			}

		} else if order.Timestamp > cur_timestamp {

			fmt.Println("[Queue Task] I will update my local ts", cur_timestamp, "to: ", order.Timestamp)
			timestamp.SetTimestamp(c, order.Timestamp) // should only increment by 1
			controller.InsertOrUpdate(c, &order, order.Timestamp)

		} else {
			// if the timestamp is less than or equal to the local timestamp
			// then we don't need to update the db
			fmt.Println("[Queue Task] Outdated sync msg, abort.", "ts=", order.Timestamp, "local ts=", cur_timestamp)
		}
	}

	heap.Push(GetTaskQueue(c), &mypq.Item{Value: task, Priority: order.Timestamp})

	c.JSON(200, gin.H{"status": "ok"})
}

func SendSyncMsg(ord *models.Order, portList []string) {
	// send http request to all the ther process
	// and update the local timestamp
	requestJSON := struct {
		ReqType string       `json:"req_type"`
		Order   models.Order `json:"order"`
	}{
		ReqType: "sync",
		Order:   *ord,
	}

	requestBody, _ := json.Marshal(requestJSON)

	for _, port := range portList {
		// post sync to the server
		_, err := http.Post(
			fmt.Sprintf("http://localhost:%s/sync", port),
			"application/json",
			bytes.NewBuffer(requestBody),
		)
		if err != nil {
			fmt.Println("Error sending for port ", port, ": ", err)
			continue
		}
	}
}

func computeOrDelay() {
	// simulate a time consuming computation
	time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
	// sleep for a random time between 0 and 2000 ms
}

func GetTaskQueue(c *gin.Context) *mypq.PriorityQueue {
	taskQueue := c.MustGet("task_queue").(*mypq.PriorityQueue)
	return taskQueue
}

func GetPortList(c *gin.Context) []string {
	portList := c.MustGet("port_list").([]string)
	return portList
}

func GetRouterLock(c *gin.Context) *sync.Mutex {
	lock := c.MustGet("router_lock").(*sync.Mutex)
	return lock
}
