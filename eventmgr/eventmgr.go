package eventmgr

import (
	. "../constants"
	"../driver"
	"strconv"
	"time"
	"fmt"
)

type Event_t struct {
	Floor     int
	EventType int
}

var eventTypes = [...]string{
	"Up",
	"Down",
	"Command",
}


func (event Event_t) String() string {
	
	if event.EventType < len(eventTypes){
		return "Floor:" + strconv.Itoa(event.Floor) + ", Type: " + eventTypes[event.EventType]
	}else{
		return "Floor:" + strconv.Itoa(event.Floor) + ", Type unknown: " + strconv.Itoa(event.EventType)
	}
}


func CheckButtons() chan Event_t{
	event := make(chan Event_t)
	var lastButtonState [N_FLOORS][N_BUTTON_TYPES]bool
	floor, button := 0, 0

	go func(){
		for true {
			pressed := driver.ReadButton(button, floor)
			if pressed != lastButtonState[floor][button] {
				lastButtonState[floor][button] = pressed
				if pressed {
					fmt.Println("button pressed!")
					var newEvent Event_t //tungvint! lag oneliner!!
					newEvent.Floor = floor
					newEvent.EventType = button
					event <- newEvent
				}
			}

			button++
			button = button % N_BUTTON_TYPES
			if button == 0 {
				floor++
				floor = floor % N_FLOORS
			}

			if floor == 0 {
				time.Sleep(time.Millisecond * 20)
			}
		}
	}()

	return event
}

func CheckFloorSignal() chan int{
	event := make(chan int)
	lastFloorState := -1

	go func(){
		for true {
			newFloorState := driver.GetFloorSignal()
			if newFloorState != lastFloorState {
				lastFloorState = newFloorState
				if newFloorState != -1 {
					event <- newFloorState
				}
			}
		}
	}()
	return event
}
