package server

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func setCert(serverKeyPath, serverPemPath, caPemPath string) grpc.ServerOption {
	cert, err := tls.LoadX509KeyPair(serverPemPath, serverKeyPath)
	if err != nil {
		log.Panic(err)
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caPemPath)
	if err != nil {
		log.Panic(err)
	}
	certPool.AppendCertsFromPEM(ca)

	// CA certificate can not be used, that is, the server certificate can be self signed
	//cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	//if err != nil {
	//	log.Panicf("Failed to parse certificate:", err)
	//}
	//certPool.AddCert(cert.Leaf)

	cred := credentials.NewTLS(&tls.Config{
		// Set the certificate chain to allow one or more certificates to be included
		Certificates: []tls.Certificate{cert},
		// The client is required to carry a certificate and the client will authenticate (multiple options are supported)
		ClientAuth: tls.RequireAndVerifyClientCert,
		// Set the collection of root certificates. The verification method uses the policy set in ClientAuth
		ClientCAs: certPool,
	})
	return grpc.Creds(cred)
}
