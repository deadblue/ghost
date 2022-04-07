package ghost

type Option interface{}

type optListen struct {
	network, address string
}

func ListenAt(network, address string) Option {
	return optListen{
		network: network,
		address: address,
	}
}
