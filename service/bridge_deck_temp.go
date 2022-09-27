package service

import (
	"bhms-ali-iot/model"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

type BridgeDeckTempService struct {
	Logger *zap.Logger
}

func (s BridgeDeckTempService) CreateBridgeDeckTemp(session *sql.DB, m *model.BridgeDeckTemp) error {
	sqlCreate := fmt.Sprintf(
		"INSERT INTO %s(%s, %s) USING %s TAGS (%s) VALUES (%s, %s)",
		m.TableName(), m.IotDataBaseColString(), m.BizColString(), m.StableName(), m.TagValString(), m.IotDataBaseValString(), m.BizValString(),
	)
	s.Logger.Debug("SQL: " + sqlCreate)
	result, errExec := session.Exec(sqlCreate)
	if errExec != nil {
		s.Logger.Debug("Exec SQL failed: '" + sqlCreate + "'")
		return errExec
	}
	rowsAffected, errRA := result.RowsAffected()
	if errRA != nil {
		return errRA
	}
	if rowsAffected != 1 {
		panic(any("Create Bridge Deck Temp error"))
	}
	return nil
}
