package config

type TDengineRest struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func (t *TDengineRest) Dsn() string {
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
		t.Port = "6041"
	}
	return t.Username + ":" + t.Password + "@http(" + t.Path + ":" + t.Port + ")/" + t.Dbname
}

func (t *TDengineRest) GetLogMode() string {
	return t.LogMode
}
