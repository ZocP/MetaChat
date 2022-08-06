package io

import (
	"MetaChat/pkg/util/cq"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type IOHandler interface {
	Handler() gin.HandlerFunc
	GetMessageCh() <-chan gjson.Result
	SendMessage(msg cq.CQResponse)
	GetOnReadyCh() <-chan bool
}
