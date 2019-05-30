package main

import (
	"fmt"
	"github.com/grafov/bcast"
	"net"
	"os"
)

const (
	LISTEN = "localhost:3334"
)

type becoming struct {
	Clients *bcast.Group
}

func (b *becoming) Become() {
	b.Clients = bcast.NewGroup()
	go b.Clients.Broadcast(0)

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
			go newClient(conn, b)
		}
	}
}

func clientReader(conn net.Conn, ch chan []byte) {
	buf := make([]byte, 256)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			return
		}
		ch <- buf
	}
}

func busReader(member *bcast.Member, ch chan []byte) {
	for {
		ch <- member.Recv().([]byte)
	}
}

func newClient(conn net.Conn, b *becoming) {
	defer conn.Close()
	member := b.Clients.Join()
	defer b.Clients.Leave(member)

	in := make(chan []byte)
	out := make(chan []byte)
	go clientReader(conn, in)
	go busReader(member, out)

	for {
		select {
		case msg := <-in:
			member.Send(msg)
		case msg := <-out:
			conn.Write(msg)
		}
	}
}

var Become becoming
