package main

import(
	."./UDP"
	"fmt"
	
)

func main() {
	kanal := make(chan interface{})//make ( chan (Event_t))
	mottaker := make(chan interface{}) //make ( chan (Event_t))
	stopp := make(chan (bool))
	MakeReciever(":20000",mottaker, stopp)

	MakeSender("10.22.70.156:20000", kanal, stopp)
	
	i:= 0
	//var melding msg
	for {
		kanal <- Event_t{1,i}
		//melding <-mottaker
		mellom := <-mottaker
		switch mellom.(type){
		case Event_t:
			fmt.Printf("Event_t: %+v \n", mellom)	

		case int:
			fmt.Printf("Int: %+v\n", mellom)	
		}
		i++
			
	}


	
}