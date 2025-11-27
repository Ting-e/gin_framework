package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name     string
		logDir   string
		logLevel string
		wantErr  bool
	}{
		{
			name:     "Info level with directory",
			logDir:   "testlogs/",
			logLevel: "info",
			wantErr:  false,
		},
		{
			name:     "Debug level",
			logDir:   "testlogs/",
			logLevel: "debug",
			wantErr:  false,
		},
		{
			name:     "Warn level",
			logDir:   "testlogs/",
			logLevel: "warn",
			wantErr:  false,
		},
		{
			name:     "Error level",
			logDir:   "testlogs/",
			logLevel: "error",
			wantErr:  false,
		},
		{
			name:     "Empty logDir should use stdout",
			logDir:   "",
			logLevel: "info",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理测试目录
			if tt.logDir != "" {
				os.MkdirAll(tt.logDir, 0755)
				defer os.RemoveAll(tt.logDir)
			}

			// 停止之前的 cron 实例
			if cronInstance != nil {
				cronInstance.Stop()
				cronInstance = nil
			}

			InitLogger(tt.logDir, tt.logLevel)

			if Logger == nil {
				t.Error("Logger should not be nil")
			}
			if Sugar == nil {
				t.Error("Sugar should not be nil")
			}

			// 测试日志写入
			Logger.Info("test message")
			Sugar.Info("test sugar message")

			// 验证日志文件是否创建（仅当 logDir 非空时）
			if tt.logDir != "" {
				logFile := filepath.Join(tt.logDir, time.Now().Format("2006-01-02")+".log")
				if _, err := os.Stat(logFile); os.IsNotExist(err) {
					t.Errorf("Log file %s should exist", logFile)
				}
			}
		})
	}
}

func TestInitLoggerPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid log level")
		}
		// 清理
		if cronInstance != nil {
			cronInstance.Stop()
			cronInstance = nil
		}
	}()

	InitLogger("testlogs/", "invalid")
}

func TestNewEncoderConfig(t *testing.T) {
	config := NewEncoderConfig()

	if config.TimeKey != "T" {
		t.Errorf("Expected TimeKey to be 'T', got %s", config.TimeKey)
	}
	if config.LevelKey != "L" {
		t.Errorf("Expected LevelKey to be 'L', got %s", config.LevelKey)
	}
	if config.NameKey != "N" {
		t.Errorf("Expected NameKey to be 'N', got %s", config.NameKey)
	}
	if config.CallerKey != "C" {
		t.Errorf("Expected CallerKey to be 'C', got %s", config.CallerKey)
	}
	if config.MessageKey != "M" {
		t.Errorf("Expected MessageKey to be 'M', got %s", config.MessageKey)
	}
	if config.StacktraceKey != "S" {
		t.Errorf("Expected StacktraceKey to be 'S', got %s", config.StacktraceKey)
	}
	if config.LineEnding != zapcore.DefaultLineEnding {
		t.Error("Expected default line ending")
	}
}

func TestTimeEncoder(t *testing.T) {
	testTime := time.Date(2024, 1, 15, 14, 30, 45, 123456789, time.UTC)

	encoder := zapcore.NewConsoleEncoder(NewEncoderConfig())
	buf, err := encoder.EncodeEntry(zapcore.Entry{
		Time: testTime,
	}, nil)

	if err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("Encoded output should not be empty")
	}

	// 验证时间格式是否包含在输出中
	// 注意：完整验证需要解析整个日志行
	t.Logf("Encoded output: %s", output)
}

func TestLogLevels(t *testing.T) {
	testDir := "testlogs/"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	tests := []struct {
		name     string
		logLevel string
		logFunc  func(string)
		message  string
	}{
		{
			name:     "Debug level",
			logLevel: "debug",
			logFunc:  func(msg string) { Logger.Debug(msg) },
			message:  "debug message",
		},
		{
			name:     "Info level",
			logLevel: "info",
			logFunc:  func(msg string) { Logger.Info(msg) },
			message:  "info message",
		},
		{
			name:     "Warn level",
			logLevel: "warn",
			logFunc:  func(msg string) { Logger.Warn(msg) },
			message:  "warn message",
		},
		{
			name:     "Error level",
			logLevel: "error",
			logFunc:  func(msg string) { Logger.Error(msg) },
			message:  "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if cronInstance != nil {
				cronInstance.Stop()
				cronInstance = nil
			}

			InitLogger(testDir, tt.logLevel)
			tt.logFunc(tt.message)

			// 验证日志文件是否存在
			logFile := filepath.Join(testDir, time.Now().Format("2006-01-02")+".log")
			if _, err := os.Stat(logFile); os.IsNotExist(err) {
				t.Errorf("Log file should exist: %s", logFile)
			}
		})
	}
}

func TestSetupLogRotation(t *testing.T) {
	testDir := "testlogs/"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// 停止之前的 cron
	if cronInstance != nil {
		cronInstance.Stop()
		cronInstance = nil
	}

	setupLogRotation(testDir, "info")

	if cronInstance == nil {
		t.Error("cronInstance should not be nil after setupLogRotation")
	}

	// 验证 cron 任务数量
	entries := cronInstance.Entries()
	if len(entries) != 1 {
		t.Errorf("Expected 1 cron entry, got %d", len(entries))
	}

	// 清理
	cronInstance.Stop()
	cronInstance = nil
}

func TestMultipleInitLoggerCalls(t *testing.T) {
	testDir := "testlogs/"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// 多次调用 InitLogger
	for i := 0; i < 3; i++ {
		InitLogger(testDir, "info")
		Logger.Info("test message", zap.Int("iteration", i))
	}

	// 验证只有一个 cron 实例在运行
	if cronInstance == nil {
		t.Error("cronInstance should exist")
	}

	entries := cronInstance.Entries()
	if len(entries) != 1 {
		t.Errorf("Expected 1 cron entry after multiple inits, got %d", len(entries))
	}

	// 清理
	if cronInstance != nil {
		cronInstance.Stop()
		cronInstance = nil
	}
}

func TestSugarLogger(t *testing.T) {
	testDir := "testlogs/"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	if cronInstance != nil {
		cronInstance.Stop()
		cronInstance = nil
	}

	InitLogger(testDir, "info")

	if Sugar == nil {
		t.Fatal("Sugar logger should not be nil")
	}

	// 测试 Sugar logger 的各种方法
	Sugar.Info("info message")
	Sugar.Infof("formatted %s", "message")
	Sugar.Infow("structured message", "key", "value")

	logFile := filepath.Join(testDir, time.Now().Format("2006-01-02")+".log")
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Log file should exist: %s", logFile)
	}
}

// 基准测试
func BenchmarkInitLogger(b *testing.B) {
	testDir := "benchlogs/"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if cronInstance != nil {
			cronInstance.Stop()
			cronInstance = nil
		}
		InitLogger(testDir, "info")
	}
}

func BenchmarkLoggerInfo(b *testing.B) {
	testDir := "benchlogs/"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	if cronInstance != nil {
		cronInstance.Stop()
		cronInstance = nil
	}
	InitLogger(testDir, "info")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Logger.Info("benchmark message", zap.Int("iteration", i))
	}
}

func BenchmarkSugarInfo(b *testing.B) {
	testDir := "benchlogs/"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	if cronInstance != nil {
		cronInstance.Stop()
		cronInstance = nil
	}
	InitLogger(testDir, "info")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sugar.Infow("benchmark message", "iteration", i)
	}
}
