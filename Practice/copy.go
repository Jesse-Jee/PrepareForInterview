package main

import "fmt"

func main() {
	s := "hello world"
	b := make([]byte, 11)
	copy(b, s)
	fmt.Println(string(b))
}
