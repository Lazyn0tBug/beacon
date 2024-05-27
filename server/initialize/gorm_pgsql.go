package initialize

import (
	"sync"

	"context"

	"github.com/Lazyn0tBug/beacon/server/config"
	"github.com/Lazyn0tBug/beacon/server/global"
	"github.com/Lazyn0tBug/beacon/server/initialize/internal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

const PGSQLDSN = "host=localhost user=zombo password=token dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"

var (
	db   *gorm.DB
	once sync.Once
)

// GormPgSql 初始化 Postgresql 数据库
// Author [piexlmax](https://github.com/piexlmax)
// Author [SliverHorn](https://github.com/SliverHorn)
func GormPgSql() *gorm.DB {
	return GormPgSqlByConfig(global.GVA_CONFIG.Pgsql)
}

// GormPgSqlByConfig 初始化 Postgresql 数据库 通过参数
func GormPgSqlByConfig(p config.Pgsql) *gorm.DB {
	once.Do(func() {
		if p.Dbname == "" {
			return
		}
		pgsqlConfig := postgres.Config{
			DSN:                  p.Dsn(), // DSN data source name
			PreferSimpleProtocol: false,
		}
		if db, err := gorm.Open(postgres.New(pgsqlConfig), internal.Gorm.Config(p.Prefix, p.Singular)); err != nil {
			db = nil
			return
			// panic(err)
		} else {
			sqlDB, err := db.DB()
			if err != nil {
				db = nil
				return
			}
			sqlDB.SetMaxIdleConns(p.MaxIdleConns)
			sqlDB.SetMaxOpenConns(p.MaxOpenConns)
		}
	})
	return db
}

func GormPostgresInit() {
	once.Do(func() {
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
	})
}

// WriteDB ...
func WriteDB(ctx context.Context) *gorm.DB {
	return db.Clauses(dbresolver.Write).WithContext(ctx)
}

// ReadDB ...
func ReadDB(ctx context.Context) *gorm.DB {
	return db.Clauses(dbresolver.Read).WithContext(ctx)
}

// DB Read write separation
func DB(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}
