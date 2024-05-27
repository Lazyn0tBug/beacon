package initialize

import (
	"context"
	"os"

	"github.com/Lazyn0tBug/beacon/server/global"
	"github.com/Lazyn0tBug/beacon/server/model"
	"github.com/Lazyn0tBug/beacon/server/model/"
	"github.com/Lazyn0tBug/beacon/server/model/system"
	"github.com/Lazyn0tBug/beacon/server/utils"
	"go.uber.org/zap"

	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	Logger := utils.GetLogger()
	var db *gorm.DB
	switch global.GVA_CONFIG.System.DbType {
	case "mysql":
		db = GormMysql()
	case "postgres":
		GormPostgresInit()
		db = DB(context.Background())
		// return GormPgSql()
	case "oracle":
		db = GormOracle()
	case "mssql":
		db = GormMssql()
	case "sqlite":
		db = GormSqlite()
	default:
		db = GormMysql()
	}

	if db == nil {
		Logger.Error("db initialized failed")
	}
	return db
}

func RegisterTables() {
	db := global.GVA_DB
	err := db.AutoMigrate(
		model.User{},
		model.Role{},
		model.Permission{},
		model.Customer{},
		model.Case{},
		// system.SysApi{},
		// system.SysUser{},
		// system.SysBaseMenu{},
		system.JwtBlacklist{},
		// system.SysAuthority{},
		// system.SysDictionary{},
		// system.SysOperationRecord{},
		// system.SysAutoCodeHistory{},
		// system.SysDictionaryDetail{},
		// system.SysBaseMenuParameter{},
		// system.SysBaseMenuBtn{},
		// system.SysAuthorityBtn{},
		// system.SysAutoCode{},
		// system.SysExportTemplate{},
		// system.Condition{},
		// system.JoinTemplate{},

		// example.ExaFile{},
		// example.ExaCustomer{},
		// example.ExaFileChunk{},
		// example.ExaFileUploadAndDownload{},
	)
	if err != nil {
		global.GVA_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.GVA_LOG.Info("register table success")
}
