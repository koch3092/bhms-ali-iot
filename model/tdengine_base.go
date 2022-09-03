package model

type TagsBase struct {
	ProductKey string `json:"product_key"`
	DeviceName string `json:"device_name"`
	IotId      string `json:"iot_id"`
	Ym         string `json:"ym"`
}

type IotDataBase struct {
	Ts        uint64 `json:"ts"`
	RequestId string `json:"request_id"`
	Yyyy      string `json:"yyyy"`
	Mm        string `json:"mm"`
	Dd        string `json:"dd"`
	Hh        string `json:"hh"`
}

type DataType struct {
	LatestDataType uint8 `json:"latest_data_type"`
}
