package basic

import (
	"fmt"
	"net/http"
	"os"
	"project/basic/component/minio"
	mysqldb "project/basic/component/mysql"
	rabbitmq "project/basic/component/rabbitmq"
	redisdb "project/basic/component/redis"
	tdenginedb "project/basic/component/tdengine"
	"project/basic/config"
	"project/basic/logger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

type App interface {
	SetComponent(components []string)
	LoadComponents()
	GetRouter() *gin.Engine
	InitPProf()
	Run() error
}

type DefaultApp struct {
	name       string
	port       int
	components []string
	router     *gin.Engine
	configPath string
	version    string
}

// InitApp 初始化应用
func InitApp(configPath string, logPath string) App {
	if configPath == "" {
		panic("未指定配置文件")
	}
	if logPath == "" {
		panic("未指定日志目录位置")
	}

	// 先加载配置
	err := config.Init(configPath)
	if err != nil {
		logger.Sugar.Error(err)
	}

	AppConfig := config.Get()

	// 验证日志目录
	if _, err := os.ReadDir(logPath); err != nil {
		panic("日志目录位置无效: " + err.Error())
	}

	// 初始化日志
	logger.InitLogger(logPath)

	app := &DefaultApp{
		name:       AppConfig.Server.Name,
		port:       AppConfig.Server.Port,
		version:    AppConfig.Server.Version,
		configPath: configPath,
	}

	// LOGO输出
	_, _ = fmt.Println(logo)
	_, _ = fmt.Println()

	// 打印启动过程
	logger.Sugar.Infof("\t[app] using config file: %s", configPath)
	logger.Sugar.Infof("\t[app] using log directory: %s", logPath)
	logger.Sugar.Info("\t[app] logger init success")
	logger.Sugar.Infof("\t[app] %s init success. version：%s", app.name, app.version)

	return app
}

func (d *DefaultApp) GetRouter() *gin.Engine {
	if d.router != nil {
		return d.router // 避免重复初始化
	}

	gin.SetMode(gin.ReleaseMode)
	d.router = gin.Default()

	// 仅在首次创建时设置中间件
	d.dealCros()

	if err := d.router.SetTrustedProxies(nil); err != nil {
		logger.Sugar.Panicf("[app] failed to set trusted proxies: %v", err)
	}

	return d.router
}

func (d *DefaultApp) dealCros() {
	d.router.Use(func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token, session, X_Requested_With, Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language, DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, X-App-Key, X-Jwt, X-Signature, X-Time-Stamp")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, Expires, Last-Modified, Pragma")
		c.Header("Access-Control-Max-Age", "172800")
		c.Header("Access-Control-Allow-Credentials", "false")
		c.Header("service_version", d.version)

		if method == http.MethodOptions {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "CORS preflight success"})
			return
		}

		c.Next()
	})
}

func (d *DefaultApp) LoadComponents() {
	for _, comp := range d.components {
		switch comp {
		case "mysql":
			if m := mysqldb.GetMysql(); m != nil {
				m.InitComponent()
			}
		case "minio":
			if m := minio.GetMinio(); m != nil {
				m.InitComponent()
			}
		case "redis":
			if r := redisdb.GetRedis(); r != nil {
				r.InitComponent()
			}
		case "tdengine":
			if t := tdenginedb.GetTdengine(); t != nil {
				t.InitComponent()
			}
		case "rabbitmq":
			if t := rabbitmq.GetRabbitMQ(); t != nil {
				t.InitComponent()
			}
		}
	}
}

// 初始化性能分析工具
func (d *DefaultApp) InitPProf() {
	if !config.Get().Debug.EnablePProf {
		return
	}

	router := d.GetRouter()
	pprof.Register(router)
}

func (d *DefaultApp) SetComponent(components []string) {
	d.components = components
}

func (d *DefaultApp) Run() error {
	if d.port <= 0 || d.port > 65535 {
		return fmt.Errorf("invalid port: %d", d.port)
	}

	addr := fmt.Sprintf(":%d", d.port)
	logger.Sugar.Infof("[app] server starting on：%s", addr)

	d.LoadComponents()

	return d.router.Run(addr)
}
