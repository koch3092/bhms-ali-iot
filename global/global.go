package global

import (
	"bhms-ali-iot/config"
	"go.uber.org/zap"
)

var (
	CONFIG config.Server
	Logger *zap.Logger
	//VP     *viper.Viper
)
