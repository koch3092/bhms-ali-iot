package config

type TDengine struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`
	Protocol  string `json:"protocol" yaml:"protocol" mapstructure:"protocol"`
}

func (t *TDengine) Dsn() string {
	if t.Username == "" {
		t.Username = "root"
	}
	if t.Password == "" {
		t.Password = "taosdata"
	}
	if t.Path == "" {
		t.Path = "localhost"
	}
	if t.Port == "" {
		t.Port = "6030"
	}
	if t.Protocol == "" {
		t.Protocol = "tcp"
	}

	return t.Username + ":" + t.Password + "@" + t.Protocol + "(" + t.Path + ":" + t.Port + ")/" + t.Dbname
}

func (t *TDengine) GetLogMode() string {
	return t.LogMode
}
