package global

import (
	"Project/config"
	"Project/utils/timer"
	"github.com/go-redis/redis/v8"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

var (
	BG_DB                  *gorm.DB
	BG_REDIS               *redis.Client
	BG_LOG                 *zap.Logger
	BG_VP                  *viper.Viper
	BG_CONFIG              config.Server
	BG_Timer               timer.Timer = timer.NewTimerTask()
	BlackCache             local_cache.Cache
	BG_DBList              map[string]*gorm.DB
	BG_Concurrency_Control = &singleflight.Group{}
)
