package model

import "Project/global"

type JwtBlacklist struct {
	global.BG_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
