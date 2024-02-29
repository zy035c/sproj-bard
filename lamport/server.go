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
	"strconv"
	"strings"
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
	/**
		Here's how the algorithm works:

	    Each process maintains a logical clock, and whenever it performs a local event, its logical clock is incremented.
	    When a process sends a message, it attaches its logical clock value to the message. Upon receiving the message, the recipient process updates its own logical clock to be greater than the timestamp in the received message.
	    If two events have the same timestamp, the tie is broken by comparing process identifiers.
	*/
	port := flag.Int("port", 8080, "port number")
	ps := flag.String("ps", "", "list of port numbers of other processes")
	flag.Parse()

	var psList []string
	if *ps != "" {
		psList = strings.Split(*ps, ",")
		for _, p := range psList {
			fmt.Println("Other process port:", p)
		}
	}
	server := gin.Default()

	// server.MaxHandlers = 1
	/* init lamport counter */
	var lamportCounter uint64 = 0
	fake_db := make(map[uint64]models.Order)
	taskQueue := mypq.PriorityQueue{}

	server.Use(func(c *gin.Context) {

		c.Set("lamport-counter", &lamportCounter)
		// db connection ...
		c.Set("db_conn", nil)

		/* As of now, we don't store data in db. Instead, we use a hashmap. */
		c.Set("fake_db", &fake_db)

		c.Set("task_queue", &taskQueue)

		c.Set("port_list", psList)

		c.Next()
	})

	server.GET("/get-data", GetOrder)
	server.GET("/get-proc-id", GetProcID)
	server.GET("/insert-data", InsertData)
	server.POST("/sync", RcvMsg)
	server.GET("/start", PollTaskQueue)

	println("-My proc id: ", os.Getpid())

	server.Run(fmt.Sprintf(":%d", *port))

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

	fmt.Println("I will fetch the db for this proc: ", pid, "at timestamp: ", timestamp.GetTimestamp(c))
	order, err := controller.GetMostRecent(c, timestamp.GetTimestamp(c))
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
	// insert (or said update) data
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

	// if reqBody.ProcID != os.Getpid() {
	// 	println("I received a request that is not for this proc id")
	// 	return
	// }

	fmt.Println("I will insert the db for id: ", reqBody.ID)

	task := func() {
		computeOrDelay()

		_, err := controller.InsertOrUpdate(c, &reqBody.Order, timestamp.IncrementTimestamp(c))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		SendSyncMsg(&reqBody.Order, GetPortList(c))
	}

	taskQueue := GetTaskQueue(c)
	heap.Push(taskQueue, task)

	c.JSON(200, gin.H{"status": "ok"})
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

	order := reqBody.Order
	println("I received a sync msg from proc: ", order.ProcID, " timestamp: ", order.Timestamp)

	/**
	 *  Lamport's logical clock algorithm
	 */
	local_timestamp := timestamp.GetTimestamp(c)
	if order.Timestamp == local_timestamp {
		println("Tie: look at procId")
		if order.ProcID > os.Getpid() {

			task := mypq.Item{Value: func() {
				controller.InsertOrUpdate(c, &order, timestamp.GetTimestamp(c))
			}}
			task.SetPriority(order.Timestamp)
			taskQueue := GetTaskQueue(c)

			heap.Push(taskQueue, &task)
		}

	} else if order.Timestamp > local_timestamp {

		task := func() {
			println("I will update my local timestamp to: ", order.Timestamp)
			timestamp.SetTimestamp(c, order.Timestamp)
			controller.InsertOrUpdate(c, &order, order.Timestamp)
		}

		taskQueue := GetTaskQueue(c)
		heap.Push(taskQueue, &mypq.Item{Value: task, Priority: order.Timestamp})
	} else {
		// if the timestamp is less than or equal to the local timestamp
		// then we don't need to update the db
		println("Received a old sync msg", order.Timestamp, local_timestamp)
	}

	c.JSON(200, gin.H{"status": "ok"})
}

func SendSyncMsg(ord *models.Order, portList []string) {
	// send http request to all the ther process
	// and update the local timestamp

	orderJSON, _ := json.Marshal(*ord)

	requestBody, _ := json.Marshal(
		map[string][]byte{
			"ReqType": []byte("sync"),
			"ProcID":  []byte(strconv.Itoa(os.Getpid())),
			"Order":   orderJSON,
		},
	)

	for _, port := range portList {
		// send request to the server

		_, err := http.Post(
			fmt.Sprintf("http://localhost:%s/sync", port),
			"application/json",
			bytes.NewBuffer(requestBody),
		)
		if err != nil {
			fmt.Println("Error: ", err)
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

func PollTaskQueue(c *gin.Context) {
	// block until the task's priority is less than or equal to the current timestamp
	// pop the task from the queue and execute it
	go func() {
		taskQueue := GetTaskQueue(c)
		for {

			if taskQueue == nil || taskQueue.Len() == 0 {
				continue
			}
			item := heap.Pop(taskQueue).(*mypq.Item)
			cur_timestamp := timestamp.GetTimestamp(c)

			if item.GetPriority() == cur_timestamp || item.GetPriority() == cur_timestamp+1 {
				// convert the value to a function and execute it
				item.Execute()
			}

		}
	}()
}
