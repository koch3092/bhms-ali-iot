package initialize

import (
	"bhms-ali-iot/global"
	"bhms-ali-iot/model"
	"bhms-ali-iot/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"pack.ag/amqp"
	"strconv"
	"time"
)

type MessageHandler struct {
	Logger *zap.Logger
}

func (h MessageHandler) HandleSaveData(ctx context.Context, rcvMessage <-chan *amqp.Message) {
	h.Logger.Info("Save data handler init success")
	// TDEngine初始化
	session := global.TDengine
	err := session.Ping()
	if err != nil {
		panic(err)
	} else {
		h.Logger.Info("Session init OK.")
	}
ProcessMessage:
	for {
		select {
		case message := <-rcvMessage:
			var dataType model.DataType
			errJson := json.Unmarshal(message.GetData(), &dataType)
			if errJson != nil {
				h.Logger.DPanic(fmt.Sprintf("Read data type error: %s", errJson.Error()))
			}

			var messageBase model.MessageBase
			errMessageBase := mapstructure.Decode(message.ApplicationProperties, &messageBase)
			if errMessageBase != nil {
				h.Logger.DPanic(fmt.Sprintf("Read message base info error: %s", errJson.Error()))
			}

			var errM, errC error
			switch dataType.LatestDataType {
			case 1: // 桥面温度
				var m model.BridgeDeckTemp
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Bridge Deck Temp error: %s", errM.Error()))
				}
				m.MessageId = messageBase.MessageId
				ms := service.BridgeDeckTempService{Logger: global.Logger}
				errC = ms.CreateBridgeDeckTemp(session, &m)
			case 2: // 环境温湿度
				var m model.AmbientTempHumidity
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Ambient Temp Humidity error. %s", errM.Error()))
				}
				m.MessageId = messageBase.MessageId
				ms := service.AmbientTempHumidityService{Logger: global.Logger}
				errC = ms.CreateAmbientTempHumidity(session, &m)
			case 3: // 主梁挠度
				var m model.Deflection
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Deflection error. %s", errM.Error()))
				}
				m.MessageId = messageBase.MessageId
				ms := service.DeflectionService{Logger: global.Logger}
				errC = ms.CreateDeflection(session, &m)
			case 4: // 索力
				var m model.CableTension
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Cable Tension error. %s", errM.Error()))
				}
				m.MessageId = messageBase.MessageId
				ms := service.CableTensionService{Logger: global.Logger}
				errC = ms.CreateCableTension(session, &m)
			case 5: // 静应变
				var m model.StaticStrain
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Static Strain error. %s", errM.Error()))
				}
				m.MessageId = messageBase.MessageId
				ms := service.StaticStrainService{Logger: global.Logger}
				errC = ms.CreateStaticStrain(session, &m)
			case 6: // 地震
				var m model.Seismic
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Seismic error. %s", errM.Error()))
				}
				m.MessageId = messageBase.MessageId
				ms := service.SeismicService{Logger: global.Logger}
				errC = ms.CreateSeismic(session, &m)
			case 7: // 车道
				var m model.Driveway
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Driveway error. %s", errM.Error()))
				}
				m.MessageId = messageBase.MessageId
				ms := service.DrivewayService{Logger: global.Logger}
				errC = ms.CreateDriveway(session, &m)
			}

			if errC != nil {
				h.Logger.DPanic(errC.Error())
				panic(any(errC))
			}
		case <-ctx.Done():
			break ProcessMessage
		}
	}
}

