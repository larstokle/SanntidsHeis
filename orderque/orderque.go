package orderque

import (
	. "../constants"
	"../driver"
	"fmt"
	"time"
)

const UNASIGNED_ID = 0

type OrderQue_t [N_FLOORS][N_ORDER_TYPES]struct {
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

func (que *OrderQue_t) AddOrder(floor int, orderType int) {
	if !que.HasOrder(floor, orderType) {
		que[floor][orderType].hasOrder = true
		que[floor][orderType].lastChangeTime = time.Now()
		que[floor][orderType].assignedToID = UNASIGNED_ID
		driver.SetButtonLight(orderType, floor, true)
	}
}

func (que *OrderQue_t) RemoveOrder(floor int, orderType int) {
	if que.HasOrder(floor, orderType) {
		que[floor][orderType].hasOrder = false
		que[floor][orderType].lastChangeTime = time.Now()
		que[floor][orderType].assignedToID = UNASIGNED_ID
		driver.SetButtonLight(orderType, floor, false)
	}
}

func (que *OrderQue_t) UnassignOrderToID(id int){
	for floor := FIRST_FLOOR; floor < N_FLOORS; floor++ {
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			if que[floor][orderType].assignedToID == id{
				que[floor][orderType].assignedToID = UNASIGNED_ID
			}
		}
	}
}

func (que *OrderQue_t) AssignOrderToID(floor int, orderType int, id int) bool{
	if !que.HasOrder(floor, orderType){
		return false
	}

	que.UnassignOrderToID(id)

	que[floor][orderType].assignedToID = id
	return true
}

func (thisQue *OrderQue_t) Sync(queToSync OrderQue_t) OrderQue_t { //add error returns?
	for floor := FIRST_FLOOR; floor < N_FLOORS; floor++ {
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			if queToSync[floor][orderType].lastChangeTime.After(thisQue[floor][orderType].lastChangeTime) {
				thisQue[floor][orderType] = queToSync[floor][orderType]
			}
		}
	}
	return *thisQue
}

func (que *OrderQue_t) HasOrder(floor int, orderType int) bool {
	return que[floor][orderType].hasOrder
}

func (que *OrderQue_t) IsEmpty() bool {
	for floor := FIRST_FLOOR; floor <= TOP_FLOOR; floor++ {
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			if que[floor][orderType].hasOrder {
				return false
			}
		}
	}
	return true
}

func (que *OrderQue_t) EarliestNonAssignedOrder() Order_t {
	if que.IsEmpty() {
		return Order_t{-1, -1}
	}

	var nonInitializedTime time.Time
	earliestOrder := Order_t{FIRST_FLOOR, DOWN} //Button doesn't exist, used as a dummy initializer only

	for floor := FIRST_FLOOR; floor <= TOP_FLOOR; floor++ {
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			if que[floor][orderType].hasOrder  && (que[floor][orderType].assignedToID == UNASIGNED_ID) {
				currentTime := que[floor][orderType].lastChangeTime
				earliestTime := que[earliestOrder.floor][earliestOrder.orderType].lastChangeTime
				if (currentTime.Before(earliestTime) || earliestTime == nonInitializedTime){
					earliestOrder = Order_t{floor: floor, orderType: orderType}
				}
			}
		}
	}
	return earliestOrder
}

func (que *OrderQue_t) NextOrderOfTypeInDir(currentFloor int, dir int, orderType int) int {
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

func (que *OrderQue_t) Print() {
	fmt.Println("OrderQue_t:")
	for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
		fmt.Printf("\tordertype: %d\n", orderType)
		for floor := FIRST_FLOOR; floor <= TOP_FLOOR; floor++ {
			fmt.Printf("\t\tFloor %d: has order = %t,last Changed: = %v\n", floor, que[floor][orderType].hasOrder, que[floor][orderType].lastChangeTime)
		}
	}
	fmt.Println()
}
