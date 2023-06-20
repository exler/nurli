package core

import "time"

type User struct {
	ID         int       `db:"id" json:"id"`
	Username   string    `db:"username" json:"username"`
	Password   string    `db:"password" json:"password"`
	DateJoined time.Time `db:"date_joined" json:"date_joined"`
}

type Tag struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Bookmark struct {
	ID        int       `db:"id" json:"id"`
	URL       string    `db:"url" json:"url"`
	Title     string    `db:"title" json:"title"`
	Tags      []Tag     `json:"tags"`
	Owner     User      `json:"owner"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
