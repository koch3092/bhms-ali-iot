package service

import (
	"bhms-ali-iot/global"
	"errors"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	"strings"
)

type AliSmsService struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SendAliSms
//@description: 发送（阿里）短信
//@return: err error

// 模板使用json字符串 {"code":"xxx"} 对应你模板里面的变量key和value
func (e *AliSmsService) SendAliSms(phones []string, templateCode string, templateParam string) (err error) {
	client := global.AliSms
	if len(phones) == 0 {
		return errors.New("请输入手机号")
	}

	phonesStr := strings.Join(phones, ",")
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String(global.CONFIG.AliSms.SignName),
		PhoneNumbers:  tea.String(phonesStr),
		TemplateCode:  tea.String(templateCode),
		TemplateParam: tea.String(templateParam),
	}
	// 复制代码运行请自行打印 API 的返回值
	_, err = client.SendSms(sendSmsRequest)
	if err != nil {
		return err
	}
	return err
}
