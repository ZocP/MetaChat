package metaChat

import (
	"MetaChat/app"
	"MetaChat/app/metaChat/config"
	"MetaChat/app/metaChat/qq"
	"MetaChat/pkg/qqBot"
	"MetaChat/pkg/signal"
	"context"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MetaChat struct {
	log    *zap.Logger
	viper  *viper.Viper
	config *config.Config

	qqMsgCh <-chan gjson.Result
	mcMsgCh <-chan gjson.Result

	stopCh chan chan bool
	stop   *signal.StopHandler
	qq     *qq.QQ
}

func (meta *MetaChat) OnStart() error {
	qqBot.AddHandler(meta.qq.MessageHandler)
	meta.stop.Add(meta)
	go func() {
		if err := meta.Listen(); err != nil {
			meta.log.Error("error while listening", zap.Error(err))
		}
	}()
	return nil
}

func (meta *MetaChat) OnStop() error {
	done := make(chan bool)
	meta.stopCh <- done
	<-done
	return nil
}

func (meta *MetaChat) Listen() error {

	for {
		select {
		case done := <-meta.stopCh:
			done <- true
		case cqMsgJson := <-meta.qqMsgCh:
			go meta.handleCQMessage(cqMsgJson)
		}
	}
}

func NewMetaChat(log *zap.Logger, viper *viper.Viper, stop *signal.StopHandler, qq *qq.QQ) app.APP {
	return &MetaChat{
		log:     log,
		viper:   viper,
		stopCh:  make(chan chan bool),
		stop:    stop,
		qqMsgCh: qq.GetThrow(),
		qq:      qq,
	}
}

func Provide() fx.Option {
	return fx.Options(
		fx.Provide(NewMetaChat),
		fx.Options(qqBot.Provide(), qq.Provide()),
		fx.Invoke(func(meta app.APP, lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return meta.OnStart()
				},
				OnStop: func(ctx context.Context) error {
					return meta.OnStop()
				},
			})
		}),
	)
}
