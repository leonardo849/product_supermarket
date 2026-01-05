package user

import (
	"context"

	"log"
	"github.com/google/uuid"
	domainUser "github.com/leonardo849/product_supermarket/internal/domain/user"
)

type CreateUserInput struct {
	ID string 
	AuthUpdatedAt  string 
	Role string 
}

type CreateUserUseCase struct {
	userRepo    domainUser.UserRepository
	userCache domainUser.Cache
}

func NewCreateUserUseCase(userRepo    domainUser.UserRepository, userCache domainUser.Cache) *CreateUserUseCase {
	return  &CreateUserUseCase{
		userRepo: userRepo,
		userCache: userCache,
	}
}

func (uc *CreateUserUseCase) Execute(input CreateUserInput) (uuid.UUID, error) {
	user, err := domainUser.New(input.ID, domainUser.Role(input.Role), input.AuthUpdatedAt) 
	if err != nil {
		return uuid.Nil, err
	}
	if err := uc.userRepo.Create(user); err != nil {
		return  uuid.Nil, err
	}
	go func() {
		if err := uc.userCache.Set(context.Background(), user); err != nil {
			log.Println(err.Error())
		}
		log.Print("user was setted in cache")
	}()
	return  user.ID, nil
}