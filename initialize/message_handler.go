package initialize

import (
	"bhms-ali-iot/global"
	"bhms-ali-iot/model"
	"bhms-ali-iot/service"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"pack.ag/amqp"
)

type MessageHandler struct {
	Logger *zap.Logger
}

func (h MessageHandler) Handle(ctx context.Context, rcvMessage <-chan *amqp.Message) {
	h.Logger.Info("Message handler init success")
	// TDEngine初始化
	session, err := TdengineSession()
	if err != nil {
		panic(any(err))
	}
	defer func() {
		err := session.Close()
		if err != nil {
			return
		}
	}()
ProcessMessage:
	for {
		select {
		case message := <-rcvMessage:
			h.Logger.Info(fmt.Sprintf("data received: %s properties: %#v", string(message.GetData()), message.ApplicationProperties))
			var dataType model.DataType
			errJson := json.Unmarshal(message.GetData(), &dataType)
			if errJson != nil {
				h.Logger.DPanic(fmt.Sprintf("Read data type error: %s", errJson.Error()))
			}

			var errM, errC error
			switch dataType.LatestDataType {
			case 1: // 桥面温度
				var m model.BridgeDeckTemp
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Bridge Deck Temp error: %s", errM.Error()))
				}
				h.Logger.Debug(fmt.Sprintf("iotDataBase: %#v, tagsBase: %#v, bdt: %#v", m.IotDataBase, m.TagsBase, m))
				ms := service.BridgeDeckTempService{Logger: global.Logger}
				errC = ms.CreateBridgeDeckTemp(session, &m)
			case 2: // 环境温湿度
				var m model.AmbientTempHumidity
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Ambient Temp Humidity error. %s", errM.Error()))
				}
				ms := service.AmbientTempHumidityService{Logger: global.Logger}
				errC = ms.CreateAmbientTempHumidity(session, &m)
			case 3: // 主梁挠度
				var m model.Deflection
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Deflection error. %s", errM.Error()))
				}
				ms := service.DeflectionService{Logger: global.Logger}
				errC = ms.CreateDeflection(session, &m)
			case 4: // 索力
				var m model.CableTension
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Cable Tension error. %s", errM.Error()))
				}
				ms := service.CableTensionService{Logger: global.Logger}
				errC = ms.CreateCableTension(session, &m)
			case 5: // 静应变
				var m model.StaticStrain
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Static Strain error. %s", errM.Error()))
				}
				ms := service.StaticStrainService{Logger: global.Logger}
				errC = ms.CreateStaticStrain(session, &m)
			case 6: // 地震
				var m model.Seismic
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Seismic error. %s", errM.Error()))
				}
				ms := service.SeismicService{Logger: global.Logger}
				errC = ms.CreateSeismic(session, &m)
			case 7: // 车道
				var m model.Driveway
				errM = json.Unmarshal(message.GetData(), &m)
				if errM != nil {
					h.Logger.DPanic(fmt.Sprintf("Read Driveway error. %s", errM.Error()))
				}
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
