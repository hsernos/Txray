package protocols

type Mode string

const (
	ModeVMess        Mode = "VMess"
	ModeVMessAEAD    Mode = "VMessAEAD"
	ModeVLESS        Mode = "VLESS"
	ModeTrojan       Mode = "Trojan"
	ModeShadowSocks  Mode = "ShadowSocks"
	ModeShadowSocksR Mode = "ShadowSocksR"
	ModeSocks        Mode = "Socks"
)

func (m Mode) String() string {
	return string(m)
}
