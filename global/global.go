package global

import (
	"bhms-ali-iot/config"
	"database/sql"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	//VP     *viper.Viper
	CONFIG   config.Server
	Logger   *zap.Logger
	TDengine *sql.DB
	Redis    *redis.Client
	AliSms   *dysmsapi20170525.Client
)
