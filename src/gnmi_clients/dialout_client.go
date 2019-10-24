package main

import (
	"google.golang.org/grpc"
	"errors"
	"proto"
	"google.golang.org/grpc/reflection"
	"crypto/tls"
	"context"
	"os"
	"os/signal"
	"fmt"
	"flag"
	"net"
)

var (
	address = flag.String("address", "0.0.0.0", "Address to listen on")
	port = flag.Int64("port", 9000, "Port to listen on")
	
)

type Server struct {
	s       *grpc.Server
	lis     net.Listener
	config  *Config
}
// Config is a collection of values for Server
type Config struct {
	// Port for the Server to listen on. If 0 or unset the Server will pick a port
	// for this Server.
	Port int64
	Addr string
}
func NewServer(config *Config, opts []grpc.ServerOption) (*Server, error) {
	if config == nil {
		return nil, errors.New("config not provided")
	}

	s := grpc.NewServer(opts...)

	reflection.Register(s)

	srv := &Server{
		s:       s,
		config:  config,
		// clients: map[string]*Client{},
	}
	var err error
	if srv.config.Port < 0 {
		srv.config.Port = 0
	}
	srv.lis, err = net.Listen("tcp", fmt.Sprintf(":%d", srv.config.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to open listener port %d: %v", srv.config.Port, err)
	}
	gnmi_sonic.RegisterGNMIDialOutServer(s, srv)
	
	log.V(1).Infof("Created Server on %s", srv.Address())

	return srv, nil
}
func main() {
	flag.Parse()
	tls_conf := tls.Config{InsecureSkipVerify: true}
	opts := []grpc.ServerOption{grpc.Creds(credentials.NewTLS(tls_conf))}



    ctx, cancel := context.WithCancel(context.Background())
    go func() {
            c := make(chan os.Signal, 1)
            signal.Notify(c, os.Interrupt)
            <-c
            cancel()
    }()
    cfg := &Config{Port: int64(*port), Addr: address}
    srv := NewServer(cfg, opts)

	if err != nil {
		panic("Failed to create collector server: %v", err)
		return
	}

	
	s.Serve() // blocks until close
	
	
	
}


