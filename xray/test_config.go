// xray/test_config.go 负责 xray 配置测试相关功能
package xray

import (
	"Txray/core"
	"Txray/core/protocols"
	"Txray/log"
	"path/filepath"
)

// 生成xray-core配置文件
// GenTestConfig 函数用于生成 xray-core 的配置文件，
// 参数 node 是一个实现了 protocols.Protocol 接口的对象，
// 返回值是生成的配置文件的路径。
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