package setting

import (
	"Txray/core/setting/key"
	"errors"
	"github.com/spf13/viper"
)

func DNSPort() int {
	return viper.GetInt(key.DNSPort)
}

func SetDNSPort(port int) error {
	if port < 0 || port > 65535 {
		return errors.New("http端口取值 0~65535")
	}
	viper.Set(key.DNSPort, port)
	return viper.WriteConfig()
}

func DNSDomestic() string {
	return viper.GetString(key.DNSDomestic)
}

func SetDNSDomestic(dns string) error {
	viper.Set(key.DNSDomestic, dns)
	return viper.WriteConfig()
}

func DNSForeign() string {
	return viper.GetString(key.DNSForeign)
}

func SetDNSForeign(dns string) error {
	viper.Set(key.DNSForeign, dns)
	return viper.WriteConfig()
}

func DNSBackup() string {
	return viper.GetString(key.DNSBackup)
}

func SetDNSBackup(dns string) error {
	viper.Set(key.DNSBackup, dns)
	return viper.WriteConfig()
}
