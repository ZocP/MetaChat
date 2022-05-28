package qqBot

import (
	"MetaChat/app/metaChat/qqBot/account"
	"MetaChat/app/metaChat/qqBot/config"
	"MetaChat/app/metaChat/qqBot/io"
	"MetaChat/app/metaChat/qqBot/io/ws"
	"MetaChat/pkg/cq"
	"MetaChat/pkg/signal"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type QQBot struct {
	log *zap.Logger

	account *account.AccountInfo

	config *config.Config

	stop   signal.StopHandler
	stopCh chan bool

	throwCh chan gjson.Result
	getCh   chan cq.CQResp

	//use for handling message with event needed
	echoHandlerCh map[string]chan gjson.Result
	//use for send and receive message
	IOHandler io.IOHandler

	isReady bool
}

func NewQQBot(log *zap.Logger, config *config.Config, handler io.IOHandler, stopHandler signal.StopHandler) *QQBot {
	return &QQBot{
		log:       log,
		config:    config,
		stop:      stopHandler,
		stopCh:    make(chan bool),
		IOHandler: handler,
	}
}

func (qq *QQBot) OnStart() {
	qq.stop.Add(qq)
	//TODO: init WS
	//loop to listen message from IOHandler
	go func() {
		for {
			select {
			case msg := <-qq.IOHandler.GetMessageCh():
				go qq.onMessage(msg)
			case <-qq.stopCh:
				qq.notifyStop()
				break
			}
		}
	}()

}

func (qq *QQBot) OnStop() error {
	qq.stopCh <- true
	return nil
}

func (qq *QQBot) onMessage(msg gjson.Result) {
	//TODO: handle event
	//handle message that is registered and wait for processing result
	if ch, ok := qq.echoHandlerCh[msg.Get(cq.ECHO).String()]; ok {
		ch <- msg
		return
	}

	qq.log.Info("OnMessage", zap.String("msg", msg.String()))
}

func (qq *QQBot) SendMessage(msg cq.CQResp) {
	qq.IOHandler.SendMessage(msg)
}

func (qq *QQBot) notifyStop() {
	//TODO: send stop message to super admin
	//qq.SendMessage(cq.GetCQResp(cq.ACTION_SEND_MESSAGE, cq.GetPrivateMessage()))
}

func (qq *QQBot) GetThrowMessageCh() <-chan gjson.Result {
	return qq.throwCh
}

func (qq *QQBot) initAccountInfo() {
	login, loginID := cq.GetCQRespEcho(cq.ACTION_GET_LOGIN_INFO, nil)
	qq.SendMessage(login)
	qq.RegisterEchoHandler(loginID)
	loginResult := qq.WaitForResult(loginID)

	group, groupID := cq.GetCQRespEcho(cq.ACTION_GET_GROUP_LIST, nil)
	qq.SendMessage(group)
	qq.RegisterEchoHandler(groupID)
	groupResult := qq.WaitForResult(groupID)

	groupList := make(map[int64]*account.Group)

	groupResult.Get(cq.DATA).ForEach(func(key, value gjson.Result) bool {
		groupList[value.Get(cq.GROUP_ID).Int()] = &account.Group{
			GroupID:   value.Get(cq.GROUP_ID).Int(),
			GroupName: value.Get(cq.GROUP_NAME).String(),
		}
		return true
	})

	friend, friendID := cq.GetCQRespEcho(cq.ACTION_GET_FRIEND_LIST, nil)
	qq.SendMessage(friend)
	qq.RegisterEchoHandler(friendID)
	friendResult := qq.WaitForResult(friendID)

	friendList := make(map[int64]*account.User)

	friendResult.Get(cq.DATA).ForEach(func(key, value gjson.Result) bool {
		friendList[value.Get(cq.USER_ID).Int()] = &account.User{
			UserID:   value.Get(cq.USER_ID).Int(),
			Nickname: value.Get(cq.NICKNAME).String(),
		}
		return true
	})

	qq.account = account.NewAccountInfo(
		loginResult.Get(cq.USER_ID).Int(),
		loginResult.Get(cq.NICKNAME).String(),
		friendList,
		groupList,
	)

}

func (qq *QQBot) IsReady() bool {
	return qq.isReady
}

func (qq *QQBot) RegisterEchoHandler(id string) {
	qq.echoHandlerCh[id] = make(chan gjson.Result)
}

func (qq *QQBot) WaitForResult(id string) gjson.Result {
	data := <-qq.echoHandlerCh[id]
	delete(qq.echoHandlerCh, id)
	return data
}

func Provide() fx.Option {
	return fx.Provide(
		NewQQBot,
		ws.NewWS,
	)
}
