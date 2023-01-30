package mongo

import (
	"RuntimeError/repo"
	"RuntimeError/types/mongo"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const usersCollectionName = "users"

// implements UserRepo
// _id will be stored as plain string, not ObjectID
type UserRepoImpl struct{}

func NewUserRepo() repo.UserRepo {
	return &UserRepoImpl{}
}

func (u *UserRepoImpl) GetAll(ctx context.Context) ([]*types.User, error) {
	result, err := GetAll(ctx, usersCollectionName, bson.M{}, UserLabel)
	if err != nil {
		return nil, err
	}

	return result.([]*types.User), nil
}

func (u *UserRepoImpl) GetById(ctx context.Context, id string) (*types.User, error) {
	result, err := GetOne(ctx, usersCollectionName, bson.M{"_id": id}, UserLabel)
	if err != nil {
		return nil, err
	}

	return result.(*types.User), nil
}

func (u *UserRepoImpl) GetByEmail(ctx context.Context, email string) (*types.User, error) {
	result, err := GetOne(ctx, usersCollectionName, bson.M{"email": email}, UserLabel)
	if err != nil {
		return nil, err
	}

	return result.(*types.User), nil
}

func (u *UserRepoImpl) Insert(ctx context.Context, user *types.User) (string, error) {
	// set id, rating
	user.Id = primitive.NewObjectID().Hex()
	user.Rating = 0
	return Insert(ctx, usersCollectionName, user, UserLabel)
}

func (u *UserRepoImpl) Update(ctx context.Context, oldUser *types.User, newUser *types.User) (string, error) {
	return Update(ctx, usersCollectionName, bson.M{"_id": oldUser.Id}, newUser)
}

func (u *UserRepoImpl) Delete(ctx context.Context, id string) error {
	return nil
}
