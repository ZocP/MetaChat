package qqBot

import (
	"MetaChat/pkg/cq"
	"MetaChat/pkg/qqBot/config"
	"MetaChat/pkg/qqBot/io"
	"MetaChat/pkg/qqBot/io/ws"
	"MetaChat/pkg/signal"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type QQBot struct {
	log    *zap.Logger
	config *config.Config

	handlerAddCh chan EventHandler
	stopCh       chan bool
	connReadyCh  <-chan bool
	throwCh      chan gjson.Result

	eventHandlers []EventHandler
	ctx           Context

	IOHandler io.IOHandler
	isReady   bool
}

func NewQQBot(log *zap.Logger, config *config.Config, handler io.IOHandler, stopHandler *signal.StopHandler) *QQBot {
	bot := &QQBot{
		log:           log,
		config:        config,
		stopCh:        make(chan bool),
		IOHandler:     handler,
		throwCh:       make(chan gjson.Result),
		isReady:       false,
		eventHandlers: EventHandlers,
	}
	bot.OnStart()
	stopHandler.Add(bot)
	return bot
}

func (qq *QQBot) OnStart() {
	qq.connReadyCh = qq.IOHandler.GetOnReadyCh()
	msgCh := qq.IOHandler.GetMessageCh()
	go func() {
		for {
			select {
			case ready := <-qq.connReadyCh:
				if ready {
					qq.isReady = true
					go qq.initContext()
				} else {
					qq.isReady = false
				}
			case msg := <-msgCh:
				if qq.isReady {
					qq.onMessage(msg)
				}
				qq.log.Info("receive message but bot not ready yet", zap.Any("msg", msg))
			case <-qq.stopCh:
				qq.notifyStop()
				break
			case handler := <-qq.handlerAddCh:
				qq.eventHandlers = append(qq.eventHandlers, handler)
			}
		}
	}()

}

func (qq *QQBot) OnStop() error {
	qq.stopCh <- true
	return nil
}

func (qq *QQBot) sendMessage(msg cq.CQResp) {
	qq.IOHandler.SendMessage(msg)
}

func (qq *QQBot) GetThrow() <-chan gjson.Result {
	return qq.throwCh
}

//将需要上层MC处理的消息放入channel，等待上层MetaChat处理
func (qq *QQBot) throw(msg gjson.Result) {
	qq.throwCh <- msg
}

func (qq *QQBot) onMessage(msg gjson.Result) {
	for _, eh := range qq.eventHandlers {
		if eh != nil {
			eh(qq.ctx, msg)
		}
	}
}

func (qq *QQBot) notifyStop() {
	//TODO: Do something when bot is stopped
	//qq.sendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetPrivateMessage()))
}

func (qq *QQBot) initContext() {
	qq.ctx = NewCtx(qq)
}

func (qq *QQBot) handlerAdder(handler ...EventHandler) {
	for _, v := range handler {
		qq.handlerAddCh <- v
	}
}

func Provide() fx.Option {
	return fx.Options(fx.Provide(
		NewQQBot,
		ws.NewWS,
		config.NewConfig,
	))
}
