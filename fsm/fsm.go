package fsm

import(
	"../driver"
	"time"
	//."../constants"
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

type ElevatorState struct{
	fsmState State
	floor int
	dir int
	destination int
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

func (elev *ElevatorState) NewDestination(destination int){
	elev.destination = destination
	if destination == elev.floor{
		elev.goToStateDoorOpen()
	} else {
		elev.goToStateMoving(int( math.Copysign(1,float64(destination-elev.floor))))
	}
	/*
	if destination - floor > 0{
		elev.goToStateMoving(DIR_UP)
	} else if destination - floor == 0 {
		elev.goToStateDoorOpen()
	} else if destination - floor < 0 {
		elev.goToStateMoving(DIR_DOWN)
	}
	*/
}

func (elev *ElevatorState) NewFloorReached(newFloor int){
	elev.floor = newFloor
	driver.SetFloorIndicator(elev.floor)
	if elev.floor == elev.destination{
		elev.destination = -1
		elev.goToStateDoorOpen()
	}
}

func (elev *ElevatorState) goToStateDoorOpen() {
	driver.RunStop()
	elev.fsmState = STATE_DOOR_OPEN
	driver.SetDoorOpen(true)

	time.AfterFunc(time.Second * 3, func(){
		driver.SetDoorOpen(false)
		elev.goToStateIdle()
		})

}

func (elev *ElevatorState) goToStateMoving(direction int) {
	elev.dir = direction
	if direction == 1 {
		driver.RunUp()
	} else {
		driver.RunDown()
	}
}

func (elev *ElevatorState) goToStateIdle(){
	driver.RunStop()
	elev.fsmState = STATE_IDLE
}
