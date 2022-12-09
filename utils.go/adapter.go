package utils

import (
	domain "RuntimeError/types/domain"
	mongo "RuntimeError/types/mongo"
)

func DomainUserToMongoUser(u *domain.User) *mongo.User {
	return &mongo.User{
		Id: u.Id,
		Email: u.Email,
		Password: u.Password,
		Rating: u.Rating,
	}
}

func MongoUserToDomainUser(u *mongo.User) *domain.User {
	return &domain.User{
		Id: "",
		Email: u.Email,
		Password: "",
		Rating: u.Rating,
	}
}