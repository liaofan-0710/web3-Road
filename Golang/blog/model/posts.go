package model

import "gorm.io/gorm"

// Post 文章
type Post struct {
	gorm.Model
	Title   string `gorm:"type:varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null"`
	Content string `gorm:"type:varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null"`
	UserID  uint
	User    User
}
