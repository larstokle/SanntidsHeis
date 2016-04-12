package main

import (
	. "./globals"
	"./eventmgr"
	"./fsm"
	"./orderque"
	"fmt"
	"math"
	"time"
)

func main() {
	var que orderque.OrderQue_t

	floorEvent := make(chan int)
	btnEvent := make(chan eventmgr.Event_t)
	eventmgr.CheckEvents(btnEvent, floorEvent)

	elevatorData := fsm.NewElevatorState(floorEvent)
	fmt.Println("State is", elevatorData.GetState())

	for true {
		select {
		case newEvent := <-floorEvent:
			fmt.Printf("New floor: %+v, State is %+v, destination is %d\n", newEvent, elevatorData.GetState(), elevatorData.GetDestination())
			elevatorData.NewFloorReached(newEvent)

		case newEvent := <-btnEvent:
			fmt.Printf("New btn pushed: %+v, State is %+v, destination is %d\n", newEvent, elevatorData.GetState(), elevatorData.GetDestination())
			que.AddOrder(newEvent.Floor, newEvent.EventType)
			dir := elevatorData.GetDir()
			closest := int(math.Min(float64(newEvent.Floor*dir), float64(elevatorData.GetDestination()*dir))) * dir
			if dir == fsm.calculateDir(newEvent.Floor, elevatorData.GetFloor()) && closest != elevatorData.GetDestination() && (dir == DIR_UP && newEvent.EventType == UP || dir == DIR_DOWN && newEvent.EventType == DOWN || newEvent.EventType == CMD) {
				elevatorData.NewDestination(newEvent.Floor)
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
	newCMDDestination := que.NearestOrderOfTypeInDir(elevator.GetFloor(), elevator.GetDir(), CMD)
	if newCMDDestination != -1 {
		fmt.Printf("CMD Destination in same dir as before = %+v\n", newCMDDestination)
		return newCMDDestination
	}

	newEarliestDestination := que.EarliestOrderInside().GetFloor()
	dir := fsm.calculateDir(newEarliestDestination, elevator.GetFloor())
	fmt.Printf("EarliestOrderInside = %+v\n", newEarliestDestination)

	if dir != elevator.GetDir() {
		fmt.Println("Change of dir")
		newCMDDestination = que.NearestOrderOfTypeInDir(elevator.GetFloor(), dir, CMD)
	}
	fmt.Printf("CMD Destination = %+v\n", newCMDDestination)

	var btn int
	if dir == DIR_UP {
		btn = UP
	} else if dir == DIR_DOWN {
		btn = DOWN
	}
	newBTNDestination := que.NearestOrderOfTypeInDir(elevator.GetFloor(), dir, btn)
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
