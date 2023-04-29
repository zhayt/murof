package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zhayt/clean-arch-tmp-forum/config"
	"time"
)

const (
	userTable = `CREATE TABLE IF NOT EXISTS user(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			password TEXT UNIQUE,
			login TEXT UNIQUE,
			username TEXT UNIQUE,
			token TEXT DEFAULT NULL,
			tokenduration DATETIME DEFAULT NULL
			);`
	postTable = `CREATE TABLE IF NOT EXISTS post(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			author TEXT,
			title TEXT,
			description TEXT,
			like INTEGER DEFAULT 0,
			dislike INTEGER DEFAULT 0,
			date TEXT,
			FOREIGN KEY(user_id) REFERENCES user(id)
			);`

	categoryTable = `CREATE TABLE IF NOT EXISTS category(
    	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    	title TEXT NOT NULL
    );
		INSERT INTO category(title)
		VALUES ("IT"), ("Sport"), ("Education"), ("News"), ("Health");
`
	postCategoryTable = `CREATE TABLE IF NOT EXISTS post_category(
    		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    		post_id INTEGER NOT NULL,
    		category_id INTEGER NOT NULL,
    		FOREIGN KEY(post_id) REFERENCES post(id)
    		FOREIGN KEY(category_id) REFERENCES category(id)
	);`

	commentTable = `CREATE TABLE IF NOT EXISTS comment(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			postId INTEGER,
			userId INTEGER,
			author TEXT,
			text TEXT,
			like INTEGER DEFAULT 0,
			dislike INTEGER DEFAULT 0,
			date TEXT,
			FOREIGN KEY(postId) REFERENCES post(id)
			);`
	likeTable = `CREATE TABLE IF NOT EXISTS like(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			postId INTEGER,
			userId INTEGER,
			commentId INTEGER,
			active INTEGER DEFAULT 0,
			FOREIGN KEY(postId) REFERENCES post(id)
		);`
	dislikeTable = `CREATE TABLE IF NOT EXISTS dislike(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			postId INTEGER,
			userId INTEGER,
			commentId INTEGER,
			active INTEGER DEFAULT 0,
			FOREIGN KEY(postId) REFERENCES post(id)
		);`
)

func Dial(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		return nil, fmt.Errorf("couldn't get db connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("couldn't ping db: %w", err)
	}
	return db, nil
}

func InnitDB(db *sql.DB) error {
	for _, table := range []string{userTable, postTable, commentTable, likeTable, dislikeTable, categoryTable, postCategoryTable} {
		_, err := db.Exec(table)
		if err != nil {
			return fmt.Errorf("couldn't init db: error: %w table: %s", err, table)
		}
	}
	return nil
}
