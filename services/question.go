package services

import (
	"RuntimeError/auth"
	"RuntimeError/db/mongo"
	internalmongo "RuntimeError/db/mongo"
	"RuntimeError/repo"
	types "RuntimeError/types/domain"
	mongotypes "RuntimeError/types/mongo"
	"RuntimeError/utils"
	helper "RuntimeError/utils"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionService interface {
	CreateQuestion(ctx context.Context, question *types.Question) (string, error)
	GetQuestion(ctx context.Context, id string) (*types.Question, error)
	GetAll(ctx context.Context) ([]*types.Question, error)
	GetQuestionsForUser(ctx context.Context, userId string) ([]*types.Question, error)
	EditContent(ctx context.Context, userId string, content *types.EditContentRequest) error
	FavoriteAnswer(ctx context.Context, questionId string, answerId string, claims *auth.JWTClaims) error
	AddReplyToAnswer(ctx context.Context, req *types.AddReplyRequest, claims *auth.JWTClaims) error
	AddAnswerToQuestion(ctx context.Context, req *types.AddAnswerRequest, claims *auth.JWTClaims) error
	DownvoteQuestion(ctx context.Context, questionId, userId string) error
	UpvoteQuestion(ctx context.Context, questionId, userId string) error
}

type QuestionServiceImpl struct {
	questionRepo repo.QuestionRepo
}

func NewQuestionService() QuestionService {
	return &QuestionServiceImpl{
		questionRepo: mongo.NewQuestionRepo(),
	}
}

func (svc *QuestionServiceImpl) AddAnswerToQuestion(ctx context.Context, req *types.AddAnswerRequest, claims *auth.JWTClaims) error {
	t := internalmongo.NewTransaction(func(ctx context.Context) (interface{}, error) {
		result, err := svc.questionRepo.GetById(ctx, req.QuestionId)
		if err != nil {
			return nil, err
		}
		question := result

		ans := mongotypes.Answer{
			Id:             primitive.NewObjectID().Hex(),
			SubmitterId:    claims.UserId,
			SubmitterEmail: claims.Email,
			Contents:       req.Contents,
			Replies:        nil,
		}
		question.Answers = append(question.Answers, ans)

		_, err = svc.questionRepo.Update(ctx, result, question)
		return nil, err
	})

	if _, err := t.Execute(ctx); err != nil {
		return err
	} else {
		return nil
	}
}

func (svc *QuestionServiceImpl) AddReplyToAnswer(ctx context.Context, req *types.AddReplyRequest, claims *auth.JWTClaims) error {
	t := internalmongo.NewTransaction(func(ctx context.Context) (interface{}, error) {
		result, err := svc.questionRepo.GetById(ctx, req.QuestionId)
		if err != nil {
			return nil, err
		}
		question := result

		reply := mongotypes.Reply{
			Id:             primitive.NewObjectID().Hex(),
			SubmitterId:    claims.UserId,
			SubmitterEmail: claims.Email,
			Contents:       req.Contents,
		}

		found := false
		for idx, answer := range question.Answers {
			if answer.Id == req.AnswerId {
				question.Answers[idx].Replies = append(answer.Replies, reply)
				found = true
				break
			}
		}

		if !found {
			return "", errors.New("Answer id does not exist")
		}

		_, err = svc.questionRepo.Update(ctx, result, question)
		return nil, err
	})

	if _, err := t.Execute(ctx); err != nil {
		return err
	} else {
		return nil
	}
}

func (svc *QuestionServiceImpl) EditContent(ctx context.Context, userId string, content *types.EditContentRequest) error {
	t := internalmongo.NewTransaction(func(ctx context.Context) (interface{}, error) {
		result, err := svc.questionRepo.GetById(ctx, content.QuestionId)
		if err != nil {
			return nil, errors.New("Not found")
		}

		if result.SubmitterId != userId {
			return nil, errors.New("Unauthorized")
		}

		question := result
		question.Title = content.Title
		question.Contents = content.Content

		_, err = svc.questionRepo.Update(ctx, result, question)
		return nil, err
	})

	if _, err := t.Execute(ctx); err != nil {
		return err
	} else {
		return nil
	}
}

func (svc *QuestionServiceImpl) UpvoteQuestion(ctx context.Context, questionId, userId string) error {
	t := internalmongo.NewTransaction(func(ctx context.Context) (interface{}, error) {
		result, err := svc.questionRepo.GetById(ctx, questionId)
		if err != nil {
			return nil, err
		}
		question := result
		for _, id := range question.Upvoters {
			if id == userId {
				// already exists
				return nil, nil
			}
		}

		for idx, downvotter := range question.Downvoters {
			if downvotter == userId {
				// remove selected item
				question.Downvoters = append(question.Downvoters[:idx], question.Downvoters[idx+1:]...)
				break
			}
		}

		question.Upvoters = append(question.Upvoters, userId)

		_, err = svc.questionRepo.Update(ctx, result, question)
		return nil, err
	})

	if _, err := t.Execute(ctx); err != nil {
		return err
	} else {
		return nil
	}
}

func (svc *QuestionServiceImpl) DownvoteQuestion(ctx context.Context, questionId, userId string) error {
	t := internalmongo.NewTransaction(func(ctx context.Context) (interface{}, error) {
		result, err := svc.questionRepo.GetById(ctx, questionId)
		if err != nil {
			return nil, err
		}
		question := result
		for _, id := range question.Downvoters {
			if id == userId {
				// already exists
				return nil, nil
			}
		}

		for idx, upvoter := range question.Upvoters {
			if upvoter == userId {
				// remove selected item
				question.Upvoters = append(question.Upvoters[:idx], question.Upvoters[idx+1:]...)
				break
			}
		}
		question.Downvoters = append(question.Downvoters, userId)

		_, err = svc.questionRepo.Update(ctx, result, question)
		return nil, err
	})

	if _, err := t.Execute(ctx); err != nil {
		return err
	} else {
		return nil
	}
}

func (svc *QuestionServiceImpl) FavoriteAnswer(ctx context.Context, questionId string, answerId string, claims *auth.JWTClaims) error {
	t := internalmongo.NewTransaction(func(ctx context.Context) (interface{}, error) {
		result, err := svc.questionRepo.GetById(ctx, questionId)
		if err != nil {
			return nil, err
		}

		question := result
		if question.BestAnswer != "" {
			errors.New("Best answer already exists.")
		}
		answerExist := false
		for _, answer := range question.Answers {
			if answer.Id == answerId {
				answerExist = true
			}
		}

		if !answerExist {
			return nil, errors.New("Answer doesn't exists")
		}

		question.BestAnswer = answerId

		return svc.questionRepo.Update(ctx, result, question)
	})

	if _, err := t.Execute(ctx); err != nil {
		return err
	} else {
		return nil
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

func (svc *QuestionServiceImpl) GetQuestionsForUser(ctx context.Context, userId string) ([]*types.Question, error) {
	mongoQuestions, err := svc.questionRepo.GetAllForUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	domainQuestions := []*types.Question{}
	for _, mongoQuestion := range mongoQuestions {
		domainQuestions = append(domainQuestions, helper.MongoQuestionToDomain(mongoQuestion))
	}
	return domainQuestions, nil
}

func (svc *QuestionServiceImpl) GetQuestion(ctx context.Context, id string) (*types.Question, error) {
	if q, err := svc.questionRepo.GetById(ctx, id); err != nil {
		return nil, err
	} else {
		return utils.MongoQuestionToDomain(q), nil
	}
}
