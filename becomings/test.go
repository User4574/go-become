package main

import "fmt"

type becomer string

func (b becomer) Become() {
	fmt.Println("Success!")
}

var Become becomer
