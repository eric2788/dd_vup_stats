package db

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
)

var (
	cli *redis.Client
	ctx = context.Background()
)

func InitRedis() {
	log.Info("正在初始化 redis..")
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Error("env REDIS_DB is not a number, use 0 as db")
		db = 0
	}
	cli = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		DB:   db,
	})
	log.Info("redis 初始化成功")
}

func Put(key string, value interface{}) error {
	err := cli.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Errorf("redis set error: %v", err)
		return err
	}
	return nil
}

func Get(key string) (interface{}, error) {
	val, err := cli.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Errorf("redis get error: %v", err)
		}
		return nil, err
	}
	return val, nil
}

func PutUserIsVup(uid int64, isVup bool) error {
	key := fmt.Sprintf("is_vup:%d", uid)
	err := Put(key, isVup)
	if err != nil {
		log.Errorf("redis set error: %v", err)
		return err
	}
	return nil
}

func GetUserIsVup(uid int64) (bool, bool) {
	val, err := Get(fmt.Sprintf("is_vup:%v", uid))
	if err != nil {
		return false, false
	}
	return val == "1", true
}
