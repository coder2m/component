package xgorm

import (
	"fmt"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-server/xinvoker"
	"gorm.io/gorm"
	"sync"
)

var db *dbInvoker

func Register(k string) xinvoker.Invoker {
	db = &dbInvoker{key: k}
	return db
}

func Invoker(key string) *gorm.DB {
	if val, ok := db.instances.Load(key); ok {
		return val.(*gorm.DB)
	}
	xlog.Panic("Application Starting",
		xlog.FieldComponentName("XInvoker"),
		xlog.FieldMethod("XInvoker.XGorm"),
		xlog.FieldDescription(fmt.Sprintf("no db(%s) invoker found", key)),
	)

	return nil
}

type dbInvoker struct {
	xinvoker.Base
	instances sync.Map
	key       string
}

func (i *dbInvoker) Init(opts ...xinvoker.Option) error {
	i.instances = sync.Map{}
	for name, cfg := range i.loadConfig() {
		db := i.newDatabaseClient(cfg)
		i.instances.Store(name, db)
	}
	return nil
}

func (i *dbInvoker) Reload(opts ...xinvoker.Option) error {
	for name, cfg := range i.loadConfig() {
		db := i.newDatabaseClient(cfg)
		i.instances.Store(name, db)
	}
	return nil
}

func (i *dbInvoker) Close(opts ...xinvoker.Option) error {
	i.instances.Range(func(key, value interface{}) bool {
		db, _ := value.(*gorm.DB).DB()
		_ = db.Close()
		i.instances.Delete(key)
		return true
	})
	return nil
}
