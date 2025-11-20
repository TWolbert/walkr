package sqldb

import (
	"context"
	"database/sql"
	_ "embed"

	_ "github.com/mattn/go-sqlite3"
	database "wlbt.nl/walkr/db/sqlc"
)

//go:embed schema.sql
var ddl string
var Queries *database.Queries

func Run() error {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "file:db.sqlite?_journal_mode=MEMORY")
	if err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	Queries = database.New(db)

	return nil
}
