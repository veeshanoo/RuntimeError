package repo

import (
	types "RuntimeError/types/mongo"
	"context"
)

type QuestionRepo interface {
	GetAll(ctx context.Context) ([]*types.Question, error)
	GetById(ctx context.Context, id string) (*types.Question, error)
	Insert(ctx context.Context, question *types.Question) (string, error)
	Update(ctx context.Context, oldQ *types.Question, newQ *types.Question) (string, error)
	Delete(ctx context.Context, id string) error
	FavoriteComment(ctx context.Context, id string, bestAnswer string) (string, error)
	EditContent(ctx context.Context, id string, newContent string) (string, error)
	AddAnswerToQuestion(ctx context.Context, id string, answer *types.Answer) (string, error)
	DownvoteQuestion(ctx context.Context, id string, downvotterId string) (string, error)
	UpvoteQuestion(ctx context.Context, id string, upvotterId string) (string, error)
	AddReplyToAnswer(ctx context.Context, id string, answerId string, reply *types.Reply) (string, error)
}
