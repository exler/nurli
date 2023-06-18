package core

type User struct {
	ID       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

type Tag struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type Bookmark struct {
	ID    int    `db:"id" json:"id"`
	URL   string `db:"url" json:"url"`
	Title string `db:"title" json:"title"`
	Tags  []Tag  `json:"tags"`
}
