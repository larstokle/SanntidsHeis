package transactionMgr

import (
	"network"
	"time"
	"fmt"
	"reflect"
	"message" 
	."globals"
	)

const port = ":20001"
const broadCastAddr = "129.241.187.255"

type Heartbeat_t struct{
	Id int
}

type transactionMgr_t struct{
	Receive chan Button_t //bad naming. it is the output from network
	netReceive chan message.Message_t
	netSend chan message.Message_t
	idOnline_timers map[int]*time.Timer
	myId int
	delegatingOrder bool
}

func New() *transactionMgr_t{
	var transMgr transactionMgr_t
	transMgr.Receive = make(chan Button_t)
	transMgr.netSend, _ = network.MakeSender(broadCastAddr + port)
	transMgr.netReceive, _ = network.MakeReceiver(port)
	transMgr.idOnline_timers = make(map[int]*time.Timer)
	transMgr.myId = network.GetLastIPByte()

	transMgr.startHeartbeat()
	go func(){
		for {
			
			receivedData := <-transMgr.netReceive
			switch receivedData.MessageId{
			case message.HEARTBEAT:
				beat := Heartbeat_t{Id: receivedData.Source}
				transMgr.NewHeartBeat(beat)
			case message.NEW_ORDER: //new event
				fmt.Printf("transMgr Received NEW_ORDER: %+v",receivedData)
				if receivedData.Source != transMgr.myId{
					transMgr.Receive <- receivedData.Button
				}
			case message.WANTS_ORDER:
				// if receivedData.Source != myId{
				// 	toClient <- struct{ID: WANT, BTN: receivedData.Button}
				// 	transMgr.delegatingOrder = true
				// }
				continue
			case message.COST:
				continue
			default:
				fmt.Printf("transactionmanager received unhandled type %+v. Received: %+v\n",reflect.TypeOf(receivedData) , receivedData)
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
	if _, exists := transMgr.idOnline_timers[beat.Id]; exists{
		transMgr.idOnline_timers[beat.Id].Reset(time.Second*1)
	}else{
		transMgr.idOnline_timers[beat.Id] = time.AfterFunc(time.Second*1, func(){transMgr.RemoveElevator(beat.Id)})
		fmt.Printf("Got New Heartbeat ID: %+v \n",beat)
	}
}

func (transMgr *transactionMgr_t) RemoveElevator(id int){
	delete(transMgr.idOnline_timers,id)
	fmt.Printf("Lost Heartbeat ID: %+v \n",id)
}

func (transMgr *transactionMgr_t) DelegateOrder(order Button_t, cost int){
	transMgr.netSend <- message.Message_t{Source: transMgr.myId, MessageId: message.NEW_ORDER, Button: order}
	fmt.Println("should now DelegateOrder like a PRO!")
}