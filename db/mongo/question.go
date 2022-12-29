package mongo

import (
	types "RuntimeError/types/mongo"
	"context"

	"go.mongodb.org/mongo-driver/bson"
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

	return Insert(ctx, questionsCollectionName, question, QuestionLabel)
}

func (q *QuestionRepoImpl) Update(ctx context.Context, oldQuestion *types.Question, newQuestion *types.Question) (string, error) {
	return Update(ctx, questionsCollectionName, bson.M{"_id": oldQuestion.Id}, newQuestion)
}

func (q *QuestionRepoImpl) Delete(ctx context.Context, id string) error {
	return nil
}

func (u *QuestionRepoImpl) GetAll(ctx context.Context) ([]*types.Question, error) {
	result, err := GetAll(ctx, questionsCollectionName, bson.M{}, QuestionLabel)
	if err != nil {
		return nil, err
	}

	return result.([]*types.Question), nil
}

func (u *QuestionRepoImpl) GetById(ctx context.Context, id string) (*types.Question, error) {
	result, err := GetOne(ctx, questionsCollectionName, bson.M{"_id": id}, QuestionLabel)
	if err != nil {
		return nil, err
	}

	return result.(*types.Question), nil
}

func (u *QuestionRepoImpl) UpvoteQuestion(ctx context.Context, id string, upvotterId string) (string, error) {
	result, err := GetOne(ctx, questionsCollectionName, bson.M{"_id": id}, QuestionLabel)
	if err != nil {
		return "", err
	}
	question := result.(*types.Question)
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
	return u.Update(ctx, result.(*types.Question), question)
}

func (u *QuestionRepoImpl) DownvoteQuestion(ctx context.Context, id string, downvotterId string) (string, error) {
	result, err := GetOne(ctx, questionsCollectionName, bson.M{"_id": id}, QuestionLabel)
	if err != nil {
		return "", err
	}
	question := result.(*types.Question)
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
	return u.Update(ctx, result.(*types.Question), question)
}

func (u *QuestionRepoImpl) AddAnswerToQuestion(ctx context.Context, id string, answer *types.Answer) (string, error) {
	result, err := GetOne(ctx, questionsCollectionName, bson.M{"_id": id}, QuestionLabel)
	if err != nil {
		return "", err
	}
	question := result.(*types.Question)
	question.Answers = append(question.Answers, *answer)
	return u.Update(ctx, result.(*types.Question), question)
}

func (u *QuestionRepoImpl) AddReplyToAnswer(ctx context.Context, id string, answerId string, reply *types.Reply) (string, error) {
	result, err := GetOne(ctx, questionsCollectionName, bson.M{"_id": id}, QuestionLabel)
	if err != nil {
		return "", err
	}
	question := result.(*types.Question)
	for _, answer := range question.Answers {
		if answer.Id == answerId {
			answer.Replies = append(answer.Replies, *reply)
		}
	}
	return u.Update(ctx, result.(*types.Question), question)
}

func (u *QuestionRepoImpl) EditContent(ctx context.Context, id string, newContent string) (string, error) {
	result, err := GetOne(ctx, questionsCollectionName, bson.M{"_id": id}, QuestionLabel)
	if err != nil {
		return "", err
	}
	question := result.(*types.Question)
	question.Contents = newContent
	return u.Update(ctx, result.(*types.Question), question)
}

func (u *QuestionRepoImpl) FavoriteComment(ctx context.Context, id string, bestAnswer string) (string, error) {
	result, err := GetOne(ctx, questionsCollectionName, bson.M{"_id": id}, QuestionLabel)
	if err != nil {
		return "", err
	}
	question := result.(*types.Question)
	question.BestAnswer = bestAnswer
	return u.Update(ctx, result.(*types.Question), question)
}
