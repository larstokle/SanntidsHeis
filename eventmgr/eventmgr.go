package eventmgr

import (
	. "../constants"
	"../driver"
	"strconv"
	"time"
)

type Event_t struct {
	Floor     int
	EventType int
}

var eventTypes = [...]string{
	"Up",
	"Down",
	"Command",
	"Floor Signal",
}

func (event Event_t) String() string {
	return "Floor:" + strconv.Itoa(event.Floor) + " " + eventTypes[event.EventType]
}

func CheckEvents(event chan Event_t) {
	driver.Init() //mulig flyttes til main?

	go checkButtons(event)

	go checkFloorSignal(event)
}

func checkButtons(event chan Event_t) {
	var lastButtonState [N_FLOORS][N_BUTTON_TYPES]bool
	floor, button := 0, 0

	for true {
		pressed := driver.ReadButton(button, floor)
		if pressed != lastButtonState[floor][button] {
			lastButtonState[floor][button] = pressed
			if pressed {
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
}

func checkFloorSignal(event chan Event_t) {
	lastFloorState := -1

	for true {
		newFloorState := driver.GetFloorSignal()
		if newFloorState != lastFloorState {
			lastFloorState = newFloorState
			if newFloorState != -1 {
				var newEvent Event_t
				newEvent.Floor = newFloorState
				newEvent.EventType = FLOOR_SIGNAL
				event <- newEvent
			}
		}
	}
}
