package net

import (
	"fmt"
	"net"
	"testing"

	"github.com/Jack301/chat_tool/chat_server/core/log"
)

func Test_listen(t *testing.T) {
	log.AddTestFlag("Test_listen")
	defer log.AddTestFlag("")

	var g_server server
	g_server = Server("tcp", "127.0.0.1:8080")
	g_server.Start()
	addr := g_server.GetLocalListenAddr()
	setting := &g_server.listener.setting

	clientConn, err := net.Dial(addr.Network(), addr.String())
	conn := NewConn(clientConn)
	conn.setting = setting
	conn.WriteData("你好！")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connect success!")
	}

	fmt.Printf("%#v\n", g_server.sessions)
}

/*
func Benchmark_listen(b *testing.B) {
	//	g_server := Server("tcp", "127.0.0.1:8080")
	//	g_server.Start()
	addr := g_server.GetLocalListenAddr()
	setting := &g_server.listener.setting
	for i := 0; i < 10; i++ {
		go func() {
			clientConn, err := net.Dial(addr.Network(), addr.String())
			conn := NewConn(clientConn)
			conn.setting = setting
			conn.WriteData("你好！")
			if err != nil {
				panic(err)
			} else {
				fmt.Println("Connect success!")
			}
		}()
	}
}
*/
