package kinglogger

import (
	"fmt"
	"os"
	"path"
	"time"
)

//  往日志文件里写日志信息

// 创建结构体
type FileLogger struct {
	level     Level         // 日志级别
	logPath   string		// 日志文件路径
	logName   string        // 日志文件名字
	stdoutFile *os.File		// 标准日志文件句柄
	errorFile *os.File 		// 错误日志文件句柄
	maxSize     int64       // 日志文件大小
}

// 为FileLogger结构体造一个构造函数
func NewFileLogger(levelStr string,logPath,logName string) *FileLogger {
	logLevel := parseLogLevel(levelStr)
	fl := &FileLogger{
		level: logLevel,
		logPath: logPath,
		logName: logName,
		maxSize: 10 * 1024 * 1024, // 10m
	}
	fl.initFile()
	return fl
}


// 初始化方法（初始化stout and error 文件句柄）
func (f *FileLogger)initFile(){
	// 1、初始化标准日志文件句柄
	stdoutLogName := path.Join(f.logPath,f.logName) // 拼接 logPath and logName 为一个完整的路径
	stoutFileObj,err := os.OpenFile(stdoutLogName,os.O_CREATE|os.O_WRONLY|os.O_APPEND,0644)
	if err != nil{
		panic(fmt.Errorf("打开文件%s错误:%v",stdoutLogName,err))
	}
	f.stdoutFile = stoutFileObj

	// 2、初始化错误日志文件句柄
	errorLogName := fmt.Sprintf("%s.err",stdoutLogName)
	errorFileObj,err := os.OpenFile(errorLogName,os.O_CREATE|os.O_WRONLY|os.O_APPEND,0644)
	if err != nil{
		panic(fmt.Errorf("打开文件%s错误:%v",errorFileObj,err))
	}
	f.errorFile = errorFileObj
}

// 封装一个切分日志文件的方法
func (f *FileLogger)splitLogFile(file *os.File) *os.File{
	// 切分文件
	fileName := file.Name() // 获取到文件的完整路径（原始日志文件名）
	backupName := fmt.Sprintf("%s_%v.back",fileName,time.Now().Unix())
	// 1.把原来的文件句柄关闭
	file.Close()
	// 2.备份原来的文件
	os.Rename(fileName,backupName)
	// 3.新建一个文件
	fileObj,err := os.OpenFile(fileName,os.O_CREATE|os.O_WRONLY|os.O_APPEND,0664)
	if err != nil{
		panic(fmt.Errorf("打开日志文件%s失败,%v",fileName,err))
	}
	return fileObj
}

// 将公共日志记录的功能封装成一个方法
func (f *FileLogger)log(level Level,format string, args ...interface{}){
	if f.level > level {
		return
	}
	// 日志格式：[时间][文件:行号][函数名][日志级别] 日志
	msg := fmt.Sprintf(format,args...) // 得到用户需要记录的日志
	nowStr := time.Now().Format("2006-01-02 15:04:05.000")
	fileName,line,funcName := getCallerInfo(3)
	logLevelStr := getLevelStr(level)
	logMsg := fmt.Sprintf("[%s][%s:%d][%s][%s] %s",nowStr,fileName,line,funcName,logLevelStr,msg)

	// 往文件里写之前要做一个检查：检查当前日志文件大小是否超过了maxSize
	if f.checkSplit(f.stdoutFile){
		f.stdoutFile = f.splitLogFile(f.stdoutFile)
	}
	fmt.Fprintln(f.stdoutFile,logMsg) // 利用fmt包将msg字符串写入f.stdoutFile文件中

	// 如果是error或者fatal级别的日志还要记录到f.errorFILE
	if level >= ErrorLevel{
		if f.checkSplit(f.errorFile){
			f.errorFile = f.splitLogFile(f.errorFile)
		}
		fmt.Fprintln(f.errorFile,logMsg)
	}
}

// 检查日志是否要拆分的功能
func (f *FileLogger)checkSplit(file *os.File)bool{
	fileinfo,_ := file.Stat()
	fileSize := fileinfo.Size()
	return fileSize > f.maxSize // 当传进来的日志文件的大小超过了maxSize则返回true
}

// Debug方法
func (f *FileLogger) Debug(format string, args ...interface{}){
	f.log(DebugLevel,format,args...)
}

// Info 方法
func (f *FileLogger) Info(format string, args ...interface{}){
	f.log(InfoLevel,format,args...)
}

// Warn 方法
func (f *FileLogger) Warn(format string, args ...interface{}){
	f.log(WarnLevel,format,args...)
}

// Error 方法
func (f *FileLogger) Error(format string, args ...interface{}){
	f.log(ErrorLevel,format,args...)
}

// Fatal 方法
func (f *FileLogger) Fatal(format string, args ...interface{}){
	f.log(FatalLevel,format,args...)
}

// Close 关闭日志文件句柄
func (f *FileLogger)Close(){
	f.stdoutFile.Close()
	f.errorFile.Close()
}