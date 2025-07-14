-- +goose Up
CREATE TABLE feed (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  feed_url TEXT NOT NULL,
  home_page_url TEXT NOT NULL
);

CREATE TABLE item (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  feed_id INTEGER NOT NULL,
  title TEXT NOT NULL,
  description TEXT,
  content TEXT,
  link TEXT NOT NULL,
  is_read BOOLEAN NOT NULL DEFAULT FALSE,
  is_starred BOOLEAN NOT NULL DEFAULT FALSE,
  published_at INTEGER
);

-- +goose Down
DROP TABLE feed;
DROP TABLE item;
