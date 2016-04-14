package globals

import "strconv"



const (
	N_FLOORS       = 4
	N_BUTTON_TYPES = 3
	N_ORDER_TYPES  = N_BUTTON_TYPES

	FIRST_FLOOR = 0
	TOP_FLOOR   = N_FLOORS - 1

	NONLEGAL_ID = -1



)

type Direction_t int32

const(
	DIR_DOWN = -1
	DIR_STOP = 0
	DIR_UP   = 1
)

type Button_t struct {
	Floor     int
	ButtonType int
}

const(
	UP    = iota
	DOWN         
	CMD          
)


var buttonTypes = [...]string{
	"Up",
	"Down",
	"Command",
}


func (btn Button_t) String() string {
	
	if btn.ButtonType < len(buttonTypes){
		return "{Floor:" + strconv.Itoa(btn.Floor) + ", Type: " + buttonTypes[btn.ButtonType] + "}"
	}else{
		return "{Floor:" + strconv.Itoa(btn.Floor) + ", Type unknown: " + strconv.Itoa(btn.ButtonType) + "}"
	}
}

var NONVALID_BUTTON = Button_t{FIRST_FLOOR, DOWN}

var DEBUG_TRNSMGR bool = true
var DEBUG_ELEVMGR bool = true
var DEBUG_FSM bool = true
var DEBUG_QUE bool = false
var DEBUG_CHANNELS bool = true



