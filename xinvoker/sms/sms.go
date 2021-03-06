/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/6 18:05
 **/
package xsms

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/coder2z/component/xcfg"
)

type (
	SmsResponse = dysmsapi.SendSmsResponse
	SmsRequest  = dysmsapi.SendSmsRequest
	Client      struct {
		SMS          *dysmsapi.Client
		signName     string
		templateCode string
	}
)

func (i *smsInvoker) newSMSClient(o *options) *Client {
	c, err := dysmsapi.NewClientWithAccessKey(o.Area, o.AccessKeyId, o.AccessSecret)
	if err != nil {
		panic(err)
	}
	return &Client{SMS: c, signName: o.SignName, templateCode: o.TemplateCode}
}

func (i *smsInvoker) loadConfig() map[string]*options {
	conf := make(map[string]*options)
	prefix := i.key
	for name := range xcfg.GetStringMap(prefix) {
		cfg := xcfg.UnmarshalWithExpect(prefix+"."+name, newSMSOptions()).(*options)
		conf[name] = cfg
	}
	return conf
}

func (ali *Client) Send(req *SmsRequest) (*SmsResponse, error) {
	if req.RpcRequest == nil {
		req.RpcRequest = new(requests.RpcRequest)
	}
	if req.TemplateCode == "" {
		req.TemplateCode = ali.templateCode
	}
	if req.SignName == "" {
		req.SignName = ali.signName
	}
	req.InitWithApiInfo("Dysmsapi", "2017-05-25", "SendSms", "dysms", "openAPI")
	rep, err := ali.SMS.SendSms(req)
	if err != nil {
		return nil, err
	}
	return rep, nil
}
