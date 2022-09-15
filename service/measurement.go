package service

import (
	"database/sql"
	"go.uber.org/zap"
)

type MeasurementService struct {
	Logger *zap.Logger
}

func (s MeasurementService) CreateMeasurement(session *sql.DB, sqlCreate string) error {
	s.Logger.Debug("Measurement SQL: " + sqlCreate)
	result, errExec := session.Exec(sqlCreate)
	if errExec != nil {
		return errExec
	}
	rowsAffected, errRA := result.RowsAffected()
	if errRA != nil {
		return errRA
	}
	if rowsAffected <= 0 {
		panic(any("Create measurement error"))
	}
	return nil
}
