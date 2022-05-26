package cq

import (
	"MetaChat/app/metaChat/cq/config"
	"MetaChat/app/metaChat/cq/ws"
	"MetaChat/app/metaChat/eventBridge"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type CQEventHandler struct {
	config       *config.Config
	conn         *websocket.Conn
	eventChannel chan eventBridge.CQEvent
	log          *zap.Logger
}

func NewCQEventHandler(viper *viper.Viper, logger *zap.Logger) *CQEventHandler {
	return &CQEventHandler{
		log:          logger,
		config:       config.Unmarshal(viper),
		eventChannel: make(chan eventBridge.CQEvent),
	}
}

func (cq *CQEventHandler) Onstart() error {
	return nil
}

func (cq *CQEventHandler) OnConnect() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := ws.Upgrade(c.Writer, c.Request)
		if err != nil {
			cq.log.Error("error while upgrading connection", zap.Error(err))
			return
		}
		cq.conn = conn
		cq.log.Info("connection with CQHttp established")
		go cq.listen()
	}
}

func (cq *CQEventHandler) listen() {
	for {
		_, message, err := cq.conn.ReadMessage()
		if err != nil {
			cq.log.Error("error while reading message", zap.Error(err))
		}
		event, err := eventBridge.UnmarshalCQEvent(message)
		if err != nil {
			cq.log.Error("error while unmarshalling message", zap.Error(err))
		}
		cq.eventChannel <- event
	}
}

func Provide() fx.Option {
	return fx.Provide(NewCQEventHandler)
}

func (cq *CQEventHandler) GetEventCh() chan eventBridge.CQEvent {
	return cq.eventChannel
}
