package user

import (
	"context"
	"log"

	domainUser "github.com/leonardo849/product_supermarket/internal/domain/user"
)

type FindUserUseByAuthIdCase struct {
	userRepo    domainUser.UserRepository
	userCache domainUser.Cache
}

func NewCreateFindUserUseCaseByAuthId(userRepo    domainUser.UserRepository, userCache domainUser.Cache) *FindUserUseByAuthIdCase {
	return  &FindUserUseByAuthIdCase{
		userRepo: userRepo,
		userCache: userCache,
	}
}

func (uc *FindUserUseByAuthIdCase) Execute(authId string) (*domainUser.User, error) {
	userCache, err := uc.userCache.GetByAuthId(context.Background(), authId)
	if err != nil {
		log.Println(err.Error())
		return nil,err
	}
	if userCache != nil {
		return  userCache, nil
	}

	user, err := uc.userRepo.FindUserByAuthID(authId)
	if err != nil {
		log.Println(err.Error())
		return nil,err
	}
	go func() {
		if err := uc.userCache.Set(context.Background(), user); err != nil {
			log.Print(err.Error())
		}
	}()
	return  user, nil
}