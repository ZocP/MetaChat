package account

import (
	"encoding/csv"
	"os"
)

type AdminPermsStorage interface {
	ReadPerms(userid string) (AdminUser, error)
	WritePerms(userid string, adminUser AdminUser) error
}

type PermsCSVImpl struct {
	filepath     string
	userPermsMap map[string]AdminUser
}

func NewPermsCSVImpl(adminCSVFilepath string) (AdminPermsStorage, error) {
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
	perms := make(map[string]AdminUser)
	for _, row := range all {
		perms[row[1]] = User{
			UserID:      row[1],
			Nickname:    row[0],
			AccountType: row[2],
		}
	}

	return &PermsCSVImpl{
		filepath:     adminCSVFilepath,
		userPermsMap: perms,
	}, nil
}

func (csvStorage *PermsCSVImpl) ReadPerms(userid string) (AdminUser, error) {
	if adminUser, ok := csvStorage.userPermsMap[userid]; ok {
		return adminUser, nil
	}

	return nil, nil
}

func (csvStorage *PermsCSVImpl) WritePerms(userid string, adminUser AdminUser) error {
	csvStorage.userPermsMap[userid] = adminUser

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
