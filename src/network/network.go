package network

import (
	"encoding/json"
	"fmt"
	. "message"
	"net"
	"strconv"
	"time"
)

//close channel to exit
func MakeSender(addr string) chan<- Message_t {
	msg := make(chan Message_t)

	go func() {
		toAddr, err := net.ResolveUDPAddr("udp", addr)
		checkAndPrintError(err, "ERROR! ResolveUDPAddr")

		waitForNetworkAvailability()
		conn, err := net.DialUDP("udp", nil, toAddr)
		checkAndPrintError(err, "ERROR! DialUDP")

		defer conn.Close()
		for newMsg := range msg {
			json_msg, err := json.Marshal(newMsg)
			checkAndPrintError(err, "ERROR! Marshal")
			_, err = conn.Write(json_msg)
			if checkAndPrintError(err, "ERROR! WriteToUDP") {
				waitForNetworkAvailability()
			}
		}

		fmt.Printf("Sender: Channel closed. function returning\n")
	}()
	return msg
}

//Send to channel for exit.
func MakeReceiver(port string) chan Message_t {
	msg := make(chan Message_t)
	go func() {
		localAddr, err := net.ResolveUDPAddr("udp", port)
		checkAndPrintError(err, "Resolve UDP error")

		conn, err := net.ListenUDP("udp", localAddr)
		if err != nil {
			fmt.Printf("ERROR: ListenUDP error: %s\n", err)
		}
		defer conn.Close()

		for {
			buf := make([]byte, 2048)
			conn.SetReadDeadline(time.Now().Add(time.Millisecond * 2000))
			n, _, err := conn.ReadFromUDP(buf)
			if !checkAndPrintError(err, "ERROR! ReadFromUDP:") {

				var recived Message_t
				err = json.Unmarshal(buf[0:n], &recived)
				checkAndPrintError(err, "ERROR! Unmarshal")

				select {
				case <-msg:
					fmt.Printf("Reciever: Received on send channel. function returning\n")
					return
				case msg <- recived:
				}
			}
			select {
			case <-msg:
				fmt.Printf("Reciever: Received on send channel. function returning\n")
				return
			default:
				continue
			}
		}
	}()
	return msg
}

func waitForNetworkAvailability() {
	if GetLocalIP()[0:2] == "::" {
		fmt.Println("Sender: No network available. Wait for connection to establish")
		for GetLocalIP()[0:2] == "::" {
			time.Sleep(time.Second * 2)
		}
		fmt.Println("Sender: Network found. Now proceeding for connection")
	}
}

func checkAndPrintError(err error, info string) bool {
	if err != nil {
		switch e := err.(type) {
		case net.Error:
			if !e.Timeout() {
				fmt.Println(info, ": ", err)
			}
		default:
			fmt.Println(info, ": ", err)
		}
		return true
	}
	return false
}

func GetLocalIP() string {
	addr, _ := net.InterfaceAddrs()
	return addr[1].String()
}

func GetLastIPByte() int {
	addr := GetLocalIP()
	if addr[0:2] == "::" {
		return -1
	}
	dot := 0
	backslash := 0
	for i, ch := range addr {
		if string(ch) == "." {
			dot = i + 1
		}
		if string(ch) == "/" {
			backslash = i
			break
		}
	}

	lastByte := addr[dot:backslash]
	num, err := strconv.Atoi(lastByte)

	if !checkAndPrintError(err, "strconv error in GetLastIPByte") {
		return num
	} else {
		return -1
	}
}
