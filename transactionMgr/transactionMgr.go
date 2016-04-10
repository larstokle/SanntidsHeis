package transactionMgr

import (
	"../network"
	"time"
	"fmt"
	"reflect"
	"../eventmgr"
 
	)

const port = ":20001"
const broadCastAddr = "129.241.187.255"

type Heartbeat_t struct{
	Id int
}

type elevatorMap_t map[int]*time.Timer

func Start() ( chan interface{} , chan interface{}){
	fromClient := make(chan interface{})
	toClient := make(chan interface{})
	
	go func(){
		netSend := make(chan []byte)
		netReceive := make(chan []byte)
		netStop := make(chan bool)
		network.MakeSender(broadCastAddr + port, netSend, netStop)
		network.MakeReceiver(port, netReceive, netStop)
		startHeartbeat(netSend)
		elevatorsOnline := make(elevatorMap_t)

		for {
			select{
				case input := <- fromClient:
					netSend <- Pack(input)

				case receivedData := <-netReceive:
					receivedUnknownType := Parse(receivedData)
					//check if sender is self
					switch received := receivedUnknownType.(type){
					case Heartbeat_t:
						elevatorsOnline.NewHeartBeat(received)
					case eventmgr.Event_t:
						toClient <- received
					default:
						fmt.Printf("transactionmanager received unhandled type %+v. Received: %+v\n",reflect.TypeOf(receivedUnknownType) , receivedUnknownType)
					}

				default:
					break
			}
			
		}
	}()

	return fromClient, toClient
}

func startHeartbeat(netSend chan []byte){
	go func(){
		beat := Heartbeat_t{Id: network.GetLastIPByte()}

		for {
			//fmt.Printf("Sending Heartbeat: %+v \n", beat)
			netSend <- Pack(beat)
			time.Sleep(time.Millisecond*200)
		}
	}()
}

func (elevatorsOnline *elevatorMap_t) NewHeartBeat(beat Heartbeat_t){
	if _, exists := (*elevatorsOnline)[beat.Id]; exists{
		(*elevatorsOnline)[beat.Id].Reset(time.Second*1)
	}else{
		(*elevatorsOnline)[beat.Id] = time.AfterFunc(time.Second*1, func(){(*elevatorsOnline).RemoveElevator(beat.Id)})
		fmt.Printf("Got New Heartbeat ID: %+v \n",beat)
	}
}

func (elevatorsOnline *elevatorMap_t) RemoveElevator(id int){
	delete((*elevatorsOnline),id)
	fmt.Printf("Lost Heartbeat ID: %+v \n",id)
}