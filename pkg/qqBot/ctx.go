package qqBot

import (
	"MetaChat/pkg/cq"
)

type ImplContext struct {
	QQBot *QQBot
}

func (i *ImplContext) SendMessage(msg cq.CQResp) {
	i.QQBot.sendMessage(msg)
}

func (i *ImplContext) GetAccountInfo() interface{} {
	return nil
}

func NewCtx(qqBot *QQBot) Context {
	return &ImplContext{
		QQBot: qqBot,
	}
}
