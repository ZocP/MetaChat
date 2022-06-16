package userMgr

import (
	"MetaChat/pkg/account/pkg/commands"
	"MetaChat/pkg/account/pkg/user"
	"MetaChat/pkg/account/status"
)

type IUserManager interface {
	OnStart() error
	OnStop() error
	CanExecuteCommand(userID string, cmd commands.Command) status.AllowStatus
	StoreUser(user user.User, accountType user.AccountType) error
	GetAccountType(userID string) user.AccountType
	SetAccountType(userID string, accountType user.AccountType) error
}
