package network

import(
	"time"
	"fmt"
	"net"
	"strconv"
)


func MakeSender(addr string, msg chan []byte, quit chan bool) {

	toAddr, err := net.ResolveUDPAddr("udp", addr)
	checkAndPrintError(err, "ResolveUDP error")

	conn, err := net.DialUDP("udp", nil, toAddr)
	checkAndPrintError(err, "DialUDP error")

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
				checkAndPrintError(err, "WriteToUDP error")
			}
		}
	}()
}


func MakeReceiver(port string, message chan []byte, quit chan bool) {

	localAddr, err := net.ResolveUDPAddr("udp", port)
	checkAndPrintError(err, "Resolve UDP error")

	conn, err := net.ListenUDP("udp", localAddr)
	checkAndPrintError(err, "ListenUDP error")
	
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
				if !checkAndPrintError(err, "ReadFromUDP error"){
					//fmt.Printf("Reciever recieved: %+v of size: %d\n",buf[0:n], n)
					message <- buf[0:n]
				}
			}
		}
	}()
}


func checkAndPrintError(err error, info string) bool {
	if err != nil {
		switch e := err.(type){
		case net.Error:
			if !e.Timeout(){
				fmt.Println(info, ": ", err)
			}
		default:
			fmt.Println(info, ": ", err)
		}
		return true
	}
	return false
}

func GetLocalIP() string{
	addr, _ := net.InterfaceAddrs()
	return addr[1].String()	
}


func GetLastIPByte() int{
	addr := GetLocalIP()
	lastByte := addr[12:15]
	i,err := strconv.Atoi(lastByte)
	
	if !checkAndPrintError(err, "strconv error in GetLastIPByte") {
		return i
	} else {
		return -1
	}
}

