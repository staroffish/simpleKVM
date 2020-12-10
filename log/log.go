package log

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"time"
)

const (
	LOG_LEVEL_DEBUG = iota
	LOG_LEVEL_INFO
)

var loglevelStringMap = map[int]string{
	LOG_LEVEL_DEBUG: "DEBUG",
	LOG_LEVEL_INFO:  "INFO",
}

type logger interface {
	PrintDebug(format string, args ...interface{})
	PrintInfo(format string, args ...interface{})
}

type baseLogger struct {
	logLevel int
	output   io.Writer
}

var log logger

func InitLog(logLevel int, output io.Writer) {
	base := &baseLogger{
		logLevel: logLevel,
		output:   output,
	}
	log = base
}

func PrintDebug(format string, args ...interface{}) {
	log.PrintDebug(format, args...)
}
func PrintInfo(format string, args ...interface{}) {
	log.PrintInfo(format, args...)
}

func (b *baseLogger) PrintDebug(format string, args ...interface{}) {
	b.printf(LOG_LEVEL_DEBUG, format, args...)
}
func (b *baseLogger) PrintInfo(format string, args ...interface{}) {
	b.printf(LOG_LEVEL_INFO, format, args...)
}

func (b *baseLogger) printf(logLevel int, format string, args ...interface{}) {
	if logLevel < b.logLevel {
		return
	}
	head := logHead(logLevel)
	fmt.Fprintf(b.output, "%s ", head)
	fmt.Fprintf(b.output, format, args...)
	fmt.Fprintf(b.output, "\n")
}

func logHead(logLevel int) string {
	_, file, line, _ := runtime.Caller(4)
	fileName := filepath.Base(file)
	// example: 2006/01/02 15:04:05 [INFO] main.go:60 XXXXXXXXXX
	return fmt.Sprintf("%s [%s] %s:%d", time.Now().Format("2006/01/02 15:04:05"), loglevelStringMap[logLevel], fileName, line)
}
