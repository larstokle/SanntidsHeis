package elevatorMgr

import (
	"orderque"
	"fsm"
	"eventmgr"
	"transactionMgr"
)



func Start(){
	que := orderque.NewOrderQue()
	localElev := fsm.NewElevator()
	btnPush := eventmgr.CheckButtons()

	transMgr := transactionMgr.New()

	go func(){//hmmm skal denne kjøre selv eller skal det go'es i main??
		for{
			select{
			case newBtn := <- btnPush:
				if !que.HasOrder(newBtn){
					que.AddOrder(newBtn)
					cost := 3 //localElev.GetCost(newBtn)
					transMgr.DelegateOrder(newBtn, cost) //Denne kan ta inn func(){que.AssignOrder(ID)} 
				}
			/*case newFromTrans := <- transMgr.Receive:
				switch newFromTrans.Id{
				case WANT:
					order := newFromTrans.Button
					cost := localElev.GetCost(order)
					transMgr.StartDelegateOrder(order, cost)
				}
				//check type
				//if btn: que.AddOrder(data)
				//cost := localElev.GetCost(newBtn)
				//cost ok for this?
				//yes: toTrans <- {order, cost}
				//no: do nothing
				//if someone wants new order
				//calculate cost and compareget 
				//toTrans <- Ok/notOK
				break*/
			default:
				if localElev.NeedNewDestination() {

				}
				break
			}
		}

	}()


}