package UDP

import( 
	"encoding/json"
	"reflect"
	//."../orderque"
	"../eventmgr"
	//"fmt"

)


func Parse(data []byte) interface{}{
	//fmt.Printf("unpack Recieved: %+v\n", data)
	newRawData := make(map[string]interface{})
	json.Unmarshal(data, &newRawData)
	for k, _:= range newRawData{
		switch k {
		case "Event_t":
			temp := make(map[string]eventmgr.Event_t)
			err := json.Unmarshal(data, &temp)
			if !checkAndPrintError(err, "Unmarshal error"){
				//fmt.Printf("parsed an Event_t to: %+vwhich is of type %+v\n", temp[k], reflect.TypeOf(temp[k]))
				return temp[k]
			}
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
	b, err := json.Marshal(newMsg)
	//fmt.Printf("%s \n",b)
	if !checkAndPrintError(err, "Marshal error"){
		return b
	}
	return nil
}