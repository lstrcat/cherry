package db

import (
	cherryGORM "github.com/cherry-game/cherry/components/gorm"
	cherryUtils "github.com/cherry-game/cherry/extend/utils"
	cherryFacade "github.com/cherry-game/cherry/facade"
	clog "github.com/cherry-game/cherry/logger"
	"gorm.io/gorm"
	"os"
)

var (
	onLoadFuncList []func() // db初始化时加载函数列表
)

type Component struct {
	cherryFacade.Component
	DB *gorm.DB
}

func (c *Component) Name() string {
	return "db_game_component"
}

// Init 组件初始化函数
// 为了简化部署的复杂性，本示例取消了数据库连接相关的逻辑
func (c *Component) Init() {
	c.getDb()
}

func (c *Component) OnAfterInit() {
	for _, fn := range onLoadFuncList {
		cherryUtils.Try(fn, func(errString string) {
			clog.Warnf(errString)
		})
	}
}

func (*Component) OnStop() {
	//组件停止时触发逻辑
}

// getDb 获取db指针
func (c *Component) getDb() {
	// 获取gorm组件
	orm := c.App().Find(cherryGORM.Name).(*cherryGORM.Component)
	if orm == nil {
		clog.DPanicf("[component = %s] not found.", cherryGORM.Name)
	}

	// 获取 db_id = "center_db_1" 的配置
	//	centerDbID := c.App().Settings().GetConfig("db_id_list").GetString("center_db_id")
	c.DB = orm.GetDb("game_db_1")
	if c.DB == nil {
		clog.Panic("game_db_1 not found")
	}

	err := c.DB.AutoMigrate(

		PlayerTable{},
	)
	if err != nil {
		clog.Warn("register table failed")
		os.Exit(0)
	}
	clog.Info("register table success")
}

func New() *Component {
	return &Component{} // register db center
}

func addOnload(fn func()) {
	onLoadFuncList = append(onLoadFuncList, fn)
}
