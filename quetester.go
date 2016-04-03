package main

import (
	. "./constants"
	"./eventmgr"
	"./orderque"
	"fmt"
)

func main() {
	var que orderque.OrderQue_t

	event := make(chan (eventmgr.Event_t))

	eventmgr.CheckEvents(event)

	fmt.Println("init done")
	for true {

		newEvent := <-event
		if newEvent.EventType < N_BUTTON_TYPES {
			fmt.Printf("new order: %+v\n", newEvent)
			que.AddOrder(newEvent.Floor, newEvent.EventType)
			que.Print()
		} else {
			fmt.Printf("new floor: %d", newEvent.Floor)
			que.RemoveOrder(newEvent.Floor, 0)
			que.RemoveOrder(newEvent.Floor, 1)
			que.RemoveOrder(newEvent.Floor, 2)
			que.Print()
			fmt.Println("que deleted")
		}

	}

}
