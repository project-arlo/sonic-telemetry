package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	gnoi_system_pb "github.com/openconfig/gnoi/system"
	"time"
	"math"
	"crypto/tls"
	"context"
	"os"
	"os/signal"
	"fmt"
	"flag"
)

var (
	module = flag.String("module", "System", "gNOI Module")
	rpc = flag.String("rpc", "Time", "rpc call in specified module to call")
)

func main() {
	flag.Parse()
	serverAddr := "localhost:8080"
	tls_conf := tls.Config{InsecureSkipVerify: true}
    opts := []grpc.DialOption{
            grpc.WithTimeout(time.Second * 3),
            grpc.WithBlock(),
            grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)),
            grpc.WithTransportCredentials(credentials.NewTLS(&tls_conf)),
    }

    ctx, cancel := context.WithCancel(context.Background())
    go func() {
            c := make(chan os.Signal, 1)
            signal.Notify(c, os.Interrupt)
            <-c
            cancel()
    }()
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		panic(err.Error())
	}
	sc := gnoi_system_pb.NewSystemClient(conn)

	switch *module {
	case "System":
		switch *rpc {
		case "Time":
			systemTime(sc, ctx)
		}
	}
}

func systemTime(sc gnoi_system_pb.SystemClient, ctx context.Context) {
	fmt.Println("System Time")
	resp,err := sc.Time(ctx, new(gnoi_system_pb.TimeRequest))
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp.Time)
}