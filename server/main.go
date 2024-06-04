// main.go

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Lazyn0tBug/beacon/server/global"
	"github.com/Lazyn0tBug/beacon/server/initialize"
	"github.com/Lazyn0tBug/beacon/server/service/system"
	"github.com/Lazyn0tBug/beacon/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	// "github.com/lestrrat-go/jwx"
)

func main() {
	// Create the logger based on the configuration
	Logger := utils.GetLogger()

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

	initialize.GormInit()
	global.GVA_DB = initialize.DB(context.Background())

	// 从db中加载未过期的jwt token
	if global.GVA_DB != nil {
		initialize.RegisterTables()
		system.LoadAll()
		db, _ := global.GVA_DB.DB()
		defer db.Close()
	}

	r := gin.Default()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	// Now you can use the logger
	Logger.Info("This is an info message")
	Logger.Error("This is an error message")

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
