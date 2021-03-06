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
	ToParent        chan message.Message_t
	ProceedOk       chan bool
	netReceive      chan message.Message_t
	netSend         chan<- message.Message_t
	heartbeatTimers map[int]*time.Timer
	heartbeatMutex  sync.Mutex
	myId            int

	orderControl map[Button_t]map[int]costAndToId_t // delegationMap/control ens?
	//=================================================
	delegationMutex sync.Mutex
}

func New() *transactionMgr_t {
	var transMgr transactionMgr_t
	transMgr.ToParent = make(chan message.Message_t)
	transMgr.ProceedOk = make(chan bool)
	transMgr.netSend = network.MakeSender(broadCastAddr + port)
	transMgr.netReceive = network.MakeReceiver(port)
	transMgr.heartbeatTimers = make(map[int]*time.Timer)
	transMgr.orderControl = make(map[Button_t]map[int]costAndToId_t)
	transMgr.myId = network.GetLastIPByte()

	heartBeat := time.Tick(time.Millisecond * 500)
	fmt.Println("transMgr: init done entering loop")
	go func() {
		for {
			var newMsg message.Message_t
			select {
			case newMsg = <-transMgr.netReceive:
			case <-heartBeat:
				transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.HEARTBEAT}
				continue
			}

			if newMsg.Source == transMgr.myId && newMsg.MessageId != message.HEARTBEAT {
				continue
			}

			switch newMsg.MessageId {

			case message.HEARTBEAT:
				beat := Heartbeat_t{Id: newMsg.Source}
				transMgr.newHeartBeat(beat)

			case message.NEW_ORDER:
				if DEBUG_TRNSMGR {
					fmt.Printf("transMgr: Received NEW_ORDER %+v from %d\n", newMsg.Button, newMsg.Source)
				}
				transMgr.ToParent <- newMsg
				transMgr.waitForParent()

			case message.COST:
				order := newMsg.Button
				cost := newMsg.Cost
				id := newMsg.Source
				if DEBUG_TRNSMGR {
					fmt.Printf("transMgr: Received COST msg on order %+v from %d with cost %d \n", order, id, cost)
				}
				transMgr.costToDelegation(order, cost, id)

			case message.DELEGATE_ORDER:
				if DEBUG_TRNSMGR {
					fmt.Printf("transMgr: Received DELEGATE_ORDER msg (%+v) from %d to %d\n", newMsg.Button, newMsg.Source, newMsg.ElevatorId)
				}
				order := newMsg.Button
				id := newMsg.Source
				toId := newMsg.ElevatorId

				transMgr.delegationMutex.Lock()
				if tempCostAndId, present := transMgr.orderControl[order][id]; present {
					tempCostAndId.toId = toId
					transMgr.orderControl[order][id] = tempCostAndId

					nDelegated := 0
					allDelegatedEqual := true
					for _, costAndToId := range transMgr.orderControl[order] {
						if costAndToId.toId == NONLEGAL_ID {
							break
						} else if transMgr.orderControl[order][transMgr.myId].toId != costAndToId.toId {
							allDelegatedEqual = false
							break
						}
						nDelegated++
					}

					if !allDelegatedEqual {
						fmt.Printf("ERROR! transMgr: allDelegatedEqual = false, delegation[%+v] = %+v\n", order, transMgr.orderControl[order])
						delete(transMgr.orderControl, order)
						transMgr.ToParent <- message.Message_t{MessageId: message.DELEGATE_ORDER, Button: order, ElevatorId: NONLEGAL_ID}

					} else if nDelegated == len(transMgr.orderControl[order]) { //<<<<<<<<<<<<<<<<<<<<<<<<DETTE MÅ VI TENKE MER GJENNOM!!!!
						if DEBUG_TRNSMGR {
							fmt.Printf("transMgr: allDelegatedEqual = true. delegated order %+v to elevator %d\n", order, transMgr.orderControl[order][transMgr.myId].toId)
						}

						transMgr.ToParent <- message.Message_t{MessageId: message.DELEGATE_ORDER, Button: order, ElevatorId: transMgr.orderControl[order][transMgr.myId].toId}
						delete(transMgr.orderControl, order)

					}
				} else {
					fmt.Printf("ERROR! transMgr: Received DELEGATE_ORDER (%+v) where cost not set from %d\n", order, id)
				}
				transMgr.delegationMutex.Unlock()

			case message.REMOVE_ORDER:
				if DEBUG_TRNSMGR {
					fmt.Printf("transMgr: Received REMOVE_ORDER %+v from \n", newMsg.Button, newMsg.Source)
				}

				transMgr.ToParent <- newMsg

				order := newMsg.Button
				order.ButtonType = UP
				transMgr.delegationMutex.Lock()
				delete(transMgr.orderControl, order)
				order.ButtonType = DOWN
				delete(transMgr.orderControl, order)
				transMgr.delegationMutex.Unlock()

			case message.UNASSIGN_ORDER:
				transMgr.ToParent <- newMsg

			case message.SYNC:
				transMgr.ToParent <- newMsg
				transMgr.waitForParent()

			default:
				fmt.Printf("transMgr: ERROR! received unhandled MessageId \n", newMsg.MessageId)
			}
		}

	}()
	return &transMgr
}

