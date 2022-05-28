package io

import (
	"MetaChat/pkg/cq"
	"github.com/tidwall/gjson"
)

type IOHandler interface {
	GetMessageCh() <-chan gjson.Result
	SendMessage(msg cq.CQResp)
}
