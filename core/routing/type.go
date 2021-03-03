package routing

type Type string

const (
	TypeProxy  Type = "Proxy"
	TypeDirect Type = "Direct"
	TypeBlock  Type = "Block"
)
