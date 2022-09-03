package config

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"
)

type AliAmqpCredential struct {
	AliAmqpHost     string `mapstructure:"ali-amqp-host" json:"ali_amqp_host" yaml:"ali-amqp-host"`             // 阿里云AMQP连接地址
	AliAmqpPort     string `mapstructure:"ali-amqp-port" json:"ali_amqp_port" yaml:"ali-amqp-port"`             // 阿里云AMQP连接端口
	AliyunUid       string `mapstructure:"aliyun-uid" json:"aliyun_uid" yaml:"access-uid"`                      // 阿里云账号ID
	IotInstanceId   string `mapstructure:"iot-instance-id" json:"iot_instance_id" yaml:"iot-instance-id"`       // 物联网实例ID
	AccessKey       string `mapstructure:"access-key" json:"access_key" yaml:"access-key"`                      // 账号Access Key
	AccessSecret    string `mapstructure:"access-secret" json:"access_secret" yaml:"access-secret"`             // 账号Access Secret
	ConsumerGroupId string `mapstructure:"consumer-group-id" json:"consumer_group_id" yaml:"consumer-group-id"` // 消费者组ID
	ClientId        string `mapstructure:"client-id" json:"client_id" yaml:"client-id"`                         // 客户端ID
}

func (acc *AliAmqpCredential) Credential() (username string, password string) {
	timestamp := time.Now().Nanosecond() / 1000000
	//userName组装方法，请参见AMQP客户端接入说明文档。
	username = fmt.Sprintf("%s|authMode=aksign,signMethod=Hmacsha1,consumerGroupId=%s,authId=%s,iotInstanceId=%s,timestamp=%d|",
		acc.ClientId, acc.ConsumerGroupId, acc.AccessKey, acc.IotInstanceId, timestamp)
	stringToSign := fmt.Sprintf("authId=%s&timestamp=%d", acc.AccessKey, timestamp)
	hmacKey := hmac.New(sha1.New, []byte(acc.AccessSecret))
	hmacKey.Write([]byte(stringToSign))
	//计算签名，password组装方法，请参见AMQP客户端接入说明文档。
	password = base64.StdEncoding.EncodeToString(hmacKey.Sum(nil))
	return username, password
}

func (acc *AliAmqpCredential) Address() (addr string) {
	if acc.AliAmqpHost == "" {
		acc.AliAmqpHost = "iot-amqp.cn-shanghai.aliyuncs.com"
	}
	if acc.AliAmqpPort == "" {
		acc.AliAmqpPort = "5671"
	}
	return fmt.Sprintf("amqps://%s.%s:%s", acc.AliyunUid, acc.AliAmqpHost, acc.AliAmqpPort)
}
