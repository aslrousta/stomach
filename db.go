package main

import (
	"database/sql"
	"io/fs"

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

func dbCheckObjects(objects map[string]fs.FileInfo) (
	changed []string,
	deleted []string,
	unstaged []string,
	err error,
) {
	rows, err := db.Query("SELECT name, mode, mtime FROM objects WHERE deleted = 0")
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var mode, mtime int64
		if err := rows.Scan(&name, &mode, &mtime); err != nil {
			return nil, nil, nil, err
		}
		info, exists := objects[name]
		if !exists {
			deleted = append(deleted, name)
		} else if int64(info.Mode()) != mode || info.ModTime().Unix() != mtime {
			changed = append(changed, name)
		}
		delete(objects, name)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, nil, err
	}
	for name := range objects {
		unstaged = append(unstaged, name)
	}
	return changed, deleted, unstaged, nil
}
