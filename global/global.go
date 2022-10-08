package global

import (
	"bhms-ali-iot/config"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	//VP     *viper.Viper
	CONFIG   config.Server
	Logger   *zap.Logger
	TDengine *sql.DB
	Redis    *redis.Client
)
