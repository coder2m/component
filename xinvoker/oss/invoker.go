package xoss

import (
	"fmt"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-server/xinvoker"
	"github.com/coder2z/g-server/xinvoker/oss/standard"
	"sync"
)

var ossI *ossInvoker

func Register(k string) xinvoker.Invoker {
	ossI = &ossInvoker{key: k}
	return ossI
}

func Invoker(key string) standard.Oss {
	if val, ok := ossI.instances.Load(key); ok {
		return val.(standard.Oss)
	}
	xlog.Panic("Application Starting",
		xlog.FieldComponentName("XInvoker"),
		xlog.FieldMethod("XInvoker.XOss"),
		xlog.FieldDescription(fmt.Sprintf("no oss(%s) invoker found", key)),
	)
	return nil
}

type ossInvoker struct {
	xinvoker.Base
	instances sync.Map
	key       string
}

func (i *ossInvoker) Init(opts ...xinvoker.Option) error {
	i.instances = sync.Map{}
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.new(cfg))
	}
	return nil
}

func (i *ossInvoker) Reload(opts ...xinvoker.Option) error {
	for name, cfg := range i.loadConfig() {
		i.instances.Store(name, i.new(cfg))
	}
	return nil
}

func (i *ossInvoker) Close(opts ...xinvoker.Option) error {
	return nil
}
