// main.go

package main

import (
	"context"

	"github.com/Lazyn0tBug/beacon/server/core"
	"github.com/Lazyn0tBug/beacon/server/global"
	"github.com/Lazyn0tBug/beacon/server/initialize"
	"github.com/Lazyn0tBug/beacon/server/service/system"
	"github.com/Lazyn0tBug/beacon/server/utils"
	// "github.com/lestrrat-go/jwx"
)

func main() {
	global.GVA_VP = core.Viper()
	// Create the logger based on the configuration
	global.GVA_LOG = utils.GetLogger() // 初始化zap日志库

	// 初始化数据库
	initialize.GormInit()
	global.GVA_DB = initialize.DB(context.Background())

	// 从db中加载未过期的jwt token
	if global.GVA_DB != nil {
		initialize.RegisterTables()
		system.LoadAll()
		db, _ := global.GVA_DB.DB()
		defer db.Close()
	}

	core.RunServer()
}
