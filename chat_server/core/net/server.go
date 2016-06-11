package net

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"

	"github.com/Jack301/chat_tool/chat_server/core/debug"
	"github.com/Jack301/chat_tool/chat_server/core/log"
)

var (
	MaxSessionId           int32 = 9999
	MaxSendBuffSize        int32 = 1024
	ErrAllSessionIdAreUsed error = errors.New("All sesion id has been uesed!")
)

type server struct {
	mutex            sync.Mutex
	listener         Listener
	maxSessionId     int32
	sendBuffSize     int32
	sessionCloseHook func(*Session)
	sessionStartHook func(*Session)
	requestHandler   requestHandler
	sessions         map[int32]*Session
	isStart          bool
	*setting
}

func Server(network, addr string) server {
	listener := Listen(network, addr)
	return NewServer(listener)
}

func NewServer(listener Listener) server {
	return server{
		listener:     listener,
		sessions:     make(map[int32]*Session, MaxSessionId),
		setting:      &listener.setting,
		sendBuffSize: MaxSendBuffSize,
	}
}

func (this *server) Start() {
	if !this.isStart {
		this.isStart = true
		go this.acceptLoop()
	}
}

func (this *server) acceptLoop() {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf(`
Stack:
%s
		    `, debug.Stack(1, " "))
		}
	}()
	for {
		conn := this.listener.Accept()
		go this.startSession(conn)
	}
}

func (this *server) startSession(conn Conn) {
	session := NewSession(this.NewSessionId(), conn, this.setting, this.sendBuffSize, this.requestHandler, this.sessionCloseHook)
	this.sessions[session.Id] = session
	if this.sessionStartHook != nil {
		this.sessionStartHook(session)
	}
	session.Start()
}

func (this *server) NewSessionId() (newSessionId int32) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	var UsedSessionIdCount int32 = 0
	for newSessionId = (this.maxSessionId + 1) % (MaxSessionId + 1); this.isSessionIdUsed(newSessionId, &UsedSessionIdCount); newSessionId = (this.maxSessionId + 1) % (MaxSessionId + 1) {
	}
	return
}

func (this *server) isSessionIdUsed(sessionId int32, countAdd *int32) bool {
	_, isUsed := this.sessions[sessionId]
	if isUsed {
		atomic.AddInt32(countAdd, 1)
		if atomic.LoadInt32(countAdd) >= MaxSessionId {
			panic(ErrAllSessionIdAreUsed)
		}
	}
	return isUsed
}

func (this *server) Stop() {
	if this.isStart {
		this.isStart = false
	}
	this.listener.Close()
	wg := new(sync.WaitGroup)
	this.CloseSessions(wg)
}

func (this *server) CloseSessions(wg *sync.WaitGroup) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	defer func() {
		if err := recover(); err != nil {
			log.Errorf(`
Stack:
%s\n
	    `, debug.Stack(1, " "))
		}
	}()
	originalSessionCloseHook := this.sessionCloseHook
	this.sessionCloseHook = func(session *Session) {
		originalSessionCloseHook(session)
		wg.Done()
	}
	for _, session := range this.sessions {
		wg.Add(1)
		session.Close()
	}
}

func (this *server) BroadCast(response Response) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	for _, session := range this.sessions {
		session.SendResponse(response)
	}
}
func (this *server) DelSession(sessionId int32) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	session := this.sessions[sessionId]
	if this.sessionCloseHook != nil {
		this.sessionCloseHook(session)
	}
	delete(this.sessions, sessionId)
}

func (this *server) SetSessionStartHook(callback func(*Session)) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.sessionStartHook = callback
}

func (this *server) SetSessionCloseHook(callback func(*Session)) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.sessionCloseHook = callback
}

func (this *server) GetLocalListenAddr() net.Addr {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.listener.GetLocalAddr()
}
func (this *server) OnSessionStart(callback func(session *Session, msg []byte)) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.requestHandler = requestHandlerStruct{callback}
}