func (h MessageHandler) HandleMeasurement(ctx context.Context, rcvMessage <-chan *amqp.Message) {
	h.Logger.Info("Measurement handler init success")
	// TDEngine初始化
	session := global.TDengine
	err := session.Ping()
	if err != nil {
		panic(err)
	} else {
		h.Logger.Info("Session init OK.")
	}
ProcessMessage:
	for {
		select {
		case message := <-rcvMessage:
			// 读取公共数据
			var tdMetricBase *model.TdMetricBase
			errTd := json.Unmarshal(message.GetData(), &tdMetricBase)
			if errTd != nil {
				h.Logger.DPanic(fmt.Sprintf("Read td metric base error: %s", errTd.Error()))
			}

			var tagsBase *model.TagsBase
			errTags := json.Unmarshal(message.GetData(), &tagsBase)
			if errTags != nil {
				h.Logger.DPanic(fmt.Sprintf("Read td metric base error: %s", errTags.Error()))
			}

			var dataType model.DataType
			errJson := json.Unmarshal(message.GetData(), &dataType)
			if errJson != nil {
				h.Logger.DPanic(fmt.Sprintf("Read data type error: %s", errJson.Error()))
			}

			// 初始化measurement实例对象
			mes := &model.Measurement{
				TdMetricBase: tdMetricBase,
				MetricsBase: &model.MetricsBase{
					Dt:         uint64(time.Now().UnixMilli()),
					MetricType: dataType.LatestDataType,
					MetricNo:   "",
					FieldName:  "",
					FieldValue: "",
					FieldUnit:  "",
				},
				AlarmBase: &model.AlarmBase{
					AlarmLevel:  0,
					AlarmCordon: 0,
				},
				TagsBase: tagsBase,
			}

			// 拼接INSERT语句，各个数据不同，则按实际情况拼接
			batchInsertSql := fmt.Sprintf(
				"INSERT INTO %s(%s, %s, %s) USING %s TAGS (%s) VALUES",
				mes.TableName(), mes.TdMetricsBaseColString(), mes.MetricsBaseColString(),
				mes.AlarmColString(), mes.StableName(), mes.TagValString(),
			)
			insertSqlValues := ""

			var errM, errC error
			switch dataType.LatestDataType {
			case 1: // 桥面温度
				var m model.BridgeDeckTemp
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Bridge Deck Temp error: %s", errM.Error()))
				}

				mes.FieldName = "bd_temperature"
				mes.FieldUnit = m.TemperatureUnit
				mes.AlarmLevel = 0
				mes.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.Temperature1 > global.CONFIG.Cordons.BridgeDeckTemp2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.BridgeDeckTemp2
				} else if m.Temperature1 > global.CONFIG.Cordons.BridgeDeckTemp1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.BridgeDeckTemp1
				}
				mes.Ts = mes.Ts + 1
				mes.Dt = mes.Dt + 0
				mes.MetricNo = strconv.Itoa(1)
				mes.FieldValue = fmt.Sprintf("%f", m.Temperature1)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				if m.Temperature2 > global.CONFIG.Cordons.BridgeDeckTemp2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.BridgeDeckTemp2
				} else if m.Temperature2 > global.CONFIG.Cordons.BridgeDeckTemp1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.BridgeDeckTemp1
				}
				mes.Ts = mes.Ts + 2
				mes.Dt = mes.Dt + 1
				mes.MetricNo = strconv.Itoa(2)
				mes.FieldValue = fmt.Sprintf("%f", m.Temperature2)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				if m.Temperature3 > global.CONFIG.Cordons.BridgeDeckTemp2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.BridgeDeckTemp2
				} else if m.Temperature3 > global.CONFIG.Cordons.BridgeDeckTemp1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.BridgeDeckTemp1
				}
				mes.Ts = mes.Ts + 3
				mes.Dt = mes.Dt + 2
				mes.MetricNo = strconv.Itoa(3)
				mes.FieldValue = fmt.Sprintf("%f", m.Temperature3)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				if m.Temperature4 > global.CONFIG.Cordons.BridgeDeckTemp2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.BridgeDeckTemp2
				} else if m.Temperature4 > global.CONFIG.Cordons.BridgeDeckTemp1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.BridgeDeckTemp1
				}
				mes.Ts = mes.Ts + 4
				mes.Dt = mes.Dt + 3
				mes.MetricNo = strconv.Itoa(4)
				mes.FieldValue = fmt.Sprintf("%f", m.Temperature4)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())
				batchInsertSql += insertSqlValues
				ms := service.MeasurementService{Logger: global.Logger}
				errC = ms.CreateMeasurement(session, batchInsertSql)
			case 2: // 环境温湿度
				var m model.AmbientTempHumidity
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Ambient Temp Humidity error. %s", errM.Error()))
				}

				mes.FieldName = "ath_temperature"
				mes.FieldUnit = m.TemperatureUnit
				mes.AlarmLevel = 0
				mes.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.Temperature1 > global.CONFIG.Cordons.AmbientTemp2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.AmbientTemp2
				} else if m.Temperature1 > global.CONFIG.Cordons.AmbientTemp1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.AmbientTemp1
				}
				mes.Ts = mes.Ts + 1
				mes.Dt = mes.Dt + 0
				mes.MetricNo = strconv.Itoa(1)
				mes.FieldValue = fmt.Sprintf("%f", m.Temperature1)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				if m.Temperature2 > global.CONFIG.Cordons.AmbientTemp2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.AmbientTemp2
				} else if m.Temperature2 > global.CONFIG.Cordons.AmbientTemp1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.AmbientTemp1
				}
				mes.Ts = mes.Ts + 2
				mes.Dt = mes.Dt + 1
				mes.MetricNo = strconv.Itoa(2)
				mes.FieldValue = fmt.Sprintf("%f", m.Temperature2)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				if m.Temperature3 > global.CONFIG.Cordons.AmbientTemp2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.AmbientTemp2
				} else if m.Temperature3 > global.CONFIG.Cordons.AmbientTemp1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.AmbientTemp1
				}
				mes.Ts = mes.Ts + 3
				mes.Dt = mes.Dt + 2
				mes.MetricNo = strconv.Itoa(3)
				mes.FieldValue = fmt.Sprintf("%f", m.Temperature3)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				mes.FieldName = "ath_humidity"
				mes.FieldUnit = m.HumidityUnit
				mes.AlarmLevel = 0
				mes.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.Humidity1 > global.CONFIG.Cordons.AmbientHumidity2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.AmbientHumidity2
				} else if m.Humidity1 > global.CONFIG.Cordons.AmbientHumidity1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.AmbientHumidity1
				}
				mes.Ts = mes.Ts + 4
				mes.Dt = mes.Dt + 3
				mes.MetricNo = strconv.Itoa(1)
				mes.FieldValue = fmt.Sprintf("%f", m.Humidity1)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				if m.Humidity2 > global.CONFIG.Cordons.AmbientHumidity2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.AmbientHumidity2
				} else if m.Humidity2 > global.CONFIG.Cordons.AmbientHumidity1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.AmbientHumidity1
				}
				mes.Ts = mes.Ts + 5
				mes.Dt = mes.Dt + 4
				mes.MetricNo = strconv.Itoa(2)
				mes.FieldValue = fmt.Sprintf("%f", m.Humidity2)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				if m.Humidity3 > global.CONFIG.Cordons.AmbientHumidity2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.AmbientHumidity2
				} else if m.Humidity3 > global.CONFIG.Cordons.AmbientHumidity1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.AmbientHumidity1
				}
				mes.Ts = mes.Ts + 6
				mes.Dt = mes.Dt + 5
				mes.MetricNo = strconv.Itoa(3)
				mes.FieldValue = fmt.Sprintf("%f", m.Humidity3)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				batchInsertSql += insertSqlValues
				ms := service.MeasurementService{Logger: global.Logger}
				errC = ms.CreateMeasurement(session, batchInsertSql)
			case 3: // 主梁挠度
				var m model.Deflection
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Deflection error. %s", errM.Error()))
				}

				mes.FieldName = "deflection"
				mes.FieldUnit = m.DeflectionUnit
				mes.AlarmLevel = 0
				mes.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.Deflection1 > global.CONFIG.Cordons.Deflection2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.Deflection2
				} else if m.Deflection1 > global.CONFIG.Cordons.Deflection1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.Deflection1
				}
				mes.Ts = mes.Ts + 1
				mes.Dt = mes.Dt + 0
				mes.MetricNo = strconv.Itoa(1)
				mes.FieldValue = fmt.Sprintf("%f", m.Deflection1)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				if m.Deflection2 > global.CONFIG.Cordons.Deflection2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.Deflection2
				} else if m.Deflection2 > global.CONFIG.Cordons.Deflection1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.Deflection1
				}
				mes.Ts = mes.Ts + 2
				mes.Dt = mes.Dt + 1
				mes.MetricNo = strconv.Itoa(2)
				mes.FieldValue = fmt.Sprintf("%f", m.Deflection2)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				if m.Deflection3 > global.CONFIG.Cordons.Deflection2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.Deflection2
				} else if m.Deflection3 > global.CONFIG.Cordons.Deflection1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.Deflection1
				}
				mes.Ts = mes.Ts + 3
				mes.Dt = mes.Dt + 2
				mes.MetricNo = strconv.Itoa(3)
				mes.FieldValue = fmt.Sprintf("%f", m.Deflection3)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				if m.Deflection4 > global.CONFIG.Cordons.Deflection2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.Deflection2
				} else if m.Deflection4 > global.CONFIG.Cordons.Deflection1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.Deflection1
				}
				mes.Ts = mes.Ts + 4
				mes.Dt = mes.Dt + 3
				mes.MetricNo = strconv.Itoa(4)
				mes.FieldValue = fmt.Sprintf("%f", m.Deflection4)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				if m.Deflection5 > global.CONFIG.Cordons.Deflection2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.Deflection2
				} else if m.Deflection5 > global.CONFIG.Cordons.Deflection1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.Deflection1
				}
				mes.Ts = mes.Ts + 5
				mes.Dt = mes.Dt + 4
				mes.MetricNo = strconv.Itoa(5)
				mes.FieldValue = fmt.Sprintf("%f", m.Deflection5)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				if m.Deflection6 > global.CONFIG.Cordons.Deflection2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.Deflection2
				} else if m.Deflection6 > global.CONFIG.Cordons.Deflection1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.Deflection1
				}
				mes.Ts = mes.Ts + 6
				mes.Dt = mes.Dt + 5
				mes.MetricNo = strconv.Itoa(6)
				mes.FieldValue = fmt.Sprintf("%f", m.Deflection6)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				batchInsertSql += insertSqlValues
				ms := service.MeasurementService{Logger: global.Logger}
				errC = ms.CreateMeasurement(session, batchInsertSql)
			case 4: // 索力
				var m model.CableTension
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Cable Tension error. %s", errM.Error()))
				}
				mes.FieldName = "cable_tension"
				mes.FieldUnit = m.CableTensionUnit
				mes.AlarmLevel = 0
				mes.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.CableTensionValue > global.CONFIG.Cordons.CableTension2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.CableTension2
				} else if m.CableTensionValue > global.CONFIG.Cordons.CableTension1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.CableTension1
				}
				mes.MetricNo = fmt.Sprintf("%d", m.CableTensionKey)
				mes.FieldValue = fmt.Sprintf("%f", m.CableTensionValue)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				batchInsertSql += insertSqlValues
				ms := service.MeasurementService{Logger: global.Logger}
				errC = ms.CreateMeasurement(session, batchInsertSql)
			case 5: // 静应变
				var m model.StaticStrain
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Static Strain error. %s", errM.Error()))
				}
				mes.FieldName = "ss_temperature"
				mes.FieldUnit = m.SSTemperatureUnit
				mes.AlarmLevel = 0
				mes.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.SSTemperatureValue > global.CONFIG.Cordons.StaticStrainTemp2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.StaticStrainTemp2
				} else if m.SSTemperatureValue > global.CONFIG.Cordons.StaticStrainTemp1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.StaticStrainTemp1
				}
				mes.Ts = mes.Ts + 1
				mes.Dt = mes.Dt + 0
				mes.MetricNo = fmt.Sprintf("%d", m.StaticStrainKey)
				mes.FieldValue = fmt.Sprintf("%f", m.SSTemperatureValue)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				mes.FieldName = "ss_strain"
				mes.FieldUnit = m.SSStrainUnit
				mes.AlarmLevel = 0
				mes.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.SSStrainValue > global.CONFIG.Cordons.StaticStrainValue2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.StaticStrainValue2
				} else if m.SSStrainValue > global.CONFIG.Cordons.StaticStrainValue1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.StaticStrainValue1
				}
				mes.Ts = mes.Ts + 2
				mes.Dt = mes.Dt + 1
				mes.MetricNo = fmt.Sprintf("%d", m.StaticStrainKey)
				mes.FieldValue = fmt.Sprintf("%f", m.SSStrainValue)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				batchInsertSql += insertSqlValues
				ms := service.MeasurementService{Logger: global.Logger}
				errC = ms.CreateMeasurement(session, batchInsertSql)
			case 6: // 地震
				var m model.Seismic
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Seismic error. %s", errM.Error()))
				}
				mes.FieldName = "seismic_x_axis"
				mes.FieldUnit = m.SeismicXUnit
				mes.AlarmLevel = 0
				mes.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.SeismicXValue > global.CONFIG.Cordons.SeismicXValue2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.SeismicXValue2
				} else if m.SeismicXValue > global.CONFIG.Cordons.SeismicXValue1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.SeismicXValue1
				}
				mes.Ts = mes.Ts + 1
				mes.Dt = mes.Dt + 0
				mes.MetricNo = fmt.Sprintf("%d", m.SeismicKey)
				mes.FieldValue = fmt.Sprintf("%f", m.SeismicXValue)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				mes.FieldName = "seismic_z_axis"
				mes.FieldUnit = m.SeismicZUnit
				mes.AlarmLevel = 0
				mes.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.SeismicZValue > global.CONFIG.Cordons.SeismicZValue2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.SeismicZValue2
				} else if m.SeismicZValue > global.CONFIG.Cordons.SeismicZValue1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.SeismicZValue1
				}
				mes.Ts = mes.Ts + 2
				mes.Dt = mes.Dt + 1
				mes.MetricNo = fmt.Sprintf("%d", m.SeismicKey)
				mes.FieldValue = fmt.Sprintf("%f", m.SeismicZValue)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				batchInsertSql += insertSqlValues
				ms := service.MeasurementService{Logger: global.Logger}
				errC = ms.CreateMeasurement(session, batchInsertSql)
			case 7: // 车道
				var m model.Driveway
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Driveway error. %s", errM.Error()))
				}
				mes.FieldName = "driveway_weight"
				mes.FieldUnit = m.DrivewayWeightUnit
				mes.AlarmLevel = 0
				mes.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.DrivewayWeightValue > global.CONFIG.Cordons.DrivewayWeight2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.DrivewayWeight2
				} else if m.DrivewayWeightValue > global.CONFIG.Cordons.DrivewayWeight1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.DrivewayWeight1
				}
				mes.Ts = mes.Ts + 1
				mes.Dt = mes.Dt + 0
				mes.MetricNo = fmt.Sprintf("%d", m.DrivewayKey)
				mes.FieldValue = fmt.Sprintf("%f", m.DrivewayWeightValue)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				mes.FieldName = "driveway_speed"
				mes.FieldUnit = m.DrivewaySpeedUnit
				mes.AlarmLevel = 0
				mes.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.DrivewaySpeedValue > global.CONFIG.Cordons.DrivewaySpeed2 {
					mes.AlarmLevel = 2
					mes.AlarmCordon = global.CONFIG.Cordons.DrivewaySpeed2
				} else if m.DrivewaySpeedValue > global.CONFIG.Cordons.DrivewaySpeed1 {
					mes.AlarmLevel = 1
					mes.AlarmCordon = global.CONFIG.Cordons.DrivewaySpeed1
				}
				mes.Ts = mes.Ts + 2
				mes.Dt = mes.Dt + 1
				mes.MetricNo = fmt.Sprintf("%d", m.DrivewayKey)
				mes.FieldValue = fmt.Sprintf("%f", m.DrivewaySpeedValue)
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				mes.FieldName = "driveway_model"
				mes.FieldUnit = m.DrivewayModelUnit
				mes.AlarmLevel = 0
				mes.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				mes.Ts = mes.Ts + 3
				mes.Dt = mes.Dt + 2
				mes.MetricNo = fmt.Sprintf("%d", m.DrivewayKey)
				mes.FieldValue = fmt.Sprintf("%d", m.DrivewayModelValue)
				mes.FieldUnit = m.DrivewayModelLabel
				insertSqlValues += fmt.Sprintf(" (%s, %s, %s)", mes.TdMetricsBaseValString(), mes.MetricsBaseValString(), mes.AlarmValString())

				batchInsertSql += insertSqlValues
				ms := service.MeasurementService{Logger: global.Logger}
				errC = ms.CreateMeasurement(session, batchInsertSql)
			}

			if errC != nil {
				h.Logger.DPanic(errC.Error())
				panic(any(errC))
			}
		case <-ctx.Done():
			break ProcessMessage
		}
	}
}

