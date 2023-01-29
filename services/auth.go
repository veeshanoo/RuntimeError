package services

import (
	"RuntimeError/auth"
	internalmongo "RuntimeError/db/mongo"
	"RuntimeError/repo"
	types "RuntimeError/types/domain"
	"RuntimeError/utils"
	"context"
	"errors"
)

type AuthService interface {
	Register(ctx context.Context, user *types.User) error
	Login(ctx context.Context, user *types.User) (string, error)
	GetUserData(ctx context.Context, id string) (*types.UserData, error)
	auth.TokenProvider
}

type AuthServiceImpl struct {
	userRepo    repo.UserRepo
	jwtProvider *auth.JWTProvider
}

func NewAuthService() AuthService {
	return &AuthServiceImpl{
		userRepo:    internalmongo.NewUserRepo(),
		jwtProvider: &auth.JWTProvider{},
	}
}

func (svc *AuthServiceImpl) Create(ctx context.Context, user *types.User) (string, error) {
	return svc.jwtProvider.Create(ctx, user)
}

func (svc *AuthServiceImpl) Verify(ctx context.Context, token string) (interface{}, error) {
	return svc.jwtProvider.Verify(ctx, token)
}

func (svc *AuthServiceImpl) Login(ctx context.Context, user *types.User) (string, error) {
	userResult, err := svc.userRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}

	if user.Password != userResult.Password {
		return "", errors.New("unauthorized")
	}

	return svc.Create(ctx, &types.User{Id: userResult.Id, Email: userResult.Email})
}

func (svc *AuthServiceImpl) Register(ctx context.Context, user *types.User) error {
	_, err := svc.userRepo.GetByEmail(ctx, user.Email)
	if err == nil {
		return errors.New("email not available")
	}

	_, err = svc.userRepo.Insert(ctx, utils.DomainUserToMongoUser(user))

	return err
}

func (svc *AuthServiceImpl) GetUserData(ctx context.Context, id string) (*types.UserData, error) {
	user, err := svc.userRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	userData := &types.UserData{
		Id:     user.Id,
		Email:  user.Email,
		Rating: user.Rating,
	}

	return userData, nil
}
