package ft_handshake

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"golang.org/x/sys/unix"
)

type FaultTolerantHandshake struct {
	localAddr      string
	remoteAddrList []string
	computation    func()
}

/*
Contructor for the FaultTolerantHandshake
*/
func New(
	localAddr string,
	remoteAddrList []string,
	cmp func(),
) *FaultTolerantHandshake {
	return &FaultTolerantHandshake{
		localAddr:      localAddr,
		remoteAddrList: remoteAddrList,
		computation:    cmp,
	}
}

func (fth *FaultTolerantHandshake) Start() {
	fmt.Println("[Start] step 1")
step2:
	cp := fth.chooseRandomNeighbor()
	fth.sendOne(cp)

step4:
	msg, w := fth.receive()
	if w == cp && msg == 1 {
		// there is a handshake between p and c(p)
		fmt.Println("[HANDSHAKE] succeed between", fth.localAddr, w)
		fth.computation()
		return
	} else {
		if w == cp && msg == 0 {
			// c(p) has chosen a node different from p
			goto step2
		} else {
			fth.sendZero(w)
			goto step4
		}
	}
}

func fdToConn(fd int) (net.Conn, error) {
	f := os.NewFile(uintptr(fd), "any_name")
	defer f.Close()
	return net.FileConn(f)
}

func (fth *FaultTolerantHandshake) chooseRandomNeighbor() string {
	addr := fth.remoteAddrList[rand.Intn(len(fth.remoteAddrList))]
	return addr
}

func (fth *FaultTolerantHandshake) sendOne(cp string) {
	conn, err := net.Dial("tcp", cp)
	if err != nil {
		fmt.Println("Error dialing to c(p):", err.Error())
		return
	}

	defer conn.Close()

	_, err = conn.Write([]byte(strconv.Itoa(1)))
	if err != nil {
		fmt.Println("Error writing to TCP:", err)
		return
	}
	fmt.Println("Sent 1 successfully!")
}

func (fth *FaultTolerantHandshake) receive() (int, string) {
	start := time.Now()

	fds := []unix.PollFd{}
	for _, listenOn := range fth.remoteAddrList {
		connRcv, err := net.Dial("tcp", listenOn)
		// possible timeout
		if err != nil {
			continue
		}
		defer connRcv.Close()

		fd, err := connRcv.(*net.TCPConn).File()
		if err != nil {
			fmt.Println("Error getting file descriptor:", err.Error())
			continue
		}

		// append the file descriptor to the list
		fds = append(
			fds,
			unix.PollFd{
				Fd:     int32(fd.Fd()),
				Events: unix.POLLIN,
			},
		)
	}

	// poll the file descriptor
	for {
		n, err := unix.Poll(fds, 1000)
		if err != nil {
			fmt.Println("Error polling:", err.Error())
			continue
		}

		for _, fd := range fds {
			if fd.Revents&unix.POLLIN == unix.POLLIN {
				fmt.Printf(
					"n=%d err=%v delay=%v flags=%016b (POLLIN=%t)\n",
					n, err, time.Since(start), fd.Revents,
					fd.Revents&unix.POLLIN != 0,
				)
				ptr := int(fd.Fd)

				// read data from the conn
				connRcv, err := fdToConn(ptr)

				if err != nil {
					fmt.Println("Error polling:", err.Error())
					continue
				}

				buf := make([]byte, 1024)
				n, err := connRcv.Read(buf)
				if err != nil {
					fmt.Println("Error reading from connRcv:", err.Error())
					continue
				}
				fmt.Printf("Received data: %s\n", buf[:n])

				return int(buf[0]), connRcv.RemoteAddr().String()
			}
		}
	}
}

func (fth *FaultTolerantHandshake) sendZero(w string) {
	conn, err := net.Dial("tcp", w)
	if err != nil {
		fmt.Println("Error dialing to w:", err.Error())
		return
	}

	defer conn.Close()

	_, err = conn.Write([]byte(strconv.Itoa(0)))
	if err != nil {
		fmt.Println("Error writing to TCP:", err)
		return
	}
	fmt.Println("Sent 0 successfully!")
}
