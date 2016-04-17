package orderque

import(
	."globals"
	"encoding/json"
	"os"
	"driver"
	"fmt"
	"time"
)

func (thisQue *orderQue_t) SyncExternal(queToSync orderQue_t) { //add error returns?
	externalButtons := []int{UP, DOWN}
	for floor := FIRST_FLOOR; floor < N_FLOORS; floor++ {
		for _, orderType := range externalButtons {
			if queToSync[floor][orderType].lastChangeTime.After(thisQue[floor][orderType].lastChangeTime) {
				thisQue[floor][orderType] = queToSync[floor][orderType]
				driver.SetButtonLight(orderType, floor, thisQue[floor][orderType].hasOrder)
			}
		}
	}
}

func (thisQue *orderQue_t) SyncInternal(queToSync orderQue_t) { //add error returns?
	for floor := FIRST_FLOOR; floor < N_FLOORS; floor++ {
		if queToSync[floor][CMD].lastChangeTime.After(thisQue[floor][CMD].lastChangeTime) {
			thisQue[floor][CMD] = queToSync[floor][CMD]
			driver.SetButtonLight(CMD, floor, thisQue[floor][CMD].hasOrder)
		}
	}
}


func (que *orderQue_t) WriteToLog(){
	f, err := os.Create(QUE_LOG_FILE)
	check(err)
	defer f.Close()
	err = json.NewEncoder(f).Encode(que)
	check(err)
}

func ReadFromLog() orderQue_t{
	que := New()
	f, err := os.Open(QUE_LOG_FILE)
	defer f.Close()
	if(os.IsNotExist(err)){
		return que
	}

	check(err)
	err = json.NewDecoder(f).Decode(&que)

	return que
}

func check(e error) {
    if e != nil {
        fmt.Printf("ERROR! que logger: ",e)
    }
}

func Encode(que orderQue_t) ([]byte, error){
	return json.Marshal(que)
}

func Decode(data []byte) (orderQue_t, error){
	var que orderQue_t
	err := json.Unmarshal(data, &que)
	return que, err
}


func (order order_t) MarshalJSON() ([]byte, error){
	return json.Marshal(&struct {
		HasOrder       bool		`json:"hasOrder"`
		LastChangeTime time.Time `json:"lastChangeTime"`
		AssignedToID int 	`json:"assignedToID"`
		}{
			HasOrder: order.hasOrder,
			LastChangeTime: order.lastChangeTime,
			AssignedToID: order.assignedToID,
		})
}

func (order *order_t) UnmarshalJSON(data []byte) error {
	temp := struct {
		HasOrder       bool		`json:"hasOrder"`
		LastChangeTime time.Time `json:"lastChangeTime"`
		AssignedToID int 	`json:"assignedToID"`
		}{}
	err := json.Unmarshal(data, &temp)

	*order = order_t{hasOrder: temp.HasOrder, lastChangeTime: temp.LastChangeTime, assignedToID: temp.AssignedToID}
	return err

}