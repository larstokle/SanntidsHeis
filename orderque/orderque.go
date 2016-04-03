package orderque

import (
	"../driver"
	"fmt"
	"time"
)

type OrderQue_t [driver.N_FLOORS][driver.N_BUTTONS]struct {
	hasOrder       bool
	lastChangeTime time.Time
	//assignedToID int //kanskje unødvendig? fjerner den encapsulation?
}

type Order_t struct {
	floor     int
	orderType int
}

func (que OrderQue_t) AddOrder(floor int, orderType int) {
	que[floor][orderType].hasOrder = true
	que[floor][orderType].lastChangeTime = time.Now()
}

func (que OrderQue_t) CompleteOrder(floor int, orderType int) {
	que[floor][orderType].hasOrder = false
	que[floor][orderType].lastChangeTime = time.Now()
}

func (thisQue OrderQue_t) Sync(queToSync OrderQue_t) OrderQue_t { //add error returns?
	for floor := 0; floor < driver.N_FLOORS; floor++ {
		for orderType := 0; orderType < driver.N_BUTTONS; orderType++ {
			if queToSync[floor][orderType].lastChangeTime.After(thisQue[floor][orderType].lastChangeTime) {
				thisQue[floor][orderType] = queToSync[floor][orderType]
			}
		}
	}
	return thisQue
}

func (que OrderQue_t) HasOrderInDir(floor int, dir int) bool {
	switch dir {
	case driver.DIR_DOWN:
		return que[floor][driver.BTN_DOWN].hasOrder || que[floor][driver.BTN_CMD].hasOrder
	case driver.DIR_UP:
		return que[floor][driver.BTN_UP].hasOrder || que[floor][driver.BTN_CMD].hasOrder
	case driver.DIR_STOP:
		return que[floor][driver.BTN_CMD].hasOrder
	default:
		return false
	}
}

func (que OrderQue_t) IsEmpty() bool {
	for floor := 0; floor < driver.N_FLOORS; floor++ {
		for orderType := 0; orderType < driver.N_BUTTONS; orderType++ {
			if que[floor][orderType].hasOrder {
				return true
			}
		}
	}
	return false
}

func (que OrderQue_t) EarliestOrderInside() Order_t {
	if que.IsEmpty() {
		return Order_t{-1, -1}
	}

	earliestOrder := Order_t{0, 0}
	for floor := 0; floor < driver.N_FLOORS; floor++ {
		for orderType := 0; orderType < driver.N_BUTTONS; orderType++ {
			if que[floor][orderType].lastChangeTime.Before(que[earliestOrder.floor][earliestOrder.orderType].lastChangeTime) {
				earliestOrder = Order_t{floor, orderType}
			}
		}
	}
	return earliestOrder
}

func (que OrderQue_t) Print() {
	fmt.Println("OrderQue_t:")
	for orderType := 0; orderType < driver.N_BUTTONS; orderType++ {
		fmt.Printf("\tordertype: %d\n", orderType)
		for floor := 0; floor < driver.N_FLOORS; floor++ {
			fmt.Printf("\t\tFloor %d: has order = %t,last Changed: = %v\n\n", floor, que[floor][orderType].hasOrder, que[floor][orderType].lastChangeTime)
		}
	}
}
