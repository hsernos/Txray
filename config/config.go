package config

import (
	"os/exec"
	"v3ray/tool"
)

type (
	// setting 基本设置
	setting struct {
		Port                  uint   `json:"port"`
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
		AllowInsecure  string `json:"allowInsecure"`
		ConfigType     int    `json:"configType"`
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

	path struct {
		Etc   string `json:"etc"`
		V2ray string `json:"v2ray"`
	}

	// Config 配置文件
	Config struct {
		exeCmd     *exec.Cmd
		Path       path     `json:"path"`
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
	configPath := tool.Join(c.getConfigPath() , "v3ray.json")
	ok := tool.PathExists(configPath)
	if ok {
		tool.ReadJSON(configPath, c)
	} else {
		c.Settings.Port = 2333
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

		c.DNS = []string{"8.8.4.4", "114.114.114.114", "localhost"}
	}

}

func (c *Config) getConfigPath() string {
	return tool.ConfigPath()
}




// NewConfig 得到实例
func NewConfig() Config {
	c := Config{}
	c.init()
	return c
}

// SaveJSON 将数据保存到json文件
func (c *Config) SaveJSON() {
	file := tool.Join(c.getConfigPath() , "v3ray.json")
	tool.WriteJSON(*c, file)
}

func nodeToVmessobj(n *node) *tool.Vmess {
	v := tool.Vmess{}
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

func vmessObjToNode(vmessObj *tool.Vmess, subid string) *node {
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
