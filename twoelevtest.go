package main

import (
	. "./constants"
	"./eventmgr"
	"./fsm"
	"./orderque"
	"fmt"
	"math"
	"time"
	"./transactionMgr"
)

func main() {
	var que orderque.OrderQue_t

	floorEvent := make(chan int)
	btnEvent := make(chan eventmgr.Event_t)
	eventmgr.CheckEvents(btnEvent, floorEvent)

	elevatorData := fsm.NewElevatorState(floorEvent)
	fmt.Println("State is", elevatorData.GetState())
	toTrans, fromTrans := transactionMgr.StartTransactionManager()
	for true {
		select {
		case newEvent := <-floorEvent:
			fmt.Printf("New floor: %+v, State is %+v, destination is %d\n", newEvent, elevatorData.GetState(), elevatorData.GetDestination())
			elevatorData.NewFloorReached(newEvent)

		case newEvent := <-btnEvent:
			fmt.Printf("New btn pushed: %+v, State is %+v, destination is %d\n", newEvent, elevatorData.GetState(), elevatorData.GetDestination())
			if newEvent.EventType != CMD{
				toTrans <- newEvent
			} else {
				//when command
			}
			// que.AddOrder(newEvent.Floor, newEvent.EventType)
			// dir := elevatorData.GetDir()
			// closest := int(math.Min(float64(newEvent.Floor*dir), float64(elevatorData.GetDestination()*dir))) * dir
			// if dir == fsm.CalculateDir(newEvent.Floor, elevatorData.GetFloor()) && closest != elevatorData.GetDestination() && (dir == DIR_UP && newEvent.EventType == UP || dir == DIR_DOWN && newEvent.EventType == DOWN || newEvent.EventType == CMD) {
			// 	elevatorData.NewDestination(newEvent.Floor)
			// }
		case newUnknownTrans := <-fromTrans:
			switch newTrans:= newUnknownTrans.(type){
			case eventmgr.Event_t:
				fmt.Printf("New REMOTE btn pushed: %+v, State is %+v, destination is %d\n", newTrans, elevatorData.GetState(), elevatorData.GetDestination())
				que.AddOrder(newTrans.Floor, newTrans.EventType)
				dir := elevatorData.GetDir()
				closest := int(math.Min(float64(newTrans.Floor*dir), float64(elevatorData.GetDestination()*dir))) * dir
				if dir == fsm.CalculateDir(newTrans.Floor, elevatorData.GetFloor()) && closest != elevatorData.GetDestination() && (dir == DIR_UP && newTrans.EventType == UP || dir == DIR_DOWN && newTrans.EventType == DOWN || newTrans.EventType == CMD) {
					elevatorData.NewDestination(newTrans.Floor)
				}
			default:
				fmt.Printf("unknown from trans: %+v \n", newTrans)
			}


		default:
			switch elevatorData.GetState() {
			case fsm.STATE_DOOR_OPEN:
				que.RemoveOrder(elevatorData.GetFloor(), UP)
				que.RemoveOrder(elevatorData.GetFloor(), DOWN)
				que.RemoveOrder(elevatorData.GetFloor(), CMD)

			case fsm.STATE_IDLE:
				if !que.IsEmpty() {

					newDestination := calculateNewDestination(que, elevatorData)

					elevatorData.NewDestination(newDestination)
				}
			default:

			}
			time.Sleep(time.Millisecond * 20)
		}

	}
}

func calculateNewDestination(que orderque.OrderQue_t, elevator fsm.ElevatorState) int {
	newCMDDestination := que.NextOrderOfTypeInDir(elevator.GetFloor(), elevator.GetDir(), CMD)
	if newCMDDestination != -1 {
		fmt.Printf("CMD Destination in same dir as before = %+v\n", newCMDDestination)
		return newCMDDestination
	}

	newEarliestDestination := que.EarliestNonAssignedOrder().GetFloor()
	dir := fsm.CalculateDir(newEarliestDestination, elevator.GetFloor())
	fmt.Printf("EarliestOrderInside = %+v\n", newEarliestDestination)

	if dir != elevator.GetDir() {
		fmt.Println("Change of dir")
		newCMDDestination = que.NextOrderOfTypeInDir(elevator.GetFloor(), dir, CMD)
	}
	fmt.Printf("CMD Destination = %+v\n", newCMDDestination)

	var btn int
	if dir == DIR_UP {
		btn = UP
	} else if dir == DIR_DOWN {
		btn = DOWN
	}
	newBTNDestination := que.NextOrderOfTypeInDir(elevator.GetFloor(), dir, btn)
	fmt.Printf("BTN Destination = %+v\n", newBTNDestination)

	newDestination := newEarliestDestination

	if newBTNDestination != -1 {
		newDestination = int(math.Min(float64(newDestination*dir), float64(newBTNDestination*dir))) * dir
	}

	if newCMDDestination != -1 {
		newDestination = int(math.Min(float64(newDestination*dir), float64(newCMDDestination*dir))) * dir
	}
	return newDestination
}
