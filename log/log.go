// log/log.go 负责日志系统的初始化、输出、级别控制等功能
package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var logger *zap.SugaredLogger

const (
	DEBUG  = zapcore.DebugLevel
	INFO   = zapcore.InfoLevel
	WARN   = zapcore.WarnLevel
	ERROR  = zapcore.ErrorLevel
	DPANIC = zapcore.DPanicLevel
	PANIC  = zapcore.PanicLevel
	FATAL  = zapcore.FatalLevel
)

func init() {
	Init(GetConsoleZapcore(INFO))
}

// 获取控制台日志核心
// 设置日志等级
func GetConsoleZapcore(level zapcore.Level) zapcore.Core {
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	encoder.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.String() + "]")
	}
	encoder.EncodeCaller = nil
	encoder.ConsoleSeparator = " "
	return zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), os.Stdout, zap.NewAtomicLevelAt(level))
}

// 获取文件日志核心
// 设置绝对路径
// 设置日志等级
// 设置日志文件最大尺寸（单位：M）
func GetFileZapcore(absPath string, level zapcore.Level, fileMaxSize int) zapcore.Core {
	syncWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:  absPath,
		MaxSize:   (1 << 20) * fileMaxSize, //20M
		LocalTime: true,
		Compress:  true,
	})
	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	return zapcore.NewCore(zapcore.NewJSONEncoder(encoder), syncWriter, zap.NewAtomicLevelAt(level))
}

// Init 初始化日志系统
// 接受多个日志核心参数，支持同时输出到多个目标
func Init(cores ...zapcore.Core) {
	core := zapcore.NewTee(cores...)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	logger = zapLogger.Sugar()
}

// Debug 输出调试级别日志
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugf 输出格式化的调试级别日志
func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

// Info 输出一般信息级别日志
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof 输出格式化的一般信息级别日志
func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

// Warn 输出警告级别日志
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Warnf 输出格式化的警告级别日志
func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

// Error 输出错误级别日志
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorf 输出格式化的错误级别日志
func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

// DPanic 输出严重错误级别日志，并导致程序崩溃
func DPanic(args ...interface{}) {
	logger.DPanic(args...)
}

// DPanicf 输出格式化的严重错误级别日志，并导致程序崩溃
func DPanicf(template string, args ...interface{}) {
	logger.DPanicf(template, args...)
}

// Panic 输出恐慌级别日志，并导致程序崩溃
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Panicf 输出格式化的恐慌级别日志，并导致程序崩溃
func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

// Fatal 输出致命错误级别日志，并导致程序退出
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Fatalf 输出格式化的致命错误级别日志，并导致程序退出
func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}
