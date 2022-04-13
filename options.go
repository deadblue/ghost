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

// TODO: Add more options
// Options in plan:
//   1. CORS option: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
//   2. TLS option
