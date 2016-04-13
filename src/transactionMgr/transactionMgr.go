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
	Receive         chan message.Message_t //bad naming. it is the output from network
	netReceive      chan message.Message_t
	netSend         chan message.Message_t
	heartbeatTimers map[int]*time.Timer
	heartbeatMutex  sync.Mutex
	myId            int
	delegatingOrder bool
	delegation      map[Button_t]map[int]costAndToId_t
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
	go func() {
		for {

			receivedData := <-transMgr.netReceive
			//are receivedData.Source in heartbeatTimers??
			switch receivedData.MessageId {
			case message.HEARTBEAT:
				beat := Heartbeat_t{Id: receivedData.Source}
				transMgr.NewHeartBeat(beat)
			case message.NEW_ORDER: //new event
				if receivedData.Source != transMgr.myId {
					fmt.Printf("transMgr: Received NEW_ORDER: %+v\n", receivedData)
					transMgr.Receive <- receivedData
				}
			case message.REQUEST_ORDER:
				if receivedData.Source != transMgr.myId {
					order := receivedData.Button
					cost := receivedData.Cost
					id := receivedData.Source
					if transMgr.delegation[order] == nil {
						fmt.Printf("transMgr: Receiving a request on order %+v, with cost %d\n ", order, cost)
						transMgr.delegation[order] = make(map[int]costAndToId_t)
						transMgr.delegation[order][id] = costAndToId_t{cost: cost, toId: NONLEGAL_ID}
						transMgr.Receive <- receivedData
					} else {
						fmt.Println("ẗransMgr: REQUEST_ORDER Receiving an already requested order")
					}

				}
			case message.COST:
				order := receivedData.Button
				id := receivedData.Source
				newCost := receivedData.Cost
				if oldCostAndToId, present := transMgr.delegation[order][id]; !present {
					transMgr.delegation[order][id] = costAndToId_t{cost: newCost, toId: NONLEGAL_ID}
					if len(transMgr.delegation[order]) == transMgr.nElevetorsOnline() {
						lowestCostId := 256
						lowestCost := 100 * N_FLOORS
						for id, costAndToId := range transMgr.delegation[order] {
							if costAndToId.cost < lowestCost || (costAndToId.cost == lowestCost && id < lowestCostId) {
								lowestCostId = id
								lowestCost = costAndToId.cost
							}
						}
						fmt.Println("DELEGATE NOW FUCKER!!!!")
						transMgr.DelegateOrder(order, lowestCostId)
						tempCostAndToId := transMgr.delegation[order][transMgr.myId]
						tempCostAndToId.toId = lowestCostId
						transMgr.delegation[order][transMgr.myId] = tempCostAndToId
					}
				} else {
					fmt.Printf("transMgr: got multiple cost on order %+v from %d. oldCostAndToId = %d, newCost got = %d\n", order, id, oldCostAndToId, newCost)
				}
			case message.DELEGATE_ORDER:
				fmt.Printf("transMgr: DELEGATE_ORDER (%+v) from %d to %d\n", receivedData.Button, receivedData.Source, receivedData.ElevatorId)
				order := receivedData.Button
				id := receivedData.Source
				toId := receivedData.ElevatorId
				if _, present := transMgr.delegation[order][id]; present {
					tempCostAndId := transMgr.delegation[order][id]
					tempCostAndId.toId = toId
					transMgr.delegation[order][id] = tempCostAndId
					nDelegated := 0
					allDelegatedEqual := true
					for _, costAndToId := range transMgr.delegation[order] {
						if costAndToId.toId != NONLEGAL_ID {
							nDelegated++
						}
						if transMgr.delegation[order][transMgr.myId].toId != costAndToId.toId {
							allDelegatedEqual = false
							break
						}
					}
					if !allDelegatedEqual {
						fmt.Printf("trasnMgr: allDelegatedEqual = false, delegation[%+v] = %+v\n", order, transMgr.delegation[order])
					} else if nDelegated == len(transMgr.delegation[order]) {
						fmt.Printf("tranMgr: DELEGATE_ORDER completed. delegated order %+v to elevator %d\n", order, transMgr.delegation[order][transMgr.myId].toId)
						transMgr.Receive <- message.Message_t{MessageId: message.DELEGATE_ORDER, Button: order, ElevatorId: transMgr.delegation[order][transMgr.myId].toId}
					}
				}

			case message.REMOVE_ORDER:
				if receivedData.Source != transMgr.myId {
					fmt.Printf("transMgr: Received REMOVE_ORDER: %+v\n", receivedData)
					transMgr.Receive <- receivedData
				}
			default:
				fmt.Printf("transMgr: received unhandled MessageId \n", receivedData.MessageId)
			}
		}

	}()

	return &transMgr
}

func (transMgr transactionMgr_t) startHeartbeat() {
	go func() {
		for {
			//fmt.Printf("Sending Heartbeat: %+v \n", beat)
			transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.HEARTBEAT}
			time.Sleep(time.Millisecond * 200)
		}
	}()
}

func (transMgr *transactionMgr_t) NewHeartBeat(beat Heartbeat_t) {
	if _, exists := transMgr.heartbeatTimers[beat.Id]; exists {
		transMgr.heartbeatMutex.Lock()
		transMgr.heartbeatTimers[beat.Id].Reset(time.Second * 1)
		transMgr.heartbeatMutex.Unlock()
	} else {
		transMgr.heartbeatMutex.Lock()
		transMgr.heartbeatTimers[beat.Id] = time.AfterFunc(time.Second*1, func() { transMgr.RemoveElevator(beat.Id) })
		transMgr.nElevetorsOnline++
		transMgr.heartbeatMutex.Unlock()
		fmt.Printf("Got New Heartbeat ID: %+v \n", beat)
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
	transMgr.nElevetorsOnline--
	delete(transMgr.heartbeatTimers, id)
	transMgr.heartbeatMutex.Unlock()
	fmt.Printf("Lost Heartbeat ID: %+v \n", id)
}

func (transMgr *transactionMgr_t) DelegateOrder(order Button_t, id int) {
	transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.DELEGATE_ORDER, Button: order, ElevatorId: id}
	fmt.Printf("transMgr: Delegate order (%+v) to id %d\n", order, id)
}

func (transMgr *transactionMgr_t) RequestOrder(order Button_t, cost int) {
	if transMgr.delegation[order] == nil {
		fmt.Printf("transMgr: Starting a request on order %+v, with cost %d\n ", order, cost)
		transMgr.delegation[order] = make(map[int]costAndToId_t)
		transMgr.delegation[order][transMgr.myId] = costAndToId_t{cost: cost, toId: NONLEGAL_ID}
		transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.REQUEST_ORDER, Button: order, Cost: cost}
	} else {
		fmt.Println("ẗransMgr: Requesting an already requested order")
	}
}

func (transMgr *transactionMgr_t) MyId() int {
	return transMgr.myId
}

func (transMgr *transactionMgr_t) NewOrder(order Button_t) {
	fmt.Printf("transMgr: sending new order on network = %+v\n", order)
	transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.NEW_ORDER, Button: order}
}

func (transMgr *transactionMgr_t) RemoveOrder(floor int) {
	fmt.Printf("transMgr: sending remove order (floor) on network = %+v\n", floor)
	order := Button_t{Floor: floor, ButtonType: UP}
	transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.REMOVE_ORDER, Button: order}
}

func (transMgr *transactionMgr_t) Cost(order Button_t, cost int) { //Rename til send cost?
	fmt.Printf("transMgr: sending cost (= %d) on order (= %+v) on network\n", cost, order)
	transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.COST, Button: order, Cost: cost}
}
