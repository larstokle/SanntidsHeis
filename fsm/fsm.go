package fsm

type State int

const (
	STATE_IDLE State = iota
	STATE_MOVING
	STATE_DOOR_OPEN
)

var states = [...]string{
	"IDLE",
	"STATE_MOVING",
	"STATE_DOOR_OPEN",
}

func (state State) String() string {
	return states[state]
}

var fsmState = STATE_IDLE

func GetState() State {
	return fsmState
}
