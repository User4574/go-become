package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"plugin"
	"sync"
)

const (
	LISTEN = "localhost:3333"
)

var (
	wait sync.WaitGroup
)

type Become interface {
	Become()
}

func main() {
	l, err := net.Listen("tcp", LISTEN)
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting: ", err.Error())
	} else {
		wait.Add(1)
		go handleRequest(conn)
	}
	wait.Wait()
}

func handleRequest(conn net.Conn) {
	defer wait.Done()
	defer conn.Close()
	file, err := ioutil.TempFile("", "become.*")
	if err != nil {
		fmt.Println("Error creating temporary file: ", err.Error())
		return
	}
	_, err = io.Copy(file, conn)
	if err != nil {
		fmt.Println("Error copying: ", err.Error())
		return
	}
	err = file.Close()
	if err != nil {
		fmt.Println("Error closing temporary file: ", err.Error())
		return
	}
	wait.Add(1)
	go becomePlugin(file.Name())
}

func becomePlugin(filename string) {
	defer wait.Done()
	defer os.Remove(filename)

	plugin, err := plugin.Open(filename)
	if err != nil {
		fmt.Println("Error opening becoming: ", err.Error())
		return
	}

	symBecoming, err := plugin.Lookup("Become")
	if err != nil {
		fmt.Println("Error looking up becoming: ", err.Error())
		return
	}

	var becoming Become
	becoming, ok := symBecoming.(Become)
	if !ok {
		fmt.Println("Becoming not of correct type")
		return
	}

	becoming.Become()
}
