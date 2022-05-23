package minecraft

import (
	"MetaChat/app/metaChat/minecraft/event"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type MCEventHandler struct {
	eventChannel chan event.MCEvent
}

func NewEventHandler() *MCEventHandler {
	return &MCEventHandler{}
}

func (mc *MCEventHandler) OnEvent() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (mc *MCEventHandler) GetEventCh() chan event.MCEvent {
	return mc.eventChannel
}

func Provide() fx.Option {
	return fx.Provide(NewEventHandler)
}
