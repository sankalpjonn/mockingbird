package main

import (
	"fmt"

	"gopkg.in/fatih/color.v1"
)

func main() {
	client := New()
	err, response := client.callServer()
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		fmt.Println(red("\nError: %s", err))
	} else {
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Println(green("\n" + response))
	}
}
