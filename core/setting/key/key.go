package key

const (
	Socks       = "socks"
	Http        = "http"
	UDP         = "udp"
	Sniffing    = "sniffing"
	FromLanConn = "from_lan_conn"
	Mux         = "mux"

	RoutingStrategy = "routing.strategy"
	RoutingBypass   = "routing.bypass" // 绕过局域网和大陆

	DNSPort     = "dns.port"
	DNSDomestic = "dns.domestic" // 国内dns
	DNSForeign  = "dns.foreign"  // 国外dns
	DNSBackup   = "dns.backup"   //备用dns

	TestURL     = "test.url"
	TestTimeout = "test.timeout"

	RunBefore = "run_before"
	PID = "pid" // 核心进程id
)
