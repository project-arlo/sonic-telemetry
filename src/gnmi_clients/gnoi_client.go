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
	"google.golang.org/grpc/metadata"
)

var (
	module = flag.String("module", "System", "gNOI Module")
	rpc = flag.String("rpc", "Time", "rpc call in specified module to call")
	target = flag.String("target", "localhost:8080", "Address:port of gNOI Server")
	username = flag.String("username", "", "Username if required")
	password = flag.String("password", "", "Password if required")
)

func main() {
	flag.Parse()
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
	conn, err := grpc.Dial(*target, opts...)
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
	ctx = setUserPass(ctx)
	resp,err := sc.Time(ctx, new(gnoi_system_pb.TimeRequest))
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp.Time)
}


func setUserPass(ctx context.Context) context.Context {
	if len(*username) > 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, "username", *username)
	}
	if len(*password) > 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, "password", *password)
	}
	return ctx
}