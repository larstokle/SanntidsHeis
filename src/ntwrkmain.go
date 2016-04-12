package main

import(
	"./network"
	"./message"

	"time"
	"fmt"

	
)



func main() {
	kanal := make(chan message.Message_t)
	mottaker := make(chan message.Message_t) 
	stopp := make(chan (bool))
	
	network.MakeReceiver(":20000", mottaker, stopp)
	network.MakeSender("10.22.68.20:20000", kanal, stopp)

	i:= 0
	
	
	for {

		kanal <- message.Message_t{Source: network.GetLastIPByte(), Message_id: message.HEARTBEAT}
		fmt.Printf("Recieved: %+v\n", <-mottaker)
		//fmt.Printf("GetLastIPByte: %d \n", network.GetLastIPByte())

		i++
		time.Sleep(time.Millisecond*2000)
	}


	
}