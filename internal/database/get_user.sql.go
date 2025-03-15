// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: get_user.sql

package database

import (
	"context"
)

const getOneUser = `-- name: GetOneUser :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red
FROM users
WHERE email=$1
`

func (q *Queries) GetOneUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getOneUser, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPassword,
		&i.IsChirpyRed,
	)
	return i, err
}
