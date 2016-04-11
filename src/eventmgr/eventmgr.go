package eventmgr

import (
	. "globals"
	"driver"
	"time"
	"fmt"
)


func CheckButtons() chan Button_t{
	event := make(chan Button_t, 100)
	var lastButtonState [N_FLOORS][N_BUTTON_TYPES]bool
	floor, button := 0, 0

	go func(){
		for true {
			pressed := driver.ReadButton(button, floor)
			if pressed != lastButtonState[floor][button] {
				lastButtonState[floor][button] = pressed
				if pressed {
					fmt.Println("button pressed!")
					var newEvent Button_t //tungvint! lag oneliner!!
					newEvent.Floor = floor
					newEvent.ButtonType = button
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
