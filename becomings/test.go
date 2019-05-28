package main

import "fmt"

type becoming string

func (b becoming) Become() {
	fmt.Println("Success!")
}

var Become becoming
