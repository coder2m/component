package xgovern

import (
	"context"
	"errors"
	"fmt"
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xdefer"
	"github.com/coder2z/g-saber/xjson"
	"github.com/coder2z/g-saber/xlog"
	"github.com/coder2z/g-saber/xtime"
	"github.com/coder2z/g-server/xapp"
	"github.com/coder2z/g-server/xcode"
	"github.com/coder2z/g-server/xmonitor"
	"net/http"
	"net/http/pprof"
	"os"
	"sync"
	"time"
)

type healthStats struct {
	IP         string `json:"ip,omitempty"`
	Hostname   string `json:"hostname,omitempty"`
	AppName    string `json:"app_name,omitempty"`
	AppId      string `json:"app_id,omitempty"`
	AppMode    string `json:"app_mode,omitempty"`
	AppDebug   bool   `json:"app_debug"`
	AppVersion string `json:"app_version,omitempty"`
	GoVersion  string `json:"go_version,omitempty"`
	Time       string `json:"time,omitempty"`
	Err        string `json:"err,omitempty"`
	Status     string `json:"status,omitempty"`
}

type h map[string]func(w http.ResponseWriter, r *http.Request)

var (
	handle       *http.ServeMux
	server       *http.Server
	HandleFuncs  = make(h)
	governConfig *Config
	once         = sync.Once{}
)

func (hm h) Run(hs *http.ServeMux) {
	for s, f := range hm {
		hs.HandleFunc(s, f)
	}
}

func HandleFunc(p string, h func(w http.ResponseWriter, r *http.Request)) {
	HandleFuncs[p] = h
}

func init() {
	xcfg.OnChange(func(*xcfg.Configuration) {
		GovernReload()
	})
}

func init() {
	HandleFunc("/debug/pprof", pprof.Index)
	HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	HandleFunc("/debug/pprof/profile", pprof.Profile)
	HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	HandleFunc("/debug/pprof/trace", pprof.Trace)

	HandleFunc("/debug/code/business", xcode.XCodeBusinessCodeHttp)
	HandleFunc("/debug/code/system", xcode.XCodeSystemCodeHttp)

	HandleFunc("/metrics", xmonitor.MonitorPrometheusHttp)

	HandleFunc("/debug/env", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_ = xjson.NewEncoder(w).Encode(os.Environ())
	})

	HandleFunc("/debug/list", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		list := make([]string, 0)
		for s, _ := range HandleFuncs {
			list = append(list, s)
		}
		_ = xjson.NewEncoder(w).Encode(list)
	})

	HandleFunc("/debug/config", func(w http.ResponseWriter, r *http.Request) {
		mm := xcfg.Traverse(".")
		w.WriteHeader(200)
		_ = xjson.NewEncoder(w).Encode(mm)
	})

	HandleFunc("/debug/health", func(w http.ResponseWriter, r *http.Request) {
		serverStats := healthStats{
			IP:         xapp.HostIP(),
			Hostname:   xapp.HostName(),
			AppName:    xapp.Name(),
			AppId:      xapp.AppId(),
			AppMode:    xapp.AppMode(),
			AppDebug:   xapp.Debug(),
			AppVersion: xapp.AppVersion(),
			GoVersion:  xapp.GoVersion(),
			Time:       xtime.Now().Format("2006-01-02 15:04:05"),
			Err:        "",
			Status:     "SUCCESS",
		}
		w.WriteHeader(200)
		_ = xjson.NewEncoder(w).Encode(serverStats)
	})
}

func GetServer() *http.ServeMux {
	if handle == nil {
		handle = http.NewServeMux()
	}
	return handle
}

func Run(opts ...Option) {
	once.Do(func() {
		c := GovernConfig()

		for _, opt := range opts {
			opt(c)
		}

		HandleFuncs.Run(GetServer())

		server = &http.Server{
			Addr:    c.Address(),
			Handler: handle,
		}

		xlog.Info("Application Starting",
			xlog.FieldComponentName("XGovern"),
			xlog.FieldMethod("XGovern.Run"),
			xlog.FieldDescription(fmt.Sprintf("Govern serve running :%v/debug/list", c.Address())),
		)

		xlog.Info("Application Starting",
			xlog.FieldComponentName("Prometheus"),
			xlog.FieldMethod("XGovern.Run"),
			xlog.FieldDescription(fmt.Sprintf("Prometheus serve running :%v/metrics", c.Address())),
		)

		xdefer.Register(func() error {
			return Shutdown()
		})

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			xlog.Error("Govern Serve ListenAndServe Error",
				xlog.FieldComponentName("XGovern"),
				xlog.FieldMethod("XGovern.Run.ListenAndServe"),
				xlog.FieldErr(err),
				xlog.FieldAddr(c.Address()))
		}
	})
}

func Shutdown() error {
	if server == nil {
		return errors.New("shutdown govern server")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		xlog.Error("Shutdown Govern Serve Error",
			xlog.FieldComponentName("XGovern"),
			xlog.FieldMethod("XGovern.Shutdown"),
			xlog.FieldErr(err))
		return err
	}
	xlog.Info("Application Stopping",
		xlog.FieldComponentName("XGovern"),
		xlog.FieldMethod("XGovern.Shutdown"),
		xlog.FieldDescription("XGovern server shutdown"),
	)
	return nil
}

func GovernConfig() *Config {
	if governConfig == nil {
		governConfig = xcfg.UnmarshalWithExpect("app.govern", DefaultConfig()).(*Config)
	}
	return governConfig
}

func GovernReload(opts ...Option) {
	xlog.Info("Application Reload",
		xlog.FieldComponentName("XGovern"),
		xlog.FieldMethod("XGovern.GovernReload"),
		xlog.FieldDescription("XGovern server reload"),
	)
	_ = Shutdown()
	once = sync.Once{}
	governConfig = nil
	server = nil
	handle = nil
	Run(opts...)
}
