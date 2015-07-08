package lib

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"
)

func GetTLSConfig(caCert, sslCert, sslKey []byte) (*tls.Config, error) {
	// TLS config
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = true
	certPool := x509.NewCertPool()

	certPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = certPool
	cert, err := tls.X509KeyPair(sslCert, sslKey)
	if err != nil {
		return &tlsConfig, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	return &tlsConfig, nil
}

func GetSslCert(sslCertPath string) ([]byte, error) {
	sslCert, err := os.Open(sslCertPath)
	if err != nil {
		log.Fatalf("unable to open ssl certificate: %s", err)
	}
	if _, err := sslCert.Stat(); err != nil {
		log.Fatalf("ssl cert is not accessible: %s", err)
	}
	sslCertData, err := ioutil.ReadAll(sslCert)
	if err != nil {
		log.Fatalf("unable to read ssl certificate: %s", err)
	}
	return sslCertData, nil
}

func GetSslKey(sslKeyPath string) ([]byte, error) {
	sslKey, err := os.Open(sslKeyPath)
	if err != nil {
		log.Fatalf("unable to open ssl key: %s", err)
	}
	if _, err := sslKey.Stat(); err != nil {
		log.Fatalf("ssl key is not accessible: %s", err)
	}
	sslKeyData, err := ioutil.ReadAll(sslKey)
	if err != nil {
		log.Fatalf("unable to read ssl key: %s", err)
	}
	return sslKeyData, nil
}

func GetCaCert(caCertPath string) ([]byte, error) {
	caCert, err := os.Open(caCertPath)
	if err != nil {
		log.Fatalf("unable to open ca certificate: %s", err)
	}
	if _, err := caCert.Stat(); err != nil {
		log.Fatalf("ca cert is not accessible: %s", err)
	}
	caCertData, err := ioutil.ReadAll(caCert)
	if err != nil {
		log.Fatalf("unable to read ca certificate: %s", err)
	}
	return caCertData, nil
}
