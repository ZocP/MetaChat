package metaChat

import (
	"MetaChat/app/metaChat/cq"
	"MetaChat/app/metaChat/eventBridge"
	"MetaChat/app/metaChat/eventBridge/response"
	"MetaChat/app/metaChat/minecraft"
	"MetaChat/app/metaChat/router"
	"MetaChat/pkg/signal"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MetaChat struct {
	log         *zap.Logger
	viper       *viper.Viper
	cqHandler   *cq.CQEventHandler
	mcHandler   *minecraft.MCEventHandler
	stopCh      chan chan bool
	stop        *signal.StopHandler
	cqch        chan gjson.Result
	mcch        chan gjson.Result
	cqreplych   chan response.CQResp
	mcreplych   chan gjson.Result
	taskHandler map[int64]func()
}

func (meta *MetaChat) OnStart() error {
	meta.stop.Add(meta)
	go meta.cqHandler.OnStart()
	go meta.mcHandler.OnStart()
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
	meta.cqch = meta.cqHandler.GetEventCh()
	meta.mcch = meta.mcHandler.GetEventCh()
	meta.cqreplych = meta.cqHandler.GetReplyCh()
	meta.mcreplych = meta.mcHandler.GetReplyCh()
	for {
		select {
		case done := <-meta.stopCh:
			//send messages
			if err := meta.cqHandler.OnStop(); err != nil {
				meta.log.Error("error while stopping cq handler", zap.Error(err))
			}
			if err := meta.mcHandler.OnStop(); err != nil {
				meta.log.Error("error while stopping mc handler", zap.Error(err))
			}
			done <- true
		case cqMsgJson := <-meta.cqch:
			meta.handleCQMessage(cqMsgJson)
		case mcMsgJson := <-meta.mcch:
			eventBridge.LogCQEvent(meta.log, mcMsgJson)
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
