package database

import (
	"time"
)

type User struct {
	ID         uint   `gorm:"primaryKey"`
	Username   string `gorm:"uniqueIndex"`
	Password   string
	DateJoined time.Time

	Bookmarks []Bookmark `gorm:"foreignKey:OwnerID;"`
	Tags      []Tag      `gorm:"foreignKey:OwnerID;"`
}

type Tag struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"uniqueIndex:idx_tag_name"`
	OwnerID int    `gorm:"uniqueIndex:idx_tag_name"`
	Owner   User
}

type Bookmark struct {
	ID          uint   `gorm:"primaryKey"`
	URL         string `gorm:"uniqueIndex:idx_bookmark_url"`
	Title       string
	Description string
	Tags        []Tag `gorm:"many2many:bookmark_tags;"`
	OwnerID     uint  `gorm:"uniqueIndex:idx_bookmark_url"`
	Owner       User
	Read        bool
	Favorite    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Session struct {
	ID        uint   `gorm:"primaryKey"`
	Token     string `gorm:"uniqueIndex:idx_session_token"`
	UserID    uint
	User      User
	CreatedAt time.Time
	ExpiresAt time.Time
}
