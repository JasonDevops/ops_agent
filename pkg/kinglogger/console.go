package kinglogger

import (
	"fmt"
	"os"
	"time"
)

//  往日志文件里写日志信息

// 创建结构体
type ConsoleLogger struct {
	level     Level         // 日志级别
}

// 为FileLogger结构体造一个构造函数
func NewConsoleLogger(levelStr string) *ConsoleLogger {
	logLevel := parseLogLevel(levelStr)
	fl := &ConsoleLogger{
		level: logLevel,
	}
	return fl
}


// 将公共日志记录的功能封装成一个方法
func (f *ConsoleLogger)log(level Level,format string, args ...interface{}){
	if f.level > level {
		return
	}
	// 日志格式：[时间][文件:行号][函数名][日志级别] 日志
	msg := fmt.Sprintf(format,args...) // 得到用户需要记录的日志
	nowStr := time.Now().Format("2006-01-02 15:04:05.000")
	fileName,line,funcName := getCallerInfo(3)
	logLevelStr := getLevelStr(level)
	logMsg := fmt.Sprintf("[%s][%s:%d][%s][%s] %s",nowStr,fileName,line,funcName,logLevelStr,msg)
	fmt.Fprintln(os.Stdout,logMsg) // 利用fmt包将msg字符串写入f.stdoutFile文件中

}



// Debug方法
func (f *ConsoleLogger) Debug(format string, args ...interface{}){
	f.log(DebugLevel,format,args...)
}

// Info 方法
func (f *ConsoleLogger) Info(format string, args ...interface{}){
	f.log(InfoLevel,format,args...)
}

// Warn 方法
func (f *ConsoleLogger) Warn(format string, args ...interface{}){
	f.log(WarnLevel,format,args...)
}

// Error 方法
func (f *ConsoleLogger) Error(format string, args ...interface{}){
	f.log(ErrorLevel,format,args...)
}

// Fatal 方法
func (f *ConsoleLogger) Fatal(format string, args ...interface{}){
	f.log(FatalLevel,format,args...)
}

// Close 操作系统终端输出不需要关闭
func (f *ConsoleLogger)Close(){

}