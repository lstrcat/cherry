package ops

import (
	cherryGORM "github.com/cherry-game/cherry/components/gorm"
	"github.com/cherry-game/cherry/examples/demo_game_cluster/internal/code"
	"github.com/cherry-game/cherry/examples/demo_game_cluster/internal/pb"
	"github.com/cherry-game/cherry/examples/demo_game_cluster/nodes/center/db"
	clog "github.com/cherry-game/cherry/logger"
	cactor "github.com/cherry-game/cherry/net/actor"
	"gorm.io/gorm"
	"time"
)

var (
	pingReturn = &pb.Bool{Value: true}
)

type (
	ActorOps struct {
		cactor.Base
		centerDB *gorm.DB
	}
)

func (p *ActorOps) AliasID() string {
	return "ops"
}

// OnInit 注册remote函数
func (p *ActorOps) OnInit() {
	// 获取gorm组件
	gorm := p.App().Find(cherryGORM.Name).(*cherryGORM.Component)
	if gorm == nil {
		clog.DPanicf("[component = %s] not found.", cherryGORM.Name)
	}

	// 获取 db_id = "center_db_1" 的配置
	centerDbID := p.App().Settings().GetConfig("db_id_list").GetString("center_db_id")
	p.centerDB = gorm.GetDb(centerDbID)
	if p.centerDB == nil {
		clog.Panic("center_db_1 not found")
	}

	// 1秒后进行一次分页查询
	p.Timer().AddOnce(1*time.Second, p.selectPagination)

	p.Remote().Register("ping", p.ping)
}

// ping 请求center是否响应
func (p *ActorOps) ping() (*pb.Bool, int32) {
	return pingReturn, code.OK
}

func (p *ActorOps) selectPagination() {
	list, count := p.pagination(1, 10)
	clog.Infof("count = %d", count)

	for _, table := range list {
		clog.Infof("%+v", table)
	}
}

// pagination 分页查询
func (p *ActorOps) pagination(page, pageSize int) ([]*db.UserBindTable, int64) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	var list []*db.UserBindTable
	var count int64

	p.centerDB.Model(&db.UserBindTable{}).Count(&count)

	if count > 0 {
		list = make([]*db.UserBindTable, pageSize)
		s := p.centerDB.Limit(pageSize).Offset((page - 1) * pageSize)
		if err := s.Find(&list).Error; err != nil {
			clog.Warn(err)
		}
	}

	return list, count
}
