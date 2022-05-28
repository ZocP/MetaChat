package metaChat

import (
	"MetaChat/app/metaChat/config"
	"MetaChat/app/metaChat/minecraft"
	"MetaChat/app/metaChat/qqBot"
	"MetaChat/app/metaChat/router"
	"MetaChat/pkg/signal"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MetaChat struct {
	log       *zap.Logger
	viper     *viper.Viper
	config    *config.Config
	mcHandler *minecraft.MCEventHandler

	stopCh chan chan bool
	stop   *signal.StopHandler

	qqBot *qqBot.QQBot
}

func (meta *MetaChat) OnStart() error {
	meta.stop.Add(meta)
	meta.qqBot.OnStart()
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
	qqMsgCh := meta.qqBot.GetMessageCh()
	for {
		select {
		case done := <-meta.stopCh:

			done <- true
		case cqMsgJson := <-qqMsgCh:
			go meta.handleCQMessage(cqMsgJson)
			//case mcMsgJson := <-meta.mcReceiveCh:
			//	eventBridge.LogCQEvent(meta.log, mcMsgJson)
		}
	}
}

func NewMetaChat(log *zap.Logger, viper *viper.Viper, mc *minecraft.MCEventHandler, stop *signal.StopHandler, bot *qqBot.QQBot) *MetaChat {
	return &MetaChat{
		log:       log,
		viper:     viper,
		mcHandler: mc,
		stopCh:    make(chan chan bool),
		stop:      stop,
		qqBot:     bot,
	}
}

func Provide() fx.Option {
	return fx.Options(
		fx.Provide(NewMetaChat),
		router.Provide(),
		fx.Options(qqBot.Provide(), minecraft.Provide()),
	)
}
