package services

import (
	types "RuntimeError/types/domain"
	"context"
)

type QuestionService interface {
	AddQuestion(ctx context.Context, userId*types.Answer)
}
