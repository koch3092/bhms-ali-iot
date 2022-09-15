package model

import (
	"bhms-ali-iot/global"
	"fmt"
)

type Measurement struct {
	*TdMetricBase
	*MetricsBase
	*AlarmBase
	*TagsBase
}

func (m *Measurement) DatabaseName() string {
	return global.CONFIG.TDengine.Dbname
}

func (m *Measurement) StableName() string {
	return "measurement"
}

func (m *Measurement) TableName() string {
	return fmt.Sprintf("measure_%d", m.MetricType)
}

// TdMetricsBaseColString 数据的基础列名字符串
func (m *Measurement) TdMetricsBaseColString() string {
	return "ts, yyyy, mm, dd, hh"
}

// TdMetricsBaseValString 数据的基础列值字符串
func (m *Measurement) TdMetricsBaseValString() string {
	// 与IotDataBaseColString中的列一一对应
	return fmt.Sprintf("%d, '%s', '%s', '%s', '%s'", m.Ts, m.Yyyy, m.Mm, m.Dd, m.Hh)
}

// MetricsBaseColString 数据的基础列名字符串
func (m *Measurement) MetricsBaseColString() string {
	return "dt, metric_type, metric_no, field_unit, field_name, field_value"
}

// MetricsBaseValString 数据的基础列值字符串
func (m *Measurement) MetricsBaseValString() string {
	// 与IotDataBaseColString中的列一一对应
	return fmt.Sprintf("%d, %d, %s, '%s', '%s', '%s'", m.Dt, m.MetricType, m.MetricNo, m.FieldUnit, m.FieldName, m.FieldValue)
}

// TagColString 表的Tag列名字符串
func (m *Measurement) TagColString() string {
	return "product_key, iot_id, device_name, ym"
}

// TagValString 表的Tag列值字符串
func (m *Measurement) TagValString() string {
	// TagColString
	return fmt.Sprintf("'%s', '%s', '%s', '%s'", m.ProductKey, m.IotId, m.DeviceName, m.Ym)
}

func (m *Measurement) AlarmColString() string {
	return "alarm_level, alarm_cordon"
}

func (m *Measurement) AlarmValString() string {
	return fmt.Sprintf("%d, %f", m.AlarmLevel, m.AlarmCordon)
}
