package initialize

import (
	"context"
	"os"
	"sync"

	"github.com/Lazyn0tBug/beacon/server/global"
	"github.com/Lazyn0tBug/beacon/server/initialize/internal"
	"github.com/Lazyn0tBug/beacon/server/model"
	"github.com/Lazyn0tBug/beacon/server/model/system"
	"github.com/Lazyn0tBug/beacon/server/utils"
	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var (
	db     *gorm.DB
	once   sync.Once
	Logger = utils.GetLogger()
)

func GormInit() {
	once.Do(func() {
		switch global.GVA_CONFIG.System.DbType {
		case "mysql":
			db = MysqlInit()
		case "postgres":
			db = PostgresInit()
		case "sqlite":
			db = SqliteInit()
		default:
			db = MysqlInit()
		}

		if db == nil {
			Logger.Error("db initialized failed")
		}
	})
}

func WriteDB(ctx context.Context) *gorm.DB {
	return db.Clauses(dbresolver.Write).WithContext(ctx)
}

func ReadDB(ctx context.Context) *gorm.DB {
	return db.Clauses(dbresolver.Read).WithContext(ctx)
}

func DB(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}

func MysqlInit() *gorm.DB {
	m := global.GVA_CONFIG.Mysql
	if m.Dbname == "" {
		db = nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config(m.Prefix, m.Singular)); err != nil {
		return nil
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}

func PostgresInit() *gorm.DB {
	p := global.GVA_CONFIG.Pgsql
	if p.Dbname == "" {
		db = nil
	}
	pgsqlConfig := postgres.Config{
		DSN:                  p.Dsn(), // DSN data source name
		PreferSimpleProtocol: false,
	}
	if DB, err := gorm.Open(postgres.New(pgsqlConfig), internal.Gorm.Config(p.Prefix, p.Singular)); err != nil {
		db = nil
		panic(err)
	} else {
		db = DB
	}
	return db
}

// GormSqlite 初始化Sqlite数据库
func SqliteInit() *gorm.DB {
	s := global.GVA_CONFIG.Sqlite
	if s.Dbname == "" {
		return nil
	}

	if db, err := gorm.Open(sqlite.Open(s.Dsn()), internal.Gorm.Config(s.Prefix, s.Singular)); err != nil {
		panic(err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(s.MaxIdleConns)
		sqlDB.SetMaxOpenConns(s.MaxOpenConns)
		return db
	}
}

func RegisterTables() {
	db := global.GVA_WriteDB
	err := db.AutoMigrate(
		model.User{},
		model.Role{},
		model.Permission{},
		model.Customer{},
		model.Case{},
		model.Appointment{},
		model.ServiceItem{},
		model.Doctor{},
		model.Hospital{},
		model.MedicalRecord{},
		// system.SysApi{},
		// system.SysUser{},
		// system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.JwtInActive{},
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
