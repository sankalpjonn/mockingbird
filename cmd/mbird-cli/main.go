package main

import (
	"fmt"
)

func main() {
	client := New()
	response := client.callServer()
	fmt.Println(response)
}
