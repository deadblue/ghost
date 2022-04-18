//go:build darwin || linux

package option

func ListenUnix(path string) Option {
	return &ListenOption{Network: "unix", Address: path}
}
