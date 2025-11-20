package models

import (
	"context"

	sqldb "wlbt.nl/walkr/db"
	database "wlbt.nl/walkr/db/sqlc"
)

func CreateUserRole(ctx context.Context, userId int64, role string) (*database.UserRole, error) {
	if roleId, err := sqldb.Queries.GetRoleByName(ctx, role); err != nil {
		return nil, err
	} else {
		if role, err := sqldb.Queries.CreateUserRole(ctx, database.CreateUserRoleParams{
			UserID: userId,
			RoleID: roleId,
		}); err != nil {
			return nil, err
		} else {
			return &role, nil
		}
	}
}
