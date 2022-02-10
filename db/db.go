package db

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
)

type Database struct {
	Client *redis.Client
}

var (
	ErrNil = errors.New("no matching record found in redis database")
	// dunno why it's using a context here but i guess we'll find out
	Ctx    = context.TODO()
)

func NewDatabase(address string) (*Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address,
		Password: "",
		DB: 0,
	})
	// ping da server
	if err := client.Ping(Ctx).Err(); err != nil {
		return nil, err
	}
	return &Database{
		Client: client,
	}, nil
}

