package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
)

const (
	LISTEN = "localhost:3334"
)

type becoming struct {
	Fibmutex sync.Mutex
	Last     [2]int
}

func (b *becoming) Become() {
	b.Last = [2]int{0, 0}
	l, err := net.Listen("tcp", LISTEN)
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
		} else {
			go handleRequest(conn, b)
		}
	}
}

func handleRequest(conn net.Conn, b *becoming) {
	defer conn.Close()
	b.Fibmutex.Lock()
	defer b.Fibmutex.Unlock()
	if b.Last == [2]int{0, 0} {
		conn.Write([]byte(strconv.Itoa(0)))
		b.Last[1] = 1
	} else {
		conn.Write([]byte(strconv.Itoa(b.Last[1])))
		sum := b.Last[0] + b.Last[1]
		b.Last[0] = b.Last[1]
		b.Last[1] = sum
	}
}

var Become becoming