func (h MessageHandler) HandleAlarm(ctx context.Context, rcvMessage <-chan *amqp.Message) {
	h.Logger.Info("Alarm handler init success")
	// Redis初始化
	rdb := global.Redis
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	} else {
		h.Logger.Info("Redis session init OK.")
	}
ProcessMessage:
	for {
		select {
		case message := <-rcvMessage:
			// 读取公共数据
			var dataType model.DataType
			errJson := json.Unmarshal(message.GetData(), &dataType)
			if errJson != nil {
				h.Logger.DPanic(fmt.Sprintf("Read data type error: %s", errJson.Error()))
			}

			// 初始化measurement实例对象
			alarm := &model.Alarm{
				MetricsBase: &model.MetricsBase{
					Dt:         uint64(time.Now().UnixMilli()),
					MetricType: dataType.LatestDataType,
					MetricNo:   "",
					FieldName:  "",
					FieldValue: "",
					FieldUnit:  "",
				},
				AlarmBase: &model.AlarmBase{
					AlarmLevel:  0,
					AlarmCordon: 0,
				},
			}

			// 初始化redis的pipeline实例对象
			pipe := rdb.TxPipeline()

			var errM, errC error
			switch dataType.LatestDataType {
			case 1: // 桥面温度
				var m model.BridgeDeckTemp
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Bridge Deck Temp error: %s", errM.Error()))
				}

				alarm.FieldName = "bd_temperature"
				alarm.FieldUnit = m.TemperatureUnit
				alarm.AlarmLevel = 0
				alarm.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.Temperature1 > global.CONFIG.Cordons.BridgeDeckTemp2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.BridgeDeckTemp2
					alarm.MetricNo = strconv.Itoa(1)
					alarm.FieldValue = fmt.Sprintf("%f", m.Temperature1)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}

				if m.Temperature2 > global.CONFIG.Cordons.BridgeDeckTemp2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.BridgeDeckTemp2
					alarm.MetricNo = strconv.Itoa(2)
					alarm.FieldValue = fmt.Sprintf("%f", m.Temperature2)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}

				if m.Temperature3 > global.CONFIG.Cordons.BridgeDeckTemp2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.BridgeDeckTemp2
					alarm.MetricNo = strconv.Itoa(3)
					alarm.FieldValue = fmt.Sprintf("%f", m.Temperature3)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}

				if m.Temperature4 > global.CONFIG.Cordons.BridgeDeckTemp2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.BridgeDeckTemp2
					alarm.MetricNo = strconv.Itoa(4)
					alarm.FieldValue = fmt.Sprintf("%f", m.Temperature4)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}
			case 2: // 环境温湿度
				var m model.AmbientTempHumidity
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Ambient Temp Humidity error. %s", errM.Error()))
				}

				alarm.FieldName = "ath_temperature"
				alarm.FieldUnit = m.TemperatureUnit
				alarm.AlarmLevel = 0
				alarm.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.Temperature1 > global.CONFIG.Cordons.AmbientTemp2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.AmbientTemp2
					alarm.MetricNo = strconv.Itoa(1)
					alarm.FieldValue = fmt.Sprintf("%f", m.Temperature1)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}

				if m.Temperature2 > global.CONFIG.Cordons.AmbientTemp2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.AmbientTemp2
					alarm.MetricNo = strconv.Itoa(2)
					alarm.FieldValue = fmt.Sprintf("%f", m.Temperature2)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}

				if m.Temperature3 > global.CONFIG.Cordons.AmbientTemp2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.AmbientTemp2
					alarm.MetricNo = strconv.Itoa(3)
					alarm.FieldValue = fmt.Sprintf("%f", m.Temperature3)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}

				alarm.FieldName = "ath_humidity"
				alarm.FieldUnit = m.HumidityUnit
				alarm.AlarmLevel = 0
				alarm.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.Humidity1 > global.CONFIG.Cordons.AmbientHumidity2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.AmbientHumidity2
					alarm.MetricNo = strconv.Itoa(1)
					alarm.FieldValue = fmt.Sprintf("%f", m.Humidity1)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}

				if m.Humidity2 > global.CONFIG.Cordons.AmbientHumidity2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.AmbientHumidity2
					alarm.MetricNo = strconv.Itoa(2)
					alarm.FieldValue = fmt.Sprintf("%f", m.Humidity2)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}

				if m.Humidity3 > global.CONFIG.Cordons.AmbientHumidity2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.AmbientHumidity2
					alarm.MetricNo = strconv.Itoa(3)
					alarm.FieldValue = fmt.Sprintf("%f", m.Humidity3)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}
			case 3: // 主梁挠度
				var m model.Deflection
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Deflection error. %s", errM.Error()))
				}

				alarm.FieldName = "deflection"
				alarm.FieldUnit = m.DeflectionUnit
				alarm.AlarmLevel = 0
				alarm.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.Deflection1 > global.CONFIG.Cordons.Deflection2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.Deflection2
					alarm.MetricNo = strconv.Itoa(1)
					alarm.FieldValue = fmt.Sprintf("%f", m.Deflection1)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}

				if m.Deflection2 > global.CONFIG.Cordons.Deflection2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.Deflection2
					alarm.MetricNo = strconv.Itoa(2)
					alarm.FieldValue = fmt.Sprintf("%f", m.Deflection2)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}

				if m.Deflection3 > global.CONFIG.Cordons.Deflection2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.Deflection2
					alarm.MetricNo = strconv.Itoa(3)
					alarm.FieldValue = fmt.Sprintf("%f", m.Deflection3)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}

				if m.Deflection4 > global.CONFIG.Cordons.Deflection2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.Deflection2
					alarm.MetricNo = strconv.Itoa(4)
					alarm.FieldValue = fmt.Sprintf("%f", m.Deflection4)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}

				if m.Deflection5 > global.CONFIG.Cordons.Deflection2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.Deflection2
					alarm.MetricNo = strconv.Itoa(5)
					alarm.FieldValue = fmt.Sprintf("%f", m.Deflection5)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}

				if m.Deflection6 > global.CONFIG.Cordons.Deflection2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.Deflection2
					alarm.MetricNo = strconv.Itoa(6)
					alarm.FieldValue = fmt.Sprintf("%f", m.Deflection6)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}
			case 4: // 索力
				var m model.CableTension
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Cable Tension error. %s", errM.Error()))
				}
				alarm.FieldName = "cable_tension"
				alarm.FieldUnit = m.CableTensionUnit
				alarm.AlarmLevel = 0
				alarm.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.CableTensionValue > global.CONFIG.Cordons.CableTension2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.CableTension2
					alarm.MetricNo = strconv.Itoa(int(m.CableTensionKey))
					alarm.FieldValue = fmt.Sprintf("%f", m.CableTensionValue)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}
			case 5: // 静应变
				var m model.StaticStrain
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Static Strain error. %s", errM.Error()))
				}
				alarm.FieldName = "ss_temperature"
				alarm.FieldUnit = m.SSTemperatureUnit
				alarm.AlarmLevel = 0
				alarm.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.SSTemperatureValue > global.CONFIG.Cordons.StaticStrainTemp2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.StaticStrainTemp2
					alarm.MetricNo = strconv.Itoa(int(m.StaticStrainKey))
					alarm.FieldValue = fmt.Sprintf("%f", m.SSTemperatureValue)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}

				alarm.FieldName = "ss_strain"
				alarm.FieldUnit = m.SSStrainUnit
				alarm.AlarmLevel = 0
				alarm.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.SSStrainValue > global.CONFIG.Cordons.StaticStrainValue2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.StaticStrainValue2
					alarm.MetricNo = strconv.Itoa(int(m.StaticStrainKey))
					alarm.FieldValue = fmt.Sprintf("%f", m.SSStrainValue)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}
			case 6: // 地震
				var m model.Seismic
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Seismic error. %s", errM.Error()))
				}
				alarm.FieldName = "seismic_x_axis"
				alarm.FieldUnit = m.SeismicXUnit
				alarm.AlarmLevel = 0
				alarm.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.SeismicXValue > global.CONFIG.Cordons.SeismicXValue2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.SeismicXValue2
					alarm.MetricNo = strconv.Itoa(int(m.SeismicKey))
					alarm.FieldValue = fmt.Sprintf("%f", m.SeismicXValue)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}

				alarm.FieldName = "seismic_z_axis"
				alarm.FieldUnit = m.SeismicZUnit
				alarm.AlarmLevel = 0
				alarm.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.SeismicZValue > global.CONFIG.Cordons.SeismicZValue2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.SeismicZValue2
					alarm.MetricNo = strconv.Itoa(int(m.SeismicKey))
					alarm.FieldValue = fmt.Sprintf("%f", m.SeismicZValue)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}
			case 7: // 车道
				var m model.Driveway
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Driveway error. %s", errM.Error()))
				}
				alarm.FieldName = "driveway_weight"
				alarm.FieldUnit = m.DrivewayWeightUnit
				alarm.AlarmLevel = 0
				alarm.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.DrivewayWeightValue > global.CONFIG.Cordons.DrivewayWeight2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.DrivewayWeight2
					alarm.MetricNo = strconv.Itoa(int(m.DrivewayKey))
					alarm.FieldValue = fmt.Sprintf("%f", m.DrivewayWeightValue)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}

				alarm.FieldName = "driveway_speed"
				alarm.FieldUnit = m.DrivewaySpeedUnit
				alarm.AlarmLevel = 0
				alarm.AlarmCordon = 0
				// 对监控点根据实际情况赋值
				if m.DrivewaySpeedValue > global.CONFIG.Cordons.DrivewaySpeed2 {
					alarm.AlarmLevel = 2
					alarm.AlarmCordon = global.CONFIG.Cordons.DrivewaySpeed2
					alarm.MetricNo = strconv.Itoa(int(m.DrivewayKey))
					alarm.FieldValue = fmt.Sprintf("%f", m.DrivewaySpeedValue)
					key := fmt.Sprintf("%d_%s_%s", alarm.MetricType, alarm.MetricNo, alarm.FieldName)
					// 读取记录，看是否已经告警过
					existAlarm, _ := rdb.HGetAll(ctx, key).Result()
					if len(existAlarm) == 0 { // 如果没有数据，则记录时间
						pipe.HSet(ctx, key, "dt", alarm.Dt)
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					} else if len(existAlarm) > 1 { // 有记录，则看持续时间，如果时间持续超过了2分钟，则报警
						pipe.HSet(ctx, key, alarm.FieldName, alarm.FieldValue)
						pipe.HSet(ctx, key, "level", alarm.AlarmLevel)
						pipe.HSet(ctx, key, "cordon", alarm.AlarmCordon)
						pipe.Expire(ctx, key, time.Second*90)
					}
				}
			}
			// TODO：收集数据，发送短信，短信需要有记录到TDengine的特定表中
			_, err := pipe.Exec(ctx)
			if err != nil {
				return
			}
			if errC != nil {
				h.Logger.DPanic(errC.Error())
				panic(errC)
			}
		case <-ctx.Done():
			break ProcessMessage
		}
	}
}
