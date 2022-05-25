package signal

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type IStop interface {
	OnStop() error
}

type StopHandler struct {
	log    *zap.Logger
	StopCB []IStop
}

func NewStopHandler(log *zap.Logger) *StopHandler {
	return &StopHandler{
		log: log,
	}
}

func Provide() fx.Option {
	return fx.Provide(
		NewStopHandler,
	)
}

func (sh *StopHandler) Add(stop IStop) {
	sh.StopCB = append(sh.StopCB, stop)
}

func (sh *StopHandler) Stop() {
	for _, v := range sh.StopCB {
		if err := v.OnStop(); err != nil {

		}
	}
}
