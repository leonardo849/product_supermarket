package redis

import (
	"context"
	"time"
)

type ErrorCache struct {
	client *Client
	ttl    time.Duration
}

func NewErrorCache(client *Client, ttl time.Duration) *ErrorCache {
	return &ErrorCache{
		client: client,
		ttl: ttl,
	}
}

func (e *ErrorCache) SetAuthError(ctx context.Context, authId string) error {
	key := "auth:error:user:" + authId
	
	return e.client.client.Set(ctx, key, 1, e.ttl).Err()
}

func (e *ErrorCache) HasAuthError(ctx context.Context, authId string) (bool, error) {
	key := "auth:error:user:" + authId
	exists, err := e.client.client.Exists(ctx, key).Result()
	if err != nil {
		return  false, err
	}
	return  exists == 1, nil
}