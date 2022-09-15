package model

import (
	"bhms-ali-iot/global"
	"fmt"
)

type BridgeDeckTemp struct {
	*IotDataBase
	*TagsBase
	TemperatureUnit string  `json:"temperature_unit"`
	Temperature1    float32 `json:"temperature1"`
	Temperature2    float32 `json:"temperature2"`
	Temperature3    float32 `json:"temperature3"`
	Temperature4    float32 `json:"temperature4"`
}

func (m *BridgeDeckTemp) DatabaseName() string {
	return global.CONFIG.TDengine.Dbname
}

func (m *BridgeDeckTemp) StableName() string {
	return "bridge_deck_temperature"
}

func (m *BridgeDeckTemp) TableName() string {
	return "bdt_" + m.IotId[0:16]
}

// IotDataBaseColString 数据的基础列名字符串
func (m *BridgeDeckTemp) IotDataBaseColString() string {
	return "ts, yyyy, mm, dd, hh, request_id, message_id"
}

// IotDataBaseValString 数据的基础列值字符串
func (m *BridgeDeckTemp) IotDataBaseValString() string {
	// 与IotDataBaseColString中的列一一对应
	return fmt.Sprintf("%d, '%s', '%s', '%s', '%s', %s, %d", m.Ts, m.Yyyy, m.Mm, m.Dd, m.Hh, m.RequestId, m.MessageId)
}

// BizColString 数据的业务列名字符串
func (m *BridgeDeckTemp) BizColString() string {
	return "temperature_unit, temperature1, temperature2, temperature3, temperature4"
}

// BizValString 数据的业务列值字符串
func (m *BridgeDeckTemp) BizValString() string {
	// 与BizColString中的列一一对应
	return fmt.Sprintf("'%s', %f, %f, %f, %f", m.TemperatureUnit, m.Temperature1, m.Temperature2, m.Temperature3, m.Temperature4)
}

// TagColString 表的Tag列名字符串
func (m *BridgeDeckTemp) TagColString() string {
	return "product_key, iot_id, device_name, ym"
}

// TagValString 表的Tag列值字符串
func (m *BridgeDeckTemp) TagValString() string {
	// TagColString
	return fmt.Sprintf("'%s', '%s', '%s', '%s'", m.ProductKey, m.IotId, m.DeviceName, m.Ym)
}
