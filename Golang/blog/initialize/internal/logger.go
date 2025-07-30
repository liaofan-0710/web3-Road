package internal

import (
	"Project/global"
	"fmt"
	"gorm.io/gorm/logger"
)

type writer struct {
	logger.Writer
}

// NewWriter writer 构造函数
func NewWriter(w logger.Writer) *writer {
	return &writer{Writer: w}
}

// Printf 格式化打印日志
func (w *writer) Printf(message string, data ...interface{}) {
	var logZap bool
	switch global.BG_CONFIG.System.DbType {
	case "mysql":
		logZap = global.BG_CONFIG.Mysql.LogZap
	case "pgsql":
		logZap = global.BG_CONFIG.Pgsql.LogZap
	}
	if logZap {
		global.BG_LOG.Info(fmt.Sprintf(message+"\n", data...))
	} else {
		w.Writer.Printf(message, data...)
	}
}
