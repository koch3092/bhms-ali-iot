package model

import (
	"bhms-ali-iot/global"
	"fmt"
)

type StaticStrain struct {
	*IotDataBase
	*TagsBase
	StaticStrainKey    uint8   `json:"static_strain_key"`
	SSTemperatureUnit  string  `json:"ss_temperature_unit"`
	SSTemperatureValue float32 `json:"ss_temperature_value"`
	SSStrainUnit       string  `json:"ss_strain_unit"`
	SSStrainValue      float32 `json:"ss_strain_value"`
}

func (m *StaticStrain) DatabaseName() string {
	return global.CONFIG.TDengine.Dbname
}

func (m *StaticStrain) StableName() string {
	return "static_strain"
}

func (m *StaticStrain) TableName() string {
	return "ss_" + m.IotId[0:16]
}

// IotDataBaseColString 数据的基础列名字符串
func (m *StaticStrain) IotDataBaseColString() string {
	return "ts, request_id, message_id, yyyy, mm, dd, hh"
}

// IotDataBaseValString 数据的基础列值字符串
func (m *StaticStrain) IotDataBaseValString() string {
	// 与IotDataBaseColString中的列一一对应
	return fmt.Sprintf("%d, %s, %d, '%s', '%s', '%s', '%s'", m.Ts, m.RequestId, m.MessageId, m.Yyyy, m.Mm, m.Dd, m.Hh)
}

// BizColString 数据的业务列名字符串
func (m *StaticStrain) BizColString() string {
	return "static_strain_key, ss_temperature_unit, ss_strain_unit, ss_temperature_value, ss_strain_value"
}

// BizValString 数据的业务列值字符串
func (m *StaticStrain) BizValString() string {
	// 与BizColString中的列一一对应
	return fmt.Sprintf(
		"'%d', '%s', '%s', %f, %f",
		m.StaticStrainKey, m.SSTemperatureUnit, m.SSStrainUnit, m.SSTemperatureValue, m.SSStrainValue,
	)
}

// TagColString 表的Tag列名字符串
func (m *StaticStrain) TagColString() string {
	return "product_key, iot_id, device_name, ym"
}

// TagValString 表的Tag列值字符串
func (m *StaticStrain) TagValString() string {
	// TagColString
	return fmt.Sprintf("'%s', '%s', '%s', '%s'", m.ProductKey, m.IotId, m.DeviceName, m.Ym)
}
