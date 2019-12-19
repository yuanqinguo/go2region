package main

import (
	"flag"
	"fmt"
	_ "github.com/icattlecoder/godaemon" //包的init函数实现，启动时加入 -d=true表示为daemon运行
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	. "go2region/config"
	"go2region/controls"
	"go2region/utils/logs"
	"go2region/web"
	"os"
	"time"
)

func main() {
	// 初始化配置文件
	flag.Parse()
	fmt.Print("InitConfig...\r")
	checkErr("InitConfig", InitConfig())
	fmt.Print("InitConfig Success!!!\n")

	baseConf := Config.BaseConf

	initLogger(baseConf)

	startService(baseConf)
}

func initLogger(baseConf BaseConf) {
	// 创建文件日志，按天分割，默认日志文件仅保留一周, 此处更新为根据配置文件决定,不支持动态更改
	// access log
	accessLog, err := rotatelogs.New(baseConf.AccessLogPath, rotatelogs.WithMaxAge(time.Duration(baseConf.LogMaxAge)*24*time.Hour))
	checkErr("CreateRotateLog", err)
	// error log
	errorLog, err := rotatelogs.New(baseConf.SystemLogPath, rotatelogs.WithMaxAge(time.Duration(baseConf.LogMaxAge)*24*time.Hour))
	checkErr("CreateRotateLog", err)

	// 设置日志, 分为access日志
	logs.LogAccess.SetOutput(accessLog)
	logs.LogAccess.SetFormatter(&logrus.JSONFormatter{})
	logs.LogAccess.SetReportCaller(false)

	//系统运行日志
	logs.LogSystem.SetOutput(errorLog)
	logs.LogSystem.SetFormatter(&logrus.JSONFormatter{})
	logs.LogSystem.SetReportCaller(false)
	logLevel, _ := logrus.ParseLevel(DynamicConf.LogLevel)
	logs.LogSystem.SetLevel(logLevel)
}

func startService(conf BaseConf) {
	// 获取IpDataInfo实例，同时加载数据
	ipdb := controls.GetInstance()
	defer ipdb.Close()

	ipdb.Reloader(conf.IpdataPath)

	// 开始运行iris框架
	fmt.Print("RunIris...\r")
	web.RunIris(conf.ServerPort)
}

// 检查错误
func checkErr(errMsg string, err error) {
	if err != nil {
		fmt.Printf("%s Error: %v\n", errMsg, err)
		os.Exit(1)
	}
}
