package account

import (
	cherryGORM "github.com/cherry-game/cherry/components/gorm"
	"github.com/cherry-game/cherry/examples/demo_game_cluster/internal/code"
	"github.com/cherry-game/cherry/examples/demo_game_cluster/internal/pb"
	"github.com/cherry-game/cherry/examples/demo_game_cluster/nodes/center/db"
	clog "github.com/cherry-game/cherry/logger"
	cactor "github.com/cherry-game/cherry/net/actor"
	"gorm.io/gorm"
	"os"
	"strings"
)

type (
	ActorAccount struct {
		cactor.Base
		centerDB *gorm.DB
	}
)

func (p *ActorAccount) AliasID() string {
	return "account"
}

// OnInit center为后端节点，不直接与客户端通信，所以了一些remote函数，供RPC调用
func (p *ActorAccount) OnInit() {
	p.getDb()

	p.Remote().Register("registerDevAccount", p.registerDevAccount)
	p.Remote().Register("getDevAccount", p.getDevAccount)
	p.Remote().Register("getUID", p.getUID)
}

// getDb 获取db指针
func (p *ActorAccount) getDb() {
	// 获取gorm组件
	orm := p.App().Find(cherryGORM.Name).(*cherryGORM.Component)
	if orm == nil {
		clog.DPanicf("[component = %s] not found.", cherryGORM.Name)
	}

	// 获取 db_id = "center_db_1" 的配置
	centerDbID := p.App().Settings().GetConfig("db_id_list").GetString("center_db_id")
	p.centerDB = orm.GetDb(centerDbID)
	if p.centerDB == nil {
		clog.Panic("center_db_1 not found")
	}

	err := p.centerDB.AutoMigrate(

		db.DevAccountTable{},
		db.UserBindTable{},
	)
	if err != nil {
		clog.Warn("register table failed")
		os.Exit(0)
	}
	clog.Info("register table success")
}

// registerDevAccount 注册开发者帐号
func (p *ActorAccount) registerDevAccount(req *pb.DevRegister) int32 {
	accountName := req.AccountName
	password := req.Password

	if strings.TrimSpace(accountName) == "" || strings.TrimSpace(password) == "" {
		return code.LoginError
	}

	if len(accountName) < 3 || len(accountName) > 18 {
		return code.LoginError
	}

	if len(password) < 3 || len(password) > 18 {
		return code.LoginError
	}

	return db.AccountRegister(p.centerDB, accountName, password, req.Ip)
}

// getDevAccount 根据帐号名获取开发者帐号表
func (p *ActorAccount) getDevAccount(req *pb.DevRegister) (*pb.Int64, int32) {
	accountName := req.AccountName
	password := req.Password

	devAccount, _ := db.AccountWithName(p.centerDB, accountName)
	if devAccount == nil || devAccount.Password != password {
		return nil, code.AccountAuthFail
	}

	return &pb.Int64{Value: devAccount.AccountId}, code.OK
}

// getUID 获取uid
func (p *ActorAccount) getUID(req *pb.User) (*pb.Int64, int32) {
	uid, ok := db.BindUIDInDB(p.centerDB, req.SdkId, req.Pid, req.OpenId)
	if uid == 0 || ok == false {
		return nil, code.AccountBindFail
	}

	return &pb.Int64{Value: uid}, code.OK
}
