package mongo

import (
	types "RuntimeError/types/mongo"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const suggestionsCollectionName = "suggestions"

// implements SuggestionRepo
type SuggestionRepoImpl struct{}

func NewSuggestionRepo() *SuggestionRepoImpl {
	return &SuggestionRepoImpl{}
}

func (q *SuggestionRepoImpl) Insert(ctx context.Context, s *types.EditSuggestion) (string, error) {
	s.Id = primitive.NewObjectID().Hex()
	s.EditStatus = "pending"
	if _, err := Insert(ctx, suggestionsCollectionName, s, SuggestionLabel); err != nil {
		return "", err
	}

	return s.Id, nil
}

func (q *SuggestionRepoImpl) Update(ctx context.Context, oldS *types.EditSuggestion, newS *types.EditSuggestion) (string, error) {
	return Update(ctx, suggestionsCollectionName, bson.M{"id": oldS.Id}, newS)
}

func (q *SuggestionRepoImpl) GetAll(ctx context.Context) ([]*types.EditSuggestion, error) {
	result, err := GetAll(ctx, suggestionsCollectionName, bson.M{}, SuggestionLabel)
	if err != nil {
		return nil, err
	}

	return result.([]*types.EditSuggestion), nil
}

func (q *SuggestionRepoImpl) GetById(ctx context.Context, id string) (*types.EditSuggestion, error) {
	result, err := GetOne(ctx, suggestionsCollectionName, bson.M{"id": id}, SuggestionLabel)
	if err != nil {
		return nil, err
	}

	return result.(*types.EditSuggestion), nil
}
