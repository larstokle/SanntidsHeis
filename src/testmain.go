package main

import(
	"elevatorMgr"
	"fmt"
)

func init(){
	fmt.Println("\n\n\n=========================== ELEVATOR STARTUP ===========================\n")
}

func main() {
	//TODO: get commandline input ad set params
	//TODO: start something like a process pair
	
	elevatorMgr.Start()
	fmt.Println("=========================== ELEVATOR STARTED ===========================\n\n\n")
	//TODO: Handle process pair
	//TODO: Handle exit
	select{}
}



