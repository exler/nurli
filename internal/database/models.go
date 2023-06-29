package database

import (
	"time"
)

type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex:idx_tag_name"`
}

type Bookmark struct {
	ID          uint   `gorm:"primaryKey"`
	URL         string `gorm:"uniqueIndex:idx_bookmark_url"`
	Title       string
	Description string
	Tags        []Tag `gorm:"many2many:bookmark_tags;"`
	Read        bool
	Favorite    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Session struct {
	ID        uint   `gorm:"primaryKey"`
	Token     string `gorm:"uniqueIndex:idx_session_token"`
	CreatedAt time.Time
	ExpiresAt time.Time
}
