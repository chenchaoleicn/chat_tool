package net

import (
	"time"
)

type setting struct {
	wTimeOut time.Duration
	rTimeOut time.Duration
	wMaxSize uint32
	rMaxSize uint32
}

func NewSeting() *setting {
	return new(setting)
}

func (set *setting) SetWriteTimeOut(wTimeOut time.Duration) {
	set.wTimeOut = wTimeOut
}
func (set *setting) SetReadTimeOut(rTimeOut time.Duration) {
	set.rTimeOut = rTimeOut
}

func (set *setting) SetWriteMaxSize(wMaxSize uint32) {
	set.wMaxSize = wMaxSize
}
func (set *setting) SetReadMaxSize(rMaxSize uint32) {
	set.rMaxSize = rMaxSize
}
func (set *setting) GetWriteTimeOut() time.Duration {
	return set.wTimeOut
}
func (set *setting) GetReadTimeOut() time.Duration {
	return set.rTimeOut
}

func (set *setting) GetWriteMaxSize() uint32 {
	return set.wMaxSize
}
func (set *setting) GetReadMaxSize() uint32 {
	return set.rMaxSize
}
