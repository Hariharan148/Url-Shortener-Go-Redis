package database

import (
	"os"
	"github.com/go-redis/redis/v8"
	"context"
	"fmt"
)

var Ctx = context.Background()

func Client(dbNo int) *redis.Client{
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASS"),
		DB: dbNo,
	})
	fmt.Println(rdb)

	return rdb

}