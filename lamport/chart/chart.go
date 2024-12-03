package chart

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Create a WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connections []*websocket.Conn

// handleWebSocket handles WebSocket connections
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	connections = append(connections, conn)
	fmt.Println("New client connected")

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Connection error:", err)
			// Remove the closed connection from the list
			for i, c := range connections {
				if c == conn {
					connections = append(connections[:i], connections[i+1:]...)
					break
				}
			}
			break
		}
	}
}

type DynamicChart struct {
	DataCh chan DataPoint
}

func (chart *DynamicChart) SendDataPoint(point DataPoint) {
	chart.DataCh <- point
}

func Init() *DynamicChart {
	return &DynamicChart{
		DataCh: make(chan DataPoint, 4096),
	}
}

func (chart *DynamicChart) Main() {
	// Create a channel to simulate data stream
	// Listen and handle data reception
	go func() {
		for value := range chart.DataCh {
			// fmt.Printf("New data point: %v\n", value)
			for _, conn := range connections {
				err := conn.WriteJSON(value)
				if err != nil {
					fmt.Println("Error sending data:", err)
					conn.Close()
					for i, c := range connections {
						if c == conn {
							connections = append(connections[:i], connections[i+1:]...)
							break
						}
					}
				}
			}
		}
	}()

	// Set up HTTP routes
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// fmt.Println("Current directory:", http.Dir("."))
		http.ServeFile(w, req, "./chart/index.html")
	})

	fmt.Println("Visit http://localhost:8084 to view the real-time chart")
	if err := http.ListenAndServe(":8084", nil); err != nil {
		panic(err)
	}
}
