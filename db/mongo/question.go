package mongo

import (
	types "RuntimeError/types/mongo"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const questionsCollectionName = "questions"

// implements QuestionRepo
type QuestionRepoImpl struct{}

func NewQuestionRepo() *QuestionRepoImpl {
	return &QuestionRepoImpl{}
}

func (q *QuestionRepoImpl) Insert(ctx context.Context, question *types.Question) (string, error) {
	question.BestAnswer = ""
	question.Answers = nil
	question.Upvoters = nil
	question.Downvoters = nil
	question.Id = primitive.NewObjectID().Hex()
	if _, err := Insert(ctx, questionsCollectionName, question, QuestionLabel); err != nil {
		return "", err
	}

	return question.Id, nil
}

func (q *QuestionRepoImpl) Update(ctx context.Context, oldQuestion *types.Question, newQuestion *types.Question) (string, error) {
	return Update(ctx, questionsCollectionName, bson.M{"id": oldQuestion.Id}, newQuestion)
}

func (q *QuestionRepoImpl) Delete(ctx context.Context, id string) error {
	return nil
}

func (q *QuestionRepoImpl) GetAll(ctx context.Context) ([]*types.Question, error) {
	result, err := GetAll(ctx, questionsCollectionName, bson.M{}, QuestionLabel)
	if err != nil {
		return nil, err
	}

	return result.([]*types.Question), nil
}

func (q *QuestionRepoImpl) GetAllForUser(ctx context.Context, userId string) ([]*types.Question, error) {
	result, err := GetAll(ctx, questionsCollectionName, bson.M{"submitterid": userId}, QuestionLabel)
	if err != nil {
		return nil, err
	}

	return result.([]*types.Question), nil
}

func (q *QuestionRepoImpl) GetById(ctx context.Context, id string) (*types.Question, error) {
	result, err := GetOne(ctx, questionsCollectionName, bson.M{"id": id}, QuestionLabel)
	if err != nil {
		return nil, err
	}

	return result.(*types.Question), nil
}
