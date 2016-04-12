package message

import(
	."globals"
	"time"
	"fmt"
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

var messageTypes = [...]string{
	"HEARTBEAT",
	"NEW_ORDER",
	"REMOVE_ORDER",
	"REQUEST_ORDER",
	"DELEGATE_ORDER",
	"COST",
	"SYNC",
}

type Message_t struct {
	Source int
	ElevatorId int
	MessageId int
	Button Button_t
	Cost int
	Time time.Time
	Data []byte
}

func (msg Message_t) String() string{
	if msg.MessageId < len(messageTypes){
		return fmt.Sprintf("Message_t{Source: %d, ElevatorId: %d, MessageId: %s, Button: Button_t{%+v}, Cost: %d", msg.Source , msg.ElevatorId ,messageTypes[msg.MessageId] ,msg.Button,msg.Cost)
	} else {
		return fmt.Sprintf("Message_t{Source: %d, ElevatorId: %d, MessageId: unknown(%d), Button: Button_t{%+v}, Cost: %d", msg.Source , msg.ElevatorId , msg.MessageId,msg.Button,msg.Cost)
	}
}

