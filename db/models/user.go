package models

import (
	"context"

	sql "wlbt.nl/walkr/db"
	database "wlbt.nl/walkr/db/sqlc"
)

func CreateUser(ctx context.Context, username, email, password string) (*database.User, error) {
	dbUser, err := sql.Queries.CreateUser(ctx, database.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: password,
	})
	if err != nil {
		return nil, err
	}

	return &dbUser, nil
}
