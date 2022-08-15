package log

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	pattern  = "/log/level"
	endpoint = ":4246"
)

var logger *zap.SugaredLogger
var atomicLevel = zap.NewAtomicLevel()

var levelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"panic": zapcore.PanicLevel,
	"fatal": zapcore.FatalLevel,
}

func toZapLevel(l string) zapcore.Level {
	if level, ok := levelMap[l]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func Init(serviceName, filePath, level string) {
	http.HandleFunc(pattern, atomicLevel.ServeHTTP)
	go func() {
		if err := http.ListenAndServe(endpoint, nil); err != nil {
			panic(err)
		}
	}()

	// error, fatal, panic
	highLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zap.ErrorLevel
	})
	// info, debug
	lowLever := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l < zap.ErrorLevel && l >= toZapLevel(level)
	})

	highCore := newCore(filePath, highLevel, "error.log")
	lowCore := newCore(filePath, lowLever, "info.log")
	log := zap.New(
		zapcore.NewTee(highCore, lowCore),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.Development(),
		zap.Fields(zap.String("func", funcName())),
		zap.Fields(zap.String("usecase", serviceName)),
	)
	logger = log.Sugar()
}

func funcName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	return filepath.Base(runtime.FuncForPC(pc).Name())
}

func newCore(filePath string, priority zap.LevelEnablerFunc, filename string) zapcore.Core {
	filename = strings.Replace(filename, ".log", "", -1) +
		"-" + time.Now().Format("2006-01-02") + ".log"
	logPath := filepath.Join(filepath.Dir(filePath), filename)
	fileWriteSyncer := &lumberjack.Logger{
		Filename:   logPath, // 日志文件存放目录
		MaxSize:    100,     // 文件大小限制,单位MB
		MaxBackups: 30,      // 最大保留日志文件数量
		MaxAge:     7,       // 日志文件保留天数
		Compress:   true,    // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "log",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(fileWriteSyncer),
		),
		priority,
	)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	logger.Panicw(msg, keysAndValues...)
}

func Sync() error {
	if logger != nil {
		return logger.Sync()
	}
	return nil
}
