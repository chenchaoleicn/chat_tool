package handleMessage

import (
	"chat_server/api/protocol/api_handleMessage"
	"chat_server/core/net"
)

func ReturnMessage(session *net.Session, in *api_handleMessage.ReturnMessage_in) *api_handleMessage.ReturnMessage_out {
	out := &api_handleMessage.ReturnMessage_out{}
	out.Message = make([]byte, len(in.Message))
	copy(out.Message, in.Message)
	return out
}
