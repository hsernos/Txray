package core

import (
	"Txray/core/protocols/mode"
	"Txray/log"
)

func (c Core) GenConfig() (string, []string) {
	exe := ""
	switch c.GetNodeMode() {
	case mode.Trojan:
		exe = "xray"
	case mode.ShadowSocks:
		exe = "xray"
	case mode.VMess:
		exe = "xray"
	//case mode.ShadowSocksR:
	default:
		log.Warn("暂不支持 ", c.GetNodeMode(), " 协议")
		return "", nil
	}

	switch exe {
	case "xray":
		if !CheckXrayFile() {
			return "", nil
		}
		return GetXrayPath(), []string{
			"-config",
			c.GenXrayConfig(),
		}
	case "v2ray":
		if !CheckV2rayFile() {
			return "", nil
		}
		return GetV2rayPath(), []string{
			"-config",
			c.GenV2rayConfig(),
		}
	}
	return "", nil
}
