package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	inputString := "Hello, OTUS!"
	fmt.Println(reverse.String(inputString))
}
