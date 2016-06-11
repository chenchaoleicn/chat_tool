//字符串在写入buffer前要转换为[]byte格式，先写入对应的[]byte长度，然后再写入字符串, 字符串长度作为两个字节长度写入

package net

import "encoding/binary"

var littleEndian = binary.LittleEndian

type Buffer struct {
	raw  []byte
	rPos int
	wPos int
}

func NewBuffer(raw []byte) *Buffer {
	return &Buffer{
		raw: raw,
	}
}
func (buf *Buffer) grows(length int) int {
	begin := buf.wPos
	if len(buf.raw)-buf.wPos < length {
		newBuffer := make([]byte, len(buf.raw)+length+100)
		copy(newBuffer, buf.raw)
		buf.raw = newBuffer
	}
	buf.wPos += length
	return begin
}
func (buf *Buffer) Set(raw []byte) {
	buf.raw = raw
}

func (buf *Buffer) Get() []byte {
	return buf.raw[:buf.wPos]
}

func (buf *Buffer) Cap() int {
	return len(buf.raw)
}
func (buf *Buffer) SetWritePos(wPos int) {
	buf.wPos = wPos
}

func (buf *Buffer) GetWritePos() int {
	return buf.wPos
}

func (buf *Buffer) SetReadPos(rPos int) {
	buf.rPos = rPos
}

func (buf *Buffer) GetReadPos() int {
	return buf.rPos
}

func (buf *Buffer) WriteBytes(bytes []byte) {
	begin := buf.grows(len(bytes))
	copy(buf.raw[begin:], bytes)
}
func (buf *Buffer) WriteUint8(num uint8) {
	begin := buf.grows(1)
	buf.raw[begin] = num
}

func (buf *Buffer) WriteUint16LE(num uint16) {
	begin := buf.grows(2)
	littleEndian.PutUint16(buf.raw[begin:buf.wPos], num)
}

func (buf *Buffer) WriteUint32LE(num uint32) {
	begin := buf.grows(4)
	littleEndian.PutUint32(buf.raw[begin:buf.wPos], num)
}

func (buf *Buffer) WriteUint64LE(num uint64) {
	begin := buf.grows(8)
	littleEndian.PutUint64(buf.raw[begin:buf.wPos], num)
}

func (buf *Buffer) WriteString(str string) {
	bytes := []byte(str)
	begin := buf.grows(len(bytes))
	copy(buf.raw[begin:], bytes)
}

func (buf *Buffer) ReadUint8() (num uint8) {
	begin := buf.rPos
	num = buf.raw[begin]
	buf.rPos += 1
	return
}

func (buf *Buffer) ReadUint16() (num uint16) {
	begin := buf.rPos
	buf.rPos += 2
	num = littleEndian.Uint16(buf.raw[begin:buf.rPos])
	return
}
func (buf *Buffer) ReadUint32() (num uint32) {
	begin := buf.rPos
	buf.rPos += 4
	num = littleEndian.Uint32(buf.raw[begin:buf.rPos])
	return
}
func (buf *Buffer) ReadUint64() (num uint64) {
	begin := buf.rPos
	buf.rPos += 8
	num = littleEndian.Uint64(buf.raw[begin:buf.rPos])
	return
}

func (buf *Buffer) ReadBytes(length uint16) []byte {
	bytes := make([]byte, length)
	begin := buf.rPos
	buf.rPos += int(length)
	copy(bytes, buf.raw[begin:buf.rPos])
	return bytes
}
