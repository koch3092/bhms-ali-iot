package initialize

import (
	"bhms-ali-iot/global"
	"context"
	"github.com/go-redis/redis/v8"
)

func InitRedis() (*redis.Client, error) {
	r := global.CONFIG.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     r.Address,
		Password: r.Password,
		DB:       r.Db,
	})

	_, err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		global.Logger.Debug("Ping redis failed: " + err.Error())
		return nil, err
	}

	return rdb, nil
}
