package UDP

func Parse(data []byte) interface{}{
	//fmt.Printf("unpack Recieved: %+v\n", data)
	newRawData := make(map[string]interface{})
	json.Unmarshal(data, &newRawData)
	for k, _:= range newRawData{
		switch k {
		case "Event_t":
			temp := make(map[string]Event_t)
			json.Unmarshal(data, &temp)
			return temp[k]
		case "int":
			temp := make(map[string]int)
			json.Unmarshal(data, &temp)
			return temp[k]
		}
	}
	return nil
}


func Pack(data interface{}) []byte{
	newMsg := make(map[string]interface{})
	newMsg[reflect.TypeOf(data).Name()] = data
	//newMsg[newMsgType] = data
	b,_ := json.Marshal(newMsg)
	//fmt.Printf("%s \n",b)
	return b
}