package main

import (
	"flag"
	"fmt"
	"ops_agent/auth"
	"ops_agent/conf"
	"ops_agent/core"
	"time"
)

var (
	FlagServer   string // 指定服务端url（格式：http://xxx）
	FlagUsername string // 指定jwt认证用户名
	FlagPassword string // 指定jwt认证密码
	FlagTimeout  int    // 指定超时时间（jwt认证和主机id获取的超时时间）
	FlagLogLevel string // 指定日志级别
	FlagLogPath  string // 指定日志存放路径
	FlagLogName  string // 指定日志名
)

func init() {
	flag.StringVar(&FlagServer, "server", "", "spec server url，eg：http://xxx")
	flag.StringVar(&FlagUsername, "username", "", "jwt auth username")
	flag.StringVar(&FlagPassword, "password", "", "jwt auth password")
	flag.StringVar(&FlagLogLevel, "log-level", "Debug", "spec log level，eg: Debug Info Warn Error Fatal，default: Debug")
	flag.StringVar(&FlagLogPath, "log-path", "/var/log", "spec log file path，eg: /var/log/asset，default:/var/log")
	flag.StringVar(&FlagLogName, "log-name", "asset.log", "spec log filename，eg: xx.log，default: asset.log")
	flag.IntVar(&FlagTimeout, "timeout", 10, "jwt auth and get asset_id request timeout，default: 10s")
}

func run() {
	//// 初始化channel(管道)
	//core.InitAssetChan()
	//core.InitMonitorChan()
	for {
		select {
		case asset_res := <-core.AssetChan:
			fmt.Println(asset_res)
			//pkg.Request_Asset(conf.ReportUrl,auth.Jwt_Token,asset_res)

		case monitor_res := <-core.MonitorChan:
			fmt.Println(monitor_res)

		default:
			time.Sleep(time.Second * 1)
			conf.Logger.Info("等待数据.....")

		}
	}
}

func main() {
	// 0、解析命令参数
	flag.Parse()

	// 1、初始化配置
	conf.Init_FlagProcess(FlagServer)
	conf.Init_Log(FlagLogLevel, FlagLogPath, FlagLogName)

	// 2、agent初始化操作
	// 发送auth认证
	auth.Jwt_Auth(conf.JwtUrl, FlagUsername, FlagPassword, FlagTimeout)
	// 获取当前系统基本数据，并发送该信息从server端获取id
	core.Asset_Host_Base("/var/asset_id", conf.IdUrl, auth.Jwt_Token, FlagTimeout)

	// 3、开始执行系统资产采集
	go core.RunAsset()
	go core.RunMonitor()

	// 4、异步执行程序
	run()

}
