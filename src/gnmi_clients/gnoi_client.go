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
)

func main() {
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
	tr := new(gnoi_system_pb.TimeRequest)
	resp,err := sc.Time(ctx, tr)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp.Time)
}