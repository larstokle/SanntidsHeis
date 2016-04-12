package elevatorMgr

import (
	"orderque"
	"fsm"
	"eventmgr"
	"transactionMgr"
	."globals"
	"fmt"
	"message"
)



func Start(){
	que := orderque.NewOrderQue()
	localElev := fsm.NewElevator()
	btnPush := eventmgr.CheckButtons()

	transMgr := transactionMgr.New()
	go func(){//hmmm skal denne kjøre selv eller skal det go'es i main??
		for{
			select{
			case floorDone := <-localElev.OrderDone:
				fmt.Printf("elevMgr: floorDone from fsm = %+v\n", floorDone)
				que.RemoveOrdersOnFloor(floorDone)
				transMgr.RemoveOrder(floorDone)//TRØBBEL! VI MÅ VITE AT DET IKKE ER EN EN CMD!! eller ikke ?? hmmm...
				if !que.IsEmpty(){
					newOrder := que.EarliestNonAssignedOrder() // switch with calculate from twoelevtest			
					if newOrder.ButtonType == CMD{
						fmt.Printf("elevMgr: setting NewDestination: %+v \n", newOrder)
						que.AssignOrderToId(newOrder, transMgr.MyId())
						localElev.NewDestination(newOrder.Floor)
					}else{
						cost := localElev.GetCost(newOrder)
						transMgr.RequestOrder(newOrder, cost)
						fmt.Printf("elevMgr: setting NewDestination which should come from delegation: %+v \n", newOrder) //DELETE WHEN DELEGATION WORKS 
						que.AssignOrderToId(newOrder, transMgr.MyId())//DELETE WHEN DELEGATION WORKS 
						localElev.NewDestination(newOrder.Floor)	//DELETE WHEN DELEGATION WORKS 
					}
				}

			case newBtn := <- btnPush:
				fmt.Printf("elevMgr: newBtn from eventmgr = %+v\n",newBtn)
				if !que.HasOrder(newBtn){
					que.AddOrder(newBtn)
					if newBtn.ButtonType == CMD{
						fmt.Printf("elevMgr: setting NewDestination: %+v \n", newBtn)
						que.AssignOrderToId(newBtn, transMgr.MyId())
						localElev.NewDestination(newBtn.Floor)
					}else{
						transMgr.NewOrder(newBtn)
						cost := localElev.GetCost(newBtn)
						if cost < fsm.INF_COST{
							transMgr.RequestOrder(newBtn, cost)
							fmt.Printf("elevMgr: setting NewDestination which should come from delegation: %+v \n", newBtn) //DELETE WHEN DELEGATION WORKS 
							que.AssignOrderToId(newBtn, transMgr.MyId())//DELETE WHEN DELEGATION WORKS 
							localElev.NewDestination(newBtn.Floor)	//DELETE WHEN DELEGATION WORKS 
						}
					}
				}

			case newMsg := <- transMgr.Receive:
				fmt.Printf("elevMgr: newMsg from transMgr = %+v\n",newMsg)
				switch newMsg.MessageId{
				case message.NEW_ORDER:
					newOrder := newMsg.Button
					fmt.Printf("elevMgr: newOrder from transMgr = %+v\n",newOrder)
					if !que.HasOrder(newOrder){
						que.AddOrder(newOrder)
						if newOrder.ButtonType == CMD{
							fmt.Printf("elevMgr: setting NewDestination: %+v \n", newOrder)
							que.AssignOrderToId(newOrder, transMgr.MyId())
							localElev.NewDestination(newOrder.Floor)
						}else{
							cost := localElev.GetCost(newOrder)
							transMgr.RequestOrder(newOrder, cost)
							fmt.Printf("elevMgr: setting NewDestination which should come from delegation: %+v \n", newOrder) //DELETE WHEN DELEGATION WORKS 
							que.AssignOrderToId(newOrder, transMgr.MyId())//DELETE WHEN DELEGATION WORKS 
							localElev.NewDestination(newOrder.Floor)	//DELETE WHEN DELEGATION WORKS 
						}
					}

				case message.DELEGATE_ORDER:
					que.AssignOrderToId(newMsg.Button, newMsg.ElevatorId)
					if newMsg.ElevatorId == transMgr.MyId(){
						localElev.NewDestination(newMsg.Button.Floor)
					} else if !que.IsIdAssigned(transMgr.MyId()){
						newOrder := que.EarliestNonAssignedOrder() // switch with calculate from twoelevtest
						cost := localElev.GetCost(newOrder)
						transMgr.RequestOrder(newOrder, cost)
					}
				case message.REMOVE_ORDER:
					toRemove := newMsg.Button
					toRemove.ButtonType = UP
					que.RemoveOrder(toRemove)
					toRemove.ButtonType = DOWN
					que.RemoveOrder(toRemove)

				case message.REQUEST_ORDER:
					order := newMsg.Button
					cost := localElev.GetCost(order)
					transMgr.Cost(order, cost)
				default:
					fmt.Printf("Unhandeled MessageId: %+v",newMsg)
				}				
			}
		}

	}()


}