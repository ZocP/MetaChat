package qq

import (
	"MetaChat/pkg/cq"
	"MetaChat/pkg/cq/condition"
	"MetaChat/pkg/qqBot"
	"github.com/rfyiamcool/go-timewheel"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sync"
	"time"
)

type QQ struct {
	qqBot.Context
	sync.Locker
	log     *zap.Logger
	throwCh chan gjson.Result

	echoHandlerMap      *sync.Map
	waitForConditionMap *sync.Map
	//	echoHandlerMap map[string]chan gjson.Result

	tw *timewheel.TimeWheel

	cmdHandlers map[string][]CMDHandler
}

func (qq *QQ) MessageHandler(ctx qqBot.Context, msg gjson.Result) {
	qq.log.Debug("收到CQ消息", zap.Any("msg", msg))

	if ch, ok := qq.echoHandlerMap.LoadAndDelete(msg.Get(cq.ECHO).String()); ok {
		ch.(chan gjson.Result) <- msg
		return
	}

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

func (qq *QQ) SetAwaitCondition(condition *condition.Condition) {
	qq.waitForConditionMap.Store(condition, make(chan gjson.Result))
}

//返回值 结果或一个bool， 如果bool是true则代表超时，应该主动退出当前的handler
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

func (qq *QQ) GetThrow() <-chan gjson.Result {
	return qq.throwCh
}

func NewQQ(log *zap.Logger, bot qqBot.Context) *QQ {
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
