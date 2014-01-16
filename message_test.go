package main

import (
	"fmt"
	"testing"
)

func TestPackUnpack(t *testing.T) {
	fmt.Println("called")
	data := []byte{103, 85, 1, 32, 0, 1, 0, 0, 0, 0, 0, 1, 3, 119, 119, 119, 5, 98, 97, 105, 100, 117, 3, 99, 111, 109, 0, 0, 1, 0, 1, 0,
		0, 41, 16, 0, 0, 0, 0, 0, 0, 0}
	var msg Message
	fmt.Println(msg.Unpack(data))
	fmt.Println(msg)
}
