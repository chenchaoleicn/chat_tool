package net

import (
	"net"
	"sync"
	"time"
)

var (
	initialReadTimeOut         = 10 * time.Second
	initialWriteTimeOut        = 10 * time.Second
	initialMaxReadSize  uint32 = 1024
	initialMaxWriteSzie uint32 = 1024
)

func Listen(netType, laddr string) Listener {
	listener, err := net.Listen(netType, laddr)
	if err != nil {
		panic(err)
	}
	return NewListener(listener)
}

type Listener struct {
	mutex    sync.Mutex
	listener net.Listener
	setting
}

func NewListener(listener net.Listener) Listener {
	return Listener{
		listener: listener,
		setting: setting{
			wTimeOut: initialWriteTimeOut,
			rTimeOut: initialReadTimeOut,
			wMaxSize: initialMaxWriteSzie,
			rMaxSize: initialMaxReadSize,
		},
	}
}

func (this *Listener) GetLocalAddr() net.Addr {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.listener.Addr()
}

func (this *Listener) Accept() Conn {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	conn, err := this.listener.Accept()
	if err != nil {
		panic(err)
	}
	this.setConnectionAttri(conn)
	Conn := NewConn(conn, &this.setting)
	return Conn
}

func (this *Listener) Close() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if err := this.listener.Close(); err != nil {
		panic(err)
	}
}
func (this *Listener) setConnectionAttri(conn net.Conn) {
	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetKeepAlive(true)
	tcpConn.SetNoDelay(true)
	tcpConn.SetKeepAlivePeriod(30 * time.Second)
}
