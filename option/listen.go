package option

import (
	"fmt"
)

type ListenOption struct {
	Network string
	Address string
}

func (o *ListenOption) isOption() {}

func Listen(network, address string) Option {
	return &ListenOption{Network: network, Address: address}
}

func ListenPort(port uint16) Option {
	return ListenTcp(fmt.Sprintf(":%d", port))
}

func ListenIpAndPort(ip string, port uint16) Option {
	return ListenTcp(fmt.Sprintf("%s:%d", ip, port))
}

func ListenTcp(addr string) Option {
	return &ListenOption{Network: "tcp", Address: addr}
}
