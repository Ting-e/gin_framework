package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

// setupTestConfig 创建临时配置文件。
func setupTestConfig(t *testing.T, content string) string {
	t.Helper()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	return configPath
}

func TestInit_ValidConfig(t *testing.T) {
	ResetForTesting()

	configContent := `
server:
  name: test-app
  port: 9090
  version: v2.0.0

log:
  path: ./test-logs/app.log
  level: debug
  maxSize: 50
  maxBackups: 3
  maxAge: 14

db:
  mysql:
    url: "user:pass@tcp(localhost:3306)/dbname"
    maxIdleConnection: 5
    maxOpenConnection: 50
  redis:
    addr: "localhost:6379"
    db: 1
    password: "secret"

storage:
  accessKey: "admin"
  secretKey: "admin123"
  bucketName: "test-bucket"
  endpoint: "localhost:9000"
  region: "us-west-1"
  source: true
`

	configPath := setupTestConfig(t, configContent)
	err := Init(configPath)
	require.NoError(t, err)

	cfg := Get()
	assert.NotNil(t, cfg)
	assert.Equal(t, "test-app", cfg.Server.Name)
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "v2.0.0", cfg.Server.Version)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, 50, cfg.Log.MaxSize)
	assert.Equal(t, "user:pass@tcp(localhost:3306)/dbname", cfg.Db.Mysql.URL)
	assert.Equal(t, 1, cfg.Db.Redis.DB)
	assert.Equal(t, "secret", cfg.Db.Redis.Password)
	assert.Equal(t, "admin123", cfg.Storage.SecretKey)
	assert.True(t, cfg.Storage.Source)
}

func TestInit_DefaultValues(t *testing.T) {
	ResetForTesting()

	configContent := `
server:
  name: test-app
  port: 8080

db:
  mysql:
    url: "user:pass@tcp(localhost:3306)/dbname"
`

	configPath := setupTestConfig(t, configContent)
	err := Init(configPath)
	require.NoError(t, err)

	cfg := Get()

	// 验证默认值
	assert.Equal(t, "v1.0.0", cfg.Server.Version)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, 100, cfg.Log.MaxSize)
	assert.Equal(t, 5, cfg.Log.MaxBackups)
	assert.Equal(t, 7, cfg.Log.MaxAge)
	assert.Equal(t, "./logs/app.log", cfg.Log.Path)
	assert.Equal(t, 10, cfg.Db.Mysql.MaxIdleConnection)
	assert.Equal(t, 100, cfg.Db.Mysql.MaxOpenConnection)
}

