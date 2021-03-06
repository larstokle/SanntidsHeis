package driver

/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: ${SRCDIR}/simulator/simelev.a /usr/lib/x86_64-linux-gnu/libphobos2.a -lpthread -lcomedi -lm
#include "./io.h"
#include "./channels.h"
*/
import "C"
import (
	. "globals"
	"fmt"
)

func init(){
	if int(C.io_init(ET_comedi)) != 1{
		fmt.Printf("ERROR! driver: init failed")
	}
	RunStop()

	for i := 0; i < N_FLOORS; i++{
		for j:= 0; j < N_BUTTON_TYPES; j++{
			SetButtonLight(j,i,false)
		}
	}
}

func ReadButton(button int, floor int) bool {
	if floor < 0 || floor >= N_FLOORS || button < 0 || button > N_BUTTON_TYPES {
		return false
	}

	var BTN_CHANNELS = [N_FLOORS][N_BUTTON_TYPES]int{
		{C.BUTTON_UP1, C.BUTTON_DOWN1, C.BUTTON_COMMAND1},
		{C.BUTTON_UP2, C.BUTTON_DOWN2, C.BUTTON_COMMAND2},
		{C.BUTTON_UP3, C.BUTTON_DOWN3, C.BUTTON_COMMAND3},
		{C.BUTTON_UP4, C.BUTTON_DOWN4, C.BUTTON_COMMAND4}}

	return (int(C.io_read_bit(C.int(BTN_CHANNELS[floor][button]))) != 0)
}

func SetButtonLight(button int, floor int, value bool) {
	channel := C.int(encodeLight(button, floor))
	if value {
		C.io_set_bit(channel)
	} else {
		C.io_clear_bit(channel)
	}
}

func encodeLight(button int, floor int) int {

	channel := C.LIGHT_COMMAND1
	if button == CMD {
		channel = channel - floor
	} else if button == UP && floor == 0 {
		channel = C.LIGHT_UP1
	} else if button == DOWN && floor == 3 {
		channel = C.LIGHT_DOWN4
	} else {
		channel = C.LIGHT_UP2
		channel = channel - button - 2*(floor-1)
	}
	return channel
}



func RunUp() {
	C.io_clear_bit(C.MOTORDIR)
	//time.Sleep(time.Second * 1)
	C.io_write_analog(C.MOTOR, 2800)
}

func RunDown() {
	C.io_set_bit(C.MOTORDIR)
	//time.Sleep(time.Second * 1)
	C.io_write_analog(C.MOTOR, 2800)
}

func RunStop() {
	C.io_write_analog(C.MOTOR, 0)
}

func SetFloorIndicator(floor int) bool {
	if floor < 0 || floor >= N_FLOORS {
		return false
	}

	if (floor & 0x02) != 0 {
		C.io_set_bit(C.LIGHT_FLOOR_IND1)
	} else {
		C.io_clear_bit(C.LIGHT_FLOOR_IND1)
	}

	if (floor & 0x01) != 0 {
		C.io_set_bit(C.LIGHT_FLOOR_IND2)
	} else {
		C.io_clear_bit(C.LIGHT_FLOOR_IND2)
	}

	return true
}

func GetFloorSignal() int {
	if C.io_read_bit(C.SENSOR_FLOOR1) != 0 {
		return 0
	} else if C.io_read_bit(C.SENSOR_FLOOR2) != 0 {
		return 1
	} else if C.io_read_bit(C.SENSOR_FLOOR3) != 0 {
		return 2
	} else if C.io_read_bit(C.SENSOR_FLOOR4) != 0 {
		return 3
	} else {
		return -1
	}
}

func SetDoorOpen(value bool) {
	if value {
		C.io_set_bit(C.LIGHT_DOOR_OPEN)
	} else {
		C.io_clear_bit(C.LIGHT_DOOR_OPEN)
	}
}
