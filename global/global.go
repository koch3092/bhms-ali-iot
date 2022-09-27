package global

import (
	"bhms-ali-iot/config"
	"database/sql"
	"go.uber.org/zap"
)

var (
	CONFIG config.Server
	Logger *zap.Logger
	//VP     *viper.Viper
	TDengine *sql.DB
)
