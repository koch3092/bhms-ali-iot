package config

type Server struct {
	TDengine    TDengine          `mapstructure:"tdengine" json:"tdengine" yaml:"tdengine"`
	AliAmqpCred AliAmqpCredential `mapstructure:"ali-amqp-cred" json:"ali_amqp_cred" yaml:"ali-amqp-cred"`
	Cordons     Cordons           `mapstructure:"cordons" json:"cordons" yaml:"cordons"`
	Redis       Redis             `mapstructure:"redis" json:"redis" yaml:"redis"`
}
