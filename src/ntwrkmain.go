package main

import(
	"network"
	."message"

	"time"
	//"fmt"

	
)



func main() {
	ticker := time.NewTicker(time.Millisecond * 3000)
	sendChan := network.MakeSender("129.241.187.255:20777")
	recChan := network.MakeReceiver(":20777")
	defer close(sendChan)
	for _ = range ticker.C{
		sendChan <- Message_t{}
		<- recChan	
	}

}