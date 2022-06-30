package storage

import "MetaChat/pkg/account/pkg/user"

type UserPermissionStorage interface {
	GetUser(userid string) (user.User, error)
	SetPermission(userid string, user user.User) error
}
