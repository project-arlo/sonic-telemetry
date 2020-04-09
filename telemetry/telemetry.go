package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"time"
	log "github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	gnmi "github.com/Azure/sonic-telemetry/gnmi_server"
	testcert "github.com/Azure/sonic-telemetry/testdata/tls"
)

var (
	userAuth = gnmi.AuthTypes{"password": true, "cert": false, "jwt": true}
	port = flag.Int("port", -1, "port to listen on")
	// Certificate files.
	caCert            = flag.String("ca_crt", "", "CA certificate for client certificate validation. Optional.")
	serverCert        = flag.String("server_crt", "", "TLS server certificate")
	serverKey         = flag.String("server_key", "", "TLS server private key")
	insecure          = flag.Bool("insecure", false, "Skip providing TLS cert and key, for testing only!")
	jwtRefInt         = flag.Uint64("jwt_refresh_int", 30, "Seconds before JWT expiry the token can be refreshed.")
	jwtValInt         = flag.Uint64("jwt_valid_int", 3600, "Seconds that JWT token is valid for.")
)

func main() {
	flag.Var(userAuth, "client_auth", "Client auth mode(s) - none,cert,password,jwt")
	flag.Parse()

	switch {
	case *port <= 0:
		log.Errorf("port must be > 0.")
		return
	}
	var certificate tls.Certificate
	var err error
	gnmi.JwtRefreshInt = time.Duration(*jwtRefInt*uint64(time.Second))
	gnmi.JwtValidInt = time.Duration(*jwtValInt*uint64(time.Second))
	if *insecure {
		certificate, err = testcert.NewCert()
		if err != nil {
			log.Exitf("could not load server key pair: %s", err)
		}
	} else {
		switch {
		case *serverCert == "":
			log.Errorf("serverCert must be set.")
			return
		case *serverKey == "":
			log.Errorf("serverKey must be set.")
			return
		}
		certificate, err = tls.LoadX509KeyPair(*serverCert, *serverKey)
		if err != nil {
			log.Exitf("could not load server key pair: %s", err)
		}
	}

	tlsCfg := &tls.Config{
		ClientAuth:   tls.RequestClientCert,
		Certificates: []tls.Certificate{certificate},
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{

			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},

	}

	if *caCert != "" {
		tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
		ca, err := ioutil.ReadFile(*caCert)
		if err != nil {
			log.Exitf("could not read CA certificate: %s", err)
		}
		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			log.Exit("failed to append CA certificate")
		}
		tlsCfg.ClientCAs = certPool
	} else {
		if userAuth.Enabled("cert") {
			log.Exit("client_auth mode cert requires ca_crt option.")
		}
	}

	opts := []grpc.ServerOption{grpc.Creds(credentials.NewTLS(tlsCfg))}
	cfg := &gnmi.Config{}
	cfg.Port = int64(*port)
	cfg.UserAuth = userAuth

	gnmi.GenerateJwtSecretKey()

	s, err := gnmi.NewServer(cfg, opts)
	if err != nil {
		log.Errorf("Failed to create gNMI server: %v", err)
		return
	}

	log.V(1).Infof("Starting RPC server on address: %s", s.Address())
	s.Serve() // blocks until close
	log.Flush()
}
