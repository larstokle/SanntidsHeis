package eventmgr

import (
	. "globals"
	"driver"
	"time"
	"fmt"
)


func CheckButtons() <-chan Button_t{
	event := make(chan Button_t, 10)
	var lastButtonState [N_FLOORS][N_BUTTON_TYPES]bool
	floor, button := 0, 0
	fmt.Println("eventmgr: CheckButtons started\n")
	go func(){
		for {
			pressed := driver.ReadButton(button, floor)
			if pressed != lastButtonState[floor][button] {
				lastButtonState[floor][button] = pressed
				if pressed {
					newButtonPressed := Button_t{Floor: floor, ButtonType: button} //tungvint! lag oneliner!!
					if(DEBUG_CHANNELS){fmt.Printf("eventMgr: button %+v pressed! \n", newButtonPressed)}
					event <- newButtonPressed
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

func CheckFloorSignal() <-chan int{
	event := make(chan int)  
	lastFloorState := -1
	fmt.Println("eventmgr: CheckFloorSignal started\n")
	go func(){
		for {
			newFloorState := driver.GetFloorSignal()
			if newFloorState != lastFloorState {
				lastFloorState = newFloorState
				if newFloorState != -1 {
					if(DEBUG_CHANNELS){fmt.Println("eventMgr: newFloorState %d ", newFloorState)}
					event <- newFloorState
				}
			}
			time.Sleep(time.Millisecond * 20)
		}
	}()
	return event
}
