package setting

import (
	"Txray/core/config"
	"Txray/log"
	"Txray/tools"
)

type (
	// 基本设置
	base struct {
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

	dns struct {
		Port    uint   `json:"port"`
		Outland string `json:"outland"`
		Inland  string `json:"inland"`
		Backup  string `json:"backup"`
	}
	test struct {
		Url     string `json:"url"`
		TimeOut uint   `json:"timeout"`
	}
)

type Setting struct {
	Base       base `json:"base"`
	KcpSetting kcp  `json:"kcp"`
	DNS        dns  `json:"dns"`
	Test       test `json:"test"`
}

var setting *Setting = initSetting()

func initSetting() *Setting {
	s := new(Setting)
	file := config.Setting
	if tools.IsFile(file) {
		err := tools.ReadJSON(file, s)
		if err != nil {
			log.Error(err)
			return nil
		}
	} else {
		s.Base.Socks = 2333
		s.Base.Http = 0
		s.Base.UDP = true
		s.Base.Sniffing = false
		s.Base.Mux = false
		s.Base.AllowLANConn = false
		s.Base.BypassLanAndContinent = true
		s.Base.DomainStrategy = "IPIfNonMatch"

		s.KcpSetting.Mtu = 1350
		s.KcpSetting.Tti = 50
		s.KcpSetting.UplinkCapacity = 12
		s.KcpSetting.DownlinkCapacity = 100
		s.KcpSetting.Congestion = false
		s.KcpSetting.ReadBufferSize = 2
		s.KcpSetting.WriteBufferSize = 2

		s.DNS.Port = 23333
		s.DNS.Outland = "1.1.1.1"
		s.DNS.Inland = "119.29.29.29"
		s.DNS.Backup = ""

		s.Test.TimeOut = 5
		s.Test.Url = "https://www.youtube.com"

		s.save()
	}
	return s
}

func (s *Setting) save() {
	file := config.Setting
	err := tools.WriteJSON(*s, file)
	if err != nil {
		log.Error(err)
	}
}

func GetSetting() *Setting {
	return setting
}
