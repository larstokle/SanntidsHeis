package elevatorMgr

import (
	"eventmgr"
	"fmt"
	"fsm"
	. "globals"
	"message"
	"orderque"
	"transactionMgr"
	//"time"

)

func Start() {
	localElev := fsm.NewElevator()
	btnPush := eventmgr.CheckButtons()
	que := orderque.NewOrderQue()

	transMgr := transactionMgr.New()

	ifLowCostThenRequest := func(order Button_t){
		if order != NONVALID_BUTTON{
			cost := localElev.GetCost(order)
			if cost < fsm.INF_COST {
				transMgr.RequestOrder(order, cost)
			}
		}
	}
	if !que.IsEmpty() {
		newOrder := que.EarliestNonAssignedOrder() // switch with calculate from twoelevtest
		ifLowCostThenRequest(newOrder)
	}
	if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: init done entering loop\n\n")}
	
	go func() { //hmmm skal denne kjøre selv eller skal det go'es i main??
		for {
			select {
			case floorDone := <-localElev.OrderDone:
				if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: floorDone case from fsm = %+v\n", floorDone)}
				que.RemoveOrdersOnFloor(floorDone)
				transMgr.RemoveOrder(floorDone) //trøbbel! vi må vite at det ikke er en en cmd!! eller ikke ?? hmmm...
				if !que.IsEmpty() {
					newOrder := que.EarliestNonAssignedOrder() // switch with calculate from twoelevtest
					ifLowCostThenRequest(newOrder)
				}
				if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: floorDone case done = %+v\n", floorDone)}				
				continue // er ikke denne unødvendig? hvordan virker select? begynner den på ny etter en case? eller tar den neste case?

			case newBtn := <-btnPush:
				if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: newBtn case from eventmgr = %+v\n", newBtn)}
				if !que.HasOrder(newBtn) {
					que.AddOrder(newBtn)
					transMgr.NewOrder(newBtn)
					ifLowCostThenRequest(newBtn)
				}
				if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: newBtn case done! = %+v\n", newBtn)}


			case newMsg := <-transMgr.ToParent:
				if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: newMsg case from transMgr = %+v\n", newMsg)}
				switch newMsg.MessageId {

				case message.NEW_ORDER:
					newOrder := newMsg.Button
					if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: newOrder from transMgr = %+v\n", newOrder)}
					if !que.HasOrder(newOrder) {
						que.AddOrder(newOrder)
						ifLowCostThenRequest(newOrder)
					} else{
						fmt.Printf("elevMgr: ERROR new order aldready in que, me or ID %d must be out of sync\n", newMsg.Source)
					}
					if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: newMsg waiting for confirmation from transMgr \n")}			
					transMgr.ProceedOk <- true //GIVE CONFIRMATION TO TRANS! COULD USE TIMER THERE, BUT THIST IS PROBABLY BETTER
					if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: newMsg got confirmation from transMgr \n")}			

				case message.DELEGATE_ORDER:
					if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: got DELEGATE_ORDER from trans on order %+v to %d\n", newMsg.Button, newMsg.ElevatorId)}
					que.AssignOrderToId(newMsg.Button, newMsg.ElevatorId)
					if newMsg.ElevatorId == transMgr.MyId() {
						localElev.NewDestination(newMsg.Button.Floor)
						continue
					} else if !que.IsIdAssigned(transMgr.MyId()) && !que.IsEmpty(){
						newOrder := que.EarliestNonAssignedOrder() // switch with calculate from twoelevtest
						if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: did not get delegation and got no destination, get new\n\n")}
						ifLowCostThenRequest(newOrder)						
					}
				case message.REMOVE_ORDER:
					if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: got REMOVE_ORDER from trans on floor %d, removing both\n", newMsg.Button.Floor)}
					toRemove := newMsg.Button
					toRemove.ButtonType = UP
					que.RemoveOrder(toRemove)
					toRemove.ButtonType = DOWN
					que.RemoveOrder(toRemove)

				case message.COST:
					order := newMsg.Button
					cost := localElev.GetCost(order)
					if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: got COST from trans on order, i calculated cost = %d \n", order, cost)}
					if !que.HasOrder(order){
						if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: SENDTE NÅ COST PÅ EN ORDRE SOM IKKE ER I KØEN ENDA...\n\n") }//DEBUG ONLY
					}
					transMgr.SendCost(order, cost)
				case message.UNASSIGN_ORDER:
					unassignId := newMsg.ElevatorId
					fmt.Printf("elevMgr: got UNASSIGN_ORDER from trans on id %d\n", unassignId)
					if unassignId != transMgr.MyId(){
						que.UnassignOrdersToID(unassignId)
						newOrder := que.EarliestNonAssignedOrder() // switch with calculate from twoelevtest
						ifLowCostThenRequest(newOrder)

					} else{
						fmt.Printf("elevMgr: ERROR Id is mine, so i did not unassign my orders\n")
					}
				default:							
					if(DEBUG_ELEVMGR){fmt.Printf("Unhandeled MessageId: %+v", newMsg)}
				}
			}
		}
	}()
}


