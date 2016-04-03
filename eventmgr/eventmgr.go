package eventmgr

import (
	"../driver"
	"strconv"
	"time"
)

type Event_t struct {
	Floor     int
	EventType int
}

const (
	UP           = 0
	DOWN         = 1
	CMD          = 2
	FLOOR_SIGNAL = 3 // flyttes?
)

var eventTypes = [...]string{
	"Up",
	"Down",
	"Command",
	"Floor Signal",
}

func (event Event_t) String() string {
	return "Floor:" + strconv.Itoa(event.floor) + " " + eventTypes[event.eventType]
}

func CheckEvents(event chan Event_t) {
	driver.Init() //mulig flyttes til main?

	go checkButtons(event)

	go checkFloorSignal(event)
}

func checkButtons(event chan Event_t) {
	var lastButtonState [driver.N_FLOORS][driver.N_BUTTONS]bool
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
		button = button % driver.N_BUTTONS
		if button == 0 {
			floor++
			floor = floor % driver.N_FLOORS
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
