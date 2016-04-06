
type msg struct{
	varType string// prøv med: type
	data interface{}
}

send(data chan interface{}){
	newData <- data
	var newMsg msg
	newMsg.data = data
	// newMsg.varType = newData.(type) //evt
	switch newData.(type){
		case orederQue_t:
			newMsg.varType = "que"
			b := json.marshal(newmsg)

	}
	broadcast(b)
}