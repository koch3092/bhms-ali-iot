package model

import (
	"bhms-ali-iot/global"
	"fmt"
)

type AmbientTempHumidity struct {
	*IotDataBase
	*TagsBase
	TemperatureUnit string  `json:"temperature_unit"`
	HumidityUnit    string  `json:"humidity_unit"`
	Temperature1    float32 `json:"temperature1"`
	Humidity1       float32 `json:"humidity1"`
	Temperature2    float32 `json:"temperature2"`
	Humidity2       float32 `json:"humidity2"`
	Temperature3    float32 `json:"temperature3"`
	Humidity3       float32 `json:"humidity3"`
}

func (m *AmbientTempHumidity) DatabaseName() string {
	return global.CONFIG.TDengine.Dbname
}

func (m *AmbientTempHumidity) StableName() string {
	return "ambient_temperature_humidity"
}

func (m *AmbientTempHumidity) TableName() string {
	return "ath_" + m.IotId[0:16]
}

// IotDataBaseColString 数据的基础列名字符串
func (m *AmbientTempHumidity) IotDataBaseColString() string {
	return "ts, request_id, yyyy, mm, dd, hh"
}

// IotDataBaseValString 数据的基础列值字符串
func (m *AmbientTempHumidity) IotDataBaseValString() string {
	// 与IotDataBaseColString中的列一一对应
	return fmt.Sprintf("%d, %s, '%s', '%s', '%s', '%s'", m.Ts, m.RequestId, m.Yyyy, m.Mm, m.Dd, m.Hh)
}

// BizColString 数据的业务列名字符串
func (m *AmbientTempHumidity) BizColString() string {
	return "temperature_unit, humidity_unit, temperature1, humidity1, temperature2, humidity2, temperature3, humidity3"
}

// BizValString 数据的业务列值字符串
func (m *AmbientTempHumidity) BizValString() string {
	// 与BizColString中的列一一对应
	return fmt.Sprintf("'%s', '%s', %f, %f, %f, %f, %f, %f", m.TemperatureUnit, m.HumidityUnit, m.Temperature1, m.Humidity1, m.Temperature2, m.Humidity2, m.Temperature3, m.Humidity3)
}

// TagColString 表的Tag列名字符串
func (m *AmbientTempHumidity) TagColString() string {
	return "product_key, iot_id, device_name, ym"
}

// TagValString 表的Tag列值字符串
func (m *AmbientTempHumidity) TagValString() string {
	// TagColString
	return fmt.Sprintf("'%s', '%s', '%s', '%s'", m.ProductKey, m.IotId, m.DeviceName, m.Ym)
}
