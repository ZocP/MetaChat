package router

import (
	"MetaChat/pkg/http"
	"MetaChat/pkg/qqBot/io"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func InitRouters(cqio io.IOHandler) http.InitRouters {
	return func(r *gin.Engine) {
		r.GET("/v1/cq", cqio.OnConnect())
		//r.POST("/v1/minecraft/event", mc.OnEvent())
	}
}

func Provide() fx.Option {
	return fx.Provide(InitRouters)
}
