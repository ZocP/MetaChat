package cq

import (
	"MetaChat/app/metaChat/cq/config"
	"MetaChat/app/metaChat/cq/ws"
	"MetaChat/app/metaChat/eventBridge/response"
	"MetaChat/pkg/signal"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type CQEventHandler struct {
	config       *config.Config
	conn         *ws.WS
	eventChannel chan gjson.Result
	replyChannel chan response.CQResp
	log          *zap.Logger
	stopHandler  *signal.StopHandler
	stopCh       chan chan bool
	readyCh      chan bool
}

func NewCQEventHandler(viper *viper.Viper, logger *zap.Logger, handler *signal.StopHandler) *CQEventHandler {
	return &CQEventHandler{
		log:          logger,
		config:       config.Unmarshal(viper),
		eventChannel: make(chan gjson.Result),
		replyChannel: make(chan response.CQResp),
		readyCh:      make(chan bool),
		stopHandler:  handler,
	}
}

func (cq *CQEventHandler) OnStart() {
	cq.log.Info("CQEventHandler started, waiting for cq connection")
	cq.stopHandler.Add(cq)
	for {
		select {
		case reply := <-cq.replyChannel:
			if err := cq.conn.WriteJSON(reply); err != nil {
				cq.log.Error("error while writing message", zap.Error(err))
			}
		}
	}
}

func (cq *CQEventHandler) OnStop() error {
	return cq.conn.Close()
}

func (cq *CQEventHandler) OnConnect() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := ws.Upgrade(c.Writer, c.Request)
		if err != nil {
			cq.log.Error("error while upgrading connection", zap.Error(err))
			return
		}
		cq.conn = ws.NewWS(conn, cq.config, cq.log)
		cq.log.Info("connection with CQHttp established")
		lc, err := cq.conn.ReadMessage()
		if err != nil {
			cq.log.Error("error while reading message", zap.Error(err))
			return
		}
		cq.log.Info("life cycle established", zap.String("message", lc.String()))
		cq.readyCh <- true

		go cq.listen()

	}
}

func (cq *CQEventHandler) listen() {
	for {
		eventJson, err := cq.conn.ReadMessage()
		if err != nil {
			cq.log.Error("error while reading message", zap.Error(err))
		}
		cq.eventChannel <- eventJson
	}
}

func (cq *CQEventHandler) IsAdmin(id int64) bool {
	admin := []int64{
		1395437934,
	}
	for _, v := range admin {
		if v == id {
			return true
		}
	}
	return false
}

func (cq *CQEventHandler) GetEventCh() chan gjson.Result {
	return cq.eventChannel
}

func (cq *CQEventHandler) GetReplyCh() chan response.CQResp {
	return cq.replyChannel
}

func (cq *CQEventHandler) GetReadyCh() chan bool {
	return cq.readyCh
}

func Provide() fx.Option {
	return fx.Provide(NewCQEventHandler)
}
