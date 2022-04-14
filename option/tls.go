package option

import (
	"crypto/tls"
	"crypto/x509"
)

type TlsOption struct {
	Config *tls.Config
}

func (o *TlsOption) isOption() {}

func (o *TlsOption) Certificate(certPem, keyPem []byte) *TlsOption {
	if o.Config == nil {
		o.Config = &tls.Config{}
	}
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err == nil {
		o.Config.Certificates = []tls.Certificate{cert}
	}
	return o
}

func (o *TlsOption) ClientCa(caPem []byte) *TlsOption {
	pool := x509.NewCertPool()
	if pool.AppendCertsFromPEM(caPem) {
		o.Config.ClientCAs = pool
	}
	return o
}
