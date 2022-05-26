package eventBridge

import (
	"MetaChat/app/metaChat/eventBridge/request"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func LogCQEvent(logger *zap.Logger, event gjson.Result) {
	if event.Get(request.META_EVENT_TYPE).String() != request.META_EVENT_TYPE_HEARTBEAT {
		logger.Info("CQ event", zap.Any("event", event.String()))
	}
}
