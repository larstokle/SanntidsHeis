package UDP

import(
	//"net"
	"encoding/json"
	.	"../constants"
	"time"
	"fmt"
	"net"
	"reflect"
)

type OrderQue_t [N_FLOORS][N_ORDER_TYPES]struct {
	hasOrder       bool
	lastChangeTime time.Time
	//assignedToID int //kanskje unødvendig? fjerner den encapsulation?
}

type Order_t struct {
	floor     int
	orderType int
}

type Event_t struct {
	Floor     int
	EventType int
}
type msg struct{
	varType string// prøv med: type
	data interface{}
}

func MakeSender(addr string, msg chan interface{}, quit chan bool) {

	toAddr, err := net.ResolveUDPAddr("udp", addr)
	CheckAndPrintError(err, "ResolveUDP error")

	conn, err := net.DialUDP("udp", nil, toAddr)
	//defer conn.Close()
	CheckAndPrintError(err, "DialUDP error")

	go func() {
		defer conn.Close()
		for {
			select {
			case q := <-quit:
				if q {
					defer func() { quit <- false }()
					defer fmt.Println("Quiting Sender")
					return
				}
			case newMsg := <-msg:
				//fmt.Println("newMsg",pack(newMsg))
				_, err := conn.Write(pack(newMsg))
				CheckAndPrintError(err, "Writing error")
			}
		}
	}()
}


func pack(data interface{})[]byte{
	var b []byte
	// newMsgType.varType = newData.(type) //evt
	/*switch data.(type){
		case OrderQue_t:
			newMsgType = "que"
		case Event_t:
			newMsgType = "Event_t"
		case int:
			newMsgType = "int"
		default:
			fmt.Println("Fuck!")

	}*/

	newMsg := make(map[string]interface{})
	newMsg[reflect.TypeOf(data).Name()] = data
	//newMsg[newMsgType] = data
	b,_ = json.Marshal(newMsg)
	//fmt.Printf("%s \n",b)
	return b
}

func MakeReciever(port string, message chan interface{}, quit chan bool) {

	localAddr, err := net.ResolveUDPAddr("udp", port)

	CheckAndPrintError(err, "Resolve UDP error")

	conn, err := net.ListenUDP("udp", localAddr)

	CheckAndPrintError(err, "ListenUDP error")

	buf := make([]byte, 1024)
	something := make(map[string]interface{})

	go func() {
		defer conn.Close()

		for {
			select {
			case q := <-quit:
				if q {
					defer func() { quit <- false }()
					defer fmt.Println("Quiting Reciever")
					return
				}
			default:
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 2000))
				n, _, err := conn.ReadFromUDP(buf)
				CheckAndPrintError(err, "ReadFromUDP error")
				if err == nil {

					json.Unmarshal(buf[0:n], &something)
					for k, v := range something{
						switch k {
						case "Event_t":
							fmt.Println("found Event_t of type",reflect.TypeOf(v))
							temp := make(map[string]Event_t)
							json.Unmarshal(buf[0:n], &temp)
							message <- temp[k]
							//vv := v.(map[string]interface{})
							//message <- Event_t{Floor: int(vv["Floor"].(float64)), EventType: int(vv["EventType"].(float64)) }
						case "int":
							fmt.Println("found int",reflect.TypeOf(v))
							message <- int(v.(float64))
						}
					}

					//message <- string(buf[0:n])
					//message <- something
				}

			}
		}
	}()
}

func CheckAndPrintError(err error, info string) {
	if err != nil && !err.(net.Error).Timeout() {
		fmt.Println(info, ": ", err)
		//exit(1) maybe??
	}
}


