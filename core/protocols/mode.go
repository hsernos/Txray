// core/protocols/mode.go 负责协议模式的定义与相关操作
package protocols

// Mode 类型用于表示不同的协议模式
type Mode string

// 定义支持的协议模式常量
const (
	ModeVMess        Mode = "VMess"        // VMess 协议
	ModeVMessAEAD    Mode = "VMessAEAD"    // VMess AEAD 协议
	ModeVLESS        Mode = "VLESS"        // VLESS 协议
	ModeTrojan       Mode = "Trojan"       // Trojan 协议
	ModeShadowSocks  Mode = "ShadowSocks"  // ShadowSocks 协议
	ModeShadowSocksR Mode = "ShadowSocksR" // ShadowSocksR 协议
	ModeSocks        Mode = "Socks"        // Socks 协议
	ModeSplitHTTP    Mode = "SplitHTTP"    // SplitHTTP 协议
	ModeXHTTP        Mode = "XHTTP"        // XHTTP 协议
	ModeHttpUpgrade  Mode = "HttpUpgrade"  // HttpUpgrade 协议
)

// String 方法将 Mode 类型转换为字符串
func (m Mode) String() string {
	return string(m)
}
