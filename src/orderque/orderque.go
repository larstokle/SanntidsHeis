package orderque

import (
	. "globals"
	"driver"
	"fmt"
	"time"
	"orderLog"
)

const UNASIGNED_ID = NONLEGAL_ID

type orderQue_t [N_FLOORS][N_ORDER_TYPES]struct {
	hasOrder       bool
	lastChangeTime time.Time
	assignedToID int //kanskje unødvendig? fjerner den encapsulation?
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
	que.syncWithInternalLog()
	if(DEBUG_QUE){fmt.Println("que: init done NewOrderQue returned\n")}
	return que
}

func (que *orderQue_t) AddOrder(button Button_t) {
	if(DEBUG_QUE){fmt.Printf("que: AddOrder = %+v\n",button)}
	if !que.HasOrder(button) {
		que[button.Floor][button.ButtonType].hasOrder = true
		que[button.Floor][button.ButtonType].lastChangeTime = time.Now()
		que[button.Floor][button.ButtonType].assignedToID = UNASIGNED_ID
		driver.SetButtonLight(button.ButtonType, button.Floor, true)
	}
	que.logInternal()

}

func (que *orderQue_t) RemoveOrder(button Button_t) {
	if(DEBUG_QUE){fmt.Printf("que: RemoveOrder = %+v\n",button)}
	if que.HasOrder(button) {
		que[button.Floor][button.ButtonType].hasOrder = false
		que[button.Floor][button.ButtonType].lastChangeTime = time.Now()
		que[button.Floor][button.ButtonType].assignedToID = UNASIGNED_ID
		driver.SetButtonLight(button.ButtonType, button.Floor, false)
	}
	que.logInternal()

}

func (que *orderQue_t) RemoveOrdersOnFloor(floor int){
	if(DEBUG_QUE){fmt.Printf("que: RemoveOrdersOnFloor = %+v\n",floor)}
	for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
		order := Button_t{Floor: floor, ButtonType: orderType}
		que.RemoveOrder(order)
	}
	if(DEBUG_QUE){fmt.Printf("Remaining orders: \n")}
	if(DEBUG_QUE){que.Print()}
}

func (que *orderQue_t) UnassignOrdersToID(id int){
	if(DEBUG_QUE){fmt.Printf("que: UnassignOrdersToId = %d\n",id)}
	for floor := FIRST_FLOOR; floor < N_FLOORS; floor++ {
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			if que[floor][orderType].assignedToID == id{
				que[floor][orderType].assignedToID = UNASIGNED_ID
				if(DEBUG_QUE){fmt.Printf("que: Unassigned order: floor = %d, orderType = %d", floor, orderType)}
			}
		}
	}
}

func (que *orderQue_t) AssignOrderToId(button Button_t, id int) bool{
	if !que.HasOrder(button){
		return false
	}

	que.UnassignOrdersToID(id)
	if(DEBUG_QUE){fmt.Printf("que: AssignOrder = %+v, ToId = %d\n",button,id)}
	que[button.Floor][button.ButtonType].assignedToID = id

	return true
}

func (que *orderQue_t) IsIdAssigned(id int) bool{
	for floor := FIRST_FLOOR; floor < N_FLOORS; floor++ {
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			if que[floor][orderType].assignedToID == id{
				return true
			}
		}
	}
	return false
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
	hasOrder := que[button.Floor][button.ButtonType].hasOrder
	//if(DEBUG_QUE){fmt.Printf("que: HasOrder %+v? %t\n", button, hasOrder)}
	return hasOrder
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

func (que *orderQue_t) EarliestNonAssignedOrder() Button_t {
	if que.IsEmpty() {
		return NONVALID_BUTTON
	}

	
	earliestOrder := NONVALID_BUTTON

	for floor := FIRST_FLOOR; floor <= TOP_FLOOR; floor++ {
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			if que[floor][orderType].hasOrder  && (que[floor][orderType].assignedToID == UNASIGNED_ID) {
				currentTime := que[floor][orderType].lastChangeTime
				earliestTime := que[earliestOrder.Floor][earliestOrder.ButtonType].lastChangeTime
				if currentTime.Before(earliestTime) || earliestOrder == NONVALID_BUTTON{
					earliestOrder = Button_t{Floor: floor, ButtonType: orderType}
				}
			}
		}
	}
	return earliestOrder
}

func (que *orderQue_t) NearestOrderOfTypeInDir(currentFloor int, dir int, orderType int) int {
	switch dir{
	case DIR_UP, DIR_STOP:
		for checkFloor := currentFloor; checkFloor <= TOP_FLOOR; checkFloor++ {
			if que[checkFloor][orderType].hasOrder {
				return checkFloor
			}
		}
	case DIR_DOWN:
		for checkFloor := currentFloor; checkFloor >= FIRST_FLOOR; checkFloor-- {
			if que[checkFloor][orderType].hasOrder {
				return checkFloor
			}
		}
	}
	
	return -1
}

func (que *orderQue_t) Print() {
	if(DEBUG_QUE){fmt.Println("orderQue_t:")}
	for floor := FIRST_FLOOR; floor <= TOP_FLOOR; floor++ {
		if(DEBUG_QUE){fmt.Printf("\tFloor: %d\n", floor)}
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			if(DEBUG_QUE){fmt.Printf("\t\t type: %d: has order = %t,last Changed: = %v\n", orderType, que[floor][orderType].hasOrder, que[floor][orderType].lastChangeTime)}
		}
	}
	if(DEBUG_QUE){fmt.Println()}
}


func (que *orderQue_t) logInternal(){
	internalOrders := make([]byte, N_FLOORS)
	
	for floor := FIRST_FLOOR; floor <= TOP_FLOOR; floor++ {
		if(que[floor][CMD].hasOrder){
			internalOrders[floor] = '1'
		} else{
			internalOrders[floor] = '0'
		}

	}
	orderLog.WriteInternalOrders(internalOrders)
}


func (que *orderQue_t) syncWithInternalLog(){
	internalOrders := orderLog.FindInternalOrders()
	for floor := FIRST_FLOOR; floor <= TOP_FLOOR; floor++ {
		if(internalOrders[floor] == '1'){
			order := Button_t{Floor: floor, ButtonType: CMD}
			que.AddOrder(order)
		} else {
			que[floor][CMD].hasOrder = false
		}
	}
}
