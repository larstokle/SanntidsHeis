package message

import(
	."globals"
	"time"
)

const(
	HEARTBEAT = iota
	NEW_ORDER
	REMOVE_ORDER
	REQUEST_ORDER
	DELEGATE_ORDER
	COST
	SYNC
)

type Message_t struct {
	Source int
	ElevatorId int
	MessageId int
	Button Button_t
	Cost int
	Time time.Time
	Data []byte
}


