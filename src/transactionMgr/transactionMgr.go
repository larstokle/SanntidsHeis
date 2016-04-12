package transactionMgr

import (
	"network"
	"time"
	"fmt"
	"message" 
	."globals"
	)

const port = ":20001"
const broadCastAddr = "129.241.187.255"

type Heartbeat_t struct{
	Id int
}

type transactionMgr_t struct{
	Receive chan message.Message_t //bad naming. it is the output from network
	netReceive chan message.Message_t
	netSend chan message.Message_t
	heartbeatTimers map[int]*time.Timer
	myId int
	delegatingOrder bool
	delegation map[Button_t]map[int]int
}

func New() *transactionMgr_t{
	var transMgr transactionMgr_t
	transMgr.Receive = make(chan message.Message_t)
	transMgr.netSend, _ = network.MakeSender(broadCastAddr + port)
	transMgr.netReceive, _ = network.MakeReceiver(port)
	transMgr.heartbeatTimers = make(map[int]*time.Timer)
	transMgr.delegation = make(map[Button_t]map[int]int)
	transMgr.myId = network.GetLastIPByte()

	transMgr.startHeartbeat()
	go func(){
		for {
			
			receivedData := <-transMgr.netReceive
			//are receivedData.Source in heartbeatTimers??
			switch receivedData.MessageId{
			case message.HEARTBEAT:
				beat := Heartbeat_t{Id: receivedData.Source}
				transMgr.NewHeartBeat(beat)
			case message.NEW_ORDER: //new event
				if receivedData.Source != transMgr.myId{
					fmt.Printf("transMgr: Received NEW_ORDER: %+v\n",receivedData)
					transMgr.Receive <- receivedData
				}
			case message.REQUEST_ORDER:
				if receivedData.Source != transMgr.myId{
					order := receivedData.Button
					cost := receivedData.Cost
					id := receivedData.Source
					if transMgr.delegation[order] == nil{
						fmt.Printf("transMgr: Receiving a request on order %+v, with cost %d\n ",order, cost)
						transMgr.delegation[order] = make(map[int]int)
						transMgr.delegation[order][id] = cost
						transMgr.Receive <- receivedData
					}else {
						fmt.Println("ẗransMgr: REQUEST_ORDER Receiving an already requested order")
					}
					
				}
				continue
			case message.COST:
				order := receivedData.Button
				id := receivedData.Source
				newCost := receivedData.Cost
				if oldCost , present := transMgr.delegation[order][id]; !present{
					transMgr.delegation[order][id] = newCost
					if len(transMgr.delegation[order]) == len(transMgr.heartbeatTimers){
						lowestCostId := 256
						lowestCost := 1000
						for id, cost := range transMgr.delegation[order]{
							if cost < lowestCost || (cost == lowestCost && id < lowestCostId){
								lowestCostId = id
								lowestCost = cost
							}
						}
						fmt.Println("DELEGATE NOW FUCKER!!!!")
						transMgr.DelegateOrder(order,lowestCostId)
						transMgr.delegation[order] = nil
					}
				} else{
					fmt.Printf("transMgr: got multiple cost on order %+v from %d. oldCost = %d, newCost got = %d\n",order, id, oldCost, newCost)
				}
			case message.DELEGATE_ORDER:
				fmt.Printf("transMgr: DELEGATE_ORDER (%+v) from %d to %d\n", receivedData.Button, receivedData.Source, receivedData.ElevatorId)
			case message.REMOVE_ORDER:
				if receivedData.Source != transMgr.myId{
					fmt.Printf("transMgr: Received REMOVE_ORDER: %+v\n",receivedData)
					transMgr.Receive <- receivedData
				}
			default:
				fmt.Printf("transMgr: received unhandled MessageId \n",receivedData.MessageId)
			}
		}		
		
	}()

	return &transMgr
}

func (transMgr transactionMgr_t) startHeartbeat(){
	go func(){
		for {
			//fmt.Printf("Sending Heartbeat: %+v \n", beat)
			transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.HEARTBEAT} 
			time.Sleep(time.Millisecond*200)
		}
	}()
}

func (transMgr *transactionMgr_t) NewHeartBeat(beat Heartbeat_t){
	if _, exists := transMgr.heartbeatTimers[beat.Id]; exists{
		transMgr.heartbeatTimers[beat.Id].Reset(time.Second*1)
	}else{
		transMgr.heartbeatTimers[beat.Id] = time.AfterFunc(time.Second*1, func(){transMgr.RemoveElevator(beat.Id)})
		fmt.Printf("Got New Heartbeat ID: %+v \n",beat)
	}
}

func (transMgr *transactionMgr_t) RemoveElevator(id int){
	delete(transMgr.heartbeatTimers,id)
	fmt.Printf("Lost Heartbeat ID: %+v \n",id)
}

func (transMgr *transactionMgr_t) DelegateOrder(order Button_t, id int){
	transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.DELEGATE_ORDER, Button: order, ElevatorId: id}
	fmt.Printf("transMgr: Delegate order (%+v) to id %d\n", order, id)
}

func  (transMgr *transactionMgr_t) RequestOrder(order Button_t, cost int){
	if transMgr.delegation[order] == nil{
		fmt.Printf("transMgr: Starting a request on order %+v, with cost %d\n ",order, cost)
		transMgr.delegation[order] = make(map[int]int)
		transMgr.delegation[order][transMgr.myId] = cost
		transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.REQUEST_ORDER, Button: order, Cost: cost}
	}else {
		fmt.Println("ẗransMgr: Requesting an already requested order")
	}
}

func (transMgr *transactionMgr_t) MyId() int{
	return transMgr.myId
}

func (transMgr *transactionMgr_t) NewOrder(order Button_t){
	fmt.Printf("transMgr: sending new order on network = %+v\n", order)
	transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.NEW_ORDER, Button: order}
}

func (transMgr *transactionMgr_t) RemoveOrder(floor int){
	fmt.Printf("transMgr: sending remove order (floor) on network = %+v\n", floor)
	order := Button_t{Floor: floor, ButtonType: UP}
	transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.REMOVE_ORDER, Button: order}
}

func (transMgr *transactionMgr_t) Cost(order Button_t, cost int){
	fmt.Printf("transMgr: sending cost (= %d) on order (= %+v) on network\n", cost, order)
	transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.COST, Button: order, Cost: cost}
}