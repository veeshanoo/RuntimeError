package services

import (
	internalmongo "RuntimeError/db/mongo"
	"RuntimeError/repo"
	"RuntimeError/types"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type LoginService interface {
	Register(ctx context.Context) (*types.User, error)
	Login() error
	Logout() error
}

type LoginServiceImpl struct {
	userRepo repo.UserRepo
}

func NewLoginService() LoginService {
	return &LoginServiceImpl{
		userRepo: internalmongo.NewUserRepo(),
	}
}

func (svc *LoginServiceImpl) Register(ctx context.Context) (*types.User, error) {
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

func (svc *LoginServiceImpl) Login() error {
	return nil
}

func (svc *LoginServiceImpl) Logout() error {
	return nil
}