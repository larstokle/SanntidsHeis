package transactionMgr

import (
	"fmt"
	. "globals"
	"message"
	"network"
	"sync"
	"time"
)

const port = ":20777"
const broadCastAddr = "129.241.187.255"

type Heartbeat_t struct {
	Id int
}

type costAndToId_t struct {
	cost int
	toId int
}

type transactionMgr_t struct {
	Receive         chan message.Message_t //bad naming. it is the msg to elevMgr
	netReceive      chan message.Message_t
	netSend         chan message.Message_t
	heartbeatTimers map[int]*time.Timer
	heartbeatMutex  sync.Mutex
	myId            int
	delegation      map[Button_t]map[int]costAndToId_t
	delegationMutex sync.Mutex
}

func New() *transactionMgr_t {
	var transMgr transactionMgr_t
	transMgr.Receive = make(chan message.Message_t)
	transMgr.netSend, _ = network.MakeSender(broadCastAddr + port)
	transMgr.netReceive, _ = network.MakeReceiver(port)
	transMgr.heartbeatTimers = make(map[int]*time.Timer)
	transMgr.delegation = make(map[Button_t]map[int]costAndToId_t)
	transMgr.myId = network.GetLastIPByte()

	transMgr.startHeartbeat()
	if(DEBUG_TRNSMGR){ fmt.Println("transMgr: init done entering loop")}
	go func() {
		for {
			receivedData := <-transMgr.netReceive
			
			switch receivedData.MessageId {

			case message.HEARTBEAT:
				beat := Heartbeat_t{Id: receivedData.Source}
				transMgr.NewHeartBeat(beat)

			case message.NEW_ORDER:
				if receivedData.Source != transMgr.myId {
					if(DEBUG_TRNSMGR){fmt.Printf("transMgr: Received NEW_ORDER: %+v\n", receivedData)}
					if(DEBUG_CHANNELS){fmt.Println("transMgr: NEW_ORDER Receive could hang")}
					transMgr.Receive <- receivedData
					if(DEBUG_CHANNELS){fmt.Println("transMgr: NEW_ORDER Receive didn't hang")}
					
					if(DEBUG_TRNSMGR){fmt.Printf("transMgr: WAIT FOR CONFIRMATION FROM ELEVMGR\n")}
					if(DEBUG_CHANNELS){fmt.Println("transMgr: NEW_ORDER wait for CONFIRMATION could hang")}
					<-transMgr.Receive //WAIT FOR CONFIRMATION FROM ELEVMGR! COULD USE TIMER, BUT THIST IS PROBABLY BETTER
					if(DEBUG_CHANNELS){fmt.Println("transMgr: NEW_ORDER wait for CONFIRMATION didn't hang")}
					if(DEBUG_TRNSMGR){fmt.Printf("transMgr: GOT CONFIRMATION FROM ELEVMGR\n")}
				}

			case message.COST:
				order := receivedData.Button
				cost := receivedData.Cost
				id := receivedData.Source
				if id != transMgr.myId {
					if(DEBUG_TRNSMGR){fmt.Printf("transMgr: Received COST msg on order %+v from %d with cost %d \n", order, id, cost)}
					transMgr.handleSetCost(order, cost, id)
				}

			case message.DELEGATE_ORDER:
				if receivedData.Source == transMgr.myId{
					continue
				}
				if(DEBUG_TRNSMGR){fmt.Printf("transMgr: Received DELEGATE_ORDER msg (%+v) from %d to %d\n", receivedData.Button, receivedData.Source, receivedData.ElevatorId)}
				order := receivedData.Button
				id := receivedData.Source
				toId := receivedData.ElevatorId
				transMgr.delegationMutex.Lock()
				if _, present := transMgr.delegation[order][id]; present {
					tempCostAndId := transMgr.delegation[order][id]
					tempCostAndId.toId = toId
					transMgr.delegation[order][id] = tempCostAndId
					nDelegated := 0
					allDelegatedEqual := true
					for _, costAndToId := range transMgr.delegation[order] {
						if costAndToId.toId == NONLEGAL_ID {
							break
						} else if transMgr.delegation[order][transMgr.myId].toId != costAndToId.toId {
							allDelegatedEqual = false
							break
						}
						nDelegated++
					}
					if !allDelegatedEqual {
						if(DEBUG_TRNSMGR){fmt.Printf("trasnMgr: allDelegatedEqual = false, delegation[%+v] = %+v\n", order, transMgr.delegation[order])}
					} else if nDelegated == len(transMgr.delegation[order]) {
						if(DEBUG_TRNSMGR){fmt.Printf("tranMgr: allDelegatedEqual = true. delegated order %+v to elevator %d\n", order, transMgr.delegation[order][transMgr.myId].toId)}
						
						if(DEBUG_CHANNELS){fmt.Println("transMgr: DELEGATE_ORDER allDelegatedEqual = true could hang")}
						transMgr.Receive <- message.Message_t{MessageId: message.DELEGATE_ORDER, Button: order, ElevatorId: transMgr.delegation[order][transMgr.myId].toId}
						if(DEBUG_CHANNELS){fmt.Println("transMgr: DELEGATE_ORDER allDelegatedEqual = true didn't hang")}
						transMgr.delegation[order] = nil
					}
				} else{
					if(DEBUG_TRNSMGR){fmt.Printf("transMgr: Received DELEGATE_ORDER (%+v) where cost not set from %d\n", order, id)}
				}
				transMgr.delegationMutex.Unlock()

			case message.REMOVE_ORDER:
				if receivedData.Source != transMgr.myId {
					if(DEBUG_TRNSMGR){fmt.Printf("transMgr: Received REMOVE_ORDER: %+v\n", receivedData)}

					if(DEBUG_CHANNELS){fmt.Println("transMgr: REMOVE_ORDER could hang")}
					transMgr.Receive <- receivedData
					if(DEBUG_CHANNELS){fmt.Println("transMgr: REMOVE_ORDER didn't hang")}


					order := receivedData.Button
					order.ButtonType = UP
					transMgr.delegationMutex.Lock()
					delete(transMgr.delegation, order)
					order.ButtonType = DOWN
					delete(transMgr.delegation, order)
					transMgr.delegationMutex.Unlock()
				}
			default:
				if(DEBUG_TRNSMGR){fmt.Printf("transMgr: received unhandled MessageId \n", receivedData.MessageId)}
			}
		}

	}()
	return &transMgr
}

