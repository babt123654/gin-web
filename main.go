package main

import (
	"embed"
	"fmt"
	"gin-web/initialize"
	"gin-web/pkg/global"
	"gin-web/pkg/service"
	"gin-web/router"
	"github.com/piupuer/go-helper/pkg/listen"
	"github.com/piupuer/go-helper/pkg/log"
	"github.com/piupuer/go-helper/pkg/tracing"
	"github.com/pkg/errors"
	_ "net/http/pprof"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

var ctx = tracing.NewId(nil)

//go:embed conf
var conf embed.FS

// @title Gin Web
// @version 1.2.1
// @description A simple RBAC admin system written by golang
// @license.name MIT
// @license.url https://github.com/piupuer/gin-web/blob/dev/LICENSE
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	defer func() {
		if err := recover(); err != nil {
			log.WithContext(ctx).WithError(errors.Errorf("%v", err)).Error("[%s]project run failed, stack: %s", global.ProName, string(debug.Stack()))
		}
	}()

	// get runtime root
	_, file, _, _ := runtime.Caller(0)
	global.RuntimeRoot = strings.TrimSuffix(file, "main.go")

	// initialize components
	initialize.Config(ctx, conf)
	initialize.Tracer()
	//initialize.Redis()
	initialize.Mysql()
	initialize.CasbinEnforcer()

	//initialize.Cron()
	//initialize.Oss()
	//CreateL0(new(gin.Context))
	// listen http
	OminMoniter()
	fmt.Println("-------------------------------------------------------------------------------------------6")
	listen.Http(
		listen.WithHttpCtx(ctx),
		listen.WithHttpProName(global.ProName),
		listen.WithHttpPort(global.Conf.System.Port),
		listen.WithHttpPprofPort(global.Conf.System.PprofPort),
		listen.WithHttpHandler(router.RegisterServers(ctx)),
		listen.WithHttpExit(func() {
			global.Tracer.Shutdown(ctx)
		}),
	)

	//go func() {

	//}()
}
func OminMoniter() {
	q := service.New(ctx)
	//q.GetL0Sibel2Git()
	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				// 每隔5分钟调用一次q.GetOminiUsdtOp1(nil)函数
				q.GetOminiUsdtOp1(nil)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// 假设你的程序需要运行一段时间后才退出
	// 在这里可以添加退出条件，例如等待某个事件发生或接收到某个信号
	// 例如，使用time.Sleep(time.Hour)来模拟程序的运行时间
	time.Sleep(time.Hour * 9999)

	// 关闭定时器，并退出循环
	close(quit)
}
