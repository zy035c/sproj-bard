package ft_handshake

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func Main() {

	var localAddr string = "127.0.0.1:9000"
	var remoteAddrList []string

	// Read port
	if localAddr_ := os.Getenv("G_ADDR"); localAddr != "" {
		fmt.Println("Using localAddr environment variable:", localAddr_)
		localAddr = localAddr_
	}

	// Read port list
	if gPs := os.Getenv("G_REMOTE_ADDR_LIST"); gPs != "" {
		fmt.Println("Using G_REMOTE_ADDR_LIST environment variable:", gPs)
		remoteAddrList = strings.Split(gPs, ",")
		for _, p := range remoteAddrList {
			fmt.Println("Other process addr:", p)
		}
	}

	thread := New(
		localAddr, remoteAddrList, func() {
			println("Yeah Computation")
		},
	)

	for {
		thread.Start()
		time.Sleep(50 * time.Millisecond)
	}
}
