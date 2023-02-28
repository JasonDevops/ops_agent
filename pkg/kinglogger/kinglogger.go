package kinglogger

import (
	"strings"
)

// 定义日志文件级别
type Level uint16

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)


// 定义一个logger接口
type Loggers interface {
	Debug(format string,args ...interface{})
	Info(format string,args ...interface{})
	Warn(format string,args ...interface{})
	Error(format string,args ...interface{})
	Fatal(format string,args ...interface{})

}

// 根据用户传入的日志级别返回字符串

func getLevelStr(level Level)string{
	switch level {
	case DebugLevel:
		return "Debug"
	case InfoLevel:
		return "Info"
	case WarnLevel:
		return "Warn"
	case ErrorLevel:
		return "Error"
	case FatalLevel:
		return "Fatal"
	default:
		return "Debug"
	}
}

// 根据用户传入的字符串类型的日志级别，解析出对应的Level
func parseLogLevel(levelStr string)Level{
	levelStr = strings.ToLower(levelStr) // 将字符串转成全小写
	switch levelStr {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return DebugLevel
	}
}

