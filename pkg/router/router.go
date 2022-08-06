package router

import (
	"MetaChat/pkg/http"
	io2 "MetaChat/pkg/minecraft/io"
	"MetaChat/pkg/qqbot_framework/commute/io"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func InitRouters(cqio io.IOHandler, mcio io2.IOHandler) http.InitRouters {
	return func(r *gin.Engine) {
		r.GET("/v1/cq", cqio.Handler())
		r.POST("/v1/mc", mcio.Handler())
	}
}

func Provide() fx.Option {
	return fx.Provide(InitRouters)
}
