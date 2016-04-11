package globals

const (
	N_FLOORS       = 4
	N_BUTTON_TYPES = 3
	N_ORDER_TYPES  = N_BUTTON_TYPES

	FIRST_FLOOR = 0
	TOP_FLOOR   = N_FLOORS - 1

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
		return "Floor:" + strconv.Itoa(btn.Floor) + ", Type: " + buttonTypes[btn.ButtonType]
	}else{
		return "Floor:" + strconv.Itoa(btn.Floor) + ", Type unknown: " + strconv.Itoa(btn.buttonType)
	}
}