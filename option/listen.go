package option

import (
	"fmt"
)

type ListenOption struct {
	Network string
	Address string
}

func (o ListenOption) isOption() {}

func ListenTcp(ip string, port uint16) Option {
	addr := fmt.Sprintf("%s:%d", ip, port)
	return ListenAddr(addr)
}

func ListenAddr(addr string) Option {
	return ListenOption{Network: "tcp", Address: addr}
}
