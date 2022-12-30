package mongo

import (
	types "RuntimeError/types/mongo"
	"context"
	"errors"

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
	return Insert(ctx, questionsCollectionName, question, QuestionLabel)
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

func (q *QuestionRepoImpl) GetById(ctx context.Context, id string) (*types.Question, error) {
	result, err := GetOne(ctx, questionsCollectionName, bson.M{"id": id}, QuestionLabel)
	if err != nil {
		return nil, err
	}

	return result.(*types.Question), nil
}

func (q *QuestionRepoImpl) UpvoteQuestion(ctx context.Context, id string, upvotterId string) (string, error) {
	result, err := q.GetById(ctx, id)
	if err != nil {
		return "", err
	}
	question := result
	for _, upvoter := range question.Upvoters {
		if upvoter == upvotterId {
			// already exists
			return "", nil
		}
	}

	for idx, downvotter := range question.Downvoters {
		if downvotter == upvotterId {
			// remove selected item
			question.Downvoters = append(question.Downvoters[:idx], question.Downvoters[idx+1:]...)
		}
	}

	question.Upvoters = append(question.Upvoters, upvotterId)
	return q.Update(ctx, result, question)
}

func (q *QuestionRepoImpl) DownvoteQuestion(ctx context.Context, id string, downvotterId string) (string, error) {
	result, err := q.GetById(ctx, id)
	if err != nil {
		return "", err
	}
	question := result
	for _, downvotter := range question.Downvoters {
		if downvotter == downvotterId {
			// already exists
			return "", nil
		}
	}

	for idx, upvoter := range question.Upvoters {
		if upvoter == downvotterId {
			// remove selected item
			question.Downvoters = append(question.Upvoters[:idx], question.Upvoters[idx+1:]...)
		}
	}
	question.Downvoters = append(question.Downvoters, downvotterId)
	return q.Update(ctx, result, question)
}

func (q *QuestionRepoImpl) AddAnswerToQuestion(ctx context.Context, id string, answer *types.Answer) (string, error) {
	result, err := q.GetById(ctx, id)
	if err != nil {
		return "", err
	}
	question := result
	question.Answers = append(question.Answers, *answer)
	return q.Update(ctx, result, question)
}

func (q *QuestionRepoImpl) AddReplyToAnswer(ctx context.Context, id string, answerId string, reply *types.Reply) (string, error) {
	result, err := q.GetById(ctx, id)
	if err != nil {
		return "", err
	}
	question := result
	found := false
	for _, answer := range question.Answers {
		if answer.Id == answerId {
			answer.Replies = append(answer.Replies, *reply)
			found = true
		}
	}
	if !found {
		return "", errors.New("No answer found.")
	}
	return q.Update(ctx, result, question)
}

func (q *QuestionRepoImpl) EditContent(ctx context.Context, id string, newContent string) (string, error) {
	result, err := q.GetById(ctx, id)
	if err != nil {
		return "", err
	}

	question := result
	question.Contents = newContent
	return q.Update(ctx, result, question)
}

func (q *QuestionRepoImpl) FavoriteComment(ctx context.Context, id string, bestAnswer string) (string, error) {
	result, err := q.GetById(ctx, id)
	if err != nil {
		return "", err
	}
	question := result
	if question.BestAnswer != "" {
		return "", errors.New("Best answer already exists.")
	}
	answerExist := false
	for _, answer := range question.Answers {
		if answer.Id == bestAnswer {
			answerExist = true
		}
	}
	if !answerExist {
		return "", errors.New("Answer doesn't exists")
	}
	question.BestAnswer = bestAnswer
	return q.Update(ctx, result, question)
}
