package initialize

import (
	"Project/global"
	"Project/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
)

// Gorm 初始化数据库并产生数据库全局变量
func Gorm() *gorm.DB {
	switch global.BG_CONFIG.System.DbType {
	case "mysql":
		return GormMysql()
	case "pgsql":
		return GormPgSql()
	default:
		return GormMysql()
	}
}

// RegisterTables 注册数据库表专用
func RegisterTables(db *gorm.DB) {
	err := db.AutoMigrate(
		model.User{},
		model.Post{},
		model.Comment{},
		model.JwtBlacklist{},
	)
	if err != nil {
		global.BG_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.BG_LOG.Info("register table success")
}
