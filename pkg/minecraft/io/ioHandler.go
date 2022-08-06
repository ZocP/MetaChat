package io

import (
	"MetaChat/pkg/util/mc"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type IOHandler interface {
	Handler() gin.HandlerFunc
	GetMessageCh() <-chan gjson.Result
	SendMessage(msg mc.MCResponse)
	GetOnReadyCh() <-chan bool
}
