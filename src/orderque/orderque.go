package orderque

import (
	. "globals"
	"driver"
	"fmt"
	"time"
)

const UNASIGNED_ID = NONLEGAL_ID

type orderQue_t [N_FLOORS][N_ORDER_TYPES]struct {
	hasOrder       bool
	lastChangeTime time.Time
	assignedToID int //kanskje unødvendig? fjerner den encapsulation?
}

type Order_t struct {
	floor     int
	orderType int
}

func (order Order_t) GetFloor() int {
	return order.floor
}

func (order Order_t) Getype() int {
	return order.orderType
}

func NewOrderQue()orderQue_t{
	var que orderQue_t
	startTime := time.Now()
	for floor := FIRST_FLOOR; floor < N_FLOORS; floor++ {
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			que[floor][orderType].assignedToID = UNASIGNED_ID
			que[floor][orderType].lastChangeTime = startTime
		}
	}
	return que
}

func (que *orderQue_t) AddOrder(button Button_t) {
	if !que.HasOrder(button) {
		que[button.Floor][button.ButtonType].hasOrder = true
		que[button.Floor][button.ButtonType].lastChangeTime = time.Now()
		que[button.Floor][button.ButtonType].assignedToID = UNASIGNED_ID
		driver.SetButtonLight(button.ButtonType, button.Floor, true)
	}
}

func (que *orderQue_t) RemoveOrder(button Button_t) {
	if que.HasOrder(button) {
		que[button.Floor][button.ButtonType].hasOrder = false
		que[button.Floor][button.ButtonType].lastChangeTime = time.Now()
		que[button.Floor][button.ButtonType].assignedToID = UNASIGNED_ID
		driver.SetButtonLight(button.ButtonType, button.Floor, false)
	}
}

func (que *orderQue_t) UnassignOrderToID(id int){
	for floor := FIRST_FLOOR; floor < N_FLOORS; floor++ {
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			if que[floor][orderType].assignedToID == id{
				que[floor][orderType].assignedToID = UNASIGNED_ID
			}
		}
	}
}

func (que *orderQue_t) AssignOrderToID(button Button_t, id int) bool{
	if !que.HasOrder(button){
		return false
	}

	que.UnassignOrderToID(id)

	que[button.Floor][button.ButtonType].assignedToID = id
	return true
}

func (thisQue *orderQue_t) Sync(queToSync orderQue_t) orderQue_t { //add error returns?
	for floor := FIRST_FLOOR; floor < N_FLOORS; floor++ {
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			if queToSync[floor][orderType].lastChangeTime.After(thisQue[floor][orderType].lastChangeTime) {
				thisQue[floor][orderType] = queToSync[floor][orderType]
			}
		}
	}
	return *thisQue
}

func (que *orderQue_t) HasOrder(button Button_t) bool {
	return que[button.Floor][button.ButtonType].hasOrder
}

func (que *orderQue_t) IsEmpty() bool {
	for floor := FIRST_FLOOR; floor <= TOP_FLOOR; floor++ {
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			if que[floor][orderType].hasOrder {
				return false
			}
		}
	}
	return true
}

func (que *orderQue_t) EarliestNonAssignedOrder() Order_t {
	if que.IsEmpty() {
		return Order_t{-1, -1}
	}

	nonValidOrder := Order_t{FIRST_FLOOR, DOWN}
	earliestOrder := nonValidOrder

	for floor := FIRST_FLOOR; floor <= TOP_FLOOR; floor++ {
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			if que[floor][orderType].hasOrder  && (que[floor][orderType].assignedToID == UNASIGNED_ID) {
				currentTime := que[floor][orderType].lastChangeTime
				earliestTime := que[earliestOrder.floor][earliestOrder.orderType].lastChangeTime
				if currentTime.Before(earliestTime) || earliestOrder == nonValidOrder{
					earliestOrder = Order_t{floor: floor, orderType: orderType}
				}
			}
		}
	}
	return earliestOrder
}

func (que *orderQue_t) NextOrderOfTypeInDir(currentFloor int, dir int, orderType int) int {
	if dir == DIR_STOP {
		dir = DIR_UP //kanskje noe annet her?
	}
	for checkFloor := currentFloor; checkFloor >= FIRST_FLOOR && checkFloor <= TOP_FLOOR; checkFloor += dir {
		if que[checkFloor][orderType].hasOrder {
			return checkFloor
		}
	}
	return -1
}

func (que *orderQue_t) Print() {
	fmt.Println("orderQue_t:")
	for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
		fmt.Printf("\tordertype: %d\n", orderType)
		for floor := FIRST_FLOOR; floor <= TOP_FLOOR; floor++ {
			fmt.Printf("\t\tFloor %d: has order = %t,last Changed: = %v\n", floor, que[floor][orderType].hasOrder, que[floor][orderType].lastChangeTime)
		}
	}
	fmt.Println()
}
