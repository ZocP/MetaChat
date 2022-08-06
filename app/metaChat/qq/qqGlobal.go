package qq

import (
	"MetaChat/pkg/qqbot_framework/commute"
	"MetaChat/pkg/util/cq"
	"MetaChat/pkg/util/cq/condition"
	"github.com/rfyiamcool/go-timewheel"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sync"
	"time"
)

type QQ struct {
	commute.Context
	sync.Locker
	log *zap.Logger

	//提交给上层的消息
	throwCh chan gjson.Result

	echoHandlerMap      *sync.Map
	waitForConditionMap *sync.Map
	//	echoHandlerMap map[string]chan gjson.Result

	//全局时间轮
	tw *timewheel.TimeWheel

	cmdHandlers map[string][]CMDHandler
}

func (qq *QQ) MessageHandler(ctx commute.Context, msg gjson.Result) {
	qq.log.Debug("收到CQ消息", zap.Any("msg", msg))

	//检查消息是否是以前的上下文，如果是则将消息交给channel处理
	if ch, ok := qq.echoHandlerMap.LoadAndDelete(msg.Get(cq.ECHO).String()); ok {
		ch.(chan gjson.Result) <- msg
		return
	}

	//检查消息是否符合注册的条件，如果是则交给相应的handler处理
	qq.waitForConditionMap.Range(func(key, value interface{}) bool {
		if key.(*condition.Condition).Fit(msg) {
			value.(chan gjson.Result) <- msg
			qq.waitForConditionMap.Delete(key)
			return false
		}
		return true
	})

	switch msg.Get(cq.POST_TYPE).String() {
	case cq.POST_TYPE_MESSAGE:
		qq.handleMessage(ctx, msg)
	case cq.POST_TYPE_REQUEST:
		qq.handleRequest(ctx, msg)
	}
}

//设置一个等待消息的通道，使用AwaitConditionResult()可以获取结果
func (qq *QQ) SetAwaitCondition(condition *condition.Condition) {
	qq.waitForConditionMap.Store(condition, make(chan gjson.Result))
}

//返回结果或是否超时， 如果bool是true则代表超时，应主动退出当前的Handler
func (qq *QQ) AwaitConditionResult(condition *condition.Condition) (gjson.Result, bool) {
	var msg gjson.Result
	ch, ok := qq.waitForConditionMap.Load(condition)
	if !ok {
		qq.log.Error("没有注册条件", zap.Any("condition", condition))
		//如果返回true则是超时退出
		return gjson.Result{}, true
	}

	//超时退出
	for {
		select {
		case msg = <-ch.(chan gjson.Result):
			return msg, false
		case <-qq.NewStop():
			return gjson.Result{}, true
		}
	}
}

func (qq *QQ) NewStop() chan bool {
	ch := make(chan bool)
	qq.tw.AfterFunc(30*time.Second, func() {
		ch <- true
		close(ch)
	})
	return ch
}

func (qq *QQ) SetAwaitMessage(echo string) {
	qq.echoHandlerMap.Store(echo, make(chan gjson.Result))
}

//same as previous
func (qq *QQ) AwaitMessage(echo string) (gjson.Result, bool) {

	var msg gjson.Result
	ch, ok := qq.echoHandlerMap.Load(echo)
	if !ok {
		qq.log.Error("没有注册echo", zap.String("echo", echo))
		return gjson.Result{}, true
	}
	for {
		select {
		case msg = <-ch.(chan gjson.Result):
			return msg, false
		case <-qq.NewStop():
			return gjson.Result{}, true
		}
	}
}

func (qq *QQ) throw(result gjson.Result) {
	qq.throwCh <- result
}

func (qq *QQ) GetThrowCh() <-chan gjson.Result {
	return qq.throwCh
}

func NewQQ(log *zap.Logger, bot commute.Context) *QQ {
	tw, err := timewheel.NewTimeWheel(1*time.Second, 360)
	if err != nil {
		log.Error("初始化时间轮失败，有些服务可能无法正常运行", zap.Error(err))
	}

	result := &QQ{
		Locker:              &sync.Mutex{},
		log:                 log,
		Context:             bot,
		throwCh:             make(chan gjson.Result),
		cmdHandlers:         make(map[string][]CMDHandler),
		echoHandlerMap:      &sync.Map{},
		waitForConditionMap: &sync.Map{},
		tw:                  tw,
	}

	commute.AddHandler(result.MessageHandler)

	result.onStart()
	return result
}

func (qq *QQ) onStart() {
	qq.tw.Start()
	qq.registerCMDHandlers()
}

func (qq *QQ) registerCMDHandlers() {
	qq.cmdHandlers["echo"] = []CMDHandler{qq.echo}
	qq.cmdHandlers["help"] = []CMDHandler{}
	qq.cmdHandlers["识图"] = []CMDHandler{qq.recognize}
}

func Provide() fx.Option {
	return fx.Provide(NewQQ)
}
