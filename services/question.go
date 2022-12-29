package services

import (
	"RuntimeError/db/mongo"
	"RuntimeError/repo"
	types "RuntimeError/types/domain"
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

	// return svc.questionRepo.Insert(ctx, question)
	return "", nil
}

// func fromDomainToMongo(question *types.Question) mongoType.Question {
// 	// replies := mongoType.Reply {
// 	// 	Id: q,
// 	// 	Contents: "",
// 	// 	SubmitterId: "",
// 	// }
// 	// answers := mongoType.Answer {
// 	// 	Id : "",
// 	// 	Contents: "",
// 	// 	SubbmiterId: "",
// 	// 	Replies: ,
// 	// }
// 	// x := mongoType.Question{
// 	// 	Id:         question.Id,
// 	// 	SumitterId: question.SumitterId,
// 	// 	Title:      question.Title,
// 	// 	Contents:   question.Contents,
// 	// 	Answers:    question.Answers,
// 	// 	BestAnswer: question.BestAnswer,
// 	// 	Upvoters:   question.Upvoters,
// 	// 	Downvoters: question.Downvoters,
// 	// }
// 	// return x
// }
