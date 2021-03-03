package setting

import (
	"Txray/log"
	"strings"
)

// 设置DNS监听端口
func SetDNSPort(port uint) {
	defer setting.save()
	setting.DNS.Port = port
	log.Info("设置本地DNS监听端口为 [", port, "]")
}

// 设置境外DNS
func SetOutlandDNS(dns string) {
	defer setting.save()
	setting.DNS.Outland = dns
	log.Info("设置境外DNS为 [", dns, "]")
}

// 设置境内DNS
func SetInlandDNS(dns string) {
	defer setting.save()
	setting.DNS.Inland = dns
	log.Info("设置境内DNS为 [", dns, "]")
}

// 设置备用DNS端口（多个用英文逗号分隔）
func SetBackupDNS(dns string) {
	defer setting.save()
	setting.DNS.Backup = dns
	log.Info("设置备用DNS为 [", dns, "]")
}

// DNS监听端口
func DNSPort() uint {
	return setting.DNS.Port
}

// 境外DNS
func OutlandDNS() string {
	return setting.DNS.Outland
}

// 境内DNS
func InlandDNS() string {
	return setting.DNS.Inland
}

// 备用DNS端口
func BackupDNS() []string {
	return strings.Split(setting.DNS.Backup, ",")
}
