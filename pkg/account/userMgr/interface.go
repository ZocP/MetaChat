package userMgr

import (
	"MetaChat/pkg/qqBot/account/status"
	"MetaChat/pkg/qqBot/pkg/commands"
	"MetaChat/pkg/qqBot/pkg/user"
)

type IUserManager interface {
	OnStart() error
	OnStop() error
	CanExecuteCommand(userID string, cmd commands.Command) status.AllowStatus
	StoreUser(user user.User, accountType user.AccountType) error
	GetAccountType(userID string) user.AccountType
	SetAccountType(userID string, accountType user.AccountType) error
}
