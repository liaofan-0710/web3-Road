package main

import (
	"Project/core"
	"Project/global"
	"Project/initialize"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

func main() {
	global.BG_VP = core.Viper()
	global.BG_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.BG_LOG)
	global.BG_DB = initialize.Gorm() // gorm连接数据库
	initialize.Timer()
	initialize.DBList()
	if global.BG_DB != nil {
		//initialize.RegisterTables(global.BG_DB) // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.BG_DB.DB()
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				fmt.Println("close db fail")
			}
		}(db)
	}
	core.RunServer()
}
