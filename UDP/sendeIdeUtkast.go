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

func MakeSender(addr string, msg chan []byte, quit chan bool) {

	toAddr, err := net.ResolveUDPAddr("udp", addr)
	CheckAndPrintError(err, "ResolveUDP error")

	conn, err := net.DialUDP("udp", nil, toAddr)
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
				//fmt.Printf("Sender sending %+v \n", newMsg)
				_, err := conn.Write(newMsg)
				CheckAndPrintError(err, "Writing error")
			}
		}
	}()
}


func MakeReciever(port string, message chan []byte, quit chan bool) {

	localAddr, err := net.ResolveUDPAddr("udp", port)

	CheckAndPrintError(err, "Resolve UDP error")

	conn, err := net.ListenUDP("udp", localAddr)

	CheckAndPrintError(err, "ListenUDP error")

	
	

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
				buf := make([]byte, 1024)
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 2000))
				n, _, err := conn.ReadFromUDP(buf)
				CheckAndPrintError(err, "ReadFromUDP error")
				if err == nil {
					//fmt.Printf("Reciever recieved: %+v of size: %d\n",buf[0:n], n)
					message <- buf[0:n]
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

func GetLocalIP() string{
	addr, _ := net.InterfaceAddrs()
	return addr[1].String()	
}


