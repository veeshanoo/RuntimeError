package repo

import (
	"RuntimeError/types"
	"context"
)

type UserRepo interface {
	GetAll(ctx context.Context) ([]*types.User, error)
	GetById(ctx context.Context, id string) (*types.User, error)
	Insert(ctx context.Context, obj *types.User) (string, error)
	Update(ctx context.Context, oldObj *types.User, newObj *types.User) (*types.User, error)
	Delete(ctx context.Context, id string) error
}