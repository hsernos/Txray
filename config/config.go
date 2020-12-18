package config

import (
	log "Tv2ray/logger"
	"Tv2ray/tools"
	"Tv2ray/vmess"
	"os/exec"
)

type (
	// setting 基本设置
	setting struct {
		Port                  uint   `json:"port"`
		Http                  uint   `json:"http"`
		UDP                   bool   `json:"udp"`
		Sniffing              bool   `json:"sniffing"`
		Mux                   bool   `json:"mux"`
		AllowLANConn          bool   `json:"allowLANConn"`
		BypassLanAndContinent bool   `json:"bypassLanAndContinent"`
		DomainStrategy        string `json:"domainStrategy"`
	}

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
		ConfigVersion  string `json:"configVersion"`
		Address        string `json:"address"`
		Port           uint   `json:"port"`
		ID             string `json:"id"`
		AlterID        int    `json:"alterId"`
		Security       string `json:"security"`
		Network        string `json:"network"`
		Remarks        string `json:"remarks"`
		HeaderType     string `json:"headerType"`
		RequestHost    string `json:"requestHost"`
		Path           string `json:"path"`
		StreamSecurity string `json:"streamSecurity"`
		TestResult     string `json:"testResult"`
		Subid          string `json:"subid"`
	}

	// sub 订阅
	sub struct {
		ID      string `json:"id"`
		Remarks string `json:"remarks"`
		URL     string `json:"url"`
		Using   bool   `json:"using"`
	}

	routing struct {
		Domain []string `json:"domain"`
		IP     []string `json:"ip"`
	}

	// Config 配置文件
	Config struct {
		exeCmd     *exec.Cmd
		Settings   setting  `json:"setting"`
		KcpSetting kcp      `json:"kcpSetting"`
		Index      uint     `json:"index"`
		Nodes      []*node  `json:"nodes"`
		Subs       []*sub   `json:"subs"`
		DNS        []string `json:"dns"`
		Proxy      routing  `json:"proxy"`
		Direct     routing  `json:"direct"`
		Block      routing  `json:"block"`
	}
)

func (c *Config) init() {
	configPath := tools.PathJoin(tools.GetRunPath(), "Tv2ray.json")
	ok := tools.IsFile(configPath)
	if ok {
		err := tools.ReadJSON(configPath, c)
		if err != nil {
			log.Error(err)
		}
	} else {
		c.Settings.Port = 2333
		c.Settings.Http = 2334
		c.Settings.UDP = true
		c.Settings.Sniffing = true
		c.Settings.Mux = true
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

		c.DNS = []string{"223.5.5.5", "223.6.6.6"}
	}

}

// 获取实例
func NewConfig() Config {
	c := Config{}
	c.init()
	return c
}

// 将数据保存到json文件
func (c *Config) SaveJSON() {
	file := tools.PathJoin(tools.GetRunPath(), "Tv2ray.json")
	err := tools.WriteJSON(*c, file)
	if err != nil {
		log.Error(err)
	}
}

// 节点数据转化
func nodeToVmessobj(n *node) *vmess.Vmess {
	v := vmess.Vmess{}
	v.V = n.ConfigVersion
	v.Type = n.HeaderType
	v.ID = n.ID
	v.Net = n.Network
	v.Path = n.Path
	v.Port = n.Port
	v.PS = n.Remarks
	v.Host = n.RequestHost
	v.TLS = n.StreamSecurity
	v.Add = n.Address
	v.Aid = n.AlterID
	return &v
}

func vmessObjToNode(vmessObj *vmess.Vmess, subid string) *node {
	n := node{}
	n.ConfigVersion = vmessObj.V
	n.HeaderType = vmessObj.Type
	n.ID = vmessObj.ID
	n.Network = vmessObj.Net
	n.Path = vmessObj.Path
	n.Port = vmessObj.Port
	n.Remarks = vmessObj.PS
	n.RequestHost = vmessObj.Host
	n.Security = "auto"
	n.StreamSecurity = vmessObj.TLS
	n.Address = vmessObj.Add
	n.AlterID = vmessObj.Aid
	n.Subid = subid
	return &n
}
