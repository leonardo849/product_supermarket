package user

import (
	"context"

	domainError "github.com/leonardo849/product_supermarket/internal/domain/error"
	domainUser "github.com/leonardo849/product_supermarket/internal/domain/user"
)

type FindIfUserIsInErrors struct {
	errorCache domainError.ErrorCache
	findUserByAuthId *FindUserUseByAuthIdCase
}

func NewFindIfUserIsInErrors(errorCache domainError.ErrorCache, findUserByAuthId *FindUserUseByAuthIdCase) *FindIfUserIsInErrors {
	return  &FindIfUserIsInErrors{
		errorCache: errorCache,
		findUserByAuthId: findUserByAuthId,
	}
}

func (uc *FindIfUserIsInErrors) Execute(authorAuthId string, issuedAt float64, targetAuthId string) (bool, error) {
	author, err := uc.findUserByAuthId.Execute(authorAuthId)
	if err != nil {
		return false, err
	}
	if author == nil {
		return  false, domainUser.ErrUserAuthorNotExists
	}
	if author.UserWasUpdatedAfterToken(issuedAt) {
		return false, domainUser.ErrUserWasUpdatedAfterToken
	}
	has, err := uc.errorCache.HasAuthError(context.Background(), targetAuthId)
	if err != nil {
		return false, err
	}
	return has, nil
}