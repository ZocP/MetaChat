package qqBot

import (
	"MetaChat/pkg/cq"
	"github.com/tidwall/gjson"
)

type Context interface {
	SendMessage(msg cq.CQResp)
	GetAccountInfo() interface{}
}

type EventHandler func(ctx Context, msg gjson.Result)

var EventHandlers = make([]EventHandler, 0)

func AddHandler(handler ...EventHandler) {
	EventHandlers = append(EventHandlers, handler...)
}
