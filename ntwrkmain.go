package main

import(
	."./network"
	"fmt"
	"reflect"
	"./eventmgr"
	"./transactionmanager"
	"time"
	
)



func main() {
	kanal := make(chan []byte)//make ( chan (Event_t))
	mottaker := make(chan []byte) //make ( chan (Event_t))
	stopp := make(chan (bool))
	MakeReceiver(":20000", mottaker, stopp)


	MakeSender("129.241.187.150:20000", kanal, stopp)

	i:= 0
	
	transactionmanager.StartTransactionManager()

	for {
		//tosend := eventmgr.Event_t{1, i}
		//kanal <- transactionmanager.Pack(tosend)

		//kanal <- Pack(i)
		newData := transactionmanager.Parse(<-mottaker)
		switch data := newData.(type){
		case eventmgr.Event_t:
			//newEvent := newData.(Event_t)
			fmt.Printf("Event_t found: %+v\n", data)	

		case int:
			//newInt
			fmt.Printf("Int found: %+v\n", newData)
		default:
			fmt.Printf("random found: %+v, with type: %+v\n", newData, reflect.TypeOf(newData))
		}
	

		i++
		time.Sleep(time.Millisecond*2000)
	}


	
}