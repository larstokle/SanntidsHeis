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
				newOrder := que.EarliestNonAssignedOrder() // switch with calculate from twoelevtest
				cost := localElev.GetCost(newOrder)
				transMgr.RequestOrder(newOrder, cost)
				if newOrder != NONVALID_BUTTON{//DELETE------------------------
					fmt.Printf("setting NewDestination: %+v \n", newOrder)
					que.AssignOrderToId(newOrder, transMgr.MyId())
					localElev.NewDestination(newOrder.Floor)
				}

			case newBtn := <- btnPush:
				fmt.Printf("elevMgr: newBtn from eventmgr = %+v\n",newBtn)
				if !que.HasOrder(newBtn){
					que.AddOrder(newBtn)
					if newBtn.ButtonType != CMD{
						transMgr.NewOrder(newBtn)
					}
					cost := localElev.GetCost(newBtn)
					if cost < fsm.INF_COST{
						transMgr.RequestOrder(newBtn, cost)
						que.AssignOrderToId(newBtn, transMgr.MyId())//DELETE------------------------
						localElev.NewDestination(newBtn.Floor)
					}
				}

			case newMsg := <- transMgr.Receive:
				fmt.Printf("elevMgr: newMsg from transMgr = %+v\n",newMsg)
				switch newMsg.MessageId{
				case message.NEW_ORDER:
					if !que.HasOrder(newMsg.Button){
						que.AddOrder(newMsg.Button)
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

				}				
			}
		}

	}()


}