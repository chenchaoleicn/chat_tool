package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	message      string = "Hello World!"
	receiveData  string
	sendTimes    int = 0
	receiveTimes int = 0
	littleEndian     = binary.LittleEndian
	signalExit       = make(chan bool)
)

func main() {
	clientConn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	go signalHandler()
	tickTime := time.Tick(1 * time.Second)
	for {
		select {
		case <-tickTime:
			select {
			case <-signalExit:
				clientConn.Close()
				return
			default:
			}
			packet := EncodePacket([]byte(message))
			if _, err = clientConn.Write(packet); err != nil {
				panic(err)
			}
			receiveData = ReadPacket(clientConn)
			if receiveData != "" {
				receiveTimes++
			}
			sendTimes++
			fmt.Printf("%v times send: %v---%v times receive:%v\n", sendTimes, message, receiveTimes, receiveData)
		}
	}
}

func ReadPacket(conn net.Conn) string {
	var (
		packetHead = make([]byte, 4)
		packetData []byte
	)
	if _, err := io.ReadFull(conn, packetHead); err != nil {
		panic(err)
	}
	packetBodyLength := littleEndian.Uint32(packetHead)
	packetData = make([]byte, packetBodyLength)
	if _, err := io.ReadFull(conn, packetData); err != nil {
		panic(err)
	}
	return string(packetData[2:])
}

func EncodePacket(msg []byte) []byte {
	packet := make([]byte, 4+1+1+2+len(msg)) //4(包长度)+1(模块id)+1(接口id)+2(包内容，指示后面的msg长度)
	littleEndian.PutUint32(packet, uint32(len(msg)+4))
	packet[4] = byte(1)
	packet[5] = byte(1)
	littleEndian.PutUint16(packet[6:], uint16(len(msg)))
	copy(packet[8:], msg)

	return packet
}

func signalHandler() {
	signalTerm := make(chan os.Signal, 1)

	go signal.Notify(signalTerm, syscall.SIGINT)
	for {
		select {
		case <-signalTerm:
			signalExit <- true
		}
	}
}
