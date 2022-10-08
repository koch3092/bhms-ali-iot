package initialize

import (
	"bhms-ali-iot/global"
	"database/sql"
)

func InitTdengine() (*sql.DB, error) {
	m := global.CONFIG.TDengine
	dsn := m.Dsn()
	global.Logger.Debug("Open TDengin dsn: " + dsn)
	tdengine, err := sql.Open("taosSql", dsn)
	if err != nil {
		global.Logger.Debug("Open TDengine failed: " + err.Error())
		return nil, err
	}

	errPing := tdengine.Ping()
	if errPing != nil {
		global.Logger.Debug("Ping TDengine failed: " + errPing.Error())
		return nil, errPing
	}
	return tdengine, nil
}
