CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    password BLOB NOT NULL,
    date_joined TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(username)
);

CREATE TABLE IF NOT EXISTS bookmarks(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    url TEXT NOT NULL,
    title TEXT NOT NULL,
    owner_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(url, owner_id),
    FOREIGN KEY (owner_id) REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS bookmarks_owner_id_FK ON bookmarks (owner_id);

CREATE TABLE IF NOT EXISTS tags(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    owner_id INTEGER NOT NULL,
    UNIQUE(name, owner_id),
    FOREIGN KEY (owner_id) REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS tags_owner_id_FK ON tags (owner_id);

CREATE TABLE IF NOT EXISTS bookmarks_tags(
    bookmark_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY(bookmark_id, tag_id),
    FOREIGN KEY (bookmark_id) REFERENCES bookmarks (id),
    FOREIGN KEY (tag_id) REFERENCES tags (id)
);

CREATE INDEX IF NOT EXISTS bookmarks_tags_bookmark_id_FK ON bookmarks_tags (bookmark_id);

CREATE INDEX IF NOT EXISTS bookmarks_tags_tag_id_FK ON bookmarks_tags (tag_id);
