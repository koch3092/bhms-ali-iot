package model

type AliSms struct {
	Phones        []string `json:"phones"`
	TemplateCode  string   `json:"templateCode"`
	TemplateParam string   `json:"templateParam"`
}
