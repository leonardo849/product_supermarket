package user

import (
	"context"

	domainUser "github.com/leonardo849/product_supermarket/internal/domain/user"
)

type DeleteUserUseCase struct {
	userRepo domainUser.UserRepository
	userCache domainUser.Cache
}

func NewDeleteUserUseCase(userRepo    domainUser.UserRepository, userCache domainUser.Cache) *DeleteUserUseCase {
	return  &DeleteUserUseCase{
		userRepo: userRepo,
		userCache: userCache,
	}
}

func (uc *DeleteUserUseCase) Execute(id string) error {
	user,err := uc.userRepo.FindUserByAuthID(id)
	if err != nil {
		return  err
	}


	if err := uc.userRepo.DeleteUserByAuthId(id); err != nil {
		return  err
	}

	if err := uc.userCache.DeleteUser(context.Background() ,user.ID.String(), user.AuthID); err != nil {
		return  err
	}

	return  nil
}