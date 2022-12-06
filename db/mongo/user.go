package mongo

import (
	"RuntimeError/repo"
	"RuntimeError/types"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const usersCollectionName = "test"

// User implements repo.UserRepo
type UserRepoImpl struct {}

func NewUserRepo() repo.UserRepo {
	return &UserRepoImpl{}
}

func (u *UserRepoImpl) GetAll(ctx context.Context) ([]*types.User, error) {
	result, err := GetAll(ctx, usersCollectionName, bson.M{}, UserLabel)
	if err != nil {
		return nil, err
	}

	users, ok := result.([]*types.User)
	if !ok {
		return nil, errors.New("failed to convert result")
	}

	return users, nil
}

func (u *UserRepoImpl) GetById(ctx context.Context, id string) (*types.User, error) {
	mongoId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := GetOne(ctx, usersCollectionName, bson.M{"_id": mongoId}, UserLabel)
	if err != nil {
		return nil, err
	}

	user, ok := result.(*types.User)
	if !ok {
		return nil, errors.New("failed to convert result")
	}

	return user, nil
}

func (u *UserRepoImpl) Insert(ctx context.Context, user *types.User) (string, error) {
	return "", nil
}

func (u *UserRepoImpl) Update(ctx context.Context, oldUser *types.User, newUser *types.User) (*types.User, error) {
	return nil, nil
}

func (u *UserRepoImpl) Delete(ctx context.Context, id string) error {
	return nil
}

