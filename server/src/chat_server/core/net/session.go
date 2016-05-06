package net

import (
	"chat_server/core/debug"
	"chat_server/core/log"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrSendToClosedSession error = errors.New("Send to closed session!")
	ErrSessionBlocked      error = errors.New("Session blocked!")
)

const (
	MaxResponseNum   int = 1000
	MaxBufferSize    int = 1024
	PacketHeadLength int = 4
)

type requestHandler interface {
	Handler(*Session, []byte)
}

type requestHandlerStruct struct {
	callback func(*Session, []byte)
}

func (Handler requestHandlerStruct) Handler(session *Session, msg []byte) {
	Handler.callback(session, msg)
}

type Response interface {
	ByteSize() int
	Encode(*Buffer)
}
type Session struct {
	Id                int32
	rMutex            sync.Mutex
	wMutex            sync.Mutex
	conn              Conn
	sendBuffSize      int32
	closeChan         chan int
	readChan          chan int
	writeDataChan     chan []byte
	writePacketChan   chan []byte
	writeResponseChan chan Response
	requestHandler    requestHandler
	closeCallback     func(*Session)
	buffer            *Buffer
	isStart           bool
	closeFlag         int32
	*setting
}

func NewSession(id int32, conn Conn, setting *setting, sendBuffSize int32, requestHandler requestHandler, closeCallback func(*Session)) *Session {
	return &Session{
		Id:                id,
		conn:              conn,
		sendBuffSize:      sendBuffSize,
		closeChan:         make(chan int),
		readChan:          make(chan int, setting.rMaxSize),
		writeDataChan:     make(chan []byte, setting.wMaxSize),
		writePacketChan:   make(chan []byte, setting.wMaxSize),
		writeResponseChan: make(chan Response, MaxResponseNum),
		closeCallback:     closeCallback,
		requestHandler:    requestHandler,
		buffer:            NewBuffer(make([]byte, MaxBufferSize)),
		setting:           setting,
	}
}

func (this *Session) Start() {
	if !this.isStart {
		this.isStart = true
		go this.readLoop()
		go this.writeLoop()
	}
}

func (this *Session) readLoop() {
	this.rMutex.Lock()
	defer this.rMutex.Unlock()

	defer func() {
		if err := recover(); err != nil {
			log.Errorf(`
Error = %v
Stack =
%s\n`,
				err,
				debug.Stack(1, "    "),
			)
		}
	}()
	defer func() {
		close(this.readChan)
	}()

L:
	for {
		select {
		case <-this.closeChan:
			break L
		default:
			msg := this.conn.Read()
			this.requestHandler.Handler(this, msg)
		}
	}
}

func (this *Session) writeLoop() {
	this.wMutex.Lock()
	defer this.wMutex.Unlock()
	defer func() {
		if err := recover(); err != nil {
			log.Errorf(`
Stack:
%s\n
	    `, debug.Stack(1, " "))
		}
	}()
	defer func() {
		close(this.writeDataChan)
		close(this.writePacketChan)
		close(this.writeResponseChan)
	}()
L:
	for {
		select {
		case data := <-this.writeDataChan:
			this.conn.WriteData(data)

		case packet := <-this.writePacketChan:
			this.conn.WritePacket(packet)

		case response := <-this.writeResponseChan:
			this.WriteResponse(response)

		case <-this.closeChan:
			break L
		}
	}
}

func (this *Session) WriteResponse(response Response) {
	buffer := this.buffer

	if buffer.Cap() < response.ByteSize()+PacketHeadLength {
		buffer.Set(make([]byte, response.ByteSize()+PacketHeadLength+512))
	}
	buffer.SetWritePos(PacketHeadLength)
	response.Encode(buffer)

	data := buffer.Get()
	littleEndian.PutUint32(data, uint32(len(data)-PacketHeadLength))

	this.conn.WritePacket(data)
}

func (this *Session) Close() {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf(`
Stack:
%s\n
	    `, debug.Stack(1, " "))
		}
	}()
	if atomic.CompareAndSwapInt32(&this.closeFlag, 0, 1) {
		close(this.closeChan)

		<-this.readChan
		<-this.writeDataChan
		<-this.writePacketChan

		this.isStart = false
		if this.closeCallback != nil {
			this.closeCallback(this)
		}
	}
}

func (this *Session) SendResponse(response Response) {
	if atomic.LoadInt32(&this.closeFlag) == 1 {
		panic(ErrSendToClosedSession)
	}

	select {
	case this.writeResponseChan <- response:

	default:
		this.Close()
		panic(ErrSessionBlocked)
	}
}
func (this *Session) SendData(data []byte) {
	if atomic.LoadInt32(&this.closeFlag) == 1 {
		panic(ErrSendToClosedSession)
	}

	select {
	case this.writeDataChan <- data:

	default:
		this.Close()
		panic(ErrSessionBlocked)
	}
}

func (this *Session) SendPacket(packet []byte) {
	if atomic.LoadInt32(&this.closeFlag) == 1 {
		panic(ErrSendToClosedSession)
	}

	select {
	case this.writePacketChan <- packet:

	default:
		this.Close()
		panic(ErrSessionBlocked)
	}
}
