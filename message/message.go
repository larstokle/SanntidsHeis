package message

const(
	HEARTBEAT = iota
	//BUTTON_PUSHED

)

type Message_t struct {
	Source int
	Message_id int
	data []byte
}