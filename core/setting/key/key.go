// core/setting/key/key.go 负责定义所有设置项的 key 常量
package key

const (
	Mixed         = "mixed"          // 混合代理
	Socks         = "socks"          // socks5 代理
	Http          = "http"           // http 代理
	UDP           = "udp"            // udp 代理
	Sniffing      = "sniffing"       // 启用嗅探
	FromLanConn   = "from_lan_conn"  // 允许局域网连接
	Mux           = "mux"            // 多路复用
	AllowInsecure = "allow_insecure" // 允许不安全连接

	RoutingStrategy = "routing.strategy" // 路由策略
	RoutingBypass   = "routing.bypass"   // 绕过局域网和大陆

	DNSPort     = "dns.port"     // DNS 端口
	DNSDomestic = "dns.domestic" // 国内 DNS
	DNSForeign  = "dns.foreign"  // 国外 DNS
	DNSBackup   = "dns.backup"   // 备用 DNS

	TestURL     = "test.url"     // 测试 URL
	TestTimeout = "test.timeout" // 测试超时时间
	TestMinTime = "test.mintime" // 测试最小时间

	RunBefore = "run_before" // 运行前脚本
	PID       = "pid"        // 核心进程 id

	VersionMin = "version.min" // 版本最小值
	VersionMax = "version.max" // 版本最大值
)
