package api_handleMessage

import (
	"chat_server/core/net"
)

type Request interface {
	Process(*net.Session)
	Decode(*net.Buffer)
}

var (
	g_processHandler processHandler
)

type processHandler interface {
	ReturnMessage(session *net.Session, in *ReturnMessage_in)
}

func SetProcessHandler(handler processHandler) {
	g_processHandler = handler
}

type ReturnMessage_in struct {
	Message []byte
}

func (this *ReturnMessage_in) Process(session *net.Session) {
	g_processHandler.ReturnMessage(session, this)
}

func (this *ReturnMessage_in) Decode(buffer *net.Buffer) {
	this.Message = buffer.ReadBytes(buffer.ReadUint16())
}

type ReturnMessage_out struct {
	Message []byte
}

func (this *ReturnMessage_out) ByteSize() int {
	return len(this.Message)
}
func (this *ReturnMessage_out) Encode(buffer *net.Buffer) {
	buffer.WriteUint8(1) //模块Id
	buffer.WriteUint8(1) //接口Id
	buffer.WriteBytes(this.Message)
}
func DecodeIn(msg []byte) Request {
	actionId := msg[0]
	buffer := net.NewBuffer(msg[1:])

	switch actionId {
	case 1:
		request := new(ReturnMessage_in)
		request.Decode(buffer)
		return request
	default:
		return nil
	}
}
