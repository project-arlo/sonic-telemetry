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
	"encoding/json"
	"strings"
	"google.golang.org/grpc/metadata"
)

var (
	module = flag.String("module", "System", "gNOI Module")
	rpc = flag.String("rpc", "Time", "rpc call in specified module to call")
	target = flag.String("target", "localhost:8080", "Address:port of gNOI Server")
	args = flag.String("jsonin", "", "RPC Arguments in json format")
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
		case "sum":
			sonicSum(sc, ctx)
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
func sonicShowTechSupport(sc spb.SonicServiceClient, ctx context.Context) {
	fmt.Println("Sonic ShowTechsupport")
	req := &spb.TechsupportRequest {
		Input: &spb.TechsupportRequest_Input{

		},
	}
	nargs := strings.Replace(string(*args), "sonic-tests:input", "input", 1)
	json.Unmarshal([]byte(nargs), &req)
	fmt.Println(req)
	resp,err := sc.ShowTechsupport(ctx, req)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp.Output.OutputFilename)
}
func sonicSum(sc spb.SonicServiceClient, ctx context.Context) {
	fmt.Println("Sonic sum")
	req := &spb.SumRequest{
		Input: &spb.SumRequest_Input{

		},
	}
	nargs := strings.Replace(string(*args), "sonic-tests:input", "input", 1)
	json.Unmarshal([]byte(nargs), &req)
	fmt.Println(req)

	resp,err := sc.Sum(ctx, req)

	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp.Output.Result)
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
