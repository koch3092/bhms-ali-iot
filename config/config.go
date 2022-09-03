package config

type Server struct {
	TDengine    TDengineRest      `mapstructure:"tdengine" json:"tdengine" yaml:"tdengine"`
	AliAmqpCred AliAmqpCredential `mapstructure:"ali-amqp-cred" json:"ali_amqp_cred" yaml:"ali-amqp-cred"`
}
