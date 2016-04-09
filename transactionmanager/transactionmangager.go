package transactionmanager

import (
	"../network"
	"time"
	"fmt"
	"reflect"
 
	)

const port = ":20001"
const broadCastAddr = "129.241.187.255"

type Hartbeat_t struct{
	Id int
}

type elevatorMap_t map[int]*time.Timer

func StartTransactionManager() ( chan interface{} , chan interface{}){
	inChannel := make(chan interface{})
	outChannel := make(chan interface{})
	
	go func(){
		sendChannel := make(chan []byte)
		receiveChannel := make(chan []byte)
		stopChannel := make(chan bool)
		network.MakeSender(broadCastAddr + port, sendChannel, stopChannel)
		network.MakeReceiver(port, receiveChannel, stopChannel)
		startHartbeat(sendChannel)
		elevators := make(elevatorMap_t)

		for {
			select{
				case input := <- inChannel:
					fmt.Println(input)

				case receivedData := <-receiveChannel:
					receivedUnknownType := Parse(receivedData)
					switch received := receivedUnknownType.(type){
					case Hartbeat_t:
						elevators.NewHartBeat(received)
					default:
						fmt.Printf("transactionmanager received unhandled type %+v. Received: %+v\n",reflect.TypeOf(receivedUnknownType) , receivedUnknownType)
					}

				default:
					break
			}
			
		}
	}()

	return inChannel, outChannel
}

func startHartbeat(sendChannel chan []byte){
	go func(){
		beat := Hartbeat_t{Id: network.GetLastIPByte()}

		for {
			//fmt.Printf("Sending hartbeat: %+v \n", beat)
			sendChannel <- Pack(beat)
			time.Sleep(time.Millisecond*200)
		}
	}()
}

func (elevators *elevatorMap_t) NewHartBeat(beat Hartbeat_t){
	if _, exists := (*elevators)[beat.Id]; exists{
		(*elevators)[beat.Id].Reset(time.Second*1)
	}else{
		(*elevators)[beat.Id] = time.AfterFunc(time.Second*1, func(){(*elevators).RemoveElevator(beat.Id)})
		fmt.Printf("Got New hartbeat ID: %+v \n",beat)
	}
}

func (elevators *elevatorMap_t) RemoveElevator(id int){
	delete((*elevators),id)
	fmt.Printf("Lost hartbeat ID: %+v \n",id)
}