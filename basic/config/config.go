package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var AppConfig *App

type App struct {
	Server   *Server   `yaml:"server"`
	Db       *Db       `yaml:"db"`
	Log      *Log      `yaml:"log"`
	Minio    *Minio    `yaml:"minio"`
	RabbitMQ *RabbitMQ `yaml:"rabbitmq"`
	Public   *Public   `yaml:"public"`
	Debug    *Debug    `yaml:"debug"`
}

type Server struct {
	Name    string `yaml:"name"`
	Port    int    `yaml:"port"`
	Version string `yaml:"version"`
}

type Db struct {
	Mysql    *Mysql    `yaml:"mysql"`
	Redis    *Redis    `yaml:"redis"`
	Tdengine *Tdengine `yaml:"tdengine"`
}

type Log struct {
	Path       string `yaml:"path"`
	Level      string `yaml:"level"`
	MaxSize    int    `yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
	MaxAge     int    `yaml:"maxAge"`
}

type Mysql struct {
	Url               string `yaml:"url"`
	MaxIdleConnection int    `yaml:"maxIdleConnection"`
	MaxOpenConnection int    `yaml:"maxOpenConnection"`
}

type Minio struct {
	AccessKey  string `yaml:"accessKey"`
	SecretKey  string `yaml:"secretKey"`
	BucketName string `yaml:"bucketName"`
	Endpoint   string `yaml:"endpoint"`
	Source     bool   `yaml:"source"`
	Region     string `yaml:"region"`
}

type Redis struct {
	DB       int    `yaml:"db"`
	Addr     string `yaml:"addr"`
	Network  string `yaml:"network"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Tdengine struct {
	Url string `yaml:"url"`
}

type Public struct {
	Ip string `yaml:"ip"`
}

type Debug struct {
	EnablePProf bool `yaml:"enablePProf"`
}

type RabbitMQ struct {
	Url string `yaml:"url"`
}

func GetConfig() *App {
	return AppConfig
}

// 读取配置文件
func InitConfig(path string) {

	configFile, err := os.ReadFile(path)
	if err != nil {
		panic("配置文件无效：" + err.Error())
	}

	_ = yaml.Unmarshal(configFile, &AppConfig)
}
