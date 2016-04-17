package transactionMgr

import (
 	"fmt"
 	"message"
 	."globals"
)

func (transMgr *transactionMgr_t) ParentReady(){
	select{
	case transMgr.ProceedOk <- true:
	default:
	}
}

func (transMgr *transactionMgr_t) MyId() int {
	return transMgr.myId
}

func (transMgr *transactionMgr_t) NewOrder(order Button_t) {
	if order.ButtonType == CMD{
		if(DEBUG_TRNSMGR){fmt.Printf("transMgr: did not send new order (%+v) since ButtonType == CMD\n", order)}
		return
	}

	numElevs := transMgr.nElevatorsOnline()
	if numElevs > 1{
		if(DEBUG_TRNSMGR){fmt.Printf("transMgr: sending new order on network = %+v to in total %d elevs\n", order, numElevs)}
		
		
		transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.NEW_ORDER, Button: order}
		

	} else{
		if(DEBUG_TRNSMGR){fmt.Printf("transMgr: did not send new order (%+v) since numElevs = %d \n",order, numElevs)}
	}
}

func (transMgr *transactionMgr_t) RequestOrder(order Button_t, cost int) /*bool*/{
	independentSender := func(){// =========================================================================<<<<<<<SUPERHACK!!
			transMgr.ToParent <- message.Message_t{MessageId: message.DELEGATE_ORDER, Button: order, ElevatorId: transMgr.myId}
	}

	if transMgr.nElevatorsOnline() <= 1{
		if(DEBUG_TRNSMGR){fmt.Printf("transMgr: Request ANY Order on %+v with cost %d and %d elevs online, TAKE IT\n", order, cost, transMgr.nElevatorsOnline())}

		go independentSender()
	} else if order.ButtonType == CMD {
		if(DEBUG_TRNSMGR){fmt.Printf("transMgr: Request CMD Order on %+v with cost %d and %d elevs online, TAKE IT\n", order, cost, transMgr.nElevatorsOnline())}
		go independentSender()

		transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.UNASSIGN_ORDER, ElevatorId: transMgr.myId}

	} else {
		if(DEBUG_TRNSMGR){fmt.Printf("transMgr: Request UP/DOWN Order on %+v with cost %d and %d elevs online\n", order, cost, transMgr.nElevatorsOnline())}
		transMgr.SendCost(order, cost)

	}

}

func (transMgr *transactionMgr_t) RemoveOrder(floor int) {
	if(DEBUG_TRNSMGR){fmt.Printf("transMgr: sending remove order (floor) on network = %+v\n", floor)}
	order := Button_t{Floor: floor, ButtonType: UP}
	if transMgr.nElevatorsOnline() > 1{
		
		transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.REMOVE_ORDER, Button: order}
		
	}
	//kan dette ligge her eller vil det kunne føre til en panic?? må det flyttes inn i ifen?
	order.ButtonType = UP
	transMgr.delegationMutex.Lock()
	delete(transMgr.delegation, order)
	order.ButtonType = DOWN
	delete(transMgr.delegation, order)
	transMgr.delegationMutex.Unlock()
}

func (transMgr *transactionMgr_t) SendCost(order Button_t, cost int) {
	if(DEBUG_TRNSMGR){fmt.Printf("transMgr: sending cost %d on order %+v on network\n", cost, order)}

	transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.COST, Button: order, Cost: cost}
	transMgr.costToDelegation(order, cost, transMgr.myId)
}

func (transMgr *transactionMgr_t) SendSync(rawQue []byte){
	if(DEBUG_TRNSMGR){fmt.Printf("transMgr: SendSync got que from elevMgr, sending que on network\n")}
	transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.SYNC, Data: rawQue}
}

