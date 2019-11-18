package main

import (
	"google.golang.org/grpc"
	gnoi_system_pb "github.com/openconfig/gnoi/system"
	spb "proto/gnoi"
	"context"
	"os"
	"os/signal"
	"fmt"
	"flag"
	"encoding/json"
	"strings"
	"google.golang.org/grpc/metadata"
	"github.com/google/gnxi/utils/credentials"
)

var (
	module = flag.String("module", "System", "gNOI Module")
	rpc = flag.String("rpc", "Time", "rpc call in specified module to call")
	target = flag.String("target", "localhost:8080", "Address:port of gNOI Server")
	args = flag.String("jsonin", "", "RPC Arguments in json format")
	jwtToken = flag.String("jwt_token", "", "JWT Token if required")
	targetName = flag.String("target_name", "hostname.com", "The target name use to verify the hostname returned by TLS handshake")
)
func setUserCreds(ctx context.Context) context.Context {
	if len(*jwtToken) > 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, "access_token", *jwtToken)
	}
	return ctx
}
func main() {
	flag.Parse()
	opts := credentials.ClientCredentials(*targetName)

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
		default:
			panic("Invalid RPC Name")
		}
	case "Sonic":
		sc := spb.NewSonicServiceClient(conn)
		switch *rpc {
		case "showtechsupport":
			sonicShowTechSupport(sc, ctx)
		case "sum":
			sonicSum(sc, ctx)
		case "saveConfig":
			saveConfig(sc, ctx)
		case "reloadConfig":
			reloadConfig(sc, ctx)
		case "authenticate":
			authenticate(sc, ctx)
		case "refresh":
			refresh(sc, ctx)
		default:
			panic("Invalid RPC Name")
		}
	default:
		panic("Invalid Module Name")
	}

}

func systemTime(sc gnoi_system_pb.SystemClient, ctx context.Context) {
	fmt.Println("System Time")
	ctx = setUserCreds(ctx)
	resp,err := sc.Time(ctx, new(gnoi_system_pb.TimeRequest))
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp.Time)
}
func sonicShowTechSupport(sc spb.SonicServiceClient, ctx context.Context) {
	fmt.Println("Sonic ShowTechsupport")
	ctx = setUserCreds(ctx)
	req := &spb.TechsupportRequest {
		Input: &spb.TechsupportRequest_Input{

		},
	}
	nargs := strings.Replace(string(*args), "sonic-show-techsupport:input", "input", 1)
	json.Unmarshal([]byte(nargs), &req)
	resp,err := sc.ShowTechsupport(ctx, req)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(resp)
}
func sonicSum(sc spb.SonicServiceClient, ctx context.Context) {
	fmt.Println("Sonic Sum")
	ctx = setUserCreds(ctx)
	req := &spb.SumRequest{
		Input: &spb.SumRequest_Input{},
	}
	nargs := strings.Replace(string(*args), "sonic-tests:input", "input", 1)
	json.Unmarshal([]byte(nargs), &req)

	resp,err := sc.Sum(ctx, req)

	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp.Output.Result)
}

func saveConfig(sc spb.SonicServiceClient, ctx context.Context) {
	fmt.Println("Sonic SaveConfig")
	ctx = setUserCreds(ctx)
	req := &spb.SaveConfigRequest{
		Input: &spb.SaveConfigRequest_Input{},
	}
	nargs := strings.Replace(string(*args), "sonic-config-mgmt:input", "input", 1)
	json.Unmarshal([]byte(nargs), &req)

	resp,err := sc.SaveConfig(ctx, req)

	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp.Output.Status)
}

func reloadConfig(sc spb.SonicServiceClient, ctx context.Context) {
	fmt.Println("Sonic ReloadConfig")
	ctx = setUserCreds(ctx)
	req := &spb.ReloadConfigRequest{
		Input: &spb.ReloadConfigRequest_Input{},
	}
	nargs := strings.Replace(string(*args), "sonic-config-mgmt:input", "input", 1)
	json.Unmarshal([]byte(nargs), &req)

	resp,err := sc.ReloadConfig(ctx, req)

	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp.Output.Status)
}

func loadMgmtConfig(sc spb.SonicServiceClient, ctx context.Context) {
	fmt.Println("Sonic LoadMgmtConfig")
	ctx = setUserCreds(ctx)
	req := &spb.LoadMgmtConfigRequest{
		Input: &spb.LoadMgmtConfigRequest_Input{},
	}
	nargs := strings.Replace(string(*args), "sonic-config-mgmt:input", "input", 1)
	json.Unmarshal([]byte(nargs), &req)

	resp,err := sc.LoadMgmtConfig(ctx, req)

	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp.Output.Status)
}

func loadMinigraph(sc spb.SonicServiceClient, ctx context.Context) {
	fmt.Println("Sonic LoadMinigraph")
	ctx = setUserCreds(ctx)
	req := &spb.LoadMinigraphRequest{
		Input: &spb.LoadMinigraphRequest_Input{},
	}
	nargs := strings.Replace(string(*args), "sonic-config-mgmt:input", "input", 1)
	json.Unmarshal([]byte(nargs), &req)

	resp,err := sc.LoadMinigraph(ctx, req)

	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp.Output.Status)
}

func authenticate(sc spb.SonicServiceClient, ctx context.Context) {
	fmt.Println("Sonic Authenticate")
	ctx = setUserCreds(ctx)
	req := &spb.AuthenticateRequest {}
	
	json.Unmarshal([]byte(*args), &req)
	fmt.Println(req)
	resp,err := sc.Authenticate(ctx, req)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp)
}

func refresh(sc spb.SonicServiceClient, ctx context.Context) {
	fmt.Println("Sonic Refresh")
	ctx = setUserCreds(ctx)
	req := &spb.RefreshRequest {}
	
	json.Unmarshal([]byte(*args), &req)
	fmt.Println(req)
	resp,err := sc.Refresh(ctx, req)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(resp)
}