package model

import (
	"bhms-ali-iot/global"
	"fmt"
)

type Deflection struct {
	*IotDataBase
	*TagsBase
	DeflectionUnit string  `json:"deflection_unit"`
	Deflection1    float32 `json:"deflection1"`
	Deflection2    float32 `json:"deflection2"`
	Deflection3    float32 `json:"deflection3"`
	Deflection4    float32 `json:"deflection4"`
	Deflection5    float32 `json:"deflection5"`
	Deflection6    float32 `json:"deflection6"`
}

func (m *Deflection) DatabaseName() string {
	return global.CONFIG.TDengine.Dbname
}

func (m *Deflection) StableName() string {
	return "deflection"
}

func (m *Deflection) TableName() string {
	return "deflection_" + m.IotId[0:16]
}

// IotDataBaseColString 数据的基础列名字符串
func (m *Deflection) IotDataBaseColString() string {
	return "ts, yyyy, mm, dd, hh, request_id, message_id"
}

// IotDataBaseValString 数据的基础列值字符串
func (m *Deflection) IotDataBaseValString() string {
	// 与IotDataBaseColString中的列一一对应
	return fmt.Sprintf("%d, '%s', '%s', '%s', '%s', %s, %d", m.Ts, m.Yyyy, m.Mm, m.Dd, m.Hh, m.RequestId, m.MessageId)
}

// BizColString 数据的业务列名字符串
func (m *Deflection) BizColString() string {
	return "deflection_unit, deflection1, deflection2, deflection3, deflection4, deflection5, deflection6"
}

// BizValString 数据的业务列值字符串
func (m *Deflection) BizValString() string {
	// 与BizColString中的列一一对应
	return fmt.Sprintf(
		"'%s', %f, %f, %f, %f, %f, %f",
		m.DeflectionUnit, m.Deflection1, m.Deflection2, m.Deflection3, m.Deflection4, m.Deflection5, m.Deflection6,
	)
}

// TagColString 表的Tag列名字符串
func (m *Deflection) TagColString() string {
	return "product_key, iot_id, device_name, ym"
}

// TagValString 表的Tag列值字符串
func (m *Deflection) TagValString() string {
	// TagColString
	return fmt.Sprintf("'%s', '%s', '%s', '%s'", m.ProductKey, m.IotId, m.DeviceName, m.Ym)
}
