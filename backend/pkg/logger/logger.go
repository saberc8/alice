package logger

import (
	"log"
	"os"
	"strings"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var (
	debugLog *log.Logger
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
	level    Level
)

// Init 初始化日志
func Init(logLevel string) {
	debugLog = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile)
	infoLog = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
	warnLog = log.New(os.Stdout, "[WARN] ", log.LstdFlags)
	errorLog = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)

	switch strings.ToLower(logLevel) {
	case "debug":
		level = LevelDebug
	case "info":
		level = LevelInfo
	case "warn":
		level = LevelWarn
	case "error":
		level = LevelError
	default:
		level = LevelInfo
	}
}

// Debug 输出调试日志
func Debug(v ...interface{}) {
	if level <= LevelDebug {
		debugLog.Println(v...)
	}
}

// Debugf 输出格式化调试日志
func Debugf(format string, v ...interface{}) {
	if level <= LevelDebug {
		debugLog.Printf(format, v...)
	}
}

// Info 输出信息日志
func Info(v ...interface{}) {
	if level <= LevelInfo {
		infoLog.Println(v...)
	}
}

// Infof 输出格式化信息日志
func Infof(format string, v ...interface{}) {
	if level <= LevelInfo {
		infoLog.Printf(format, v...)
	}
}

// Warn 输出警告日志
func Warn(v ...interface{}) {
	if level <= LevelWarn {
		warnLog.Println(v...)
	}
}

// Warnf 输出格式化警告日志
func Warnf(format string, v ...interface{}) {
	if level <= LevelWarn {
		warnLog.Printf(format, v...)
	}
}

// Error 输出错误日志
func Error(v ...interface{}) {
	if level <= LevelError {
		errorLog.Println(v...)
	}
}

// Errorf 输出格式化错误日志
func Errorf(format string, v ...interface{}) {
	if level <= LevelError {
		errorLog.Printf(format, v...)
	}
}
