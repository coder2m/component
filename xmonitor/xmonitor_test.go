/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/25 19:25
 **/
package xmonitor

import (
	"github.com/myxy99/component/xgovern"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func Monitor(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		h(w, r)

		ServerHandleHistogram.WithLabelValues(TypeHTTP, "http", "0").Observe(time.Since(start).Seconds())
	}

}

func Query(w http.ResponseWriter, r *http.Request) {

	//模拟业务查询耗时0~1s

	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	_, _ = io.WriteString(w, "some results")

}

func TestName(t *testing.T) {

	go xgovern.Run()

	http.HandleFunc("/query", Monitor(Query))

	t.Log(http.ListenAndServe(":8080", nil))
}
