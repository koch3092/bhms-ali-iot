package core

import (
	"bhms-ali-iot/global"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

// Viper 读取配置文件，并赋值给全局配置变量global.CONFIG
func Viper() *viper.Viper {
	v := viper.New()
	workingDir, _ := os.Getwd()
	v.AddConfigPath("/etc/bhms")
	v.AddConfigPath(workingDir)
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvPrefix("app")
	// 读取配置文件
	err := v.ReadInConfig()
	if err != nil {
		panic(any("Fatal error read config.\n"))
	}
	// 将配置赋值
	if err = v.Unmarshal(&global.CONFIG); err != nil {
		fmt.Println(err)
	}

	return v
}
