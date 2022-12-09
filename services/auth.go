package services

import (
	"RuntimeError/auth"
	internalmongo "RuntimeError/db/mongo"
	"RuntimeError/repo"
	"RuntimeError/types/domain"
	"RuntimeError/utils.go"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService interface {
	Register(ctx context.Context, user *types.User) error
	Login(ctx context.Context, user *types.User) (string, error)
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

	return svc.Create(ctx, &types.User{Id: userResult.Id})
}

func (svc *AuthServiceImpl) Register(ctx context.Context, user *types.User) error {
	_, err := svc.userRepo.GetByEmail(ctx, user.Email)
	if err == nil {
		return errors.New("email not available")
	}

	_, err = svc.userRepo.Insert(ctx, utils.DomainUserToMongoUser(user))

	return err
}

func (svc *AuthServiceImpl) Transaction(ctx context.Context) (*types.User, error) {
	t := internalmongo.NewTransaction(func(ctx context.Context) (interface{}, error) {
		switch ctx.(type) {
		case mongo.SessionContext:
			fmt.Println("WOW")
		}

		_, err := svc.userRepo.GetById(ctx, "637e7f29bc95879e58d567fb")
		if err == nil {
			fmt.Println("WORKS 1")
		}
		usr, err := svc.userRepo.GetById(ctx, "637e7f29bc95879e58d567fb")
		if err == nil {
			fmt.Println("WORKS 2")
		}

		return usr, nil
	})

	if res, err := t.Execute(ctx); err != nil {
		return nil, err
	} else {
		return res.(*types.User), nil
	}
}
