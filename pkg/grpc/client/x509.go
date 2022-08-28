package client

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
)

func setCert(clientKeyPath, clientPemPath, caPemPath, commonName string) (opt grpc.DialOption, err error) {
	cert, err := tls.LoadX509KeyPair(clientPemPath, clientKeyPath)
	if err != nil {
		return
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caPemPath)
	if err != nil {
		return
	}
	certPool.AppendCertsFromPEM(ca)

	cred := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   commonName,
		RootCAs:      certPool,
	})
	return grpc.WithTransportCredentials(cred), nil
}
