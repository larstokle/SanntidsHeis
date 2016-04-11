package fsm

import (
	"../driver"
	"time"
	"fmt"
	"math"
	"../eventmgr"

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

type Direction int

const(
	DIR_DOWN = -1
	DIR_STOP = 0
	DIR_UP   = 1
)

const UNDEFINED_DESTINATION = -1

func (state State) String() string {
	return states[state]
}

type ElevatorState struct {
	fsmState    State
	floor       int
	dir         Direction
	destination int
}

func NewElevator() *ElevatorState {
	driver.Init()
	var elev ElevatorState
	elev.destination = UNDEFINED_DESTINATION
	
	floorEvent  := eventmgr.CheckFloorSignal()
	driver.RunDown()
	elev.NewFloorReached(<-floorEvent)
	driver.RunStop()

	elev.goToStateIdle()

	go func(){
		for {
			fmt.Println("Loop!" )
			newFloor := <- floorEvent
			elev.NewFloorReached(newFloor)
		}
	}()
	fmt.Println("did it Loop?" )

	return &elev
}

func (elev *ElevatorState) State() State {
	return elev.fsmState
}
func (elev ElevatorState) Floor() int {
	return elev.floor
}
func (elev ElevatorState) Dir() Direction {
	return elev.dir
}
func (elev ElevatorState) Destination() int {
	return elev.destination
}

//Settes i en annen modul? Dette er bare en sign funksjon
func CalculateDir(destination int, currentFloor int) Direction {
	return Direction(math.Copysign(1, float64(destination-currentFloor)))
}

func (elev *ElevatorState) NewDestination(destination int) {
	fmt.Printf("new destination = %d\n", destination)
	elev.destination = destination
	if destination == elev.floor {
		elev.goToStateDoorOpen()
	} else if elev.fsmState == STATE_IDLE{
		elev.goToStateMoving(CalculateDir(destination, elev.floor))
	}
}

func (elev *ElevatorState) NewFloorReached(newFloor int) {
	elev.floor = newFloor
	fmt.Println(elev.floor)
	driver.SetFloorIndicator(elev.floor)
	if elev.floor == elev.destination {
		elev.destination = UNDEFINED_DESTINATION
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
		if elev.destination == UNDEFINED_DESTINATION{
			elev.goToStateIdle()
		} else {
			elev.goToStateMoving(CalculateDir(elev.destination, elev.floor))
		}
	})

}

func (elev *ElevatorState) goToStateMoving(direction Direction) {
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

func (elev *ElevatorState) NeedNewDestination() bool{
	return elev.fsmState == STATE_IDLE || elev.fsmState == STATE_DOOR_OPEN
}
