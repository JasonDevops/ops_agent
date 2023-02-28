package kinglogger

var (
	LogStore Loggers
)
func Init(){
	LogStore = NewFileLogger("info","./logs/","app.log")
}