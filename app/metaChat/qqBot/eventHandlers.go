package qqBot

import (
	"MetaChat/pkg/cq"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func (qq *QQBot) RegisterHandlers() {
	qq.eventHandlers = map[string]map[string]func(msg gjson.Result){
		cq.POST_TYPE: {
			cq.POST_TYPE_MESSAGE: qq.onPostMessage,
			cq.POST_TYPE_NOTICE:  qq.onPostNotice,
		},
	}
}

func (qq *QQBot) onPostMessage(msg gjson.Result) {
	qq.log.Debug("onPostMessage", zap.Any("msg", msg))
}

func (qq *QQBot) onPostNotice(msg gjson.Result) {
	qq.log.Debug("onPostNotice", zap.Any("msg", msg))
}
