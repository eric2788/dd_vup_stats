package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"strconv"
)

var (
	cli *redis.Client
	ctx = context.Background()
)

const VupListKey = "vup_list"
const VupBlackListKey = "vup_blacklist"

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

func SetAdd(key string, value string) error {
	return cli.SAdd(ctx, key, value).Err()
}

func SetRemove(key string, value string) error {
	return cli.SRem(ctx, key, value).Err()
}

func SetContain(key string, value string) (bool, error) {
	return cli.SIsMember(ctx, key, value).Result()
}

func SetGet(key string) ([]string, error) {
	return cli.SMembers(ctx, key).Result()
}

func PutMap(key string, dict interface{}) error {
	return cli.HSet(ctx, key, dict).Err()
}

func GetMap(key string, dict *map[string]string) error {
	d, err := cli.HGetAll(ctx, key).Result()
	if err != nil {
		return err
	}
	*dict = d
	return nil
}

func GetMapFields(key string, dict *map[string]string, fields ...string) ([]string, error) {

	var errorField []string
	values := make(map[string]string)

	for _, field := range fields {
		value, err := cli.HGet(ctx, key, field).Result()
		if err != nil {
			if err != redis.Nil {
				log.Warnf("Redis 獲取 %s 中的 %s 值時出現錯誤: %v", key, field, err)
				return nil, err
			}
			errorField = append(errorField, field)
		} else {
			values[field] = value
		}
	}

	*dict = values
	return errorField, nil
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
