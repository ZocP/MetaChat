package commute

import (
	"MetaChat/pkg/util/cq"
	"MetaChat/pkg/util/cq/condition"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"sync"
)

type Context interface {
	Log() *zap.Logger
	OnStart()
	SendMessage(msg cq.CQResponse)
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
