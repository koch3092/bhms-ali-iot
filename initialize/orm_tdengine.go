package initialize

import (
	"bhms-ali-iot/global"
	"database/sql"
	"errors"
)

func InitTdengine() error {
	m := global.CONFIG.TDengine
	dsn := m.Dsn()
	tdengine, err := sql.Open("taosRestful", dsn)
	if err != nil {
		return err
	}

	defer func() {
		errClose := tdengine.Close()
		if errClose != nil {
			return
		}
	}()

	errPing := tdengine.Ping()
	if errPing != nil {
		return errPing
	}

	_, errUseDB := tdengine.Exec("USE " + global.CONFIG.TDengine.Dbname + ";")
	if errUseDB != nil {
		return errors.New(errUseDB.Error() + " =>> " + global.CONFIG.TDengine.Dbname)
	}
	return nil
}

func TdengineSession() (*sql.DB, error) {
	m := global.CONFIG.TDengine
	dsn := m.Dsn()
	return sql.Open("taosRestful", dsn)
}
