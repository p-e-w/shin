// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2022  Philipp Emanuel Weidmann <pew@worldwidemann.com>

package main

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	_ "github.com/mattn/go-sqlite3"
)

type history struct {
	*sql.DB
}

func openHistory() (*history, error) {
	dataDirectory := filepath.Join(xdg.DataHome, "shin")
	historyFile := filepath.Join(dataDirectory, "history.db")

	// Note that the database file is only created when the first query is executed.
	db, err := sql.Open("sqlite3", historyFile)

	if err != nil {
		return nil, err
	}

	_, err = os.Stat(historyFile)

	// The history database must be opened before the input engine
	// can be initialized, so a fast path based on the existence of
	// the database file is used to minimize the amount of disk I/O
	// required in case the database has already been created.
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(dataDirectory, os.ModePerm)

		if err != nil {
			return nil, err
		}

		_, err = db.Exec(`

		CREATE TABLE history (
			command TEXT NOT NULL,
			time DATETIME NOT NULL,
			count INTEGER NOT NULL
		);

		CREATE UNIQUE INDEX command_index ON history (command);

		CREATE INDEX time_index ON history (time);

		CREATE INDEX count_index ON history (count);

		`)

		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &history{db}, nil
}

func (h *history) addCommand(command string) error {
	_, err := h.Exec(`

	INSERT INTO history (command, time, count)
	VALUES (?, CURRENT_TIMESTAMP, 1)
	ON CONFLICT (command) DO UPDATE
	SET time = CURRENT_TIMESTAMP, count = count + 1;

	`, command)

	return err
}

func (h *history) getRecentCommand(prefix string, index uint32) (string, error) {
	var command string

	// Searching for commands matching the prefix using BETWEEN
	// is much better than using LIKE because it is case sensitive,
	// doesn't require escaping, and makes use of the index.
	err := h.QueryRow(`

	SELECT command FROM history
	WHERE command BETWEEN ? AND ?
	ORDER BY time DESC, rowid
	LIMIT 1 OFFSET ?;

	`, prefix, prefix+"\u00FF", index).Scan(&command)

	return command, err
}
