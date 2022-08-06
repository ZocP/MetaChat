package http

import (
	"MetaChat/pkg/minecraft/config"
	"MetaChat/pkg/minecraft/io"
	"MetaChat/pkg/signal"
	"MetaChat/pkg/util"
	"MetaChat/pkg/util/mc"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"io/ioutil"
)

type HTTPHandler struct {
	log       *zap.Logger
	config    *config.Config
	rawCh     chan []byte
	stopCh    chan bool
	messageCh chan gjson.Result
}

func (H *HTTPHandler) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			H.log.Error("read request body error", zap.Error(err))
			c.AbortWithStatus(500)
			return
		}
		H.rawCh <- raw
	}
}

func (H *HTTPHandler) OnStop() error {
	H.stopCh <- true
	return nil
}

func (H *HTTPHandler) GetMessageCh() <-chan gjson.Result {
	return H.messageCh
}

func (H *HTTPHandler) SendMessage(msg mc.MCResponse) {
	data, _ := json.Marshal(msg)
	result, err := util.NewJsonPostRequest(H.config.RemoteAddress, data)
	if err != nil {
		H.log.Error("send message error", zap.Error(err))
		return
	}
	H.log.Info("send message success", zap.Any("result", result))
}
func (H *HTTPHandler) GetOnReadyCh() <-chan bool {
	return nil
}

func NewHttpHandler(logger *zap.Logger, stop *signal.StopHandler) io.IOHandler {
	result := &HTTPHandler{
		log:       logger,
		rawCh:     make(chan []byte),
		stopCh:    make(chan bool),
		messageCh: make(chan gjson.Result),
	}
	stop.Add(result)
	go result.listen()
	return result
}

func (H *HTTPHandler) listen() {
	for {
		select {
		case <-H.stopCh:
			return
		case msg := <-H.rawCh:
			H.messageCh <- gjson.ParseBytes(msg)
		}
	}
}

func Provide() fx.Option {
	return fx.Provide(
		NewHttpHandler,
	)
}
