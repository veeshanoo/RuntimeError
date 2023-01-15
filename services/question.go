package services

import (
	"RuntimeError/db/mongo"
	"RuntimeError/repo"
	types "RuntimeError/types/domain"
	mongotypes "RuntimeError/types/mongo"
	helper "RuntimeError/utils"
	"context"
)

type QuestionService interface {
	CreateQuestion(ctx context.Context, question *types.Question) (string, error)
	GetAll(ctx context.Context) ([]*types.Question, error)
	EditConent(ctx context.Context, id string, newContent string) (string, error)
	FavoriteComment(ctx context.Context, id string, bestAnswer string) (string, error)
	AddReplyToAnswer(ctx context.Context, id string, answerId string, reply *types.Reply) (string, error)
	AddAnswerToQuestion(ctx context.Context, id string, answer *types.Answer) (string, error)
	DownvoteQuestion(ctx context.Context, id string, downvotterId string) (string, error)
	UpvoteQuestion(ctx context.Context, id string, upvotterId string) (string, error)
}

type QuestionServiceImpl struct {
	questionRepo repo.QuestionRepo
}

func NewQuestionService() QuestionService {
	return &QuestionServiceImpl{
		questionRepo: mongo.NewQuestionRepo(),
	}
}

func (svc *QuestionServiceImpl) AddAnswerToQuestion(ctx context.Context, id string, answer *types.Answer) (string, error) {
	return svc.questionRepo.AddAnswerToQuestion(ctx, id, &mongotypes.Answer{
		Id:          answer.Id,
		SubmitterId: answer.SubmitterId,
		Contents:    answer.Contents,
		Replies:     helper.DomainReplyToMongo(answer.Replies),
	})
}

func (svc *QuestionServiceImpl) AddReplyToAnswer(ctx context.Context, id string, answerId string, reply *types.Reply) (string, error) {
	return svc.questionRepo.AddReplyToAnswer(ctx, id, answerId, &mongotypes.Reply{
		Id:          reply.Id,
		SubmitterId: reply.SubmitterId,
		Contents:    reply.Contents,
	})
}

func (svc *QuestionServiceImpl) EditConent(ctx context.Context, id string, newContent string) (string, error) {
	return svc.questionRepo.EditContent(ctx, id, newContent)
}

func (svc *QuestionServiceImpl) DownvoteQuestion(ctx context.Context, id string, downvotterId string) (string, error) {
	return svc.questionRepo.DownvoteQuestion(ctx, id, downvotterId)
}

func (svc *QuestionServiceImpl) UpvoteQuestion(ctx context.Context, id string, upvotterId string) (string, error) {
	return svc.questionRepo.UpvoteQuestion(ctx, id, upvotterId)
}

func (svc *QuestionServiceImpl) FavoriteComment(ctx context.Context, id string, bestAnswer string) (string, error) {
	return svc.questionRepo.FavoriteComment(ctx, id, bestAnswer)
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
