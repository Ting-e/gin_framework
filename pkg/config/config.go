// Package config 提供应用配置管理。
package config

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	appConfig *App
	once      sync.Once
	initErr   error
)

// App 是顶层配置结构。
type App struct {
	Server     *Server   `mapstructure:"server"`
	Db         *Db       `mapstructure:"db"`
	Log        *Log      `mapstructure:"log"`
	JWT        *JWT      `mapstructure:"jwt"`
	Storage    *Storage  `mapstructure:"storage"`
	RabbitMQ   *RabbitMQ `mapstructure:"rabbitmq"`
	Public     *Public   `mapstructure:"public"`
	Debug      *Debug    `mapstructure:"debug"`
	Components []string  `mapstructure:"components"`
}

// Server 配置。
type Server struct {
	Name        string `mapstructure:"name"`
	Port        int    `mapstructure:"port"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

// JWT JWT 配置
type JWT struct {
	Secret             string `mapstructure:"secret"`
	Issuer             string `mapstructure:"issuer"`
	ExpiresHours       int    `mapstructure:"expires_hours"`
	RefreshExpiresDays int    `mapstructure:"refresh_expires_days"`
}

// Db 数据库配置。
type Db struct {
	Mysql    *Mysql    `mapstructure:"mysql"`
	Redis    *Redis    `mapstructure:"redis"`
	Tdengine *Tdengine `mapstructure:"tdengine"`
}

// Log 日志配置。
type Log struct {
	Path       string `mapstructure:"path"`
	Level      string `mapstructure:"level"`
	MaxSize    int    `mapstructure:"maxSize"`    // MB
	MaxBackups int    `mapstructure:"maxBackups"` // 文件数
	MaxAge     int    `mapstructure:"maxAge"`     // 天
}

// Mysql 配置。
type Mysql struct {
	URL               string `mapstructure:"url"`
	MaxIdleConnection int    `mapstructure:"maxIdleConnection"`
	MaxOpenConnection int    `mapstructure:"maxOpenConnection"`
}

// Storage 配置。
type Storage struct {
	AccessKey  string `mapstructure:"accessKey"`
	SecretKey  string `mapstructure:"secretKey"`
	BucketName string `mapstructure:"bucketName"`
	Endpoint   string `mapstructure:"endpoint"`
	Source     bool   `mapstructure:"source"`
	Region     string `mapstructure:"region"`
}

// Redis 配置。
type Redis struct {
	DB       int    `mapstructure:"db"`
	Addr     string `mapstructure:"addr"`
	Network  string `mapstructure:"network"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// Tdengine 配置。
type Tdengine struct {
	URL string `mapstructure:"url"`
}

// Public 公共配置。
type Public struct {
	IP string `mapstructure:"ip"`
}

// Debug 调试配置。
type Debug struct {
	EnablePProf bool `mapstructure:"enablePProf"`
}

// RabbitMQ 配置。
type RabbitMQ struct {
	URL string `mapstructure:"url"`
}

// Get 返回当前加载的配置（只读）。
// 必须先调用 Init。
func Get() *App {
	if appConfig == nil {
		panic("config not initialized; call config.Init() first")
	}
	return appConfig
}

// MustGet 返回配置，如果未初始化则 panic。
func MustGet() *App {
	return Get()
}

// IsInitialized 检查配置是否已初始化。
func IsInitialized() bool {
	return appConfig != nil
}

// SafeCopy 返回脱敏后的配置副本，用于日志或调试。
func (a *App) SafeCopy() *App {
	if a == nil {
		return nil
	}

	cp := *a

	// 深拷贝嵌套结构
	if a.Server != nil {
		server := *a.Server
		cp.Server = &server
	}

	// 深拷贝嵌套结构
	if a.JWT != nil {
		jwt := *a.JWT
		cp.JWT = &jwt
	}

	if a.Db != nil {
		db := *a.Db
		cp.Db = &db

		if a.Db.Redis != nil {
			redis := *a.Db.Redis
			redis.Password = "***REDACTED***"
			cp.Db.Redis = &redis
		}

		if a.Db.Mysql != nil {
			mysql := *a.Db.Mysql
			mysql.URL = redactDSN(mysql.URL)
			cp.Db.Mysql = &mysql
		}

		if a.Db.Tdengine != nil {
			td := *a.Db.Tdengine
			td.URL = redactDSN(td.URL)
			cp.Db.Tdengine = &td
		}
	}

	if a.Storage != nil {
		storage := *a.Storage
		storage.SecretKey = "***REDACTED***"
		storage.AccessKey = redactKey(storage.AccessKey)
		cp.Storage = &storage
	}

	if a.RabbitMQ != nil {
		rmq := *a.RabbitMQ
		rmq.URL = redactDSN(rmq.URL)
		cp.RabbitMQ = &rmq
	}

	if a.Log != nil {
		log := *a.Log
		cp.Log = &log
	}

	if a.Public != nil {
		public := *a.Public
		cp.Public = &public
	}

	if a.Debug != nil {
		debug := *a.Debug
		cp.Debug = &debug
	}

	return &cp
}

// redactDSN 脱敏 DSN 字符串中的密码部分。
func redactDSN(dsn string) string {
	if dsn == "" {
		return ""
	}

	// 简单处理：查找 ":" 和 "@" 之间的内容
	if idx1 := strings.Index(dsn, "://"); idx1 != -1 {
		rest := dsn[idx1+3:]
		if idx2 := strings.Index(rest, ":"); idx2 != -1 {
			if idx3 := strings.Index(rest[idx2:], "@"); idx3 != -1 {
				return dsn[:idx1+3+idx2+1] + "***REDACTED***" + rest[idx2+idx3:]
			}
		}
	}

	return dsn
}

// redactKey 部分脱敏 key（显示前后几位）。
func redactKey(key string) string {
	if len(key) <= 8 {
		return "***REDACTED***"
	}
	return key[:4] + "***" + key[len(key)-4:]
}

// Validate 验证必要配置项是否有效。
func (a *App) Validate() error {
	if a.Server == nil {
		return errors.New("missing 'server' section")
	}

	if a.JWT != nil {
		if err := a.JWT.Validate(); err != nil {
			return fmt.Errorf("jwt config: %w", err)
		}
	}

	if err := a.Server.Validate(); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	if a.Db != nil {
		if err := a.Db.Validate(); err != nil {
			return fmt.Errorf("db config: %w", err)
		}
	}

	if a.Storage != nil {
		if err := a.Storage.Validate(); err != nil {
			return fmt.Errorf("storage config: %w", err)
		}
	}

	if a.Log != nil {
		if err := a.Log.Validate(); err != nil {
			return fmt.Errorf("log config: %w", err)
		}
	}

	if a.Components != nil {
		if len(a.Components) <= 0 {
			return fmt.Errorf("components config: can`t null")
		}
	}

	return nil
}

// Validate 验证 Server 配置。
func (s *Server) Validate() error {
	if s.Port <= 0 || s.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be 1-65535)", s.Port)
	}
	if s.Name == "" {
		return errors.New("server name is required")
	}
	return nil
}

