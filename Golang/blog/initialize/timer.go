package initialize

import (
	"Project/config"
	"Project/global"
	"Project/utils"
	"fmt"
	"github.com/robfig/cron/v3"
)

func Timer() {
	if global.BG_CONFIG.Timer.Start {
		for i := range global.BG_CONFIG.Timer.Detail {
			go func(detail config.Detail) {
				var option []cron.Option
				if global.BG_CONFIG.Timer.WithSeconds {
					option = append(option, cron.WithSeconds())
				}
				_, err := global.BG_Timer.AddTaskByFunc("ClearDB", global.BG_CONFIG.Timer.Spec, func() {
					err := utils.ClearTable(global.BG_DB, detail.TableName, detail.CompareField, detail.Interval)
					if err != nil {
						fmt.Println("timer error:", err)
					}
				}, option...)
				if err != nil {
					fmt.Println("add timer error:", err)
				}
			}(global.BG_CONFIG.Timer.Detail[i])
		}
	}
}
