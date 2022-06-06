package account

import (
	"MetaChat/pkg/qqBot/account/userMgr"
	"go.uber.org/fx"
)

type AccountManager interface {
	userMgr.IUserManager
}

func Provide() fx.Option {
	return fx.Provide(
		NewUserManager,
	)
}
