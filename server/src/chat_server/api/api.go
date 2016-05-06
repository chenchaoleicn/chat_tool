package api

import (
	"chat_server/api/protocol/api_handleMessage"
	"chat_server/core/net"
)

type Request interface {
	Process(*net.Session)
	Decode(*net.Buffer)
}

func DecodeIn(msg []byte) Request {
	moduleId := msg[0]

	switch moduleId {
	case 1:
		return api_handleMessage.DecodeIn(msg[1:])
	default:
		return nil
	}
}
