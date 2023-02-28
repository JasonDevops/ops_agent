package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"io/ioutil"
	"net/http"
	"ops_agent/conf"
	"ops_agent/pkg"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// 获取系统资源入口

var (
	Asset_ID  int // 资产ID
	AssetChan chan pkg.Store_Asset_Data
)

const (
	Cpu    = "cpu"
	Memory = "memory"
	Disk   = "disk"
)

// 获取系统基础信息
func Asset_Host_Base(FileName, Url, Token string, timeout int) {
	// 1.先获取到主机基础信息
	hostnamae_cmd, err := exec.Command("/bin/bash", "-c", "hostname").Output()
	if err != nil {
		conf.Logger.Error("获取主机名失败，程序退出...")
		os.Exit(1)
	}
	sn_cmd, err := exec.Command("/bin/bash", "-c", `dmidecode -t 1 | grep "UUID" | awk -F "[: ]+" '{print $2}'`).Output()
	if err != nil {
		conf.Logger.Error("获取sn码失败，程序退出...")
		os.Exit(1)
	}

	_hostname := string(hostnamae_cmd)
	_sn := string(sn_cmd)

	// 最后拿到主机名和sn
	hostname := _hostname
	sn := _sn

	conf.Logger.Info("获取 主机名: %s，sn：%s 成功", _hostname, _sn)

	// 先判断本地是否有保存asset_id的文件（返回nil表示文件存在）
	_, err = os.Stat(FileName)
	if err == nil {
		// 读取文件内容
		conf.Logger.Info("ID存储文件存在，路径为：%s，直接从文件中读取ID", FileName)
		read_cmd := fmt.Sprintf("cat  %v | xargs echo -n", FileName)
		id_cmd, _ := exec.Command("/bin/bash", "-c", read_cmd).Output()
		id_str := string(id_cmd)
		// 将ID从string类型转换为int类型
		id, err := strconv.Atoi(id_str)
		if err != nil {
			conf.Logger.Error("ID从string类型转换为int类型失败，程序退出...")
			os.Exit(1)
		}
		Asset_ID = id
		conf.Logger.Info("读取文件的ID为：%v", Asset_ID)
		return
	} else if os.IsNotExist(err) {
		conf.Logger.Info("ID存储文件路径：%s 不存在，需向server端请求获取ID", FileName)
		for ; timeout > 1; timeout-- {
			time.Sleep(time.Second * 1)
			data := make(map[string]interface{})
			data["hostname"] = hostname
			data["sn"] = sn
			bs, err := json.Marshal(data)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			reader := bytes.NewReader(bs)
			request, err := http.NewRequest("POST", Url, reader)
			request.Header.Add("token", Token)
			request.Header.Set("Content-Type", "application/json;charset=UTF-8")
			client := http.Client{}
			resp, err := client.Do(request)
			if err != nil {
				continue
			}
			body, _ := ioutil.ReadAll(resp.Body)

			var resp_json pkg.Server_Resp_ID
			json.Unmarshal(body, &resp_json)

			// 服务端返回code
			// -1 ：获取id失败
			// 0  ：获取id成功
			if resp_json.Code == -1 {
				conf.Logger.Error("从server端获取ID失败,尝试再次获取")
				continue
			}
			// 将id写到指定路径的文件中
			write_cmd := fmt.Sprintf("echo %v > %v", resp_json.Data, FileName)
			_, err = exec.Command("/bin/bash", "-c", write_cmd).Output()
			if err != nil {
				conf.Logger.Error("ID信息写入保存文件 %s 失败", FileName)
				os.Exit(1)
			}
			conf.Logger.Error("创建文件 %s，并写入ID", FileName)
			Asset_ID = resp_json.Data
			break
			defer resp.Body.Close()
		}
	}
}

type Asset struct {
}

