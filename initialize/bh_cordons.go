package initialize

import (
	"bhms-ali-iot/global"
	"errors"
)

func InitCordons() error {
	c := global.CONFIG.Cordons
	if c.BridgeDeckTemp1 == 0 || c.BridgeDeckTemp2 == 0 {
		return errors.New("bridge deck temperature cordon must be set")
	}
	if c.AmbientTemp1 == 0 || c.AmbientTemp2 == 0 {
		return errors.New("ambient temperature cordon must be set")
	}
	if c.AmbientHumidity1 == 0 || c.AmbientHumidity2 == 0 {
		return errors.New("ambient humidity cordon must be set")
	}
	if c.Deflection1 == 0 || c.Deflection2 == 0 {
		return errors.New("deflection cordon must be set")
	}
	if c.CableTension1 == 0 || c.CableTension2 == 0 {
		return errors.New("cable tension cordon must be set")
	}
	if c.StaticStrainTemp1 == 0 || c.StaticStrainTemp2 == 0 {
		return errors.New("static strain temperature cordon must be set")
	}
	if c.StaticStrainValue1 == 0 || c.StaticStrainValue2 == 0 {
		return errors.New("static strain value cordon must be set")
	}
	if c.SeismicXValue1 == 0 || c.SeismicXValue2 == 0 {
		return errors.New("seismic x value cordon must be set")
	}
	if c.SeismicZValue1 == 0 || c.SeismicZValue2 == 0 {
		return errors.New("seismic z value cordon must be set")
	}
	if c.DrivewayWeight1 == 0 || c.DrivewayWeight2 == 0 {
		return errors.New("driveway weight cordon must be set")
	}
	if c.DrivewaySpeed1 == 0 || c.DrivewaySpeed2 == 0 {
		return errors.New("driveway speed cordon must be set")
	}
	if c.BridgeDeckTemp1 == 0 || c.BridgeDeckTemp2 == 0 {
		return errors.New("bridge deck temperature cordon must be set")
	}

	return nil
}
