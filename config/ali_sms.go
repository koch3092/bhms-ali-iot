package config

type AliSms struct {
	AccessKey    string `mapstructure:"access-key" json:"access_key" yaml:"access-key"`          // 账号Access Key
	AccessSecret string `mapstructure:"access-secret" json:"access_secret" yaml:"access-secret"` // 账号Access Secret
	SignName     string `mapstructure:"sign-name" json:"sign_name" yaml:"sign-name"`             // 短信的 SignName
}
