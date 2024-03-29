package commute

import (
	"MetaChat/pkg/util/cq"
	"go.uber.org/zap"
)

type ImplContext struct {
	QQBot *QQBot
}

func (i *ImplContext) OnStart() {
	i.QQBot.OnStart()
}

func (i *ImplContext) SendMessage(msg cq.CQResponse) {
	i.QQBot.sendMessage(msg)
}

func (i *ImplContext) GetAccountInfo() interface{} {
	return nil
}

func (i *ImplContext) Log() *zap.Logger {
	return i.QQBot.log
}

func NewCtx(qqBot *QQBot) Context {
	return &ImplContext{
		QQBot: qqBot,
	}
}
