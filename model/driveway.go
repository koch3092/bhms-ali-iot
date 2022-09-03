package model

import (
	"bhms-ali-iot/global"
	"fmt"
)

type Driveway struct {
	*IotDataBase
	*TagsBase
	DrivewayKey         uint8   `json:"driveway_key"`
	DrivewayWeightUnit  string  `json:"driveway_weight_unit"`
	DrivewayWeightValue float32 `json:"driveway_weight_value"`
	DrivewaySpeedUnit   string  `json:"driveway_speed_unit"`
	DrivewaySpeedValue  float32 `json:"driveway_speed_value"`
	DrivewayModelUnit   string  `json:"driveway_model_unit"`
	DrivewayModelValue  int8    `json:"driveway_model_value"`
	DrivewayModelLabel  string  `json:"driveway_model_label"`
}

func (m *Driveway) DatabaseName() string {
	return global.CONFIG.TDengine.Dbname
}

func (m *Driveway) StableName() string {
	return "driveway"
}

func (m *Driveway) TableName() string {
	return "driveway_" + m.IotId[0:16]
}

// IotDataBaseColString 数据的基础列名字符串
func (m *Driveway) IotDataBaseColString() string {
	return "ts, request_id, yyyy, mm, dd, hh"
}

// IotDataBaseValString 数据的基础列值字符串
func (m *Driveway) IotDataBaseValString() string {
	// 与IotDataBaseColString中的列一一对应
	return fmt.Sprintf("%d, %s, '%s', '%s', '%s', '%s'", m.Ts, m.RequestId, m.Yyyy, m.Mm, m.Dd, m.Hh)
}

// BizColString 数据的业务列名字符串
func (m *Driveway) BizColString() string {
	return "driveway_key, driveway_weight_unit, driveway_speed_unit, driveway_model_unit, driveway_weight_value, driveway_speed_value, driveway_model_value, driveway_model_label"
}

// BizValString 数据的业务列值字符串
func (m *Driveway) BizValString() string {
	// 与BizColString中的列一一对应
	return fmt.Sprintf(
		"'%d', '%s', '%s', '%s', %f, %f, %d, '%s'",
		m.DrivewayKey, m.DrivewayWeightUnit, m.DrivewaySpeedUnit, m.DrivewayModelUnit, m.DrivewayWeightValue,
		m.DrivewaySpeedValue, m.DrivewayModelValue, m.DrivewayModelLabel,
	)
}

// TagColString 表的Tag列名字符串
func (m *Driveway) TagColString() string {
	return "product_key, iot_id, device_name, ym"
}

// TagValString 表的Tag列值字符串
func (m *Driveway) TagValString() string {
	// TagColString
	return fmt.Sprintf("'%s', '%s', '%s', '%s'", m.ProductKey, m.IotId, m.DeviceName, m.Ym)
}
