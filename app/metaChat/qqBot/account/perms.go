package account

import (
	"fmt"

	"github.com/spf13/viper"
)

type Command interface {
	GetName() string
}

type AdminUser interface {
	GetUserID() string
	GetNickName() string
	GetAccountType() string
}

type accountTypeMap map[string]Perms

type Perms map[string]interface{}

func (perms Perms) AllowCall(cmd string) bool {
	if _, ok := perms[cmd]; ok {
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

type PermMgr struct {
	accountTypeMap accountTypeMap
	permsStorage   AdminPermsStorage
}

func (mgr *PermMgr) Check(cmd Command, admin AdminUser) bool {
	if perms, ok := mgr.accountTypeMap[admin.GetAccountType()]; ok {
		return perms.AllowCall(cmd.GetName())
	}

	return false
}

func (mgr *PermMgr) CheckByUserID(cmd Command, userid string) bool {
	if adminUser, err := mgr.permsStorage.ReadPerms(userid); err != nil {
		// TODO: logger err
		return false
	} else if adminUser == nil {
		return false
	} else {
		return mgr.Check(cmd, adminUser)
	}
}

func NewPermMgr(permsConfigFilepath string, permsStorage AdminPermsStorage) (*PermMgr, error) {
	typeMap, err := parsePerm(permsConfigFilepath)
	if err != nil {
		return nil, err
	}

	return &PermMgr{
		accountTypeMap: typeMap,
		permsStorage:   permsStorage,
	}, nil
}

func parsePerm(permsConfigFilepath string) (accountTypeMap, error) {
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

	typeMap := make(accountTypeMap)
	for _, accountInfo := range permsCfg.AccountTypes {
		perms := make(Perms)
		for _, cmd := range accountInfo.AllowCommands {
			if cfg, ok := commandMap[cmd]; ok {
				perms[cmd] = cfg
			}
		}

		typeMap[accountInfo.TypeName] = perms
	}

	return typeMap, nil
}
