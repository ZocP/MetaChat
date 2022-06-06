package storage

import (
	"MetaChat/pkg/qqBot/pkg/user"
	"encoding/csv"
	"os"
)



type PermsCSVImpl struct {
	filepath     string
	userMap map[string]*user.User
}

func NewPermsCSVImpl(adminCSVFilepath string) (UserPermissionStorage, error) {
	file, err := os.Open(adminCSVFilepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	all, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	// TODO: 优化初始化 增加健壮性
	perms := make(map[string]*user.User)
	for _, row := range all {
		accountType := row[2]
		perms[row[1]] = user.NewUser(row[1],row[2].(string).(user.AccountType),row[0])
		}
	}

	return &PermsCSVImpl{
		filepath:     adminCSVFilepath,
		userPermsMap: perms,
	}, nil
}

func (csvStorage *PermsCSVImpl) GetUser(userid string) (user.User, error) {
	if adminUser, ok := csvStorage.userPermsMap[userid]; ok {
		return adminUser, nil
	}
	return nil, nil
}

func (csvStorage *PermsCSVImpl) WriteUser(userid string, u, *user.User) error {
	csvStorage.userPermsMap[userid] = u

	all := make([][]string, 0, 10)
	for _, user := range csvStorage.userPermsMap {
		all = append(all, []string{user.GetUserID(), user.GetNickName(), user.GetAccountType()})
	}

	file, err := os.Open(csvStorage.filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	cWriter := csv.NewWriter(file)
	if err = cWriter.WriteAll(all); err != nil {
		return err
	}
	cWriter.Flush()

	return nil
}
