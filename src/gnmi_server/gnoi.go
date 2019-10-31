package gnmi_server

import (
	"context"
	gnoi_system_pb "github.com/openconfig/gnoi/system"
	log "github.com/golang/glog"
	"time"
	spb "proto/gnoi"
	transutil "transl_utils"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

func (srv *Server) Reboot(context.Context, *gnoi_system_pb.RebootRequest) (*gnoi_system_pb.RebootResponse, error) {
	log.V(1).Info("gNOI: Reboot")
	return nil, nil
}
func (srv *Server) RebootStatus(context.Context, *gnoi_system_pb.RebootStatusRequest) (*gnoi_system_pb.RebootStatusResponse, error) {
	log.V(1).Info("gNOI: RebootStatus")
	return nil, nil
}
func (srv *Server) CancelReboot(context.Context, *gnoi_system_pb.CancelRebootRequest) (*gnoi_system_pb.CancelRebootResponse, error) {
	log.V(1).Info("gNOI: CancelReboot")
	return nil, nil
}
func (srv *Server) Ping(*gnoi_system_pb.PingRequest, gnoi_system_pb.System_PingServer) error {
	log.V(1).Info("gNOI: Ping")
	return nil
}
func (srv *Server) Traceroute(*gnoi_system_pb.TracerouteRequest, gnoi_system_pb.System_TracerouteServer) error {
	log.V(1).Info("gNOI: Traceroute")
	return nil
}
func (srv *Server) SetPackage(gnoi_system_pb.System_SetPackageServer) error {
	log.V(1).Info("gNOI: SetPackage")
	return nil
}
func (srv *Server) SwitchControlProcessor(context.Context, *gnoi_system_pb.SwitchControlProcessorRequest) (*gnoi_system_pb.SwitchControlProcessorResponse, error) {
	log.V(1).Info("gNOI: SwitchControlProcessor")
	return nil, nil
}
func (srv *Server) Time(context.Context, *gnoi_system_pb.TimeRequest) (*gnoi_system_pb.TimeResponse, error) {
	log.V(1).Info("gNOI: Time")
	var tm gnoi_system_pb.TimeResponse
	tm.Time = uint64(time.Now().UnixNano())
	return &tm, nil
}


func (srv *Server) ShowTechsupport(context.Context, *spb.TechsupportRequest) (*spb.TechsupportResponse, error) {
	log.V(1).Info("gNOI: Sonic ShowTechsupport")
	var resp spb.TechsupportResponse
	// resp.OutputFilename = "test"
	jsresp, err:= transutil.TranslProcessAction("sonic-show-techsupport:sonic-show-techsupport-info", []byte("{\"sonic-show-techsupport-info:input\": {\"date\": \"2019-10-31T09:00:00Z\"}"))
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	resp.OutputFilename = string(jsresp)
	return &resp, nil
}

func (srv *Server) MyEcho(context.Context, *spb.MyEchoRequest) (*spb.MyEchoResponse, error) {
	log.V(1).Info("gNOI: Sonic MyEcho")
	var resp spb.MyEchoResponse
	// resp.OutputFilename = "test"
	jsresp, err:= transutil.TranslProcessAction("/restconf/operations/api-tests:my-echo", []byte("{\"api-tests:input\": {\"message\": \"2019-10-31T09:00:00Z\"}"))
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	resp.Message = string(jsresp)
	return &resp, nil
}