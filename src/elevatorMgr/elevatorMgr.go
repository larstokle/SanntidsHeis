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
	que := orderque.NewOrderQue()
	localElev := fsm.NewElevator()
	btnPush := eventmgr.CheckButtons()

	transMgr := transactionMgr.New()

	if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: init done entering loop\n\n")}
	go func() { //hmmm skal denne kjøre selv eller skal det go'es i main??
		for {
			select {
			case floorDone := <-localElev.OrderDone:
				if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: floorDone from fsm = %+v\n", floorDone)}
				que.RemoveOrdersOnFloor(floorDone)
				transMgr.RemoveOrder(floorDone) //trøbbel! vi må vite at det ikke er en en cmd!! eller ikke ?? hmmm...
				if !que.IsEmpty() {
					newOrder := que.EarliestNonAssignedOrder() // switch with calculate from twoelevtest
					if newOrder.ButtonType == CMD {
						if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: setting NewDestination: %+v \n", newOrder)}
						que.AssignOrderToId(newOrder, transMgr.MyId())
						localElev.NewDestination(newOrder.Floor)
					} else if newOrder != NONVALID_BUTTON {
						cost := localElev.GetCost(newOrder)
						canTakeAtOnce := transMgr.RequestOrder(newOrder, cost)
						if canTakeAtOnce{ //DENNE PRØVER Å LØSE PROBLEMET LENGRE NED MED REQUEST FØR COST OG MOTSATT, VED Å TILLATE UBUFFRET FRA TRANS: LITT ADD HOOK
							que.AssignOrderToId(newOrder, transMgr.MyId())
							localElev.NewDestination(newOrder.Floor)
						}
					} else {
						if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: que returned NONVALID_BUTTON while !que.IsEmpty() => all orders taken\n")}
					}
				}
				if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: floorDone case done = %+v\n", floorDone)}				
				continue // er ikke denne unødvendig?

			case newBtn := <-btnPush:
				if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: newBtn from eventmgr = %+v\n", newBtn)}
				if !que.HasOrder(newBtn) {
					que.AddOrder(newBtn)
					cost := localElev.GetCost(newBtn)
					if newBtn.ButtonType == CMD && cost < fsm.INF_COST {
						if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: setting NewDestination: %+v \n", newBtn)}
						que.AssignOrderToId(newBtn, transMgr.MyId())
						localElev.NewDestination(newBtn.Floor)
					} else if newBtn.ButtonType != CMD {
						transMgr.NewOrder(newBtn)
						if cost < fsm.INF_COST {
							canTakeAtOnce := transMgr.RequestOrder(newBtn, cost)
							if canTakeAtOnce{ //DENNE PRØVER Å LØSE PROBLEMET LENGRE NED MED REQUEST FØR COST OG MOTSATT, VED Å TILLATE UBUFFRET FRA TRANS: LITT ADD HOOK
								if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: confirmation to take order at once from transMgr \n")}			
								que.AssignOrderToId(newBtn, transMgr.MyId())
								localElev.NewDestination(newBtn.Floor)								
							}
						}
					}	
				}
				if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: newBtn case done! = %+v\n", newBtn)}


			case newMsg := <-transMgr.Receive:
				if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: newMsg from transMgr = %+v\n", newMsg)}
				switch newMsg.MessageId {

				case message.NEW_ORDER:
					newOrder := newMsg.Button
					if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: newOrder from transMgr = %+v\n", newOrder)}
					if !que.HasOrder(newOrder) {
						que.AddOrder(newOrder)
						cost := localElev.GetCost(newOrder)
						if cost < fsm.INF_COST {
							canTakeAtOnce := transMgr.RequestOrder(newOrder, cost)
							if canTakeAtOnce{ //DENNE PRØVER Å LØSE PROBLEMET LENGRE NED MED REQUEST FØR COST OG MOTSATT, VED Å TILLATE UBUFFRET FRA TRANS: LITT ADD HOOK
								que.AssignOrderToId(newOrder, transMgr.MyId())
								localElev.NewDestination(newOrder.Floor)
							}
						}
					}
					if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: newMsg waiting for confirmation from transMgr \n")}			
					transMgr.Receive <- message.Message_t{} //GIVE CONFIRMATION TO TRANS! COULD USE TIMER THER, BUT THIST IS PROBABLY BETTER
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
						if newOrder != NONVALID_BUTTON {
							cost := localElev.GetCost(newOrder)
							if cost < fsm.INF_COST {
								canTakeAtOnce := transMgr.RequestOrder(newOrder, cost)
								if canTakeAtOnce{ //DENNE PRØVER Å LØSE PROBLEMET LENGRE NED MED REQUEST FØR COST OG MOTSATT, VED Å TILLATE UBUFFRET FRA TRANS: LITT ADD HOOK
									que.AssignOrderToId(newOrder, transMgr.MyId())
									localElev.NewDestination(newOrder.Floor)
								}
							} else {
								if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: inf cost when received DELEGATE_ORDER and i got no order\n")}
							}
						} else {
							if(DEBUG_ELEVMGR){fmt.Printf("elevMgr: que returned NONVALID_BUTTON while !que.IsEmpty() => all orders taken\n")}
						}
						
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
					transMgr.SendCost(order, cost) /*MULIG DETTE ER ET PROBLEM!!!! DEN ANDRE HEISEN REKKER Å PULLE COSTEN FØR DENNE HEISEN HAR LAGT ORDREN I KØEN.
													SÅ DENNE SENDER EN REQUEST SOM GJØR ORDREN FERDIG DELEGERT FØR DEN SENDER DENNE*/
				default:							
					if(DEBUG_ELEVMGR){fmt.Printf("Unhandeled MessageId: %+v", newMsg)}
				}

			/*default:
				if que.IsEmpty(){
					//fmt.Printf("Que is empty... \n")
					time.Sleep(time.Millisecond*100)
				} else if localElev.State()== fsm.STATE_IDLE {
					nextOrder := que.EarliestNonAssignedOrder()
					if nextOrder != NONVALID_BUTTON{

						//DENNNE MÅ SJEKKE AT VI ER ALENE!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
						localElev.NewDestination(nextOrder.Floor)
						fmt.Printf("elevMgr default  made me take%+v \n", nextOrder)
					} else{
						fmt.Printf("All orders assigned... \n")
					}
					time.Sleep(time.Millisecond*100)

				}
			*/	
		}
	}

	}()

}


