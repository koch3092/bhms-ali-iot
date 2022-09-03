package model

import (
	"bhms-ali-iot/global"
	"fmt"
)

type Seismic struct {
	*IotDataBase
	*TagsBase
	SeismicKey    uint8   `json:"seismic_key"`
	SeismicXUnit  string  `json:"seismic_x_unit"`
	SeismicXValue float32 `json:"seismic_x_value"`
	SeismicZUnit  string  `json:"seismic_z_unit"`
	SeismicZValue float32 `json:"seismic_z_value"`
}

func (m *Seismic) DatabaseName() string {
	return global.CONFIG.TDengine.Dbname
}

func (m *Seismic) StableName() string {
	return "seismic"
}

func (m *Seismic) TableName() string {
	return "seismic_" + m.IotId[0:16]
}

// IotDataBaseColString 数据的基础列名字符串
func (m *Seismic) IotDataBaseColString() string {
	return "ts, request_id, yyyy, mm, dd, hh"
}

// IotDataBaseValString 数据的基础列值字符串
func (m *Seismic) IotDataBaseValString() string {
	// 与IotDataBaseColString中的列一一对应
	return fmt.Sprintf("%d, %s, '%s', '%s', '%s', '%s'", m.Ts, m.RequestId, m.Yyyy, m.Mm, m.Dd, m.Hh)
}

// BizColString 数据的业务列名字符串
func (m *Seismic) BizColString() string {
	return "seismic_key, seismic_x_unit, seismic_z_unit, seismic_x_value, seismic_z_value"
}

// BizValString 数据的业务列值字符串
func (m *Seismic) BizValString() string {
	// 与BizColString中的列一一对应
	return fmt.Sprintf(
		"'%d', '%s', '%s', %f, %f",
		m.SeismicKey, m.SeismicXUnit, m.SeismicZUnit, m.SeismicXValue, m.SeismicZValue,
	)
}

// TagColString 表的Tag列名字符串
func (m *Seismic) TagColString() string {
	return "product_key, iot_id, device_name, ym"
}

// TagValString 表的Tag列值字符串
func (m *Seismic) TagValString() string {
	// TagColString
	return fmt.Sprintf("'%s', '%s', '%s', '%s'", m.ProductKey, m.IotId, m.DeviceName, m.Ym)
}
