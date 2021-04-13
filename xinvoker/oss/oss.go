package xoss

import (
	"errors"
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-server/xinvoker/oss/alioss"
	"github.com/coder2z/g-server/xinvoker/oss/file"
	"github.com/coder2z/g-server/xinvoker/oss/standard"
)

func (i *ossInvoker) loadConfig() map[string]*options {
	conf := make(map[string]*options)

	prefix := i.key
	for name := range xcfg.GetStringMap(prefix) {
		cfg := xcfg.UnmarshalWithExpect(prefix+"."+name, newOssOptions()).(*options)
		conf[name] = cfg
	}
	return conf
}

func (i *ossInvoker) new(o *options) (client standard.Oss) {
	var err error
	switch o.Mode {
	case "aliOss":
		client, err = alioss.NewOss(o.Addr, o.AccessKeyID, o.AccessKeySecret, o.OssBucket, o.IsDeleteSrcPath)
	case "file":
		client, err = file.NewOss(o.CdnName, o.FileBucket, o.IsDeleteSrcPath)
	default:
		err = errors.New("oss mode not exist")
	}
	if err != nil {
		xlog.Panic("new oss", xlog.FieldErr(err))
	}
	return
}
