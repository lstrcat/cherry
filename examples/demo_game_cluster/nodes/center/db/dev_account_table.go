package db

import (
	cherryError "github.com/cherry-game/cherry/error"
	"github.com/cherry-game/cherry/examples/demo_game_cluster/internal/code"
	"github.com/cherry-game/cherry/examples/demo_game_cluster/internal/guid"
	cherryString "github.com/cherry-game/cherry/extend/string"
	cherryTime "github.com/cherry-game/cherry/extend/time"
	"github.com/cherry-game/cherry/logger"
	"gorm.io/gorm"
)

// DevAccountTable 开发模式的帐号信息表(platform.TypeDevMode)
type DevAccountTable struct {
	AccountId   int64  `gorm:"column:account_id;primary_key;comment:'帐号id'" json:"accountId"`
	AccountName string `gorm:"column:account_name;comment:'帐号名'" json:"accountName"`
	Password    string `gorm:"column:password;comment:'密码'" json:"-"`
	CreateIP    string `gorm:"column:create_ip;comment:'创建ip'" json:"createIP"`
	CreateTime  int64  `gorm:"column:create_time;comment:'创建时间'" json:"createTime"`
}

func (*DevAccountTable) TableName() string {
	return "dev_account"
}

func AccountRegister(DB *gorm.DB, accountName, password, ip string) int32 {
	user := DevAccountTable{}

	result := DB.Where("account_name = ?", accountName).First(&user)
	if result.RowsAffected > 0 {
		return code.AccountNameIsExist
	}
	devAccountTable := &DevAccountTable{
		AccountId:   guid.Next(),
		AccountName: accountName,
		Password:    password,
		CreateIP:    ip,
		CreateTime:  cherryTime.Now().Unix(),
	}

	err := DB.Create(&devAccountTable).Error
	if err != nil {
		cherryLogger.Error("Save DevAccountTable failed")
		return code.Error
	}
	return code.OK
}

func AccountWithName(DB *gorm.DB, accountName string) (*DevAccountTable, error) {
	user := DevAccountTable{}

	result := DB.Where("account_name = ?", accountName).First(&user)

	if result.RowsAffected <= 0 {
		return nil, cherryError.Error("account not found")
	}
	return &user, nil
}

func DevAccountRegister(accountName, password, ip string) int32 {
	devAccount, _ := DevAccountWithName(accountName)
	if devAccount != nil {
		return code.AccountNameIsExist
	}

	devAccountTable := &DevAccountTable{
		AccountId:   guid.Next(),
		AccountName: accountName,
		Password:    password,
		CreateIP:    ip,
		CreateTime:  cherryTime.Now().Unix(),
	}

	devAccountCache.Put(accountName, devAccountTable)
	// TODO 保存db

	return code.OK
}

func DevAccountWithName(accountName string) (*DevAccountTable, error) {
	val, found := devAccountCache.GetIfPresent(accountName)
	if found == false {
		return nil, cherryError.Error("account not found")
	}

	return val.(*DevAccountTable), nil
}

func initGuid() {
	guid.InitNextId()
}

// loadDevAccount 节点启动时预加载帐号数据
func loadDevAccount() {
	// 演示用，直接手工构建几个帐号
	for i := 1; i <= 10; i++ {
		index := cherryString.ToString(i)

		devAccount := &DevAccountTable{
			AccountId:   guid.Next(),
			AccountName: "test" + index,
			Password:    "test" + index,
			CreateIP:    "127.0.0.1",
			CreateTime:  cherryTime.Now().ToMillisecond(),
		}

		devAccountCache.Put(devAccount.AccountName, devAccount)
	}

	cherryLogger.Info("preload DevAccountTable")
}
