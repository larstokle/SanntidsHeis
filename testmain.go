package main

import (
	"./eventmgr"
	"./fsm"
	"fmt"
)

func main() {
	fmt.Println("Hello")
	fmt.Println("State is", fsm.GetState())

	event := make(chan (eventmgr.Event_t))

	eventmgr.CheckEvents(event)

	for true {
		fmt.Printf("event: %+v\n", <-event)
	}
}
