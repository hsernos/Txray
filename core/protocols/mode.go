package protocols

type Mode string

const (
	ModeVMess        Mode = "vmess"
	ModeVLESS             = "vless"
	ModeTrojan            = "trojan"
	ModeShadowSocks       = "ss"
	ModeShadowSocksR      = "ssr"
)
