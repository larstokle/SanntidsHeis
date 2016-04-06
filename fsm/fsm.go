package fsm

import (
	"../driver"
	"time"
	//."../constants"
	"fmt"
	"math"
)

type State int

const (
	STATE_IDLE State = iota
	STATE_MOVING
	STATE_DOOR_OPEN
)

var states = [...]string{
	"STATE_IDLE",
	"STATE_MOVING",
	"STATE_DOOR_OPEN",
}

func (state State) String() string {
	return states[state]
}

type ElevatorState struct {
	fsmState    State
	floor       int
	dir         int
	destination int
}

func NewElevatorState(floorEvent chan int) ElevatorState {
	var elev ElevatorState
	//elev.floor = -1
	driver.RunDown()
	elev.floor = <-floorEvent
	elev.goToStateIdle()
	return elev
}

func (elev *ElevatorState) GetState() State {
	return elev.fsmState
}
func (elev ElevatorState) GetFloor() int {
	return elev.floor
}
func (elev ElevatorState) GetDir() int {
	return elev.dir
}
func (elev ElevatorState) GetDestination() int {
	return elev.destination
}

//Settes i en annen modul? Dette er bare en sign funksjon
func CalculateDir(destination int, currentFloor int) int {
	return int(math.Copysign(1, float64(destination-currentFloor)))
}

func (elev *ElevatorState) NewDestination(destination int) {
	fmt.Printf("new destination = %d\n", destination)
	elev.destination = destination
	if destination == elev.floor {
		elev.goToStateDoorOpen()
	} else {
		elev.goToStateMoving(CalculateDir(destination, elev.floor))
	}
}

func (elev *ElevatorState) NewFloorReached(newFloor int) {
	elev.floor = newFloor
	driver.SetFloorIndicator(elev.floor)
	if elev.floor == elev.destination {
		elev.destination = -1
		elev.goToStateDoorOpen()
	}
}

func (elev *ElevatorState) goToStateDoorOpen() {
	fmt.Println("Door Opening")
	driver.RunStop()
	elev.fsmState = STATE_DOOR_OPEN
	driver.SetDoorOpen(true)

	time.AfterFunc(time.Second*3, func() {
		driver.SetDoorOpen(false)
		elev.goToStateIdle()
	})

}

func (elev *ElevatorState) goToStateMoving(direction int) {
	fmt.Println("Starting to move")
	elev.dir = direction
	if direction == 1 {
		driver.RunUp()
	} else {
		driver.RunDown()
	}
	elev.fsmState = STATE_MOVING
}

func (elev *ElevatorState) goToStateIdle() {
	fmt.Println("Going idle")
	driver.RunStop()
	elev.fsmState = STATE_IDLE
}
