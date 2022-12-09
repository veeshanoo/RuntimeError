package auth

import (
	"RuntimeError/types/domain"
	"context"
)

type TokenProvider interface {
	Create(ctx context.Context, user *types.User) (string, error)
	Verify(ctx context.Context, token string) (interface{}, error)
}