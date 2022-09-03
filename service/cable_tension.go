package service

import (
	"bhms-ali-iot/model"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

type CableTensionService struct {
	Logger *zap.Logger
}

func (s CableTensionService) CreateCableTension(session *sql.DB, m *model.CableTension) error {
	sqlCreate := fmt.Sprintf(
		"INSERT INTO %s.%s(%s, %s) USING %s TAGS (%s) VALUES (%s, %s)",
		m.DatabaseName(), m.TableName(), m.IotDataBaseColString(), m.BizColString(), m.StableName(), m.TagValString(), m.IotDataBaseValString(), m.BizValString(),
	)
	s.Logger.Debug("SQL: " + sqlCreate)
	result, errExec := session.Exec(sqlCreate)
	if errExec != nil {
		return errExec
	}
	rowsAffected, errRA := result.RowsAffected()
	if errRA != nil {
		return errRA
	}
	if rowsAffected != 1 {
		panic(any("Create Cable Tension error"))
	}
	return nil
}
