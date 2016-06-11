package net

import (
	"fmt"
	"testing"

	"github.com/Jack301/chat_tool/chat_server/core/log"
)

func Test_buffer(t *testing.T) {
	log.AddTestFlag("Test_buffer")
	defer log.AddTestFlag("")
	buffer := NewBuffer()

	buffer.WriteUint8(1)
	buffer.WriteUint16LE(2)
	buffer.WriteUint32LE(3)
	buffer.WriteUint64LE(4)
	buffer.WriteString("Hello World!")

	fmt.Println(buffer.ReadUint8())
	fmt.Println(buffer.ReadUint16())
	fmt.Println(buffer.ReadUint32())
	fmt.Println(buffer.ReadUint64())
	fmt.Println(buffer.ReadString())
}
