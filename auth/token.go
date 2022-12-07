package auth

import (
	"RuntimeError/types"
	"context"
)

type TokenProvider interface {
	Create(ctx context.Context, user *types.User) (string, error)
	Verify(ctx context.Context, token string) (interface{}, error)
}