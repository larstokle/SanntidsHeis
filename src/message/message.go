package message

import(
	."globals"
	"time"
)

const(
	HEARTBEAT = iota
	NEW_ORDER
	WANTS_ORDER
	COST
)

type Message_t struct {
	Source int
	MessageId int
	Button Button_t
	Time time.Time
	Data []byte
}