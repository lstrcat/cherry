package game

import (
	"github.com/cherry-game/cherry"
	cherryCron "github.com/cherry-game/cherry/components/cron"
	cherryGops "github.com/cherry-game/cherry/components/gops"
	cherryGORM "github.com/cherry-game/cherry/components/gorm"
	checkCenter "github.com/cherry-game/cherry/examples/demo_game_cluster/internal/component/check_center"
	"github.com/cherry-game/cherry/examples/demo_game_cluster/internal/data"
	"github.com/cherry-game/cherry/examples/demo_game_cluster/nodes/game/db"
	"github.com/cherry-game/cherry/examples/demo_game_cluster/nodes/game/module/player"
	cherrySnowflake "github.com/cherry-game/cherry/extend/snowflake"
	cstring "github.com/cherry-game/cherry/extend/string"
	cherryUtils "github.com/cherry-game/cherry/extend/utils"
)

func Run(profileFilePath, nodeId string) {
	if cherryUtils.IsNumeric(nodeId) == false {
		panic("node parameter must is number.")
	}

	// snowflake global id
	serverId, _ := cstring.ToInt64(nodeId)
	cherrySnowflake.SetDefaultNode(serverId)

	// 配置cherry引擎
	app := cherry.Configure(profileFilePath, nodeId, false, cherry.Cluster)

	// 注册gorm组件，数据库具体配置请查看 config/profile-dev.json文件
	app.Register(cherryGORM.NewComponent())
	// diagnose
	app.Register(cherryGops.New())
	// 注册调度组件
	app.Register(cherryCron.New())
	// 注册数据配置组件
	app.Register(data.New())
	// 注册检测中心节点组件，确认中心节点启动后，再启动当前节点
	app.Register(checkCenter.New())
	// 注册db组件
	app.Register(db.New())

	app.AddActors(
		&player.ActorPlayers{},
	)

	app.Startup()
}
