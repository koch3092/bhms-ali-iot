package main

import (
	"bhms-ali-iot/core"
	"bhms-ali-iot/global"
	"bhms-ali-iot/initialize"
	"context"
	_ "github.com/taosdata/driver-go/v2/taosRestful"
	"pack.ag/amqp"
)

func main() {
	_ = core.Viper()                                  // 初始化Viper，将读取配置，并赋值给global.CONFIG
	global.Logger = core.Zap()                        // 初始化Zap，将配置日志打印程序，并赋值给global.Logger
	if err := initialize.InitTdengine(); err != nil { // 初始化TDengine，保证连通性
		global.Logger.DPanic("Init TDengine Error: " + err.Error())
		panic(any(err))
	}

	ctx := context.Background()

	// 用于转发消息的Channel
	rcvMessage := make(chan *amqp.Message)

	// 阿里云AMQP凭证对象
	aliCred := global.CONFIG.AliAmqpCred
	address := aliCred.Address()
	username, password := aliCred.Credential()
	amqpManager := initialize.AmqpManager{
		Address:  address,
		Username: username,
		Password: password,
		Logger:   global.Logger,
	}
	// 开启阿里云AMQP客户端
	go amqpManager.StartReceiveMessage(ctx, rcvMessage)

	// 从Channel中提取并处理数据
	msgHandler := initialize.MessageHandler{
		Logger: global.Logger,
	}
	go msgHandler.Handle(ctx, rcvMessage)

	<-ctx.Done()
}
