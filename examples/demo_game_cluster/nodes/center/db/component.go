package db

import (
	"context"
	cherryGORM "github.com/cherry-game/cherry/components/gorm"
	cherryUtils "github.com/cherry-game/cherry/extend/utils"
	cherryFacade "github.com/cherry-game/cherry/facade"
	cherryLogger "github.com/cherry-game/cherry/logger"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"os"
)

var (
	onLoadFuncList []func() // db初始化时加载函数列表
)

type Component struct {
	cherryFacade.Component
	DB  *gorm.DB
	RDB *redis.Client
}

func (c *Component) Name() string {
	return "db_center_component"
}

// Init 组件初始化函数
// 为了简化部署的复杂性，本示例取消了数据库连接相关的逻辑
func (c *Component) Init() {
	cherryLogger.Infof("center-db-component 组件Init")

	c.getDb()

	// 连接redis
	c.RDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()

	_, err := c.RDB.Ping(ctx).Result()
	if err != nil {
		cherryLogger.Warn("redis 连接错误")
	}
}

func (c *Component) OnAfterInit() {
	//	addOnload(loadDevAccount)
	addOnload(initGuid)

	for _, fn := range onLoadFuncList {
		cherryUtils.Try(fn, func(errString string) {
			cherryLogger.Warnf(errString)
		})
	}
}

func (*Component) OnStop() {
	//组件停止时触发逻辑
}

func New() *Component {
	return &Component{}
}

// getDb 获取db指针
func (p *Component) getDb() {
	// 获取gorm组件
	orm := p.App().Find(cherryGORM.Name).(*cherryGORM.Component)
	if orm == nil {
		cherryLogger.DPanicf("[component = %s] not found.", cherryGORM.Name)
	}

	// 获取 db_id = "center_db_1" 的配置
	centerDbID := p.App().Settings().GetConfig("db_id_list").GetString("center_db_id")
	p.DB = orm.GetDb(centerDbID)
	if p.DB == nil {
		cherryLogger.Panic("center_db_1 not found")
	}

	err := p.DB.AutoMigrate(

		DevAccountTable{},
		UserBindTable{},
	)
	if err != nil {
		cherryLogger.Warn("register table failed")
		os.Exit(0)
	}
	cherryLogger.Info("register table success")
}

func addOnload(fn func()) {
	onLoadFuncList = append(onLoadFuncList, fn)
}
