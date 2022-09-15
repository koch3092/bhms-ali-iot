package model

type TagsBase struct {
	ProductKey string `json:"product_key"`
	DeviceName string `json:"device_name"`
	IotId      string `json:"iot_id"`
	Ym         string `json:"ym"`
}

type TdMetricBase struct {
	Ts   uint64 `json:"ts"`
	Yyyy string `json:"yyyy"`
	Mm   string `json:"mm"`
	Dd   string `json:"dd"`
	Hh   string `json:"hh"`
}

type MetricsBase struct {
	Dt         uint64 `json:"dt"`          // 度量时间，使用时间戳表示
	MetricType uint8  `json:"metric_type"` // 度量类型，如1表示bridge_deck_temp等
	MetricNo   string `json:"metric_no"`   // 同类度量的位置序号，如1表示temperature1等
	FieldName  string `json:"field_name"`  // 度量的名称，如temperature/humidity等
	FieldValue string `json:"field_value"` // 度量值
	FieldUnit  string `json:"field_unit"`  // 度量单位
}

type IotDataBase struct {
	*TdMetricBase
	RequestId string `json:"request_id"`
	MessageId uint64 `json:"message_id"`
}

type AlarmBase struct {
	AlarmLevel  int8    `json:"alarm_level"`
	AlarmCordon float32 `json:"alarm_cordon"`
}

type DataType struct {
	LatestDataType uint8 `json:"latest_data_type"`
}

type MessageBase struct {
	MessageId    uint64 `json:"message_id"`
	Topic        string `json:"topic"`
	GenerateTime uint32 `json:"generate_time"`
	Qos          uint8  `json:"qos"`
}
