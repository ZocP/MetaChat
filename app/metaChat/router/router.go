package router

import (
	"MetaChat/app/metaChat/minecraft"
	"MetaChat/app/metaChat/qqBot/io"
	"MetaChat/pkg/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func InitRouters(cqio io.IOHandler, mc *minecraft.MCEventHandler) http.InitRouters {
	return func(r *gin.Engine) {
		r.GET("/v1/cq", cqio.OnConnect())
		r.POST("/v1/minecraft/event", mc.OnEvent())
	}
}

func Provide() fx.Option {
	return fx.Provide(InitRouters)
}
