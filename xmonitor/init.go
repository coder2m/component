/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/25 16:29
 **/
package xmonitor

import (
	xapp "github.com/myxy99/component"
	cfg "github.com/myxy99/component/xcfg"
	"github.com/myxy99/component/xgovern"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func init() {
	BuildInfoGauge.WithLabelValues(
		xapp.Name(),
		cfg.GetString("app.mode"),
		xapp.AppVersion(),
		xapp.GoVersion(),
		xapp.StartTime(),
	).Set(float64(time.Now().UnixNano() / 1e6))

	xgovern.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})
}
