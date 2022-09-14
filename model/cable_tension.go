package model

import (
	"bhms-ali-iot/global"
	"fmt"
)

type CableTension struct {
	*IotDataBase
	*TagsBase
	CableTensionKey   uint8   `json:"cable_tension_key"`
	CableTensionUnit  string  `json:"cable_tension_unit"`
	CableTensionValue float32 `json:"cable_tension_value"`
}

func (m *CableTension) DatabaseName() string {
	return global.CONFIG.TDengine.Dbname
}

func (m *CableTension) StableName() string {
	return "cable_tension"
}

func (m *CableTension) TableName() string {
	return "ct_" + m.IotId[0:16]
}

// IotDataBaseColString 数据的基础列名字符串
func (m *CableTension) IotDataBaseColString() string {
	return "ts, request_id, message_id, yyyy, mm, dd, hh"
}

// IotDataBaseValString 数据的基础列值字符串
func (m *CableTension) IotDataBaseValString() string {
	// 与IotDataBaseColString中的列一一对应
	return fmt.Sprintf("%d, %s, %d, '%s', '%s', '%s', '%s'", m.Ts, m.RequestId, m.MessageId, m.Yyyy, m.Mm, m.Dd, m.Hh)
}

// BizColString 数据的业务列名字符串
func (m *CableTension) BizColString() string {
	return "cable_tension_key, cable_tension_unit, cable_tension_value"
}

// BizValString 数据的业务列值字符串
func (m *CableTension) BizValString() string {
	// 与BizColString中的列一一对应
	return fmt.Sprintf("'%d', '%s', %f", m.CableTensionKey, m.CableTensionUnit, m.CableTensionValue)
}

// TagColString 表的Tag列名字符串
func (m *CableTension) TagColString() string {
	return "product_key, iot_id, device_name, ym"
}

// TagValString 表的Tag列值字符串
func (m *CableTension) TagValString() string {
	// TagColString
	return fmt.Sprintf("'%s', '%s', '%s', '%s'", m.ProductKey, m.IotId, m.DeviceName, m.Ym)
}