// Validate 验证 Storage 配置。
func (j *JWT) Validate() error {
	if j.Secret == "" {
		return errors.New("secret is required")
	}

	if j.Issuer == "" {
		return errors.New("issuer is required")
	}

	if j.ExpiresHours == 0 {
		return errors.New("expires_hours is required")
	}

	if j.RefreshExpiresDays == 0 {
		return errors.New("refresh_expires_days is required")
	}
	return nil
}

// Validate 验证 Db 配置。
func (d *Db) Validate() error {
	if d.Mysql != nil {
		if d.Mysql.URL == "" {
			return errors.New("mysql.url is required")
		}
		if d.Mysql.MaxIdleConnection < 0 {
			return errors.New("mysql.maxIdleConnection must be >= 0")
		}
		if d.Mysql.MaxOpenConnection <= 0 {
			return errors.New("mysql.maxOpenConnection must be > 0")
		}
	}

	if d.Redis != nil && d.Redis.Addr == "" {
		return errors.New("redis.addr is required")
	}

	if d.Tdengine != nil && d.Tdengine.URL == "" {
		return errors.New("tdengine.url is required")
	}

	return nil
}

// Validate 验证 Storage 配置。
func (s *Storage) Validate() error {
	if s.AccessKey == "" {
		return errors.New("accessKey is required")
	}
	if s.SecretKey == "" {
		return errors.New("secretKey is required")
	}
	if s.BucketName == "" {
		return errors.New("bucketName is required")
	}
	if s.Endpoint == "" {
		return errors.New("endpoint is required")
	}
	return nil
}

