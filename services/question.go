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
	GetAll(ctx context.Context) ([]*types.Question, error)
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
}

func (svc *QuestionServiceImpl) GetAll(ctx context.Context) ([]*types.Question, error) {

	mongoQuestions, err := svc.questionRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	domainQuestions := []*types.Question{}
	for _, mongoQuestion := range mongoQuestions {
		domainQuestions = append(domainQuestions, helper.MongoQuestionToDomain(mongoQuestion))
	}
	return domainQuestions, nil
}
