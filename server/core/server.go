package core

import (
	"fmt"

	"time"

	"github.com/Lazyn0tBug/beacon/server/global"
	"github.com/Lazyn0tBug/beacon/server/initialize"
	"github.com/Lazyn0tBug/beacon/server/service/system"
	"github.com/Lazyn0tBug/beacon/server/utils"
	"go.uber.org/zap"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

var (
	Logger = utils.GetLogger()
)

type server interface {
	ListenAndServe() error
}

func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 20 * time.Second
	s.WriteTimeout = 20 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
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
			Logger.Error("MongoDB connection failed", zap.String("err", err.Error()))
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

	Logger.Info("server run success on ", zap.String("address", address))

	Logger.Error(s.ListenAndServe().Error())
}
