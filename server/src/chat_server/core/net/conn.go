package net

import (
	"errors"
	"io"
	"net"
	"sync"
	"time"
)

var ErrDataTooLang error = errors.New("The data  is too lang!")

type Request interface {
	Process(*Session)
	Decode(*Buffer)
}

type Conn struct {
	mutex sync.Mutex
	conn  net.Conn

	rHeadBuff []byte
	wHeadBuff []byte
	rMutex    sync.Mutex
	wMutex    sync.Mutex

	*setting
}

func NewConn(conn net.Conn, setting *setting) Conn {
	return Conn{
		conn:      conn,
		rHeadBuff: make([]byte, PacketHeadLength),
		wHeadBuff: make([]byte, PacketHeadLength),
		setting:   setting,
	}
}

func (conn *Conn) LocalAddr() net.Addr {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	return conn.conn.LocalAddr()
}

func (conn *Conn) RemoteAddr() net.Addr {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	return conn.conn.RemoteAddr()
}

func (conn *Conn) getNewSpace(length uint32) (data []byte) {
	data = make([]byte, length)
	return
}

func (conn *Conn) Read() []byte {
	conn.rMutex.Lock()
	defer conn.rMutex.Unlock()

	if conn.rTimeOut != 0 {
		conn.conn.SetReadDeadline(time.Now().Add(conn.rTimeOut))
	} else {
		conn.conn.SetReadDeadline(time.Time{})
	}

	if _, err := io.ReadFull(conn.conn, conn.rHeadBuff); err != nil {
		panic(err)
	}
	size := littleEndian.Uint32(conn.rHeadBuff)
	if conn.rMaxSize > 0 && size > conn.rMaxSize {
		panic(ErrDataTooLang)
	}

	receiveData := conn.getNewSpace(size)
	if _, err := io.ReadFull(conn.conn, receiveData); err != nil {
		panic(err)
	}
	return receiveData
}

//write data,not an entire packet, no packet head.
func (conn *Conn) WriteData(data []byte) {
	conn.wMutex.Lock()
	defer conn.wMutex.Unlock()

	if conn.wMaxSize > 0 && uint32(len(data)) > conn.wMaxSize {
		panic(ErrDataTooLang)
	}

	if conn.wTimeOut != 0 {
		conn.conn.SetWriteDeadline(time.Now().Add(conn.wTimeOut))
	} else {
		conn.conn.SetWriteDeadline(time.Time{})
	}

	littleEndian.PutUint32(conn.wHeadBuff, uint32(len(data)))
	if _, err := conn.conn.Write(conn.wHeadBuff); err != nil {
		panic(err)
	}

	if _, err := conn.conn.Write(data); err != nil {
		panic(err)
	}
}

//write packet, an entire packet
func (conn *Conn) WritePacket(packet []byte) {
	conn.wMutex.Lock()
	defer conn.wMutex.Unlock()

	if conn.wMaxSize > 0 && uint32(len(packet)) > conn.wMaxSize {
		panic(ErrDataTooLang)
	}

	if conn.wTimeOut != 0 {
		conn.conn.SetWriteDeadline(time.Now().Add(conn.wTimeOut))
	} else {
		conn.conn.SetWriteDeadline(time.Time{})
	}
	if _, err := conn.conn.Write(packet); err != nil {
		panic(err)
	}
}
