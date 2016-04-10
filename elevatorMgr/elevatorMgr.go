package elevatorMgr

import (
	"../orderque"
	"../fsm"
	"../eventmgr"
	"../transactionMgr"
)



func Start(){
	que := orderque.NewOrderQue()
	localElev := fsm.NewElevator()
	btnPush := eventmgr.CheckButtons()

	toTrans, fromTrans := transactionMgr.Start()

	go func(){//hmmm skal denne kjøre selv eller skal det go'es i main??
		for{
			select{
			case newBtn := <- btnPush:
				que.AddOrder(newBtn)
				cost := localElev.GetCost(newBtn)
				toTrans <- struct{event: newBtn, cost: cost} //Dette MÅ fikses! egen type? cost/ikke cost?

			case newFromTrans := <- fromTrans:
				//check type
				//if btn: que.AddOrder(data)
				//cost := localElev.GetCost(newBtn)
				//cost ok for this?
				//yes: toTrans <- {order, cost}
				//no: do nothing
				//if someone wants new order
				//calculate cost and compare
				//toTrans <- Ok/notOK
				break
			case newFromFsm := <- localElev.NeedCommand:
				break
			}
		}

	}()


}