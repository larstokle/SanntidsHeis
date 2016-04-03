package main 

import (
	"./fsm"
	"./eventmgr"
	"./orderque"
	."./constants"
	"fmt"
)

func main() {
	var que orderque.OrderQue_t
	var elevatorData fsm.ElevatorState
	event := make(chan (eventmgr.Event_t))
	
	eventmgr.CheckEvents(event)
	fmt.Println("State is", elevatorData.GetState())

	for true{
		newEvent := <- event
		fmt.Printf("event: %+v\nState is %+v\n", newEvent, elevatorData.GetState())
		if newEvent.EventType == FLOOR_SIGNAL{
			elevatorData.NewFloorReached(newEvent.Floor)

		} else {
			que.AddOrder(newEvent.Floor, newEvent.EventType)
		}

		switch elevatorData.GetState() {
			case fsm.STATE_DOOR_OPEN:
				que.RemoveOrder(elevatorData.GetFloor(),UP)
				que.RemoveOrder(elevatorData.GetFloor(),DOWN)
				que.RemoveOrder(elevatorData.GetFloor(),CMD)
				
			case fsm.STATE_IDLE:
				if !que.IsEmpty(){
					newOrder := que.EarliestOrderInside()
					elevatorData.NewDestination(newOrder.GetFloor())
				}
		}
		que.Print()
		
	}
}