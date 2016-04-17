package main

import(
	//"encoding/json"
	"network"
	"message"
	"fmt"
	//"bytes"
	//"time"
	"orderque"
	."globals"

)


func main() {
	sendChan := network.MakeSender("129.241.187.146:20777")
	recChan := network.MakeReceiver(":20777")

	q1 := orderque.GetFromLog()
	q1.AddOrder(Button_t{0,0})
	q1.Print()

	msg := message.Message_t{MessageId: message.SYNC}
	var err error
	msg.Data, err = orderque.Encode(q1)
	if err != nil{
		fmt.Println(err)
	}
	

	fmt.Println("Sending\n",string(msg.Data), "\n\n\n")
	sendChan <- msg

	msg = <- recChan
	fmt.Print("Received\n",string(msg.Data), "\n\n\n")

	q2 := orderque.New()
	q3, err := orderque.Decode(msg.Data)
	if err != nil{
		fmt.Println(err)
	}
	q2.Sync(q3)

	q2.Print()

	//q2.Log()

}






