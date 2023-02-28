package auth

import (
	"encoding/json"
	"fmt"
	"ops_agent/conf"
	"ops_agent/pkg"
	"time"
)

var (
	Jwt_Token string // 存储jwt token
)

// 实现jwt认证入口
func Jwt_Auth(url, username, password string, timeout int) {
	for ; timeout > 1; timeout-- {
		json_data := make(map[string]interface{})
		// 传入jwt认证用户名和密码
		json_data["username"] = username
		json_data["password"] = password
		bs, err := json.Marshal(json_data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		resp, err := pkg.Request_Evo(url, "application/json", "POST", bs)
		if err != nil {
			continue
		}
		// 解析服务端返回的json数据
		// 判断token是否获取，如果没有token则再次重试
		var resp_json pkg.Server_Resp_Data
		json.Unmarshal(*resp, &resp_json)
		if resp_json.Data != "" {
			// 将获取到的token赋值给Jwt_Token
			Jwt_Token = resp_json.Data
			conf.Logger.Info("jwt认证成功, msg: %s", resp_json.Msg)
			break
		}
		conf.Logger.Info("jwt认证失败, msg: %s", resp_json.Msg)
		time.Sleep(time.Second)
	}
}
