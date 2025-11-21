package logger

import (
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

func InitLogger(logDir string, logLevel string) {

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logDir + time.Now().Format("2006-01-02") + ".log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
	})

	// 判断配置文件配置的日志等级
	var zapLevel zapcore.Level

	switch logLevel {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	default:
		panic("level: debug, info, warn, error")
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout),
			w),
		zapLevel,
	)

	Logger = zap.New(core, zap.AddCaller())
	Sugar = Logger.Sugar()

	if logDir == "" {
		// 输出到 stdout
		logger, _ := zap.NewProduction()
		Sugar = logger.Sugar()
		return
	}

	setupLogRotation(logDir, logLevel)
}

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

var cronInstance *cron.Cron

func setupLogRotation(logDir string, logLevel string) {

	if cronInstance != nil {
		cronInstance.Stop() // 避免重复启动
	}

	cronInstance = cron.New()
	// 每天凌晨1点分割日志
	cronInstance.AddFunc("0 0 1 * * *", func() {
		InitLogger(logDir, logLevel)
	})
	cronInstance.Start()
}
