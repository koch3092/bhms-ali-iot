package service

import (
	"bhms-ali-iot/model"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

type DeflectionService struct {
	Logger *zap.Logger
}

func (s DeflectionService) CreateDeflection(session *sql.DB, m *model.Deflection) error {
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
		panic(any("Create Deflection error"))
	}
	return nil
}