func (transMgr transactionMgr_t) startHeartbeat() {
	go func() {
		for {
			//if(DEBUG_TRNSMGR){fmt.Printf("Sending Heartbeat: %+v \n", beat)}

			transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.HEARTBEAT}
			time.Sleep(time.Millisecond * 500)
		}
	}()
}

func (transMgr *transactionMgr_t) NewHeartBeat(beat Heartbeat_t) {
	transMgr.heartbeatMutex.Lock()
	if _, exists := transMgr.heartbeatTimers[beat.Id]; exists {
		transMgr.heartbeatTimers[beat.Id].Reset(time.Millisecond * 1500)
		transMgr.heartbeatMutex.Unlock()
	} else {
		transMgr.heartbeatTimers[beat.Id] = time.AfterFunc(time.Millisecond * 1500, func() { transMgr.RemoveElevator(beat.Id) })
		transMgr.heartbeatMutex.Unlock()
		if(DEBUG_TRNSMGR){fmt.Printf("===Got New Heartbeat ID: %+v, now have %d elevs===\n", beat, transMgr.nElevatorsOnline())}
	}
}

func (transMgr *transactionMgr_t) nElevatorsOnline() int {
	transMgr.heartbeatMutex.Lock()
	nElevators := len(transMgr.heartbeatTimers)
	transMgr.heartbeatMutex.Unlock()
	return nElevators
}

func (transMgr *transactionMgr_t) RemoveElevator(id int) {
	transMgr.heartbeatMutex.Lock()
	delete(transMgr.heartbeatTimers, id)
	transMgr.heartbeatMutex.Unlock()
	if(DEBUG_TRNSMGR){fmt.Printf("===Lost Heartbeat ID: %+v, now have %d elevs===\n", id, transMgr.nElevatorsOnline())}
}

func (transMgr *transactionMgr_t) DelegateOrder(order Button_t) {
	lowestCostId := 256
	lowestCost := 100 * N_FLOORS
	transMgr.delegationMutex.Lock()
	for id, costAndToId := range transMgr.delegation[order] {
		if costAndToId.cost < lowestCost || (costAndToId.cost == lowestCost && id < lowestCostId) {
			lowestCostId = id
			lowestCost = costAndToId.cost
		}
	}
	tempCostAndToId := transMgr.delegation[order][transMgr.myId]
	tempCostAndToId.toId = lowestCostId
	transMgr.delegation[order][transMgr.myId] = tempCostAndToId

	if(DEBUG_CHANNELS){fmt.Println("transMgr: DelegateOrder could hang")}
	transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.DELEGATE_ORDER, Button: order, ElevatorId: lowestCostId}
	if(DEBUG_CHANNELS){fmt.Println("transMgr: DelegateOrder didn't hang")}
	
	transMgr.delegationMutex.Unlock()
	if(DEBUG_TRNSMGR){fmt.Printf("transMgr: Delegate order %+v to id %d\n", order, lowestCostId)}
}

