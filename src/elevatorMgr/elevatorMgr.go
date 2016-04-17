package elevatorMgr

import (
	"eventmgr"
	"fmt"
	"fsm"
	. "globals"
	"message"
	"orderque"
	"time"
	"transactionMgr"
)

func Start() {
	localElev := fsm.NewElevator()
	btnPush := eventmgr.CheckButtons()
	que := orderque.New()
	transMgr := transactionMgr.New()

	loggedQue := orderque.ReadFromLog()
	que.SyncInternal(loggedQue)
	que.SyncExternal(loggedQue)
	que.UnassignAllOrders()
	que.Print()

	retryTimer := time.NewTimer(time.Second * 5)

	lastOrderRequested := NONVALID_BUTTON //ENDRING: FOR Å SENDE MINDRE REQUESTS PÅ SAMME ORDRE

	if DEBUG_ELEVMGR {
		fmt.Printf("elevMgr: init done entering loop\n\n")
	}

	go func() { //hmmm skal denne kjøre selv eller skal det go'es i main??
		for {
			transMgr.ParentReady()

			select {
			case floorDone := <-localElev.OrderDone:
				fmt.Printf("elevMgr: Local elevator done with floor %+v\n", floorDone)
				que.RemoveOrdersOnFloor(floorDone)
				transMgr.RemoveOrder(floorDone)
				lastOrderRequested = NONVALID_BUTTON

			case newBtn := <-btnPush:
				if DEBUG_ELEVMGR {
					fmt.Printf("elevMgr: newBtn case from eventmgr = %+v\n", newBtn)
				}
				if !que.HasOrder(newBtn) {
					que.AddOrder(newBtn)
					transMgr.NewOrder(newBtn)
				}

			case newMsg := <-transMgr.ToParent:
				if DEBUG_ELEVMGR {
					fmt.Printf("elevMgr: newMsg case from transMgr = %+v\n", newMsg)
				}
				switch newMsg.MessageId {

				case message.NEW_ORDER:
					newOrder := newMsg.Button
					if DEBUG_ELEVMGR {
						fmt.Printf("elevMgr: newOrder from transMgr = %+v\n", newOrder)
					}
					if !que.HasOrder(newOrder) {
						que.AddOrder(newOrder)
					} else {
						fmt.Printf("ERROR! elevMgr: new order aldready in que, me or ID %d must be out of sync\n", newMsg.Source)
					}

				case message.DELEGATE_ORDER:
					if DEBUG_ELEVMGR {
						fmt.Printf("elevMgr: got DELEGATE_ORDER from trans on order %+v to %d\n", newMsg.Button, newMsg.ElevatorId)
					}
					que.AssignOrderToId(newMsg.Button, newMsg.ElevatorId)
					if newMsg.ElevatorId == transMgr.MyId() {
						fmt.Printf("elevMgr: Local elevator got delegated order %+v \n", newMsg.Button)
						//lastOrderRequested = NONVALID_BUTTON //=====================================================USIKKER PÅ DENNE, MEN GIR MEG NÅ!!!
						localElev.NewDestination(newMsg.Button.Floor) //<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<>KANAL AV DENNE?? kun brukt her, kan løse litt locks, ikke at vi har noen men						continue
					} else if !que.IsIdAssigned(transMgr.MyId()) && !que.IsEmpty() {
						if DEBUG_ELEVMGR {
							fmt.Printf("elevMgr: did not get delegation and got no destination\n\n")
						}
						//lastOrderRequested = NONVALID_BUTTON //Lurer på om den skal inn her!
					}

				case message.REMOVE_ORDER:
					if DEBUG_ELEVMGR {
						fmt.Printf("elevMgr: got REMOVE_ORDER from trans on floor %d, removing both\n", newMsg.Button.Floor)
					}
					toRemove := newMsg.Button
					toRemove.ButtonType = UP
					que.RemoveOrder(toRemove)
					toRemove.ButtonType = DOWN
					que.RemoveOrder(toRemove)

				case message.COST:
					order := newMsg.Button
					cost := localElev.GetCost(order)
					if DEBUG_ELEVMGR {
						fmt.Printf("elevMgr: got COST from trans on order %+v, i calculated cost = %d \n", order, cost)
					}
					if !que.HasOrder(order) {
						fmt.Printf("ERROR! elevMgr: Reveived COST request from transMgr without having order %+v in queue \n\n", order)
					}
					transMgr.SendCost(order, cost)

				case message.UNASSIGN_ORDER:
					unassignId := newMsg.ElevatorId
					if DEBUG_ELEVMGR {
						fmt.Printf("elevMgr: got UNASSIGN_ORDER from trans on id %d\n", unassignId)
					}

					if unassignId != transMgr.MyId() {
						que.UnassignOrdersToID(unassignId)
					} else {
						fmt.Printf("ERROR! elevMgr: Id is mine, so i did not unassign my orders\n")
					}

				case message.SYNC:
					if DEBUG_ELEVMGR {
						fmt.Printf("elevMgr: got SYNC from trans from elevator %d\n", newMsg.Source)
					}
					queToSyncWith, err := orderque.Decode(newMsg.Data)
					if err != nil {
						fmt.Println("ERROR! orderque.Decode: ", err)
					}
					que.SyncExternal(queToSyncWith)
					que.WriteToLog()
					fmt.Printf("elevMgr: synced with que from id %d\n", newMsg.Source) //====================SKAL VI HA DEBUGPRINT PÅ DENNE ? en ny NOPRINTATALL?
					continue

				case message.HEARTBEAT:
					if DEBUG_ELEVMGR {
						fmt.Printf("elevMgr: got HEARTBEAT from trans, Encoding que\n")
					}
					rawQue, err := orderque.Encode(que)
					if err != nil {
						fmt.Println("ERROR! orderque.Encode: ", err)
					}
					transMgr.SendSync(rawQue)

				default:
					fmt.Printf("ERROR! elevMgr: Unhandled MessageId : %+v", newMsg)
				}
			case <-retryTimer.C:
				if DEBUG_ELEVMGR {
					fmt.Printf("elevMgr: timed out, checking for better orders\n")
				}
			}

			if !que.IsEmpty() { //Try to get a better destination
				lowestCost := fsm.INF_COST
				bestOrder := NONVALID_BUTTON
				for orderType := 0; orderType < N_BUTTON_TYPES; orderType++ {
					for floor := FIRST_FLOOR; floor < N_FLOORS; floor++ {
						order := Button_t{Floor: floor, ButtonType: orderType}
						if que.HasOrder(order) && !que.IsOrderAssigned(order) {
							cost := localElev.GetCost(order)
							if cost <= lowestCost {
								lowestCost = cost
								bestOrder = order
							}
						}
					}
				}

				lastOrderCost := localElev.GetCost(lastOrderRequested)
				if bestOrder == lastOrderRequested && lowestCost < fsm.INF_COST { // ============ NOT DONE!!! NEED SOME EXTRAS
					fmt.Printf("ERROR! elevMgr: could not get the requested order %+v\n", lastOrderRequested)

				} else if lowestCost < fsm.INF_COST && lowestCost < lastOrderCost && bestOrder != lastOrderRequested {
					fmt.Printf("elevMgr: Local elevator requesting order %+v, with cost %d\n", bestOrder, localElev.GetCost(bestOrder)) // ====================SKAL VI HA DEBUGPRINT PÅ DENNE ? en ny NOPRINTATALL?
					transMgr.RequestOrder(bestOrder, lowestCost)
					lastOrderRequested = bestOrder
				}
			}
			retryTimer.Reset(time.Second * 5)

		}
	}()
}
