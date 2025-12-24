package user

import (
	domainUser "github.com/leonardo849/product_supermarket/internal/domain/user"
)

type DeleteUserUseCase struct {
	userRepo domainUser.UserRepository
}

func NewDeleteUserUseCase(userRepo    domainUser.UserRepository) *DeleteUserUseCase {
	return  &DeleteUserUseCase{
		userRepo: userRepo,
	}
}

func (uc *DeleteUserUseCase) Execute(id string) error {
	if _,err := uc.userRepo.FindUserByAuthID(id); err != nil {
		return  err
	}

	if err := uc.userRepo.DeleteUserByAuthId(id); err != nil {
		return  err
	}
	return  nil
}