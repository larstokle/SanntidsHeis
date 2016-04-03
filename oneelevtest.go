package main 

import (
	"./fsm"
	"./eventmgr"
	"./orderque"
	."./constants"
	"fmt"
	"time"
)

func main() {
	var que orderque.OrderQue_t
	var elevatorData fsm.ElevatorState
	event := make(chan (eventmgr.Event_t))
	
	eventmgr.CheckEvents(event)
	fmt.Println("State is", elevatorData.GetState())

	//elevatorData.Init()

	for true{
		select{
		case newEvent := <- event:
			fmt.Printf("event: %+v\nState is %+v, destination is %d\n", newEvent, elevatorData.GetState(), elevatorData.GetDestination())
			if newEvent.EventType == FLOOR_SIGNAL{
				elevatorData.NewFloorReached(newEvent.Floor)

			} else {
				que.AddOrder(newEvent.Floor, newEvent.EventType)
			}

		default:
			switch elevatorData.GetState() {
				case fsm.STATE_DOOR_OPEN:
					fmt.Println("im open!!!")
					que.RemoveOrder(elevatorData.GetFloor(),UP)
					que.RemoveOrder(elevatorData.GetFloor(),DOWN)
					que.RemoveOrder(elevatorData.GetFloor(),CMD)
					
				case fsm.STATE_IDLE:
					if !que.IsEmpty(){
						newDestination := que.EarliestOrderInside().GetFloor()
						fmt.Printf("EarliestOrderInside = %+v\n", newDestination)
						
						elevatorData.NewDestination(newDestination)
						
					}
				default:
					fmt.Println("im... IDK!!!")
			}
			time.Sleep(time.Millisecond*20)
		}

	}
}