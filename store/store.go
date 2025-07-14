package store

import (
	"database/sql"
	"embed"
	"errors"
	"log/slog"
	"path/filepath"

	"github.com/lakerszhy/rssx/rss"
	"github.com/pressly/goose/v3"
)

//go:embed migration/*.sql
var migrations embed.FS

var errFeedExist = errors.New("feed already exist")

type Store struct {
	db     *sql.DB
	logger *slog.Logger
}

func New(dir string, logger *slog.Logger) (*Store, error) {
	name := filepath.Join(dir, "rssx.db")
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	goose.SetBaseFS(migrations)
	if err = goose.SetDialect("sqlite3"); err != nil {
		return nil, err
	}

	if err = goose.Up(db, "migration"); err != nil {
		return nil, err
	}

	return &Store{
		db:     db,
		logger: logger,
	}, nil
}

func (s *Store) AddFeed(f rss.Feed) (rss.Feed, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return f, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			s.logger.Error("rollback add feed failed", "error", err)
		}
	}()

	f, err = s.addFeed(tx, f)
	if err != nil {
		return f, err
	}

	itemSQL := `INSERT INTO item (feed_id, title, description, content, link, published_at) VALUES (?, ?, ?, ?, ?, ?);`
	itemSTMT, err := tx.Prepare(itemSQL)
	if err != nil {
		return f, err
	}
	defer itemSTMT.Close()

	for i := range f.Items {
		item := f.Items[i]
		ret, err := itemSTMT.Exec(f.ID, item.Title, item.Description, item.Content, item.Link, item.PublishedAt.Unix())
		if err != nil {
			return f, err
		}
		f.Items[i].ID, err = ret.LastInsertId()
		if err != nil {
			return f, err
		}
	}

	return f, tx.Commit()
}

func (s *Store) AddFeeds(feeds []rss.Feed) ([]rss.Feed, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			s.logger.Error("rollback add feeds failed", "error", err)
		}
	}()

	inserted := []rss.Feed{}
	for _, i := range feeds {
		f, err := s.addFeed(tx, i)
		if err != nil {
			if errors.Is(err, errFeedExist) {
				continue
			}
			return nil, err
		}
		inserted = append(inserted, f)
	}

	return inserted, tx.Commit()
}

func (s *Store) addFeed(tx *sql.Tx, f rss.Feed) (rss.Feed, error) {
	exist, err := s.isFeedExist(tx, f.FeedURL)
	if err != nil {
		return f, err
	}
	if exist {
		return f, errFeedExist
	}

	feedSQL := `INSERT INTO feed (name, feed_url, home_page_url) VALUES (?, ?, ?);`
	ret, err := tx.Exec(feedSQL, f.Name, f.FeedURL, f.HomePageURL)
	if err != nil {
		return f, err
	}

	f.ID, err = ret.LastInsertId()
	if err != nil {
		return f, err
	}

	return f, err
}

func (s *Store) InsertItems(feedID int64, items []rss.FeedItem) ([]rss.FeedItem, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			s.logger.Error("rollback insert items failed", "error", err)
		}
	}()

	itemSQL := `INSERT INTO item (feed_id, title, description, content, link, published_at) VALUES (?, ?, ?, ?, ?, ?);`
	itemSTMT, err := tx.Prepare(itemSQL)
	if err != nil {
		return nil, err
	}
	defer itemSTMT.Close()

	for i := range items {
		item := items[i]
		var ret sql.Result
		ret, err = itemSTMT.Exec(feedID, item.Title, item.Description,
			item.Content, item.Link, item.PublishedAt.Unix())
		if err != nil {
			return nil, err
		}
		items[i].ID, err = ret.LastInsertId()
		if err != nil {
			return nil, err
		}
	}

	return items, tx.Commit()
}

func (s *Store) DeleteFeed(id int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err = tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			s.logger.Error("rollback delete feed failed",
				"id", id, "error", err)
		}
	}()

	_, err = tx.Exec(`DELETE FROM feed WHERE id = ?;`, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM item WHERE feed_id = ?;`, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) GetAllFeeds() ([]rss.Feed, error) {
	feedSQL := `SELECT id, name, feed_url, home_page_url FROM feed;`
	feedRows, err := s.db.Query(feedSQL)
	if err != nil {
		return nil, err
	}
	defer feedRows.Close()

	var feeds []feed
	for feedRows.Next() {
		var f feed
		if err = feedRows.Scan(&f.id, &f.name, &f.feedURL, &f.homePageURL); err != nil {
			return nil, err
		}
		feeds = append(feeds, f)
	}

	if err = feedRows.Err(); err != nil {
		return nil, err
	}

	itemSQL := `SELECT id, feed_id, title, description, content, link, is_read, is_starred, published_at FROM item`
	itemSTMT, err := s.db.Prepare(itemSQL)
	if err != nil {
		return nil, err
	}
	defer itemSTMT.Close()

	itemRows, err := itemSTMT.Query()
	if err != nil {
		return nil, err
	}
	defer itemRows.Close()

	var items []feedItem
	for itemRows.Next() {
		var i feedItem
		if err = itemRows.Scan(&i.id, &i.feedID, &i.title, &i.description, &i.content, &i.link, &i.isRead, &i.isStarred, &i.publishedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	if err = itemRows.Err(); err != nil {
		return nil, err
	}

	rssFeeds := make([]rss.Feed, 0, len(feeds))
	for i := range feeds {
		f := feeds[i].toFeed()
		for _, j := range items {
			if j.feedID == f.ID {
				f.Items = append(f.Items, j.toItem())
			}
		}
		rssFeeds = append(rssFeeds, f)
	}

	return rssFeeds, nil
}

func (s *Store) ToogleRead(id int64) error {
	itemSQL := `UPDATE item SET is_read = NOT is_read WHERE id = ?;`
	_, err := s.db.Exec(itemSQL, id)
	return err
}

func (s *Store) ToogleStarred(id int64) error {
	itemSQL := `UPDATE item SET is_starred = NOT is_starred WHERE id = ?;`
	_, err := s.db.Exec(itemSQL, id)
	return err
}

func (s *Store) RenameFeed(id int64, name string) error {
	feedSQL := `UPDATE feed SET name = ? WHERE id = ?;`
	_, err := s.db.Exec(feedSQL, name, id)
	return err
}

func (s *Store) isFeedExist(tx *sql.Tx, feedURL string) (bool, error) {
	feedSQL := `SELECT 1 FROM feed WHERE feed_url = ?;`
	row, err := tx.Query(feedSQL, feedURL)
	if err != nil {
		return false, err
	}
	defer row.Close()

	if err = row.Err(); err != nil {
		return false, err
	}

	return row.Next(), nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
