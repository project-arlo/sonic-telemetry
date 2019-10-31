package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	gnoi_system_pb "github.com/openconfig/gnoi/system"
	spb "proto/gnoi"
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
	target = flag.String("target", "localhost:8080", "Address:port of gNOI Server")
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
	
	switch *module {
	case "System":
		sc := gnoi_system_pb.NewSystemClient(conn)
		switch *rpc {
		case "Time":

			systemTime(sc, ctx)
		}
	case "Sonic":
		sc := spb.NewSonicServiceClient(conn)
		switch *rpc {
		case "showtechsupport":
			sonicShowTechSupport(sc, ctx)
		case "my-echo":
			sonicMyEcho(sc, ctx)
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
func sonicShowTechSupport(sc spb.SonicServiceClient, ctx context.Context) {
	fmt.Println("Sonic ShowTechsupport")
	resp,err := sc.ShowTechsupport(ctx, new(spb.TechsupportRequest))
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp.OutputFilename)
}
func sonicMyEcho(sc spb.SonicServiceClient, ctx context.Context) {
	fmt.Println("Sonic my-echo")
	resp,err := sc.MyEcho(ctx, new(spb.MyEchoRequest))
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp.Message)
}