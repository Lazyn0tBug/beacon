package core

import (
	"fmt"

	"github.com/Lazyn0tBug/beacon/server/global"
	"github.com/Lazyn0tBug/beacon/server/initialize"
	"github.com/Lazyn0tBug/beacon/server/service/system"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {
	if global.GVA_CONFIG.System.UseMultipoint || global.GVA_CONFIG.System.UseRedis {
		// 初始化redis链接
		initialize.InitializeRedis()
	}

	if global.GVA_CONFIG.System.UseMongo {
		// 初始化mongo链接
		err := initialize.Mongo.InitializeMongo()
		if err != nil {
			global.GVA_LOG.Error("MongoDB connection failed", zap.String("err", err.Error()))
			panic(err)
		}
	}

	// 从db加载jwt数据
	if global.GVA_DB != nil {
		system.LoadAll()
	}

	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, Router)

	global.GVA_LOG.Info("server run success on ", zap.String("address", address))

	global.GVA_LOG.Error(s.ListenAndServe().Error())
}
