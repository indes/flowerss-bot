package log

import (
	"strings"

	"github.com/indes/flowerss-bot/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Logger 日志对象
	Logger    *zap.Logger
	zapConfig zap.Config
)

func init() {
	logLevel := config.GetString("log.level")
	if strings.ToLower(logLevel) == "debug" {
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		zapConfig.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		zapConfig.EncoderConfig = zap.NewProductionEncoderConfig()
	}

	//日志时间戳人类可读
	zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	logFile := config.GetString("log.file")
	if logFile != "" {
		zapConfig.Sampling = &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		}
		zapConfig.Encoding = "json"
		zapConfig.OutputPaths = []string{logFile}
		zapConfig.ErrorOutputPaths = []string{logFile}

	} else {
		zapConfig.OutputPaths = []string{"stderr"}
		zapConfig.ErrorOutputPaths = []string{"stderr"}
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.Encoding = "console"
	}

	Logger, _ = zapConfig.Build()
	zap.ReplaceGlobals(Logger)
}
