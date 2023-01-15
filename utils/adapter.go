package utils

import (
	domain "RuntimeError/types/domain"
	mongo "RuntimeError/types/mongo"
)

func DomainUserToMongoUser(u *domain.User) *mongo.User {
	return &mongo.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Rating:   u.Rating,
	}
}

func MongoUserToDomainUser(u *mongo.User) *domain.User {
	return &domain.User{
		Id:       "",
		Email:    u.Email,
		Password: "",
		Rating:   u.Rating,
	}
}

func DomainReplyToMongo(replies []domain.Reply) []mongo.Reply {
	x := []mongo.Reply{}
	for _, reply := range replies {
		x = append(x, mongo.Reply{
			Id:          reply.Id,
			Contents:    reply.Contents,
			SubmitterId: reply.SubmitterId,
		})
	}
	return x
}

func DomainAnswerToMongo(answers []domain.Answer) []mongo.Answer {
	x := []mongo.Answer{}
	for _, answer := range answers {
		x = append(x, mongo.Answer{
			Id:          answer.Id,
			Contents:    answer.Contents,
			SubmitterId: answer.SubmitterId,
			Replies:     DomainReplyToMongo(answer.Replies),
		})
	}
	return x
}

func DomainQuestionToMongo(question *domain.Question) *mongo.Question {

	x := &mongo.Question{
		Id:         question.Id,
		SumitterId: question.SumitterId,
		Title:      question.Title,
		Contents:   question.Contents,
		Answers:    DomainAnswerToMongo(question.Answers),
		BestAnswer: question.BestAnswer,
		Upvoters:   question.Upvoters,
		Downvoters: question.Downvoters,
	}
	return x
}

func MongoReplyToDomain(replies []mongo.Reply) []domain.Reply {
	x := []domain.Reply{}
	for _, reply := range replies {
		x = append(x, domain.Reply{
			Id:          reply.Id,
			Contents:    reply.Contents,
			SubmitterId: reply.SubmitterId,
		})
	}
	return x
}

func MongoAnswerToDomain(answers []mongo.Answer) []domain.Answer {
	x := []domain.Answer{}
	for _, answer := range answers {
		x = append(x, domain.Answer{
			Id:          answer.Id,
			Contents:    answer.Contents,
			SubmitterId: answer.SubmitterId,
			Replies:     MongoReplyToDomain(answer.Replies),
		})
	}
	return x
}

func MongoQuestionToDomain(question *mongo.Question) *domain.Question {

	x := &domain.Question{
		Id:         question.Id,
		SumitterId: question.SumitterId,
		Title:      question.Title,
		Contents:   question.Contents,
		Answers:    MongoAnswerToDomain(question.Answers),
		BestAnswer: question.BestAnswer,
		Upvoters:   question.Upvoters,
		Downvoters: question.Downvoters,
	}
	return x
}
