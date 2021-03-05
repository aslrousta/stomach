package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const (
	dbFileName = ".stomach"
	migration  = `BEGIN;
	CREATE TABLE IF NOT EXISTS objects (
		id      INT     PRIMARY KEY,
		name    TEXT    NOT NULL UNIQUE,
		mode    INT     NOT NULL,
		ctime   INT     NOT NULL,
		mtime   INT     NOT NULL,
		atime   INT     NOT NULL,
		content TEXT    NOT NULL,
		deleted TINYINT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS i_objects_deleted ON objects (deleted);
	CREATE TABLE IF NOT EXISTS commits (
		id       INT     PRIMARY KEY,
		ctime    INT     NOT NULL,
		hash     TEXT    NOT NULL UNIQUE,
		message  TEXT    NOT NULL,
		accepted TINYINT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS i_commits_accepted ON commits (accepted);
	CREATE TABLE IF NOT EXISTS edits (
		id       INT     PRIMARY KEY,
		o_id     INT     NOT NULL REFERENCES objects (id),
		c_id     INT     NOT NULL REFERENCES commits (id),
		chtype   TINYINT NOT NULL,
		lineno   INT     NOT NULL,
		line     TEXT
	);
	CREATE TABLE IF NOT EXISTS fschanges (
		id      INT     PRIMARY KEY,
		o_id    INT     NOT NULL REFERENCES objects (id),
		c_id    INT     NOT NULL REFERENCES commits (id),
		mode    INT     NOT NULL,
		ctime   INT     NOT NULL,
		mtime   INT     NOT NULL,
		atime   INT     NOT NULL,
		deleted TINYINT NOT NULL
	);
	COMMIT;`
)

func dbOpen() (err error) {
	db, err = sql.Open("sqlite3", dbFileName)
	if err != nil {
		return err
	}
	_, err = db.Exec(migration)
	return err
}

func dbClose() {
	if db != nil {
		db.Close()
	}
}
