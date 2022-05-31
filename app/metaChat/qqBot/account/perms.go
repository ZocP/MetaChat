package account

import (
	"fmt"
	"go.uber.org/zap"

	"github.com/spf13/viper"
)

type Command interface {
	GetName() string
}

type User interface {
	GetUserID() string
	GetNickName() string
	GetAccountType() string
}

type accountPermissions map[string]Permissions

type Permissions map[string]interface{}

func (perms Permissions) AllowExecuteParam(param string) bool {
	if _, ok := perms[param]; ok {
		return true
	}
	return false
}

type permsConfig struct {
	AccountTypes []struct {
		TypeName      string   `json:"typeName"`
		AllowCommands []string `json:"allowCommands"`
	} `json:"accountTypes"`
	Commands []struct {
		CommandsName string `json:"commandsName"`
		Config       string `json:"config"`
	} `json:"commands"`
}

type PermissionManager struct {
	log               *zap.Logger
	accountTypeMap    accountPermissions
	permissionStorage UserPermissionStorage
}

func (mgr *PermissionManager) CheckExecutable(cmd Command, user User) bool {
	if perms, ok := mgr.accountTypeMap[user.GetAccountType()]; ok {
		return perms.AllowExecuteParam(cmd.GetName())
	}
	return false
}

func (mgr *PermissionManager) CheckExecutableByUserID(cmd Command, userid string) bool {
	if user, err := mgr.permissionStorage.GetUser(userid); err != nil {
		mgr.log.Error("get user failed", zap.Error(err))
		return false
	} else if user == nil {
		return false
	} else {
		return mgr.CheckExecutable(cmd, user)
	}
}

func NewPermissionManager(permsConfigFilepath string, permsStorage UserPermissionStorage) (*PermissionManager, error) {
	typeMap, err := parsePerm(permsConfigFilepath)
	if err != nil {
		return nil, err
	}
	return &PermissionManager{
		accountTypeMap:    typeMap,
		permissionStorage: permsStorage,
	}, nil
}

func parsePerm(permsConfigFilepath string) (accountPermissions, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(permsConfigFilepath)

	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("读取config 失败 %w", err)
	}

	permsCfg := permsConfig{}
	if err := v.Unmarshal(&permsCfg); err != nil {
		return nil, err
	}

	commandMap := make(map[string]interface{})
	for _, command := range permsCfg.Commands {
		commandMap[command.CommandsName] = command.Config
	}

	typeMap := make(accountPermissions)
	for _, accountInfo := range permsCfg.AccountTypes {
		perms := make(Permissions)
		for _, cmd := range accountInfo.AllowCommands {
			if cfg, ok := commandMap[cmd]; ok {
				perms[cmd] = cfg
			}
		}

		typeMap[accountInfo.TypeName] = perms
	}

	return typeMap, nil
}
