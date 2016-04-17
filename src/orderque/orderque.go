package orderque

import (
	. "globals"
	"driver"
	"fmt"
	"time"
)

const UNASIGNED_ID = NONLEGAL_ID
const QUE_LOG_FILE = "QueLog.txt"

type order_t struct {
	hasOrder       bool
	lastChangeTime time.Time
	assignedToID int 
}

type orderQue_t [N_FLOORS][N_ORDER_TYPES] order_t


func New()orderQue_t{
	var que orderQue_t
	for floor := FIRST_FLOOR; floor < N_FLOORS; floor++ {
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			que[floor][orderType].assignedToID = UNASIGNED_ID
		}
	}
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
	que.WriteToLog()
	if(DEBUG_QUE){que.Print()}
}

func (que *orderQue_t) RemoveOrder(button Button_t) {
	if(DEBUG_QUE){fmt.Printf("que: RemoveOrder = %+v\n",button)}
	if que.HasOrder(button) {
		que[button.Floor][button.ButtonType].hasOrder = false
		que[button.Floor][button.ButtonType].lastChangeTime = time.Now()
		que[button.Floor][button.ButtonType].assignedToID = UNASIGNED_ID
		driver.SetButtonLight(button.ButtonType, button.Floor, false)
	}
	que.WriteToLog()
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
	if id == UNASIGNED_ID{
		return
	}
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

func (que *orderQue_t) UnassignAllOrders(){
	for floor := FIRST_FLOOR; floor < N_FLOORS; floor++ {
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			que[floor][orderType].assignedToID = UNASIGNED_ID
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

func (que *orderQue_t) IsOrderAssigned(order Button_t) bool{
	return que[order.Floor][order.ButtonType].assignedToID != UNASIGNED_ID
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

func (que *orderQue_t) PrintVerbose() {
	fmt.Println("orderQue_t:{")
	for floor := FIRST_FLOOR; floor <= TOP_FLOOR; floor++ {
		fmt.Printf("\tFloor: %d\n", floor)
		for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
			fmt.Printf("\t  type: %d: has order = %t, assigned to id: %d\n \t\tlast Changed: = %v\n", orderType, que[floor][orderType].hasOrder,que[floor][orderType].assignedToID, que[floor][orderType].lastChangeTime)
		}
	}
	fmt.Println("}")
}


func (que *orderQue_t) Print(){
	qPrint := fmt.Sprintln("orderQue_t:{")
	for floor := FIRST_FLOOR; floor <= TOP_FLOOR; floor++ {
		qPrint = qPrint + fmt.Sprintf("\tFloor %d: ", floor)
		for orderType, typeWord :=  range []string{"UP", "DOWN", "CMD"}{
			has := 0
			if que[floor][orderType].hasOrder{
				has = 1
			}
			qPrint = qPrint + fmt.Sprintf("%s{ %d, Id: %d}, ", typeWord, has, que[floor][orderType].assignedToID)
		}
		qPrint = qPrint + fmt.Sprintln()
	}
	qPrint = qPrint + fmt.Sprintln("}")
	fmt.Print(qPrint)
}
	


