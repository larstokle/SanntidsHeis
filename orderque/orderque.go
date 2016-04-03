package orderque

import (
	. "../constants"
	"../driver"
	"fmt"
	"time"
)

type OrderQue_t [N_FLOORS][N_ORDER_TYPES]struct {
	hasOrder       bool
	lastChangeTime time.Time
	//assignedToID int //kanskje unødvendig? fjerner den encapsulation?
}

type Order_t struct {
	floor     int
	orderType int
}

func (order Order_t)GetFloor() int {
	return order.floor
}

func (order Order_t)Getype() int {
	return order.orderType
}

func (que *OrderQue_t) AddOrder(floor int, orderType int) {
	if !que.HasOrder(floor, orderType) {
		que[floor][orderType].hasOrder = true
		que[floor][orderType].lastChangeTime = time.Now()
		driver.SetButtonLight(orderType, floor, true)
	}
}

func (que *OrderQue_t) RemoveOrder(floor int, orderType int) {
	if que.HasOrder(floor, orderType) {
		que[floor][orderType].hasOrder = false
		que[floor][orderType].lastChangeTime = time.Now()
		driver.SetButtonLight(orderType, floor, false)
	}
}

func (thisQue *OrderQue_t) Sync(queToSync OrderQue_t) OrderQue_t { //add error returns?
	for floor := 0; floor < N_FLOORS; floor++ {
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
	for floor := 0; floor < N_FLOORS; floor++ {
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			if que[floor][orderType].hasOrder {
				return false
			}
		}
	}
	return true
}

func (que *OrderQue_t) EarliestOrderInside() Order_t {
	if que.IsEmpty() {
		return Order_t{-1, -1}
	}

	earliestOrder := Order_t{0, 0}
	for floor := 0; floor < N_FLOORS; floor++ {
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			if que[floor][orderType].lastChangeTime.Before(que[earliestOrder.floor][earliestOrder.orderType].lastChangeTime) {
				earliestOrder = Order_t{floor, orderType}
			}
		}
	}
	return earliestOrder
}

func (que *OrderQue_t) Print() {
	fmt.Println("OrderQue_t:")
	for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
		fmt.Printf("\tordertype: %d\n", orderType)
		for floor := 0; floor < N_FLOORS; floor++ {
			fmt.Printf("\t\tFloor %d: has order = %t,last Changed: = %v\n", floor, que[floor][orderType].hasOrder, que[floor][orderType].lastChangeTime)
		}
	}
	fmt.Println()
}
