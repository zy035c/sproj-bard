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
	"sync/atomic"
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
	taskQueue := &mypq.PriorityQueue{}

	routerLock := sync.Mutex{}

	server.Use(func(c *gin.Context) {

		c.Set("lamport-counter", &lamportCounter)
		// db connection ...
		c.Set("db_conn", nil)

		/* As of now, we don't store data in db. Instead, we use a hashmap. */
		c.Set("fake_db", &fake_db)

		c.Set("task_queue", taskQueue)

		c.Set("port_list", psList)

		c.Set("router_lock", &routerLock)
		// c.Next()
	})

	go func() {
		for {
			if taskQueue.Len() == 0 {
				continue
			}
			item := heap.Pop(taskQueue).(mypq.Item)
			cur_timestamp := atomic.LoadUint64(&lamportCounter)

			if item.GetPriority() == cur_timestamp || item.GetPriority() == cur_timestamp+1 {
				// convert the value to a function and execute it
				fmt.Println("Will execute item. Current timestamp", cur_timestamp, "; item's timestamp", item.GetPriority())
				item.Execute()
			} else {
				heap.Push(taskQueue, &item)
			}
		}
	}()

	server.GET("/get-data", GetOrder)
	server.GET("/get-proc-id", GetProcID)
	server.POST("/insert-data", InsertData)
	server.POST("/sync", RcvMsg)

	server.GET("/get-port", func(c *gin.Context) {
		c.JSON(200, gin.H{"port": *port})
	})

	println("-- My proc id: ", os.Getpid())

	server.Run(fmt.Sprintf(":%d", *port))
}

type RequestBody struct {
	ID      uint64 `json:"id"`
	ProcID  int    `json:"proc-id"`
	Port    int    `json:"port"`
	ReqType string `json:"req-type" binding:"required"`

	Order models.Order `json:"order"`
}

/* router func: get by id */
func GetOrder(c *gin.Context) {
	// get by current timestamp
	// parse request by bind to a struct
	// return the struct

	// c.Query("proc-id")

	// pid := os.Getpid()
	// if reqBody.ProcID != pid {
	// 	println("I received a request that is not for this proc id")
	// 	return
	// }

	fmt.Println("Will fetch the db for this proc:", os.Getpid(), "at timestamp:", timestamp.GetTimestamp(c))
	order, err := controller.GetMostRecent(c, timestamp.GetTimestamp(c))
	fmt.Println(err)
	if err != nil {
		c.JSON(404, gin.H{"error": err})
		return
	}

	fmt.Println("Fetched data for: ", os.Getpid(), " Data: ", *order)

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
	lock := GetRouterLock(c)
	lock.Lock()
	defer lock.Unlock()

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

	task := func() {
		computeOrDelay()

		_, err := controller.InsertOrUpdate(c, &reqBody.Order, timestamp.IncrementTimestamp(c))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("Insertion ok. Now timestamp:", timestamp.GetTimestamp(c))
		reqBody.Order.Timestamp = timestamp.GetTimestamp(c) // ?
		SendSyncMsg(&reqBody.Order, GetPortList(c))
		fmt.Println("Send sync msg ok.")
	}

	fmt.Println("Add Task: Will insert to db. Priority:", timestamp.GetTimestamp(c))
	taskQueue := GetTaskQueue(c)
	heap.Push(taskQueue, &mypq.Item{Value: task, Priority: timestamp.GetTimestamp(c)})

	c.JSON(200, gin.H{"status": "ok"})
}

func RcvMsg(c *gin.Context) {
	lock := GetRouterLock(c)
	lock.Lock()
	defer lock.Unlock()

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
	local_timestamp := timestamp.GetTimestamp(c)
	fmt.Println("Received a sync msg from proc: ", reqBody.ProcID, " timestamp: ", order.Timestamp,
		"local_timestamp: ", local_timestamp)

	/**
	 *  Lamport's logical clock algorithm
		Here's how the algorithm works:

	    Each process maintains a logical clock, and whenever it performs a local event, its logical clock is incremented.
	    When a process sends a message, it attaches its logical clock value to the message. Upon receiving the message,
		the recipient process updates its own logical clock to be greater than the timestamp in the received message.
	    If two events have the same timestamp, the tie is broken by comparing process identifiers.
	*/
	if order.Timestamp == local_timestamp {
		fmt.Println("Tie: look at procId", reqBody.ProcID, " ", os.Getpid())
		if reqBody.ProcID > os.Getpid() {

			fmt.Println("Tie broken")
			task := mypq.Item{Value: func() {
				controller.InsertOrUpdate(c, &order, order.Timestamp)
			}}
			task.SetPriority(order.Timestamp)
			taskQueue := GetTaskQueue(c)

			heap.Push(taskQueue, &task)
			fmt.Println("Add Task: update entry to", order)
		}

	} else if order.Timestamp > local_timestamp {

		task := func() {
			println("I will update my local timestamp to: ", order.Timestamp)
			timestamp.SetTimestamp(c, order.Timestamp) // ?
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
	requestJSON := struct {
		ReqType string       `json:"req-type"`
		ProcID  int          `json:"proc-id"`
		Order   models.Order `json:"order"`
	}{
		ReqType: "sync",
		ProcID:  os.Getpid(),
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