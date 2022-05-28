package io

import (
	"MetaChat/pkg/cq"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type IOHandler interface {
	OnConnect() gin.HandlerFunc
	GetMessageCh() <-chan gjson.Result
	SendMessage(msg cq.CQResp)
}
