package qqBot

import (
	"MetaChat/pkg/cq"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"regexp"
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
	groupid := msg.Get(cq.GROUP_ID).Int()
	user := msg.Get(cq.USER_ID).Int()
	qq.log.Info("receive & processing group message from", zap.Int64("group", groupid), zap.Int64("user", user), zap.Any("msg", msg.Get(cq.MESSAGE).String()))
	message := msg.Get(cq.MESSAGE).String()
	compiler, err := regexp.Compile("^//")
	if err != nil {
		panic("compiler error")
	}
	if compiler.MatchString(message) {
		qq.log.Debug("On Bot Command", zap.Int64("group", groupid), zap.Int64("user", user), zap.Any("msg", message))
		qq.onBotCommand(msg, msg.Get(cq.MESSAGE_TYPE).String())
		return
	}

}

func (qq *QQBot) onPostNotice(msg gjson.Result) {
	qq.log.Debug("onPostNotice", zap.Any("msg", msg))
}
