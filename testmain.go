package main

import (
	"./eventmgr"
	"./fsm"
	"fmt"
)
var elev fsm.ElevatorState

func main() {
	fmt.Println("Hello")
	fmt.Println("State is", elev.GetState())

	event := make(chan (eventmgr.Event_t))

	eventmgr.CheckEvents(event)

	for true {
		testEvent := <-event
		fmt.Printf("event: %+v\nState is %+v\n", testEvent, elev.GetState())
		if testEvent.EventType == 3{
			elev.NewFloorReached(testEvent.Floor)
		} 
	}
}
