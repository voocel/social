package tcp

import (
	"fmt"
	"log"
	"os"
)

// Won't compile if Logger can't be realized by a DefaultLogger
var _ Logger = (*DefaultLogger)(nil)

var Flog Logger = newLogger()

type Level uint8

const (
	DEBUG Level = iota
	INFO
	ERROR
	FATAL
	PANIC
)

type Logger interface {
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
}

type DefaultLogger struct {
	rawLogger *log.Logger
}

func newLogger() *DefaultLogger {
	return &DefaultLogger{
		rawLogger: log.New(os.Stderr, "【SOCIAL】", log.LstdFlags),
	}
}

// Info logs a message at level Info on the standard log
func (d *DefaultLogger) Info(v ...interface{}) {
	d.rawLogger.Printf("[INFO] %s", v...)
}

// Infof logs a message at level Info on the standard log
func (d *DefaultLogger) Infof(format string, v ...interface{}) {
	d.rawLogger.Printf("[INFO] %s", fmt.Sprintf(format, v...))
}

// Debug logs a message at level Debug on the standard log
func (d *DefaultLogger) Debug(v ...interface{}) {
	d.rawLogger.Printf("[DEBUG] %s", v...)
}

// Debugf logs a message at level Debug on the standard log
func (d *DefaultLogger) Debugf(format string, v ...interface{}) {
	d.rawLogger.Printf("[DEBUG] %s", fmt.Sprintf(format, v...))
}

// Error logs a message at level Error on the standard log
func (d *DefaultLogger) Error(v ...interface{}) {
	d.rawLogger.Printf("[ERROR] %s", v...)
}

// Errorf logs a message at level Error on the standard log
func (d *DefaultLogger) Errorf(format string, v ...interface{}) {
	d.rawLogger.Printf("[ERROR] %s", fmt.Sprintf(format, v...))
}

// Fatal logs a message at level Fatal on the standard log then the process will exit with status set to 1
func (d *DefaultLogger) Fatal(v ...interface{}) {
	d.rawLogger.Fatalf("[FATAL] %v", v)
}

// Fatalf logs a message at level Fatal on the standard log then the process will exit with status set to 1
func (d *DefaultLogger) Fatalf(format string, v ...interface{}) {
	d.rawLogger.Fatalf("[FATAL] %s", fmt.Sprintf(format, v...))
}

// Panic logs a message at level Panic on the standard log
func (d *DefaultLogger) Panic(v ...interface{}) {
	d.rawLogger.Panicf("[PANIC] %v", v)
}

// Panicf logs a message at level Panic on the standard log
func (d *DefaultLogger) Panicf(format string, v ...interface{}) {
	d.rawLogger.Panicf("[PANIC] %s", fmt.Sprintf(format, v...))
}

// SetLogger custom log
func SetLogger(l Logger) {
	if l == nil {
		panic("log is nil")
	}
	Flog = l
}
