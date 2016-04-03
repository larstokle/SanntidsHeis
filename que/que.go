package orderque

import (
	"time"
)

type OrderQue_t [N_FLOORS][N_buttons]struct {
	hasOrder       bool
	lastChangeTime time.Time
	//assignedToID int //kanskje unødvendig? fjerner den encapsulation?
}

type Order_t struct {
	floor     int
	orderType int
}

func (que OrderQue_t) AddOrder(newOrder Order_t) {
	que[newOrder.floor][newOrder.orderType].hasOrder = true
	que[newOrder.floor][newOrder.orderType].lastChangeTime = time.Now()
}

func (que OrderQue_t) CompleteOrder(doneOrder Order_t) {
	que[doneOrder.floor][doneOrder.orderType].hasOrder = false
	que[doneOrder.floor][doneOrder.orderType].lastChangeTime = time.Now()
}

func (thisQue OrderQue_t) Sync(queToSync OrderQue_t) OrderQue_t { //add error returns?
	for floor := 0; floor < N_FLOORS; floor++ {
		for orderType := 0; orderType < N_BUTTONS; orderType++ {
			if queToSync[floor][orderType].lastChangeTime > thisQue[floor][orderType].lastChangeTime {
				thisQue[floor][orderType] = queToSync[floor][orderType]
			}
		}
	}
	return thisQue
}

func (que OrderQue_t) HasOrderInDir(floor int, dir int) bool {
	switch dir {
	case DIR_DOWN:
		return que[floor][BTN_DOWN].hasOrder || que[floor][BTN_CMD].hasOrder
	case DIR_UP:
		return que[floor][BTN_UP].hasOrder || que[floor][BTN_CMD].hasOrder
	case DIR_STOP:
		return que[floor][BTN_CMD].hasOrder
	default:
		return false
	}
}

func (que OrderQue_t) IsEmpty() bool {
	for floor := 0; floor < N_FLOORS; floor++ {
		for orderType := 0; orderType < N_BUTTONS; orderType++ {
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
	for floor := 0; floor < N_FLOORS; floor++ {
		for orderType := 0; orderType < N_BUTTONS; orderType++ {
			if que[floor][orderType].lastChangeTime < que[earliestOrder.floor][earliestOrder.orderType].lastChangeTime {
				earliestOrder = Order_t{floor, orderType}
			}
		}
	}
	return earliestOrder
}
