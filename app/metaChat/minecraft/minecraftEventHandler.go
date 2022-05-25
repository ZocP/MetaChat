package minecraft

import (
	"MetaChat/app/metaChat/eventBridge"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MCEventHandler struct {
	eventChannel chan eventBridge.MCEvent
	log          *zap.Logger
}

func NewEventHandler(log *zap.Logger) *MCEventHandler {
	return &MCEventHandler{
		eventChannel: make(chan eventBridge.MCEvent),
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
		mc.eventChannel <- eventBridge.MCEvent{}
	}
}

func (mc *MCEventHandler) GetEventCh() chan eventBridge.MCEvent {
	return mc.eventChannel
}

func Provide() fx.Option {
	return fx.Provide(NewEventHandler)
}
