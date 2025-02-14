package tls

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"fmt"
)

var (
	//go:embed ca.pem
	CACert []byte
)

var (
	//go:embed cli.pem
	CliCert []byte
	//go:embed cli-key.pem
	CliKey []byte
)

func CliTLSConfig(serverName string) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(CliCert, CliKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load client cert pair: %v", err)
	}
	ca := x509.NewCertPool()
	if ok := ca.AppendCertsFromPEM(CACert); !ok {
		return nil, fmt.Errorf("failed to load client CA cert")
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      ca,
		ServerName:   serverName,
	}
	return tlsConfig, nil
}
