package services

import (
	"RuntimeError/db/mongo"
	internalmongo "RuntimeError/db/mongo"
	"RuntimeError/repo"
	types "RuntimeError/types/domain"
	"RuntimeError/utils"
	"context"
	"errors"
)

type SuggestionsService interface {
	CreateSuggestion(ctx context.Context, s *types.EditSuggestion) (string, error)
	GetIncomingSuggestions(ctx context.Context, userId string) ([]*types.EditSuggestion, error)
	GetOutgoingSuggestions(ctx context.Context, userId string) ([]*types.EditSuggestion, error)
	ApproveSuggestion(ctx context.Context, userId, suId string) error
	RejectSuggestion(ctx context.Context, userId, suId string) error
}

type SuggestionsServiceImpl struct {
	suggestionRepo repo.SuggestionRepo
	questionRepo   repo.QuestionRepo
}

func NewSuggestionsService() SuggestionsService {
	return &SuggestionsServiceImpl{
		suggestionRepo: mongo.NewSuggestionRepo(),
		questionRepo:   mongo.NewQuestionRepo(),
	}
}

func (svc *SuggestionsServiceImpl) CreateSuggestion(ctx context.Context, s *types.EditSuggestion) (string, error) {
	t := internalmongo.NewTransaction(func(ctx context.Context) (interface{}, error) {
		// check for questionid
		q, err := svc.questionRepo.GetById(ctx, s.QuestionId)
		if err != nil {
			return nil, err
		}

		s.ApproverId = q.SubmitterId
		return svc.suggestionRepo.Insert(ctx, utils.DomainSToMongoS(s))
	})

	if id, err := t.Execute(ctx); err != nil {
		return "", err
	} else {
		return id.(string), nil
	}
}

func (svc *SuggestionsServiceImpl) GetIncomingSuggestions(ctx context.Context, userId string) ([]*types.EditSuggestion, error) {
	sug, err := svc.suggestionRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var filteredSug []*types.EditSuggestion
	for _, v := range sug {
		if v.ApproverId == userId {
			filteredSug = append(filteredSug, utils.MongoSToDomainS(v))
		}
	}

	return filteredSug, nil
}

func (svc *SuggestionsServiceImpl) GetOutgoingSuggestions(ctx context.Context, userId string) ([]*types.EditSuggestion, error) {
	sug, err := svc.suggestionRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var filteredSug []*types.EditSuggestion
	for _, v := range sug {
		if v.SubmitterId == userId {
			filteredSug = append(filteredSug, utils.MongoSToDomainS(v))
		}
	}

	return filteredSug, nil
}

func (svc *SuggestionsServiceImpl) ApproveSuggestion(ctx context.Context, userId, suId string) error {
	t := internalmongo.NewTransaction(func(ctx context.Context) (interface{}, error) {
		su, err := svc.suggestionRepo.GetById(ctx, suId)
		if err != nil {
			return nil, err
		}

		if su.EditStatus != "pending" {
			return nil, errors.New("Not permitted")
		}

		if su.ApproverId != userId {
			return nil, errors.New("Unauthorized")
		}

		q, err := svc.questionRepo.GetById(ctx, su.QuestionId)
		if err != nil {
			return nil, errors.New("Not found")
		}

		q.Contents = su.Contents
		_, err = svc.questionRepo.Update(ctx, q, q)

		if err != nil {
			return nil, err
		}

		su.EditStatus = "approved"
		return svc.suggestionRepo.Update(ctx, su, su)
	})

	if _, err := t.Execute(ctx); err != nil {
		return err
	} else {
		return nil
	}
}

func (svc *SuggestionsServiceImpl) RejectSuggestion(ctx context.Context, userId, suId string) error {
	t := internalmongo.NewTransaction(func(ctx context.Context) (interface{}, error) {
		su, err := svc.suggestionRepo.GetById(ctx, suId)
		if err != nil {
			return nil, err
		}

		if su.EditStatus != "pending" {
			return nil, errors.New("Not permitted")
		}

		if su.ApproverId != userId {
			return nil, errors.New("Unauthorized")
		}

		su.EditStatus = "rejected"
		return svc.suggestionRepo.Update(ctx, su, su)
	})

	if _, err := t.Execute(ctx); err != nil {
		return err
	} else {
		return nil
	}
}
