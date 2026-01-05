package user

import (
	"context"
	"log"

	domainUser "github.com/leonardo849/product_supermarket/internal/domain/user"
)

type FindUserUseCaseById struct {
	userRepo    domainUser.UserRepository
	userCache domainUser.Cache
}

func NewCreateFindUserUseCaseById(userRepo    domainUser.UserRepository, userCache domainUser.Cache) *FindUserUseCaseById {
	return  &FindUserUseCaseById{
		userRepo: userRepo,
		userCache: userCache,
	}
}

func (uc *FindUserUseCaseById) Execute(id string) (*domainUser.User, error) {
	userCache, err := uc.userCache.Get(context.Background(), id)
	if err != nil {
		log.Println(err.Error())
		return nil,err
	}
	if userCache != nil {
		return  userCache, nil
	}

	user, err := uc.userRepo.FindUserById(id)
	if err != nil {
		log.Println(err.Error())
		return nil,err
	}
	return  user, nil
}