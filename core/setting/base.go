package setting

import "Txray/log"

// 本地socks监听端口
func SocksPort() uint {
	return setting.Base.Socks
}

// 本地http监听端口
func HttpPort() uint {
	return setting.Base.Http
}

func IsUDP() bool {
	return setting.Base.UDP
}

// 设置本地监听端口
func SetSocksPort(port uint) {
	defer setting.save()
	setting.Base.Socks = port
	log.Info("设置本地socks5监听端口为 [", port, "]")
}

// 设置本地监听端口
func SetHttpPort(port uint) {
	defer setting.save()
	setting.Base.Http = port
	log.Info("设置本地http监听端口为 [", port, "]")
}

// 设置是否允许udp转发
func SetUDP(udp bool) {
	defer setting.save()
	setting.Base.UDP = udp
	if udp {
		log.Info("开启UDP转发")
	} else {
		log.Info("关闭UDP转发")
	}

}

// 设置是否允许流量监听
func SetSniffing(sniffing bool) {
	defer setting.save()
	setting.Base.Sniffing = sniffing
	if sniffing {
		log.Info("开启流量监听")
	} else {
		log.Info("关闭流量监听")
	}
}

// 设置是否允许多路复用
func SetMux(mux bool) {
	defer setting.save()
	setting.Base.Mux = mux
	if mux {
		log.Info("开启多路复用")
	} else {
		log.Info("关闭多路复用")
	}
}

// 设置是否允许局域网的连接
func SetLANConn(conn bool) {
	defer setting.save()
	setting.Base.AllowLANConn = conn
	if conn {
		log.Info("允许来着局域网的连接")
	} else {
		log.Info("禁止来着局域网的连接")
	}
}

// 设置是否绕过局域网和大陆
func SetBypassLanAndContinent(bypass bool) {
	defer setting.save()
	setting.Base.BypassLanAndContinent = bypass
	log.Info("设置是否绕过局域网和大陆为: ", bypass)
	if bypass {
		log.Info("开启绕过局域网和大陆")
	} else {
		log.Info("关闭绕过局域网和大陆")
	}
}

// 设置路由策略
func SetDomainStrategy(mode int) {
	defer setting.save()
	domainStrategy := setting.Base.DomainStrategy
	switch mode {
	case 1:
		domainStrategy = "AsIs"
	case 2:
		domainStrategy = "IPIfNonMatch"
	case 3:
		domainStrategy = "IPOnDemand"
	}
	setting.Base.DomainStrategy = domainStrategy
	log.Info("设置路由策略为 [", domainStrategy, "]")
}

func Base() base {
	return setting.Base
}
