package user

import "errors"

var (
	ErrAuthUpdatedAtAfterNow = errors.New("auth updated at is after now")
	ErrRoleInvalid = errors.New("role is invalid")
	ErrItIsNotAMongoID = errors.New("id isn't a mongo id")
	ErrUserNotFound = errors.New("user wasn't found")
	ErrUserCannotCreateAProduct = errors.New("user can't create a product")
	ErrUserWasUpdatedAfterToken = errors.New("user was updated after token credential_version")
	ErrUserAuthorNotExists = errors.New("author doesn't exist")
)