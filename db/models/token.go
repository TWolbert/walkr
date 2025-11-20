package models

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	sqldb "wlbt.nl/walkr/db"
	database "wlbt.nl/walkr/db/sqlc"
)

func generateSessionToken() (string, error) {
	// Generate 32 random bytes (256 bits)
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Encode to base64 for URL-safe string
	return base64.URLEncoding.EncodeToString(b), nil
}

func CreateToken(ctx context.Context, user_id int64) (*database.Token, error) {
	if err := sqldb.Queries.DeleteTokenById(ctx, user_id); err != nil {
		return nil, err
	}

	expiresAt := time.Now().AddDate(0, 0, 30)

	if token, err := generateSessionToken(); err != nil {
		return nil, err
	} else {
		if data, err := sqldb.Queries.CreateToken(ctx, database.CreateTokenParams{
			Token:     token,
			UserID:    user_id,
			ExpiresAt: expiresAt,
		}); err != nil {
			return nil, err
		} else {
			return &data, nil
		}
	}
}

func GetUserByToken(ctx context.Context, token string) (*database.User, error) {
	if data, err := sqldb.Queries.GetUserFromToken(ctx, token); err != nil {
		return nil, err
	} else {
		return &data, nil
	}
}
