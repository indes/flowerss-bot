package log

import (
	"fmt"
	"github.com/indes/flowerss-bot/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	// Logger 日志对象
	Logger    *zap.Logger
	zapConfig zap.Config
)

func init() {
	logLevel := config.GetString("log.level")
	if logLevel == "Debug" || logLevel == "debug" {
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		zapConfig.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		zapConfig.EncoderConfig = zap.NewProductionEncoderConfig()
	}

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
}

// Println calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	Logger.WithOptions(zap.AddCallerSkip(1)).Info(fmt.Sprint(v...))
}

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	Logger.WithOptions(zap.AddCallerSkip(1)).Info(fmt.Sprintf(format, v...))
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	Logger.WithOptions(zap.AddCallerSkip(1)).Fatal(fmt.Sprint(v...))
}

// InfoWithMessage logs a message at InfoLevel with telegram message. The message includes telegram message info
// at the log site, as well as any fields accumulated on the logger.
func InfoWithMessage(m *tb.Message, v ...interface{}) {
	Logger.WithOptions(zap.AddCallerSkip(1)).Info(fmt.Sprint(v...),
		zap.Field{Key: "telegram message", Type: zapcore.ReflectType, Interface: m},
	)
}

// InfoWithField logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func InfoWithField(msg string, fields ...zapcore.Field) {
	Logger.WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(msg string, keysAndValues ...interface{}) {
	Logger.WithOptions(zap.AddCallerSkip(1)).Sugar().Infow(msg, keysAndValues...)
}

// Debug logs a message at DebugLevel.
func Debug(v ...interface{}) {
	Logger.WithOptions(zap.AddCallerSkip(1)).Debug(fmt.Sprint(v...))
}

// DebugWithField logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func DebugWithField(msg string, fields ...zapcore.Field) {
	Logger.WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Debugw(msg string, keysAndValues ...interface{}) {
	Logger.WithOptions(zap.AddCallerSkip(1)).Sugar().Debugw(msg, keysAndValues...)
}

// DebugWithMessage logs a message at DebugLevel. The message includes telegram message info
// at the log site, as well as any fields accumulated on the logger.
func DebugWithMessage(m *tb.Message, v ...interface{}) {
	Logger.WithOptions(zap.AddCallerSkip(1)).Debug(fmt.Sprint(v...),
		zap.Field{Key: "telegram message", Type: zapcore.ReflectType, Interface: m},
	)
}

// Warn logs a message at WarnLevel.
func Warn(v ...interface{}) {
	Logger.WithOptions(zap.AddCallerSkip(1)).Warn(fmt.Sprint(v...))
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keysAndValues ...interface{}) {
	Logger.WithOptions(zap.AddCallerSkip(1)).Sugar().Errorw(msg, keysAndValues...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(v ...interface{}) {
	Logger.WithOptions(zap.AddCallerSkip(1)).Error(fmt.Sprint(v...))
}