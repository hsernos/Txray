package core

import (
	"Txray/core/protocols"
	"Txray/log"
	"Txray/tools"
	"os/exec"
)

type (
	// 基本设置
	setting struct {
		Socks                 uint   `json:"socks"`
		Http                  uint   `json:"http"`
		UDP                   bool   `json:"udp"`
		Sniffing              bool   `json:"sniffing"`
		Mux                   bool   `json:"mux"`
		AllowLANConn          bool   `json:"allowLANConn"`
		BypassLanAndContinent bool   `json:"bypassLanAndContinent"`
		DomainStrategy        string `json:"domainStrategy"`
	}

	// kcp 设置
	kcp struct {
		Mtu              int  `json:"mtu"`
		Tti              int  `json:"tti"`
		UplinkCapacity   int  `json:"uplinkCapacity"`
		DownlinkCapacity int  `json:"downlinkCapacity"`
		Congestion       bool `json:"congestion"`
		ReadBufferSize   int  `json:"readBufferSize"`
		WriteBufferSize  int  `json:"writeBufferSize"`
	}

	// node 节点
	node struct {
		data       protocols.Protocol
		Link       string `json:"link"`
		TestResult string `json:"testResult"`
		Subid      string `json:"subid"`
	}

	// sub 订阅
	sub struct {
		ID      string `json:"id"`
		Remarks string `json:"remarks"`
		URL     string `json:"url"`
		Using   bool   `json:"using"`
	}

	dns struct {
		Port    uint   `json:"port"`
		Outland string `json:"outland"`
		Inland  string `json:"inland"`
		Backup  string `json:"backup"`
	}

	routing struct {
		Data string `json:"data"`
		Mode string `json:"mode"`
	}

	Core struct {
		exeCmd     *exec.Cmd
		Settings   setting    `json:"setting"`
		KcpSetting kcp        `json:"kcpSetting"`
		Index      uint       `json:"index"`
		Nodes      []*node    `json:"nodes"`
		Subs       []*sub     `json:"subs"`
		DNS        dns        `json:"dns"`
		Proxy      []*routing `json:"proxy"`
		Direct     []*routing `json:"direct"`
		Block      []*routing `json:"block"`
	}
)

func (c *Core) init() {
	configPath := tools.PathJoin(tools.GetRunPath(), "data.json")
	ok := tools.IsFile(configPath)
	if ok {
		err := tools.ReadJSON(configPath, c)
		if err != nil {
			log.Error(err)
			return
		}
		c.initNodesData()
	} else {
		c.Settings.Socks = 2333
		c.Settings.Http = 0
		c.Settings.UDP = true
		c.Settings.Sniffing = false
		c.Settings.Mux = false
		c.Settings.AllowLANConn = false
		c.Settings.BypassLanAndContinent = true
		c.Settings.DomainStrategy = "IPIfNonMatch"

		c.KcpSetting.Mtu = 1350
		c.KcpSetting.Tti = 50
		c.KcpSetting.UplinkCapacity = 12
		c.KcpSetting.DownlinkCapacity = 100
		c.KcpSetting.Congestion = false
		c.KcpSetting.ReadBufferSize = 2
		c.KcpSetting.WriteBufferSize = 2
		c.Index = 0

		c.DNS.Port = 23333
		c.DNS.Outland = "1.1.1.1"
		c.DNS.Inland = "223.5.5.5"
		c.DNS.Backup = ""
		c.Save()
	}

}

// 将数据保存到json文件
func (c *Core) Save() {
	file := tools.PathJoin(tools.GetRunPath(), "data.json")
	err := tools.WriteJSON(*c, file)
	if err != nil {
		log.Error(err)
	}
}

// 获取实例
func New() Core {
	c := Core{}
	c.init()
	return c
}
