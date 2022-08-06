package userMgr

import (
	"MetaChat/pkg/account/pkg/user"
	"MetaChat/pkg/qqbot_framework/command"
)

type IUserManager interface {
	CheckExecutable(userID string, cmd command.Command) (bool, error)
	SetPermissionLevel(userID string, tp user.AccountType) error
}
