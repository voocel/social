package log

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	pattern  = "/log/level"
	endpoint = ":4246"
)

var levelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
	"panic": zapcore.PanicLevel,
}

var (
	srv    *http.Server
	logger *zap.SugaredLogger
)

func toZapLevel(l string) zapcore.Level {
	if level, ok := levelMap[l]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func Init(serviceName, level string, logPaths ...string) {
	var atomicLevel = zap.NewAtomicLevel()
	atomicLevel.SetLevel(toZapLevel(level))

	var logPath string
	if len(logPaths) == 0 {
		_, logPath, _, _ = runtime.Caller(0)
		logPath = filepath.Dir(filepath.Dir(filepath.Dir(logPath)))
		logPath = filepath.Join(logPath, "logs", serviceName)
	} else {
		logPath = logPaths[0]
	}

	mux := http.NewServeMux()
	mux.HandleFunc(pattern, atomicLevel.ServeHTTP)
	srv = &http.Server{
		Addr:    endpoint,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	//// error, fatal, panic
	//highLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
	//	return l >= zap.ErrorLevel
	//})
	//// info, debug
	//lowLever := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
	//	return l < zap.ErrorLevel && l >= toZapLevel(level)
	//})
	//
	//highCore := newCore(filePath, highLevel, "error.log")
	//lowCore := newCore(filePath, lowLever, "info.log")

	core := newCore(logPath, atomicLevel)
	log := zap.New(
		zapcore.NewTee(core),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.Development(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		//zap.Fields(zap.String("func", funcName())),
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

func newCore(logPath string, atomicLevel zap.AtomicLevel) zapcore.Core {
	filename := time.Now().Format("2006-01-02") + ".log"
	logPath = filepath.Join(logPath, filename)
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
		EncodeTime:     zapcore.RFC3339TimeEncoder,
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
		atomicLevel,
	)
}

type DefaultPair struct {
	key   string
	value interface{}
}

func Pair(key string, v interface{}) DefaultPair {
	return DefaultPair{
		key:   key,
		value: v,
	}
}

func spread(kvs ...DefaultPair) []interface{} {
	s := make([]interface{}, 0, len(kvs))
	for _, v := range kvs {
		s = append(s, v.key, v.value)
	}
	return s
}

// Debug .
func Debug(args ...interface{}) {
	logger.Debug(args...)
}
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}
func Debugw(msg string, kvs ...DefaultPair) {
	args := spread(kvs...)
	logger.Debugw(msg, args...)
}

// Info .
func Info(args ...interface{}) {
	logger.Info(args...)
}
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}
func Infow(msg string, kvs ...DefaultPair) {
	args := spread(kvs...)
	logger.Infow(msg, args...)
}

// Warn .
func Warn(args ...interface{}) {
	logger.Warn(args...)
}
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}
func Warnw(msg string, kvs ...DefaultPair) {
	args := spread(kvs...)
	logger.Warnw(msg, args...)
}

// Error .
func Error(args ...interface{}) {
	logger.Error(args...)
}
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}
func Errorw(msg string, kvs ...DefaultPair) {
	args := spread(kvs...)
	logger.Errorw(msg, args...)
}

// Fatal .
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
func Fatalw(msg string, kvs ...DefaultPair) {
	args := spread(kvs...)
	logger.Fatalw(msg, args...)
}

// Panic .
func Panic(args ...interface{}) {
	logger.Panic(args...)
}
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}
func Panicw(msg string, kvs ...DefaultPair) {
	args := spread(kvs...)
	logger.Panicw(msg, args...)
}

func Sync() error {
	if srv != nil && logger != nil {
		srv.Shutdown(context.Background())
		return logger.Sync()
	}
	return nil
}
