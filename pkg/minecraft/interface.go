package minecraft

import (
	"MetaChat/pkg/minecraft/io/http"
	"MetaChat/pkg/util/mc"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
)

type Context interface {
	GetThrowCh() chan gjson.Result
	SendMessage(resp mc.MCResponse)
}

func Provide() fx.Option {
	return fx.Options(fx.Provide(NewController), http.Provide())
}
