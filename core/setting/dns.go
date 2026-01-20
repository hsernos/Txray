// core/setting/dns.go 负责 DNS 相关设置项的定义与操作
package setting

import (
	"Txray/core/setting/key"
	"errors"
	"github.com/spf13/viper"
)

// DNSPort 返回 DNS 服务端口
func DNSPort() int {
	return viper.GetInt(key.DNSPort)
}

// SetDNSPort 设置 DNS 服务端口
// port: 需要设置的端口值
// 返回错误信息（如果有的话）
func SetDNSPort(port int) error {
	if port < 0 || port > 65535 {
		return errors.New("http端口取值 0~65535")
	}
	viper.Set(key.DNSPort, port)
	return viper.WriteConfig()
}

// DNSDomestic 返回国内 DNS 服务器地址
func DNSDomestic() string {
	return viper.GetString(key.DNSDomestic)
}

// SetDNSDomestic 设置国内 DNS 服务器地址
// dns: 需要设置的 DNS 服务器地址
// 返回错误信息（如果有的话）
func SetDNSDomestic(dns string) error {
	viper.Set(key.DNSDomestic, dns)
	return viper.WriteConfig()
}

// DNSForeign 返回国外 DNS 服务器地址
func DNSForeign() string {
	return viper.GetString(key.DNSForeign)
}

// SetDNSForeign 设置国外 DNS 服务器地址
// dns: 需要设置的 DNS 服务器地址
// 返回错误信息（如果有的话）
func SetDNSForeign(dns string) error {
	viper.Set(key.DNSForeign, dns)
	return viper.WriteConfig()
}

// DNSBackup 返回备用 DNS 服务器地址
func DNSBackup() string {
	return viper.GetString(key.DNSBackup)
}

// SetDNSBackup 设置备用 DNS 服务器地址
// dns: 需要设置的 DNS 服务器地址
// 返回错误信息（如果有的话）
func SetDNSBackup(dns string) error {
	viper.Set(key.DNSBackup, dns)
	return viper.WriteConfig()
}
