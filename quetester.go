package main

import (
	"./eventmgr"
	"./orderque"
	"fmt"
)

func main() {
	var que orderque.OrderQue_t

	event := make(chan (eventmgr.Event_t))

	eventmgr.CheckEvents(event)

	fmt.Println("init don")
	for true {

		newEvent := <-event
		if newEvent.EventType < eventmgr.FLOOR_SIGNAL {
			fmt.Printf("new order: %+v\n", newEvent)
			que.AddOrder(newEvent.Floor, newEvent.EventType)
			que.Print()
		} else {
			fmt.Printf("new floor: %d", newEvent.Floor)
			que.CompleteOrder(newEvent.Floor, 0)
			que.CompleteOrder(newEvent.Floor, 1)
			que.CompleteOrder(newEvent.Floor, 2)
			que.Print()
			fmt.Println("gue deleted")
		}

	}

}
