package db

import (
	cherryTime "github.com/cherry-game/cherry/extend/time"
	cherryLogger "github.com/cherry-game/cherry/logger"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// UserBindTable uid绑定第三方平台表
type UserBindTable struct {
	UID      int64  `gorm:"column:uid;primary_key;comment:'用户唯一id'" json:"uid"`
	SdkId    int32  `gorm:"column:sdk_id;comment:'sdk id'" json:"sdkId"`
	PID      int32  `gorm:"column:pid;comment:'平台id'" json:"pid"`
	OpenId   string `gorm:"column:open_id;comment:'平台帐号open_id'" json:"openId"`
	BindTime int64  `gorm:"column:bind_time;comment:'绑定时间'" json:"bindTime"`
}

func (*UserBindTable) TableName() string {
	return "user_bind"
}

func GetUIDFromOpenId(DB *gorm.DB, pid int32, openId string) (int64, bool) {
	user := UserBindTable{}

	result := DB.Where("pid = ? AND open_id = ?", pid, openId).First(&user)

	if result.RowsAffected <= 0 {
		return 0, false
	}
	return user.UID, true
}

func BindUIDInDB(rdb *redis.Client, DB *gorm.DB, sdkId, pid int32, openId string) (int64, bool) {
	uid, ok := GetUIDFromOpenId(DB, pid, openId)
	if ok {
		return uid, true
	}

	userBind := &UserBindTable{
		SdkId:    sdkId,
		PID:      pid,
		OpenId:   openId,
		BindTime: cherryTime.Now().ToMillisecond(),
	}

	err := DB.Create(&userBind).Error
	if err != nil {
		cherryLogger.Error("Create UserBindTable failed")
	}
	return userBind.UID, true
}
