package center

import (
	"github.com/cherry-game/cherry"
	cherryCron "github.com/cherry-game/cherry/components/cron"
	cherryGORM "github.com/cherry-game/cherry/components/gorm"
	"github.com/cherry-game/cherry/examples/demo_game_cluster/internal/data"
	"github.com/cherry-game/cherry/examples/demo_game_cluster/nodes/center/db"
	"github.com/cherry-game/cherry/examples/demo_game_cluster/nodes/center/module/account"
	"github.com/cherry-game/cherry/examples/demo_game_cluster/nodes/center/module/ops"
)

func Run(profileFilePath, nodeId string) {
	app := cherry.Configure(
		profileFilePath,
		nodeId,
		false,
		cherry.Cluster,
	)
	// 注册gorm组件，数据库具体配置请查看 config/profile-dev.json文件
	app.Register(cherryGORM.NewComponent())
	app.Register(cherryCron.New())
	app.Register(data.New())
	app.Register(db.New())

	app.AddActors(
		&account.ActorAccount{},
		&ops.ActorOps{},
	)

	app.Startup()
}
