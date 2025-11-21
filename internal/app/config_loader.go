package app

import (
	"flag"
	"os"
)

type ConfigOptions struct {
	ConfigPath string
	LogPath    string
}

// LoadConfig 加载配置路径和日志路径
func LoadConfig() ConfigOptions {
	var opts ConfigOptions

	// 默认值
	defaultConf := "./configs/app-dev.yaml"
	defaultLog := "./log/"

	// 从环境变量读取（Docker/K8s）
	if envConf := os.Getenv("APP_CONFIG"); envConf != "" {
		defaultConf = envConf
	}
	if envLog := os.Getenv("APP_LOG_DIR"); envLog != "" {
		defaultLog = envLog
	}

	// 命令行参数可以覆盖
	flag.StringVar(&opts.ConfigPath, "conf", defaultConf, "配置文件路径")
	flag.StringVar(&opts.LogPath, "log", defaultLog, "日志目录路径")
	flag.Parse()

	return opts
}
