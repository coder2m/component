/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 11:18
 */
package xsms

type options struct {
	Area         string `mapStructure:"area"`
	AccessKeyId  string `mapStructure:"accessKeyId"`
	AccessSecret string `mapStructure:"accessSecret"`
	SignName     string `mapStructure:"signName"`
	TemplateCode string `mapStructure:"templateCode"`
}

func newSMSOptions() *options {
	return &options{
		Area:         "ap-guangzhou",
		AccessKeyId:  "",
		AccessSecret: "",
		SignName:     "",
		TemplateCode: "",
	}
}
