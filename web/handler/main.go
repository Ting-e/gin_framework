/*
 * @Author: 廷益 15137125570@163.com
 * @Date: 2025-05-09 15:18:13
 * @LastEditors: 廷益 15137125570@163.com
 * @LastEditTime: 2025-05-09 15:19:09
 * @FilePath: \backend_framework\web\handler\main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"flag"
	"project/basic"
	"project/basic/logger"
	"project/web/handler/controller"

	"github.com/robfig/cron"
)

// 全局变量定义
var (
	app  basic.App
	log  string
	conf string
)

// setupLogRotation 设置日志分割
// 使用cron包定时分割日志文件
func setupLogRotation() {
	// 日志定时分割
	logCron := cron.New()
	logCron.AddFunc("0 0 1 * *", func() {
		logger.InitLogger(log)
	})
	logCron.Start()
}

// init 初始化函数
func init() {
	initializeFlags()
	initializeApp()
	setupLogRotation()
}

// initializeFlags 初始化命令行参数
// 通过flag包获取命令行参数
//
//	-conf: 配置文件路径
//	-log: 日志文件路径
func initializeFlags() {
	flag.StringVar(&conf, "conf", "./conf/app-dev.yaml", "配置文件路径")
	flag.StringVar(&log, "log", "./log/", "日志文件路径")
	flag.Parse()
}

// initializeApp 初始化应用程序
// 通过命令行参数获取配置文件路径和日志文件路径
func initializeApp() {
	// 初始化app
	app = basic.InitApp(conf, log)

	// 设置需要加载的组件
	app.SetComponent([]string{"mysql"})
}

func main() {

	//获取app的路由
	r := app.GetRouter()

	// /端/功能
	login := r.Group("/applet/ceshi")
	login.POST("/obtain-data", controller.GetData)

	// 运行app
	if err := app.Run(); err != nil {
		logger.Sugar.Error(err)
	}
}
