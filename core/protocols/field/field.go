package field

type Field struct {
	Key   string // 字段名
	Value string // 默认值
}

func NewField(key, value string) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}

var (
	VLessEncryption = NewField("encryption", "none") // 加密, VLESS可选值 none
	VMessEncryption = NewField("encryption", "auto") // 加密,  VMess可选值 auto/aes-128-gcm/chacha20-poly1305/none
	Flow            = NewField("flow", "")           // XTLS 的流控方式，可选值xtls-rprx-direct/xtls-rprx-splice

	// ==================================  协议的传输方式 =====================================
	NetworkType Field = NewField("type", "raw") // 协议的传输方式, 可选值 raw(tcp)/kcp/ws/http/quic/grpc/xhttp
	// Raw
	RawHeaderType = NewField("headerType", "none")
	RawHost       = NewField("host", "")
	RawPath       = NewField("path", "")

	// WebSocket
	WsPath = NewField("path", "/")
	WsHost = NewField("host", "")

	// H2
	H2Path = NewField("path", "/")
	H2Host = NewField("host", "")

	// mKCP
	KcpMtu = NewField("mtu", "1350")

	// gRPC
	GrpcServiceName = NewField("serviceName", "")
	GrpcAuthority   = NewField("authority", "")
	GrpcMode        = NewField("mode", "gun") // gRPC 的传输模式, 可选值 gun/multi/guna

	// xhttp
	XhttpHost  = NewField("host", "")
	XhttpPath  = NewField("path", "")
	XhttpMode  = NewField("mode", "")
	XhttpExtra = NewField("extra", "")

	// HTTPUpgrade
	HttpUpgradeHost = NewField("host", "")
	HttpUpgradePath = NewField("path", "")

	// ====================================   底层传输所使用的TLS类型===========================
	TlsSecurity = NewField("security", "none") // 设定底层传输所使用的 TLS 类型, 可选值有 none/tls/reality
	// 公共
	SNI         = NewField("sni", "") // TLS SNI
	FingerPrint = NewField("fp", "")  // TLS Client Hello 指纹

	// TLS
	Alpn = NewField("alpn", "") // alpn 多选 h2,http/1.1

	// REALITY
	PublicKey     = NewField("pbk", "") // REALITY的公钥
	ShortId       = NewField("sid", "") // REALITY 的 ID
	SpiderX       = NewField("spx", "") // REALITY 的爬虫
	Mldsa65Verify = NewField("pqv", "") // REALITY  mldsa65签名验证使用的公钥

	// finalmask
	Finalmask   = NewField("fm", "")
	Mport       = NewField("mport", "")
	Hopinterval = NewField("hopinterval", "60")
)
