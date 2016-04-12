package fsm

import (
	"driver"
	"time"
	"fmt"
	"eventmgr"
	."globals"
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


const NO_DESTINATION = -1
const INF_COST = 255

func (state State) String() string {
	return states[state]
}

type ElevatorState struct {
	fsmState    State
	floor       int
	dir         Direction_t
	destination int
	OrderDone chan int
	doorTimer *time.Timer
}

func NewElevator() *ElevatorState {
	driver.Init()
	var elev ElevatorState
	elev.destination = NO_DESTINATION
	elev.OrderDone = make(chan int, 1)
	floorEvent  := eventmgr.CheckFloorSignal()
	driver.RunDown()
	elev.NewFloorReached(<-floorEvent)
	driver.RunStop()

	elev.goToStateIdle()

	go func(){
		for {
			newFloor := <- floorEvent
			elev.NewFloorReached(newFloor)
			
		}
	}()


	return &elev
}

func (elev *ElevatorState) State() State {
	return elev.fsmState
}
func (elev ElevatorState) Floor() int {
	return elev.floor
}
func (elev ElevatorState) Dir() Direction_t {
	return elev.dir
}
func (elev ElevatorState) Destination() int {
	return elev.destination
}



func (elev *ElevatorState) NewDestination(destination int) {
	fmt.Printf("fsm: new destination = %d\n", destination)
	elev.destination = destination
	if destination == elev.floor {//&& elev.fsmState != STATE_DOOR_OPEN{ //forslag til nytt: reset the afterfunc timer in goToStateDoorOpen() if STATE_DOOR_OPEN
		elev.destinationReaced()
	} else if elev.fsmState == STATE_IDLE{
		elev.goToStateMoving(calculateDir(destination, elev.floor))
	}
}

func (elev *ElevatorState) destinationReaced(){
	fmt.Println("fsm: Destination reached")
	elev.destination = NO_DESTINATION
	elev.goToStateDoorOpen()
}


func (elev *ElevatorState) NewFloorReached(newFloor int) {
	fmt.Printf("fsm: New floor reached= %d\n", newFloor)
	elev.floor = newFloor
	driver.SetFloorIndicator(elev.floor)
	if elev.floor == elev.destination {
		elev.destinationReaced()
	}
}

func (elev *ElevatorState) goToStateDoorOpen() {
	fmt.Printf("fsm: Door Opening in floor = %d\n",elev.floor)

	if elev.fsmState == STATE_DOOR_OPEN{
		fmt.Println("fsm: Door alredy open")
		elev.doorTimer.Reset(time.Second*3)
		elev.OrderDone <- elev.floor
		return
	}

	driver.RunStop()
	elev.fsmState = STATE_DOOR_OPEN
	elev.OrderDone <- elev.floor
	driver.SetDoorOpen(true)

	doorClose := func() {
		fmt.Printf("fsm: Door closing\n")
		driver.SetDoorOpen(false)
		if elev.destination == NO_DESTINATION{
			elev.goToStateIdle()
		} else {
			elev.goToStateMoving(calculateDir(elev.destination, elev.floor))
		}
	}
	 
	elev.doorTimer = time.AfterFunc(time.Second*3, doorClose)
	
	
}

func (elev *ElevatorState) goToStateMoving(direction Direction_t) {
	fmt.Printf("fsm: Starting to move in dir = %d against destination = %d\n", direction, elev.destination)
	switch direction{
	case DIR_UP:
		driver.RunUp()
	case DIR_DOWN:
		driver.RunDown()
	default:
		fmt.Printf("fsm: unknown direction")
		return
	}
	elev.dir = direction
	elev.fsmState = STATE_MOVING
	
	
}

func (elev *ElevatorState) goToStateIdle() {
	fmt.Println("fsm: Going idle in floor = ", elev.floor)
	driver.RunStop()
	elev.fsmState = STATE_IDLE
}

func (elev *ElevatorState) NeedNewDestination() bool{
	return elev.destination == NO_DESTINATION
}

func (elev *ElevatorState) GetCost(order Button_t) int{
	newDistance := distance(order.Floor,elev.floor)

	if elev.destination == NO_DESTINATION{
		return newDistance 
	}


	currentDistance := distance(elev.destination, elev.floor)
	
	directionToOrder := calculateDir(order.Floor, elev.floor)
	orderIsInSameDir := directionToOrder == elev.dir || elev.dir == DIR_STOP

	if newDistance > currentDistance || !orderIsInSameDir{
		return INF_COST
	}

	
	buttonDir := buttonTypeToDirection(order.ButtonType)

	if order.ButtonType == CMD{
		return newDistance
	} else if buttonDir == elev.dir || elev.destination == order.Floor{
		return newDistance + N_FLOORS
	}else{
		fmt.Printf("fsm: Unhandlet GetCost! ELevator = %+v, order to get cost for = %+v\n",*elev, order)
		return INF_COST
	}
}


func buttonTypeToDirection(buttonType int) Direction_t{
	switch buttonType{
	case UP:
		return DIR_UP
	case DOWN:
		return DIR_DOWN
	default:
		return DIR_STOP
	}
}

