package error

import "context"

type ErrorCache interface {
	SetAuthError(ctx context.Context, authId string) error
	HasAuthError(ctx context.Context, authId string) (bool, error)
}