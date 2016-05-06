package handleMessage

import (
	"chat_server/api/protocol/api_handleMessage"
	"chat_server/core/net"
)

func init() {
	api_handleMessage.SetProcessHandler(handleMessage{})
}

type handleMessage struct {
}

func (this handleMessage) ReturnMessage(session *net.Session, in *api_handleMessage.ReturnMessage_in) {
	out := &api_handleMessage.ReturnMessage_out{} //为什么去掉&就报错了
	out = ReturnMessage(session, in)
	session.SendResponse(out)
}
