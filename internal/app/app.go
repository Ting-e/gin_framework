package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project/pkg/database"
	"project/pkg/queue"

	"project/pkg/config"
	"project/pkg/logger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// App 定义应用接口
type App interface {
	SetComponent(components []string)
	LoadComponents() error
	GetRouter() *gin.Engine
	InitPProf()
	Run() error
	Shutdown(ctx context.Context) error
}

// DefaultApp 默认应用实现
type DefaultApp struct {
	name       string
	port       int
	components []string
	router     *gin.Engine
	configPath string
	version    string
	server     *http.Server
}

// InitApp 初始化应用
func InitApp() (App, error) {

	// 加载配置路径和日志路径
	opts := LoadConfig()

	// 参数验证
	if opts.ConfigPath == "" {
		return nil, fmt.Errorf("config path is required")
	}
	if opts.LogPath == "" {
		return nil, fmt.Errorf("log path is required")
	}

	// 验证日志目录
	if _, err := os.Stat(opts.LogPath); err != nil {
		if os.IsNotExist(err) {
			// 尝试创建日志目录
			if err := os.MkdirAll(opts.LogPath, 0755); err != nil {
				return nil, fmt.Errorf("failed to create log directory: %w", err)
			}
		} else {
			return nil, fmt.Errorf("invalid log directory: %w", err)
		}
	}

	// 加载配置
	if err := config.Init(opts.ConfigPath); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	appConfig := config.Get()

	// 初始化日志
	logger.InitLogger(opts.LogPath, appConfig.Log.Level)

	app := &DefaultApp{
		name:       appConfig.Server.Name,
		port:       appConfig.Server.Port,
		version:    appConfig.Server.Version,
		configPath: opts.ConfigPath,
	}

	app.SetComponent(appConfig.Components)

	// 打印 Logo
	fmt.Println(logo)
	fmt.Println()

	// 打印启动信息
	logger.Sugar.Infof("[app] using config file: %s", opts.ConfigPath)
	logger.Sugar.Infof("[app] using log directory: %s", opts.LogPath)
	logger.Sugar.Info("[app] logger initialized successfully")
	logger.Sugar.Infof("[app] %s initialized successfully, version: %s", app.name, app.version)

	return app, nil
}

// MustInitApp 初始化应用，失败时 panic
func MustInitApp() App {
	app, err := InitApp()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize app: %v", err))
	}
	return app
}

// GetRouter 获取 Gin 路由实例
func (d *DefaultApp) GetRouter() *gin.Engine {
	if d.router != nil {
		return d.router
	}

	// // 根据配置设置 Gin 模式
	// if config.Get().Log.Level == "debug" {
	// 	gin.SetMode(gin.DebugMode)
	// } else {
	gin.SetMode(gin.ReleaseMode)
	// }

	d.router = gin.New()

	// 添加中间件
	d.router.Use(gin.Recovery())
	d.router.Use(d.loggerMiddleware())
	d.router.Use(d.corsMiddleware())

	// 设置信任的代理
	if err := d.router.SetTrustedProxies(nil); err != nil {
		logger.Sugar.Errorf("[app] failed to set trusted proxies: %v", err)
	}

	// 健康检查接口
	d.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": d.version,
			"name":    d.name,
		})
	})

	return d.router
}

