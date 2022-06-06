package storage

import "MetaChat/app/metaChat/qqBot/account/pkg/user"

type UserPermissionStorage interface {
	GetUser(userid string) (user.User, error)
	SetPermission(userid string, user user.User) error
}
