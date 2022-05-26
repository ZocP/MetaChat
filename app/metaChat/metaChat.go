package metaChat

import (
	"MetaChat/app/metaChat/cq"
	"MetaChat/app/metaChat/eventBridge"
	"MetaChat/app/metaChat/minecraft"
	"MetaChat/app/metaChat/router"
	"MetaChat/pkg/signal"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MetaChat struct {
	log       *zap.Logger
	viper     *viper.Viper
	cqHandler *cq.CQEventHandler
	mcHandler *minecraft.MCEventHandler
	stopCh    chan chan bool
	stop      *signal.StopHandler
}

func (meta *MetaChat) OnStart() error {
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
	cqch := meta.cqHandler.GetEventCh()
	mcch := meta.mcHandler.GetEventCh()
	for {
		select {
		case done := <-meta.stopCh:
			//send messages
			done <- true
		case cqMsg := <-cqch:
			eventBridge.LogCQEvent(meta.log, cqMsg)
		case mcMsg := <-mcch:
			meta.log.Info("received mc message", zap.Any("mcMsg", mcMsg))
		}
	}
}

func NewMetaChat(log *zap.Logger, viper *viper.Viper, cq *cq.CQEventHandler, mc *minecraft.MCEventHandler, stop *signal.StopHandler) *MetaChat {
	return &MetaChat{
		log:       log,
		viper:     viper,
		cqHandler: cq,
		mcHandler: mc,
		stopCh:    make(chan chan bool),
		stop:      stop,
	}
}

func Provide() fx.Option {
	return fx.Options(
		fx.Provide(NewMetaChat),
		router.Provide(),
		fx.Options(cq.Provide(), minecraft.Provide()),
	)
}
