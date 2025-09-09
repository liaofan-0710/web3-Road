package model

import "gorm.io/gorm"

// Comment 评论
type Comment struct {
	gorm.Model
	Content string `gorm:"type:varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}
