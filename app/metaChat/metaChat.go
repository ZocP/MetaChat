package metaChat

import (
	"MetaChat/app/metaChat/cq"
	"MetaChat/app/metaChat/minecraft"
	"MetaChat/app/metaChat/router"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MetaChat struct {
	log       *zap.Logger
	viper     *viper.Viper
	cqHandler cq.CQEventHandler
	mcHandler minecraft.MCEventHandler
}

func (meta *MetaChat) OnStart() error {
	go func() {
		if err := meta.Listen(); err != nil {
			meta.log.Error("error while listening", zap.Error(err))
		}
	}()
}

func (meta *MetaChat) Listen() error {
	cqch := meta.cqHandler.GetEventCh()
	mcch := meta.mcHandler.GetEventCh()

	for {
		select {
		//停止处理
		// case <- stopCh:

		case cqMsg := <-cqch:

		case mcMsg := <-mcch:

		}
	}
}

func NewMetaChat(log *zap.Logger, viper *viper.Viper, cq cq.CQEventHandler, mc minecraft.MCEventHandler) *MetaChat {
	return &MetaChat{
		log:       log,
		viper:     viper,
		cqHandler: cq,
		mcHandler: mc,
	}
}

func Provide() fx.Option {
	return fx.Provide(
		NewMetaChat,
		router.Provide(),
		fx.Options(cq.Provide(), minecraft.Provide()),
	)
}
