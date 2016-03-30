package main

import (
	"./eventmgr"
	"fmt"
)

func main() {
	fmt.Println("Hello")

	event := make(chan (eventmgr.Event_t))

	eventmgr.CheckEvents(event)

	for true {
		fmt.Printf("event: %+v\n", <-event)
	}
}
