package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}


func NewClient(addr, password string, db int) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("redis is connected")

	return &Client{client: client}, nil
}

func (c *Client) Close() error {
	return  c.client.Close()
}