// loggerMiddleware 日志中间件
func (d *DefaultApp) loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		if query != "" {
			path = path + "?" + query
		}

		logger.Sugar.Infof("[GIN] %s | %3d | %13v | %15s | %-7s %s",
			time.Now().Format("2006-01-02 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}

// corsMiddleware CORS 中间件
func (d *DefaultApp) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token, session, X-Requested-With, Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language, DNT, X-CustomHeader, Keep-Alive, User-Agent, If-Modified-Since, Cache-Control, Content-Type, X-App-Key, X-Jwt, X-Signature, X-Time-Stamp")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, Expires, Last-Modified, Pragma")
		c.Header("Access-Control-Max-Age", "172800")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("X-Service-Version", d.version)

		if method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// LoadComponents 加载组件
func (d *DefaultApp) LoadComponents() error {
	if len(d.components) == 0 {
		logger.Sugar.Info("[app] no components to load")
		return nil
	}

	logger.Sugar.Infof("[app] loading %d components...", len(d.components))

	for _, comp := range d.components {
		if err := d.loadComponent(comp); err != nil {
			return fmt.Errorf("failed to load component '%s': %w", comp, err)
		}
	}

	logger.Sugar.Info("[app] all components loaded successfully")
	return nil
}

// loadComponent 加载单个组件
func (d *DefaultApp) loadComponent(comp string) error {
	logger.Sugar.Infof("[app] loading component: %s", comp)

	switch comp {
	case "mysql":
		m := database.GetMysql()
		if m == nil {
			return fmt.Errorf("mysql component is nil")
		}
		m.InitComponent()

	case "gorm":
		gm := database.GetGormMysql()
		if gm == nil {
			return fmt.Errorf("gorm_mysql component is nil")
		}
		gm.InitComponent()

	case "redis":
		r := database.GetRedis()
		if r == nil {
			return fmt.Errorf("redis component is nil")
		}
		r.InitComponent()

	case "tdengine":
		t := database.GetTdengine()
		if t == nil {
			return fmt.Errorf("tdengine component is nil")
		}
		t.InitComponent()

	case "rabbitmq":
		rb := queue.GetRabbitMQ()
		if rb == nil {
			return fmt.Errorf("rabbitmq component is nil")
		}
		rb.InitComponent()

	default:
		return fmt.Errorf("unknown component: %s", comp)
	}

	logger.Sugar.Infof("[app] component '%s' loaded successfully", comp)
	return nil
}

// InitPProf 初始化性能分析工具
func (d *DefaultApp) InitPProf() {
	if !config.Get().Debug.EnablePProf {
		logger.Sugar.Info("[app] pprof is disabled")
		return
	}

	router := d.GetRouter()
	pprof.Register(router)
	logger.Sugar.Info("[app] pprof enabled, visit /debug/pprof/")
}

// SetComponent 设置组件列表
func (d *DefaultApp) SetComponent(components []string) {
	d.components = components
	logger.Sugar.Infof("[app] components set: %v", components)
}

// Run 启动应用（阻塞式，带优雅关闭）
func (d *DefaultApp) Run() error {
	if d.port <= 0 || d.port > 65535 {
		return fmt.Errorf("invalid port: %d (must be 1-65535)", d.port)
	}

	// 加载组件
	if err := d.LoadComponents(); err != nil {
		return fmt.Errorf("failed to load components: %w", err)
	}

	addr := fmt.Sprintf(":%d", d.port)

	// 创建 HTTP 服务器
	d.server = &http.Server{
		Addr:           addr,
		Handler:        d.router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 监听系统信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 在 goroutine 中启动服务器
	go func() {
		logger.Sugar.Infof("[app] server starting on %s", addr)
		logger.Sugar.Infof("[app] visit: http://localhost%s", addr)
		if err := d.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Sugar.Fatalf("[app] failed to start server: %v", err)
		}
	}()

	// 等待退出信号
	sig := <-quit
	logger.Sugar.Infof("[app] received signal: %v, shutting down...", sig)

	// 关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := d.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	logger.Sugar.Info("[app] server exited gracefully")
	return nil
}

// Shutdown 优雅关闭应用
func (d *DefaultApp) Shutdown(ctx context.Context) error {
	if d.server == nil {
		return nil
	}

	logger.Sugar.Info("[app] shutting down server...")

	// 关闭 HTTP 服务器
	if err := d.server.Shutdown(ctx); err != nil {
		logger.Sugar.Errorf("[app] server shutdown error: %v", err)
		return err
	}

	// 关闭数据库连接等资源
	// TODO: 添加组件清理逻辑
	for _, comp := range d.components {
		logger.Sugar.Infof("[app] cleaning up component: %s", comp)
		// 调用各组件的清理方法（如果有）
	}

	logger.Sugar.Info("[app] shutdown completed")
	return nil
}

// GetVersion 获取版本号
func (d *DefaultApp) GetVersion() string {
	return d.version
}

// GetName 获取应用名称
func (d *DefaultApp) GetName() string {
	return d.name
}

// GetPort 获取端口号
func (d *DefaultApp) GetPort() int {
	return d.port
}
