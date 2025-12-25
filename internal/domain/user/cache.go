package user

import "context"

type Cache interface {
	Get(ctx context.Context, id string) (*User, error)
	Set(ctx context.Context, u *User) error
	GetByAuthId(ctx context.Context, authId string) (*User, error)
	DeleteUser(ctx context.Context, id string, authId string) (error)
}