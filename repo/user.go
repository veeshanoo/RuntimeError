package repo

import (
	"RuntimeError/types/mongo"
	"context"
)

type UserRepo interface {
	GetAll(ctx context.Context) ([]*types.User, error)
	GetById(ctx context.Context, id string) (*types.User, error)
	GetByEmail(ctx context.Context, email string) (*types.User, error)
	Insert(ctx context.Context, user *types.User) (string, error)
	Update(ctx context.Context, oldUser *types.User, newUser *types.User) (*types.User, error)
	Delete(ctx context.Context, id string) error
}
