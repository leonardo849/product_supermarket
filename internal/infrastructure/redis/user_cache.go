package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	user "github.com/leonardo849/product_supermarket/internal/domain/user"
	goredis "github.com/redis/go-redis/v9"
)

type UserCache struct {
	client *Client
	ttl    time.Duration
}


func NewUserCache(client *Client, ttl time.Duration) *UserCache {
	return &UserCache{
		client: client,
		ttl: ttl,
	}
}

func (u *UserCache) Get(ctx context.Context, id string) (*user.User, error) {
	key := "user:" + id
	
	val, err := u.client.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var user user.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserCache) GetByAuthId(ctx context.Context, authId string) (*user.User, error) {
	key := "user:auth_id:" + authId

	val, err := u.client.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var user user.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserCache) Set(ctx context.Context,  user *user.User) error {
	key := "user:" + user.ID.String()

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	pipe := u.client.client.TxPipeline()

	keyAuthId := "user:auth_id:" + user.AuthID
	pipe.Set(ctx, key, data, u.ttl)
	pipe.Set(ctx, keyAuthId, data, u.ttl)

	_, err = pipe.Exec(ctx)
	return  err
}

func (u *UserCache) DeleteUser(ctx context.Context, id string, authId string) error {
	pipe := u.client.client.TxPipeline()

	pipe.Del(ctx, "user:" + id, "user:auth_id:" + authId)

	_, err := pipe.Exec(ctx)
	return  err
}