func (transMgr *transactionMgr_t) RequestOrder(order Button_t, cost int) bool{ //BYGD OM TIL Å RETURNERE TRUE OM HEIS KAN TA DEN MED EN GANG: UNGÅR KANALBRUK OG DERMED LOCK
	if transMgr.nElevatorsOnline() <= 1{
		if(DEBUG_TRNSMGR){fmt.Printf("transMgr: RequestOrder on %+v with cost %d, but no other elevs\n", order, cost)}
	 return true //transMgr.Receive <- message.Message_t{MessageId: message.DELEGATE_ORDER, Button: order, ElevatorId: transMgr.myId}
	} else{
		if(DEBUG_TRNSMGR){fmt.Printf("transMgr: RequestOrder on %+v with cost %d and %d elevs online\n", order, cost, transMgr.nElevatorsOnline())}
		transMgr.SendCost(order, cost)
		return false
	}
}

func (transMgr *transactionMgr_t) MyId() int {
	return transMgr.myId
}

func (transMgr *transactionMgr_t) NewOrder(order Button_t) {
	numElevs := transMgr.nElevatorsOnline()
	if numElevs > 1{
		if(DEBUG_TRNSMGR){fmt.Printf("transMgr: sending new order on network = %+v to in total %d elevs\n", order, numElevs)}
		
		if(DEBUG_CHANNELS){fmt.Println("transMgr: NewOrder could hang")}
		transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.NEW_ORDER, Button: order}
		if(DEBUG_CHANNELS){fmt.Println("transMgr: NewOrder didn't hang")}

	} else{
		if(DEBUG_TRNSMGR){fmt.Printf("transMgr: did not send new order (%+v) since numElevs = %d \n",order, numElevs)}
	}
}

func (transMgr *transactionMgr_t) RemoveOrder(floor int) {
	if(DEBUG_TRNSMGR){fmt.Printf("transMgr: sending remove order (floor) on network = %+v\n", floor)}
	order := Button_t{Floor: floor, ButtonType: UP}
	
	if(DEBUG_CHANNELS){fmt.Println("transMgr: RemoveOrder netSend could hang")}
	transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.REMOVE_ORDER, Button: order}
	if(DEBUG_CHANNELS){fmt.Println("transMgr: RemoveOrder netSend didn't hang")}

	order.ButtonType = UP
	transMgr.delegationMutex.Lock()
	delete(transMgr.delegation, order)
	order.ButtonType = DOWN
	delete(transMgr.delegation, order)
	transMgr.delegationMutex.Unlock()
}

func (transMgr *transactionMgr_t) SendCost(order Button_t, cost int) {
	if(DEBUG_TRNSMGR){fmt.Printf("transMgr: sending cost %d on order %+v on network\n", cost, order)}
	
	if(DEBUG_CHANNELS){fmt.Println("transMgr: SendCost netSend could hang")}
	transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.COST, Button: order, Cost: cost}
	if(DEBUG_CHANNELS){fmt.Println("transMgr: SendCost netSend didn't hang")}

	transMgr.handleSetCost(order, cost, transMgr.myId)
}

func (transMgr *transactionMgr_t) handleSetCost(order Button_t, cost int, id int){
	transMgr.delegationMutex.Lock()
	if transMgr.delegation[order] == nil {
		if(DEBUG_TRNSMGR){fmt.Printf("transMgr: New delegation sequence on order %+v, setting cost %d to id %d\n", order, cost, id)}
		transMgr.delegation[order] = make(map[int]costAndToId_t)
		transMgr.delegation[order][id] = costAndToId_t{cost: cost, toId: NONLEGAL_ID}
		transMgr.delegationMutex.Unlock()
		if id != transMgr.myId{
			if(DEBUG_CHANNELS){fmt.Println("transMgr: handleSetCost Receive could hang")}
			transMgr.Receive <- message.Message_t{MessageId: message.COST, Button: order}
			if(DEBUG_CHANNELS){fmt.Println("transMgr: handleSetCost Receive didn't hang")}
		}
	} else if oldCostAndToId, present := transMgr.delegation[order][id]; !present {
		if(DEBUG_TRNSMGR){fmt.Printf("transMgr: Setting cost on order %+v, with cost %d to %d in existing delegation sequence\n", order, cost, id)}
		transMgr.delegation[order][id] = costAndToId_t{cost: cost, toId: NONLEGAL_ID}
		pendingElevs := transMgr.nElevatorsOnline() - len(transMgr.delegation[order])
		if pendingElevs <= 0 {
			transMgr.delegationMutex.Unlock()
			transMgr.DelegateOrder(order)
		} else {
			transMgr.delegationMutex.Unlock()
			if(DEBUG_TRNSMGR){fmt.Printf("transMgr: No delegation yet, still waiting for %d other elevs\n",pendingElevs)}
		}
	} else {
		transMgr.delegationMutex.Unlock()
		if(DEBUG_TRNSMGR){fmt.Printf("transMgr: got multiple cost on order %+v from %d. oldCostAndToId = %d, newCost got = %d\n", order, id, oldCostAndToId, cost)}
	}
}
