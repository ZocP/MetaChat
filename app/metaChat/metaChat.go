package metaChat

import (
	"MetaChat/app/metaChat/cq"
	"MetaChat/app/metaChat/eventBridge"
	"MetaChat/app/metaChat/eventBridge/request"
	"MetaChat/app/metaChat/minecraft"
	"MetaChat/app/metaChat/router"
	"MetaChat/pkg/signal"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
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
	cqch := meta.cqHandler.GetEventCh()
	mcch := meta.mcHandler.GetEventCh()
	//cqreplych := meta.cqHandler.GetReplyCh()
	//mcreplych := meta.mcHandler.GetReplyCh()
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
		case cqMsgJson := <-cqch:
			//eventBridge.LogCQEvent(meta.log, cqMsgJson)
			meta.handleCQMessage(cqMsgJson)
		case mcMsgJson := <-mcch:
			eventBridge.LogCQEvent(meta.log, mcMsgJson)
		}
	}
}

func (meta *MetaChat) handleCQMessage(msg gjson.Result) {
	postType := msg.Get(request.POST_TYPE).String()
	switch postType {
	case request.POST_TYPE_MESSAGE:
		meta.handleCQPostMsg(msg)
	case request.POST_TYPE_REQUEST:

	}
}

func (meta *MetaChat) handleCQPostMsg(msg gjson.Result) {
	switch msg.Get(request.MESSAGE_TYPE).String() {
	case request.MESSAGE_TYPE_GROUP:
		meta.handleCQPostMsgGroup(msg)
	case request.MESSAGE_TYPE_PRIVATE:
		meta.handleCQPostMsgPrivate(msg)

	}
}

func (meta *MetaChat) handleCQPostMsgGroup(msg gjson.Result) {
	meta.log.Info("receive group message", zap.Any("msg", msg.Get(request.MESSAGE).String()))
}

func (meta *MetaChat) handleCQPostMsgPrivate(msg gjson.Result) {
	meta.log.Info("receive private message", zap.Any("msg", msg.Get(request.MESSAGE).String()))
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
