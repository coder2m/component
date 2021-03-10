/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/25 18:20
 **/
package xapp

import (
	"github.com/coder2z/g-saber/xcfg"
	"github.com/coder2z/g-saber/xconsole"
	"github.com/coder2z/g-saber/xnet"
	"github.com/coder2z/g-saber/xstring"
	"github.com/coder2z/g-server/xversion"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	dAppName    = "MyApp"
	dAppVersion = "v0.1.0"
)

func init() {
	xconsole.Blue(`   _____ ____  _____  ______ _____  ___  ______`)
	xconsole.Blue(`  / ____/ __ \|  __ \|  ____|  __ \|__ \|___  /`)
	xconsole.Blue(` | |   | |  | | |  | | |__  | |__) |  ) |  / /`)
	xconsole.Blue(` | |   | |  | | |  | |  __| |  _  /  / /  / / `)
	xconsole.Blue(` | |___| |__| | |__| | |____| | \ \ / /_ / /__ `)
	xconsole.Blue(`  \_____\____/|_____/|______|_|  \_|____/_____|`)
	xconsole.Blue(`									--version = ` + xversion.Version)
	startTime = time.Now().Format("2006-01-02 15:04:05")
	goVersion = runtime.Version()
}

var (
	startTime       string
	goVersion       string
	appName         string
	hostName        string
	buildAppVersion string
	buildHost       string
	debug           = true
	one             = sync.Once{}
	appId           = xstring.GenerateID()
)

// Name gets application name.
func Name() string {
	if appName == "" {
		if appName = xcfg.GetString("app.name"); appName == "" {
			appName = dAppName
		}
	}
	return appName
}

// Debug gets application debug.
func Debug() bool {
	one.Do(func() {
		if data := xcfg.GetString("app.debug"); data == "false" {
			debug = false
		}
	})
	return debug
}

//AppVersion get buildAppVersion
func AppVersion() string {
	if buildAppVersion == "" {
		if buildAppVersion = xcfg.GetString("app.version"); buildAppVersion == "" {
			buildAppVersion = dAppVersion
		}
	}
	return buildAppVersion
}

//BuildHost get buildHost
func BuildHost() string {
	if buildHost == "" {
		var err error
		if buildHost, err = xnet.GetLocalIP(); err != nil {
			hostName = "0.0.0.0"
		}
	}
	return buildHost
}

// HostName get host name
func HostName() string {
	if hostName == "" {
		var err error
		if hostName, err = os.Hostname(); err != nil {
			hostName = "unknown"
		}
	}
	return hostName
}

//StartTime get start time
func StartTime() string {
	return startTime
}

//GoVersion get go version
func GoVersion() string {
	return goVersion
}

func AppId() string {
	return appId
}

func PrintVersion() {
	xconsole.Greenf("app name:", Name())
	xconsole.Greenf("app id:", AppId())
	xconsole.Greenf("host name:", HostName())
	xconsole.Greenf("app debug:", Debug())
	xconsole.Greenf("app version:", AppVersion())
	xconsole.Greenf("build host:", BuildHost())
	xconsole.Greenf("start time:", StartTime())
	xconsole.Greenf("go version:", GoVersion())
}
