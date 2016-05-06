package main

import (
	"chat_server/api"
	"chat_server/core/debug"
	"chat_server/core/log"
	"chat_server/core/net"
	"database/sql"
	_ "mysql"
	"os"
	"os/signal"
	"syscall"
)

import (
	_ "chat_server/module/handleMessage"
)

var (
	signalExit = make(chan bool)
)

func main() {
	g_server := net.Server("tcp", "127.0.0.1:8080")

	logPrefix := "chat_"
	log.SetLogFilePrefix(logPrefix)
	log.Setup("log", true)

	g_server.OnSessionStart(requestHandlerFunc)
	g_server.Start()
	go signalHandler()

	_, err := sql.Open("mysql", "root:@/test")
	if err != nil {
		panic(err)
	}

	<-signalExit

	g_server.Stop()
	return
}

func signalHandler() {
	signalTerm := make(chan os.Signal, 1)

	go signal.Notify(signalTerm, syscall.SIGKILL)

	for {
		select {
		case <-signalTerm:
			signalExit <- true
		}
	}
}

//TODO
func requestHandlerFunc(session *net.Session, msg []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf(`
Error:%v
Stack:
%s
`, err, debug.Stack(1, "	    "))
		}
	}()
	request := api.DecodeIn(msg)
	request.Process(session)
}
