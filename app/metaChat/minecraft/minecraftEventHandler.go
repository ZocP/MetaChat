package minecraft

import (
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MCEventHandler struct {
	eventChannel chan gjson.Result
	replyChannel chan interface{}
	log          *zap.Logger
}

func NewEventHandler(log *zap.Logger) *MCEventHandler {
	return &MCEventHandler{
		eventChannel: make(chan gjson.Result),
		replyChannel: make(chan interface{}),
		log:          log,
	}
}

func (mc *MCEventHandler) OnStart() {
	for {
		select {
		case event := <-mc.replyChannel:
			mc.log.Debug("Received event", zap.Any("event", event))
		}
	}
}

func (mc *MCEventHandler) OnStop() error {
	return nil
}

func (mc *MCEventHandler) OnEvent() gin.HandlerFunc {
	return func(c *gin.Context) {
		//unwrap eventBridge and send eventBridge to eventBridge channel
		mc.eventChannel <- gjson.Result{}
	}
}

func (mc *MCEventHandler) GetEventCh() chan gjson.Result {
	return mc.eventChannel
}

func (mc *MCEventHandler) GetReplyCh() chan interface{} {
	return mc.replyChannel
}

func Provide() fx.Option {
	return fx.Provide(NewEventHandler)
}
