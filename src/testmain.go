package main

import(
	"./elevatorMgr"
	"fmt"
)


func main() {
	fmt.Println("\n\n\n=========================== ELEVATOR STARTUP ===========================\n")
	elevatorMgr.Start()
	fmt.Println("=========================== ELEVATOR STARTED ===========================\n\n\n")
	select{}
}



