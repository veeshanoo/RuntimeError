package repo

import (
	types "RuntimeError/types/mongo"
	"context"
)

type QuestionRepo interface {
	GetAll(ctx context.Context) ([]*types.Question, error)
	GetById(ctx context.Context, id string) (*types.Question, error)
	Insert(ctx context.Context, question *types.Question) (string, error)
	Update(ctx context.Context, oldQ *types.Question, newQ *types.Question) (*types.Question, error)
	Delete(ctx context.Context, id string) error
}
