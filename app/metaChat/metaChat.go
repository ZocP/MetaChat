package metaChat

import (
	"MetaChat/app/metaChat/cq"
	"MetaChat/app/metaChat/cq/group"
	"MetaChat/app/metaChat/cq/user"
	"MetaChat/app/metaChat/eventBridge"
	"MetaChat/app/metaChat/eventBridge/request"
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
	log   *zap.Logger
	viper *viper.Viper

	cqHandler *cq.CQEventHandler
	mcHandler *minecraft.MCEventHandler

	cqReadyCh chan bool
	stopCh    chan chan bool
	stop      *signal.StopHandler

	cqReceiveCh chan gjson.Result
	mcReceiveCh chan gjson.Result
	cqreplych   chan response.CQResp
	//mcreplych   chan response.MCResp
	mcreplych chan interface{}

	echoHandlerCh map[string]chan gjson.Result

	qqBot *cq.QQBot

	botIsReady bool
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
	meta.cqReadyCh = meta.cqHandler.GetReadyCh()
	meta.cqReceiveCh = meta.cqHandler.GetEventCh()
	meta.mcReceiveCh = meta.mcHandler.GetEventCh()
	meta.cqreplych = meta.cqHandler.GetReplyCh()
	meta.mcreplych = meta.mcHandler.GetReplyCh()
	for {
		select {
		case done := <-meta.stopCh:
			//send messages
			if meta.botIsReady {
				for k, _ := range meta.qqBot.GetAdminList() {
					meta.SendToQQ(response.GetCQResp(response.ACTION_SEND_MESSAGE, response.GetPrivateMessage(k, "Bot服务器已关闭")))
				}
			}
			if err := meta.cqHandler.OnStop(); err != nil {
				meta.log.Error("error while stopping cq handler", zap.Error(err))
			}
			if err := meta.mcHandler.OnStop(); err != nil {
				meta.log.Error("error while stopping mc handler", zap.Error(err))
			}
			done <- true
		case <-meta.cqReadyCh:
			go meta.initBot()
		case cqMsgJson := <-meta.cqReceiveCh:
			meta.handleCQMessage(cqMsgJson)
		case mcMsgJson := <-meta.mcReceiveCh:
			eventBridge.LogCQEvent(meta.log, mcMsgJson)
		}
	}
}
func (meta *MetaChat) SendToQQ(msg response.CQResp) {
	meta.cqreplych <- msg
}

func (meta *MetaChat) deleteCQEchoHandler(id ...string) {
	for _, v := range id {
		close(meta.echoHandlerCh[v])
		delete(meta.echoHandlerCh, v)
	}
}

func (meta *MetaChat) initBot() {
	aRaw, aEcho := response.GetCQRespEcho(response.ACTION_GET_LOGIN_INFO, nil)
	meta.SendToQQ(aRaw)
	meta.echoHandlerCh[aEcho] = make(chan gjson.Result)
	accountInfo := <-meta.echoHandlerCh[aEcho]
	raw, gEcho := response.GetCQRespEcho(response.ACTION_GET_GROUP_LIST, nil)
	meta.SendToQQ(raw)
	meta.echoHandlerCh[gEcho] = make(chan gjson.Result)
	groupReady := make(chan interface{})
	go func(ch chan interface{}) {
		msg := <-meta.echoHandlerCh[gEcho]
		groupInfo := msg
		groupList := make(map[int64]*group.Group)
		groupInfo.Get(request.DATA).ForEach(func(key, value gjson.Result) bool {
			groupList[value.Get("group_id").Int()] = &group.Group{

				GroupID:   value.Get("group_id").Int(),
				GroupName: value.Get("group_name").String(),
				BotMode:   group.MODE_REPEAT,
			}
			return true
		})
		ch <- groupList
		close(ch)
	}(groupReady)
	groupList := <-groupReady

	rawList, listEcho := response.GetCQRespEcho(response.ACTION_GET_FRIEND_LIST, nil)
	meta.SendToQQ(rawList)
	meta.echoHandlerCh[listEcho] = make(chan gjson.Result)
	friendReady := make(chan interface{})
	go func() {
		msg := <-meta.echoHandlerCh[listEcho]
		friendList := make(map[int64]*user.User)
		msg.Get(request.DATA).ForEach(func(key, value gjson.Result) bool {
			character := user.NORMAL
			//if user is admin
			usert := &user.User{
				Nickname:  value.Get("nickname").String(),
				Character: character,
			}
			friendList[value.Get("user_id").Int()] = usert
			return true
		})
		friendReady <- friendList
		close(friendReady)
	}()

	friendList := <-friendReady

	meta.log.Debug("QQBot is ready")
	meta.botIsReady = true
	meta.qqBot = &cq.QQBot{
		AccountId:  accountInfo.Get("user_id").Int(),
		Nickname:   accountInfo.Get("nickname").String(),
		FriendList: friendList.(map[int64]*user.User),
		GroupList:  groupList.(map[int64]*group.Group),
	}

	meta.deleteCQEchoHandler(aEcho, listEcho, gEcho)
}

func NewMetaChat(log *zap.Logger, viper *viper.Viper, cq *cq.CQEventHandler, mc *minecraft.MCEventHandler, stop *signal.StopHandler) *MetaChat {
	return &MetaChat{
		log:           log,
		viper:         viper,
		cqHandler:     cq,
		mcHandler:     mc,
		stopCh:        make(chan chan bool),
		stop:          stop,
		echoHandlerCh: make(map[string]chan gjson.Result),
	}
}

func Provide() fx.Option {
	return fx.Options(
		fx.Provide(NewMetaChat),
		router.Provide(),
		fx.Options(cq.Provide(), minecraft.Provide()),
	)
}
