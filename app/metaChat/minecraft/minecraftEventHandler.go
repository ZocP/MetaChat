package minecraft

import (
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MCEventHandler struct {
	eventChannel chan gjson.Result
	replyChannel chan gjson.Result
	log          *zap.Logger
}

func NewEventHandler(log *zap.Logger) *MCEventHandler {
	return &MCEventHandler{
		eventChannel: make(chan gjson.Result),
		log:          log,
	}
}

func (mc *MCEventHandler) OnStart() {
	//process on start event
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

func (mc *MCEventHandler) GetReplyCh() chan gjson.Result {
	return mc.replyChannel
}

func Provide() fx.Option {
	return fx.Provide(NewEventHandler)
}
