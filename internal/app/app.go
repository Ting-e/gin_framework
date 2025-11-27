package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project/internal/middleware"
	"project/pkg/config"
	"project/pkg/logger"
	"project/pkg/response"

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

// localhost页面参数
type PageData struct {
	ServiceName   string // 服务名称
	Now           string // 服务器时间
	Environment   string // 环境 (Development/Production)
	Version       string // 版本号
	StatusMessage string // 状态消息
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
	logger.Sugar.Infof("\t[app] using config file: %s", opts.ConfigPath)
	logger.Sugar.Infof("\t[app] using log directory: %s", opts.LogPath)
	logger.Sugar.Info("\t[app] logger initialized successfully")
	logger.Sugar.Infof("\t[app] %s initialized successfully, version: %s", app.name, app.version)

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

	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	d.router = gin.New()

	// 注册中间件
	middleware.RegisterDefaultMiddlewares(d.router, d.version)

	// 设置信任的代理
	if err := d.router.SetTrustedProxies(nil); err != nil {
		logger.Sugar.Errorf("[app] failed to set trusted proxies: %v", err)
	}

	// localhost
	d.router.GET("/", func(c *gin.Context) {

		response.Success(c, PageData{
			ServiceName:   config.Get().Server.Name,
			Now:           time.Now().Format("2006-01-02 15:04:05"),
			Environment:   config.Get().Server.Environment,
			Version:       config.Get().Server.Version,
			StatusMessage: config.Get().Server.Name + " is running",
		})
	})

	return d.router
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
		logger.Sugar.Infof("\t[app] server starting on %s", addr)
		logger.Sugar.Infof("\t[app] visit: http://localhost%s", addr)
		if err := d.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Sugar.Fatalf("[app] failed to start server: %v", err)
		}
	}()

	// 等待退出信号
	sig := <-quit
	logger.Sugar.Infof("\t[app] received signal: %v, shutting down...", sig)

	// 关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := d.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	// 加载组件
	if err := d.CloseComponents(); err != nil {
		return fmt.Errorf("failed to load components: %w", err)
	}

	logger.Sugar.Info("[app] server exited gracefully")
	return nil
}

// Shutdown 优雅关闭应用
func (d *DefaultApp) Shutdown(ctx context.Context) error {
	if d.server == nil {
		return nil
	}

	logger.Sugar.Info("\t[app] shutting down server...")

	// 关闭 HTTP 服务器
	if err := d.server.Shutdown(ctx); err != nil {
		logger.Sugar.Errorf("[app] server shutdown error: %v", err)
		return err
	}

	logger.Sugar.Info("\t[app] shutdown completed")
	return nil
}
