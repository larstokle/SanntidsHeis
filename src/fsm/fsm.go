package fsm

import (
	"driver"
	"eventmgr"
	"fmt"
	. "globals"
	"sync/atomic"
	"sync"
	"time"
	"os"
)

type State_t int32

const (
	STATE_IDLE State_t = iota
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
const MOTOR_ERROR_SEC = 5

func (state State_t) String() string {
	return states[state]
}

type ElevatorState struct {
	fsmState    State_t
	floor       int32
	dir         Direction_t
	destination int32
	OrderDone   chan int
	doorTimer   *time.Timer
	doorTimerMutex sync.Mutex
	motorErrorTimer *time.Timer
}

func NewElevator() *ElevatorState {
	var elev ElevatorState
	elev.setDestination(NO_DESTINATION)
	elev.OrderDone = make(chan int, 1)
	floorEvent := eventmgr.CheckFloorSignal()
	driver.RunDown()
	elev.motorErrorTimer = time.AfterFunc(time.Second*MOTOR_ERROR_SEC, motorError)
	elev.NewFloorReached(<-floorEvent)
	driver.RunStop()
	elev.motorErrorTimer.Stop()

	elev.goToStateIdle()

	go func() {
		for {
			newFloor := <-floorEvent
			elev.NewFloorReached(newFloor)
		}
	}()
	fmt.Printf("fsm: init done at floor %d, NewElevator returned\n\n", elev.Floor())

	return &elev
}

func (elev *ElevatorState) State() State_t {
	return State_t(atomic.LoadInt32((*int32)(&elev.fsmState)))
}

func (elev *ElevatorState) setState(state State_t) {
	atomic.StoreInt32((*int32)(&elev.fsmState), int32(state))
}

func (elev *ElevatorState) Floor() int { //BURDE RETURNERE int32 ?
	return int(atomic.LoadInt32((*int32)(&elev.floor)))
}

func (elev *ElevatorState) setFloor(floor int){
	atomic.StoreInt32((*int32)(&elev.floor),int32(floor))
}

func (elev *ElevatorState) Dir() Direction_t {
	return Direction_t(atomic.LoadInt32((*int32)(&elev.dir)))
}

func (elev *ElevatorState) setDir(dir Direction_t){
	atomic.StoreInt32((*int32)(&elev.dir),int32(dir))
}

func (elev *ElevatorState) Destination() int { //BURDE RETURNERE int32 ?
	return int(atomic.LoadInt32((*int32)(&elev.destination)))
}

func (elev *ElevatorState) setDestination(dest int){
	atomic.StoreInt32((*int32)(&elev.destination),int32(dest))
}


func (elev *ElevatorState) NewDestination(destination int) {
	if(DEBUG_FSM){fmt.Printf("fsm: new destination = %d\n", destination)}

	if destination < FIRST_FLOOR || destination > TOP_FLOOR{
		destination = NO_DESTINATION
		return
	}

	elev.setDestination(destination)
	if destination == elev.Floor() {
		elev.destinationReaced()
	} else if elev.State() == STATE_IDLE {
		elev.goToStateMoving(calculateDir(destination, elev.Floor()))
	}
}

func (elev *ElevatorState) destinationReaced() {
	if(DEBUG_FSM){fmt.Println("fsm: Destination reached")}
	elev.setDestination(NO_DESTINATION)
	elev.goToStateDoorOpen()
}

func (elev *ElevatorState) NewFloorReached(newFloor int) {
	if(DEBUG_FSM){fmt.Printf("fsm: New floor reached= %d\n", newFloor)}
	elev.setFloor(newFloor)
	driver.SetFloorIndicator(int(elev.floor))
	elev.motorErrorTimer.Reset(time.Second*MOTOR_ERROR_SEC)
	if elev.Floor() == elev.Destination() {
		elev.destinationReaced()
	}else if elev.Destination() == NO_DESTINATION{
		elev.goToStateIdle()
	}
}

func (elev *ElevatorState) goToStateDoorOpen() {
	if(DEBUG_FSM){fmt.Printf("fsm: Door Opening in floor = %d\n", elev.floor)}

	if elev.State() == STATE_DOOR_OPEN {
		if(DEBUG_FSM){fmt.Println("fsm: Door alredy open")}
		elev.doorTimerMutex.Lock()
		elev.doorTimer.Reset(time.Second * 3)
		elev.doorTimerMutex.Unlock()
		
		elev.OrderDone <- int(elev.floor)
		
		return
	}

	driver.RunStop()
	elev.motorErrorTimer.Stop()
	elev.setState(STATE_DOOR_OPEN)
	
	elev.OrderDone <- int(elev.floor)
	

	driver.SetDoorOpen(true)

	doorClose := func() {
		if(DEBUG_FSM){fmt.Printf("fsm: Door closing, currentDestination = %d\n",elev.Destination())}
		driver.SetDoorOpen(false)
		if elev.Destination() == NO_DESTINATION {
			elev.goToStateIdle()
		} else {
			elev.goToStateMoving(calculateDir(elev.Destination(), elev.Floor()))
		}
	}
	elev.doorTimerMutex.Lock()
	elev.doorTimer = time.AfterFunc(time.Second*3, doorClose)
	elev.doorTimerMutex.Unlock()
}

func (elev *ElevatorState) goToStateMoving(direction Direction_t) {
	if(DEBUG_FSM){fmt.Printf("fsm: Starting to move in dir = %d against destination = %d\n\n\n", direction, elev.destination)}
	switch direction {
	case DIR_UP:
		driver.RunUp()
		elev.motorErrorTimer.Reset(time.Second*MOTOR_ERROR_SEC)
	case DIR_DOWN:
		driver.RunDown()
		elev.motorErrorTimer.Reset(time.Second*MOTOR_ERROR_SEC)
	default:
		if(DEBUG_FSM){fmt.Printf("fsm: unknown direction")}
		return
	}
	elev.setDir(direction)
	elev.setState(STATE_MOVING)

}

func (elev *ElevatorState) goToStateIdle() {
	if(DEBUG_FSM){fmt.Println("fsm: Going idle in floor = ", elev.floor)}
	driver.RunStop()
	elev.motorErrorTimer.Stop()
	elev.setState(STATE_IDLE)
}

func (elev *ElevatorState) GetCost(order Button_t) int {
	if order == NONVALID_BUTTON || order == NONVALID_BUTTON2 || order.Floor < FIRST_FLOOR || order.Floor > TOP_FLOOR || order.ButtonType < 0 || order.ButtonType >= N_BUTTON_TYPES{
		return INF_COST
	}

	newDistance := distance(order.Floor, elev.Floor())

	if elev.Destination() == NO_DESTINATION {
		return newDistance
	}

	currentDistance := distance(elev.Destination(), elev.Floor())

	directionToOrder := calculateDir(order.Floor, elev.Floor())
	orderIsInSameDir := directionToOrder == elev.Dir() || elev.Dir() == DIR_STOP

	if newDistance > currentDistance || !orderIsInSameDir {
		return INF_COST
	}

	buttonDir := buttonTypeToDirection(order.ButtonType)

	if order.ButtonType == CMD { // TROR DET ER EN BUG HER, SWITCHER MELLOM DOWN, OG CMD I SAMME ETASJE
		return newDistance
	} else if buttonDir == elev.Dir() {
		return newDistance //+ N_FLOORS
	} else if buttonDir != elev.Dir() {
		return INF_COST
	} else {
		if(DEBUG_FSM){fmt.Printf("ERROR!! fsm: Unhandlet GetCost! ELevator = %+v, order to get cost for = %+v\n", *elev, order)}
		return INF_COST
	}
}

func buttonTypeToDirection(buttonType int) Direction_t {
	switch buttonType {
	case UP:
		return DIR_UP
	case DOWN:
		return DIR_DOWN
	default:
		return DIR_STOP
	}
}

func motorError(){
	driver.RunStop()
	fmt.Printf("\n\n================== ERROR, ELEVATOR MOTOR NOT WORKING PROPERLY ==========\n\n")
	if floor := driver.GetFloorSignal(); floor != -1{
		fmt.Printf("Elevator at floor %d, door opened\n", floor)
	}
	fmt.Printf("took %d seconds to reach new floor when running, try fixing the problem and restart\n", MOTOR_ERROR_SEC)
	fmt.Printf("Program exiting now!\n")
	fmt.Printf("\n===============================================================================\n\n\n")
	os.Exit(1)
}