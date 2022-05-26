package cq

import (
	"MetaChat/app/metaChat/cq/config"
	"MetaChat/app/metaChat/cq/ws"
	"MetaChat/pkg/signal"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type CQEventHandler struct {
	config       *config.Config
	conn         *websocket.Conn
	eventChannel chan gjson.Result
	replyChannel chan gjson.Result
	log          *zap.Logger
	stopHandler  *signal.StopHandler
	stopCh       chan chan bool
	qqBot        *QQBot
}

func NewCQEventHandler(viper *viper.Viper, logger *zap.Logger, handler *signal.StopHandler) *CQEventHandler {
	return &CQEventHandler{
		log:          logger,
		config:       config.Unmarshal(viper),
		eventChannel: make(chan gjson.Result),
		stopHandler:  handler,
	}
}

func (cq *CQEventHandler) OnStart() {
	cq.stopHandler.Add(cq)
	for {
		select {
		case reply := <-cq.replyChannel:
			if err := cq.conn.WriteJSON(reply); err != nil {
				cq.log.Error("error while writing message", zap.Error(err))
			}

			//case stop := <-cq.stopCh:
			//	if err := cq.conn.WriteJSON(); err != nil{
			//		cq.log.Error("error while writing message", zap.Error(err))
			//	}
			//

		}
	}
}

func (cq *CQEventHandler) OnStop() error {
	cq.stopCh <- make(chan bool)
	<-cq.stopCh
	cq.conn.Close()
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

func (cq *CQEventHandler) GetBot() *QQBot {
	return cq.qqBot
}

func (cq *CQEventHandler) listen() {
	for {
		_, message, err := cq.conn.ReadMessage()
		if err != nil {
			cq.log.Error("error while reading message", zap.Error(err))
		}
		eventJson := gjson.Parse(string(message))
		if err != nil {
			cq.log.Error("error while unmarshalling message", zap.Error(err))
		}
		cq.eventChannel <- eventJson
	}
}

func Provide() fx.Option {
	return fx.Provide(NewCQEventHandler)
}

func (cq *CQEventHandler) GetEventCh() chan gjson.Result {
	return cq.eventChannel
}

func (cq *CQEventHandler) GetReplyCh() chan gjson.Result {
	return cq.replyChannel
}

func getStopMessage() {

}
