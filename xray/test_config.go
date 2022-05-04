package xray

import (
	"Txray/core"
	"Txray/core/protocols"
	"Txray/log"
	"path/filepath"
)

// 生成xray-core配置文件
func GenTestConfig(node protocols.Protocol) string {
	path := filepath.Join(core.GetConfigDir(), "config.json")
	var conf = map[string]interface{}{
		"inbounds":  inboundsConfig(),
		"outbounds": outboundConfig(node),
	}
	err := core.WriteJSON(conf, path)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return path
}