// Validate 验证 Log 配置。
func (l *Log) Validate() error {
	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true, "fatal": true,
	}
	if !validLevels[strings.ToLower(l.Level)] {
		return fmt.Errorf("invalid log level: %s (must be debug/info/warn/error/fatal)", l.Level)
	}
	if l.MaxSize <= 0 {
		return errors.New("maxSize must be > 0")
	}
	if l.MaxBackups < 0 {
		return errors.New("maxBackups must be >= 0")
	}
	if l.MaxAge < 0 {
		return errors.New("maxAge must be >= 0")
	}
	return nil
}

// Init 从指定路径加载配置文件，并支持环境变量覆盖。
// 环境变量命名规则：APP_ 前缀 + 大写 + 下划线，例如 APP_SERVER_PORT。
func Init(configPath string) error {
	once.Do(func() {
		v := viper.New()
		v.SetConfigFile(configPath)
		v.SetConfigType("yaml")

		// 设置环境变量前缀和替换规则
		v.SetEnvPrefix("APP")
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		v.AutomaticEnv()

		// 设置默认值
		setViperDefaults(v)

		// 读取配置文件
		if err := v.ReadInConfig(); err != nil {
			initErr = fmt.Errorf("failed to read config file %q: %w", configPath, err)
			return
		}

		var cfg App
		if err := v.Unmarshal(&cfg); err != nil {
			initErr = fmt.Errorf("failed to unmarshal config: %w", err)
			return
		}

		// 验证配置
		if err := cfg.Validate(); err != nil {
			initErr = fmt.Errorf("invalid config: %w", err)
			return
		}

		appConfig = &cfg
	})

	return initErr
}

// setViperDefaults 设置 viper 默认值。
func setViperDefaults(v *viper.Viper) {
	// Server 默认值
	v.SetDefault("server.name", "my-app")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.version", "v1.0.0")

	// JWT 默认值
	v.SetDefault("jwt.secret", "woaifcll")
	v.SetDefault("jwt.issuer", "tinge")
	v.SetDefault("jwt.expires_hours", 2)
	v.SetDefault("jwt.refresh_expires_days", 7)

	// Log 默认值
	v.SetDefault("log.path", "./log/")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.maxSize", 100)
	v.SetDefault("log.maxBackups", 5)
	v.SetDefault("log.maxAge", 7)

	// Redis 默认值
	v.SetDefault("db.redis.db", 0)
	v.SetDefault("db.redis.addr", "localhost:6379")
	v.SetDefault("db.redis.network", "tcp")

	// Mysql 默认值
	v.SetDefault("db.mysql.maxIdleConnection", 10)
	v.SetDefault("db.mysql.maxOpenConnection", 100)

	// Storage 默认值
	v.SetDefault("storage.source", false)
	v.SetDefault("storage.region", "us-east-1")

	// Debug 默认值
	v.SetDefault("debug.enablePProf", false)
}

// ResetForTesting 重置配置状态，仅供测试使用。
// 生产环境请勿调用此函数。
func ResetForTesting() {
	appConfig = nil
	initErr = nil
	once = sync.Once{}
}
