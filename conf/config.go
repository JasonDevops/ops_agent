package conf

// 配置初始化
import (
	"fmt"
	"ops_agent/pkg/kinglogger"
)

var (
	Logger                   kinglogger.Loggers // 初始化日志对象
	JwtUrl, IdUrl, ReportUrl string
)

// 初始化日志对象
func Init_Log(levelStr, logPath, logName string) {
	Logger = kinglogger.NewFileLogger(levelStr, logPath, logName)
}

// 将用户传入的url给拼接成新url（三个url,一个做jwt认证，一个是获取主机ID，一个是发送资产数据接口）
func Init_FlagProcess(url string) {
	JwtUrl = fmt.Sprintf("%s/api/v1/jwt/token", url)
	IdUrl = fmt.Sprintf("%s/api/v1/asset/id", url)
	ReportUrl = fmt.Sprintf("%s/api/v1/asset/report", url)
}
