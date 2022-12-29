package services

import (
	"RuntimeError/db/mongo"
	"RuntimeError/repo"
	types "RuntimeError/types/domain"
	helper "RuntimeError/utils"
	"context"
)

type QuestionService interface {
	CreateQuestion(ctx context.Context, question *types.Question) (string, error)
}

type QuestionServiceImpl struct {
	questionRepo repo.QuestionRepo
}

func NewQuestionService() QuestionService {
	return &QuestionServiceImpl{
		questionRepo: mongo.NewQuestionRepo(),
	}
}

func (svc *QuestionServiceImpl) CreateQuestion(ctx context.Context, question *types.Question) (string, error) {

	return svc.questionRepo.Insert(ctx, helper.DomainQuestionToMongo(question))
	// return "", nil
}
