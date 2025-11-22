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
	if user, err := GetUserByName(ctx, username); user != nil || err != nil {
		return nil, errors.New("Username already taken")
	}

	if user, err := GetUserByEmail(ctx, email); user != nil || err != nil {
		return nil, errors.New("Email already taken")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return nil, err
	}

	if dbUser, err := sqldb.Queries.CreateUser(ctx, database.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: string(passwordHash),
	}); err != nil {
		return nil, err
	} else {
		if _, err := CreateUserRole(ctx, dbUser.ID, "user"); err != nil {
			return nil, err
		}
		return &dbUser, nil
	}
}

func GetUserByName(ctx context.Context, username string) (*database.User, error) {
	if user, err := sqldb.Queries.GetUserByUsername(ctx, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	} else {
		return &user, nil
	}
}

func GetUserByEmail(ctx context.Context, email string) (*database.User, error) {
	if user, err := sqldb.Queries.GetUserByEmail(ctx, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	} else {
		return &user, nil
	}
}

func GetUserById(ctx context.Context, id int64) (*database.User, error) {
	if user, err := sqldb.Queries.GetUserById(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	} else {
		return &user, nil
	}
}

func UserGetRole(ctx context.Context, user *database.User) (*database.Role, error) {
	data, err := sqldb.Queries.GetUserRoleByUserId(ctx, user.ID)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func UserGetToken(ctx context.Context, user *database.User) (*database.Token, error) {
	data, err := sqldb.Queries.UserGetToken(ctx, user.ID)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

type UserWithRelations struct {
	User  *database.User
	Token *database.Token
	Role  *database.Role
}

func UserLoadRelation(ctx context.Context, user *database.User, relations ...string) (*UserWithRelations, error) {
	newUser := &UserWithRelations{
		User: user,
	}

	for _, relation := range relations {
		switch relation {
		case "role":
			data, err := UserGetRole(ctx, user)
			if err != nil {
				return nil, err
			}

			if data != nil {
				newUser.Role = data
			}
		case "token":
			data, err := UserGetToken(ctx, user)

			if err != nil {
				return nil, err
			}

			if data != nil {
				newUser.Token = data
			}
		}
	}

	return newUser, nil
}
