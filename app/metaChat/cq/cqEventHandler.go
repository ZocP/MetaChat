package cq

import (
	"MetaChat/app/metaChat/cq/config"
	"MetaChat/app/metaChat/cq/event"
	"MetaChat/app/metaChat/cq/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type CQEventHandler struct {
	config       *config.Config
	conn         *websocket.Conn
	eventChannel chan event.CQEvent
}

func NewCQEventHandler(viper *viper.Viper) *CQEventHandler {
	return &CQEventHandler{
		config: config.Unmarshal(viper),
	}
}

func (cq *CQEventHandler) OnConnect() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := ws.Upgrade(c.Writer, c.Request)
		if err != nil {
			return
		}
		cq.conn = conn
	}
}

func Provide() fx.Option {
	return fx.Provide(NewCQEventHandler)
}

func (cq *CQEventHandler) GetEventCh() chan event.CQEvent {
	return cq.eventChannel
}
