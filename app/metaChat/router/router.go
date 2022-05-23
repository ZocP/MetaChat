package router

import (
	"MetaChat/app/metaChat/cq"
	"MetaChat/app/metaChat/minecraft"
	"MetaChat/pkg/http"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func InitRouters(viper *viper.Viper, cq cq.CQEventHandler, mc minecraft.EventHandler) http.InitRouters {
	return func(r *gin.Engine) {
		r.GET("/v1/cq", cq.OnConnect())
		r.POST("/v1/minecraft/event", mc.OnEvent())
	}
}

func Provide() fx.Option {
	return fx.Provide(InitRouters)
}
