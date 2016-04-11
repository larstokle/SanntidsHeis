package transactionMgr

import( 
	"encoding/json"
	"reflect"
	//."../orderque"
	"../eventmgr"
	"fmt"
	//
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
		case "Hartbeat_t":
			temp := make(map[string]Hartbeat_t)
			json.Unmarshal(data, &temp)
			//fmt.Printf("parsed a hartbeat_t: %+v \n", temp[k])
			return temp[k]
		}
	}
	return nil
}


func Pack(data interface{}) []byte{
	newMsg := make(map[string]interface{})
	newMsg[reflect.TypeOf(data).Name()] = data
	b, err := json.Marshal(newMsg)
	//fmt.Printf("Packed %+v as %s \n",newMsg, b)
	if !checkAndPrintError(err, "Marshal error"){
		return b
	}
	return nil
}


func checkAndPrintError(err error, info string) bool {
	if err != nil {
		switch e := err.(type){
		default:
			fmt.Println(info, ": ", e)
		}
		return true
	}
	return false
}