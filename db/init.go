package sqldb

import (
	"context"
	"database/sql"
	_ "embed"

	_ "modernc.org/sqlite"
	database "wlbt.nl/walkr/db/sqlc"
)

//go:embed schema.sql
var ddl string
var Queries *database.Queries

func Run() error {
	ctx := context.Background()

	db, err := sql.Open("sqlite", "db.sqlite")
	if err != nil {
		return err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	Queries = database.New(db)

	return nil
}