// gorouting运行函数获取的数据写入channel管道
// channel管道初始化时支持四种模式（slow、hight、fast 、quick）
func (a *Asset) InitAssetChan(mode string) {
	// 根据模式获取对应的管道缓冲区大小
	level := pkg.ChanModeLevel(mode)
	AssetChan = make(chan pkg.Store_Asset_Data, level) // 通道初始化(带缓冲区)
}

func RunAsset() {
	exec_asset := Asset{}
	// 初始化管道
	exec_asset.InitAssetChan("slow")

	for {
		go exec_asset.AssetCPU()
		go exec_asset.AssetMemory()
		go exec_asset.AssetDisk()
		time.Sleep(time.Second * 10)
	}
}

func (a *Asset) AssetCPU() {
	/*
		cpu_logic_number: cpu逻辑个数
		cpu_physics_number: cpu物理个数
		cpu_model: cpu型号

		data3 = {
		    "id":id,
		    "type": "cpu",
		    "data": [
		        {"cpu_logic_number": "4"},
		        {"cpu_physics_number": "8"},
		        {"cpu_model": "inter 233"}
		    ]
		}

	*/
	// data用来存储获取的到信息
	// list_data 主要是将data里的数据添加到列表中
	data := make(map[string]interface{})
	var list_data []map[string]interface{}
	// 数据结构：data = {cpu_logic_number:xx,cpu_physics_number:xx,cpu_model:xx}
	// 将map追加到列表中：[{cpu_logic_number:xx,cpu_physics_number:xx,cpu_model:xx},]

	cpu_sum, _ := cpu.Info()
	conf.Logger.Info("获取CPU信息成功 %s", list_data)
	for _, val := range cpu_sum {
		data["cpu"] = val.CPU
		data["cpu_cores"] = val.Cores
		data["cpu_model"] = val.ModelName
		list_data = append(list_data, data)
	}
	commit_data := pkg.Store_Asset_Data{
		Id:   Asset_ID,
		Type: Cpu,
		Data: list_data,
	}

	// channel存储数据
	AssetChan <- commit_data
}

func (a *Asset) AssetMemory() {
	/*
			memory_total:
			swap_total:
			memory_cached:
		    memory_buffer：
	*/
	mem_sum, _ := mem.VirtualMemory()
	data := make(map[string]interface{})
	var list_data []map[string]interface{}
	data["total"] = (mem_sum.Total / 1024 / 1024)
	data["used"] = (mem_sum.Used / 1024 / 1024)
	data["free"] = (mem_sum.Free / 1024 / 1024)
	data["cached"] = (mem_sum.Cached / 1024 / 1024)
	data["buffers"] = (mem_sum.Buffers / 1024 / 1024)

	list_data = append(list_data, data)
	conf.Logger.Info("获取Memory信息成功 %s", data)

	commit_data := pkg.Store_Asset_Data{
		Id:   Asset_ID,
		Type: Memory,
		Data: list_data,
	}

	// 测试channel存储数据
	AssetChan <- commit_data
}

func (a *Asset) AssetDisk() {

	// 获取到本机所有disk信息
	var list_data []map[string]interface{}
	ds, _ := disk.Partitions(true)

	for _, val := range ds {
		// 跳过一些无用的磁盘信息
		match := strings.HasPrefix(val.Device, "/")
		if !match {
			continue
		}
		// 获取指定磁盘的具体信息
		var ds_info *disk.UsageStat
		ds_info, _ = disk.Usage(val.Device)
		data := make(map[string]interface{})
		data["disk_path"] = ds_info.Path
		data["disk_total"] = strconv.Itoa(int(ds_info.Total))
		data["disk_free"] = strconv.Itoa(int(ds_info.Free))
		data["disk_system_type"] = ds_info.Fstype
		data["disk_mount"] = val.Mountpoint
		list_data = append(list_data, data)
	}
	conf.Logger.Info("获取磁盘信息成功: %s", list_data)

	commit_data := pkg.Store_Asset_Data{
		Id:   Asset_ID,
		Type: Disk,
		Data: list_data,
	}
	// channel存储disk数据
	AssetChan <- commit_data
	time.Sleep(time.Second * 10)
}
