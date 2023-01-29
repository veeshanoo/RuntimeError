package repo

import (
	types "RuntimeError/types/mongo"
	"context"
)

type SuggestionRepo interface {
	GetAll(ctx context.Context) ([]*types.EditSuggestion, error)
	GetById(ctx context.Context, id string) (*types.EditSuggestion, error)
	Insert(ctx context.Context, s *types.EditSuggestion) (string, error)
	Update(ctx context.Context, oldS *types.EditSuggestion, newS *types.EditSuggestion) (string, error)
}