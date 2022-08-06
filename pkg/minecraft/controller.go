package minecraft

import (
	io2 "MetaChat/pkg/minecraft/io"
	"MetaChat/pkg/signal"
	"MetaChat/pkg/util/mc"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

type Controller struct {
	log *zap.Logger
	io  io2.IOHandler

	msgCh <-chan gjson.Result

	throw chan gjson.Result

	stop chan chan bool
}

func (c *Controller) GetThrowCh() chan gjson.Result {
	return c.throw
}

func (c *Controller) onStart() {
	go c.listen()
}

func (c *Controller) OnStop() error {
	done := make(chan bool)
	c.stop <- done

	<-done
	return nil
}
func (c *Controller) listen() {
	ch := c.io.GetMessageCh()
	for {
		select {
		case msg := <-ch:
			go c.handle(msg)
		case stop := <-c.stop:
			//notify stop
			stop <- true
		}
	}
}

func (c *Controller) SendMessage(response mc.MCResponse) {
	c.io.SendMessage(response)
}

func NewController(log *zap.Logger, io io2.IOHandler, stop *signal.StopHandler) Context {
	res := &Controller{
		log:   log,
		io:    io,
		stop:  make(chan chan bool),
		throw: make(chan gjson.Result),
	}
	stop.Add(res)
	res.onStart()
	return res
}

func (c *Controller) handle(msg gjson.Result) {
	c.throw <- msg
}
