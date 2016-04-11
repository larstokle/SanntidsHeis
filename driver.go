package driver

/*
#cgo CFLAGS: -std=c11
#include "../simulator_2/client/elev.h"
#include "./channels.h"
*/
import "C"
import (
	. "../constants"
	"time"
)
func Init() {
	C.elev_init(ET_simulation)
	RunStop()
	/*
	for i := 0; i < N_FLOORS; i++{
		for j:= 0; j < N_BUTTON_TYPES; j++{
			SetButtonLight(j,i,false)
		}
	}
	
	return returnVal
	*/
}
func RunUp() {
	/*
	C.io_clear_bit(C.MOTORDIR)
	//time.Sleep(time.Second * 1)
	C.io_write_analog(C.MOTOR, 2800)
	*/
	C.elev_set_motor_direction(C.DIRN_UP)
    
}

func RunDown() {
	/*
	C.io_set_bit(C.MOTORDIR)
	//time.Sleep(time.Second * 1)
	C.io_write_analog(C.MOTOR, 2800)
	*/
	C.elev_set_motor_direction(C.DIRN_DOWN)
}

func RunStop() {
	//C.io_write_analog(C.MOTOR, 0)
	C.elev_set_motor_direction(C.DIRN_DOWN)

}

func SetButtonLight(button int, floor int, value bool) {
	//channel := C.int(encodeLight(button, floor))
	if value {
		C.elev_set_button_lamp(C.int(button),C.int(floor),1)
		//C.io_set_bit(channel)
	} else {
		C.elev_set_button_lamp(C.int(button),C.int(floor),0)
		//C.io_clear_bit(channel)
	}
}

/*
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
*/

func SetFloorIndicator(floor int) {
	C.elev_set_floor_indicator( C.int(floor) );
	/*
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
	*/
}

func SetDoorOpen(value bool) {
	C.elev_set_door_open_lamp(C.int(value))
	/*
	if value {
		C.io_set_bit(C.LIGHT_DOOR_OPEN)
	} else {
		C.io_clear_bit(C.LIGHT_DOOR_OPEN)
	}
	*/
}

func ReadButton(button int, floor int) bool {
	return int( elev_get_button_signal(C.elev_button_type_t(button) , C.int(floor) ) )
	/*
	if floor < 0 || floor >= N_FLOORS || button < 0 || button > N_BUTTON_TYPES {
		return false
	}

	var BTN_CHANNELS = [N_FLOORS][N_BUTTON_TYPES]int{
		{C.BUTTON_UP1, C.BUTTON_DOWN1, C.BUTTON_COMMAND1},
		{C.BUTTON_UP2, C.BUTTON_DOWN2, C.BUTTON_COMMAND2},
		{C.BUTTON_UP3, C.BUTTON_DOWN3, C.BUTTON_COMMAND3},
		{C.BUTTON_UP4, C.BUTTON_DOWN4, C.BUTTON_COMMAND4}}

	return (int(C.io_read_bit(C.int(BTN_CHANNELS[floor][button]))) != 0)
	*/
}


func GetFloorSignal() int {
	return int(elev_get_floor_sensor_signal(void) )
	/*
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
	*/
}

