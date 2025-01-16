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
	//go:embed scheduler.pem
	SchedulerServerCert []byte
	//go:embed scheduler-key.pem
	SchedulerServerKey []byte
)

var (
	//go:embed worker-server.pem
	WorkerServerCert []byte
	//go:embed worker-server-key.pem
	WorkerServerKey []byte
)

var (
	//go:embed worker.pem
	WorkerCert []byte
	//go:embed worker-key.pem
	WorkerKey []byte
)

var (
	//go:embed cli.pem
	CliCert []byte
	//go:embed cli-key.pem
	CliKey []byte
)

func SchedulerTLSConfig(serverName string) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(SchedulerServerCert, SchedulerServerKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load server cert pair: %v", err)
	}
	ca := x509.NewCertPool()
	if ok := ca.AppendCertsFromPEM(CACert); !ok {
		return nil, fmt.Errorf("failed to load server CA cert")
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    ca,
		ServerName:   serverName,
	}
	return tlsConfig, nil
}

func WorkerServerTLSConfig(serverName string) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(WorkerServerCert, WorkerServerKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load server cert pair: %v", err)
	}
	ca := x509.NewCertPool()
	if ok := ca.AppendCertsFromPEM(CACert); !ok {
		return nil, fmt.Errorf("failed to load server CA cert")
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    ca,
		ServerName:   serverName,
	}
	return tlsConfig, nil
}

func WorkerTLSConfig(serverName string) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(WorkerCert, WorkerKey)
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
