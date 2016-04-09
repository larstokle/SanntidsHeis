package main

import(
	."./UDP"
	"fmt"
	"reflect"
	"./eventmgr"
	
)



func main() {
	kanal := make(chan []byte)//make ( chan (Event_t))
	mottaker := make(chan []byte) //make ( chan (Event_t))
	stopp := make(chan (bool))
	MakeReciever(":20000", mottaker, stopp)

	MakeSender("129.241.187.150:20000", kanal, stopp)
	fmt.Println(GetOwnID())
	i:= 0
	//var melding msg
	for {
		tosend := eventmgr.Event_t{1, i}
		kanal <- Pack(tosend)

		//kanal <- Pack(i)
		newData := Parse(<-mottaker)
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
			
	}


	
}