// core/setting/init.go 负责设置模块的初始化流程
package setting

import (
	"Txray/core"
	"Txray/core/setting/key"
	"Txray/log"
	"os"

	"github.com/spf13/viper"
)

func init() {
	// 配置文件不存在则创建
	if _, err := os.Stat(core.SettingFile); os.IsNotExist(err) {
		file, _ := os.Create(core.SettingFile)
		_ = file.Close()
	}
	viper.SetConfigName("setting")
	viper.SetConfigType("toml")
	viper.AddConfigPath(core.GetConfigDir())
	// 设置默认值
	// Xray-core 一直都支持代理协议来端口转发 https://github.com/XTLS/Xray-core/pull/4968
	viper.SetDefault(key.Mixed, 1025)
	viper.SetDefault(key.Socks, 0)
	viper.SetDefault(key.Http, 0)
	viper.SetDefault(key.UDP, true)
	viper.SetDefault(key.Sniffing, true)
	viper.SetDefault(key.FromLanConn, false)
	viper.SetDefault(key.Mux, false)
	viper.SetDefault(key.AllowInsecure, false)

	viper.SetDefault(key.RoutingStrategy, "IPIfNonMatch") //路由策略
	viper.SetDefault(key.RoutingBypass, true)             // 绕过局域网和大陆

	viper.SetDefault(key.DNSPort, 13500)
	viper.SetDefault(key.DNSForeign, "1.1.1.1")
	viper.SetDefault(key.DNSDomestic, "223.6.6.6")
	viper.SetDefault(key.DNSBackup, "114.114.114.114")

	viper.SetDefault(key.TestURL, "https://www.youtube.com")
	viper.SetDefault(key.TestTimeout, 5)
	viper.SetDefault(key.TestMinTime, 1000)
	viper.SetDefault(key.RunBefore, "")

	viper.SetDefault(key.PID, 0)

	viper.SetDefault(key.VersionMin, "") // 版本最小值，默认为空（不限制）
	viper.SetDefault(key.VersionMax, "") // 版本最大值，默认为空（不限制）

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err)
	}
}
