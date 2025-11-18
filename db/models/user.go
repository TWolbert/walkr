package models

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
	sqldb "wlbt.nl/walkr/db"
	database "wlbt.nl/walkr/db/sqlc"
)

func CreateUser(ctx context.Context, username string, email, password string) (*database.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if user, err := GetUserByName(ctx, username); user != nil || err != nil {
		return nil, errors.New("Username already taken")
	}

	if user, err := GetUserByEmail(ctx, email); user != nil || err != nil {
		return nil, errors.New("Email already taken")
	}

	if err != nil {
		return nil, err
	}

	dbUser, err := sqldb.Queries.CreateUser(ctx, database.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: string(passwordHash),
	})

	if err != nil {
		return nil, err
	}

	return &dbUser, nil
}

func GetUserByName(ctx context.Context, username string) (*database.User, error) {
	user, err := sqldb.Queries.GetUserByUsername(ctx, username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(ctx context.Context, email string) (*database.User, error) {
	user, err := sqldb.Queries.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
