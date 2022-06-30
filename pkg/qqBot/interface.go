package qqBot

import (
	"MetaChat/pkg/cq"
	"MetaChat/pkg/cq/condition"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"sync"
)

type Context interface {
	Log() *zap.Logger
	OnStart()
	SendMessage(msg cq.CQResp)
}

type EventHandler func(ctx Context, msg gjson.Result)
type ConditionHandler func(ctx Context, msg gjson.Result)

var EventHandlers = make([]EventHandler, 0)

var ConditionHandlers = &sync.Map{}

func AddHandler(handler ...EventHandler) {
	EventHandlers = append(EventHandlers, handler...)
}

func AddConditionHandler(condition *condition.Condition, handler ...ConditionHandler) {
	ConditionHandlers.Store(condition, handler)
}
