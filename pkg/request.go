package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"ops_agent/conf"
)

// 服务端返回的json数据结构
type Server_Resp_Data struct {
	Data  string `json:data`
	Msg   string `json:msg`
	Error string `json:error`
	Code  int    `json:code`
}

// 服务端返回的ID
type Server_Resp_ID struct {
	Data  int    `json:data`
	Msg   string `json:msg`
	Error string `json:error`
	Code  int    `json:code`
}

// 客户端json数据结构（资产数据）
type Store_Asset_Data struct {
	Id   int                      `json:"id"`
	Type string                   `json:"type"`
	Data []map[string]interface{} `json:"data"`
}

// 客户端json数据结构（监控数据）
type Store_Monitor_Data struct {
	Id   int                      `json:"id"`
	Type string                   `json:"type"`
	Data []map[string]interface{} `json:"data"`
}

// 封装jwt认证发送请求函数
func Request_Evo(url string, ct string, method string, data []byte) (resp *[]byte, err error) {
	// 实现post请求
	conf.Logger.Info("发送http请求--url：%s 方法：%s", url, method)
	if method == "POST" || method == "post" {
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			conf.Logger.Info("发送http请求失败")
			return nil, err
		}
		body, _ := ioutil.ReadAll(resp.Body)
		return &body, nil
	}
	return
}

// 封装发送资产数据的提交函数
func Request_Asset(url, Token string, data Store_Asset_Data) {
	bs, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bs)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		conf.Logger.Error("请求server端失败...")
		return
	}
	request.Header.Add("token", Token)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	body, _ := ioutil.ReadAll(resp.Body)

	var resp_json *Server_Resp_Data
	json.Unmarshal(body, &resp_json)

}
