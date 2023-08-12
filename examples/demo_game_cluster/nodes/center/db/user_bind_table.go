package db

import (
	"fmt"
	"github.com/cherry-game/cherry/examples/demo_game_cluster/internal/guid"
	cherryTime "github.com/cherry-game/cherry/extend/time"
	cherryLogger "github.com/cherry-game/cherry/logger"
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

func BindUIDInDB(DB *gorm.DB, sdkId, pid int32, openId string) (int64, bool) {
	uid, ok := GetUIDFromOpenId(DB, pid, openId)
	if ok {
		return uid, true
	}

	userBind := &UserBindTable{
		UID:      guid.Next(),
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

func GetUID(pid int32, openId string) (int64, bool) {
	cacheKey := fmt.Sprintf(uidKey, pid, openId)

	val, found := uidCache.GetIfPresent(cacheKey)
	if found == false {
		return 0, false
	}

	return val.(int64), true
}

// BindUID 绑定UID
func BindUID(sdkId, pid int32, openId string) (int64, bool) {
	// TODO 根据 platformType的配置要求，决定查询UID的方式：
	// 条件1: platformType + openId查询，是否存在uid
	// 条件2: pid + openId查询，是否存在uid

	uid, ok := GetUID(pid, openId)
	if ok {
		return uid, true
	}

	userBind := &UserBindTable{
		UID:      guid.Next(),
		SdkId:    sdkId,
		PID:      pid,
		OpenId:   openId,
		BindTime: cherryTime.Now().ToMillisecond(),
	}

	cacheKey := fmt.Sprintf(uidKey, pid, openId)
	uidCache.Put(cacheKey, userBind.UID)

	return userBind.UID, true
}
