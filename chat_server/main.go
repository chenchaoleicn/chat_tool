package main

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jack301/chat_tool/chat_server/api"
	"github.com/Jack301/chat_tool/chat_server/core/debug"
	"github.com/Jack301/chat_tool/chat_server/core/log"
	"github.com/Jack301/chat_tool/chat_server/core/net"
	_ "github.com/Jack301/chat_tool/chat_server/module/handleMessage"
	_ "github.com/Jack301/chat_tool/mysql"
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
