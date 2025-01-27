// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: users.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createUsername = `-- name: CreateUsername :exec
INSERT INTO public.profiles (id, username)
VALUES ($1, $2)
`

type CreateUsernameParams struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func (q *Queries) CreateUsername(ctx context.Context, arg CreateUsernameParams) error {
	_, err := q.db.Exec(ctx, createUsername, arg.ID, arg.Username)
	return err
}