func (transMgr *transactionMgr_t) newHeartBeat(beat Heartbeat_t) {
	transMgr.heartbeatMutex.Lock()
	if _, exists := transMgr.heartbeatTimers[beat.Id]; exists {
		transMgr.heartbeatTimers[beat.Id].Reset(time.Millisecond * 1500)
		transMgr.heartbeatMutex.Unlock()
	} else {
		transMgr.heartbeatTimers[beat.Id] = time.AfterFunc(time.Millisecond*1500, func() { transMgr.lostHeartBeat(beat.Id) })
		transMgr.heartbeatMutex.Unlock()
		fmt.Printf("===Got New Heartbeat ID: %+v, now have %d elevs===\n", beat, transMgr.nElevatorsOnline())
		if beat.Id != transMgr.myId {
			transMgr.ToParent <- message.Message_t{MessageId: message.HEARTBEAT}
		}
	}
}

func (transMgr *transactionMgr_t) lostHeartBeat(id int) {
	transMgr.restartDelegations()
	transMgr.heartbeatMutex.Lock()
	delete(transMgr.heartbeatTimers, id)
	transMgr.heartbeatMutex.Unlock()

	fmt.Printf("===Lost Heartbeat ID: %+v, now have %d elevs===\n", id, transMgr.nElevatorsOnline())

	if id != transMgr.myId {
		transMgr.ToParent <- message.Message_t{Source: transMgr.myId, MessageId: message.UNASSIGN_ORDER, ElevatorId: id}

	} else {
		fmt.Printf("===transMgr: Lost my own heartbeat, all alone in the world\n")
	}
}

func (transMgr *transactionMgr_t) nElevatorsOnline() int {
	transMgr.heartbeatMutex.Lock()
	nElevators := len(transMgr.heartbeatTimers)
	transMgr.heartbeatMutex.Unlock()
	return nElevators
}

func (transMgr *transactionMgr_t) setLowestCostId(order Button_t) {
	lowestCostId := 256
	lowestCost := 100 * N_FLOORS
	transMgr.delegationMutex.Lock()
	for id, costAndToId := range transMgr.orderControl[order] {
		if costAndToId.cost < lowestCost || (costAndToId.cost == lowestCost && id < lowestCostId) {
			lowestCostId = id
			lowestCost = costAndToId.cost
		}
	}
	tempCostAndToId := transMgr.orderControl[order][transMgr.myId]
	tempCostAndToId.toId = lowestCostId
	transMgr.orderControl[order][transMgr.myId] = tempCostAndToId

	transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.DELEGATE_ORDER, Button: order, ElevatorId: lowestCostId}

	transMgr.delegationMutex.Unlock()
	if DEBUG_TRNSMGR {
		fmt.Printf("transMgr: lowest cost on order %+v found on id %d, delegation sendt on net\n", order, lowestCostId)
	}
}

func (transMgr *transactionMgr_t) costToDelegation(order Button_t, cost int, id int) {
	transMgr.delegationMutex.Lock()

	if transMgr.orderControl[order] == nil {
		if DEBUG_TRNSMGR {
			fmt.Printf("transMgr: New delegation sequence on order %+v, setting cost %d to id %d\n", order, cost, id)
		}

		transMgr.orderControl[order] = make(map[int]costAndToId_t)
		transMgr.orderControl[order][id] = costAndToId_t{cost: cost, toId: NONLEGAL_ID}
		transMgr.delegationMutex.Unlock()

		if id != transMgr.myId {
			transMgr.ToParent <- message.Message_t{MessageId: message.COST, Button: order}
		}

	} else if oldCostAndToId, present := transMgr.orderControl[order][id]; !present {
		if DEBUG_TRNSMGR {
			fmt.Printf("transMgr: Setting cost on order %+v, with cost %d to %d in existing delegation sequence\n", order, cost, id)
		}

		transMgr.orderControl[order][id] = costAndToId_t{cost: cost, toId: NONLEGAL_ID}

		pendingElevs := transMgr.nElevatorsOnline() - len(transMgr.orderControl[order])
		transMgr.delegationMutex.Unlock()
		if pendingElevs <= 0 {
			transMgr.setLowestCostId(order)

		} else if DEBUG_TRNSMGR {
			fmt.Printf("transMgr: No delegation yet, still waiting for %d other elevs\n", pendingElevs)
		}

	} else if pendingElevs := transMgr.nElevatorsOnline() - len(transMgr.orderControl[order]); pendingElevs <= 0 {
		fmt.Printf("transMgr: got multiple cost, but no pending elevs\n")
		transMgr.delegationMutex.Unlock()
		transMgr.setLowestCostId(order)
	} else {
		transMgr.delegationMutex.Unlock()
		fmt.Printf("ERROR! transMgr: got multiple cost on order %+v from %d. oldCostAndToId = %d, newCost got = %d\n", order, id, oldCostAndToId, cost)
	}
}

func (transMgr *transactionMgr_t) waitForParent() {
	<-transMgr.ProceedOk
}

func (transMgr *transactionMgr_t) restartDelegations() {
	transMgr.delegationMutex.Lock()
	for key, _ := range transMgr.orderControl {
		delete(transMgr.orderControl, key)
	}
	transMgr.delegationMutex.Unlock()
}
