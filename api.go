package ghost

const (
	DefaultNetwork = "tcp"
	DefaultAddress = "127.0.0.1:9057"
)

// Born creates a Shell with your ghost, which will listen at default
// network and address.
func Born[Ghost any](ghost Ghost) Shell {
	return BornAt(ghost, DefaultNetwork, DefaultAddress)
}

// BornAt creates a Shell with your ghost, which will listen at the
// specific network and address.
func BornAt[Ghost any](ghost Ghost, network, address string) Shell {
	return createShell(network, address, createKernel(ghost))
}
