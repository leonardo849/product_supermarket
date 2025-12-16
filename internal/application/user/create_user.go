package user

import (

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
}

func NewCreateUserUseCase(userRepo    domainUser.UserRepository) *CreateUserUseCase {
	return  &CreateUserUseCase{
		userRepo: userRepo,
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
	return  user.ID, nil
}