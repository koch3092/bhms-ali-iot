package config

type Cordons struct {
	BridgeDeckTemp1    float32  `json:"bridge_deck_temp_1" yaml:"bridge-deck-temp-1" mapstructure:"bridge-deck-temp-1"`
	BridgeDeckTemp2    float32  `json:"bridge_deck_temp_2" yaml:"bridge-deck-temp-2" mapstructure:"bridge-deck-temp-2"`
	AmbientTemp1       float32  `json:"ambient_temp_1" yaml:"ambient-temp-1" mapstructure:"ambient-temp-1"`
	AmbientHumidity1   float32  `json:"ambient_humidity_1" yaml:"ambient-humidity-1" mapstructure:"ambient-humidity-1"`
	AmbientTemp2       float32  `json:"ambient_temp_2" yaml:"ambient-temp-2" mapstructure:"ambient-temp-2"`
	AmbientHumidity2   float32  `json:"ambient_humidity_2" yaml:"ambient-humidity-2" mapstructure:"ambient-humidity-2"`
	Deflection1        float32  `json:"deflection_1" yaml:"deflection-1" mapstructure:"deflection-1"`
	Deflection2        float32  `json:"deflection_2" yaml:"deflection-2" mapstructure:"deflection-2"`
	CableTension1      float32  `json:"cable_tension_1" yaml:"cable-tension-1" mapstructure:"cable-tension-1"`
	CableTension2      float32  `json:"cable_tension_2" yaml:"cable-tension-2" mapstructure:"cable-tension-2"`
	StaticStrainTemp1  float32  `json:"static_strain_temp_1" yaml:"static-strain-temp-1" mapstructure:"static-strain-temp-1"`
	StaticStrainTemp2  float32  `json:"static_strain_temp_2" yaml:"static-strain-temp-2" mapstructure:"static-strain-temp-2"`
	StaticStrainValue1 float32  `json:"static_strain_value_1" yaml:"static-strain-value-1" mapstructure:"static-strain-value-1"`
	StaticStrainValue2 float32  `json:"static_strain_value_2" yaml:"static-strain-value-2" mapstructure:"static-strain-value-2"`
	SeismicXValue1     float32  `json:"seismic_x_value_1" yaml:"seismic-x-value-1" mapstructure:"seismic-x-value-1"`
	SeismicXValue2     float32  `json:"seismic_x_value_2" yaml:"seismic-x-value-2" mapstructure:"seismic-x-value-2"`
	SeismicZValue1     float32  `json:"seismic_z_value_1" yaml:"seismic-z-value-1" mapstructure:"seismic-z-value-1"`
	SeismicZValue2     float32  `json:"seismic_z_value_2" yaml:"seismic-z-value-2" mapstructure:"seismic-z-value-2"`
	DrivewayWeight1    float32  `json:"driveway_weight_1" yaml:"driveway-weight-1" mapstructure:"driveway-weight-1"`
	DrivewayWeight2    float32  `json:"driveway_weight_2" yaml:"driveway-weight-2" mapstructure:"driveway-weight-2"`
	DrivewaySpeed1     float32  `json:"driveway_speed_1" yaml:"driveway-speed-1" mapstructure:"driveway-speed-1"`
	DrivewaySpeed2     float32  `json:"driveway_speed_2" yaml:"driveway-speed-2" mapstructure:"driveway-speed-2"`
	AlarmContacts      []string `json:"alarm_contacts" yaml:"alarm-contacts" mapstructure:"alarm-contacts"`
	AlarmSmsTemplateId string   `json:"alarm_sms_template_id" yaml:"alarm-sms-template-id" mapstructure:"alarm-sms-template-id"`
}
