package userMgr

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type UserManager struct {
	log *zap.Logger
}

func NewUserManager(log *zap.Logger) UserManager {
	return UserManager{
		log: log,
	}
}

func Provide() fx.Option {
	return fx.Options(
		fx.Provide(NewUserManager),
	)
}