func TestInit_InvalidPort(t *testing.T) {
	ResetForTesting()

	configContent := `
server:
  name: test-app
  port: 99999
`

	configPath := setupTestConfig(t, configContent)
	err := Init(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid port")
}

func TestInit_MissingRequiredFields(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		expectedErr string
	}{
		{
			name: "missing mysql url",
			config: `
server:
  name: test-app
  port: 8080
db:
  mysql:
    maxIdleConnection: 10
`,
			expectedErr: "mysql.url is required",
		},
		{
			name: "missing storage accessKey",
			config: `
server:
  name: test-app
  port: 8080
storage:
  secretKey: "secret"
  bucketName: "bucket"
  endpoint: "localhost:9000"
`,
			expectedErr: "accessKey is required",
		},
		{
			name: "invalid log level",
			config: `
server:
  name: test-app
  port: 8080
log:
  level: invalid
`,
			expectedErr: "invalid log level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResetForTesting()
			configPath := setupTestConfig(t, tt.config)
			err := Init(configPath)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestInit_FileNotFound(t *testing.T) {
	ResetForTesting()

	err := Init("/nonexistent/config.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config file")
}

func TestInit_InvalidYAML(t *testing.T) {
	ResetForTesting()

	configContent := `
server:
  name: test-app
  port: [invalid yaml
`

	configPath := setupTestConfig(t, configContent)
	err := Init(configPath)
	assert.Error(t, err)
}

func TestGet_NotInitialized(t *testing.T) {
	ResetForTesting()

	assert.Panics(t, func() {
		Get()
	})
}

func TestIsInitialized(t *testing.T) {
	ResetForTesting()

	assert.False(t, IsInitialized())

	configContent := `
server:
  name: test-app
  port: 8080
db:
  mysql:
    url: "user:pass@tcp(localhost:3306)/dbname"
`

	configPath := setupTestConfig(t, configContent)
	err := Init(configPath)
	require.NoError(t, err)

	assert.True(t, IsInitialized())
}

func TestSafeCopy(t *testing.T) {
	ResetForTesting()

	configContent := `
server:
  name: test-app
  port: 8080

db:
  mysql:
    url: "user:password123@tcp(localhost:3306)/dbname"
  redis:
    addr: "localhost:6379"
    password: "redis-secret"

storage:
  accessKey: "AKIAIOSFODNN7EXAMPLE"
  secretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  bucketName: "test-bucket"
  endpoint: "localhost:9000"

rabbitmq:
  url: "amqp://user:pass@localhost:5672/"
`

	configPath := setupTestConfig(t, configContent)
	err := Init(configPath)
	require.NoError(t, err)

	cfg := Get()
	safeCfg := cfg.SafeCopy()

	// 验证脱敏
	assert.Equal(t, "***REDACTED***", safeCfg.Storage.SecretKey)
	assert.Contains(t, safeCfg.Storage.AccessKey, "***")
	assert.NotEqual(t, cfg.Storage.AccessKey, safeCfg.Storage.AccessKey)
	assert.Equal(t, "***REDACTED***", safeCfg.Db.Redis.Password)
	assert.Contains(t, safeCfg.Db.Mysql.URL, "***REDACTED***")
	assert.Contains(t, safeCfg.RabbitMQ.URL, "***REDACTED***")

	// 验证原始配置未被修改
	assert.Equal(t, "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", cfg.Storage.SecretKey)
	assert.Equal(t, "redis-secret", cfg.Db.Redis.Password)
}

func TestSafeCopy_Nil(t *testing.T) {
	var cfg *App
	safeCfg := cfg.SafeCopy()
	assert.Nil(t, safeCfg)
}

func TestRedactDSN(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "mysql dsn",
			input:    "mysql://user:password@localhost:3306/dbname",
			expected: "mysql://user:***REDACTED***@localhost:3306/dbname",
		},
		{
			name:     "postgres dsn",
			input:    "postgresql://admin:secret123@db.example.com:5432/mydb",
			expected: "postgresql://admin:***REDACTED***@db.example.com:5432/mydb",
		},
		{
			name:     "rabbitmq url",
			input:    "amqp://guest:guest@localhost:5672/",
			expected: "amqp://guest:***REDACTED***@localhost:5672/",
		},
		{
			name:     "no password",
			input:    "redis://localhost:6379",
			expected: "redis://localhost:6379",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactDSN(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRedactKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "long key",
			input:    "AKIAIOSFODNN7EXAMPLE",
			expected: "AKIA***MPLE",
		},
		{
			name:     "short key",
			input:    "short",
			expected: "***REDACTED***",
		},
		{
			name:     "exactly 8 chars",
			input:    "12345678",
			expected: "***REDACTED***",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactKey(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidate_ServerConfig(t *testing.T) {
	tests := []struct {
		name        string
		server      *Server
		expectError bool
	}{
		{
			name:        "valid config",
			server:      &Server{Name: "app", Port: 8080, Version: "v1.0.0"},
			expectError: false,
		},
		{
			name:        "invalid port - too low",
			server:      &Server{Name: "app", Port: 0, Version: "v1.0.0"},
			expectError: true,
		},
		{
			name:        "invalid port - too high",
			server:      &Server{Name: "app", Port: 99999, Version: "v1.0.0"},
			expectError: true,
		},
		{
			name:        "missing name",
			server:      &Server{Name: "", Port: 8080, Version: "v1.0.0"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.server.Validate()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidate_LogConfig(t *testing.T) {
	tests := []struct {
		name        string
		log         *Log
		expectError bool
	}{
		{
			name:        "valid config",
			log:         &Log{Level: "info", MaxSize: 100, MaxBackups: 5, MaxAge: 7},
			expectError: false,
		},
		{
			name:        "invalid level",
			log:         &Log{Level: "invalid", MaxSize: 100, MaxBackups: 5, MaxAge: 7},
			expectError: true,
		},
		{
			name:        "invalid maxSize",
			log:         &Log{Level: "info", MaxSize: 0, MaxBackups: 5, MaxAge: 7},
			expectError: true,
		},
		{
			name:        "negative maxBackups",
			log:         &Log{Level: "info", MaxSize: 100, MaxBackups: -1, MaxAge: 7},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.log.Validate()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEnvOverride(t *testing.T) {
	ResetForTesting()

	configContent := `
server:
  name: test-app
  port: 8080
`

	// 设置环境变量
	os.Setenv("APP_SERVER_PORT", "9090")
	os.Setenv("APP_SERVER_NAME", "env-app")
	defer func() {
		os.Unsetenv("APP_SERVER_PORT")
		os.Unsetenv("APP_SERVER_NAME")
	}()

	configPath := setupTestConfig(t, configContent)
	err := Init(configPath)
	require.NoError(t, err)

	cfg := Get()
	// 注意：环境变量覆盖在优化版本中使用 APP_ 前缀
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, "env-app", cfg.Server.Name)
}

func TestResetForTesting(t *testing.T) {
	configContent := `
server:
  name: test-app
  port: 8080
db:
  mysql:
    url: "user:pass@tcp(localhost:3306)/dbname"
`

	configPath := setupTestConfig(t, configContent)
	err := Init(configPath)
	require.NoError(t, err)
	assert.True(t, IsInitialized())

	ResetForTesting()
	assert.False(t, IsInitialized())

	// 再次初始化应该成功
	err = Init(configPath)
	require.NoError(t, err)
	assert.True(t, IsInitialized())
}
