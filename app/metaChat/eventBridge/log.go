package eventBridge

import "go.uber.org/zap"

func LogCQEvent(logger *zap.Logger, event CQEvent) {
	if event.MetaEventType != META_EVENT_TYPE_HEARTBEAT {
		logger.Info("CQEvent", zap.Any("event", event))
	}
}
